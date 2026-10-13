package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"

	"gorm.io/gorm"
)

type updateIngredientRequest struct {
	Name      string   `json:"name"`
	TypeCode  string   `json:"typeCode"`
	Allergens []string `json:"allergens"`
}

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
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var req updateIngredientRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" || req.TypeCode == "" {
		resp.Error(w, http.StatusBadRequest, "name and typeCode are required")
		return
	}
	for _, code := range req.Allergens {
		if !entities.IsAllergenCode(code) {
			resp.Error(w, http.StatusBadRequest, "unknown allergen: "+code)
			return
		}
	}

	var ing entities.Ingredient
	if err := conn.First(&ing, id).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "ingredient not found")
		return
	}
	var dup int64
	conn.Model(&entities.Ingredient{}).
		Where("id <> ? AND parent_id IS NOT DISTINCT FROM ? AND lower(name) = lower(?)", ing.ID, ing.ParentID, name).
		Count(&dup)
	if dup > 0 {
		resp.Error(w, http.StatusConflict, "ingredient already exists")
		return
	}

	var ingType entities.IngredientType
	if err := conn.Where("code = ?", req.TypeCode).First(&ingType).Error; err != nil {
		resp.Error(w, http.StatusBadRequest, "unknown typeCode")
		return
	}

	if err := conn.Transaction(func(tx *gorm.DB) error {
		ing.Name = name
		ing.TypeID = ingType.ID
		if err := tx.Omit("Type", "Parent", "Image", "Allergens").Save(&ing).Error; err != nil {
			return err
		}
		if err := tx.Where("ingredient_id = ?", ing.ID).Delete(&entities.IngredientAllergen{}).Error; err != nil {
			return err
		}
		for _, code := range req.Allergens {
			if err := tx.Create(&entities.IngredientAllergen{IngredientID: ing.ID, Code: code}).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	conn.Preload("Type").Preload("Image").Preload("Allergens").First(&ing, ing.ID)
	resp.JSON(w, http.StatusOK, mappers.MapIngredientToDto(ing))
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

	if err := conn.Delete(&ing).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.JSON(w, http.StatusOK, map[string]any{"deleted": ing.ID})
}
