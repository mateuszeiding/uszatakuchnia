package dtos

// RecipeEquipmentDto is an equipment item in a recipe (ordered by sort_order).
type RecipeEquipmentDto struct {
	EquipmentID uint    `json:"equipmentId"`
	Name        string  `json:"name"`
	Spec        *string `json:"spec"`
	Optional    bool    `json:"optional"`
}

// EquipmentDto is an equipment dictionary item; UsedIn = number of recipes using it.
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
