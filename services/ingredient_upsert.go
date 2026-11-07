package services

import (
	"net/http"
	"strings"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// UpsertResult — wynik UpsertIngredient: Status HTTP i ciało odpowiedzi.
type UpsertResult struct {
	Status int
	Body   any
}

// UpsertIngredient — wspólna ścieżka POST (id 0) i PUT /api/ingredients: walidacja, zapis w transakcji,
// przeliczenie makro przepisów, pełne DTO w odpowiedzi.
func UpsertIngredient(conn *gorm.DB, id uint, req dtos.IngredientRequest, okStatus int) UpsertResult {
	fail := func(status int, msg string) UpsertResult { return UpsertResult{status, map[string]any{"error": msg}} }

	name := strings.TrimSpace(req.Name)
	if name == "" || req.TypeCode == "" {
		return fail(http.StatusBadRequest, "name and typeCode are required")
	}
	for _, code := range req.Allergens {
		if !entities.IsAllergenCode(code) {
			return fail(http.StatusBadRequest, "unknown allergen: "+code)
		}
	}

	var ing entities.Ingredient
	if id != 0 {
		if err := conn.First(&ing, id).Error; err != nil {
			return fail(http.StatusNotFound, "ingredient not found")
		}
	}
	var ingType entities.IngredientType
	if err := conn.Where("code = ?", req.TypeCode).First(&ingType).Error; err != nil {
		return fail(http.StatusBadRequest, "unknown typeCode")
	}

	ix, err := LoadNutritionIndex(conn)
	if err != nil {
		return fail(http.StatusInternalServerError, err.Error())
	}
	if fieldErrs := ValidateIngredientRequest(req, ix, id); len(fieldErrs) > 0 {
		return UpsertResult{http.StatusUnprocessableEntity, map[string]any{"error": "validation failed", "fields": fieldErrs}}
	}

	var dup int64
	conn.Model(&entities.Ingredient{}).
		Where("id <> ? AND parent_id IS NOT DISTINCT FROM ? AND lower(name) = lower(?)", id, req.ParentID, name).
		Count(&dup)
	if dup > 0 {
		return fail(http.StatusConflict, "ingredient already exists")
	}

	if err := conn.Transaction(func(tx *gorm.DB) error {
		return SaveIngredient(tx, &ing, req, ingType.ID)
	}); err != nil {
		return fail(http.StatusInternalServerError, err.Error())
	}

	conn.Preload("Type").Preload("Image").Preload("Allergens").First(&ing, ing.ID)
	if ix, err = LoadNutritionIndex(conn); err != nil {
		return fail(http.StatusInternalServerError, err.Error())
	}
	out := mappers.MapIngredientToDto(ing)
	FillIngredientNutrition(&out, ix)
	return UpsertResult{okStatus, out}
}
