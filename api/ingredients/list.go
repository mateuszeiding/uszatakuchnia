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
		Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := mappers.MapIngredientArrayToDto(list)
	resp.JSON(w, http.StatusOK, dtos)
}
