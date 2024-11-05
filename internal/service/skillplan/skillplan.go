package skillplan

import (
	"bufio"
	"fmt"
	"github.com/gambtho/canifly/internal/utils/xlog"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gambtho/canifly/internal/embed"
	"github.com/gambtho/canifly/internal/model"
)

var SkillPlans map[string]model.SkillPlan // Holds all loaded skill plans

// GetWritablePlansPath returns the writable directory path for storing skill plans.
func GetWritablePlansPath() (string, error) {
	configDir, err := os.UserConfigDir()

	if err != nil {
		return "", fmt.Errorf("failed to retrieve writeable directory: %w", err)
	}

	pathSuffix := os.Getenv("PATH_SUFFIX")
	planPath := filepath.Join(configDir, "canifly", "plans")
	if pathSuffix != "" {
		planPath = filepath.Join(configDir, "canifly", "plans", pathSuffix)
	}

	if err := os.MkdirAll(planPath, os.ModePerm); err != nil {
		return "", fmt.Errorf("failed to create writable plans directory: %w", err)
	}

	return planPath, nil
}

// ProcessSkillPlans loads all skill plans, copying static plans from embed to the writable path if needed.
func ProcessSkillPlans() error {
	// Get paths for writable and static directories
	writablePath, err := GetWritablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable plans path: %w", err)
	}

	// Copy embedded skill plans to the writable directory if needed
	if err := CopyEmbeddedPlansToWritable(writablePath); err != nil {
		return fmt.Errorf("failed to copy embedded plans: %w", err)
	}

	// Verify that all files in the writable directory are accessible
	entries, _ := os.ReadDir(writablePath)
	for _, entry := range entries {
		if !entry.IsDir() {
			_, err := os.Open(filepath.Join(writablePath, entry.Name()))
			if err != nil {
				xlog.Logf("Access check failed for file: %s - Error: %v", entry.Name(), err)
			}
		}
	}

	// Load all skill plans from the writable directory into the SkillPlans map
	SkillPlans, err = loadSkillPlans(writablePath)
	if err != nil {
		return fmt.Errorf("failed to load skill plans from %s: %w", writablePath, err)
	}

	return nil
}

// loadSkillPlans reads all skill files in a directory and returns a map of skill plans.
func loadSkillPlans(dir string) (map[string]model.SkillPlan, error) {
	skillPlans := make(map[string]model.SkillPlan)

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Process only files, not directories
		if !info.IsDir() {
			planName := strings.TrimSuffix(filepath.Base(path), ".txt") // Remove .txt extension for the key
			if !strings.HasSuffix(path, ".txt") {
				// Log or handle the case if the file does not end with .txt
				xlog.Logf("Skipping file without .txt extension: %s", path)
				return nil // Skip files without .txt
			}

			skills, err := readSkillsFromFile(path)
			if err != nil {
				return err
			}
			skillPlans[planName] = model.SkillPlan{
				Name:   planName,
				Skills: skills,
			}
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to load skill plans from %s: %w", dir, err)
	}

	return skillPlans, nil
}

// CopyEmbeddedPlansToWritable copies plans from the embedded filesystem to the writable path, skipping existing files.
func CopyEmbeddedPlansToWritable(writableDir string) error {
	entries, err := embed.StaticFiles.ReadDir("static/plans")
	if err != nil {
		return fmt.Errorf("failed to read embedded plans: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			fileName := entry.Name()
			destPath := filepath.Join(writableDir, fileName)

			// Copy and confirm existence immediately
			if err := copyEmbeddedFile("static/plans/"+fileName, destPath); err != nil {
				return fmt.Errorf("failed to copy embedded plan %s to %s: %w", fileName, destPath, err)
			}

			// Verify that the file exists after copying
			_, statErr := os.Stat(destPath)
			if statErr != nil {
				xlog.Logf("Error: file %s was copied but does not exist at %s", fileName, destPath)
				return fmt.Errorf("file %s was copied but does not exist at %s", fileName, destPath)
			}
		}
	}
	return nil
}

// copyEmbeddedFile copies a file from the embedded filesystem to the destination path.
func copyEmbeddedFile(srcPath, destPath string) error {
	// Open the embedded file
	srcFile, err := embed.StaticFiles.Open(srcPath)
	if err != nil {
		return fmt.Errorf("failed to open embedded file %s: %w", srcPath, err)
	}
	defer srcFile.Close()

	// Create the destination file
	destFile, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file %s: %w", destPath, err)
	}
	defer destFile.Close()

	// Copy contents from the embedded file to the destination file
	if _, err := io.Copy(destFile, srcFile); err != nil {
		return fmt.Errorf("failed to copy content to destination file %s: %w", destPath, err)
	}

	return nil
}

// readSkillsFromFile reads a single skill file and returns a map of unique skills with their highest levels.
func readSkillsFromFile(filePath string) (map[string]model.Skill, error) {
	// Open directly from the writable path, not the embedded system
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

		// Skip empty lines to prevent parsing errors
		if strings.TrimSpace(line) == "" {
			continue
		}

		// Split line into words
		parts := strings.Fields(line)

		// Ensure there are at least two parts: skill name and level
		if len(parts) < 2 {
			return nil, fmt.Errorf("invalid format in file %s at line %d: %s", filePath, lineNumber, line)
		}

		// Assume last part is the skill level and the rest is the skill name
		skillLevelStr := parts[len(parts)-1]
		skillName := strings.Join(parts[:len(parts)-1], " ")

		// Parse skill level as an integer
		skillLevel, err := strconv.Atoi(skillLevelStr)
		if err != nil {
			return nil, fmt.Errorf("invalid skill level in file %s at line %d: %s", filePath, lineNumber, skillLevelStr)
		}

		// Store only the highest skill level
		if currentSkill, exists := skills[skillName]; !exists || skillLevel > currentSkill.Level {
			skills[skillName] = model.Skill{Name: skillName, Level: skillLevel}
		}
	}

	// Check for scanner errors
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	return skills, nil
}

// SaveSkillPlan saves a new skill plan to the writable path and updates SkillPlans in memory.
func SaveSkillPlan(planName string, skills map[string]model.Skill) error {
	// Check if the skills map is empty
	if len(skills) == 0 {
		return fmt.Errorf("cannot save an empty skill plan for planName: %s", planName)
	}

	writablePath, err := GetWritablePlansPath()
	if err != nil {
		return fmt.Errorf("failed to get writable path: %w", err)
	}

	planFilePath := filepath.Join(writablePath, planName+".txt")
	file, err := os.Create(planFilePath)
	if err != nil {
		return fmt.Errorf("failed to create plan file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	for skillName, skill := range skills {
		line := fmt.Sprintf("%s %d\n", skillName, skill.Level)
		if _, err := writer.WriteString(line); err != nil {
			return fmt.Errorf("failed to write skill to file: %w", err)
		}
	}
	if err := writer.Flush(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	// Update SkillPlans in memory
	planName = strings.TrimSuffix(planName, ".txt")
	SkillPlans[planName] = model.SkillPlan{Name: planName, Skills: skills}

	return nil
}
