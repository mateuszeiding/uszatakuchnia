package dtos

import "uszatakuchnia/db/entities"

// RecipeTaxonomyDto: classification values grouped by axis code (dish, base, cuisine, diet, practical).
type RecipeTaxonomyDto map[string][]TaxonomyValueDto

type RecipeDto struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Servings    int               `json:"servings"`
	Description *string           `json:"description,omitempty"`
	Tagline     *string           `json:"tagline,omitempty"`
	Tip         *string           `json:"tip,omitempty"`
	Cautions    []string          `json:"cautions"`
	TimeMin     int               `json:"timeMin"`
	TimeMax     *int              `json:"timeMax,omitempty"`
	Difficulty  *int              `json:"difficulty,omitempty"`
	Status      string            `json:"status"`
	NeedsPrep   bool              `json:"needsPrep"`
	Taxonomy    RecipeTaxonomyDto `json:"taxonomy"`
	Allergens   []string          `json:"allergens"`
	// Per-serving macros computed from the ingredients; null = no data.
	Nutrition   *entities.RecipeNutrition `json:"nutrition"`
	Photo       *RecipePhotoDto           `json:"photo,omitempty"`
	Steps       []RecipeStepDto           `json:"steps"`
	Ingredients []RecipeIngredientDto     `json:"ingredients"`
	Equipment   []RecipeEquipmentDto      `json:"equipment"`
	// Ingredient rows of other recipes that link to this one (content of the delete modal).
	LinkedFrom []RecipeLinkSourceDto `json:"linkedFrom"`
}

type RecipeLinkSourceDto struct {
	RecipeID       uint   `json:"recipeId"`
	RecipeName     string `json:"recipeName"`
	IngredientName string `json:"ingredientName"`
}

type RecipeBaseDto struct {
	ID         uint              `json:"id"`
	Name       string            `json:"name"`
	Tagline    *string           `json:"tagline,omitempty"`
	TimeMin    int               `json:"timeMin"`
	TimeMax    *int              `json:"timeMax,omitempty"`
	Difficulty *int              `json:"difficulty,omitempty"`
	Status     string            `json:"status"`
	NeedsPrep  bool              `json:"needsPrep"`
	Taxonomy   RecipeTaxonomyDto `json:"taxonomy"`
	Photo      *RecipePhotoDto   `json:"photo,omitempty"`
	// Kcal per serving (range); null = no data.
	Kcal *entities.SumRange `json:"kcal"`
	// Ingredient names: the list search matches them too.
	IngredientNames []string `json:"ingredientNames"`
	// Substitute names (deduplicated): the list search matches them too ("pancetta" → Ragù).
	SubstituteNames []string `json:"substituteNames"`
}

type UpsertRecipeRequest struct {
	Name        string  `json:"name"`
	Servings    int     `json:"servings"`
	Description *string `json:"description"`
	Tagline     *string `json:"tagline"`
	Tip         *string `json:"tip"`
	TimeMin     int     `json:"timeMin"` // minutes; 0 = not given (drafts only)
	TimeMax     *int    `json:"timeMax"`
	Difficulty  *int    `json:"difficulty"`
	PhotoUrl    *string `json:"photoUrl"`
	Status      string  `json:"status"`
	// Notes; nil = unchanged (older client), on create nil = [].
	Cautions *[]string `json:"cautions"`

	TaxonomyValueIDs []uint `json:"taxonomyValueIds"`

	Steps       []UpsertStepRequest       `json:"steps"`
	Ingredients []UpsertIngredientRequest `json:"ingredients"`
	// nil = unchanged (older client); [] = recipe without equipment.
	Equipment *[]UpsertEquipmentRequest `json:"equipment"`
}

type UpsertStepRequest struct {
	StepNo   int      `json:"stepNo"`
	Section  *string  `json:"section"`
	Text     string   `json:"text"`
	Prep     bool     `json:"prep"`
	PrepLead *string  `json:"prepLead"`
	Hints    []string `json:"hints"`
}

type UpsertIngredientRequest struct {
	SortOrder    int     `json:"sortOrder"`
	Section      *string `json:"section"`
	IngredientId uint    `json:"ingredientId"`
	Quantity     *string `json:"quantity"`
	UnitId       *uint   `json:"unitId"`
	Note         *string `json:"note"`
	Optional     bool    `json:"optional"`
	// Link to the recipe for this ingredient; nil = none.
	LinkedRecipeId *uint `json:"linkedRecipeId"`
	// Substitutes in order; nil / [] = none.
	Substitutes []UpsertSubstituteRequest `json:"substitutes"`
}

type UpsertSubstituteRequest struct {
	IngredientId uint `json:"ingredientId"`
	Mix          bool `json:"mix"`
}

// FieldError is a validation error attached to a form field (e.g. "ingredients[2].quantity").
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}
