package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func extractTags(tags []entities.RecipeTag) (categories []string, diet []string, practical []string) {
	categories = []string{}
	diet = []string{}
	practical = []string{}
	for _, t := range tags {
		switch t.GroupName {
		case "category":
			categories = append(categories, t.Tag)
		case "diet":
			diet = append(diet, t.Tag)
		case "practical":
			practical = append(practical, t.Tag)
		}
	}
	return
}

func MapRecipeToDto(e entities.Recipe) dtos.RecipeDto {
	cats, diet, practical := extractTags(e.Tags)
	return dtos.RecipeDto{
		ID:             e.ID,
		Name:           e.Name,
		Servings:       e.Servings,
		Description:    e.Description,
		Tagline:        e.Tagline,
		Region:         e.Region,
		TimeMinutes:    e.TimeMinutes,
		Difficulty:     e.Difficulty,
		KcalPerServing: e.KcalPerServing,
		Status:         e.Status,
		NeedsPrep:      e.NeedsPrep,
		Categories:     cats,
		DietTags:       diet,
		PracticalTags:  practical,
		Photo:          MapRecipePhotoToDto(e.Photo),
		Steps:          MapRecipeStepArrayToDto(e.Steps),
		Ingredients:    MapRecipeIngredientArrayToDto(e.Ingredients),
	}
}

func MapRecipeBaseToDto(e entities.RecipeBase) dtos.RecipeBaseDto {
	cats, diet, practical := extractTags(e.Tags)
	return dtos.RecipeBaseDto{
		ID:            e.ID,
		Name:          e.Name,
		Tagline:       e.Tagline,
		Region:        e.Region,
		TimeMinutes:   e.TimeMinutes,
		Difficulty:    e.Difficulty,
		Status:        e.Status,
		NeedsPrep:     e.NeedsPrep,
		Categories:    cats,
		DietTags:      diet,
		PracticalTags: practical,
		Photo:         MapRecipePhotoToDto(e.Photo),
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
