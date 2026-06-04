package dtos

type RecipeDto struct {
	ID             uint                  `json:"id"`
	Name           string                `json:"name"`
	Servings       int                   `json:"servings"`
	Description    *string               `json:"description,omitempty"`
	Tagline        *string               `json:"tagline,omitempty"`
	Category       *string               `json:"category,omitempty"`
	Region         *string               `json:"region,omitempty"`
	TimeMinutes    *int                  `json:"timeMinutes,omitempty"`
	Difficulty     *int                  `json:"difficulty,omitempty"`
	KcalPerServing *int                  `json:"kcalPerServing,omitempty"`
	Photo          *RecipePhotoDto       `json:"photo,omitempty"`
	Steps          []RecipeStepDto       `json:"steps"`
	Ingredients    []RecipeIngredientDto `json:"ingredients"`
}

type RecipeBaseDto struct {
	ID          uint            `json:"id"`
	Name        string          `json:"name"`
	Tagline     *string         `json:"tagline,omitempty"`
	Category    *string         `json:"category,omitempty"`
	Region      *string         `json:"region,omitempty"`
	TimeMinutes *int            `json:"timeMinutes,omitempty"`
	Difficulty  *int            `json:"difficulty,omitempty"`
	Photo       *RecipePhotoDto `json:"photo,omitempty"`
}

type UpsertRecipeRequest struct {
	Name           string  `json:"name"`
	Servings       int     `json:"servings"`
	Description    *string `json:"description"`
	Tagline        *string `json:"tagline"`
	Category       *string `json:"category"`
	Region         *string `json:"region"`
	TimeMinutes    *int    `json:"timeMinutes"`
	Difficulty     *int    `json:"difficulty"`
	KcalPerServing *int    `json:"kcalPerServing"`
	PhotoUrl       *string `json:"photoUrl"`

	Steps       []UpsertStepRequest       `json:"steps"`
	Ingredients []UpsertIngredientRequest `json:"ingredients"`
}

type UpsertStepRequest struct {
	StepNo  int     `json:"stepNo"`
	Section *string `json:"section"`
	Text    string  `json:"text"`
}

type UpsertIngredientRequest struct {
	SortOrder    int      `json:"sortOrder"`
	Section      *string  `json:"section"`
	IngredientId uint     `json:"ingredientId"`
	Amount       *float64 `json:"amount"`
	Unit         *string  `json:"unit"`
	AmountText   *string  `json:"amountText"`
	Note         *string  `json:"note"`
}
