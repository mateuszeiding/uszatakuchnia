package dtos

type TaxonomyValueDto struct {
	ID    uint    `json:"id"`
	Label string  `json:"label"`
	Color *string `json:"color,omitempty"`
	// Number of recipes with this value; only in the dictionaries (GET /api/taxonomy/list).
	RecipeCount *int `json:"recipeCount,omitempty"`
}

type TaxonomyAxisDto struct {
	ID         uint               `json:"id"`
	Code       string             `json:"code"`
	Label      string             `json:"label"`
	Multi      bool               `json:"multi"`
	Required   bool               `json:"required"`
	HasColor   bool               `json:"hasColor"`
	FilterMode string             `json:"filterMode"`
	Values     []TaxonomyValueDto `json:"values"`
}

type UnitDto struct {
	ID    uint     `json:"id"`
	Name  string   `json:"name"`
	Forms []string `json:"forms,omitempty"` // [1, 2–4, 5+]; empty = doesn't inflect
	NoQty bool     `json:"noQty"`
	// Number of recipe ingredient rows using it.
	UsedIn int `json:"usedIn"`
}

type AllergenDto struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// DictionariesDto: all dictionaries in one response.
type DictionariesDto struct {
	Axes             []TaxonomyAxisDto    `json:"axes"`
	Units            []UnitDto            `json:"units"`
	Allergens        []AllergenDto        `json:"allergens"`
	IngredientTypes  []TypeDto            `json:"ingredientTypes"`
	Equipment        []EquipmentDto       `json:"equipment"`
	ChallengeFormats []ChallengeFormatDto `json:"challengeFormats"`
}
