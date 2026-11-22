package services

import (
	"strings"
	"testing"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func sp(s string) *string { return &s }
func ipt(n int) *int      { return &n }

var portionUnits = map[uint]entities.Unit{
	1: {ID: 1, Name: "g"},
	2: {ID: 2, Name: "placek", FormOne: sp("placek"), FormFew: sp("placki"), FormMany: sp("placków")},
}

func portionErrs(p dtos.UpsertPortionRequest, servings int, publish bool) []string {
	var out []string
	validatePortion(p, servings, publish, portionUnits, func(f, m string) { out = append(out, f+": "+m) })
	return out
}

func TestValidatePortion(t *testing.T) {
	pizza := dtos.UpsertPortionRequest{
		Mode: entities.PortionSized, PieceUnitId: up(2),
		Sizes: []dtos.PieceSizeDto{{DiameterCm: 25, DoughG: 190}, {DiameterCm: 30, DoughG: 255, IsBase: true}},
	}
	form := dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRound, DCm: ipt(24), Pieces: ipt(12)}}
	rect := dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRect, ACm: ipt(25), BCm: ipt(35)}}

	ok := []struct {
		name     string
		p        dtos.UpsertPortionRequest
		servings int
	}{
		{"servings", dtos.UpsertPortionRequest{Mode: entities.PortionServings}, 4},
		{"sized", pizza, 2},
		{"form round", form, 0}, // servings not needed for a form
		{"form rect", rect, 0},
	}
	for _, c := range ok {
		if e := portionErrs(c.p, c.servings, true); len(e) > 0 {
			t.Errorf("%s: unexpected %v", c.name, e)
		}
	}

	bad := []struct {
		name     string
		p        dtos.UpsertPortionRequest
		servings int
		publish  bool
		want     string
	}{
		{"unknown mode", dtos.UpsertPortionRequest{Mode: "x"}, 2, false, "portion.mode"},
		{"servings missing", dtos.UpsertPortionRequest{Mode: entities.PortionServings}, 0, true, "servings"},
		{"sized pieces missing", pizza, 0, true, "servings"},
		{"sized no base", dtos.UpsertPortionRequest{Mode: entities.PortionSized, PieceUnitId: up(2),
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 30, DoughG: 255}}}, 2, true, "portion.sizes"},
		{"sized no sizes", dtos.UpsertPortionRequest{Mode: entities.PortionSized, PieceUnitId: up(2)}, 2, true, "portion.sizes"},
		{"sized no unit", dtos.UpsertPortionRequest{Mode: entities.PortionSized,
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 30, DoughG: 255, IsBase: true}}}, 2, true, "portion.pieceUnitId"},
		{"unit without forms", dtos.UpsertPortionRequest{Mode: entities.PortionServings, PieceUnitId: up(1)}, 2, false, "portion.pieceUnitId"},
		{"two bases (draft too)", dtos.UpsertPortionRequest{Mode: entities.PortionServings,
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 25, DoughG: 1, IsBase: true}, {DiameterCm: 30, DoughG: 1, IsBase: true}}}, 2, false, "portion.sizes"},
		{"duplicate size", dtos.UpsertPortionRequest{Mode: entities.PortionSized,
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 30, DoughG: 1}, {DiameterCm: 30, DoughG: 2}}}, 2, false, "portion.sizes[1].diameterCm"},
		{"size out of range", dtos.UpsertPortionRequest{Mode: entities.PortionSized,
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 61, DoughG: 1}}}, 2, false, "portion.sizes[0].diameterCm"},
		{"no dough", dtos.UpsertPortionRequest{Mode: entities.PortionSized,
			Sizes: []dtos.PieceSizeDto{{DiameterCm: 30}}}, 2, false, "portion.sizes[0].doughG"},
		{"form missing", dtos.UpsertPortionRequest{Mode: entities.PortionForm}, 0, true, "portion.form.shape"},
		{"round without d", dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRound}}, 0, true, "portion.form.dCm"},
		{"rect half", dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRect, ACm: ipt(25)}}, 0, true, "portion.form.aCm"},
		{"form too small", dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRound, DCm: ipt(9)}}, 0, false, "portion.form.dCm"},
		{"pieces zero", dtos.UpsertPortionRequest{Mode: entities.PortionForm, Form: &dtos.FormDto{Shape: entities.FormRound, DCm: ipt(24), Pieces: ipt(0)}}, 0, false, "portion.form.pieces"},
	}
	for _, c := range bad {
		errs := portionErrs(c.p, c.servings, c.publish)
		found := false
		for _, e := range errs {
			if strings.HasPrefix(e, c.want+":") {
				found = true
			}
		}
		if !found {
			t.Errorf("%s: want error on %s, got %v", c.name, c.want, errs)
		}
	}
}

func TestApplyPortionKeepsOnlyChosenShapeDims(t *testing.T) {
	var e entities.Recipe
	applyPortion(&e, &dtos.UpsertPortionRequest{Mode: entities.PortionForm,
		Form: &dtos.FormDto{Shape: entities.FormRect, DCm: ipt(24), ACm: ipt(25), BCm: ipt(35), Pieces: ipt(20)}})
	if e.FormDCm != nil || *e.FormACm != 25 || *e.FormBCm != 35 || *e.FormPieces != 20 || *e.FormShape != entities.FormRect {
		t.Errorf("unexpected form fields: %+v", e)
	}
	var fresh entities.Recipe
	applyPortion(&fresh, nil)
	if fresh.PortionMode != entities.PortionServings {
		t.Errorf("nil portion on create should default to servings, got %q", fresh.PortionMode)
	}
}

func TestNutritionDivisor(t *testing.T) {
	cases := []struct {
		r    entities.Recipe
		want int
	}{
		{entities.Recipe{Servings: 4}, 4},
		{entities.Recipe{Servings: 4, PortionMode: entities.PortionServings}, 4},
		{entities.Recipe{Servings: 2, PortionMode: entities.PortionSized}, 2},
		{entities.Recipe{Servings: 6, PortionMode: entities.PortionForm, FormPieces: ipt(12)}, 12},
		{entities.Recipe{Servings: 6, PortionMode: entities.PortionForm}, 1},
	}
	for _, c := range cases {
		if got := c.r.NutritionDivisor(); got != c.want {
			t.Errorf("%+v: got %d, want %d", c.r, got, c.want)
		}
	}
}
