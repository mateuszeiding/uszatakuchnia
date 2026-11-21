package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// recipeAllergens: union of the ingredients' allergens, in EU list order.
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
		if s.Prep && !s.Optional {
			return true
		}
	}
	return false
}

// MapRecipeToDto needs these preloads: Photo, Steps, Ingredients.Ingredient.Allergens, Ingredients.Ingredient.Type, Ingredients.Unit, Taxonomy.Value.Axis.
// Ingredient gramsPerUnit is filled in by services.FillGramsPerUnit.
func MapRecipeToDto(e entities.Recipe) dtos.RecipeDto {
	return dtos.RecipeDto{
		ID:          e.ID,
		Name:        e.Name,
		Servings:    e.Servings,
		Description: e.Description,
		Tagline:     e.Tagline,
		Tip:         e.Tip,
		Cautions:    append([]string{}, e.Cautions...),
		TimeMin:     e.TimeMin,
		TimeMax:     e.TimeMax,
		Difficulty:  e.Difficulty,
		Status:      e.Status,
		NeedsPrep:   recipeNeedsPrep(e.Steps),
		Taxonomy:    MapRecipeTaxonomyToDto(e.Taxonomy),
		Allergens:   recipeAllergens(e.Ingredients),
		Nutrition:   e.Nutrition,
		Photo:       MapRecipePhotoToDto(e.Photo),
		Steps:       MapRecipeStepArrayToDto(e.Steps),
		Ingredients: MapRecipeIngredientArrayToDto(e.Ingredients),
		Equipment:   MapRecipeEquipmentToDto(e.Equipment),
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
		Kcal:       recipeKcal(e.Nutrition),
	}
}

func recipeKcal(n *entities.RecipeNutrition) *entities.SumRange {
	if n == nil {
		return nil
	}
	if k, ok := n.PerPortion["kcal"]; ok {
		return &k
	}
	return nil
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
