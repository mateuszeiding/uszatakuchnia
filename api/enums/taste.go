package enums

import (
	"encoding/json"
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	dbConn := db.DB()

	var tastes []entities.TasteType
	if err := dbConn.Find(&tastes).Error; err != nil {
		http.Error(w, "Failed to fetch tastes", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tastes)
}
