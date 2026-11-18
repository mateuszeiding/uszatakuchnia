package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"
)

func List(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()

	var list []entities.Ingredient
	if err := conn.
		Preload("Type").
		Preload("Image").
		Preload("Allergens").
		Order("name").
		Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	// W ilu przepisach jest składnik (też jako zamiennik) — panel blokuje usuwanie używanych.
	usedIn, err := services.IngredientUsage(conn)
	if err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	ix, err := services.LoadNutritionIndex(conn)
	if err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := mappers.MapIngredientArrayToDto(list)
	for i := range dtos {
		dtos[i].UsedIn = usedIn[dtos[i].ID]
		services.FillIngredientNutrition(&dtos[i], ix)
	}
	resp.JSON(w, http.StatusOK, dtos)
}
