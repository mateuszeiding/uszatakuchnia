package handlers

import (
	"encoding/json"
	"testing"
)

func TestValueRequestIcon(t *testing.T) {
	cases := []struct {
		body  string
		set   bool
		value string // "" = nil
		valid bool
	}{
		{`{"label":"zupa"}`, false, "", true}, // PUT without the field keeps the icon
		{`{"label":"zupa","icon":null}`, true, "", true},
		{`{"label":"zupa","icon":""}`, true, "", true},
		{`{"label":"zupa","icon":"zupa"}`, true, "zupa", true},
		{`{"label":"zupa","icon":"pizza"}`, true, "pizza", false},
		// a repeated member: the last one wins, also when it clears the icon
		{`{"label":"zupa","icon":"zupa","icon":null}`, true, "", true},
		{`{"label":"zupa","icon":"zupa","icon":""}`, true, "", true},
		{`{"label":"zupa","icon":"zupa","icon":"deser"}`, true, "deser", true},
	}
	for _, c := range cases {
		var req valueRequest
		if err := json.Unmarshal([]byte(c.body), &req); err != nil {
			t.Fatalf("%s: %v", c.body, err)
		}
		got := ""
		if req.Icon.Value != nil {
			got = *req.Icon.Value
		}
		if req.Icon.Set != c.set || got != c.value || req.Icon.valid() != c.valid {
			t.Errorf("%s: set=%v value=%q valid=%v", c.body, req.Icon.Set, got, req.Icon.valid())
		}
	}
	var req valueRequest
	if err := json.Unmarshal([]byte(`{"icon":3}`), &req); err == nil {
		t.Error("a non-string icon must be rejected")
	}
}
