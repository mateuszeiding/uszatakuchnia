package entities

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
)

// NutrientKeys — klucze makro w stałej kolejności (kcal, białko, tłuszcz, węglowodany, błonnik).
var NutrientKeys = []string{"kcal", "protein", "fat", "carbs", "fiber"}

// MacroRange — wartość orientacyjna (na 100 g albo gramy jednostki). Max nil = jedna wartość.
type MacroRange struct {
	Min float64  `json:"min"`
	Max *float64 `json:"max"`
}

// Upper — górna granica; dla pojedynczej wartości równa Min.
func (r MacroRange) Upper() float64 {
	if r.Max == nil {
		return r.Min
	}
	return *r.Max
}

// SumRange — wartość wyliczona dla przepisu (zawsze para).
type SumRange struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}

// RecipeNutrition — cache makro przepisu w recipes.nutrition (jsonb).
type RecipeNutrition struct {
	PerPortion map[string]SumRange `json:"perPortion"`
	Missing    []string            `json:"missing"`
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
