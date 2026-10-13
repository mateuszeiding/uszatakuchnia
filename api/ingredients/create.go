package api

import (
	"encoding/json"
	"net/http"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
)

type createIngredientRequest struct {
	Name      string   `json:"name"`
	TypeCode  string   `json:"typeCode"`
	Allergens []string `json:"allergens"`
}

func Create(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req createIngredientRequest
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

	conn := db.DB()

	var ingType entities.IngredientType
	if err := conn.Where("code = ?", req.TypeCode).First(&ingType).Error; err != nil {
		resp.Error(w, http.StatusBadRequest, "unknown typeCode")
		return
	}

	var dup int64
	conn.Model(&entities.Ingredient{}).Where("parent_id IS NULL AND lower(name) = lower(?)", name).Count(&dup)
	if dup > 0 {
		resp.Error(w, http.StatusConflict, "ingredient already exists")
		return
	}

	ing := entities.Ingredient{Name: name, TypeID: ingType.ID}
	for _, code := range req.Allergens {
		ing.Allergens = append(ing.Allergens, entities.IngredientAllergen{Code: code})
	}
	// Create zapisuje też Allergens (asocjacja has-many) w jednej transakcji.
	if err := conn.Create(&ing).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	conn.Preload("Type").Preload("Image").Preload("Allergens").First(&ing, ing.ID)
	resp.JSON(w, http.StatusCreated, mappers.MapIngredientToDto(ing))
}
