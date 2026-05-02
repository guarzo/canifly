package model

// CharacterUpdate is a partial update for a character. Nil pointers are no-ops.
// JSON keys match the legacy map-based payload accepted by PATCH /api/characters/:id.
type CharacterUpdate struct {
	Role *string `json:"Role,omitempty"`
	MCT  *bool   `json:"MCT,omitempty"`
}
