package dtos

type RecipeIngredientDto struct {
	SortOrder int     `json:"sortOrder"`
	Section   *string `json:"section,omitempty"`

	IngredientID   uint   `json:"ingredientId"`
	IngredientName string `json:"ingredientName"`

	Quantity *string `json:"quantity,omitempty"`
	UnitID   *uint   `json:"unitId,omitempty"`
	Unit     *string `json:"unit,omitempty"` // nazwa jednostki (forma słownikowa); odmiana po stronie frontu
	Note     *string `json:"note,omitempty"`
	Optional bool    `json:"optional"`

	Allergens []string `json:"allergens"` // kody alergenów składnika (do listy „Alergeny” w przepisie)
}
