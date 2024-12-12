package eve

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// EveProfilesStore manages EVE-specific settings file operations
type EveProfilesStore struct {
	logger interfaces.Logger
}

// NewEveProfilesStore returns a new instance of EveProfilesStore
func NewEveProfilesStore(logger interfaces.Logger) *EveProfilesStore {
	return &EveProfilesStore{
		logger: logger,
	}
}

// ListSettingsFiles returns raw file information for character and user files in the given subDir.
func (e *EveProfilesStore) ListSettingsFiles(subDir, settingsDir string) ([]model.RawFileInfo, error) {
	directory := filepath.Join(settingsDir, subDir)
	entries, err := os.ReadDir(directory)
	if err != nil {
		return nil, err
	}

	var results []model.RawFileInfo
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

		// Check for char file pattern: core_char_<charId>.dat
		if strings.HasPrefix(file, "core_char_") && strings.HasSuffix(file, ".dat") {
			charId := strings.TrimSuffix(strings.TrimPrefix(file, "core_char_"), ".dat")
			if matched, _ := regexp.MatchString(`^\d+$`, charId); matched {
				results = append(results, model.RawFileInfo{
					FileName:     file,
					CharOrUserID: charId,
					IsChar:       true,
					Mtime:        mtime,
				})
			}
		} else if strings.HasPrefix(file, "core_user_") && strings.HasSuffix(file, ".dat") {
			userId := strings.TrimSuffix(strings.TrimPrefix(file, "core_user_"), ".dat")
			if matched, _ := regexp.MatchString(`^\d+$`, userId); matched {
				results = append(results, model.RawFileInfo{
					FileName:     file,
					CharOrUserID: userId,
					IsChar:       false,
					Mtime:        mtime,
				})
			}
		}
	}

	return results, nil
}

// BackupDirectory creates a backup tar.gz of directories under targetDir that start with "settings_"
func (e *EveProfilesStore) BackupDirectory(targetDir, backupDir string) error {
	e.logger.Infof("Starting backup of settings directories from %s to %s", targetDir, backupDir)

	subDirs, err := e.GetSubDirectories(targetDir)
	if err != nil {
		e.logger.Errorf("Failed to get subdirectories from %s: %v", targetDir, err)
		return err
	}

	if len(subDirs) == 0 {
		errMsg := fmt.Sprintf("No settings_ subdirectories found in %s", targetDir)
		e.logger.Warnf(errMsg)
		return fmt.Errorf(errMsg)
	}

	now := time.Now()
	formattedDate := now.Format("2006-01-02_15-04-05")

	backupFileName := fmt.Sprintf("%s_%s.bak.tar.gz", filepath.Base(targetDir), formattedDate)
	backupFilePath := filepath.Join(backupDir, backupFileName)

	e.logger.Infof("Creating backup file at %s", backupFilePath)
	f, err := os.Create(backupFilePath)
	if err != nil {
		e.logger.Errorf("Failed to create backup file %s: %v", backupFilePath, err)
		return err
	}
	defer f.Close()

	gz := gzip.NewWriter(f)
	defer gz.Close()

	tw := tar.NewWriter(gz)
	defer tw.Close()

	for _, dir := range subDirs {
		fullPath := filepath.Join(targetDir, dir)
		e.logger.Infof("Backing up subdirectory: %s", fullPath)
		err = filepath.Walk(fullPath, func(p string, info os.FileInfo, err error) error {
			if err != nil {
				e.logger.Errorf("Error walking through %s: %v", p, err)
				return err
			}
			relPath, _ := filepath.Rel(filepath.Dir(targetDir), p)

			header, err := tar.FileInfoHeader(info, relPath)
			if err != nil {
				e.logger.Errorf("Failed to get FileInfoHeader for %s: %v", p, err)
				return err
			}
			header.Name = relPath

			if err = tw.WriteHeader(header); err != nil {
				e.logger.Errorf("Failed to write tar header for %s: %v", p, err)
				return err
			}

			if !info.IsDir() {
				srcFile, err := os.Open(p)
				if err != nil {
					e.logger.Errorf("Failed to open file %s: %v", p, err)
					return err
				}
				defer srcFile.Close()
				if _, err := io.Copy(tw, srcFile); err != nil {
					e.logger.Errorf("Failed to copy file %s into tar: %v", p, err)
					return err
				}
			}
			return nil
		})
		if err != nil {
			e.logger.Errorf("Error walking subdirectory %s: %v", dir, err)
			return err
		}
	}

	e.logger.Infof("Backup completed successfully: %s", backupFilePath)
	return nil
}

// GetSubDirectories returns subdirs in settingsDir that start with "settings_".
func (e *EveProfilesStore) GetSubDirectories(settingsDir string) ([]string, error) {
	entries, err := os.ReadDir(settingsDir)
	if err != nil {
		e.logger.Errorf("failed to read %s, with error %v", settingsDir, err)
		return nil, err
	}

	var dirs []string
	for _, ent := range entries {
		if ent.IsDir() && strings.HasPrefix(ent.Name(), "settings_") {
			dirs = append(dirs, ent.Name())
		}
	}
	return dirs, nil
}

func (e *EveProfilesStore) SyncSubdirectory(subDir, userId, charId, settingsDir string) (int, int, error) {
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

	return e.applyContentToSubDir(subDirPath, userFileName, charFileName, userContent, charContent)
}

func (e *EveProfilesStore) SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir string) (int, int, error) {
	e.logger.Infof("Starting SyncAllSubdirectories with baseSubDir=%s, userId=%s, charId=%s", baseSubDir, userId, charId)

	baseSubDirPath := filepath.Join(settingsDir, baseSubDir)
	if _, err := os.Stat(baseSubDirPath); os.IsNotExist(err) {
		e.logger.Errorf("Base subdirectory does not exist: %s", baseSubDirPath)
		return 0, 0, fmt.Errorf("base subdirectory does not exist: %s", baseSubDirPath)
	}

	userFileName := "core_user_" + userId + ".dat"
	charFileName := "core_char_" + charId + ".dat"
	userFilePath := filepath.Join(baseSubDirPath, userFileName)
	charFilePath := filepath.Join(baseSubDirPath, charFileName)

	e.logger.Infof("Reading user file: %s", userFilePath)
	userContent, userErr := os.ReadFile(userFilePath)
	if userErr != nil {
		e.logger.Errorf("Failed to read user file %s: %v", userFilePath, userErr)
		return 0, 0, fmt.Errorf("failed to read user file %s: %v", userFilePath, userErr)
	}

	e.logger.Infof("Reading character file: %s", charFilePath)
	charContent, charErr := os.ReadFile(charFilePath)
	if charErr != nil {
		e.logger.Errorf("Failed to read char file %s: %v", charFilePath, charErr)
		return 0, 0, fmt.Errorf("failed to read char file %s: %v", charFilePath, charErr)
	}

	e.logger.Infof("Retrieving all settings_ subdirectories from %s", settingsDir)
	subDirs, err := e.GetSubDirectories(settingsDir)
	if err != nil {
		e.logger.Errorf("Failed to get subdirectories from %s: %v", settingsDir, err)
		return 0, 0, fmt.Errorf("failed to get subdirectories: %v", err)
	}

	totalUserCopied := 0
	totalCharCopied := 0

	for _, otherSubDir := range subDirs {
		if otherSubDir == baseSubDir {
			continue
		}
		e.logger.Infof("Applying content to subdir: %s", otherSubDir)
		otherSubDirPath := filepath.Join(settingsDir, otherSubDir)
		uCopied, cCopied, err := e.applyContentToSubDir(otherSubDirPath, userFileName, charFileName, userContent, charContent)
		if err != nil {
			e.logger.Warnf("Error applying content to subdir %s: %v", otherSubDir, err)
			continue
		}
		e.logger.Infof("Successfully applied content to %s: %d user files, %d char files copied.", otherSubDir, uCopied, cCopied)
		totalUserCopied += uCopied
		totalCharCopied += cCopied
	}

	e.logger.Infof("SyncAllSubdirectories complete: %d total user files, %d total char files copied.", totalUserCopied, totalCharCopied)
	return totalUserCopied, totalCharCopied, nil
}

func (e *EveProfilesStore) applyContentToSubDir(
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

		if strings.HasPrefix(fName, "core_user_") && strings.HasSuffix(fName, ".dat") && fName != userFileName && userContent != nil {
			if err := os.WriteFile(fPath, userContent, 0644); err == nil {
				userFilesCopied++
			} else {
				e.logger.Warnf("Failed to write user file %s: %v", fPath, err)
			}
		}

		if strings.HasPrefix(fName, "core_char_") && strings.HasSuffix(fName, ".dat") && fName != charFileName && charContent != nil {
			if err := os.WriteFile(fPath, charContent, 0644); err == nil {
				charFilesCopied++
			} else {
				e.logger.Warnf("Failed to write char file %s: %v", fPath, err)
			}
		}
	}

	return userFilesCopied, charFilesCopied, nil
}
