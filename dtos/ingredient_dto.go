package dtos

type IngredientDto struct {
	ID         uint      `json:"id"`
	Name       string    `json:"name"`
	ParentID   *uint     `json:"parentId"`
	Aromas     []TypeDto `json:"aromas"`
	IsAllergen bool      `json:"isAllergen"`
	Type       TypeDto   `json:"type"`
	Tastes     []TypeDto `json:"tastes"`
	Image      *ImageDto `json:"image"`
}
