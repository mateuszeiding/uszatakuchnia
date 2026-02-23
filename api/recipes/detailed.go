package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
)

func Detailed(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var entity entities.Recipe
	if err := conn.
		Preload("Steps").
		Preload("Ingredients").
		Find(&entity, id).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := mappers.MapRecipeToDto(entity)
	resp.JSON(w, http.StatusOK, dtos)
}
