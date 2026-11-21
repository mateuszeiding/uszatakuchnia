package services

import (
	"reflect"
	"strings"
	"testing"

	"uszatakuchnia/dtos"
)

func TestParseStepRefs(t *testing.T) {
	got := ParseStepRefs("Dodaj mąkę {{ing:maka}} i połowę masła {{ing:i3|1/2}}, potem {{term:deglasowanie}}. {{bez}} {x}")
	want := []StepRef{{"ing", "maka", ""}, {"ing", "i3", "1/2"}, {"term", "deglasowanie", ""}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %#v", got)
	}
	if got := ParseStepRefs("Zwykły tekst bez odwołań."); len(got) != 0 {
		t.Errorf("plain text: %#v", got)
	}
}

func collect(fn func(add func(f, m string))) []dtos.FieldError {
	errs := []dtos.FieldError{}
	fn(func(f, m string) { errs = append(errs, dtos.FieldError{Field: f, Message: m}) })
	return errs
}

func TestValidateStepRefs(t *testing.T) {
	req := dtos.UpsertRecipeRequest{
		Ingredients: []dtos.UpsertIngredientRequest{{Key: "maka"}, {Key: "i2"}},
		Steps: []dtos.UpsertStepRequest{
			{Text: "Wsyp mąkę {{ing:maka}} i resztę {{ing:i2|rest}}."},
			{Text: "Bez odwołań."},
		},
	}
	for _, publish := range []bool{false, true} {
		if errs := collect(func(add func(f, m string)) { validateStepRefs(req, publish, add) }); len(errs) != 0 {
			t.Errorf("valid refs (publish=%v): %v", publish, errs)
		}
	}

	missing := req
	missing.Steps = []dtos.UpsertStepRequest{{Text: "x"}, {Text: "Dodaj {{ing:usuniety}}."}}
	if errs := collect(func(add func(f, m string)) { validateStepRefs(missing, false, add) }); len(errs) != 0 {
		t.Errorf("draft with a missing row must save: %v", errs)
	}
	errs := collect(func(add func(f, m string)) { validateStepRefs(missing, true, add) })
	if len(errs) != 1 || errs[0].Field != "steps[1].text" || !strings.HasPrefix(errs[0].Message, "Krok 2:") {
		t.Errorf("publish with a missing row: %#v", errs)
	}

	bad := req
	bad.Steps = []dtos.UpsertStepRequest{{Text: "{{ingr:maka}} {{ing:maka|5/8}} {{ing:}}"}}
	if errs := collect(func(add func(f, m string)) { validateStepRefs(bad, false, add) }); len(errs) != 3 {
		t.Errorf("malformed refs are errors also in a draft: %#v", errs)
	}
}

func TestValidateIngredientKeys(t *testing.T) {
	req := dtos.UpsertRecipeRequest{Ingredients: []dtos.UpsertIngredientRequest{
		{Key: "maka-zaczyn"}, {Key: ""}, {Key: "Mąka"}, {Key: "maka-zaczyn"}, {Key: strings.Repeat("a", 33)},
	}}
	errs := collect(func(add func(f, m string)) { validateIngredientKeys(req, add) })
	fields := []string{}
	for _, e := range errs {
		fields = append(fields, e.Field)
	}
	want := []string{"ingredients[2].key", "ingredients[3].key", "ingredients[4].key"}
	if !reflect.DeepEqual(fields, want) {
		t.Errorf("got %v, want %v", fields, want)
	}
}

func TestIngredientKeys(t *testing.T) {
	keys := ingredientKeys([]dtos.UpsertIngredientRequest{{Key: "maka"}, {}, {Key: " i2 "}, {}})
	if keys[0] != "maka" || keys[2] != "i2" {
		t.Errorf("given keys must stay: %v", keys)
	}
	seen := map[string]bool{}
	for _, k := range keys {
		if !rowKeyRe.MatchString(k) || seen[k] {
			t.Errorf("bad or duplicate key %q in %v", k, keys)
		}
		seen[k] = true
	}
}
