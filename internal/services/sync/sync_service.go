package sync

import (
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// SyncService handles settings synchronization
type SyncService struct {
	eveDataService interfaces.EVEDataService
	eveRepo        interfaces.EveProfilesRepository
	configService  interfaces.ConfigurationService
	logger         interfaces.Logger
}

// NewSyncService creates a new sync service
func NewSyncService(
	eveData interfaces.EVEDataService,
	eveRepo interfaces.EveProfilesRepository,
	config interfaces.ConfigurationService,
	logger interfaces.Logger,
) interfaces.SyncService {
	return &SyncService{
		eveDataService: eveData,
		eveRepo:        eveRepo,
		configService:  config,
		logger:         logger,
	}
}

// SyncDirectory syncs a specific subdirectory for a character/user
func (s *SyncService) SyncDirectory(subDir, charId, userId string) (int, int, error) {
	s.logger.Infof("Syncing directory %s for char %s and user %s", subDir, charId, userId)
	return s.eveDataService.SyncDir(subDir, charId, userId)
}

// SyncAllDirectories syncs all subdirectories for a character/user
func (s *SyncService) SyncAllDirectories(baseSubDir, charId, userId string) (int, int, error) {
	s.logger.Infof("Syncing all directories from base %s for char %s and user %s", baseSubDir, charId, userId)
	return s.eveDataService.SyncAllDir(baseSubDir, charId, userId)
}

// GetSyncDirectories gets list of available sync directories
func (s *SyncService) GetSyncDirectories() ([]string, error) {
	settingsDir, err := s.configService.GetSettingsDir()
	if err != nil {
		return nil, err
	}

	return s.eveRepo.GetSubDirectories(settingsDir)
}

// BackupDirectory backs up a directory
func (s *SyncService) BackupDirectory(targetDir, backupDir string) error {
	s.logger.Infof("Backing up %s to %s", targetDir, backupDir)
	return s.eveDataService.BackupDir(targetDir, backupDir)
}
