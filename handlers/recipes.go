package handlers

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

// RecipeList — GET /api/recipes/list[?admin=true]. admin=true — także szkice, tylko dla właściciela.
func RecipeList(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()

	// admin=true — także szkice, tylko dla właściciela; reszta widzi opublikowane
	isAdmin := r.URL.Query().Get("admin") == "true"
	if isAdmin && !resp.IsOwner(r) {
		resp.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	query := conn.Model(&entities.Recipe{}).
		Select(`id, name, tagline, time_min, time_max, difficulty, status, nutrition,
			EXISTS (SELECT 1 FROM recipe_steps s WHERE s.recipe_id = recipes.id AND s.prep) AS needs_prep`).
		Preload("Photo").
		Preload("Taxonomy.Value.Axis").
		Order("id DESC")

	if !isAdmin {
		query = query.Where("status = ?", entities.StatusPublished)
	}

	var list []entities.RecipeBase
	if err := query.Find(&list).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}

	// Nazwy składników per przepis (do wyszukiwania po składnikach).
	var rows []struct {
		RecipeID uint
		Name     string
	}
	if err := conn.Table("recipe_ingredients ri").
		Select("ri.recipe_id, i.name").
		Joins("JOIN ingredients i ON i.id = ri.ingredient_id").
		Order("ri.recipe_id, ri.sort_order").
		Scan(&rows).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	names := map[uint][]string{}
	for _, r := range rows {
		names[r.RecipeID] = append(names[r.RecipeID], r.Name)
	}

	// Nazwy zamienników per przepis (wyszukiwanie: „pancetta” znajduje przepis z boczkiem).
	var subRows []struct {
		RecipeID uint
		Name     string
	}
	if err := conn.Table("recipe_ingredient_substitutes s").
		Select("s.recipe_id, i.name, MIN(s.sort_order) AS so").
		Joins("JOIN ingredients i ON i.id = s.ingredient_id").
		Group("s.recipe_id, i.name").
		Order("s.recipe_id, so, i.name").
		Scan(&subRows).Error; err != nil {
		resp.JSON(w, http.StatusInternalServerError, map[string]any{"error": err.Error()})
		return
	}
	subs := map[uint][]string{}
	for _, r := range subRows {
		subs[r.RecipeID] = append(subs[r.RecipeID], r.Name)
	}

	dtos := mappers.MapRecipeBaseArrayToDto(list)
	for i := range dtos {
		dtos[i].IngredientNames = names[dtos[i].ID]
		if dtos[i].IngredientNames == nil {
			dtos[i].IngredientNames = []string{}
		}
		dtos[i].SubstituteNames = subs[dtos[i].ID]
		if dtos[i].SubstituteNames == nil {
			dtos[i].SubstituteNames = []string{}
		}
	}
	resp.JSON(w, http.StatusOK, dtos)
}

// RecipeCreate — POST /api/recipes. Bez statusu = szkic.
func RecipeCreate(w http.ResponseWriter, r *http.Request) {
	var req dtos.UpsertRecipeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.Status == "" {
		req.Status = entities.StatusDraft
	}

	conn := db.DB()

	fieldErrs, err := services.ValidateRecipe(conn, req, 0)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(fieldErrs) > 0 {
		resp.JSON(w, http.StatusUnprocessableEntity, map[string]any{"error": "validation failed", "fields": fieldErrs})
		return
	}

	var recipe entities.Recipe
	services.ApplyRecipeRequest(&recipe, req)

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&recipe).Error; err != nil {
			return err
		}
		if err := services.ReplaceRecipeChildren(tx, recipe.ID, req); err != nil {
			return err
		}
		return services.RecomputeRecipes(tx, recipe.ID)
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusCreated, map[string]any{"id": recipe.ID})
}

// RecipeGet — GET /api/recipes/{id}. Szkic widzi tylko właściciel; reszta dostaje 404.
func RecipeGet(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := pathID(r)

	var entity entities.Recipe
	if err := conn.
		Preload("Photo").
		Preload("Steps", func(tx *gorm.DB) *gorm.DB { return tx.Order("step_no") }).
		Preload("Ingredients", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order") }).
		Preload("Ingredients.Ingredient.Allergens").
		Preload("Ingredients.Ingredient.Type").
		Preload("Ingredients.Unit").
		Preload("Ingredients.LinkedRecipe.Photo").
		Preload("Ingredients.Substitutes", func(tx *gorm.DB) *gorm.DB { return tx.Order("position") }).
		Preload("Ingredients.Substitutes.Ingredient.Allergens").
		Preload("Ingredients.Substitutes.Ingredient.Type").
		Preload("Taxonomy.Value.Axis").
		Preload("Equipment", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order") }).
		Preload("Equipment.Equipment").
		Find(&entity, id).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// Szkic widzi tylko właściciel; reszta dostaje 404, jakby przepisu nie było.
	owner := resp.IsOwner(r)
	if entity.ID == 0 || (!owner && entity.Status != entities.StatusPublished) {
		resp.Error(w, http.StatusNotFound, "recipe not found")
		return
	}

	dto := mappers.MapRecipeToDto(entity)
	if !owner {
		// link do przepisu-składnika, który jest szkicem, zdradzałby szkic
		for i := range dto.Ingredients {
			if l := dto.Ingredients[i].LinkedRecipe; l != nil && l.Status != entities.StatusPublished {
				dto.Ingredients[i].LinkedRecipe = nil
			}
		}
	}
	ix, err := services.LoadNutritionIndex(conn)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	services.FillGramsPerUnit(dto.Ingredients, ix)
	services.FillLineNutrition(dto.Ingredients, entity, ix)

	// Kto prowadzi do tego przepisu (modal usuwania): „Przepis · składnik”.
	dto.LinkedFrom = []dtos.RecipeLinkSourceDto{}
	if err := conn.Table("recipe_ingredients ri").
		Select("r.id AS recipe_id, r.name AS recipe_name, i.name AS ingredient_name").
		Joins("JOIN recipes r ON r.id = ri.recipe_id").
		Joins("JOIN ingredients i ON i.id = ri.ingredient_id").
		Where("ri.linked_recipe_id = ? AND ri.recipe_id <> ?", entity.ID, entity.ID).
		Where("? OR r.status = ?", owner, entities.StatusPublished).
		Order("r.name, ri.sort_order").
		Scan(&dto.LinkedFrom).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, dto)
}

// RecipeUpdate — PUT /api/recipes/{id}. Bez statusu = status bez zmian.
func RecipeUpdate(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := pathID(r)

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

// RecipeDelete — DELETE /api/recipes/{id}.
func RecipeDelete(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := pathID(r)

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
