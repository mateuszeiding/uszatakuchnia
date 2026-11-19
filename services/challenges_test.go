package services

import (
	"strings"
	"testing"

	"uszatakuchnia/dtos"
)

func ptr[T any](v T) *T { return &v }

func fields(errs []dtos.FieldError) map[string]bool {
	out := map[string]bool{}
	for _, e := range errs {
		out[e.Field] = true
	}
	return out
}

func TestNormalizeChallengeRequest(t *testing.T) {
	req := NormalizeChallengeRequest(dtos.UpsertChallengeRequest{
		Title:       "  Ryba   w 20 minut ",
		Description: "Pierwszy akapit.  \r\n\r\n\r\n\r\nDrugi\r\nz łamaniem.\n\n",
		Source:      ptr("   "),
		Ingredients: []dtos.ChallengeIngredientDto{
			{Name: " filet z dorsza ", Required: true},
			{Name: "   "},
			{Name: ""},
			{Name: "cytryna"},
		},
	})
	if req.Title != "Ryba w 20 minut" {
		t.Errorf("title = %q", req.Title)
	}
	if req.Description != "Pierwszy akapit.\n\nDrugi\nz łamaniem." {
		t.Errorf("description = %q", req.Description)
	}
	if req.Source != nil {
		t.Errorf("empty source should be nil, got %q", *req.Source)
	}
	if len(req.Ingredients) != 2 || req.Ingredients[0].Name != "filet z dorsza" || !req.Ingredients[0].Required ||
		req.Ingredients[1].Name != "cytryna" {
		t.Errorf("ingredients = %+v", req.Ingredients)
	}
}

func TestValidateChallengeFields(t *testing.T) {
	ok := dtos.UpsertChallengeRequest{Title: "T", FormatID: 1, Description: "D", Difficulty: 3}
	if errs := validateChallengeFields(ok); len(errs) != 0 {
		t.Fatalf("valid request: %+v", errs)
	}
	byName := ok
	byName.FormatID, byName.Format = 0, "na temat"
	if errs := validateChallengeFields(byName); len(errs) != 0 {
		t.Fatalf("format by name: %+v", errs)
	}

	cases := []struct {
		name  string
		edit  func(r *dtos.UpsertChallengeRequest)
		field string
	}{
		{"no title", func(r *dtos.UpsertChallengeRequest) { r.Title = "" }, "title"},
		{"long title", func(r *dtos.UpsertChallengeRequest) { r.Title = strings.Repeat("ą", 256) }, "title"},
		{"no description", func(r *dtos.UpsertChallengeRequest) { r.Description = "" }, "description"},
		{"no format", func(r *dtos.UpsertChallengeRequest) { r.FormatID = 0 }, "formatId"},
		{"difficulty 0", func(r *dtos.UpsertChallengeRequest) { r.Difficulty = 0 }, "difficulty"},
		{"difficulty 6", func(r *dtos.UpsertChallengeRequest) { r.Difficulty = 6 }, "difficulty"},
		{"time 0", func(r *dtos.UpsertChallengeRequest) { r.TimeLimitMin = ptr(0) }, "timeLimitMin"},
		{"time negative", func(r *dtos.UpsertChallengeRequest) { r.TimeLimitMin = ptr(-5) }, "timeLimitMin"},
		{"long source", func(r *dtos.UpsertChallengeRequest) { r.Source = ptr(strings.Repeat("x", 256)) }, "source"},
		{"long ingredient", func(r *dtos.UpsertChallengeRequest) {
			r.Ingredients = []dtos.ChallengeIngredientDto{{Name: "ok"}, {Name: strings.Repeat("x", 121)}}
		}, "ingredients[1].name"},
	}
	for _, c := range cases {
		r := ok
		c.edit(&r)
		if got := fields(validateChallengeFields(r)); !got[c.field] {
			t.Errorf("%s: expected error on %s, got %v", c.name, c.field, got)
		}
	}

	withTime := ok
	withTime.TimeLimitMin = ptr(45)
	if errs := validateChallengeFields(withTime); len(errs) != 0 {
		t.Errorf("45 min should pass: %+v", errs)
	}
}
