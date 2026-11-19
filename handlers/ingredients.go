package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

// IngredientList — GET /api/ingredients/list.
func IngredientList(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()

	var list []entities.Ingredient
	if err := conn.
		Preload("Type").
		Preload("Image").
		Preload("Allergens").
		Order("name").
		Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	// W ilu przepisach jest składnik (też jako zamiennik) — panel blokuje usuwanie używanych.
	usedIn, err := services.IngredientUsage(conn)
	if err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	ix, err := services.LoadNutritionIndex(conn)
	if err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	dtos := mappers.MapIngredientArrayToDto(list)
	for i := range dtos {
		dtos[i].UsedIn = usedIn[dtos[i].ID]
		services.FillIngredientNutrition(&dtos[i], ix)
	}
	resp.JSON(w, http.StatusOK, dtos)
}

// IngredientCreate — POST /api/ingredients.
func IngredientCreate(w http.ResponseWriter, r *http.Request) {
	var req dtos.IngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	res := services.UpsertIngredient(db.DB(), 0, req, http.StatusCreated)
	resp.JSON(w, res.Status, res.Body)
}

// IngredientUpdate — PUT /api/ingredients/{id}.
func IngredientUpdate(w http.ResponseWriter, r *http.Request) {
	var req dtos.IngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	res := services.UpsertIngredient(db.DB(), pathID(r), req, http.StatusOK)
	resp.JSON(w, res.Status, res.Body)
}

// IngredientDelete — DELETE /api/ingredients/{id}. Używanego w przepisach składnika nie da się
// usunąć (409 + liczba przepisów).
func IngredientDelete(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()

	var ing entities.Ingredient
	if err := conn.First(&ing, pathID(r)).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "ingredient not found")
		return
	}

	usage, err := services.IngredientUsage(conn)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	used := usage[ing.ID] // też jako zamiennik
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

// IngredientTypes — GET /api/enums/ingredient.
func IngredientTypes(w http.ResponseWriter, r *http.Request) {
	var list []entities.IngredientType
	db.DB().Find(&list)
	resp.JSON(w, http.StatusOK, mappers.MapIngredientTypeArrayToDto(list))
}
