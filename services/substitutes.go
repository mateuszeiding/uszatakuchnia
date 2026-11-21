package services

import (
	"fmt"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// validateSubstitutes: the substitute exists, differs from the main ingredient, no duplicates in the row.
func validateSubstitutes(tx *gorm.DB, req dtos.UpsertRecipeRequest, add func(field, msg string)) error {
	ids := []uint{}
	for _, it := range req.Ingredients {
		for _, s := range it.Substitutes {
			ids = append(ids, s.IngredientId)
		}
	}
	if len(ids) == 0 {
		return nil
	}
	known := map[uint]bool{}
	var found []uint
	if err := tx.Model(&entities.Ingredient{}).Where("id IN ?", uniqueUints(ids)).Pluck("id", &found).Error; err != nil {
		return err
	}
	for _, id := range found {
		known[id] = true
	}
	type option struct {
		id  uint
		mix bool
	}
	for i, it := range req.Ingredients {
		seen := map[option]bool{}
		for j, s := range it.Substitutes {
			f := fmt.Sprintf("ingredients[%d].substitutes[%d]", i, j)
			o := option{s.IngredientId, s.Mix}
			switch {
			case !known[s.IngredientId]:
				add(f, "Nie ma takiego składnika")
			case s.IngredientId == it.IngredientId:
				add(f, "Zamiennik musi być innym składnikiem niż główny")
			case seen[o]:
				add(f, "Ten zamiennik już jest przy składniku")
			}
			seen[o] = true
		}
	}
	return nil
}

// createSubstitutes saves a row's substitutes (the row exists already; the old ones went with the cascade when rows were replaced).
func createSubstitutes(tx *gorm.DB, recipeID uint, sortOrder int, subs []dtos.UpsertSubstituteRequest) error {
	for k, s := range subs {
		row := entities.RecipeIngredientSubstitute{
			RecipeID: recipeID, SortOrder: sortOrder, Position: k + 1, IngredientID: s.IngredientId, Mix: s.Mix,
		}
		if err := tx.Omit("Ingredient").Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

// nutritionLineOf turns a recipe row into a macro calculation line (needs Unit preloaded).
func nutritionLineOf(it entities.RecipeIngredient) NutritionLine {
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
	return l
}

// FillLineNutrition sets the macro share of every row and every substitute (mix = half main + half
// substitute), plus grams per unit and "has macro data" for substitutes. items and e.Ingredients in the same order.
func FillLineNutrition(items []dtos.RecipeIngredientDto, e entities.Recipe, ix *NutritionIndex) {
	for i := range items {
		if i >= len(e.Ingredients) {
			return
		}
		main := nutritionLineOf(e.Ingredients[i])
		items[i].Nutrition = ix.LineNutrition([]NutritionLine{main}, e.Servings)
		for k := range items[i].Substitutes {
			sub := &items[i].Substitutes[k]
			alt := main
			alt.IngredientID = sub.IngredientID
			lines := []NutritionLine{alt}
			if sub.Mix {
				lines = []NutritionLine{main.Half(), alt.Half()}
			}
			sub.Nutrition = ix.LineNutrition(lines, e.Servings)
			sub.HasNutrition = len(ix.ResolveMacros(sub.IngredientID)) > 0
			sub.GramsPerUnit = nil
			if main.NoQty {
				continue
			}
			if g, ok := ix.GramsPerUnit(sub.IngredientID, main.UnitID, main.UnitName); ok {
				sub.GramsPerUnit = &dtos.GramsRangeDto{Min: g.Min, Max: g.Upper()}
			}
		}
	}
}

// IngredientUsage: number of recipes using each ingredient (as an ingredient or a substitute).
func IngredientUsage(tx *gorm.DB) (map[uint]int, error) {
	var rows []struct {
		IngredientID uint
		N            int
	}
	if err := tx.Raw(`SELECT ingredient_id, COUNT(DISTINCT recipe_id) AS n FROM (
			SELECT ingredient_id, recipe_id FROM recipe_ingredients
			UNION ALL
			SELECT ingredient_id, recipe_id FROM recipe_ingredient_substitutes
		) u GROUP BY ingredient_id`).Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := map[uint]int{}
	for _, r := range rows {
		out[r.IngredientID] = r.N
	}
	return out, nil
}
