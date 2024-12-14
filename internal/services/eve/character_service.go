// services/character/character_service.go
package eve

import (
	"fmt"
	"strconv"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/interfaces"
)

var _ interfaces.CharacterService = (*characterService)(nil)

// CharacterService orchestrates character data updates by using ESIService
type characterService struct {
	esi            interfaces.ESIService
	logger         interfaces.Logger
	sysRepo        interfaces.SystemRepository
	skillService   interfaces.SkillService
	accountService interfaces.AccountService
	configService  interfaces.ConfigService
}

func NewCharacterService(esi interfaces.ESIService,
	logger interfaces.Logger,
	sys interfaces.SystemRepository, sk interfaces.SkillService,
	as interfaces.AccountService, c interfaces.ConfigService) interfaces.CharacterService {
	return &characterService{
		esi:            esi,
		logger:         logger,
		sysRepo:        sys,
		skillService:   sk,
		accountService: as,
		configService:  c,
	}
}

func (c *characterService) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	c.logger.Debugf("Processing identity for character ID: %d", charIdentity.Character.CharacterID)

	user, err := c.esi.GetUserInfo(&charIdentity.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	c.logger.Debugf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	characterResponse, err := c.esi.GetCharacter(strconv.FormatInt(charIdentity.Character.CharacterID, 10))
	if err != nil {
		c.logger.Warnf("Failed to get character %s: %v", charIdentity.Character.CharacterName, err)
	}

	skills, err := c.esi.GetCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		c.logger.Warnf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}
	c.logger.Debugf("Fetched %d skills for character %d", len(skills.Skills), charIdentity.Character.CharacterID)

	skillQueue, err := c.esi.GetCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		c.logger.Warnf("Failed to get eve queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	c.logger.Debugf("Fetched %d eve queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	characterLocation, err := c.esi.GetCharacterLocation(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		c.logger.Warnf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	}

	corporationName := ""
	allianceName := ""
	if characterResponse != nil {
		characterCorporation, err := c.esi.GetCorporation(int64(characterResponse.CorporationID), &charIdentity.Token)
		if err != nil {
			c.logger.Warnf("Failed to get corporation for corporation %d: %v", characterResponse.CorporationID, err)
		} else {
			corporationName = characterCorporation.Name
		}
		if characterCorporation != nil && characterCorporation.AllianceID != 0 {
			characterAlliance, err := c.esi.GetAlliance(int64(characterCorporation.AllianceID), &charIdentity.Token)
			if err != nil {
				c.logger.Warnf("Failed to get alliance for character %s: %v", characterCorporation.AllianceID, err)
			} else {
				allianceName = characterAlliance.Name
			}
		}
	}

	c.logger.Debugf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	// Update charIdentity with fetched data
	c.logger.Debugf("updating %s", user.CharacterName)
	charIdentity.Character.UserInfoResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = c.sysRepo.GetSystemName(charIdentity.Character.Location)
	charIdentity.MCT = c.isCharacterTraining(*skillQueue)
	if charIdentity.MCT {
		charIdentity.Training = c.skillService.GetSkillName(charIdentity.Character.SkillQueue[0].SkillID)
	}
	charIdentity.CorporationName = corporationName
	charIdentity.AllianceName = allianceName

	// Initialize maps if nil
	if charIdentity.Character.QualifiedPlans == nil {
		charIdentity.Character.QualifiedPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingPlans == nil {
		charIdentity.Character.PendingPlans = make(map[string]bool)
	}
	if charIdentity.Character.PendingFinishDates == nil {
		charIdentity.Character.PendingFinishDates = make(map[string]*time.Time)
	}
	if charIdentity.Character.MissingSkills == nil {
		charIdentity.Character.MissingSkills = make(map[string]map[string]int32)
	}

	err = c.esi.SaveEsiCache()
	if err != nil {
		c.logger.WithError(err).Infof("failed to save esi cache after processing identity")
	}

	return charIdentity, nil
}

func (c *characterService) isCharacterTraining(queue []model.SkillQueue) bool {
	for _, q := range queue {
		if q.StartDate != nil && q.FinishDate != nil && q.FinishDate.After(time.Now()) {
			c.logger.Debugf("training - start %s, finish %s, eve %d", q.StartDate, q.FinishDate, q.SkillID)
			return true
		}
	}
	return false
}

func (c *characterService) DoesCharacterExist(characterID int64) (bool, *model.CharacterIdentity, error) {
	accounts, err := c.accountService.FetchAccounts()
	if err != nil {
		return false, nil, fmt.Errorf("failed to fetch accounts: %w", err)
	}
	charIdentity := c.findCharacterInAccounts(accounts, characterID)
	if charIdentity == nil {
		return false, nil, nil
	}
	return true, charIdentity, nil
}

func (c *characterService) UpdateCharacterFields(characterID int64, updates map[string]interface{}) error {
	accounts, err := c.accountService.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to fetch accounts: %w", err)
	}

	charIdentity := c.findCharacterInAccounts(accounts, characterID)
	if charIdentity == nil {
		return fmt.Errorf("character not found")
	}

	for key, value := range updates {
		switch key {
		case "Role":
			roleStr, ok := value.(string)
			if !ok {
				return fmt.Errorf("role must be a string")
			}
			// Update roles via configService
			if err := c.configService.UpdateRoles(roleStr); err != nil {
				c.logger.Infof("Failed to update roles: %v", err)
			}
			charIdentity.Role = roleStr

		case "MCT":
			mctBool, ok := value.(bool)
			if !ok {
				return fmt.Errorf("MCT must be boolean")
			}
			charIdentity.MCT = mctBool

		default:
			return fmt.Errorf("unknown update field: %s", key)
		}
	}

	if err := c.accountService.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts: %w", err)
	}

	return nil
}

func (c *characterService) RemoveCharacter(characterID int64) error {
	accounts, err := c.accountService.FetchAccounts()
	if err != nil {
		return fmt.Errorf("failed to fetch accounts: %w", err)
	}

	accountIndex, charIndex, found := c.findCharacterIndices(accounts, characterID)
	if !found {
		return fmt.Errorf("character not found in accounts")
	}

	accountName := accounts[accountIndex].Name
	c.logger.Infof("Removing character %d from account %s", characterID, accountName)

	// Remove character
	accounts[accountIndex].Characters = append(
		accounts[accountIndex].Characters[:charIndex],
		accounts[accountIndex].Characters[charIndex+1:]...,
	)

	if err := c.accountService.SaveAccounts(accounts); err != nil {
		return fmt.Errorf("failed to save accounts after character removal: %w", err)
	}

	// Optional verification log if needed
	updatedAccounts, err := c.accountService.FetchAccounts()
	if err == nil {
		idx, found := c.findAccountIndex(updatedAccounts, accountName)
		if found {
			c.logger.Infof("after removal - length of %s characters is %d", accountName, len(updatedAccounts[idx].Characters))
		}
	}

	return nil
}

func (c *characterService) findCharacterInAccounts(accounts []model.Account, characterID int64) *model.CharacterIdentity {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				return &accounts[i].Characters[j]
			}
		}
	}
	return nil
}

func (c *characterService) findCharacterIndices(accounts []model.Account, characterID int64) (int, int, bool) {
	for i := range accounts {
		for j := range accounts[i].Characters {
			if accounts[i].Characters[j].Character.CharacterID == characterID {
				return i, j, true
			}
		}
	}
	return 0, 0, false
}

func (c *characterService) findAccountIndex(accounts []model.Account, accountName string) (int, bool) {
	for i, acc := range accounts {
		if acc.Name == accountName {
			return i, true
		}
	}
	return -1, false
}
