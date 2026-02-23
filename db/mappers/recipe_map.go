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

func MapRecipeArrayToDto(a []entities.Recipe) []dtos.RecipeDto {
	out := make([]dtos.RecipeDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeToDto(i))
	}
	return out
}
