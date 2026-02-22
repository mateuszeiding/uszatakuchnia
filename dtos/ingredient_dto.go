package dtos

type IngredientDto struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	ParentID   *uint     `json:"parentId"`
	IsAllergen bool      `json:"isAllergen"`
	Type       TypeDto   `json:"type"`
	Image      *ImageDto `json:"image"`
}
