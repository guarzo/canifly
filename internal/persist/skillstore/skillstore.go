package skillstore

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.SkillRepository = (*SkillStore)(nil)

const (
	plansDir      = "plans"
	skillTypeFile = "static/invTypes.csv"
)

// SkillStore implements interfaces.SkillRepository
type SkillStore struct {
	logger        interfaces.Logger
	fs            persist.FileSystem
	basePath      string
	skillPlans    map[string]model.SkillPlan
	skillTypes    map[string]model.SkillType
	skillIdToType map[string]model.SkillType
	mut           sync.RWMutex
}

// NewSkillStore now accepts a FileSystem and a basePath for writable directories.
func NewSkillStore(logger interfaces.Logger, fs persist.FileSystem, basePath string) *SkillStore {
	return &SkillStore{
		logger:     logger,
		fs:         fs,
		basePath:   basePath,
		skillPlans: make(map[string]model.SkillPlan),
		skillTypes: make(map[string]model.SkillType),
	}
}

func (s *SkillStore) LoadSkillPlans() error {
	writableDir := filepath.Join(s.basePath, plansDir)
	if err := s.fs.MkdirAll(writableDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to ensure plans directory: %w", err)
	}

	// Copy embedded plans if needed
	if err := s.copyEmbeddedPlansToWritable(writableDir); err != nil {
		return fmt.Errorf("failed to copy embedded plans: %w", err)
	}

	plans, err := s.loadSkillPlans(writableDir)
	if err != nil {
		return fmt.Errorf("failed to load skill plans: %w", err)
	}

	s.mut.Lock()
	s.skillPlans = plans
	s.mut.Unlock()

	s.logger.Debugf("Loaded %d skill plans", len(plans))
	return nil
}

func (s *SkillStore) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	if len(skills) == 0 {
		return fmt.Errorf("cannot save an empty skill plan for planName: %s", planName)
	}

	planFilePath := filepath.Join(s.basePath, plansDir, planName+".txt")

	var sb strings.Builder
	for skillName, skill := range skills {
		sb.WriteString(fmt.Sprintf("%s %d\n", skillName, skill.Level))
	}

	if err := s.fs.WriteFile(planFilePath, []byte(sb.String()), 0644); err != nil {
		return fmt.Errorf("failed to write plan file: %w", err)
	}

	planKey := planName
	s.mut.Lock()
	s.skillPlans[planKey] = model.SkillPlan{Name: planKey, Skills: skills}
	s.mut.Unlock()
	s.logger.Infof("Saved skill plan %s with %d skills", planKey, len(skills))
	return nil
}

func (s *SkillStore) DeleteSkillPlan(planName string) error {
	planFilePath := filepath.Join(s.basePath, plansDir, planName+".txt")

	if err := s.fs.Remove(planFilePath); err != nil {
		if os.IsNotExist(err) {
			s.logger.Warnf("Skill plan %s does not exist", planName)
			return fmt.Errorf("skill plan does not exist: %w", err)
		}
		return fmt.Errorf("failed to delete skill plan file: %w", err)
	}

	s.mut.Lock()
	delete(s.skillPlans, planName)
	s.mut.Unlock()

	s.logger.Infof("Deleted skill plan %s", planName)
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
	s.logger.Infof("Attempting to serve skill plan file: %s", planName)

	skillPlanDir := filepath.Join(s.basePath, plansDir)

	filePath := filepath.Join(skillPlanDir, planName)
	return os.ReadFile(filePath)
}

func (s *SkillStore) copyEmbeddedPlansToWritable(writableDir string) error {
	entries, err := embed.StaticFiles.ReadDir("static/plans")
	if err != nil {
		return fmt.Errorf("failed to read embedded plans: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			destPath := filepath.Join(writableDir, fileName)

			if err := s.copyEmbeddedFile("static/plans/"+fileName, destPath); err != nil {
				return fmt.Errorf("failed to copy embedded plan %s: %w", fileName, err)
			}
		}
	}
	return nil
}

func (s *SkillStore) copyEmbeddedFile(srcPath, destPath string) error {
	srcFile, err := embed.StaticFiles.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	data, err := io.ReadAll(srcFile)
	if err != nil {
		return fmt.Errorf("failed to read embedded file %s: %w", srcPath, err)
	}

	dir := filepath.Dir(destPath)
	if err := s.fs.MkdirAll(dir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory for %s: %w", destPath, err)
	}

	if err := s.fs.WriteFile(destPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %w", destPath, err)
	}
	return nil
}

func (s *SkillStore) loadSkillPlans(dir string) (map[string]model.SkillPlan, error) {
	plans := make(map[string]model.SkillPlan)

	// We need to list files in dir. Since we're using fs abstraction for reading,
	// we might still rely on os.ReadDir if fs does not provide a listing method.
	// If needed, extend FileSystem or handle that logic outside.
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill plans directory: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".txt") {
			continue
		}

		planName := strings.TrimSuffix(entry.Name(), ".txt")
		path := filepath.Join(dir, entry.Name())

		skills, err := s.readSkillsFromFile(path)
		if err != nil {
			return nil, fmt.Errorf("failed to read skills from %s: %w", path, err)
		}
		plans[planName] = model.SkillPlan{Name: planName, Skills: skills}
	}

	return plans, nil
}

func (s *SkillStore) readSkillsFromFile(filePath string) (map[string]model.Skill, error) {
	data, err := s.fs.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read skill plan file %s: %w", filePath, err)
	}

	skills := make(map[string]model.Skill)
	scanner := bufio.NewScanner(strings.NewReader(string(data)))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if strings.TrimSpace(line) == "" {
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
			return nil, fmt.Errorf("invalid skill level in %s at line %d: %s", filePath, lineNumber, skillLevelStr)
		}

		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error scanning file %s: %w", filePath, err)
	}

	s.logger.Debugf("Read %d skills from %s", len(skills), filePath)
	return skills, nil
}

func (s *SkillStore) LoadSkillTypes() error {
	file, err := embed.StaticFiles.Open(skillTypeFile)
	if err != nil {
		return fmt.Errorf("failed to open skill type file %s: %w", skillTypeFile, err)
	}
	defer file.Close()

	records, err := persist.ReadCsvRecords(file)
	if err != nil {
		return fmt.Errorf("failed to read CSV records from %s: %w", skillTypeFile, err)
	}

	skillTypes, skillIDTypes, err := s.parseSkillTypes(records)
	if err != nil {
		return fmt.Errorf("failed to parse skill types: %w", err)
	}

	s.mut.Lock()
	s.skillTypes = skillTypes
	s.skillIdToType = skillIDTypes
	s.mut.Unlock()

	s.logger.Debugf("Loaded %d skill types", len(skillTypes))
	return nil
}

func (s *SkillStore) parseSkillTypes(records [][]string) (map[string]model.SkillType, map[string]model.SkillType, error) {
	if len(records) == 0 {
		return nil, nil, fmt.Errorf("no data in skill type file")
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
			s.logger.Warnf("Skipping malformed row %d in skill types", lineNumber)
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