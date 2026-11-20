package services

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ValidateRecipe: a draft needs only a name and valid data; publishing needs everything
// (name, at least 1 ingredient, at least 1 step, required axes, difficulty, servings).
// selfID is the edited recipe's id (0 on create); an ingredient can't link to its own recipe.
func ValidateRecipe(tx *gorm.DB, req dtos.UpsertRecipeRequest, selfID uint) ([]dtos.FieldError, error) {
	errs := []dtos.FieldError{}
	add := func(field, msg string) { errs = append(errs, dtos.FieldError{Field: field, Message: msg}) }

	if req.Status != entities.StatusDraft && req.Status != entities.StatusPublished {
		add("status", "Nieznany status")
	}
	publish := req.Status == entities.StatusPublished

	if strings.TrimSpace(req.Name) == "" {
		add("name", "Podaj nazwę przepisu")
	}
	if req.Servings < 0 {
		add("servings", "Liczba porcji nie może być ujemna")
	}
	if req.TimeMin < 0 {
		add("time", "Czas nie może być ujemny")
	}
	if req.TimeMax != nil && *req.TimeMax < req.TimeMin {
		add("time", "Górna granica czasu musi być większa od dolnej")
	}
	if req.Difficulty != nil && (*req.Difficulty < 1 || *req.Difficulty > 3) {
		add("difficulty", "Trudność musi być w zakresie 1–3")
	}

	// Ingredients: a parsable quantity, an existing unit, "do smaku" without a quantity.
	units := map[uint]entities.Unit{}
	var unitRows []entities.Unit
	if err := tx.Find(&unitRows).Error; err != nil {
		return nil, err
	}
	for _, u := range unitRows {
		units[u.ID] = u
	}
	for i, it := range req.Ingredients {
		f := fmt.Sprintf("ingredients[%d]", i)
		if it.IngredientId == 0 {
			add(f+".ingredientId", "Wybierz składnik z bazy")
		}
		qty := ""
		if it.Quantity != nil {
			qty = *it.Quantity
		}
		q, ok := ParseQuantity(qty)
		if !ok {
			add(f+".quantity", "Niepoprawna ilość (np. 2, ½, 1,5, 2–3)")
		}
		if it.UnitId != nil {
			u, found := units[*it.UnitId]
			if !found {
				add(f+".unitId", "Nieznana jednostka")
			} else if u.NoQty && ok && !q.Empty {
				add(f+".quantity", fmt.Sprintf("Jednostka „%s” nie przyjmuje ilości", u.Name))
			}
		}
	}

	// Recipe links: not to itself, only to existing recipes.
	linked := []uint{}
	for _, it := range req.Ingredients {
		if it.LinkedRecipeId != nil {
			linked = append(linked, *it.LinkedRecipeId)
		}
	}
	existing := map[uint]bool{}
	if len(linked) > 0 {
		var ids []uint
		if err := tx.Model(&entities.Recipe{}).Where("id IN ?", uniqueUints(linked)).Pluck("id", &ids).Error; err != nil {
			return nil, err
		}
		for _, id := range ids {
			existing[id] = true
		}
	}
	for i, it := range req.Ingredients {
		if it.LinkedRecipeId == nil {
			continue
		}
		f := fmt.Sprintf("ingredients[%d].linkedRecipeId", i)
		if selfID != 0 && *it.LinkedRecipeId == selfID {
			add(f, "Przepis nie może prowadzić do siebie")
		} else if !existing[*it.LinkedRecipeId] {
			add(f, "Nie ma takiego przepisu")
		}
	}

	if req.Cautions != nil {
		validateCautions(*req.Cautions, add)
	}
	validateSteps(req, add)
	if err := validateSubstitutes(tx, req, add); err != nil {
		return nil, err
	}
	if err := validateEquipment(tx, req, add); err != nil {
		return nil, err
	}

	// Classification: values exist, single-value axes have at most 1 value.
	perAxis := map[uint]int{}
	if len(req.TaxonomyValueIDs) > 0 {
		var vals []entities.TaxonomyValue
		if err := tx.Where("id IN ?", req.TaxonomyValueIDs).Find(&vals).Error; err != nil {
			return nil, err
		}
		if len(vals) != len(uniqueUints(req.TaxonomyValueIDs)) {
			add("taxonomy", "Nieznana wartość klasyfikacji")
		}
		for _, v := range vals {
			perAxis[v.AxisID]++
		}
	}
	var axes []entities.TaxonomyAxis
	if err := tx.Find(&axes).Error; err != nil {
		return nil, err
	}
	for _, a := range axes {
		if !a.Multi && perAxis[a.ID] > 1 {
			add("taxonomy."+a.Code, fmt.Sprintf("%s: wybierz jedną wartość", a.Label))
		}
		if publish && a.Required && perAxis[a.ID] == 0 {
			add("taxonomy."+a.Code, fmt.Sprintf("%s: wybierz co najmniej jedną wartość", a.Label))
		}
	}

	if publish {
		if len(req.Ingredients) == 0 {
			add("ingredients", "Dodaj co najmniej jeden składnik")
		}
		hasStep := false
		for _, s := range req.Steps {
			if strings.TrimSpace(s.Text) != "" {
				hasStep = true
				break
			}
		}
		if !hasStep {
			add("steps", "Dodaj co najmniej jeden krok")
		}
		if req.Difficulty == nil {
			add("difficulty", "Wybierz trudność")
		}
		if req.Servings <= 0 {
			add("servings", "Podaj liczbę porcji")
		}
		if req.TimeMin <= 0 {
			add("time", "Podaj czas")
		}
	}

	return errs, nil
}

func uniqueUints(a []uint) []uint {
	seen := map[uint]bool{}
	out := []uint{}
	for _, x := range a {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
}

// ApplyRecipeRequest copies request fields onto the entity (without saving).
func ApplyRecipeRequest(e *entities.Recipe, req dtos.UpsertRecipeRequest) {
	e.Name = strings.TrimSpace(req.Name)
	e.Servings = req.Servings
	e.Description = req.Description
	e.Tagline = req.Tagline
	e.Tip = req.Tip
	e.TimeMin = req.TimeMin
	e.TimeMax = req.TimeMax
	e.Difficulty = req.Difficulty
	e.Status = req.Status
	if req.Cautions != nil {
		e.Cautions = CleanCautions(*req.Cautions)
	}
}

// ReplaceRecipeChildren replaces the recipe's photo, steps, ingredients and classification. Call inside a transaction.
func ReplaceRecipeChildren(tx *gorm.DB, recipeID uint, req dtos.UpsertRecipeRequest) error {
	if req.PhotoUrl != nil {
		if *req.PhotoUrl == "" {
			if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipePhoto{}).Error; err != nil {
				return err
			}
		} else {
			photo := entities.RecipePhoto{RecipeID: recipeID, ImageURL: *req.PhotoUrl}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "recipe_id"}},
				DoUpdates: clause.Assignments(map[string]any{"image_url": photo.ImageURL}),
			}).Create(&photo).Error; err != nil {
				return err
			}
		}
	}

	if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipeStep{}).Error; err != nil {
		return err
	}
	for i, s := range req.Steps {
		var lead *string
		if s.Prep {
			lead = s.PrepLead
		}
		step := entities.RecipeStep{
			RecipeID: recipeID, StepNo: i + 1, Section: s.Section, Text: s.Text, Prep: s.Prep, PrepLead: lead,
			Hints: CleanHints(s.Hints),
		}
		if err := tx.Create(&step).Error; err != nil {
			return err
		}
	}

	if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipeIngredient{}).Error; err != nil {
		return err
	}
	for i, it := range req.Ingredients {
		var qty *string
		if it.Quantity != nil && strings.TrimSpace(*it.Quantity) != "" {
			v := strings.TrimSpace(*it.Quantity)
			qty = &v
		}
		item := entities.RecipeIngredient{
			RecipeID: recipeID, SortOrder: i + 1, Section: it.Section,
			IngredientID: it.IngredientId, Quantity: qty, UnitID: it.UnitId, Note: it.Note,
			Optional: it.Optional, LinkedRecipeID: it.LinkedRecipeId,
		}
		if err := tx.Omit("Substitutes").Create(&item).Error; err != nil {
			return err
		}
		if err := createSubstitutes(tx, recipeID, item.SortOrder, it.Substitutes); err != nil {
			return err
		}
	}

	if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipeTaxonomy{}).Error; err != nil {
		return err
	}
	for _, vid := range uniqueUints(req.TaxonomyValueIDs) {
		if err := tx.Create(&entities.RecipeTaxonomy{RecipeID: recipeID, ValueID: vid}).Error; err != nil {
			return err
		}
	}
	return replaceRecipeEquipment(tx, recipeID, req)
}

// MaxHintLen is the longest step hint (characters, not bytes).
const MaxHintLen = 300

func validateSteps(req dtos.UpsertRecipeRequest, add func(field, msg string)) {
	for i, s := range req.Steps {
		if s.PrepLead != nil && !s.Prep {
			add(fmt.Sprintf("steps[%d].prepLead", i), "„Ile wcześniej” dotyczy tylko kroku przygotowania")
		}
		for _, h := range CleanHints(s.Hints) {
			if utf8.RuneCountInString(h) > MaxHintLen {
				add("steps", fmt.Sprintf("Krok %d: wskazówka może mieć najwyżej %d znaków", i+1, MaxHintLen))
				break
			}
		}
	}
}

// CleanHints trims hints and drops empty ones. Always returns a non-nil slice (stored as []).
func CleanHints(in []string) []string {
	out := make([]string, 0, len(in))
	for _, h := range in {
		if h = strings.TrimSpace(h); h != "" {
			out = append(out, h)
		}
	}
	return out
}
