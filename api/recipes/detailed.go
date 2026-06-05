package api

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm/clause"
	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
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
		Preload("Steps").
		Preload("Ingredients").
		Preload("Tags").
		Find(&entity, id).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, mappers.MapRecipeToDto(entity))
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

	entity.Name = req.Name
	entity.Servings = req.Servings
	entity.Description = req.Description
	entity.Tagline = req.Tagline
	entity.Region = req.Region
	entity.TimeMinutes = req.TimeMinutes
	entity.Difficulty = req.Difficulty
	entity.KcalPerServing = req.KcalPerServing
	entity.NeedsPrep = req.NeedsPrep
	if req.Status != "" {
		entity.Status = req.Status
	}

	if err := conn.Save(&entity).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.PhotoUrl != nil {
		photo := entities.RecipePhoto{RecipeID: entity.ID, ImageURL: *req.PhotoUrl}
		conn.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "recipe_id"}},
			DoUpdates: clause.Assignments(map[string]any{"image_url": photo.ImageURL}),
		}).Create(&photo)
	}

	if err := conn.Where("recipe_id = ?", entity.ID).Delete(&entities.RecipeStep{}).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, s := range req.Steps {
		conn.Create(&entities.RecipeStep{RecipeID: entity.ID, StepNo: s.StepNo, Section: s.Section, Text: s.Text})
	}

	if err := conn.Where("recipe_id = ?", entity.ID).Delete(&entities.RecipeIngredient{}).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	for _, i := range req.Ingredients {
		conn.Create(&entities.RecipeIngredient{
			RecipeID: entity.ID, SortOrder: i.SortOrder, Section: i.Section,
			IngredientID: i.IngredientId, Amount: i.Amount, Unit: i.Unit,
			AmountText: i.AmountText, Note: i.Note,
		})
	}

	// replace tags
	conn.Where("recipe_id = ?", entity.ID).Delete(&entities.RecipeTag{})
	for _, tag := range req.Categories {
		conn.Create(&entities.RecipeTag{RecipeID: entity.ID, Tag: tag, GroupName: "category"})
	}
	for _, tag := range req.DietTags {
		conn.Create(&entities.RecipeTag{RecipeID: entity.ID, Tag: tag, GroupName: "diet"})
	}
	for _, tag := range req.PracticalTags {
		conn.Create(&entities.RecipeTag{RecipeID: entity.ID, Tag: tag, GroupName: "practical"})
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
