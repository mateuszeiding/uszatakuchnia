package services

import (
	"uszatakuchnia/db/entities"

	"gorm.io/gorm"
)

// LoadNutritionIndex wczytuje wszystkie składniki, przeliczniki i jednostki (tabele są małe).
func LoadNutritionIndex(tx *gorm.DB) (*NutritionIndex, error) {
	var ings []entities.Ingredient
	if err := tx.Omit("Type", "Parent", "Image").Find(&ings).Error; err != nil {
		return nil, err
	}
	var convs []entities.IngredientUnitConversion
	if err := tx.Find(&convs).Error; err != nil {
		return nil, err
	}
	var units []entities.Unit
	if err := tx.Find(&units).Error; err != nil {
		return nil, err
	}

	ix := &NutritionIndex{
		Own:          map[uint]Macros{},
		Children:     map[uint][]uint{},
		Parent:       map[uint]*uint{},
		Conversions:  map[uint]map[uint]entities.MacroRange{},
		Names:        map[uint]string{},
		UnitIDByName: map[string]uint{},
	}
	for _, ing := range ings {
		ix.Names[ing.ID] = ing.Name
		ix.Parent[ing.ID] = ing.ParentID
		if ing.ParentID != nil {
			ix.Children[*ing.ParentID] = append(ix.Children[*ing.ParentID], ing.ID)
		}
		own := Macros{}
		for _, k := range entities.NutrientKeys {
			if r := ing.MacroOf(k); r != nil {
				own[k] = *r
			}
		}
		if len(own) > 0 {
			ix.Own[ing.ID] = own
		}
	}
	for _, c := range convs {
		if ix.Conversions[c.IngredientID] == nil {
			ix.Conversions[c.IngredientID] = map[uint]entities.MacroRange{}
		}
		ix.Conversions[c.IngredientID][c.UnitID] = entities.MacroRange{Min: c.GramsMin, Max: c.GramsMax}
	}
	for _, u := range units {
		ix.UnitIDByName[u.Name] = u.ID
	}
	return ix, nil
}

// RecomputeRecipes przelicza i zapisuje recipes.nutrition. Bez argumentów — wszystkie przepisy
// (zmiana składnika wpływa na przodków i potomków w drzewie, a przepisów jest niewiele).
func RecomputeRecipes(tx *gorm.DB, recipeIDs ...uint) error {
	ix, err := LoadNutritionIndex(tx)
	if err != nil {
		return err
	}

	q := tx.Model(&entities.Recipe{}).Select("id", "servings").
		Preload("Ingredients", func(db *gorm.DB) *gorm.DB { return db.Order("sort_order") }).
		Preload("Ingredients.Unit")
	if len(recipeIDs) > 0 {
		q = q.Where("id IN ?", recipeIDs)
	}
	var recipes []entities.Recipe
	if err := q.Find(&recipes).Error; err != nil {
		return err
	}

	for _, r := range recipes {
		lines := make([]NutritionLine, 0, len(r.Ingredients))
		for _, it := range r.Ingredients {
			raw := ""
			if it.Quantity != nil {
				raw = *it.Quantity
			}
			qty, ok := ParseQuantity(raw)
			if !ok {
				qty = Quantity{Empty: true}
			}
			l := NutritionLine{IngredientID: it.IngredientID, Qty: qty, UnitID: it.UnitID, Optional: it.Optional}
			if it.Unit != nil {
				l.UnitName = it.Unit.Name
				l.NoQty = it.Unit.NoQty
			}
			lines = append(lines, l)
		}
		n := ix.Compute(lines, r.Servings)
		if err := tx.Model(&entities.Recipe{}).Where("id = ?", r.ID).Update("nutrition", n).Error; err != nil {
			return err
		}
	}
	return nil
}
