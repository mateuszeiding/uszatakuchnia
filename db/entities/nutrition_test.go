package entities

import (
	"reflect"
	"testing"
)

func f(v float64) *float64 { return &v }

func TestRecipeNutritionScanValue(t *testing.T) {
	in := RecipeNutrition{
		PerPortion: map[string]SumRange{"kcal": {Min: 420, Max: 451}, "fiber": {Min: 3.5, Max: 4}},
		Missing:    []string{"oregano"},
	}
	raw, err := in.Value()
	if err != nil {
		t.Fatal(err)
	}
	var out RecipeNutrition
	if err := out.Scan(raw.([]byte)); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(in, out) {
		t.Errorf("round-trip: got %+v, want %+v", out, in)
	}

	var empty RecipeNutrition
	if err := empty.Scan(nil); err != nil {
		t.Errorf("Scan(nil) error: %v", err)
	}
}

func TestIngredientMacroOf(t *testing.T) {
	ing := Ingredient{FatMin: f(2)}
	if got := ing.MacroOf("kcal"); got != nil {
		t.Errorf("kcal: got %+v, want nil", got)
	}
	got := ing.MacroOf("fat")
	if got == nil || got.Min != 2 || got.Max != nil {
		t.Errorf("fat: got %+v, want {2 nil}", got)
	}
	ing.FatMax = f(3)
	if got := ing.MacroOf("fat"); got == nil || got.Max == nil || *got.Max != 3 {
		t.Errorf("fat max: got %+v", got)
	}
}

func TestIngredientSetMacro(t *testing.T) {
	var ing Ingredient
	ing.SetMacro("protein", &MacroRange{Min: 20, Max: f(22)})
	if ing.ProteinMin == nil || *ing.ProteinMin != 20 || ing.ProteinMax == nil || *ing.ProteinMax != 22 {
		t.Errorf("set: got %v %v", ing.ProteinMin, ing.ProteinMax)
	}
	ing.SetMacro("protein", nil)
	if ing.ProteinMin != nil || ing.ProteinMax != nil {
		t.Errorf("clear: got %v %v", ing.ProteinMin, ing.ProteinMax)
	}
}
