package services

import (
	"fmt"
	"strings"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ValidateRecipe — szkic wymaga tylko nazwy i poprawnych danych; publikacja wymaga kompletu
// (nazwa, min. 1 składnik, min. 1 krok, wymagane osie, trudność, porcje).
func ValidateRecipe(tx *gorm.DB, req dtos.UpsertRecipeRequest) ([]dtos.FieldError, error) {
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
	if req.Difficulty != nil && (*req.Difficulty < 1 || *req.Difficulty > 3) {
		add("difficulty", "Trudność musi być w zakresie 1–3")
	}

	// Składniki: ilość parsowalna, jednostka istnieje, "do smaku" bez ilości.
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

	for i, s := range req.Steps {
		if s.PrepLead != nil && !s.Prep {
			add(fmt.Sprintf("steps[%d].prepLead", i), "„Ile wcześniej” dotyczy tylko kroku przygotowania")
		}
	}

	// Klasyfikacja: wartości istnieją, osie pojedyncze mają max 1 wartość.
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

// ApplyRecipeRequest przepisuje pola z requestu na encję (bez zapisu).
func ApplyRecipeRequest(e *entities.Recipe, req dtos.UpsertRecipeRequest) {
	e.Name = strings.TrimSpace(req.Name)
	e.Servings = req.Servings
	e.Description = req.Description
	e.Tagline = req.Tagline
	e.Tip = req.Tip
	e.TimeMinutes = req.TimeMinutes
	e.Difficulty = req.Difficulty
	e.KcalPerServing = req.KcalPerServing
	e.Status = req.Status
}

// ReplaceRecipeChildren zastępuje zdjęcie, kroki, składniki i klasyfikację przepisu. Wołać w transakcji.
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
			Optional: it.Optional,
		}
		if err := tx.Create(&item).Error; err != nil {
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
	return nil
}
