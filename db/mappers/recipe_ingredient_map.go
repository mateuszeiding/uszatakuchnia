package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapRecipeIngredientToDto(e entities.RecipeIngredient) dtos.RecipeIngredientDto {
	var unit *string
	if e.Unit != nil {
		unit = &e.Unit.Name
	}
	return dtos.RecipeIngredientDto{
		SortOrder:      e.SortOrder,
		Section:        e.Section,
		IngredientID:   e.IngredientID,
		IngredientName: e.Ingredient.Name,
		Quantity:       e.Quantity,
		UnitID:         e.UnitID,
		Unit:           unit,
		Note:           e.Note,
		Optional:       e.Optional,
		Allergens:      mapIngredientAllergenCodes(e.Ingredient.Allergens),
		IngredientType: dtos.IngredientTypeRefDto{Code: e.Ingredient.Type.Code, Label: e.Ingredient.Type.Name},
		SkipShopping:   e.Ingredient.SkipShopping,
		LinkedRecipe:   mapLinkedRecipe(e.LinkedRecipe),
		Substitutes:    mapSubstitutes(e.Substitutes),
	}
}

// mapSubstitutes wymaga preloadu Ingredients.Substitutes.Ingredient.{Allergens,Type}.
// gramsPerUnit, hasNutrition i nutrition uzupełnia services.FillLineNutrition.
func mapSubstitutes(a []entities.RecipeIngredientSubstitute) []dtos.RecipeSubstituteDto {
	out := make([]dtos.RecipeSubstituteDto, 0, len(a))
	for _, s := range a {
		out = append(out, dtos.RecipeSubstituteDto{
			IngredientID:   s.IngredientID,
			Name:           s.Ingredient.Name,
			Mix:            s.Mix,
			Allergens:      mapIngredientAllergenCodes(s.Ingredient.Allergens),
			IngredientType: dtos.IngredientTypeRefDto{Code: s.Ingredient.Type.Code, Label: s.Ingredient.Type.Name},
			SkipShopping:   s.Ingredient.SkipShopping,
		})
	}
	return out
}

// mapLinkedRecipe wymaga preloadu Ingredients.LinkedRecipe.Photo.
func mapLinkedRecipe(r *entities.Recipe) *dtos.LinkedRecipeDto {
	if r == nil || r.ID == 0 {
		return nil
	}
	var image *string
	if r.Photo != nil && r.Photo.ImageURL != "" {
		image = &r.Photo.ImageURL
	}
	return &dtos.LinkedRecipeDto{ID: r.ID, Title: r.Name, Status: r.Status, Image: image}
}

func MapRecipeIngredientArrayToDto(a []entities.RecipeIngredient) []dtos.RecipeIngredientDto {
	out := make([]dtos.RecipeIngredientDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeIngredientToDto(i))
	}
	return out
}
