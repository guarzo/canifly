package interfaces

// SyncService handles settings synchronization operations
type SyncService interface {
	// Sync a specific subdirectory for a character/user
	SyncDirectory(subDir, charId, userId string) (int, int, error)
	
	// Sync all subdirectories for a character/user
	SyncAllDirectories(baseSubDir, charId, userId string) (int, int, error)
	
	// Get list of available sync directories
	GetSyncDirectories() ([]string, error)
	
	// Backup directory
	BackupDirectory(targetDir, backupDir string) error
}