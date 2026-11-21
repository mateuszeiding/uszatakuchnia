package dtos

import "uszatakuchnia/db/entities"

type RecipeIngredientDto struct {
	SortOrder int `json:"sortOrder"`
	// Stable row key; steps refer to the row with {{ing:<key>}}.
	Key     string  `json:"key"`
	Section *string `json:"section,omitempty"`

	IngredientID   uint   `json:"ingredientId"`
	IngredientName string `json:"ingredientName"`

	Quantity *string `json:"quantity,omitempty"`
	UnitID   *uint   `json:"unitId,omitempty"`
	Unit     *string `json:"unit,omitempty"` // unit name (dictionary form); the frontend inflects it
	Note     *string `json:"note,omitempty"`
	Optional bool    `json:"optional"`

	Allergens []string `json:"allergens"`

	// Store aisle (shopping list).
	IngredientType IngredientTypeRefDto `json:"ingredientType"`
	// Skipped on the shopping list (e.g. water).
	SkipShopping bool `json:"skipShopping"`
	// Grams per 1 unit of this row; nil = can't convert (the shopping list adds up mixed units via grams).
	GramsPerUnit *GramsRangeDto `json:"gramsPerUnit"`

	// Drafts are returned too (with their status): the public page hides them, the form shows them.
	LinkedRecipe *LinkedRecipeDto `json:"linkedRecipe"`

	// This row's share of the per-serving macros (shares add up to RecipeDto.nutrition); nil = the row doesn't count
	// (optional, to taste, no quantity). The frontend computes the chosen variant's macros from these shares.
	Nutrition *entities.RecipeNutrition `json:"nutrition"`
	// Substitutes in order; [] = none.
	Substitutes []RecipeSubstituteDto `json:"substitutes"`
}

// RecipeSubstituteDto is a row's substitute with data for allergens, the shopping list and macros (no extra request).
type RecipeSubstituteDto struct {
	IngredientID   uint                 `json:"ingredientId"`
	Name           string               `json:"name"`
	Mix            bool                 `json:"mix"` // "half and half" with the main ingredient
	Allergens      []string             `json:"allergens"`
	IngredientType IngredientTypeRefDto `json:"ingredientType"`
	SkipShopping   bool                 `json:"skipShopping"`
	// Grams per 1 unit of the row for the substitute; nil = can't convert.
	GramsPerUnit *GramsRangeDto `json:"gramsPerUnit"`
	// The ingredient has any macro data (otherwise "no nutrition data").
	HasNutrition bool `json:"hasNutrition"`
	// The variant's share of the per-serving macros (mix = half main + half substitute); nil = the row doesn't count.
	Nutrition *entities.RecipeNutrition `json:"nutrition"`
}

type LinkedRecipeDto struct {
	ID     uint    `json:"id"`
	Title  string  `json:"title"`
	Status string  `json:"status"`
	Image  *string `json:"image"`
}

type IngredientTypeRefDto struct {
	Code  string `json:"code"`
	Label string `json:"label"`
}

type GramsRangeDto struct {
	Min float64 `json:"min"`
	Max float64 `json:"max"`
}
