package dtos

type RecipeIngredientDto struct {
	SortOrder int     `json:"sortOrder"`
	Section   *string `json:"section,omitempty"`

	IngredientID   uint   `json:"ingredientId"`
	IngredientName string `json:"ingredientName"`

	Amount     *float64 `json:"amount,omitempty"`
	Unit       *string  `json:"unit,omitempty"`
	AmountText *string  `json:"amountText,omitempty"`
	Note       *string  `json:"note,omitempty"`
}
