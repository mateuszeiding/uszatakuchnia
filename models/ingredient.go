package models

import (
	enums "uszatakuchnia/enums/ingredient"
)

type Ingredient struct {
	Id         int
	ParentId   int
	Name       string
	Type       enums.IngredientType
	IsAllergen bool
}
