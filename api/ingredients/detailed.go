package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

// Detailed — PUT/DELETE /api/ingredients/:id
func Detailed(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	switch r.Method {
	case http.MethodPut:
		putIngredient(w, r)
	case http.MethodDelete:
		deleteIngredient(w, r)
	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func putIngredient(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid id")
		return
	}
	var req dtos.IngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	res := services.UpsertIngredient(db.DB(), uint(id), req, http.StatusOK)
	resp.JSON(w, res.Status, res.Body)
}

// deleteIngredient — używanego w przepisach składnika nie da się usunąć (409 + liczba przepisów).
func deleteIngredient(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var ing entities.Ingredient
	if err := conn.First(&ing, id).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "ingredient not found")
		return
	}

	var used int64
	conn.Model(&entities.RecipeIngredient{}).Where("ingredient_id = ?", ing.ID).
		Distinct("recipe_id").Count(&used)
	if used > 0 {
		resp.JSON(w, http.StatusConflict, map[string]any{
			"error":  fmt.Sprintf("ingredient used in %d recipes", used),
			"usedIn": used,
		})
		return
	}

	// Usunięcie zmienia drzewko (dziedziczone makro i przeliczniki) — przelicz makro przepisów.
	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&ing).Error; err != nil {
			return err
		}
		return services.RecomputeRecipes(tx)
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.JSON(w, http.StatusOK, map[string]any{"deleted": ing.ID})
}
