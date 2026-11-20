package dtos

import "uszatakuchnia/db/entities"

type IngredientDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parentId"`
	Allergens []string  `json:"allergens"` // codes from entities.Allergens
	UsedIn    int       `json:"usedIn"`    // number of recipes (list only)
	Type      TypeDto   `json:"type"`
	Image     *ImageDto `json:"image"`
	// Never added to the shopping list (e.g. water).
	SkipShopping bool `json:"skipShopping"`

	// Macros per 100 g: own and resolved (own or inherited from sub-ingredients), keyed by entities.NutrientKeys.
	Nutrition         map[string]entities.MacroRange `json:"nutrition"`
	ResolvedNutrition map[string]entities.MacroRange `json:"resolvedNutrition"`
	// Unit conversions: own + inherited from ancestors (inherited) for units without an own one.
	Conversions []ConversionDto `json:"conversions"`
}

type ConversionDto struct {
	UnitID    uint     `json:"unitId"`
	UnitName  string   `json:"unitName"`
	GramsMin  float64  `json:"gramsMin"`
	GramsMax  *float64 `json:"gramsMax"`
	Inherited bool     `json:"inherited"`
}

// IngredientRequest creates or edits an ingredient (admin panel). A missing key in Nutrition = no own data.
type IngredientRequest struct {
	Name         string                         `json:"name"`
	TypeCode     string                         `json:"typeCode"`
	Allergens    []string                       `json:"allergens"`
	ParentID     *uint                          `json:"parentId"`
	SkipShopping bool                           `json:"skipShopping"`
	Nutrition    map[string]entities.MacroRange `json:"nutrition"`
	Conversions  []ConversionRequest            `json:"conversions"`
}

type ConversionRequest struct {
	UnitID   uint     `json:"unitId"`
	GramsMin float64  `json:"gramsMin"`
	GramsMax *float64 `json:"gramsMax"`
}
