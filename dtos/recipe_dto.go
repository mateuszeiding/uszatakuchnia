package dtos

import "uszatakuchnia/db/entities"

type RecipeTaxonomyDto map[string][]TaxonomyValueDto

type RecipeDto struct {
	ID          uint              `json:"id"`
	Name        string            `json:"name"`
	Servings    int               `json:"servings"`
	Portion     PortionDto        `json:"portion"`
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
	// Glossary terms marked in the steps and hints (missing ones skipped: the page shows the bare word).
	Terms []GlossaryTermDto `json:"terms"`
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
	Servings   int               `json:"servings"`
	Portion    PortionDto        `json:"portion"`
	// Kcal per serving / piece of the base size / piece of the form (whole form without pieces); null = no data.
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
	// nil = unchanged (older client, content/ scripts); on create nil = "servings".
	Portion *UpsertPortionRequest `json:"portion"`

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
	Optional bool     `json:"optional"`
	Hints    []string `json:"hints"`
}

type UpsertIngredientRequest struct {
	SortOrder int `json:"sortOrder"`
	// Stable row key ([a-z0-9-], up to 32); empty = the server picks one. Keep it on a round trip, or the
	// step refs {{ing:<key>}} to this row break.
	Key          string  `json:"key"`
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
