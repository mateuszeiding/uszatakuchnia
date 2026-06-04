package api

import (
	"encoding/json"
	"net/http"

	"gorm.io/gorm/clause"
	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
)

func Create(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}

	if r.Method != http.MethodPost {
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	var req dtos.UpsertRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Name == "" || req.Servings <= 0 {
		resp.Error(w, http.StatusBadRequest, "name and servings are required")
		return
	}

	status := req.Status
	if status == "" {
		status = "draft"
	}

	conn := db.DB()

	recipe := entities.Recipe{
		Name:           req.Name,
		Servings:       req.Servings,
		Description:    req.Description,
		Tagline:        req.Tagline,
		Category:       req.Category,
		Region:         req.Region,
		TimeMinutes:    req.TimeMinutes,
		Difficulty:     req.Difficulty,
		KcalPerServing: req.KcalPerServing,
		Status:         status,
		NeedsPrep:      req.NeedsPrep,
	}

	if err := conn.Create(&recipe).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	if req.PhotoUrl != nil && *req.PhotoUrl != "" {
		photo := entities.RecipePhoto{RecipeID: recipe.ID, ImageURL: *req.PhotoUrl}
		conn.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "recipe_id"}},
			DoUpdates: clause.Assignments(map[string]any{"image_url": photo.ImageURL}),
		}).Create(&photo)
	}

	for _, s := range req.Steps {
		conn.Create(&entities.RecipeStep{RecipeID: recipe.ID, StepNo: s.StepNo, Section: s.Section, Text: s.Text})
	}

	for _, i := range req.Ingredients {
		conn.Create(&entities.RecipeIngredient{
			RecipeID: recipe.ID, SortOrder: i.SortOrder, Section: i.Section,
			IngredientID: i.IngredientId, Amount: i.Amount, Unit: i.Unit,
			AmountText: i.AmountText, Note: i.Note,
		})
	}

	for _, tag := range req.DietTags {
		conn.Create(&entities.RecipeTag{RecipeID: recipe.ID, Tag: tag, GroupName: "diet"})
	}
	for _, tag := range req.PracticalTags {
		conn.Create(&entities.RecipeTag{RecipeID: recipe.ID, Tag: tag, GroupName: "practical"})
	}

	resp.JSON(w, http.StatusCreated, map[string]any{"id": recipe.ID})
}
