package api

import (
	"net/http"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"
)

// List — wszystkie słowniki: osie klasyfikacji z wartościami, jednostki, alergeny, typy składników, sprzęt,
// formaty wyzwań.
func List(w http.ResponseWriter, r *http.Request) {
	if resp.Begin(w, r) {
		return
	}
	conn := db.DB()

	var axes []entities.TaxonomyAxis
	var values []entities.TaxonomyValue
	var units []entities.Unit
	var types []entities.IngredientType
	var equipment []entities.Equipment
	var formats []entities.ChallengeFormat

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
	if err := conn.Order("lower(name), id").Find(&equipment).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	eqUsage, err := services.EquipmentUsage(conn)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if err := conn.Order("sort_order, id").Find(&formats).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	formatUsage, err := services.ChallengeFormatUsage(conn)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	formatsDto := make([]dtos.ChallengeFormatDto, 0, len(formats))
	for _, f := range formats {
		d := mappers.MapChallengeFormatToDto(f)
		d.UsedIn = formatUsage[f.ID]
		formatsDto = append(formatsDto, d)
	}
	equipmentDto := make([]dtos.EquipmentDto, 0, len(equipment))
	for _, e := range equipment {
		d := mappers.MapEquipmentToDto(e)
		d.UsedIn = eqUsage[e.ID]
		equipmentDto = append(equipmentDto, d)
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
		Axes:             axesDto,
		Units:            unitsDto,
		Allergens:        mappers.MapAllergensToDto(),
		IngredientTypes:  mappers.MapIngredientTypeArrayToDto(types),
		Equipment:        equipmentDto,
		ChallengeFormats: formatsDto,
	})
}
