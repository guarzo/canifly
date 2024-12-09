// persist/filesystem.go
package persist

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
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

// GetFilesForDropdown returns available character and user files for a subdir.
func (ds *DataStore) GetFilesForDropdown(subDir, settingsDir string) ([]model.CharFile, []model.UserFile, error) {
	configData, err := ds.FetchConfigData()
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

// BackupDirectory creates a backup tar.gz of the settingsDir
func (ds *DataStore) BackupDirectory(settingsDir string) (bool, string) {
	userDataPath, _ := os.UserConfigDir()
	now := time.Now()
	formattedDate := now.Format("2006-01-02_15-04-05")

	backupFileName := fmt.Sprintf("%s_%s.bak.tar.gz", filepath.Base(settingsDir), formattedDate)
	backupFilePath := filepath.Join(userDataPath, backupFileName)

	f, err := os.Create(backupFilePath)
	if err != nil {
		return false, err.Error()
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	baseDir := filepath.Dir(settingsDir)
	baseName := filepath.Base(settingsDir)

	filepath.Walk(filepath.Join(baseDir, baseName), func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, _ := filepath.Rel(baseDir, p)

		header, err := tar.FileInfoHeader(info, relPath)
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		if !info.IsDir() {
			srcFile, err := os.Open(p)
			if err != nil {
				return err
			}
			defer srcFile.Close()
			if _, err := io.Copy(tw, srcFile); err != nil {
				return err
			}
		}
		return nil
	})

	return true, "Backup created successfully at: " + backupFilePath
}

func (ds *DataStore) GetSubDirectories(settingsDir string) ([]string, error) {
	entries, err := os.ReadDir(settingsDir)
	if err != nil {
		ds.logger.Errorf("failed to read %s, with error %v", settingsDir, err)
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

func (ds *DataStore) SaveUserSelections(selections model.UserSelections) error {
	selectionsPath, err := ds.getUserSelectionsFilePath()
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(selections, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(selectionsPath, data, 0644)
}

func (ds *DataStore) LoadUserSelections() (model.UserSelections, error) {
	selectionsPath, err := ds.getUserSelectionsFilePath()
	if err != nil {
		ds.logger.Infof("no selection file path %v", err)
		return nil, err
	}
	data, err := os.ReadFile(selectionsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return make(model.UserSelections), nil
		}
		ds.logger.Warnf("failed in reading file, %v", err)
		return nil, err
	}

	var selections model.UserSelections
	if err = json.Unmarshal(data, &selections); err != nil {
		ds.logger.Infof("failed to unmarshal selections %v", err)
		return nil, err
	}
	return selections, nil
}

func (ds *DataStore) getUserSelectionsFilePath() (string, error) {
	identityPath, err := ds.GetWriteablePath()
	if err != nil {
		return "", err
	}
	return filepath.Join(identityPath, "userSelections.json"), nil
}

func (ds *DataStore) SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error) {
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

	return ds.applyContentToSubDir(subDirPath, userFileName, charFileName, userContent, charContent)
}

func (ds *DataStore) SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error) {
	baseSubDirPath := filepath.Join(settingsDir, baseSubDir)
	if _, err := os.Stat(baseSubDirPath); os.IsNotExist(err) {
		return 0, 0, fmt.Errorf("base subdirectory does not exist: %s", baseSubDirPath)
	}

	userFileName := "core_user_" + userId + ".dat"
	charFileName := "core_char_" + charId + ".dat"

	userFilePath := filepath.Join(baseSubDirPath, userFileName)
	charFilePath := filepath.Join(baseSubDirPath, charFileName)

	userContent, userErr := os.ReadFile(userFilePath)
	if userErr != nil {
		return 0, 0, fmt.Errorf("failed to read user file %s: %v", userFilePath, userErr)
	}

	charContent, charErr := os.ReadFile(charFilePath)
	if charErr != nil {
		return 0, 0, fmt.Errorf("failed to read char file %s: %v", charFilePath, charErr)
	}

	subDirs, err := ds.GetSubDirectories(settingsDir)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get subdirectories: %v", err)
	}

	totalUserCopied := 0
	totalCharCopied := 0

	for _, otherSubDir := range subDirs {
		if otherSubDir == baseSubDir {
			continue
		}

		otherSubDirPath := filepath.Join(settingsDir, otherSubDir)
		uCopied, cCopied, err := ds.applyContentToSubDir(otherSubDirPath, userFileName, charFileName, userContent, charContent)
		if err != nil {
			ds.logger.Warnf("Error applying content to subdir %s: %v", otherSubDir, err)
			continue
		}

		totalUserCopied += uCopied
		totalCharCopied += cCopied
	}

	return totalUserCopied, totalCharCopied, nil
}

func (ds *DataStore) applyContentToSubDir(
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
					ds.logger.Warnf("Failed to write user file %s: %v", fPath, err)
				}
			}
		}

		if strings.HasPrefix(fName, "core_char_") && strings.HasSuffix(fName, ".dat") && fName != charFileName {
			if charContent != nil {
				if err := os.WriteFile(fPath, charContent, 0644); err == nil {
					charFilesCopied++
				} else {
					ds.logger.Warnf("Failed to write char file %s: %v", fPath, err)
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
