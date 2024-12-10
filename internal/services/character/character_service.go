// services/character/character_service.go
package character

import (
	"fmt"
	"github.com/guarzo/canifly/internal/services/interfaces"
	"time"

	"github.com/guarzo/canifly/internal/auth"
	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/persist"
	"github.com/guarzo/canifly/internal/services/esi"
	"github.com/sirupsen/logrus"
)

// CharacterService orchestrates character data updates by using ESIService and AuthClient.
type characterService struct {
	esi         esi.ESIService
	auth        auth.AuthClient
	logger      *logrus.Logger
	sysResolver persist.SystemNameResolver
}

func NewCharacterService(esi esi.ESIService, auth auth.AuthClient, logger *logrus.Logger, resolver persist.SystemNameResolver) interfaces.CharacterService {
	return &characterService{
		esi:         esi,
		auth:        auth,
		logger:      logger,
		sysResolver: resolver,
	}
}

// ProcessIdentity previously lived in esiService. Now it's here.
func (cs *characterService) ProcessIdentity(charIdentity *model.CharacterIdentity) (*model.CharacterIdentity, error) {
	cs.logger.Debugf("Processing identity for character ID: %d", charIdentity.Character.CharacterID)

	// Refresh token if needed
	newToken, err := cs.auth.RefreshToken(charIdentity.Token.RefreshToken)
	if err != nil {
		return nil, fmt.Errorf("failed to refresh token for character %d: %v", charIdentity.Character.CharacterID, err)
	}
	charIdentity.Token = *newToken
	cs.logger.Debugf("Token refreshed for character %d", charIdentity.Character.CharacterID)

	user, err := cs.esi.GetUserInfo(&charIdentity.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	cs.logger.Debugf("Fetched user info for character %s (ID: %d)", user.CharacterName, user.CharacterID)

	skills, err := cs.esi.GetCharacterSkills(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		cs.logger.Warnf("Failed to get skills for character %d: %v", charIdentity.Character.CharacterID, err)
		skills = &model.CharacterSkillsResponse{Skills: []model.SkillResponse{}}
	}
	cs.logger.Debugf("Fetched %d skills for character %d", len(skills.Skills), charIdentity.Character.CharacterID)

	skillQueue, err := cs.esi.GetCharacterSkillQueue(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		cs.logger.Warnf("Failed to get skill queue for character %d: %v", charIdentity.Character.CharacterID, err)
		skillQueue = &[]model.SkillQueue{}
	}
	cs.logger.Debugf("Fetched %d skill queue entries for character %d", len(*skillQueue), charIdentity.Character.CharacterID)

	characterLocation, err := cs.esi.GetCharacterLocation(charIdentity.Character.CharacterID, &charIdentity.Token)
	if err != nil {
		cs.logger.Warnf("Failed to get location for character %d: %v", charIdentity.Character.CharacterID, err)
		characterLocation = 0
	}
	cs.logger.Debugf("Character %d is located at %d", charIdentity.Character.CharacterID, characterLocation)

	// Update charIdentity with fetched data
	charIdentity.Character.UserInfoResponse = *user
	charIdentity.Character.CharacterSkillsResponse = *skills
	charIdentity.Character.SkillQueue = *skillQueue
	charIdentity.Character.Location = characterLocation
	charIdentity.Character.LocationName = cs.sysResolver.GetSystemName(charIdentity.Character.Location)

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

	return charIdentity, nil
}
