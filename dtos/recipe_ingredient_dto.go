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

	// Dział w sklepie (lista zakupów).
	IngredientType IngredientTypeRefDto `json:"ingredientType"`
	// Ile gramów to 1 jednostka tej pozycji; nil = nie da się przeliczyć (lista zakupów sumuje różne jednostki przez gramy).
	GramsPerUnit *GramsRangeDto `json:"gramsPerUnit"`

	// Przepis, do którego prowadzi składnik; nil = brak. Szkic też jest zwracany (ze statusem) —
	// strona publiczna go nie pokazuje, formularz tak.
	LinkedRecipe *LinkedRecipeDto `json:"linkedRecipe"`
}

type LinkedRecipeDto struct {
	ID     uint    `json:"id"`
	Title  string  `json:"title"`
	Status string  `json:"status"`
	Image  *string `json:"image"` // adres zdjęcia przepisu (miniatura); nil = brak zdjęcia
}

type IngredientTypeRefDto struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

// GramsRangeDto — gramy na jednostkę; Max = Min dla pojedynczej wartości.
type GramsRangeDto struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
