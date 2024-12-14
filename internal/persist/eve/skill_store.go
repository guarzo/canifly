package eve

import (
	"bufio"
	"encoding/json"
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
	s.logger.Infof("load skill plans")

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

	if err := s.fs.WriteFile(planFilePath, []byte(sb.String()), 0644); err != nil {
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
		return nil, fmt.Errorf("failed to read eve plans directory: %w", err)
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
		return nil, fmt.Errorf("failed to read eve plan file %s: %w", filePath, err)
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
			return nil, fmt.Errorf("invalid eve level in %s at line %d: %s", filePath, lineNumber, skillLevelStr)
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
	s.logger.Infof("load skill types")
	file, err := embed.StaticFiles.Open(skillTypeFile)
	if err != nil {
		return fmt.Errorf("failed to open eve type file %s: %w", skillTypeFile, err)
	}
	defer file.Close()

	records, err := persist.ReadCsvRecords(file)
	if err != nil {
		return fmt.Errorf("failed to read CSV records from %s: %w", skillTypeFile, err)
	}

	skillTypes, skillIDTypes, err := s.parseSkillTypes(records)
	if err != nil {
		return fmt.Errorf("failed to parse eve types: %w", err)
	}

	s.mut.Lock()
	s.skillTypes = skillTypes
	s.skillIdToType = skillIDTypes
	s.mut.Unlock()

	s.logger.Debugf("Loaded %d eve types", len(skillTypes))
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

// A helper function to load deleted embedded plans from a JSON or text file
func (s *SkillStore) loadDeletedEmbeddedPlans() (map[string]bool, error) {
	deletedPlans := make(map[string]bool)
	deletedListPath := filepath.Join(s.basePath, plansDir, "deleted_embedded_plans.json")

	if _, err := os.Stat(deletedListPath); os.IsNotExist(err) {
		// No file, so no deleted plans
		return deletedPlans, nil
	}

	data, err := s.fs.ReadFile(deletedListPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read deleted embedded plans: %w", err)
	}

	// Suppose we store them as a simple JSON array of plan names
	var planNames []string
	if err := json.Unmarshal(data, &planNames); err != nil {
		return nil, fmt.Errorf("failed to parse deleted embedded plans: %w", err)
	}
	for _, name := range planNames {
		deletedPlans[name] = true
	}
	return deletedPlans, nil
}

func (s *SkillStore) saveDeletedEmbeddedPlans(deletedPlans map[string]bool) error {
	planNames := make([]string, 0, len(deletedPlans))
	for name := range deletedPlans {
		planNames = append(planNames, name)
	}
	data, err := json.MarshalIndent(planNames, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal deleted embedded plans: %w", err)
	}

	deletedListPath := filepath.Join(s.basePath, plansDir, "deleted_embedded_plans.json")
	if err := s.fs.WriteFile(deletedListPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write deleted embedded plans file: %w", err)
	}
	return nil
}

// Enhance copyEmbeddedPlansToWritable to skip deleted plans
func (s *SkillStore) copyEmbeddedPlansToWritable(writableDir string) error {
	deletedPlans, err := s.loadDeletedEmbeddedPlans()
	if err != nil {
		return err
	}

	entries, err := embed.StaticFiles.ReadDir("static/plans")
	if err != nil {
		return fmt.Errorf("failed to read embedded plans: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			planName := strings.TrimSuffix(fileName, ".txt")

			// If the user previously deleted this embedded plan, skip copying it
			if deletedPlans[planName] {
				s.logger.Debugf("Skipping previously deleted embedded plan: %s", planName)
				continue
			}

			destPath := filepath.Join(writableDir, fileName)
			// Only copy if file does not exist to avoid overwriting custom user changes
			if _, err := s.fs.Stat(destPath); os.IsNotExist(err) {
				if err := s.copyEmbeddedFile("static/plans/"+fileName, destPath); err != nil {
					return fmt.Errorf("failed to copy embedded plan %s: %w", fileName, err)
				}
			}
		}
	}
	return nil
}

// Modify DeleteSkillPlan to record deletions of embedded plans
func (s *SkillStore) DeleteSkillPlan(planName string) error {
	planFilePath := filepath.Join(s.basePath, plansDir, planName+".txt")

	if err := s.fs.Remove(planFilePath); err != nil {
		if os.IsNotExist(err) {
			s.logger.Warnf("Skill plan %s does not exist", planName)
			return fmt.Errorf("eve plan does not exist: %w", err)
		}
		return fmt.Errorf("failed to delete eve plan file: %w", err)
	}

	// If this plan was embedded originally, add it to the deleted list
	// One way is to check if it matches an embedded plan name
	// For a robust solution, store the original embedded plan list somewhere.
	// Here, we read from embed and check if planName was one of them.
	entries, err := embed.StaticFiles.ReadDir("static/plans")
	if err == nil {
		embeddedPlan := false
		for _, entry := range entries {
			if !entry.IsDir() && strings.TrimSuffix(entry.Name(), ".txt") == planName {
				embeddedPlan = true
				break
			}
		}
		if embeddedPlan {
			deletedPlans, err := s.loadDeletedEmbeddedPlans()
			if err != nil {
				return err
			}
			deletedPlans[planName] = true
			if err := s.saveDeletedEmbeddedPlans(deletedPlans); err != nil {
				return err
			}
		}
	}

	s.mut.Lock()
	delete(s.skillPlans, planName)
	s.mut.Unlock()

	s.logger.Infof("Deleted eve plan %s", planName)
	return nil
}
