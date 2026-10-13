package services

import (
	"math"
	"testing"
)

func TestParseQuantity(t *testing.T) {
	ok := []struct {
		in       string
		min, max float64
	}{
		{"2", 2, 2},
		{"1,5", 1.5, 1.5},
		{"1.5", 1.5, 1.5},
		{"½", 0.5, 0.5},
		{"1½", 1.5, 1.5},
		{"1 ½", 1.5, 1.5},
		{"1/4", 0.25, 0.25},
		{"1 1/2", 1.5, 1.5},
		{"2–3", 2, 3},
		{"2-3", 2, 3},
		{"2 - 3", 2, 3},
		{"½–1", 0.5, 1},
	}
	for _, c := range ok {
		q, valid := ParseQuantity(c.in)
		if !valid || math.Abs(q.Min-c.min) > 1e-9 || math.Abs(q.Max-c.max) > 1e-9 {
			t.Errorf("ParseQuantity(%q) = %+v, %v; want %v–%v", c.in, q, valid, c.min, c.max)
		}
	}

	if q, valid := ParseQuantity("  "); !valid || !q.Empty {
		t.Errorf("empty input should be valid and Empty, got %+v %v", q, valid)
	}

	for _, in := range []string{"abc", "0", "3–2", "2–2", "1/0", "2–", "1-2-3"} {
		if _, valid := ParseQuantity(in); valid {
			t.Errorf("ParseQuantity(%q) should be invalid", in)
		}
	}
}
