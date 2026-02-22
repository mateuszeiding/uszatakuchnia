package db

import (
	"fmt"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/services"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func runSeed(gdb *gorm.DB) error {
	ingredientTypes := []entities.IngredientType{
		{Code: "VEGETABLE", Name: "Warzywo"},
		{Code: "FRUIT", Name: "Owoc"},
		{Code: "HERB", Name: "Zioło"},
		{Code: "SPICE", Name: "Przyprawa"},
		{Code: "MEAT", Name: "Mięso"},
		{Code: "FISH", Name: "Ryba"},
		{Code: "OTHER", Name: "Inne"},
	}

	if err := gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&ingredientTypes).Error; err != nil {
		return err
	}

	if err := SeedAll(gdb); err != nil {
		return err
	}

	return nil
}

// --------------------
// Ingredients
// --------------------

type IngredientImageSeed struct {
	PostID string // Unsplash photo id
}

type IngredientSeed struct {
	Name       string
	TypeCode   string
	IsAllergen bool
	ParentName *string
	Image      *IngredientImageSeed
}

func loadIngredientTypeMap(tx *gorm.DB) (map[string]uint, error) {
	out := map[string]uint{}
	var rows []entities.IngredientType
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.Code] = r.ID
	}
	return out, nil
}

func findRootIngredientID(tx *gorm.DB, name string) (*uint, error) {
	var ing entities.Ingredient
	if err := tx.Where("parent_id IS NULL AND name = ?", name).First(&ing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("parent not found (root ingredient): %s", name)
		}
		return nil, err
	}
	return &ing.ID, nil
}

func SeedOneIngredient(tx *gorm.DB, seed IngredientSeed) (uint, error) {
	typeMap, err := loadIngredientTypeMap(tx)
	if err != nil {
		return 0, err
	}

	typeID, ok := typeMap[seed.TypeCode]
	if !ok {
		return 0, fmt.Errorf("unknown ingredient type code: %s", seed.TypeCode)
	}

	var parentID *uint
	if seed.ParentName != nil {
		pid, err := findRootIngredientID(tx, *seed.ParentName)
		if err != nil {
			return 0, err
		}
		parentID = pid
	}

	var ing entities.Ingredient
	if err := tx.Where(entities.Ingredient{Name: seed.Name, ParentID: parentID}).
		Assign(entities.Ingredient{TypeID: typeID, IsAllergen: seed.IsAllergen}).
		FirstOrCreate(&ing).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create ingredient: %w", err)
	}

	// Unsplash image (opcjonalny) — tylko dla ingredientów
	if seed.Image != nil {
		if err := services.UpsertImageFromUnsplash(tx, ing.ID, entities.IngredientRef, seed.Image.PostID); err != nil {
			return 0, err
		}
	}

	return ing.ID, nil
}

// --------------------
// Recipes
// --------------------

type RecipePhotoSeed struct {
	URL string // publiczny URL z Vercel Blob (opcjonalnie)
}

type RecipeSeed struct {
	Name        string
	Servings    int
	Description *string
	Photo       *RecipePhotoSeed

	Steps []RecipeStepSeed
	Items []RecipeIngredientSeed
}

type RecipeStepSeed struct {
	StepNo  int
	Section *string
	Text    string
}

type RecipeIngredientSeed struct {
	SortOrder int
	Section   *string

	IngredientName       string
	IngredientParentName *string

	Amount     *float64
	Unit       *string
	AmountText *string
	Note       *string
}

// kompatybilne z SQLite/MySQL/Postgres: rozgałęzienie po parentID == nil
func findIngredientID(tx *gorm.DB, name string, parentName *string) (uint, error) {
	var parentID *uint
	if parentName != nil {
		pid, err := findRootIngredientID(tx, *parentName)
		if err != nil {
			return 0, err
		}
		parentID = pid
	}

	var ing entities.Ingredient
	q := tx.Where("name = ?", name)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}

	if err := q.First(&ing).Error; err != nil {
		return 0, fmt.Errorf("ingredient not found: %s (parent=%v): %w", name, parentName, err)
	}

	return ing.ID, nil
}

func upsertRecipePhoto(tx *gorm.DB, recipeID uint, url string) error {
	photo := entities.RecipePhoto{
		RecipeID: recipeID,
		ImageURL: url,
	}

	// PK = recipe_id, więc upsert po recipe_id
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "recipe_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"image_url": photo.ImageURL,
		}),
	}).Create(&photo).Error
}

func SeedOneRecipe(tx *gorm.DB, seed RecipeSeed) (uint, error) {
	if seed.Servings <= 0 {
		return 0, fmt.Errorf("recipe servings must be > 0: %s", seed.Name)
	}

	var r entities.Recipe
	if err := tx.Where(entities.Recipe{Name: seed.Name}).
		Assign(entities.Recipe{Servings: seed.Servings, Description: seed.Description}).
		FirstOrCreate(&r).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create recipe: %w", err)
	}

	// Photo (opcjonalne) — zapis do recipe_photos, NIE do images/unsplash
	if seed.Photo != nil && seed.Photo.URL != "" {
		if err := upsertRecipePhoto(tx, r.ID, seed.Photo.URL); err != nil {
			return 0, err
		}
	}

	// Steps: PK (recipe_id, step_no)
	for _, s := range seed.Steps {
		step := entities.RecipeStep{
			RecipeID: r.ID,
			StepNo:   s.StepNo,
			Section:  s.Section,
			Text:     s.Text,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "recipe_id"}, {Name: "step_no"}},
			DoUpdates: clause.Assignments(map[string]any{
				"section": step.Section,
				"text":    step.Text,
			}),
		}).Create(&step).Error; err != nil {
			return 0, err
		}
	}

	// Ingredients: PK (recipe_id, sort_order)
	for _, it := range seed.Items {
		ingID, err := findIngredientID(tx, it.IngredientName, it.IngredientParentName)
		if err != nil {
			return 0, err
		}

		item := entities.RecipeIngredient{
			RecipeID:     r.ID,
			SortOrder:    it.SortOrder,
			Section:      it.Section,
			IngredientID: ingID,
			Amount:       it.Amount,
			Unit:         it.Unit,
			AmountText:   it.AmountText,
			Note:         it.Note,
		}

		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "recipe_id"}, {Name: "sort_order"}},
			DoUpdates: clause.Assignments(map[string]any{
				"section":       item.Section,
				"ingredient_id": item.IngredientID,
				"amount":        item.Amount,
				"unit":          item.Unit,
				"amount_text":   item.AmountText,
				"note":          item.Note,
			}),
		}).Create(&item).Error; err != nil {
			return 0, err
		}
	}

	return r.ID, nil
}

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		// ---------- Ingredients ----------
		ings := []IngredientSeed{
			{Name: "Cebula", TypeCode: "VEGETABLE", IsAllergen: false, Image: &IngredientImageSeed{PostID: "bC1fXU1v98U"}},
			{Name: "Czosnek", TypeCode: "VEGETABLE", IsAllergen: false, Image: &IngredientImageSeed{PostID: "_b8n1KTZ8rw"}},
			{Name: "Cytryna", TypeCode: "FRUIT", IsAllergen: false, Image: &IngredientImageSeed{PostID: "7WAGthfGJ9w"}},
			{Name: "Pomidor", TypeCode: "FRUIT", IsAllergen: false, Image: &IngredientImageSeed{PostID: "Jy7e7RjOFVo"}},
			{Name: "Bazylia", TypeCode: "HERB", IsAllergen: false, Image: &IngredientImageSeed{PostID: "DDKdfRe1GxA"}},
			{Name: "Mięta", TypeCode: "HERB", IsAllergen: false, Image: &IngredientImageSeed{PostID: "s3C-iXNQIsQ"}},
			{Name: "Papryczka chili", TypeCode: "VEGETABLE", IsAllergen: false, Image: &IngredientImageSeed{PostID: "iBsmi-wCXNE"}},
			{Name: "Pietruszka (nać)", TypeCode: "HERB", IsAllergen: false, Image: &IngredientImageSeed{PostID: "r6BUzN_jTHg"}},
			{Name: "Wołowina", TypeCode: "MEAT", IsAllergen: false, Image: &IngredientImageSeed{PostID: "u2QnkcLNsBY"}},
			{Name: "Kurczak", TypeCode: "MEAT", IsAllergen: false, Image: &IngredientImageSeed{PostID: "pNcFMdEe09Q"}},
		}

		for _, s := range ings {
			if _, err := SeedOneIngredient(tx, s); err != nil {
				return err
			}
		}

		// ---------- Recipe example ----------
		secSauce := "Sos"
		secMain := "Główne"

		desc := "Prosty przepis startowy pod MVP — do sprawdzenia UI sekcji i kolejności."
		recipe := RecipeSeed{
			Name:        "Pomidorowy sos z czosnkiem i bazylią",
			Servings:    2,
			Description: &desc,
			Photo:       nil, // zostaw nil — zdjęcia z Vercel Blob ogarniesz później
			Steps: []RecipeStepSeed{
				{StepNo: 1, Section: &secSauce, Text: "Posiekaj czosnek, podsmaż krótko na tłuszczu."},
				{StepNo: 2, Section: &secSauce, Text: "Dodaj pomidory i duś kilka minut."},
				{StepNo: 3, Section: &secSauce, Text: "Dopraw, na koniec dorzuć bazylię."},
			},
			Items: []RecipeIngredientSeed{
				{SortOrder: 1, Section: &secSauce, IngredientName: "Czosnek", AmountText: strPtr("2 ząbki")},
				{SortOrder: 2, Section: &secSauce, IngredientName: "Pomidor", AmountText: strPtr("3–4 szt.")},
				{SortOrder: 3, Section: &secSauce, IngredientName: "Bazylia", AmountText: strPtr("kilka listków")},
				{SortOrder: 4, Section: &secMain, IngredientName: "Cebula", AmountText: strPtr("1/2 szt."), Note: strPtr("opcjonalnie")},
			},
		}

		if _, err := SeedOneRecipe(tx, recipe); err != nil {
			return err
		}

		return nil
	})
}

func strPtr(s string) *string { return &s }
