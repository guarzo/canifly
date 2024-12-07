// persist/skillplan.go
package persist

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/guarzo/canifly/internal/embed"
	"github.com/guarzo/canifly/internal/model"
)

func (ds *DataStore) ProcessSkillPlans() error {
	writablePath, err := ds.GetWriteablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable plans path: %w", err)
	}

	if err := ds.copyEmbeddedPlansToWritable(writablePath); err != nil {
		return fmt.Errorf("failed to copy embedded plans: %w", err)
	}

	ds.logger.Infof("Loading skill plans from %s", writablePath)
	skillPlans, err := ds.loadSkillPlans(writablePath)
	if err != nil {
		return fmt.Errorf("failed to load skill plans: %w", err)
	}
	ds.skillPlans = skillPlans
	ds.logger.Infof("Loaded %d skill plans", len(ds.skillPlans))
	return nil
}

func (ds *DataStore) GetSkillPlans() map[string]model.SkillPlan {
	return ds.skillPlans
}

func (ds *DataStore) SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	if len(skills) == 0 {
		return fmt.Errorf("cannot save an empty skill plan for planName: %s", planName)
	}

	writablePath, err := ds.GetWriteablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable path: %w", err)
	}

	planFilePath := filepath.Join(writablePath, planName+".txt")
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
	ds.skillPlans[planKey] = model.SkillPlan{Name: planKey, Skills: skills}
	ds.logger.Infof("Saved skill plan %s with %d skills", planKey, len(skills))
	return nil
}

func (ds *DataStore) DeleteSkillPlan(planName string) error {
	writablePath, err := ds.GetWriteablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable path: %w", err)
	}

	planFilePath := filepath.Join(writablePath, planName+".txt")
	if err := os.Remove(planFilePath); err != nil {
		if os.IsNotExist(err) {
			ds.logger.Warnf("Skill plan %s does not exist", planName)
			return fmt.Errorf("skill plan does not exist: %w", err)
		}
		ds.logger.WithError(err).Errorf("Failed to delete skill plan file %s", planFilePath)
		return err
	}

	delete(ds.skillPlans, planName)
	ds.logger.Infof("Deleted skill plan %s", planName)
	return nil
}

// Convert helper functions to methods on DataStore as needed

func (ds *DataStore) GetWriteablePlansPath() (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	planPath := filepath.Join(configDir, "canifly", "plans")
	if pathSuffix != "" {
		planPath = filepath.Join(planPath, pathSuffix)
	}

	if err := os.MkdirAll(planPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create writable plans directory: %w", err)
	}

	return planPath, nil
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

	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy content to %s: %w", destPath, err)
	}

	ds.logger.Infof("Copied embedded file %s to %s", srcPath, destPath)
	return nil
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
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	skills := make(map[string]model.Skill)
	scanner := bufio.NewScanner(file)
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
			return nil, fmt.Errorf("invalid skill level in file %s at line %d: %s", filePath, lineNumber, skillLevelStr)
		}

		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}

	if err = scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	ds.logger.Infof("Read %d skills from %s", len(skills), filePath)
	return skills, nil
}
