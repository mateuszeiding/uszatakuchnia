package services

import (
	"fmt"
	"slices"
	"strings"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

const msgMaxBelowMin = "Max nie może być mniejsze niż min"

// ValidateIngredientRequest — makro, przeliczniki i rodzic składnika selfID (0 = nowy).
// „Max bez min” nie występuje w JSON-ie (min jest wymagane w parze), pilnuje tego formularz.
func ValidateIngredientRequest(req dtos.IngredientRequest, ix *NutritionIndex, selfID uint) []dtos.FieldError {
	errs := []dtos.FieldError{}
	add := func(field, msg string) { errs = append(errs, dtos.FieldError{Field: field, Message: msg}) }

	for k, r := range req.Nutrition {
		f := "nutrition." + k
		if !slices.Contains(entities.NutrientKeys, k) {
			add(f, "Nieznana wartość odżywcza")
			continue
		}
		validateRange(r.Min, r.Max, f+".min", f+".max", add)
	}

	units := map[uint]string{}
	for name, id := range ix.UnitIDByName {
		units[id] = name
	}
	seen := map[uint]bool{}
	for i, c := range req.Conversions {
		f := fmt.Sprintf("conversions[%d]", i)
		name, known := units[c.UnitID]
		switch {
		case !known:
			add(f+".unitId", "Nieznana jednostka")
		case fixedGrams[name] != 0 || ix.NoQtyUnits[c.UnitID]:
			add(f+".unitId", fmt.Sprintf("Jednostka „%s” nie potrzebuje przelicznika", name))
		case seen[c.UnitID]:
			add(f+".unitId", "Ta jednostka ma już przelicznik")
		}
		seen[c.UnitID] = true
		validateRange(c.GramsMin, c.GramsMax, f+".gramsMin", f+".gramsMax", add)
	}

	if req.ParentID != nil {
		if _, exists := ix.Names[*req.ParentID]; !exists {
			add("parentId", "Nieznany składnik nadrzędny")
		} else if selfID != 0 && ix.isSelfOrDescendant(*req.ParentID, selfID) {
			add("parentId", "Składnik nie może być własnym podskładnikiem")
		}
	}
	return errs
}

func validateRange(min float64, max *float64, fMin, fMax string, add func(string, string)) {
	if min < 0 {
		add(fMin, "Wartość nie może być ujemna")
	}
	if max != nil && *max < 0 {
		add(fMax, "Wartość nie może być ujemna")
	} else if max != nil && *max < min {
		add(fMax, msgMaxBelowMin)
	}
}

// isSelfOrDescendant — czy id to root albo leży w jego poddrzewie (idąc w górę od id).
func (ix *NutritionIndex) isSelfOrDescendant(id, root uint) bool {
	seen := map[uint]bool{}
	for cur := &id; cur != nil && !seen[*cur]; cur = ix.Parent[*cur] {
		if *cur == root {
			return true
		}
		seen[*cur] = true
	}
	return false
}

// SaveIngredient zapisuje składnik z requestu (pola, rodzic, makro, alergeny, przeliczniki) i przelicza makro
// wszystkich przepisów. Wołać w transakcji, po walidacji.
func SaveIngredient(tx *gorm.DB, ing *entities.Ingredient, req dtos.IngredientRequest, typeID uint) error {
	ing.Name = strings.TrimSpace(req.Name)
	ing.TypeID = typeID
	ing.ParentID = req.ParentID
	for _, k := range entities.NutrientKeys {
		if r, ok := req.Nutrition[k]; ok {
			ing.SetMacro(k, &r)
		} else {
			ing.SetMacro(k, nil)
		}
	}
	if err := tx.Omit("Type", "Parent", "Image", "Allergens", "Conversions").Save(ing).Error; err != nil {
		return err
	}

	if err := tx.Where("ingredient_id = ?", ing.ID).Delete(&entities.IngredientAllergen{}).Error; err != nil {
		return err
	}
	for _, code := range req.Allergens {
		if err := tx.Create(&entities.IngredientAllergen{IngredientID: ing.ID, Code: code}).Error; err != nil {
			return err
		}
	}

	if err := tx.Where("ingredient_id = ?", ing.ID).Delete(&entities.IngredientUnitConversion{}).Error; err != nil {
		return err
	}
	for _, c := range req.Conversions {
		row := entities.IngredientUnitConversion{IngredientID: ing.ID, UnitID: c.UnitID, GramsMin: c.GramsMin, GramsMax: c.GramsMax}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	return RecomputeRecipes(tx)
}

// FillIngredientNutrition uzupełnia DTO o makro własne, wynikowe i przeliczniki (z dziedziczeniem).
func FillIngredientNutrition(d *dtos.IngredientDto, ix *NutritionIndex) {
	d.Nutrition = map[string]entities.MacroRange(ix.Own[d.ID])
	if d.Nutrition == nil {
		d.Nutrition = map[string]entities.MacroRange{}
	}
	d.ResolvedNutrition = ix.ResolveMacros(d.ID)

	unitNames := map[uint]string{}
	for name, id := range ix.UnitIDByName {
		unitNames[id] = name
	}
	d.Conversions = []dtos.ConversionDto{}
	taken := map[uint]bool{}
	seen := map[uint]bool{}
	id := d.ID
	for cur := &id; cur != nil && !seen[*cur]; cur = ix.Parent[*cur] {
		seen[*cur] = true
		unitIDs := make([]uint, 0, len(ix.Conversions[*cur]))
		for u := range ix.Conversions[*cur] {
			unitIDs = append(unitIDs, u)
		}
		slices.Sort(unitIDs)
		for _, u := range unitIDs {
			if taken[u] {
				continue
			}
			taken[u] = true
			r := ix.Conversions[*cur][u]
			d.Conversions = append(d.Conversions, dtos.ConversionDto{
				UnitID: u, UnitName: unitNames[u], GramsMin: r.Min, GramsMax: r.Max, Inherited: *cur != d.ID,
			})
		}
	}
}
