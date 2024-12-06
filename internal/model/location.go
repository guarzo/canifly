package model

type Station struct {
	SystemID int64  `json:"system_id"`
	ID       int64  `json:"station_id"`
	Name     string `json:"station_name"`
}

type Structure struct {
	Name     string `json:"name"`
	OwnerID  int64  `json:"owner_id"`
	SystemID int64  `json:"solar_system_id"`
	TypeID   int64  `json:"type_id"`
}

type CharacterLocation struct {
	SolarSystemID int64 `json:"solar_system_id"`
	StructureID   int64 `json:"structure_id"`
}

type CloneLocation struct {
	HomeLocation struct {
		LocationID   int64  `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"home_location"`
	JumpClones []struct {
		Implants     []int  `json:"implants"`
		JumpCloneID  int64  `json:"jump_clone_id"`
		LocationID   int64  `json:"location_id"`
		LocationType string `json:"location_type"`
	} `json:"jump_clones"`
}
