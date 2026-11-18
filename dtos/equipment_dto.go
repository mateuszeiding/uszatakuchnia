package dtos

// RecipeEquipmentDto — pozycja sprzętu w przepisie (kolejność sort_order).
type RecipeEquipmentDto struct {
	EquipmentID uint    `json:"equipmentId"`
	Name        string  `json:"name"`
	Spec        *string `json:"spec"`
	Optional    bool    `json:"optional"`
}

// EquipmentDto — pozycja słownika sprzętu; UsedIn = w ilu przepisach.
type EquipmentDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	UsedIn int    `json:"usedIn"`
}

type UpsertEquipmentRequest struct {
	EquipmentId uint    `json:"equipmentId"`
	Spec        *string `json:"spec"`
	Optional    bool    `json:"optional"`
}
