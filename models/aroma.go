package models

import (
	enums "uszatakuchnia/enums/aroma"
)

type Aroma struct {
	Id           int
	IngredientId int
	Name         string
	Type         enums.AromaType
	Intensity    int
}
