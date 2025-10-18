package db

import (
	"uszatakuchnia/db/entities"

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

	return nil
}
