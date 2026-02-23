package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapRecipeStepToDto(e entities.RecipeStep) dtos.RecipeStepDto {
	return dtos.RecipeStepDto{
		StepNo:  e.StepNo,
		Section: e.Section,
		Text:    e.Text,
	}
}

func MapRecipeStepArrayToDto(a []entities.RecipeStep) []dtos.RecipeStepDto {
	out := make([]dtos.RecipeStepDto, 0, len(a))
	for _, i := range a {
		out = append(out, MapRecipeStepToDto(i))
	}
	return out
}
