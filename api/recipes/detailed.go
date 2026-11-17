package api

import (
	"encoding/json"
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func Detailed(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}

	switch r.Method {
	case http.MethodGet:
		getDetailed(w, r)
	case http.MethodPut:
		putDetailed(w, r)
	case http.MethodDelete:
		deleteDetailed(w, r)
	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func getDetailed(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var entity entities.Recipe
	if err := conn.
		Preload("Photo").
		Preload("Steps", func(tx *gorm.DB) *gorm.DB { return tx.Order("step_no") }).
		Preload("Ingredients", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order") }).
		Preload("Ingredients.Ingredient.Allergens").
		Preload("Ingredients.Ingredient.Type").
		Preload("Ingredients.Unit").
		Preload("Ingredients.LinkedRecipe.Photo").
		Preload("Taxonomy.Value.Axis").
		Find(&entity, id).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if entity.ID == 0 {
		resp.Error(w, http.StatusNotFound, "recipe not found")
		return
	}

	dto := mappers.MapRecipeToDto(entity)
	ix, err := services.LoadNutritionIndex(conn)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	services.FillGramsPerUnit(dto.Ingredients, ix)

	// Kto prowadzi do tego przepisu (modal usuwania): „Przepis · składnik”.
	dto.LinkedFrom = []dtos.RecipeLinkSourceDto{}
	if err := conn.Table("recipe_ingredients ri").
		Select("r.id AS recipe_id, r.name AS recipe_name, i.name AS ingredient_name").
		Joins("JOIN recipes r ON r.id = ri.recipe_id").
		Joins("JOIN ingredients i ON i.id = ri.ingredient_id").
		Where("ri.linked_recipe_id = ? AND ri.recipe_id <> ?", entity.ID, entity.ID).
		Order("r.name, ri.sort_order").
		Scan(&dto.LinkedFrom).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, dto)
}

func putDetailed(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var req dtos.UpsertRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	var entity entities.Recipe
	if err := conn.Find(&entity, id).Error; err != nil || entity.ID == 0 {
		resp.Error(w, http.StatusNotFound, "recipe not found")
		return
	}
	if req.Status == "" {
		req.Status = entity.Status
	}

	fieldErrs, err := services.ValidateRecipe(conn, req, entity.ID)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(fieldErrs) > 0 {
		resp.JSON(w, http.StatusUnprocessableEntity, map[string]any{"error": "validation failed", "fields": fieldErrs})
		return
	}

	services.ApplyRecipeRequest(&entity, req)

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Save(&entity).Error; err != nil {
			return err
		}
		if err := services.ReplaceRecipeChildren(tx, entity.ID, req); err != nil {
			return err
		}
		return services.RecomputeRecipes(tx, entity.ID)
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, map[string]any{"id": entity.ID})
}

func deleteDetailed(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := r.URL.Query().Get("id")

	var entity entities.Recipe
	if err := conn.Find(&entity, id).Error; err != nil || entity.ID == 0 {
		resp.Error(w, http.StatusNotFound, "recipe not found")
		return
	}

	if err := conn.Select(clause.Associations).Delete(&entity).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, map[string]any{"deleted": entity.ID})
}
