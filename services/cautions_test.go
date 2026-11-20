package services

import (
	"strings"
	"testing"

	"uszatakuchnia/dtos"
)

func TestCleanCautions(t *testing.T) {
	got := CleanCautions([]string{"  Mięso  tylko\nświeże. ", "", "   ", "\t", "Drugie"})
	if len(got) != 2 || got[0] != "Mięso tylko świeże." || got[1] != "Drugie" {
		t.Errorf("CleanCautions = %#v", got)
	}
	if got := CleanCautions(nil); got == nil || len(got) != 0 {
		t.Errorf("CleanCautions(nil) = %#v, want empty non-nil", got)
	}
}

func TestValidateCautionsLimit(t *testing.T) {
	run := func(in []string) []dtos.FieldError {
		errs := []dtos.FieldError{}
		validateCautions(in, func(f, m string) { errs = append(errs, dtos.FieldError{Field: f, Message: m}) })
		return errs
	}
	ok := strings.Repeat("ż", MaxCautionLen) // znaki, nie bajty
	if errs := run([]string{ok}); len(errs) != 0 {
		t.Errorf("220 znaków: %v", errs)
	}
	// pusta pozycja nie liczy się do indeksu (formularz też ją pomija)
	errs := run([]string{"", "a", ok + "x"})
	if len(errs) != 1 || errs[0].Field != "cautions[1]" {
		t.Errorf("221 znaków: %v", errs)
	}
}

func TestChallengeCautions(t *testing.T) {
	req := NormalizeChallengeRequest(dtos.UpsertChallengeRequest{Cautions: &[]string{" a ", "  "}})
	if req.Cautions == nil || len(*req.Cautions) != 1 || (*req.Cautions)[0] != "a" {
		t.Errorf("cautions = %#v", req.Cautions)
	}
	if req := NormalizeChallengeRequest(dtos.UpsertChallengeRequest{}); req.Cautions != nil {
		t.Errorf("brak pola = bez zmian, got %#v", req.Cautions)
	}
	long := strings.Repeat("x", MaxCautionLen+1)
	errs := validateChallengeFields(dtos.UpsertChallengeRequest{Cautions: &[]string{long}})
	if !fields(errs)["cautions[0]"] {
		t.Errorf("errs = %v", errs)
	}
}
