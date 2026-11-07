package api

import (
	"encoding/json"
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"
)

func Create(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dtos.IngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	res := services.UpsertIngredient(db.DB(), 0, req, http.StatusCreated)
	resp.JSON(w, res.Status, res.Body)
}
