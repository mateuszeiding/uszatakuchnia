package services

import (
	"fmt"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// validateSubstitutes — zamiennik istnieje w bazie, różni się od składnika głównego, bez duplikatów w wierszu.
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

// createSubstitutes zapisuje zamienniki wiersza (wiersz już istnieje; stare usunęła kaskada przy zastąpieniu wierszy).
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

// nutritionLineOf — pozycja przepisu jako linia do wyliczeń makro (wymaga preloadu Unit).
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

// FillLineNutrition — wkład makro każdej pozycji i każdego zamiennika (mieszanka = połowa głównego + połowa
// zamiennika) oraz gramy na jednostkę i „ma dane makro” dla zamienników. items i e.Ingredients w tej samej kolejności.
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

// IngredientUsage — w ilu przepisach użyty każdy składnik (jako składnik albo zamiennik).
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
