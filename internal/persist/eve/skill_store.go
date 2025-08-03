package eve

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"github.com/guarzo/canifly/internal/services/skillplans"
)

var _ interfaces.SkillRepository = (*SkillStore)(nil)

const (
	plansDir = "plans"
)

// SkillStore implements interfaces.SkillRepository
type SkillStore struct {
	logger           interfaces.Logger
	fs               persist.FileSystem
	basePath         string
	skillPlans       map[string]model.SkillPlan
	skillTypes       map[string]model.SkillType
	skillIdToType    map[string]model.SkillType
	githubDownloader *skillplans.GitHubDownloader
	mut              sync.RWMutex
}

// NewSkillStore now accepts a FileSystem and a basePath for writable directories.
// The githubDownloader parameter is optional - pass it to enable downloading skill plans from GitHub.
func NewSkillStore(logger interfaces.Logger, fs persist.FileSystem, basePath string, opts ...func(*SkillStore)) *SkillStore {
	store := &SkillStore{
		logger:        logger,
		fs:            fs,
		basePath:      basePath,
		skillPlans:    make(map[string]model.SkillPlan),
		skillTypes:    make(map[string]model.SkillType),
		skillIdToType: make(map[string]model.SkillType),
	}

	// Apply options
	for _, opt := range opts {
		opt(store)
	}

	return store
}

// WithGitHubDownloader is an option to set the GitHub downloader
func WithGitHubDownloader(downloader *skillplans.GitHubDownloader) func(*SkillStore) {
	return func(s *SkillStore) {
		s.githubDownloader = downloader
	}
}

func (s *SkillStore) LoadSkillPlans() error {
	s.logger.Infof("load skill plans")

	writableDir := filepath.Join(s.basePath, plansDir)
	if err := s.fs.MkdirAll(writableDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to ensure plans directory: %w", err)
	}

	// Try to download latest plans from GitHub
	if s.githubDownloader != nil {
		if err := s.githubDownloader.DownloadPlans(writableDir); err != nil {
			s.logger.Warnf("Failed to download plans from GitHub: %v", err)
			// Continue with whatever local files exist
		}
	}

	plans, err := s.loadSkillPlans(writableDir)
	if err != nil {
		return fmt.Errorf("failed to load eve plans: %w", err)
	}

	s.mut.Lock()
	s.skillPlans = plans
	s.mut.Unlock()

	s.logger.Debugf("Loaded %d eve plans", len(plans))
	return nil
}

func (s *SkillStore) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	if len(skills) == 0 {
		return fmt.Errorf("cannot save an empty eve plan for planName: %s", planName)
	}

	planFilePath := filepath.Join(s.basePath, plansDir, planName+".txt")

	var sb strings.Builder
	for skillName, skill := range skills {
		sb.WriteString(fmt.Sprintf("%s %d\n", skillName, skill.Level))
	}

	if err := persist.AtomicWriteFile(s.fs, planFilePath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write plan file: %w", err)
	}

	planKey := planName
	s.mut.Lock()
	s.skillPlans[planKey] = model.SkillPlan{Name: planKey, Skills: skills}
	s.mut.Unlock()
	s.logger.Infof("Saved eve plan %s with %d skills", planKey, len(skills))
	return nil
}

func (s *SkillStore) GetSkillPlans() map[string]model.SkillPlan {
	s.mut.RLock()
	defer s.mut.RUnlock()

	// Return a copy if needed, or the original if safe
	plansCopy := make(map[string]model.SkillPlan, len(s.skillPlans))
	for k, v := range s.skillPlans {
		plansCopy[k] = v
	}
	return plansCopy
}

func (s *SkillStore) GetSkillPlanFile(planName string) ([]byte, error) {
	planName += ".txt"
	s.logger.Infof("Attempting to serve eve plan file: %s", planName)

	skillPlanDir := filepath.Join(s.basePath, plansDir)

	filePath := filepath.Join(skillPlanDir, planName)
	return os.ReadFile(filePath)
}

func (s *SkillStore) loadSkillPlans(dir string) (map[string]model.SkillPlan, error) {
	plans := make(map[string]model.SkillPlan)

	// We need to list files in dir. Since we're using fs abstraction for reading,
	// we might still rely on os.ReadDir if fs does not provide a listing method.
	// If needed, extend FileSystem or handle that logic outside.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read eve plans directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		planName := strings.TrimSuffix(entry.Name(), ".txt")
		path := filepath.Join(dir, entry.Name())

		parsed, err := s.readSkillsFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read skills from %s: %w", path, err)
		}
		plans[planName] = model.SkillPlan{
			Name:   planName,
			Skills: parsed.Skills,
			Icon:   parsed.Icon,
		}
	}

	return plans, nil
}

// ParsedSkillPlan holds the result of parsing a skill plan file
type ParsedSkillPlan struct {
	Skills map[string]model.Skill
	Icon   string // Optional icon identifier
}

func (s *SkillStore) readSkillsFromFile(filePath string) (*ParsedSkillPlan, error) {
	data, err := s.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read eve plan file %s: %w", filePath, err)
	}

	result := &ParsedSkillPlan{
		Skills: make(map[string]model.Skill),
	}

	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Check for icon directive
		if strings.HasPrefix(trimmed, "# icon:") {
			result.Icon = strings.TrimSpace(strings.TrimPrefix(trimmed, "# icon:"))
			continue
		}

		// Skip empty lines and other comments
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format in file %s at line %d: %s", filePath, lineNumber, line)
		}

		skillLevelStr := parts[len(parts)-1]
		skillName := strings.Join(parts[:len(parts)-1], " ")
		skillLevel, err := strconv.Atoi(skillLevelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid eve level in %s at line %d: %s", filePath, lineNumber, skillLevelStr)
		}

		if currentSkill, exists := result.Skills[skillName]; !exists || skillLevel > currentSkill.Level {
			result.Skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", filePath, err)
	}

	s.logger.Debugf("Read %d skills from %s (icon: %s)", len(result.Skills), filePath, result.Icon)
	return result, nil
}

func (s *SkillStore) LoadSkillTypes() error {
	s.logger.Infof("load skill types")

	// Load from downloaded Fuzzworks data only
	fuzzworksPath := filepath.Join(s.basePath, "config", "fuzzworks", "invTypes.csv")
	file, err := s.fs.Open(fuzzworksPath)
	if err != nil {
		return fmt.Errorf("skill types not found - ensure Fuzzworks data is downloaded: %w", err)
	}
	defer file.Close()

	s.logger.Infof("Loading skill types from Fuzzworks data: %s", fuzzworksPath)

	records, err := persist.ReadCsvRecords(file)
	if err != nil {
		return fmt.Errorf("failed to read Fuzzworks CSV: %w", err)
	}

	skillTypes, skillIDTypes, err := s.parseSkillTypes(records)
	if err != nil {
		return fmt.Errorf("failed to parse Fuzzworks data: %w", err)
	}

	s.mut.Lock()
	s.skillTypes = skillTypes
	s.skillIdToType = skillIDTypes
	s.mut.Unlock()

	s.logger.Debugf("Loaded %d skill types from Fuzzworks data", len(skillTypes))
	return nil
}

func (s *SkillStore) parseSkillTypes(records [][]string) (map[string]model.SkillType, map[string]model.SkillType, error) {
	if len(records) == 0 {
		return nil, nil, fmt.Errorf("no data in eve type file")
	}

	headers := records[0]
	records = records[1:] // skip header

	colIndices := map[string]int{"typeID": -1, "typeName": -1, "description": -1}
	for i, header := range headers {
		switch strings.TrimSpace(header) {
		case "typeID":
			colIndices["typeID"] = i
		case "typeName":
			colIndices["typeName"] = i
		case "description":
			colIndices["description"] = i
		}
	}

	if colIndices["typeID"] == -1 || colIndices["typeName"] == -1 {
		return nil, nil, fmt.Errorf("required columns (typeID, typeName) are missing")
	}

	skillTypes := make(map[string]model.SkillType)
	skillIDTypes := make(map[string]model.SkillType)

	lineNumber := 1
	for _, row := range records {
		lineNumber++
		if len(row) < 2 {
			s.logger.Warnf("Skipping malformed row %d in eve types", lineNumber)
			continue
		}

		typeID := strings.TrimSpace(row[colIndices["typeID"]])
		typeName := strings.TrimSpace(row[colIndices["typeName"]])

		if typeName == "" {
			continue
		}

		desc := ""
		if di := colIndices["description"]; di != -1 && di < len(row) {
			desc = strings.TrimSpace(row[di])
		}

		st := model.SkillType{
			TypeID:      typeID,
			TypeName:    typeName,
			Description: desc,
		}

		skillTypes[typeName] = st
		skillIDTypes[typeID] = st
	}

	return skillTypes, skillIDTypes, nil
}

func (s *SkillStore) GetSkillTypes() map[string]model.SkillType {
	s.mut.RLock()
	defer s.mut.RUnlock()
	// return a copy if needed
	cpy := make(map[string]model.SkillType, len(s.skillTypes))
	for k, v := range s.skillTypes {
		cpy[k] = v
	}
	return cpy
}

func (s *SkillStore) GetSkillTypeByID(id string) (model.SkillType, bool) {
	s.mut.RLock()
	defer s.mut.RUnlock()
	st, ok := s.skillIdToType[id]
	return st, ok
}

func (s *SkillStore) DeleteSkillPlan(planName string) error {
	planFilePath := filepath.Join(s.basePath, plansDir, planName+".txt")

	if err := s.fs.Remove(planFilePath); err != nil {
		if os.IsNotExist(err) {
			s.logger.Warnf("Skill plan %s does not exist", planName)
			return fmt.Errorf("eve plan does not exist: %w", err)
		}
		return fmt.Errorf("failed to delete eve plan file: %w", err)
	}

	s.mut.Lock()
	delete(s.skillPlans, planName)
	s.mut.Unlock()

	s.logger.Infof("Deleted eve plan %s", planName)
	return nil
}
