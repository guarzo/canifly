package profile

import (
	"fmt"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

// Service implements interfaces.ProfileService.
type Service struct {
	repo        interfaces.EveProfilesRepository
	config      interfaces.ConfigurationService
	accountMgmt interfaces.AccountManagementService
	esi         interfaces.ESIAPIService
	logger      interfaces.Logger
}

// NewService constructs a ProfileService with narrow dependencies.
func NewService(
	repo interfaces.EveProfilesRepository,
	config interfaces.ConfigurationService,
	accountMgmt interfaces.AccountManagementService,
	esi interfaces.ESIAPIService,
	logger interfaces.Logger,
) *Service {
	return &Service{
		repo:        repo,
		config:      config,
		accountMgmt: accountMgmt,
		esi:         esi,
		logger:      logger,
	}
}

var _ interfaces.ProfileService = (*Service)(nil)

func (s *Service) LoadCharacterSettings() ([]model.EveProfile, error) {
	settingsDir, err := s.config.GetSettingsDir()
	if err != nil {
		s.logger.Errorf("Failed to get settings directory: %v", err)
		return nil, err
	}
	s.logger.Infof("Loading character settings from directory: %s", settingsDir)

	subDirs, err := s.repo.GetSubDirectories(settingsDir)
	if err != nil {
		s.logger.Errorf("Failed to get subdirectories from %s: %v", settingsDir, err)
		return nil, err
	}
	s.logger.Infof("Found %d subdirectories in settings directory", len(subDirs))
	for _, dir := range subDirs {
		s.logger.Infof("  - %s", dir)
	}

	var settingsData []model.EveProfile
	allCharIDs := make(map[string]struct{})

	for _, sd := range subDirs {
		s.logger.Infof("Processing subdirectory: %s", sd)
		rawFiles, err := s.repo.ListSettingsFiles(sd, settingsDir)
		if err != nil {
			s.logger.Warnf("Error fetching settings files for subDir %s: %v", sd, err)
			continue
		}
		s.logger.Infof("Found %d settings files in %s", len(rawFiles), sd)

		var charFiles []model.CharFile
		var userFiles []model.UserFile

		for _, rf := range rawFiles {
			if rf.IsChar {
				// Just record charId for later ESI resolution
				allCharIDs[rf.CharOrUserID] = struct{}{}
				charFiles = append(charFiles, model.CharFile{
					File:   rf.FileName,
					CharId: rf.CharOrUserID,
					Name:   "CharID:" + rf.CharOrUserID, // Temporary name, will update after ESI lookup
					Mtime:  rf.Mtime,
				})
			} else {
				friendlyName := rf.CharOrUserID
				if savedName, ok := s.accountMgmt.GetAccountNameByID(rf.CharOrUserID); ok {
					friendlyName = savedName
				}
				userFiles = append(userFiles, model.UserFile{
					File:   rf.FileName,
					UserId: rf.CharOrUserID,
					Name:   friendlyName,
					Mtime:  rf.Mtime,
				})
			}
		}

		settingsData = append(settingsData, model.EveProfile{
			Profile:            sd,
			AvailableCharFiles: charFiles,
			AvailableUserFiles: userFiles,
		})
	}

	// Resolve character names via ESI
	var charIdList []string
	for id := range allCharIDs {
		charIdList = append(charIdList, id)
	}
	charIdToName, err := s.esi.ResolveCharacterNames(charIdList)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve character names: %w", err)
	}

	// Update character files with resolved names
	for si, sd := range settingsData {
		var filteredChars []model.CharFile
		for _, cf := range sd.AvailableCharFiles {
			if name, ok := charIdToName[cf.CharId]; ok && name != "" {
				cf.Name = name
				filteredChars = append(filteredChars, cf)
			}
		}
		settingsData[si].AvailableCharFiles = filteredChars
	}

	return settingsData, nil
}

func (s *Service) BackupDir(targetDir, backupDir string) error {
	// Backup the EVE "settings_" directories
	if err := s.repo.BackupDirectory(targetDir, backupDir); err != nil {
		return err
	}

	// Also backup JSON files from config
	return s.config.BackupJSONFiles(backupDir)
}

func (s *Service) SyncDir(subDir, charId, userId string) (int, int, error) {
	settingsDir, err := s.config.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}

	return s.repo.SyncSubdirectory(subDir, userId, charId, settingsDir)
}

func (s *Service) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	settingsDir, err := s.config.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}
	if settingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return s.repo.SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir)
}
