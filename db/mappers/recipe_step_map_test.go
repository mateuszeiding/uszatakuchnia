package mappers

import (
	"encoding/json"
	"testing"

	"uszatakuchnia/db/entities"
)

func TestMapStepHintsNilIsEmpty(t *testing.T) {
	dto := MapRecipeStepToDto(entities.RecipeStep{StepNo: 1, Text: "x"})
	b, _ := json.Marshal(dto)
	var m map[string]any
	_ = json.Unmarshal(b, &m)
	h, ok := m["hints"].([]any)
	if !ok || len(h) != 0 {
		t.Errorf("hints = %#v, want []", m["hints"])
	}
}

func TestMapStepHints(t *testing.T) {
	dto := MapRecipeStepToDto(entities.RecipeStep{Hints: entities.StringList{"a"}})
	if len(dto.Hints) != 1 || dto.Hints[0] != "a" {
		t.Errorf("hints = %v", dto.Hints)
	}
}

func TestMapCautionsNilIsEmpty(t *testing.T) {
	for name, v := range map[string]any{
		"recipe":    MapRecipeToDto(entities.Recipe{}),
		"challenge": MapChallengeToDto(entities.Challenge{}),
	} {
		b, _ := json.Marshal(v)
		var m map[string]any
		_ = json.Unmarshal(b, &m)
		if c, ok := m["cautions"].([]any); !ok || len(c) != 0 {
			t.Errorf("%s: cautions = %#v, want []", name, m["cautions"])
		}
	}
}
