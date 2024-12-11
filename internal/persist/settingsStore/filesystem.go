// persist/filesystem.go
package settingsStore

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"github.com/guarzo/canifly/internal/persist"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/model"
)

const selectionsFileName = "selections.json"

// GetFilesForDropdown returns available character and user files for a subdir.
func (s *SettingsStore) GetFilesForDropdown(subDir, settingsDir string) ([]model.CharFile, []model.UserFile, error) {
	configData, err := s.FetchConfigData()
	if err != nil {
		return nil, nil, err
	}

	directory := filepath.Join(settingsDir, subDir)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, nil, err
	}

	var charFiles []model.CharFile
	var userFiles []model.UserFile

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		file := entry.Name()
		fullPath := filepath.Join(directory, file)

		info, err := os.Stat(fullPath)
		if err != nil {
			continue
		}
		mtime := info.ModTime().Format(time.RFC3339)

		if strings.HasPrefix(file, "core_char_") && strings.HasSuffix(file, ".dat") {
			charId := strings.TrimSuffix(strings.TrimPrefix(file, "core_char_"), ".dat")
			if matched, _ := regexp.MatchString(`^\d+$`, charId); matched {
				charFiles = append(charFiles, model.CharFile{
					File:   file,
					CharId: charId,
					Name:   fmt.Sprintf("CharID:%s", charId),
					Mtime:  mtime,
				})
			}
		} else if strings.HasPrefix(file, "core_user_") && strings.HasSuffix(file, ".dat") {
			userId := strings.TrimSuffix(strings.TrimPrefix(file, "core_user_"), ".dat")
			if matched, _ := regexp.MatchString(`^\d+$`, userId); matched {
				userFiles = append(userFiles, model.UserFile{
					File:   file,
					UserId: userId,
					Name:   getSavedNameForUserID(userId, configData),
					Mtime:  mtime,
				})
			}
		}
	}

	return charFiles, userFiles, nil
}

func getSavedNameForUserID(userId string, configData *model.ConfigData) string {
	name := userId
	if savedName, ok := configData.UserAccount[userId]; ok {
		name = savedName
	}
	return name
}

// BackupDirectory creates a backup tar.gz of only the directories under targetDir that start with "settings_"
func (s *SettingsStore) BackupDirectory(targetDir, backupDir string) error {
	s.logger.Infof("Starting backup of settings directories from %s to %s", targetDir, backupDir)

	subDirs, err := s.GetSubDirectories(targetDir)
	if err != nil {
		s.logger.Errorf("Failed to get subdirectories from %s: %v", targetDir, err)
		return err
	}

	if len(subDirs) == 0 {
		errMsg := fmt.Sprintf("No settings_ subdirectories found in %s", targetDir)
		s.logger.Warnf(errMsg)
		return fmt.Errorf(errMsg)
	}

	now := time.Now()
	formattedDate := now.Format("2006-01-02_15-04-05")

	backupFileName := fmt.Sprintf("%s_%s.bak.tar.gz", filepath.Base(targetDir), formattedDate)
	backupFilePath := filepath.Join(backupDir, backupFileName)

	s.logger.Infof("Creating backup file at %s", backupFilePath)
	f, err := os.Create(backupFilePath)
	if err != nil {
		s.logger.Errorf("Failed to create backup file %s: %v", backupFilePath, err)
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	for _, dir := range subDirs {
		fullPath := filepath.Join(targetDir, dir)
		s.logger.Infof("Backing up subdirectory: %s", fullPath)
		err = filepath.Walk(fullPath, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				s.logger.Errorf("Error walking through %s: %v", p, err)
				return err
			}
			relPath, _ := filepath.Rel(filepath.Dir(targetDir), p)

			header, err := tar.FileInfoHeader(info, relPath)
			if err != nil {
				s.logger.Errorf("Failed to get FileInfoHeader for %s: %v", p, err)
				return err
			}
			header.Name = relPath

			if err = tw.WriteHeader(header); err != nil {
				s.logger.Errorf("Failed to write tar header for %s: %v", p, err)
				return err
			}

			if !info.IsDir() {
				srcFile, err := os.Open(p)
				if err != nil {
					s.logger.Errorf("Failed to open file %s: %v", p, err)
					return err
				}
				defer srcFile.Close()
				if _, err := io.Copy(tw, srcFile); err != nil {
					s.logger.Errorf("Failed to copy file %s into tar: %v", p, err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			s.logger.Errorf("Error walking subdirectory %s: %v", dir, err)
			return err
		}
	}

	s.logger.Infof("Backup completed successfully: %s", backupFilePath)
	return nil
}

func (s *SettingsStore) GetSubDirectories(settingsDir string) ([]string, error) {
	entries, err := os.ReadDir(settingsDir)
	if err != nil {
		s.logger.Errorf("failed to read %s, with error %v", settingsDir, err)
		return nil, err
	}

	var dirs []string
	for _, e := range entries {
		if e.IsDir() && strings.HasPrefix(e.Name(), "settings_") {
			dirs = append(dirs, e.Name())
		}
	}
	return dirs, nil
}

func (s *SettingsStore) SaveUserSelections(selections model.UserSelections) error {
	selectionsPath, err := s.getUserSelectionsFilePath()
	if err != nil {
		return err
	}
	return persist.OldSaveJson(selectionsPath, selections)
}

func (s *SettingsStore) FetchUserSelections() (model.UserSelections, error) {
	selectionsPath, err := s.getUserSelectionsFilePath()
	if err != nil {
		s.logger.Infof("no selection file path %v", err)
		return nil, err
	}

	if exists, _ := persist.FileExists(selectionsPath); !exists {
		return make(model.UserSelections), nil
	}

	var selections model.UserSelections
	if err = persist.OldReadJson(selectionsPath, &selections); err != nil {
		s.logger.Infof("failed to load selections %v", err)
		return nil, err
	}
	return selections, nil
}

func (s *SettingsStore) getUserSelectionsFilePath() (string, error) {
	return getConfigFileName(selectionsFileName)
}

func (s *SettingsStore) SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error) {
	subDirPath := filepath.Join(settingsDir, subDir)
	if _, err := os.Stat(subDirPath); os.IsNotExist(err) {
		return 0, 0, fmt.Errorf("subdirectory does not exist: %s", subDirPath)
	}

	userFileName := "core_user_" + userId + ".dat"
	charFileName := "core_char_" + charId + ".dat"

	userFilePath := filepath.Join(subDirPath, userFileName)
	charFilePath := filepath.Join(subDirPath, charFileName)

	userContent, userErr := os.ReadFile(userFilePath)
	charContent, charErr := os.ReadFile(charFilePath)

	if userErr != nil {
		return 0, 0, fmt.Errorf("failed to read user file %s: %v", userFilePath, userErr)
	}
	if charErr != nil {
		return 0, 0, fmt.Errorf("failed to read char file %s: %v", charFilePath, charErr)
	}

	return s.applyContentToSubDir(subDirPath, userFileName, charFileName, userContent, charContent)
}

func (s *SettingsStore) SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error) {
	s.logger.Infof("Starting SyncAllSubdirectories with baseSubDir=%s, userId=%s, charId=%s", baseSubDir, userId, charId)

	baseSubDirPath := filepath.Join(settingsDir, baseSubDir)
	if _, err := os.Stat(baseSubDirPath); os.IsNotExist(err) {
		s.logger.Errorf("Base subdirectory does not exist: %s", baseSubDirPath)
		return 0, 0, fmt.Errorf("base subdirectory does not exist: %s", baseSubDirPath)
	}

	userFileName := "core_user_" + userId + ".dat"
	charFileName := "core_char_" + charId + ".dat"
	userFilePath := filepath.Join(baseSubDirPath, userFileName)
	charFilePath := filepath.Join(baseSubDirPath, charFileName)

	s.logger.Infof("Reading user file: %s", userFilePath)
	userContent, userErr := os.ReadFile(userFilePath)
	if userErr != nil {
		s.logger.Errorf("Failed to read user file %s: %v", userFilePath, userErr)
		return 0, 0, fmt.Errorf("failed to read user file %s: %v", userFilePath, userErr)
	}

	s.logger.Infof("Reading character file: %s", charFilePath)
	charContent, charErr := os.ReadFile(charFilePath)
	if charErr != nil {
		s.logger.Errorf("Failed to read char file %s: %v", charFilePath, charErr)
		return 0, 0, fmt.Errorf("failed to read char file %s: %v", charFilePath, charErr)
	}

	s.logger.Infof("Retrieving all settings_ subdirectories from %s", settingsDir)
	subDirs, err := s.GetSubDirectories(settingsDir)
	if err != nil {
		s.logger.Errorf("Failed to get subdirectories from %s: %v", settingsDir, err)
		return 0, 0, fmt.Errorf("failed to get subdirectories: %v", err)
	}

	totalUserCopied := 0
	totalCharCopied := 0

	for _, otherSubDir := range subDirs {
		if otherSubDir == baseSubDir {
			continue
		}
		s.logger.Infof("Applying content to subdir: %s", otherSubDir)
		otherSubDirPath := filepath.Join(settingsDir, otherSubDir)
		uCopied, cCopied, err := s.applyContentToSubDir(otherSubDirPath, userFileName, charFileName, userContent, charContent)
		if err != nil {
			s.logger.Warnf("Error applying content to subdir %s: %v", otherSubDir, err)
			continue
		}
		s.logger.Infof("Successfully applied content to %s: %d user files, %d char files copied.", otherSubDir, uCopied, cCopied)
		totalUserCopied += uCopied
		totalCharCopied += cCopied
	}

	s.logger.Infof("SyncAllSubdirectories complete: %d total user files, %d total char files copied.", totalUserCopied, totalCharCopied)
	return totalUserCopied, totalCharCopied, nil
}

func (s *SettingsStore) applyContentToSubDir(
	dirPath string,
	userFileName string,
	charFileName string,
	userContent []byte,
	charContent []byte,
) (int, int, error) {
	entries, err := os.ReadDir(dirPath)
	if err != nil {
		return 0, 0, err
	}

	userFilesCopied := 0
	charFilesCopied := 0

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		fName := entry.Name()
		fPath := filepath.Join(dirPath, fName)

		if strings.HasPrefix(fName, "core_user_") && strings.HasSuffix(fName, ".dat") && fName != userFileName {
			if userContent != nil {
				if err := os.WriteFile(fPath, userContent, 0644); err == nil {
					userFilesCopied++
				} else {
					s.logger.Warnf("Failed to write user file %s: %v", fPath, err)
				}
			}
		}

		if strings.HasPrefix(fName, "core_char_") && strings.HasSuffix(fName, ".dat") && fName != charFileName {
			if charContent != nil {
				if err := os.WriteFile(fPath, charContent, 0644); err == nil {
					charFilesCopied++
				} else {
					s.logger.Warnf("Failed to write char file %s: %v", fPath, err)
				}
			}
		}
	}

	return userFilesCopied, charFilesCopied, nil
}

func isWSL() bool {
	if runtime.GOOS == "linux" {
		data, err := os.ReadFile("/proc/version")
		if err == nil && strings.Contains(string(data), "microsoft") {
			return true
		}
	}
	return false
}

func getWindowsHomeInWSL() (string, error) {
	cmd := "cmd.exe"
	args := []string{"/C", "echo", "%USERPROFILE%"}
	out, err := runCommand(cmd, args)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve Windows home directory in WSL: %w", err)
	}
	windowsHome := strings.TrimSpace(out)
	windowsHome = strings.ReplaceAll(windowsHome, "\\", "/")

	out2, err := runCommand("wslpath", []string{"-u", windowsHome})
	if err != nil {
		return "", fmt.Errorf("failed to convert Windows home path to WSL format: %w", err)
	}
	return strings.TrimSpace(out2), nil
}

func runCommand(name string, args []string) (string, error) {
	cmd := exec.Command(name, args...)
	output, err := cmd.Output()
	return string(output), err
}
