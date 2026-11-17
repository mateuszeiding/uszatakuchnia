package db

import "testing"

func TestSchemaFingerprintStable(t *testing.T) {
	a, b := SchemaFingerprint(), SchemaFingerprint()
	if a == "" || a != b {
		t.Fatalf("odcisk niestabilny: %q vs %q", a, b)
	}
}
