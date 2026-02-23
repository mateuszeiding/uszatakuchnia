package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
)

func Ingredient(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()

	var list []entities.IngredientType
	conn.Find(&list)

	var dtos = mappers.MapIngredientTypeArrayToDto(list)
	resp.JSON(w, http.StatusOK, dtos)
}
