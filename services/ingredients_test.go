package services

import (
	"testing"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func fieldMsgs(errs []dtos.FieldError) map[string]string {
	out := map[string]string{}
	for _, e := range errs {
		out[e.Field] = e.Message
	}
	return out
}

func validReq() dtos.IngredientRequest {
	return dtos.IngredientRequest{
		Name:     "Fusilli",
		TypeCode: "GRAIN",
		ParentID: up(1),
		Nutrition: map[string]entities.MacroRange{
			"kcal":    rng(330, 345),
			"protein": rng(22),
		},
		Conversions: []dtos.ConversionRequest{{UnitID: unitLyzka, GramsMin: 12, GramsMax: fp(15)}},
	}
}

func TestValidateIngredientOK(t *testing.T) {
	if errs := ValidateIngredientRequest(validReq(), testIndex(), 3); len(errs) != 0 {
		t.Errorf("errs = %v", errs)
	}
}

func TestValidateIngredientMaxBelowMin(t *testing.T) {
	req := validReq()
	req.Nutrition["kcal"] = rng(345, 330)
	if got := fieldMsgs(ValidateIngredientRequest(req, testIndex(), 3))["nutrition.kcal.max"]; got != "Max nie może być mniejsze niż min" {
		t.Errorf("msg = %q", got)
	}
}

func TestValidateIngredientNegative(t *testing.T) {
	req := validReq()
	req.Nutrition["fat"] = rng(-1)
	req.Conversions[0].GramsMin = -2
	m := fieldMsgs(ValidateIngredientRequest(req, testIndex(), 3))
	if m["nutrition.fat.min"] == "" || m["conversions[0].gramsMin"] == "" {
		t.Errorf("errs = %v", m)
	}
}

func TestValidateIngredientUnknownNutrient(t *testing.T) {
	req := validReq()
	req.Nutrition["sugar"] = rng(5)
	if m := fieldMsgs(ValidateIngredientRequest(req, testIndex(), 3)); m["nutrition.sugar"] == "" {
		t.Errorf("errs = %v", m)
	}
}

func TestValidateIngredientDuplicateUnit(t *testing.T) {
	req := validReq()
	req.Conversions = append(req.Conversions, dtos.ConversionRequest{UnitID: unitLyzka, GramsMin: 10})
	if got := fieldMsgs(ValidateIngredientRequest(req, testIndex(), 3))["conversions[1].unitId"]; got != "Ta jednostka ma już przelicznik" {
		t.Errorf("msg = %q", got)
	}
}

func TestValidateIngredientFixedUnits(t *testing.T) {
	ix := testIndex()
	for _, u := range []uint{unitG, unitKg, unitDoSmaku, 999} {
		req := validReq()
		req.Conversions = []dtos.ConversionRequest{{UnitID: u, GramsMin: 1}}
		if m := fieldMsgs(ValidateIngredientRequest(req, ix, 3)); m["conversions[0].unitId"] == "" {
			t.Errorf("unit %d: brak błędu", u)
		}
	}
}

func TestValidateIngredientParentCycle(t *testing.T) {
	ix := testIndex()
	ix.Children[3] = []uint{7}
	ix.Parent[7] = up(3)
	ix.Names[7] = "fusilli mini"
	for _, parent := range []uint{3, 7} { // sam siebie, wnuk
		req := validReq()
		req.ParentID = up(parent)
		if got := fieldMsgs(ValidateIngredientRequest(req, ix, 3))["parentId"]; got != "Składnik nie może być własnym podskładnikiem" {
			t.Errorf("parent %d: msg = %q", parent, got)
		}
	}
	req := validReq()
	req.ParentID = up(999)
	if m := fieldMsgs(ValidateIngredientRequest(req, ix, 3)); m["parentId"] == "" {
		t.Errorf("nieistniejący rodzic: brak błędu")
	}
}

func TestFillIngredientNutritionInherits(t *testing.T) {
	ix := testIndex()
	ix.Conversions[2] = map[uint]entities.MacroRange{unitSzt: rng(5)}
	d := dtos.IngredientDto{ID: 2}
	FillIngredientNutrition(&d, ix)
	if len(d.Nutrition) != 2 || !eqRange(d.ResolvedNutrition["kcal"], 350, fp(360)) {
		t.Errorf("nutrition = %+v / %+v", d.Nutrition, d.ResolvedNutrition)
	}
	want := []dtos.ConversionDto{
		{UnitID: unitSzt, UnitName: "sztuka", GramsMin: 5, Inherited: false},
		{UnitID: unitLyzka, UnitName: "łyżka", GramsMin: 15, Inherited: true},
	}
	if len(d.Conversions) != 2 {
		t.Fatalf("conversions = %+v", d.Conversions)
	}
	for i, w := range want {
		g := d.Conversions[i]
		if g.UnitID != w.UnitID || g.UnitName != w.UnitName || g.GramsMin != w.GramsMin || g.Inherited != w.Inherited {
			t.Errorf("conversions[%d] = %+v, want %+v", i, g, w)
		}
	}
}
