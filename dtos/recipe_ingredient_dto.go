package dtos

import "uszatakuchnia/db/entities"

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
	// Składnik pomijany na liście zakupów (np. woda).
	SkipShopping bool `json:"skipShopping"`
	// Ile gramów to 1 jednostka tej pozycji; nil = nie da się przeliczyć (lista zakupów sumuje różne jednostki przez gramy).
	GramsPerUnit *GramsRangeDto `json:"gramsPerUnit"`

	// Przepis, do którego prowadzi składnik; nil = brak. Szkic też jest zwracany (ze statusem) —
	// strona publiczna go nie pokazuje, formularz tak.
	LinkedRecipe *LinkedRecipeDto `json:"linkedRecipe"`

	// Wkład pozycji w makro na porcję (suma wkładów = RecipeDto.nutrition); nil = pozycja się nie liczy
	// (opcjonalna, „do smaku”, bez ilości). Front liczy z wkładów makro wybranego wariantu.
	Nutrition *entities.RecipeNutrition `json:"nutrition"`
	// Zamienniki w kolejności; [] = brak.
	Substitutes []RecipeSubstituteDto `json:"substitutes"`
}

// RecipeSubstituteDto — zamiennik pozycji z danymi do alergenów, listy zakupów i makro (bez kolejnego zapytania).
type RecipeSubstituteDto struct {
	IngredientID   uint                 `json:"ingredientId"`
	Name           string               `json:"name"`
	Mix            bool                 `json:"mix"` // „pół na pół” ze składnikiem głównym
	Allergens      []string             `json:"allergens"`
	IngredientType IngredientTypeRefDto `json:"ingredientType"`
	SkipShopping   bool                 `json:"skipShopping"`
	// Gramy na 1 jednostkę pozycji dla zamiennika; nil = nie da się przeliczyć.
	GramsPerUnit *GramsRangeDto `json:"gramsPerUnit"`
	// Składnik ma jakiekolwiek dane makro (inaczej „bez danych o wartościach odżywczych”).
	HasNutrition bool `json:"hasNutrition"`
	// Wkład wariantu w makro na porcję (mieszanka = połowa głównego + połowa zamiennika); nil = pozycja się nie liczy.
	Nutrition *entities.RecipeNutrition `json:"nutrition"`
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
