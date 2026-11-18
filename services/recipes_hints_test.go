package services

import (
	"reflect"
	"strings"
	"testing"

	"uszatakuchnia/dtos"
)

func TestCleanHints(t *testing.T) {
	got := CleanHints([]string{"  x ", "", "   ", "y"})
	if !reflect.DeepEqual(got, []string{"x", "y"}) {
		t.Errorf("got %#v", got)
	}
	if got := CleanHints(nil); got == nil || len(got) != 0 {
		t.Errorf("CleanHints(nil) = %#v, want empty non-nil", got)
	}
}

func TestValidateHintTooLong(t *testing.T) {
	run := func(hint string) []dtos.FieldError {
		errs := []dtos.FieldError{}
		add := func(f, m string) { errs = append(errs, dtos.FieldError{Field: f, Message: m}) }
		validateSteps(dtos.UpsertRecipeRequest{Steps: []dtos.UpsertStepRequest{
			{Text: "a"},
			{Text: "b", Hints: []string{hint}},
		}}, add)
		return errs
	}
	if errs := run(strings.Repeat("ż", 300)); len(errs) != 0 {
		t.Errorf("300 znaków: %v", errs)
	}
	errs := run(strings.Repeat("ż", 301))
	if len(errs) != 1 || errs[0].Field != "steps" || errs[0].Message != "Krok 2: wskazówka może mieć najwyżej 300 znaków" {
		t.Errorf("301 znaków: %#v", errs)
	}
	if errs := run("   " + strings.Repeat("a", 300) + "   "); len(errs) != 0 {
		t.Errorf("spacje wokół nie liczą się: %v", errs)
	}
}
