package services

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// Macros: macros per 100 g keyed by entities.NutrientKeys; a missing key = no data.
type Macros map[string]entities.MacroRange

// NutritionIndex holds the data for the calculations: own macros, the ingredient tree and unit conversions.
type NutritionIndex struct {
	Own          map[uint]Macros
	Children     map[uint][]uint
	Parent       map[uint]*uint
	Conversions  map[uint]map[uint]entities.MacroRange // ingredient → unit → grams
	Names        map[uint]string
	NoQtyUnits   map[uint]bool // "do smaku" (to taste) units
	UnitIDByName map[string]uint
}

// Units converted without rows in ingredient_unit_conversions. ml / l (water density) only when the ingredient has no ml or l conversion.
var (
	fixedGrams    = map[string]float64{"g": 1, "kg": 1000}
	fallbackGrams = map[string]float64{"ml": 1, "l": 1000}
	// When ml / l has no conversion (own or inherited), we derive it from the other volume unit.
	volumeSiblings = map[string]struct {
		unit   string
		factor float64
	}{"l": {"ml", 1000}, "ml": {"l", 0.001}}
)

func scaleRange(r entities.MacroRange, f float64) entities.MacroRange {
	out := entities.MacroRange{Min: r.Min * f}
	if r.Max != nil {
		m := *r.Max * f
		out.Max = &m
	}
	return out
}

// ResolveMacros: per value, the own one, or if missing, from sub-ingredients (lowest min, highest max).
func (ix *NutritionIndex) ResolveMacros(id uint) Macros {
	return ix.resolveMacros(id, map[uint]bool{})
}

func (ix *NutritionIndex) resolveMacros(id uint, visiting map[uint]bool) Macros {
	if visiting[id] {
		return Macros{} // a cycle in the tree: shouldn't happen, parent validation blocks it
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

// ResolveConversion: the unit's own conversion, or if missing, the nearest ancestor's.
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

// GramsPerUnit: how many grams one unit of an ingredient weighs. A row without a unit counts as "sztuka" (piece).
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
	if v, ok := volumeSiblings[unitName]; ok {
		if sibID, ok := ix.UnitIDByName[v.unit]; ok {
			if r, ok := ix.ResolveConversion(id, sibID); ok {
				return scaleRange(r, v.factor), true
			}
		}
	}
	if g, ok := fallbackGrams[unitName]; ok {
		return entities.MacroRange{Min: g}, true
	}
	return entities.MacroRange{}, false
}

// FillGramsPerUnit sets grams per 1 unit for every recipe row (the shopping list converts mixed units through them).
// Same rules as for macros: g / kg fixed, ml / l from a conversion or as water, no conversion → nil. "Do smaku" → nil.
func FillGramsPerUnit(items []dtos.RecipeIngredientDto, ix *NutritionIndex) {
	for i := range items {
		it := &items[i]
		it.GramsPerUnit = nil
		if it.UnitID != nil && ix.NoQtyUnits[*it.UnitID] {
			continue
		}
		name := ""
		if it.Unit != nil {
			name = *it.Unit
		}
		g, ok := ix.GramsPerUnit(it.IngredientID, it.UnitID, name)
		if !ok {
			continue
		}
		it.GramsPerUnit = &dtos.GramsRangeDto{Min: g.Min, Max: g.Upper()}
	}
}

// NutritionLine is a recipe row for the macro calculation.
type NutritionLine struct {
	IngredientID uint
	Qty          Quantity
	UnitID       *uint
	UnitName     string
	NoQty        bool // "do smaku" (to taste) unit
	Optional     bool
}

// Compute returns per-serving macros as a sum of ranges (min×min, max×max). Skips optional, "to taste" and quantity-less rows;
// rows without macros or a conversion go to Missing, and rows counted without some values go to Incomplete.
// Nil when nothing could be computed.
func (ix *NutritionIndex) Compute(lines []NutritionLine, servings int) *entities.RecipeNutrition {
	n, _, contributed := ix.accumulate(lines, servings)
	if !contributed {
		return nil
	}
	return n
}

// LineNutrition is a row's (or a "half and half" pair's) share of the per-serving macros, in the same format as Compute.
// The shares of all rows add up to Compute (the frontend computes the chosen variant's macros from them).
// A row without macros / conversion: empty PerPortion and its name in Missing. Nil when the row doesn't count
// (optional, to taste, no quantity).
func (ix *NutritionIndex) LineNutrition(lines []NutritionLine, servings int) *entities.RecipeNutrition {
	n, counted, _ := ix.accumulate(lines, servings)
	if !counted {
		return nil
	}
	return n
}

// accumulate is the shared summing for Compute and LineNutrition. counted: some row counts (included or
// skipped for missing data); contributed: some row added something.
func (ix *NutritionIndex) accumulate(lines []NutritionLine, servings int) (n *entities.RecipeNutrition, counted, contributed bool) {
	sum := map[string]*entities.SumRange{}
	missing := []string{}
	seenMissing := map[string]bool{}
	incomplete := map[string][]string{}
	seenIncomplete := map[[2]string]bool{} // [value, name]

	for _, l := range lines {
		if l.Optional || l.NoQty || l.Qty.Empty {
			continue
		}
		counted = true
		grams, okGrams := ix.GramsPerUnit(l.IngredientID, l.UnitID, l.UnitName)
		macros := ix.ResolveMacros(l.IngredientID)
		name := ix.Names[l.IngredientID]
		if !okGrams || len(macros) == 0 {
			if !seenMissing[name] {
				seenMissing[name] = true
				missing = append(missing, name)
			}
			continue
		}
		contributed = true
		for _, k := range entities.NutrientKeys {
			if _, ok := macros[k]; ok || seenIncomplete[[2]string{k, name}] {
				continue
			}
			seenIncomplete[[2]string{k, name}] = true
			incomplete[k] = append(incomplete[k], name)
		}
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

	if servings < 1 {
		servings = 1
	}
	per := map[string]entities.SumRange{}
	for k, s := range sum {
		per[k] = entities.SumRange{Min: s.Min / float64(servings), Max: s.Max / float64(servings)}
	}
	return &entities.RecipeNutrition{PerPortion: per, Missing: missing, Incomplete: incomplete}, counted, contributed
}

// Half returns the same row with half the quantity ("half and half").
func (l NutritionLine) Half() NutritionLine {
	h := l
	h.Qty.Min /= 2
	h.Qty.Max /= 2
	return h
}
