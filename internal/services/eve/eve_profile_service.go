package eve

import (
	"fmt"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.EveProfilesService = (*eveProfileService)(nil)

type eveProfileService struct {
	logger         interfaces.Logger
	eveRepo        interfaces.EveProfilesRepository
	configService  interfaces.ConfigService
	esiService     interfaces.ESIService
	accountService interfaces.AccountService
}

func NewEveProfileservice(
	logger interfaces.Logger,
	eveRepo interfaces.EveProfilesRepository, ac interfaces.AccountService,
	esi interfaces.ESIService, c interfaces.ConfigService,
) interfaces.EveProfilesService {
	return &eveProfileService{
		logger:         logger,
		eveRepo:        eveRepo,
		esiService:     esi,
		configService:  c,
		accountService: ac,
	}
}

func (e *eveProfileService) LoadCharacterSettings() ([]model.EveProfile, error) {
	settingsDir, err := e.configService.GetSettingsDir()
	if err != nil {
		return nil, err
	}

	subDirs, err := e.eveRepo.GetSubDirectories(settingsDir)
	if err != nil {
		return nil, fmt.Errorf("failed to get subdirectories: %w", err)
	}

	var settingsData []model.EveProfile
	allCharIDs := make(map[string]struct{})

	for _, sd := range subDirs {
		rawFiles, err := e.eveRepo.ListSettingsFiles(sd, settingsDir)
		if err != nil {
			e.logger.Warnf("Error fetching settings files for subDir %s: %v", sd, err)
			continue
		}

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
				//
				friendlyName := rf.CharOrUserID
				if savedName, ok := e.accountService.GetAccountNameByID(rf.CharOrUserID); ok {
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
	charIdToName, err := e.esiService.ResolveCharacterNames(charIdList)
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

func (e *eveProfileService) SyncDir(subDir, charId, userId string) (int, int, error) {
	settingsDir, err := e.configService.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}

	return e.eveRepo.SyncSubdirectory(subDir, userId, charId, settingsDir)
}

func (e *eveProfileService) SyncAllDir(baseSubDir, charId, userId string) (int, int, error) {
	settingsDir, err := e.configService.GetSettingsDir()
	if err != nil {
		return 0, 0, err
	}
	if settingsDir == "" {
		return 0, 0, fmt.Errorf("SettingsDir not set")
	}

	return e.eveRepo.SyncAllSubdirectories(baseSubDir, userId, charId, settingsDir)
}

func (e *eveProfileService) BackupDir(targetDir, backupDir string) error {
	err := e.eveRepo.BackupDirectory(targetDir, backupDir)
	if err != nil {
		return err
	}

	err = e.configService.UpdateBackupDir(backupDir)
	if err != nil {
		e.logger.Infof("backup succeeded, but updating config failed: %v", err)
	}

	return nil
}
