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

	// admin=true query param means show all; otherwise only published
	isAdmin := r.URL.Query().Get("admin") == "true"

	query := conn.Model(&entities.Recipe{}).
		Select("id, name, tagline, category, region, time_minutes, difficulty, status, needs_prep").
		Preload("Photo").
		Preload("Tags")

	if !isAdmin {
		query = query.Where("status = ?", "published")
	}

	var list []entities.RecipeBase
	if err := query.Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := mappers.MapRecipeBaseArrayToDto(list)
	resp.JSON(w, http.StatusOK, dtos)
}
