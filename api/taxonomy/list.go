package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
)

// List — wszystkie słowniki: osie klasyfikacji z wartościami, jednostki, alergeny, typy składników.
func List(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	conn := db.DB()

	var axes []entities.TaxonomyAxis
	var values []entities.TaxonomyValue
	var units []entities.Unit
	var types []entities.IngredientType

	if err := conn.Order("sort_order, id").Find(&axes).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := conn.Order("sort_order, id").Find(&values).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := conn.Order("sort_order, id").Find(&units).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := conn.Order("id").Find(&types).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	// Liczby użyć (panel: usuwanie używanej wartości = przeniesienie przepisów).
	type count struct {
		ID uint
		N  int
	}
	var valueCounts, unitCounts []count
	if err := conn.Model(&entities.RecipeTaxonomy{}).Select("value_id AS id, COUNT(*) AS n").
		Group("value_id").Scan(&valueCounts).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := conn.Model(&entities.RecipeIngredient{}).Select("unit_id AS id, COUNT(*) AS n").
		Where("unit_id IS NOT NULL").Group("unit_id").Scan(&unitCounts).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	byValue := map[uint]int{}
	for _, c := range valueCounts {
		byValue[c.ID] = c.N
	}
	byUnit := map[uint]int{}
	for _, c := range unitCounts {
		byUnit[c.ID] = c.N
	}

	axesDto := mappers.MapTaxonomyAxesToDto(axes, values)
	for ai := range axesDto {
		for vi := range axesDto[ai].Values {
			n := byValue[axesDto[ai].Values[vi].ID]
			axesDto[ai].Values[vi].RecipeCount = &n
		}
	}
	unitsDto := mappers.MapUnitArrayToDto(units)
	for i := range unitsDto {
		unitsDto[i].UsedIn = byUnit[unitsDto[i].ID]
	}

	resp.JSON(w, http.StatusOK, dtos.DictionariesDto{
		Axes:            axesDto,
		Units:           unitsDto,
		Allergens:       mappers.MapAllergensToDto(),
		IngredientTypes: mappers.MapIngredientTypeArrayToDto(types),
	})
}
