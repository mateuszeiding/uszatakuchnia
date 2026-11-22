package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

var NutrientKeys = []string{"kcal", "protein", "fat", "carbs", "fiber"}

// MacroRange: per 100 g, or grams per unit. Max nil = single value.
type MacroRange struct {
	Min float64  `json:"min"`
	Max *float64 `json:"max"`
}

func (r MacroRange) Upper() float64 {
	if r.Max == nil {
		return r.Min
	}
	return *r.Max
}

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

// Value returns a string, not []byte: pgx sends a []byte parameter as bytea, which a
// jsonb column rejects (see db.Open).
func (n RecipeNutrition) Value() (driver.Value, error) {
	b, err := json.Marshal(n)
	return string(b), err
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
