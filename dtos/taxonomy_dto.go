package dtos

type TaxonomyValueDto struct {
	ID    uint    `json:"id"`
	Label string  `json:"label"`
	Color *string `json:"color,omitempty"`
	// Liczba przepisów z tą wartością — tylko w słownikach (GET /api/taxonomy/list).
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
	Forms []string `json:"forms,omitempty"` // [1, 2–4, 5+]; brak = nieodmienna
	NoQty bool     `json:"noQty"`
	// W ilu składnikach przepisów użyta.
	UsedIn int `json:"usedIn"`
}

type AllergenDto struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// DictionariesDto — wszystkie słowniki w jednym zapytaniu (GET /api/taxonomy).
type DictionariesDto struct {
	Axes            []TaxonomyAxisDto `json:"axes"`
	Units           []UnitDto         `json:"units"`
	Allergens       []AllergenDto     `json:"allergens"`
	IngredientTypes []TypeDto         `json:"ingredientTypes"`
}
