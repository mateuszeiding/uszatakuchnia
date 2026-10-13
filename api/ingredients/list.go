package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
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

	// W ilu przepisach jest składnik — panel blokuje usuwanie używanych.
	var counts []struct {
		IngredientID uint
		N            int
	}
	if err := conn.Model(&entities.RecipeIngredient{}).
		Select("ingredient_id, COUNT(DISTINCT recipe_id) AS n").
		Group("ingredient_id").Scan(&counts).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	usedIn := map[uint]int{}
	for _, c := range counts {
		usedIn[c.IngredientID] = c.N
	}

	dtos := mappers.MapIngredientArrayToDto(list)
	for i := range dtos {
		dtos[i].UsedIn = usedIn[dtos[i].ID]
	}
	resp.JSON(w, http.StatusOK, dtos)
}
