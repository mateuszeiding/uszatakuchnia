package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// recipeAllergens — suma alergenów składników, w kolejności listy UE.
func recipeAllergens(items []entities.RecipeIngredient) []string {
	present := map[string]bool{}
	for _, it := range items {
		for _, a := range it.Ingredient.Allergens {
			present[a.Code] = true
		}
	}
	out := []string{}
	for _, a := range entities.Allergens {
		if present[a.Code] {
			out = append(out, a.Code)
		}
	}
	return out
}

func recipeNeedsPrep(steps []entities.RecipeStep) bool {
	for _, s := range steps {
		if s.Prep {
			return true
		}
	}
	return false
}

// MapRecipeToDto wymaga preloadów: Photo, Steps, Ingredients.Ingredient.Allergens, Ingredients.Unit, Taxonomy.Value.Axis.
func MapRecipeToDto(e entities.Recipe) dtos.RecipeDto {
	return dtos.RecipeDto{
		ID:             e.ID,
		Name:           e.Name,
		Servings:       e.Servings,
		Description:    e.Description,
		Tagline:        e.Tagline,
		Tip:            e.Tip,
		TimeMin:        e.TimeMin,
		TimeMax:        e.TimeMax,
		Difficulty:     e.Difficulty,
		KcalPerServing: e.KcalPerServing,
		Status:         e.Status,
		NeedsPrep:      recipeNeedsPrep(e.Steps),
		Taxonomy:       MapRecipeTaxonomyToDto(e.Taxonomy),
		Allergens:      recipeAllergens(e.Ingredients),
		Photo:          MapRecipePhotoToDto(e.Photo),
		Steps:          MapRecipeStepArrayToDto(e.Steps),
		Ingredients:    MapRecipeIngredientArrayToDto(e.Ingredients),
	}
}

func MapRecipeBaseToDto(e entities.RecipeBase) dtos.RecipeBaseDto {
	return dtos.RecipeBaseDto{
		ID:         e.ID,
		Name:       e.Name,
		Tagline:    e.Tagline,
		TimeMin:    e.TimeMin,
		TimeMax:    e.TimeMax,
		Difficulty: e.Difficulty,
		Status:     e.Status,
		NeedsPrep:  e.NeedsPrep,
		Taxonomy:   MapRecipeTaxonomyToDto(e.Taxonomy),
		Photo:      MapRecipePhotoToDto(e.Photo),
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
