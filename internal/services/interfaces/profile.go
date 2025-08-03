package interfaces

import "github.com/guarzo/canifly/internal/model"

// ProfileService handles EVE profile management
type ProfileService interface {
	LoadCharacterSettings() ([]model.EveProfile, error)
	BackupDir(targetDir, backupDir string) error
	SyncDir(subDir, charId, userId string) (int, int, error)
	SyncAllDir(baseSubDir, charId, userId string) (int, int, error)
}
