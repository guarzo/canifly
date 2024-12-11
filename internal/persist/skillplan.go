// persist/skillplan.go
package persist

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

func (ds *DataStore) ProcessSkillPlans() error {
	writablePath, err := getWriteablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable plans path: %w", err)
	}

	if err := ds.copyEmbeddedPlansToWritable(writablePath); err != nil {
		return fmt.Errorf("failed to copy embedded plans: %w", err)
	}

	ds.logger.Debugf("Loading skill plans from %s", writablePath)
	skillPlans, err := ds.loadSkillPlans(writablePath)
	if err != nil {
		return fmt.Errorf("failed to load skill plans: %w", err)
	}
	ds.mut.Lock()
	ds.skillPlans = skillPlans
	ds.mut.Unlock()
	ds.logger.Debugf("Loaded %d skill plans", len(ds.skillPlans))
	return nil
}

func (ds *DataStore) GetWriteablePlansPath() (string, error) {
	return getWriteablePlansPath()
}

func (ds *DataStore) GetSkillPlans() map[string]model.SkillPlan {
	ds.mut.RLock()
	defer ds.mut.RUnlock()
	return ds.skillPlans
}

func (ds *DataStore) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	if len(skills) == 0 {
		return fmt.Errorf("cannot save an empty skill plan for planName: %s", planName)
	}

	planFilePath, err := getPlanFileName(planName + ".txt")
	file, err := os.Create(planFilePath)
	if err != nil {
		ds.logger.WithError(err).Errorf("Failed to create plan file %s", planFilePath)
		return err
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for skillName, skill := range skills {
		line := fmt.Sprintf("%s %d\n", skillName, skill.Level)
		if _, err := writer.WriteString(line); err != nil {
			ds.logger.WithError(err).Errorf("Failed to write skill to file %s", planFilePath)
			return err
		}
	}
	if err := writer.Flush(); err != nil {
		ds.logger.WithError(err).Errorf("Failed to flush writer for %s", planFilePath)
		return err
	}

	planKey := strings.TrimSuffix(planName, ".txt")
	ds.mut.Lock()
	ds.skillPlans[planKey] = model.SkillPlan{Name: planKey, Skills: skills}
	ds.mut.Unlock()
	ds.logger.Infof("Saved skill plan %s with %d skills", planKey, len(skills))
	return nil
}

func (ds *DataStore) DeleteSkillPlan(planName string) error {
	planFilePath, err := getPlanFileName(planName + ".txt")
	if err != nil {
		return err
	}
	if err := os.Remove(planFilePath); err != nil {
		if os.IsNotExist(err) {
			ds.logger.Warnf("Skill plan %s does not exist", planName)
			return fmt.Errorf("skill plan does not exist: %w", err)
		}
		ds.logger.WithError(err).Errorf("Failed to delete skill plan file %s", planFilePath)
		return err
	}

	ds.mut.Lock()
	delete(ds.skillPlans, planName)
	ds.mut.Unlock()

	ds.logger.Infof("Deleted skill plan %s", planName)
	return nil
}

func (ds *DataStore) copyEmbeddedPlansToWritable(writableDir string) error {
	entries, err := embed.StaticFiles.ReadDir("static/plans")
	if err != nil {
		return fmt.Errorf("failed to read embedded plans: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			destPath := filepath.Join(writableDir, fileName)

			if err := ds.copyEmbeddedFile("static/plans/"+fileName, destPath); err != nil {
				return fmt.Errorf("failed to copy embedded plan %s: %w", fileName, err)
			}

			if _, statErr := os.Stat(destPath); statErr != nil {
				ds.logger.WithError(statErr).Errorf("File %s was copied but does not exist at %s", fileName, destPath)
				return fmt.Errorf("file %s does not exist after copying", fileName)
			}
		}
	}
	return nil
}

func (ds *DataStore) copyEmbeddedFile(srcPath, destPath string) error {
	srcFile, err := embed.StaticFiles.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	return copyReaderToFile(destPath, srcFile)
}

func (ds *DataStore) loadSkillPlans(dir string) (map[string]model.SkillPlan, error) {
	skillPlans := make(map[string]model.SkillPlan)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			ds.logger.WithError(err).Errorf("Error walking through %s", path)
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".txt") {
			planName := strings.TrimSuffix(filepath.Base(path), ".txt")

			skills, err := ds.readSkillsFromFile(path)
			if err != nil {
				ds.logger.WithError(err).Errorf("Failed to read skills from %s", path)
				return err
			}
			skillPlans[planName] = model.SkillPlan{Name: planName, Skills: skills}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load skill plans from %s: %w", dir, err)
	}

	return skillPlans, nil
}

func (ds *DataStore) readSkillsFromFile(filePath string) (map[string]model.Skill, error) {
	skills := make(map[string]model.Skill)
	lineNumber := 0

	err := processFileLines(filePath, func(line string) error {
		lineNumber++
		if strings.TrimSpace(line) == "" {
			return nil
		}

		parts := strings.Fields(line)
		if len(parts) < 2 {
			return fmt.Errorf("invalid format in file %s at line %d: %s", filePath, lineNumber, line)
		}

		skillLevelStr := parts[len(parts)-1]
		skillName := strings.Join(parts[:len(parts)-1], " ")

		skillLevel, err := strconv.Atoi(skillLevelStr)
		if err != nil {
			return fmt.Errorf("invalid skill level in file %s at line %d: %s", filePath, lineNumber, skillLevelStr)
		}

		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	ds.logger.Debugf("Read %d skills from %s", len(skills), filePath)
	return skills, nil
}
