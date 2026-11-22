package entities

import (
	"os"
	"regexp"
	"slices"
	"testing"
)

// The API accepts only icons the frontend can draw: TaxonomyIcons must list exactly the DISH_ICONS keys.
func TestTaxonomyIconsMatchFrontend(t *testing.T) {
	src, err := os.ReadFile("../../web/src/shared/dishIcons.ts")
	if err != nil {
		t.Fatal(err)
	}
	block := regexp.MustCompile(`(?s)export const DISH_ICONS = \{(.*?)\n\} satisfies`).FindSubmatch(src)
	if block == nil {
		t.Fatal("DISH_ICONS not found in dishIcons.ts")
	}
	var keys []string
	for _, m := range regexp.MustCompile(`(?m)^    (\w+): \[`).FindAllSubmatch(block[1], -1) {
		keys = append(keys, string(m[1]))
	}
	got, want := slices.Sorted(slices.Values(TaxonomyIcons)), slices.Sorted(slices.Values(keys))
	if !slices.Equal(got, want) {
		t.Fatalf("TaxonomyIcons %v, DISH_ICONS keys %v", got, want)
	}
}
