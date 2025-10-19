package ingredients

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
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
		Preload("Aromas").
		Preload("Aromas.Type").
		Preload("IngredientTastes.Taste.Type").
		Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := entities.IngredientList(list).ListToDto()

	resp.JSON(w, http.StatusOK, dtos)
}
