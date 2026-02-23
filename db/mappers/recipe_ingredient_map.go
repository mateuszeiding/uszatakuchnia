package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapRecipeIngredientToDto(e entities.RecipeIngredient) dtos.RecipeIngredientDto {
	return dtos.RecipeIngredientDto{
		SortOrder:      e.SortOrder,
		Section:        e.Section,
		IngredientID:   e.IngredientID,
		IngredientName: e.Ingredient.Name,
		Amount:         e.Amount,
		Unit:           e.Unit,
		AmountText:     e.AmountText,
		Note:           e.Note,
	}
}

func MapRecipeIngredientArrayToDto(a []entities.RecipeIngredient) []dtos.RecipeIngredientDto {
	out := make([]dtos.RecipeIngredientDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeIngredientToDto(i))
	}
	return out
}
