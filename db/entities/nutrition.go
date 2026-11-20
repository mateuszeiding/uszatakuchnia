package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NutrientKeys: macro keys in a fixed order (kcal, protein, fat, carbs, fiber).
var NutrientKeys = []string{"kcal", "protein", "fat", "carbs", "fiber"}

// MacroRange is an approximate value (per 100 g, or grams per unit). Max nil = single value.
type MacroRange struct {
	Min float64  `json:"min"`
	Max *float64 `json:"max"`
}

// Upper returns the upper bound; equals Min for a single value.
func (r MacroRange) Upper() float64 {
	if r.Max == nil {
		return r.Min
	}
	return *r.Max
}

// SumRange is a value computed for a recipe (always a pair).
type SumRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// RecipeNutrition is the recipe macro cache in recipes.nutrition (jsonb).
type RecipeNutrition struct {
	PerPortion map[string]SumRange `json:"perPortion"`
	Missing    []string            `json:"missing"` // skipped entirely (no macros or no conversion)
	// Counted in the sums but missing a value: nutrient → ingredient names ("bez: tofu naturalne").
	Incomplete map[string][]string `json:"incomplete"`
}

func (n RecipeNutrition) Value() (driver.Value, error) {
	return json.Marshal(n)
}

func (n *RecipeNutrition) Scan(src any) error {
	switch v := src.(type) {
	case nil:
		*n = RecipeNutrition{}
		return nil
	case []byte:
		return json.Unmarshal(v, n)
	case string:
		return json.Unmarshal([]byte(v), n)
	default:
		return fmt.Errorf("RecipeNutrition: nieobsługiwany typ %T", src)
	}
}
