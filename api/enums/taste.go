package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	resp "uszatakuchnia/http"
)

func Taste(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()

	var list []entities.TasteType
	conn.Find(&list)
	resp.JSON(w, http.StatusOK, list)
}
