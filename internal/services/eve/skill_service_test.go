package eve_test

import (
	"errors"
	"testing"
	"time"

	"github.com/guarzo/canifly/internal/model"
	"github.com/guarzo/canifly/internal/services/eve"
	"github.com/guarzo/canifly/internal/testutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSkillPlanFile_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	expectedData := []byte("skillplan file content")
	skillRepo.On("GetSkillPlanFile", "myplan").Return(expectedData, nil).Once()

	data, err := svc.GetSkillPlanFile("myplan")
	assert.NoError(t, err)
	assert.Equal(t, expectedData, data)

	skillRepo.AssertExpectations(t)
}

func TestGetSkillPlanFile_Error(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	skillRepo.On("GetSkillPlanFile", "unknown").Return([]byte(nil), errors.New("not found")).Once()

	data, err := svc.GetSkillPlanFile("unknown")
	assert.Error(t, err)
	assert.Nil(t, data)

	skillRepo.AssertExpectations(t)
}

func TestDeleteSkillPlan_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	skillRepo.On("DeleteSkillPlan", "myplan").Return(nil).Once()

	err := svc.DeleteSkillPlan("myplan")
	assert.NoError(t, err)
	skillRepo.AssertExpectations(t)
}

func TestDeleteSkillPlan_Error(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	skillRepo.On("DeleteSkillPlan", "badplan").Return(errors.New("delete error")).Once()

	err := svc.DeleteSkillPlan("badplan")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "delete error")

	skillRepo.AssertExpectations(t)
}

func TestParseAndSaveSkillPlan_Success(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	contents := `
Skill A 3
Skill B IV
Skill C 1
` // Skill B "IV" should parse as 4

	// After parsing:
	// Skill A -> 3
	// Skill B -> 4 (Roman numeral)
	// Skill C -> 1
	// Verify that these are passed to SaveSkillPlan
	expectedSkills := map[string]model.Skill{
		"Skill A": {Name: "Skill A", Level: 3},
		"Skill B": {Name: "Skill B", Level: 4},
		"Skill C": {Name: "Skill C", Level: 1},
	}

	skillRepo.On("SaveSkillPlan", "myplan", mock.MatchedBy(func(skills map[string]model.Skill) bool {
		if len(skills) != len(expectedSkills) {
			return false
		}
		for k, v := range expectedSkills {
			if sv, ok := skills[k]; !ok || sv.Level != v.Level {
				return false
			}
		}
		return true
	})).Return(nil).Once()

	err := svc.ParseAndSaveSkillPlan(contents, "myplan")
	assert.NoError(t, err)

	skillRepo.AssertExpectations(t)
}

func TestParseAndSaveSkillPlan_Empty(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	contents := `
InvalidLine
` // No valid lines with space/level
	// Should result in an empty skills map
	skillRepo.On("SaveSkillPlan", "emptyplan", mock.Anything).Return(errors.New("cannot save empty skillplan")).Once()

	err := svc.ParseAndSaveSkillPlan(contents, "emptyplan")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty skillplan")
	skillRepo.AssertExpectations(t)
}

func TestGetSkillPlans(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	plans := map[string]model.SkillPlan{"p1": {Name: "p1"}}
	skillRepo.On("GetSkillPlans").Return(plans).Once()

	result := svc.GetSkillPlans()
	assert.Equal(t, plans, result)

	skillRepo.AssertExpectations(t)
}

func TestGetSkillTypes(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	typesMap := map[string]model.SkillType{"skillX": {TypeID: "1001"}}
	skillRepo.On("GetSkillTypes").Return(typesMap).Once()

	result := svc.GetSkillTypes()
	assert.Equal(t, typesMap, result)

	skillRepo.AssertExpectations(t)
}

func TestGetSkillTypeByID(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	st := model.SkillType{TypeID: "2002", TypeName: "Afterburner"}
	skillRepo.On("GetSkillTypeByID", "2002").Return(st, true).Once()

	skill, found := svc.GetSkillTypeByID("2002")
	assert.True(t, found)
	assert.Equal(t, "Afterburner", skill.TypeName)

	skillRepo.AssertExpectations(t)
}

func TestGetSkillName_NotFound(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	// Skill not found
	skillRepo.On("GetSkillTypeByID", "9999").Return(model.SkillType{}, false).Once()

	name := svc.GetSkillName(9999)
	assert.Equal(t, "", name) // no skill found means empty name

	skillRepo.AssertExpectations(t)
}

func TestGetMatchingSkillPlans(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	// Define some skill plans
	skillPlans := map[string]model.SkillPlan{
		"PlanA": {
			Name: "PlanA",
			Skills: map[string]model.Skill{
				"Afterburner":         {Name: "Afterburner", Level: 3},
				"Small Hybrid Turret": {Name: "Small Hybrid Turret", Level: 2},
			},
		},
		"PlanB": {
			Name: "PlanB",
			Skills: map[string]model.Skill{
				"Small Hybrid Turret": {Name: "Small Hybrid Turret", Level: 5},
			},
		},
	}

	// Define skill types so we can map skill names to IDs and check character skills
	skillTypes := map[string]model.SkillType{
		"Afterburner":         {TypeID: "1001", TypeName: "Afterburner"},
		"Small Hybrid Turret": {TypeID: "1002", TypeName: "Small Hybrid Turret"},
	}

	// Define a character with some skills
	charSkillList := []model.SkillResponse{
		{SkillID: 1001, TrainedSkillLevel: 3}, // Afterburner level 3
		{SkillID: 1002, TrainedSkillLevel: 2}, // Small Hybrid Turret level 2
	}

	// Define a skill queue: let's say character is training Small Hybrid Turret to level 5
	finishDate := time.Now().Add(2 * time.Hour)
	charQueue := []model.SkillQueue{
		{
			SkillID:       1002,
			FinishedLevel: 5,
			FinishDate:    &finishDate,
		},
	}

	character := model.Character{
		UserInfoResponse:        model.UserInfoResponse{CharacterName: "MyChar"},
		CharacterSkillsResponse: model.CharacterSkillsResponse{Skills: charSkillList},
		SkillQueue:              charQueue,
	}

	account := model.Account{
		Name: "MyAccount",
		Characters: []model.CharacterIdentity{
			{Character: character},
		},
		Status: model.Alpha,
	}

	accounts := []model.Account{account}

	// Mock skillRepo to return the skill types as requested
	skillRepo.On("GetSkillTypeByID", "1001").Return(skillTypes["Afterburner"], true).Maybe()
	skillRepo.On("GetSkillTypeByID", "1002").Return(skillTypes["Small Hybrid Turret"], true).Maybe()

	// Call GetPlanAndConversionData
	updatedPlans, _ := svc.GetPlanAndConversionData(accounts, skillPlans, skillTypes)

	// Let's analyze what we expect:
	// PlanA: requires Afterburner (3) and Small Hybrid Turret (2)
	// Character has exactly these levels, so qualifies. Not pending.
	// PlanB: requires Small Hybrid Turret (5), character is training to 5 -> Pending.

	planAResult := updatedPlans["PlanA"]
	assert.Equal(t, "PlanA", planAResult.Name)
	assert.Contains(t, planAResult.QualifiedCharacters, "MyChar")
	assert.Len(t, planAResult.PendingCharacters, 0)
	assert.Empty(t, planAResult.MissingSkills["MyChar"])
	assert.Len(t, planAResult.Characters, 1)
	assert.Equal(t, "Qualified", planAResult.Characters[0].Status)

	planBResult := updatedPlans["PlanB"]
	assert.Equal(t, "PlanB", planBResult.Name)
	assert.Len(t, planBResult.QualifiedCharacters, 0)
	assert.Contains(t, planBResult.PendingCharacters, "MyChar")
	assert.Empty(t, planBResult.MissingSkills["MyChar"])
	assert.Len(t, planBResult.Characters, 1)
	assert.Equal(t, "Pending", planBResult.Characters[0].Status)

	skillRepo.AssertExpectations(t)
}

func TestGetMatchingSkillPlans_MissingSkill(t *testing.T) {
	logger := &testutil.MockLogger{}
	skillRepo := &testutil.MockSkillRepository{}
	svc := eve.NewSkillService(logger, skillRepo)

	// PlanC requires Afterburner level 5
	skillPlans := map[string]model.SkillPlan{
		"PlanC": {
			Name: "PlanC",
			Skills: map[string]model.Skill{
				"Afterburner": {Name: "Afterburner", Level: 5},
			},
		},
	}

	// Known skill type for Afterburner
	skillTypes := map[string]model.SkillType{
		"Afterburner": {TypeID: "1001", TypeName: "Afterburner"},
	}

	// The character only has Afterburner at level 1, needs 5
	charSkillList := []model.SkillResponse{
		{SkillID: 1001, TrainedSkillLevel: 1},
	}
	character := model.Character{
		UserInfoResponse:        model.UserInfoResponse{CharacterName: "MyChar2"},
		CharacterSkillsResponse: model.CharacterSkillsResponse{Skills: charSkillList},
		SkillQueue:              []model.SkillQueue{},
	}

	account := model.Account{
		Name: "MyAccount2",
		Characters: []model.CharacterIdentity{
			{Character: character},
		},
		Status: model.Alpha,
	}

	accounts := []model.Account{account}

	// Since we have Afterburner as a known skill, the service might query skillRepo.GetSkillTypeByID
	skillRepo.On("GetSkillTypeByID", "1001").Return(skillTypes["Afterburner"], true).Maybe()

	updatedPlans, _ := svc.GetPlanAndConversionData(accounts, skillPlans, skillTypes)

	planCResult := updatedPlans["PlanC"]
	assert.Equal(t, "PlanC", planCResult.Name)
	assert.Empty(t, planCResult.QualifiedCharacters)
	assert.Empty(t, planCResult.PendingCharacters)
	// MissingSkills should have an entry for MyChar2 and Afterburner
	assert.NotEmpty(t, planCResult.MissingSkills["MyChar2"])
	assert.Equal(t, "Not Qualified", planCResult.Characters[0].Status)

	skillRepo.AssertExpectations(t)
}

/*
If you needed to mock GetPlanAndConversionData in a test with a mock service, you could do something like:

mockSkillService := &MockSkillService{}
mockSkillService.On("GetPlanAndConversionData", mock.Anything, mock.Anything, mock.Anything).
    Return(expectedPlanStatusMap, expectedConversionsMap)

Then in your code:

plans, conv := mockSkillService.GetPlanAndConversionData(accounts, skillPlans, skillTypes)
assert.Equal(t, expectedPlanStatusMap, plans)
assert.Equal(t, expectedConversionsMap, conv)

*/
