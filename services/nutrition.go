package services

import "uszatakuchnia/db/entities"

// Macros — makro na 100 g po kluczu z entities.NutrientKeys; brak klucza = brak danych.
type Macros map[string]entities.MacroRange

// NutritionIndex — dane potrzebne do wyliczeń: własne makro, drzewko składników i przeliczniki jednostek.
type NutritionIndex struct {
	Own          map[uint]Macros
	Children     map[uint][]uint
	Parent       map[uint]*uint
	Conversions  map[uint]map[uint]entities.MacroRange // składnik → jednostka → gramy
	Names        map[uint]string
	UnitIDByName map[string]uint
}

// Jednostki przeliczane bez wierszy w ingredient_unit_conversions. ml / l tylko gdy składnik nie ma własnego przelicznika.
var (
	fixedGrams    = map[string]float64{"g": 1, "kg": 1000}
	fallbackGrams = map[string]float64{"ml": 1, "l": 1000}
)

// ResolveMacros — per wartość: własna, a jeśli jej brak, z podskładników (najniższe min, najwyższe max).
func (ix *NutritionIndex) ResolveMacros(id uint) Macros {
	return ix.resolveMacros(id, map[uint]bool{})
}

func (ix *NutritionIndex) resolveMacros(id uint, visiting map[uint]bool) Macros {
	if visiting[id] {
		return Macros{} // cykl w drzewie — nie powinien wystąpić, walidacja rodzica go blokuje
	}
	visiting[id] = true
	defer delete(visiting, id)

	out := Macros{}
	for k, v := range ix.Own[id] {
		out[k] = v
	}
	if len(out) == len(entities.NutrientKeys) {
		return out
	}

	var children []Macros
	for _, c := range ix.Children[id] {
		children = append(children, ix.resolveMacros(c, visiting))
	}
	for _, k := range entities.NutrientKeys {
		if _, own := out[k]; own {
			continue
		}
		found := false
		var lo, hi float64
		for _, cm := range children {
			v, ok := cm[k]
			if !ok {
				continue
			}
			if !found || v.Min < lo {
				lo = v.Min
			}
			if !found || v.Upper() > hi {
				hi = v.Upper()
			}
			found = true
		}
		if !found {
			continue
		}
		r := entities.MacroRange{Min: lo}
		if hi > lo {
			h := hi
			r.Max = &h
		}
		out[k] = r
	}
	return out
}

// ResolveConversion — przelicznik jednostki: własny, a jeśli brak, u najbliższego przodka.
func (ix *NutritionIndex) ResolveConversion(id, unitID uint) (entities.MacroRange, bool) {
	seen := map[uint]bool{}
	for cur := &id; cur != nil && !seen[*cur]; cur = ix.Parent[*cur] {
		seen[*cur] = true
		if r, ok := ix.Conversions[*cur][unitID]; ok {
			return r, true
		}
	}
	return entities.MacroRange{}, false
}

// GramsPerUnit — ile gramów to jedna jednostka składnika. Pozycja bez jednostki liczona jak „sztuka”.
func (ix *NutritionIndex) GramsPerUnit(id uint, unitID *uint, unitName string) (entities.MacroRange, bool) {
	if unitID == nil {
		sztuka, ok := ix.UnitIDByName["sztuka"]
		if !ok {
			return entities.MacroRange{}, false
		}
		unitID, unitName = &sztuka, "sztuka"
	}
	if g, ok := fixedGrams[unitName]; ok {
		return entities.MacroRange{Min: g}, true
	}
	if r, ok := ix.ResolveConversion(id, *unitID); ok {
		return r, true
	}
	if g, ok := fallbackGrams[unitName]; ok {
		return entities.MacroRange{Min: g}, true
	}
	return entities.MacroRange{}, false
}

// NutritionLine — pozycja przepisu do wyliczenia makro.
type NutritionLine struct {
	IngredientID uint
	Qty          Quantity
	UnitID       *uint
	UnitName     string
	NoQty        bool // jednostka „do smaku”
	Optional     bool
}

// Compute — makro na porcję jako suma zakresów (min×min, max×max). Pomija opcjonalne, „do smaku” i pozycje bez ilości;
// pozycje bez makro lub bez przelicznika trafiają do Missing. Nil, gdy nic nie dało się policzyć.
func (ix *NutritionIndex) Compute(lines []NutritionLine, servings int) *entities.RecipeNutrition {
	sum := map[string]*entities.SumRange{}
	missing := []string{}
	seenMissing := map[string]bool{}
	contributed := false

	for _, l := range lines {
		if l.Optional || l.NoQty || l.Qty.Empty {
			continue
		}
		grams, okGrams := ix.GramsPerUnit(l.IngredientID, l.UnitID, l.UnitName)
		macros := ix.ResolveMacros(l.IngredientID)
		if !okGrams || len(macros) == 0 {
			name := ix.Names[l.IngredientID]
			if !seenMissing[name] {
				seenMissing[name] = true
				missing = append(missing, name)
			}
			continue
		}
		contributed = true
		gMin, gMax := l.Qty.Min*grams.Min, l.Qty.Max*grams.Upper()
		for k, v := range macros {
			s := sum[k]
			if s == nil {
				s = &entities.SumRange{}
				sum[k] = s
			}
			s.Min += gMin * v.Min / 100
			s.Max += gMax * v.Upper() / 100
		}
	}
	if !contributed {
		return nil
	}

	if servings < 1 {
		servings = 1
	}
	per := map[string]entities.SumRange{}
	for k, s := range sum {
		per[k] = entities.SumRange{Min: s.Min / float64(servings), Max: s.Max / float64(servings)}
	}
	return &entities.RecipeNutrition{PerPortion: per, Missing: missing}
}
