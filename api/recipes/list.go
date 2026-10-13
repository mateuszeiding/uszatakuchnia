package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"
)

func List(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()

	// admin=true query param means show all; otherwise only published
	isAdmin := r.URL.Query().Get("admin") == "true"

	query := conn.Model(&entities.Recipe{}).
		Select(`id, name, tagline, time_minutes, difficulty, status,
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

	dtos := mappers.MapRecipeBaseArrayToDto(list)
	for i := range dtos {
		dtos[i].IngredientNames = names[dtos[i].ID]
		if dtos[i].IngredientNames == nil {
			dtos[i].IngredientNames = []string{}
		}
	}
	resp.JSON(w, http.StatusOK, dtos)
}
