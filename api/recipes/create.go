package api

import (
	"encoding/json"
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

func Create(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dtos.UpsertRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Status == "" {
		req.Status = entities.StatusDraft
	}

	conn := db.DB()

	fieldErrs, err := services.ValidateRecipe(conn, req)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(fieldErrs) > 0 {
		resp.JSON(w, http.StatusUnprocessableEntity, map[string]any{"error": "validation failed", "fields": fieldErrs})
		return
	}

	var recipe entities.Recipe
	services.ApplyRecipeRequest(&recipe, req)

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&recipe).Error; err != nil {
			return err
		}
		return services.ReplaceRecipeChildren(tx, recipe.ID, req)
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusCreated, map[string]any{"id": recipe.ID})
}
