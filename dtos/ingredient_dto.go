package dtos

import "uszatakuchnia/db/entities"

type IngredientDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parentId"`
	Allergens []string  `json:"allergens"` // kody z entities.Allergens
	UsedIn    int       `json:"usedIn"`    // liczba przepisów (tylko lista)
	Type      TypeDto   `json:"type"`
	Image     *ImageDto `json:"image"`

	// Makro na 100 g: własne oraz wynikowe (własne albo odziedziczone z podskładników), po kluczu z entities.NutrientKeys.
	Nutrition         map[string]entities.MacroRange `json:"nutrition"`
	ResolvedNutrition map[string]entities.MacroRange `json:"resolvedNutrition"`
	// Przeliczniki jednostek: własne + odziedziczone od przodków (inherited) dla jednostek bez własnego.
	Conversions []ConversionDto `json:"conversions"`
}

type ConversionDto struct {
	UnitID    uint     `json:"unitId"`
	UnitName  string   `json:"unitName"`
	GramsMin  float64  `json:"gramsMin"`
	GramsMax  *float64 `json:"gramsMax"`
	Inherited bool     `json:"inherited"`
}

// IngredientRequest — tworzenie i edycja składnika (panel). Brak klucza w Nutrition = brak własnych danych.
type IngredientRequest struct {
	Name        string                         `json:"name"`
	TypeCode    string                         `json:"typeCode"`
	Allergens   []string                       `json:"allergens"`
	ParentID    *uint                          `json:"parentId"`
	Nutrition   map[string]entities.MacroRange `json:"nutrition"`
	Conversions []ConversionRequest            `json:"conversions"`
}

type ConversionRequest struct {
	UnitID   uint     `json:"unitId"`
	GramsMin float64  `json:"gramsMin"`
	GramsMax *float64 `json:"gramsMax"`
}
