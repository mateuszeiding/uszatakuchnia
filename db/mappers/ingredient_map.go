package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapIngredientToDto(e entities.Ingredient) dtos.IngredientDto {
	ingredientType := dtos.TypeDto{
		Name: e.Type.Name,
		Code: e.Type.Code,
	}

	return dtos.IngredientDto{
		ID:         e.ID,
		Name:       e.Name,
		ParentID:   e.ParentID,
		IsAllergen: e.IsAllergen,
		Type:       ingredientType,
		Image:      MapImageToDto(&e.Image),
	}
}

func MapIngredientArrayToDto(a []entities.Ingredient) []dtos.IngredientDto {
	out := make([]dtos.IngredientDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapIngredientToDto(i))
	}
	return out
}
