package entities

import "testing"

func TestStringListRoundTrip(t *testing.T) {
	in := StringList{"a", "b"}
	raw, err := in.Value()
	if err != nil {
		t.Fatalf("Value: %v", err)
	}
	var out StringList
	if err := out.Scan([]byte(raw.(string))); err != nil {
		t.Fatalf("Scan: %v", err)
	}
	if len(out) != 2 || out[0] != "a" || out[1] != "b" {
		t.Errorf("got %v", out)
	}
}

func TestStringListScanNil(t *testing.T) {
	out := StringList{"x"}
	if err := out.Scan(nil); err != nil {
		t.Fatalf("Scan(nil): %v", err)
	}
	if len(out) != 0 {
		t.Errorf("Scan(nil) = %v, want empty", out)
	}
	raw, err := StringList(nil).Value()
	if err != nil || raw != "[]" {
		t.Errorf("nil Value = %v, %v; want \"[]\"", raw, err)
	}
}
