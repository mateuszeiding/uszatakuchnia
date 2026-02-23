package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapRecipeToDto(e entities.Recipe) dtos.RecipeDto {
	return dtos.RecipeDto{
		ID:          e.ID,
		Name:        e.Name,
		Servings:    e.Servings,
		Description: e.Description,
		Photo:       MapRecipePhotoToDto(e.Photo),
		Steps:       MapRecipeStepArrayToDto(e.Steps),
		Ingredients: MapRecipeIngredientArrayToDto(e.Ingredients),
	}
}

func MapRecipeBaseToDto(e entities.RecipeBase) dtos.RecipeBaseDto {
	return dtos.RecipeBaseDto{
		ID:    e.ID,
		Name:  e.Name,
		Photo: MapRecipePhotoToDto(e.Photo),
	}
}

func MapRecipeArrayToDto(a []entities.Recipe) []dtos.RecipeDto {
	out := make([]dtos.RecipeDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeToDto(i))
	}
	return out
}

func MapRecipeBaseArrayToDto(a []entities.RecipeBase) []dtos.RecipeBaseDto {
	out := make([]dtos.RecipeBaseDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeBaseToDto(i))
	}
	return out
}
