package dtos

// RecipeTaxonomyDto — wartości klasyfikacji pogrupowane po kodzie osi (dish, base, cuisine, diet, practical).
type RecipeTaxonomyDto map[string][]TaxonomyValueDto

type RecipeDto struct {
	ID             uint                  `json:"id"`
	Name           string                `json:"name"`
	Servings       int                   `json:"servings"`
	Description    *string               `json:"description,omitempty"`
	Tagline        *string               `json:"tagline,omitempty"`
	Tip            *string               `json:"tip,omitempty"`
	TimeMinutes    *int                  `json:"timeMinutes,omitempty"`
	Difficulty     *int                  `json:"difficulty,omitempty"`
	KcalPerServing *int                  `json:"kcalPerServing,omitempty"`
	Status         string                `json:"status"`
	NeedsPrep      bool                  `json:"needsPrep"`
	Taxonomy       RecipeTaxonomyDto     `json:"taxonomy"`
	Allergens      []string              `json:"allergens"`
	Photo          *RecipePhotoDto       `json:"photo,omitempty"`
	Steps          []RecipeStepDto       `json:"steps"`
	Ingredients    []RecipeIngredientDto `json:"ingredients"`
}

type RecipeBaseDto struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Tagline     *string           `json:"tagline,omitempty"`
	TimeMinutes *int              `json:"timeMinutes,omitempty"`
	Difficulty  *int              `json:"difficulty,omitempty"`
	Status      string            `json:"status"`
	NeedsPrep   bool              `json:"needsPrep"`
	Taxonomy    RecipeTaxonomyDto `json:"taxonomy"`
	Photo       *RecipePhotoDto   `json:"photo,omitempty"`
	// Nazwy składników — wyszukiwarka listy szuka też po nich.
	IngredientNames []string `json:"ingredientNames"`
}

type UpsertRecipeRequest struct {
	Name           string  `json:"name"`
	Servings       int     `json:"servings"`
	Description    *string `json:"description"`
	Tagline        *string `json:"tagline"`
	Tip            *string `json:"tip"`
	TimeMinutes    *int    `json:"timeMinutes"`
	Difficulty     *int    `json:"difficulty"`
	KcalPerServing *int    `json:"kcalPerServing"`
	PhotoUrl       *string `json:"photoUrl"`
	Status         string  `json:"status"`

	TaxonomyValueIDs []uint `json:"taxonomyValueIds"`

	Steps       []UpsertStepRequest       `json:"steps"`
	Ingredients []UpsertIngredientRequest `json:"ingredients"`
}

type UpsertStepRequest struct {
	StepNo   int     `json:"stepNo"`
	Section  *string `json:"section"`
	Text     string  `json:"text"`
	Prep     bool    `json:"prep"`
	PrepLead *string `json:"prepLead"`
}

type UpsertIngredientRequest struct {
	SortOrder    int     `json:"sortOrder"`
	Section      *string `json:"section"`
	IngredientId uint    `json:"ingredientId"`
	Quantity     *string `json:"quantity"`
	UnitId       *uint   `json:"unitId"`
	Note         *string `json:"note"`
	Optional     bool    `json:"optional"`
}

// FieldError — błąd walidacji przypięty do pola formularza (np. "ingredients[2].quantity").
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
