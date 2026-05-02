package model_test

import (
	"encoding/json"
	"testing"

	"github.com/guarzo/canifly/internal/model"
)

func TestCharacterUpdate_DecodesPartial(t *testing.T) {
	raw := []byte(`{"Role":"main","MCT":true}`)
	var u model.CharacterUpdate
	if err := json.Unmarshal(raw, &u); err != nil {
		t.Fatal(err)
	}
	if u.Role == nil || *u.Role != "main" {
		t.Fatalf("Role=%v, want 'main'", u.Role)
	}
	if u.MCT == nil || *u.MCT != true {
		t.Fatalf("MCT=%v, want true", u.MCT)
	}
}

func TestCharacterUpdate_OmitsAbsent(t *testing.T) {
	var u model.CharacterUpdate
	if err := json.Unmarshal([]byte(`{"Role":"alt"}`), &u); err != nil {
		t.Fatal(err)
	}
	if u.MCT != nil {
		t.Fatalf("MCT=%v, want nil", u.MCT)
	}
}
