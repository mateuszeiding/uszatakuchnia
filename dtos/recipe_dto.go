package dtos

type RecipeDto struct {
	ID          uint                  `json:"id"`
	Name        string                `json:"name"`
	Servings    int                   `json:"servings"`
	Description *string               `json:"description,omitempty"`
	Photo       *RecipePhotoDto       `json:"photo,omitempty"`
	Steps       []RecipeStepDto       `json:"steps"`
	Ingredients []RecipeIngredientDto `json:"ingredients"`
}

type RecipeBaseDto struct {
	ID    uint            `json:"id"`
	Name  string          `json:"name"`
	Photo *RecipePhotoDto `json:"photo,omitempty"`
}
