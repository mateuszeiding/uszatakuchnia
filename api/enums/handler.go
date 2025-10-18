package enums

import (
	"net/http"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	resp "uszatakuchnia/http"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	path := strings.TrimPrefix(r.URL.Path, "/api/enums/")

	switch path {
	case "taste":
		var list []entities.TasteType
		conn.Find(&list)
		resp.JSON(w, http.StatusOK, list)

	case "aroma":
		var list []entities.AromaType
		conn.Find(&list)
		resp.JSON(w, http.StatusOK, list)

	case "ingredient":
		var list []entities.IngredientType
		conn.Find(&list)
		resp.JSON(w, http.StatusOK, list)

	default:
		http.Error(w, "unknown enum", http.StatusNotFound)
	}
}
