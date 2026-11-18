package services

import (
	"testing"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// Suma wkładów pozycji = Compute (front składa makro wariantu z wkładów; bez wyboru wynik musi być dzisiejszy).
func TestLineNutritionSumsToCompute(t *testing.T) {
	ix := testIndex()
	lines := []NutritionLine{
		{IngredientID: 2, Qty: Quantity{Min: 200, Max: 250}, UnitID: up(unitG), UnitName: "g"},
		{IngredientID: 4, Qty: Quantity{Min: 2, Max: 2}, UnitID: up(unitSzt), UnitName: "sztuka"},
		{IngredientID: 5, Qty: Quantity{Min: 2, Max: 2}, UnitID: up(unitLyzka), UnitName: "łyżka"},
		{IngredientID: 6, Qty: Quantity{Min: 1, Max: 1}, UnitID: up(unitG), UnitName: "g"},                         // bez danych → Missing
		{IngredientID: 5, Qty: Quantity{Min: 1, Max: 1}, UnitID: up(unitLyzka), UnitName: "łyżka", Optional: true}, // nie liczy się
	}
	want := ix.Compute(lines, 4)

	sum := map[string]entities.SumRange{}
	missing := []string{}
	for i, l := range lines {
		n := ix.LineNutrition([]NutritionLine{l}, 4)
		if i == 4 {
			if n != nil {
				t.Fatalf("pozycja opcjonalna ma wkład: %+v", n)
			}
			continue
		}
		missing = append(missing, n.Missing...)
		for k, v := range n.PerPortion {
			s := sum[k]
			s.Min += v.Min
			s.Max += v.Max
			sum[k] = s
		}
	}
	for k, w := range want.PerPortion {
		if !near(sum[k].Min, w.Min) || !near(sum[k].Max, w.Max) {
			t.Errorf("%s: suma wkładów %+v, Compute %+v", k, sum[k], w)
		}
	}
	if len(sum) != len(want.PerPortion) {
		t.Errorf("klucze: %v vs %v", sum, want.PerPortion)
	}
	if len(missing) != 1 || missing[0] != "oregano" || len(want.Missing) != 1 {
		t.Errorf("missing: %v (Compute %v)", missing, want.Missing)
	}
}

// Mieszanka = połowa głównego + połowa zamiennika; zamiennik bez danych → Missing i hasNutrition false.
func TestFillLineNutritionSubstitutes(t *testing.T) {
	ix := testIndex()
	qty := "200"
	recipe := entities.Recipe{Servings: 2, Ingredients: []entities.RecipeIngredient{{
		IngredientID: 2, Quantity: &qty, UnitID: up(unitG), Unit: &entities.Unit{ID: unitG, Name: "g"},
	}}}
	items := []dtos.RecipeIngredientDto{{
		IngredientID: 2,
		Substitutes: []dtos.RecipeSubstituteDto{
			{IngredientID: 3},            // fusilli zamiast
			{IngredientID: 3, Mix: true}, // pół na pół
			{IngredientID: 6},            // oregano — bez danych
		},
	}}
	FillLineNutrition(items, recipe, ix)

	main := items[0].Nutrition.PerPortion["kcal"]
	alt := items[0].Substitutes[0].Nutrition.PerPortion["kcal"]
	mix := items[0].Substitutes[1].Nutrition.PerPortion["kcal"]
	if !near(main.Min, 200*350/100.0/2) || !near(alt.Min, 200*330/100.0/2) {
		t.Fatalf("główny %+v, zamiennik %+v", main, alt)
	}
	if !near(mix.Min, (main.Min+alt.Min)/2) || !near(mix.Max, (main.Max+alt.Max)/2) {
		t.Errorf("mieszanka %+v ≠ połowa %+v i %+v", mix, main, alt)
	}
	none := items[0].Substitutes[2]
	if none.HasNutrition || len(none.Nutrition.PerPortion) != 0 || len(none.Nutrition.Missing) != 1 {
		t.Errorf("zamiennik bez danych: %+v / %+v", none, none.Nutrition)
	}
	if items[0].Substitutes[0].GramsPerUnit == nil || items[0].Substitutes[0].GramsPerUnit.Min != 1 {
		t.Errorf("gramsPerUnit zamiennika: %+v", items[0].Substitutes[0].GramsPerUnit)
	}
}
