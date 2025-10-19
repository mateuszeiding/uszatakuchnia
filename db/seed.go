package db

import (
	"fmt"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/services"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func runSeed(gdb *gorm.DB) error {
	tasteTypes := []entities.TasteType{
		{Code: "SWEET", Name: "Słodki"},
		{Code: "SOUR", Name: "Kwaśny"},
		{Code: "SALTY", Name: "Słony"},
		{Code: "BITTER", Name: "Gorzki"},
		{Code: "UMAMI", Name: "Umami"},
		{Code: "UNKNOWN", Name: "Nieznany"},
	}

	gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&tasteTypes)

	ingredientTypes := []entities.IngredientType{
		{Code: "VEGETABLE", Name: "Warzywo"},
		{Code: "FRUIT", Name: "Owoc"},
		{Code: "HERB", Name: "Zioło"},
		{Code: "SPICE", Name: "Przyprawa"},
		{Code: "MEAT", Name: "Mięso"},
		{Code: "FISH", Name: "Ryba"},
		{Code: "OTHER", Name: "Inne"},
	}

	gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&ingredientTypes)

	aromaTypes := []entities.AromaType{
		{Code: "FRUITY", Name: "Owocowy"},
		{Code: "HERBAL", Name: "Ziołowy"},
		{Code: "SPICY", Name: "Korzenny"},
		{Code: "EARTHY", Name: "Ziemisty"},
		{Code: "SMOKY", Name: "Dymny"},
		{Code: "CREAMY", Name: "Kremowy"},
		{Code: "OTHER", Name: "Inny"},
	}

	gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&aromaTypes)

	if err := SeedAll(gdb); err != nil {
		fmt.Println(err)
	}

	return nil
}

func loadDictMaps(tx *gorm.DB) (ingredientTypeID map[string]uint, aromaTypeID map[string]uint, tasteID map[string]uint, err error) {
	ingredientTypeID = map[string]uint{}
	var itypes []entities.IngredientType
	if err = tx.Find(&itypes).Error; err != nil {
		return
	}
	for _, t := range itypes {
		ingredientTypeID[t.Code] = t.ID
	}

	aromaTypeID = map[string]uint{}
	var atypes []entities.AromaType
	if err = tx.Find(&atypes).Error; err != nil {
		return
	}
	for _, a := range atypes {
		aromaTypeID[a.Code] = a.ID
	}

	tasteID = map[string]uint{}
	type row struct {
		ID   uint
		Code string
	}
	var rows []row
	if err = tx.Table("taste").
		Select("taste.id, taste_types.code").
		Joins("JOIN taste_types ON taste.type_id = taste_types.id").
		Scan(&rows).Error; err != nil {
		return
	}
	for _, r := range rows {
		tasteID[r.Code] = r.ID
	}

	return
}

func EnsureTasteForTypes(tx *gorm.DB) error {
	var ttypes []entities.TasteType
	if err := tx.Find(&ttypes).Error; err != nil {
		return err
	}
	for _, tt := range ttypes {
		t := entities.Taste{TypeID: tt.ID}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "type_id"}},
			DoNothing: true,
		}).Create(&t).Error; err != nil {
			return err
		}
	}
	return nil
}

type TasteSeed struct {
	TypeCode  string
	Intensity int
}

type AromaSeed struct {
	Name, TypeCode string
	Intensity      int
}
type ImageSeed struct {
	PostID  string
	RefType entities.RefType
}

func SeedOneIngredient(
	tx *gorm.DB,
	name string,
	typeCode string,
	isAllergen bool,
	parentName *string,
	tastes []TasteSeed,
	aromas []AromaSeed,
	image ImageSeed,
) error {
	ingTypeID, aromaTypeID, tasteID, err := loadDictMaps(tx)
	if err != nil {
		return err
	}

	tID, ok := ingTypeID[typeCode]
	if !ok {
		return fmt.Errorf("unknown ingredient type code: %s", typeCode)
	}

	var parentID *uint
	if parentName != nil {
		var parent entities.Ingredient
		if err := tx.Where("parent_id IS NULL AND name = ?", *parentName).First(&parent).Error; err != nil {
			return fmt.Errorf("parent not found: %w", err)
		}
		parentID = &parent.ID
	}

	// Przygotuj dane składnika
	ing := entities.Ingredient{}

	// Użyj Assign, aby zapewnić, że TypeID i IsAllergen są ustawione
	// zarówno przy tworzeniu, jak i przy aktualizacji.
	// FirstOrCreate znajdzie rekord po Name i ParentID lub utworzy nowy z pełnymi danymi.
	if err := tx.Where(entities.Ingredient{Name: name, ParentID: parentID}).
		Assign(entities.Ingredient{TypeID: tID, IsAllergen: isAllergen}).
		FirstOrCreate(&ing).Error; err != nil {
		return fmt.Errorf("failed to find or create ingredient: %w", err)
	}

	for _, t := range tastes {
		tid, ok := tasteID[t.TypeCode]
		if !ok {
			return fmt.Errorf("unknown taste type code: %s", t.TypeCode)
		}
		it := entities.IngredientTaste{
			IngredientID: ing.ID,
			TasteID:      tid,
			Intensity:    t.Intensity,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "ingredient_id"}, {Name: "taste_id"}},
			DoUpdates: clause.Assignments(map[string]any{"intensity": it.Intensity}),
		}).Create(&it).Error; err != nil {
			return err
		}
	}

	for _, a := range aromas {
		atid, ok := aromaTypeID[a.TypeCode]
		if !ok {
			return fmt.Errorf("unknown aroma type code: %s", a.TypeCode)
		}
		ar := entities.Aroma{
			IngredientID: ing.ID,
			Name:         a.Name,
			TypeID:       atid,
			Intensity:    a.Intensity,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "ingredient_id"}, {Name: "name"}},
			DoUpdates: clause.Assignments(map[string]any{
				"type_id":   ar.TypeID,
				"intensity": ar.Intensity,
			}),
		}).Create(&ar).Error; err != nil {
			return err
		}
	}

	if err := services.UpsertImageFromUnsplash(tx, ing.ID, image.RefType, image.PostID); err != nil {
		return err
	}

	return nil
}
func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		if err := EnsureTasteForTypes(tx); err != nil {
			return err
		}

		// Cebula
		if err := SeedOneIngredient(
			tx,
			"Cebula", "VEGETABLE", false, nil,
			[]TasteSeed{
				{TypeCode: "SWEET", Intensity: 20},
				{TypeCode: "UMAMI", Intensity: 10},
			},
			[]AromaSeed{
				{Name: "Ziemisty", TypeCode: "EARTHY", Intensity: 30},
				{Name: "Dymny", TypeCode: "SMOKY", Intensity: 20},
			},
			ImageSeed{PostID: "bC1fXU1v98U", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Czosnek
		if err := SeedOneIngredient(
			tx,
			"Czosnek", "VEGETABLE", false, nil,
			[]TasteSeed{
				{TypeCode: "UMAMI", Intensity: 25},
				{TypeCode: "BITTER", Intensity: 5},
			},
			[]AromaSeed{
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 40},
				{Name: "Ziemisty", TypeCode: "EARTHY", Intensity: 15},
			},
			ImageSeed{PostID: "_b8n1KTZ8rw", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Cytryna
		if err := SeedOneIngredient(
			tx,
			"Cytryna", "FRUIT", false, nil,
			[]TasteSeed{
				{TypeCode: "SOUR", Intensity: 95},
				{TypeCode: "BITTER", Intensity: 5},
			},
			[]AromaSeed{
				{Name: "Owocowy", TypeCode: "FRUITY", Intensity: 85},
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 10},
			},
			ImageSeed{PostID: "7WAGthfGJ9w", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Pomidor
		if err := SeedOneIngredient(
			tx,
			"Pomidor", "FRUIT", false, nil,
			[]TasteSeed{
				{TypeCode: "UMAMI", Intensity: 35},
				{TypeCode: "SWEET", Intensity: 20},
				{TypeCode: "SOUR", Intensity: 10},
			},
			[]AromaSeed{
				{Name: "Ziemisty", TypeCode: "EARTHY", Intensity: 25},
				{Name: "Owocowy", TypeCode: "FRUITY", Intensity: 35},
			},
			ImageSeed{PostID: "Jy7e7RjOFVo", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Bazylia
		if err := SeedOneIngredient(
			tx,
			"Bazylia", "HERB", false, nil,
			[]TasteSeed{
				{TypeCode: "BITTER", Intensity: 10},
			},
			[]AromaSeed{
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 80},
				{Name: "Korzenny", TypeCode: "SPICY", Intensity: 20},
			},
			ImageSeed{PostID: "DDKdfRe1GxA", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Mięta
		if err := SeedOneIngredient(
			tx,
			"Mięta", "HERB", false, nil,
			[]TasteSeed{
				{TypeCode: "BITTER", Intensity: 5},
				{TypeCode: "SWEET", Intensity: 5},
			},
			[]AromaSeed{
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 70},
				{Name: "Kremowy", TypeCode: "CREAMY", Intensity: 15},
			},
			ImageSeed{PostID: "s3C-iXNQIsQ", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Papryczka chili
		if err := SeedOneIngredient(
			tx,
			"Papryczka chili", "VEGETABLE", false, nil,
			[]TasteSeed{
				{TypeCode: "BITTER", Intensity: 10},
				{TypeCode: "SOUR", Intensity: 5},
			},
			[]AromaSeed{
				{Name: "Korzenny", TypeCode: "SPICY", Intensity: 70},
				{Name: "Dymny", TypeCode: "SMOKY", Intensity: 20},
			},
			ImageSeed{PostID: "iBsmi-wCXNE", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Pietruszka (nać)
		if err := SeedOneIngredient(
			tx,
			"Pietruszka (nać)", "HERB", false, nil,
			[]TasteSeed{
				{TypeCode: "BITTER", Intensity: 10},
				{TypeCode: "SWEET", Intensity: 5},
			},
			[]AromaSeed{
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 70},
				{Name: "Ziemisty", TypeCode: "EARTHY", Intensity: 30},
			},
			ImageSeed{PostID: "r6BUzN_jTHg", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Wołowina
		if err := SeedOneIngredient(
			tx,
			"Wołowina", "MEAT", false, nil,
			[]TasteSeed{
				{TypeCode: "UMAMI", Intensity: 60},
				{TypeCode: "BITTER", Intensity: 10},
			},
			[]AromaSeed{
				{Name: "Dymny", TypeCode: "SMOKY", Intensity: 40},
				{Name: "Ziemisty", TypeCode: "EARTHY", Intensity: 25},
			},
			ImageSeed{PostID: "u2QnkcLNsBY", RefType: "ingredient"},
		); err != nil {
			return err
		}

		// Kurczak
		if err := SeedOneIngredient(
			tx,
			"Kurczak", "MEAT", false, nil,
			[]TasteSeed{
				{TypeCode: "UMAMI", Intensity: 40},
				{TypeCode: "SWEET", Intensity: 10},
			},
			[]AromaSeed{
				{Name: "Kremowy", TypeCode: "CREAMY", Intensity: 30},
				{Name: "Ziołowy", TypeCode: "HERBAL", Intensity: 15},
			},
			ImageSeed{PostID: "pNcFMdEe09Q", RefType: "ingredient"},
		); err != nil {
			return err
		}

		return nil
	})
}
