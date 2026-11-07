package services

import (
	"math"
	"reflect"
	"testing"

	"uszatakuchnia/db/entities"
)

func fp(v float64) *float64 { return &v }
func up(v uint) *uint       { return &v }

func rng(min float64, max ...float64) entities.MacroRange {
	r := entities.MacroRange{Min: min}
	if len(max) > 0 {
		r.Max = fp(max[0])
	}
	return r
}

const (
	unitG uint = iota + 1
	unitKg
	unitMl
	unitL
	unitSzt
	unitLyzka
	unitDoSmaku
)

// Drzewko: makaron(1) → spaghetti(2), fusilli(3); jajko(4); oliwa(5) → oliwa extra virgin(8); oregano(6, bez danych); mleko(7, przelicznik tylko dla l).
func testIndex() *NutritionIndex {
	return &NutritionIndex{
		Own: map[uint]Macros{
			2: {"kcal": rng(350, 360), "fiber": rng(3, 3.5)},
			3: {"kcal": rng(330, 345), "fiber": rng(6, 11)},
			4: {"kcal": rng(140), "protein": rng(12, 13)},
			5: {"kcal": rng(884), "fat": rng(100)},
		},
		Children: map[uint][]uint{1: {2, 3}},
		Parent:   map[uint]*uint{2: up(1), 3: up(1), 8: up(5)},
		Conversions: map[uint]map[uint]entities.MacroRange{
			1: {unitLyzka: rng(15)},
			4: {unitSzt: rng(50, 65)},
			5: {unitLyzka: rng(10), unitMl: rng(0.92)},
			7: {unitL: rng(1030, 1040)},
		},
		NoQtyUnits: map[uint]bool{unitDoSmaku: true},
		Names:      map[uint]string{1: "makaron", 2: "spaghetti", 3: "fusilli", 4: "jajko", 5: "oliwa", 6: "oregano", 7: "mleko", 8: "oliwa extra virgin"},
		UnitIDByName: map[string]uint{
			"g": unitG, "kg": unitKg, "ml": unitMl, "l": unitL, "sztuka": unitSzt, "łyżka": unitLyzka, "do smaku": unitDoSmaku,
		},
	}
}

func near(a, b float64) bool { return math.Abs(a-b) < 1e-9 }

func eqRange(got entities.MacroRange, min float64, max *float64) bool {
	if !near(got.Min, min) {
		return false
	}
	if max == nil {
		return got.Max == nil
	}
	return got.Max != nil && near(*got.Max, *max)
}

func TestResolveMacrosInherits(t *testing.T) {
	ix := testIndex()
	ix.Own[3]["kcal"] = rng(330) // pojedyncza wartość liczona jako 330–330
	m := ix.ResolveMacros(1)
	if !eqRange(m["kcal"], 330, fp(360)) {
		t.Errorf("kcal = %+v, want 330–360", m["kcal"])
	}
	if !eqRange(m["fiber"], 3, fp(11)) {
		t.Errorf("fiber = %+v, want 3–11", m["fiber"])
	}
	if _, ok := m["protein"]; ok {
		t.Errorf("protein: dzieci nie mają danych, a jest %+v", m["protein"])
	}
}

func TestResolveMacrosSingleWhenEqual(t *testing.T) {
	ix := testIndex()
	ix.Own[2]["kcal"] = rng(350)
	ix.Own[3]["kcal"] = rng(350)
	if m := ix.ResolveMacros(1); !eqRange(m["kcal"], 350, nil) {
		t.Errorf("kcal = %+v, want 350 (single)", m["kcal"])
	}
}

func TestResolveMacrosOwnWins(t *testing.T) {
	ix := testIndex()
	ix.Own[1] = Macros{"kcal": rng(355)}
	m := ix.ResolveMacros(1)
	if !eqRange(m["kcal"], 355, nil) {
		t.Errorf("kcal = %+v, want własne 355", m["kcal"])
	}
	if !eqRange(m["fiber"], 3, fp(11)) {
		t.Errorf("fiber = %+v, want odziedziczone 3–11", m["fiber"])
	}
}

func TestResolveMacrosCycle(t *testing.T) {
	ix := testIndex()
	ix.Children[2] = []uint{1} // 1 → 2 → 1
	_ = ix.ResolveMacros(1)    // nie może się zawiesić
}

func TestGramsPerUnit(t *testing.T) {
	ix := testIndex()
	cases := []struct {
		name     string
		ing      uint
		unit     *uint
		unitName string
		ok       bool
		min      float64
		max      *float64
	}{
		{"g", 2, up(unitG), "g", true, 1, nil},
		{"kg", 2, up(unitKg), "kg", true, 1000, nil},
		{"ml bez przelicznika", 2, up(unitMl), "ml", true, 1, nil},
		{"l bez przelicznika", 2, up(unitL), "l", true, 1000, nil},
		{"ml z przelicznikiem", 5, up(unitMl), "ml", true, 0.92, nil},
		{"l z przelicznika ml", 5, up(unitL), "l", true, 920, nil},
		{"ml z przelicznika l", 7, up(unitMl), "ml", true, 1.03, fp(1.04)},
		{"l z przelicznika ml u rodzica", 8, up(unitL), "l", true, 920, nil},
		{"łyżka u rodzica", 2, up(unitLyzka), "łyżka", true, 15, nil},
		{"brak jednostki = sztuka", 4, nil, "", true, 50, fp(65)},
		{"brak przelicznika", 6, up(unitLyzka), "łyżka", false, 0, nil},
	}
	for _, c := range cases {
		got, ok := ix.GramsPerUnit(c.ing, c.unit, c.unitName)
		if ok != c.ok || (ok && !eqRange(got, c.min, c.max)) {
			t.Errorf("%s: got %+v, %v", c.name, got, ok)
		}
	}
}

func line(ing uint, qty string, unit *uint, unitName string) NutritionLine {
	q, _ := ParseQuantity(qty)
	return NutritionLine{IngredientID: ing, Qty: q, UnitID: unit, UnitName: unitName}
}

func TestComputeRanges(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{
		line(2, "200", up(unitG), "g"),
		line(4, "2", up(unitSzt), "sztuka"),
	}, 2)
	if n == nil {
		t.Fatal("nil")
	}
	// spaghetti 700–720 + jajka 2×(50–65) g × 140/100 = 140–182 → 840–902 / 2
	if got := n.PerPortion["kcal"]; !near(got.Min, 420) || !near(got.Max, 451) {
		t.Errorf("kcal = %+v, want 420–451", got)
	}
	// białko tylko z jajek: 100–130 g × 12–13 / 100 = 12–16,9 → /2
	if got := n.PerPortion["protein"]; !near(got.Min, 6) || !near(got.Max, 8.45) {
		t.Errorf("protein = %+v, want 6–8.45", got)
	}
	if len(n.Missing) != 0 {
		t.Errorf("missing = %v", n.Missing)
	}
}

func TestComputeQuantityRange(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{line(2, "100–200", up(unitG), "g")}, 1)
	if got := n.PerPortion["kcal"]; !near(got.Min, 350) || !near(got.Max, 720) {
		t.Errorf("kcal = %+v, want 350–720", got)
	}
}

func TestComputeSkipsAndMissing(t *testing.T) {
	ix := testIndex()
	optional := line(5, "2", up(unitLyzka), "łyżka")
	optional.Optional = true
	toTaste := line(6, "", up(unitDoSmaku), "do smaku")
	toTaste.NoQty = true
	n := ix.Compute([]NutritionLine{
		line(2, "100", up(unitG), "g"),
		optional,
		toTaste,
		line(5, "", up(unitLyzka), "łyżka"),  // brak ilości — pomijane
		line(6, "1", up(unitG), "g"),         // brak makro
		line(4, "1", up(unitLyzka), "łyżka"), // brak przelicznika
		line(6, "2", up(unitG), "g"),         // duplikat nazwy
	}, 1)
	if got := n.PerPortion["kcal"]; !near(got.Min, 350) || !near(got.Max, 360) {
		t.Errorf("kcal = %+v, want 350–360 (tylko spaghetti)", got)
	}
	if !reflect.DeepEqual(n.Missing, []string{"oregano", "jajko"}) {
		t.Errorf("missing = %v, want [oregano jajko]", n.Missing)
	}
}

func TestComputePartialFields(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{
		line(5, "100", up(unitG), "g"), // oliwa: kcal + fat, bez błonnika
		line(2, "100", up(unitG), "g"), // spaghetti: kcal + fiber
	}, 1)
	if len(n.Missing) != 0 {
		t.Errorf("missing = %v, want []", n.Missing)
	}
	if got := n.PerPortion["fiber"]; !near(got.Min, 3) || !near(got.Max, 3.5) {
		t.Errorf("fiber = %+v", got)
	}
	if got := n.PerPortion["fat"]; !near(got.Min, 100) || !near(got.Max, 100) {
		t.Errorf("fat = %+v", got)
	}
	if _, ok := n.PerPortion["protein"]; ok {
		t.Errorf("protein nie powinno istnieć")
	}
}

func TestComputeServingsZero(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{line(2, "100", up(unitG), "g")}, 0)
	if got := n.PerPortion["kcal"]; !near(got.Min, 350) {
		t.Errorf("kcal = %+v, want całość 350–360", got)
	}
}

func TestComputeNothing(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{line(6, "1", up(unitG), "g")}, 1)
	if n != nil {
		t.Errorf("want nil, got %+v", n)
	}
}

func TestComputeIncomplete(t *testing.T) {
	ix := testIndex()
	ix.Own[9] = Macros{"kcal": rng(120, 145)}
	ix.Names[9] = "tofu naturalne"
	n := ix.Compute([]NutritionLine{
		line(9, "200", up(unitG), "g"),
		line(9, "100", up(unitG), "g"), // duplikat — nazwa raz
	}, 2)
	if got := n.PerPortion["kcal"]; !near(got.Min, 180) || !near(got.Max, 217.5) {
		t.Errorf("kcal = %+v, want 180–217,5", got)
	}
	want := map[string][]string{}
	for _, k := range entities.NutrientKeys {
		if k != "kcal" {
			want[k] = []string{"tofu naturalne"}
		}
	}
	if !reflect.DeepEqual(n.Incomplete, want) {
		t.Errorf("incomplete = %v, want %v", n.Incomplete, want)
	}
	if len(n.Missing) != 0 {
		t.Errorf("missing = %v", n.Missing)
	}
}

func TestComputeIncompleteMixed(t *testing.T) {
	ix := testIndex()
	n := ix.Compute([]NutritionLine{
		line(4, "1", up(unitSzt), "sztuka"), // jajko: kcal + białko
		line(5, "10", up(unitG), "g"),       // oliwa: kcal + tłuszcz
	}, 1)
	if _, ok := n.Incomplete["kcal"]; ok {
		t.Errorf("kcal mają wszystkie, a incomplete = %v", n.Incomplete["kcal"])
	}
	if !reflect.DeepEqual(n.Incomplete["protein"], []string{"oliwa"}) {
		t.Errorf("protein = %v, want [oliwa]", n.Incomplete["protein"])
	}
	if !reflect.DeepEqual(n.Incomplete["fiber"], []string{"jajko", "oliwa"}) {
		t.Errorf("fiber = %v, want [jajko oliwa]", n.Incomplete["fiber"])
	}
}
