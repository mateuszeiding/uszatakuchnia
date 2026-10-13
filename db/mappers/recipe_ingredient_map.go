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
	}
}

func MapRecipeIngredientArrayToDto(a []entities.RecipeIngredient) []dtos.RecipeIngredientDto {
	out := make([]dtos.RecipeIngredientDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeIngredientToDto(i))
	}
	return out
}
