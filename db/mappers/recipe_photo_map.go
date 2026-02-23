package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapRecipePhotoToDto(e *entities.RecipePhoto) *dtos.RecipePhotoDto {
	if e == nil {
		return nil
	}

	return &dtos.RecipePhotoDto{
		URL: e.ImageURL,
	}
}
