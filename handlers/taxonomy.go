package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

func Dictionaries(w http.ResponseWriter, r *http.Request) {
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

	// Usage counts (admin panel: deleting a used value = moving its recipes).
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

// ─── Classification axis values (admin → Classification) ─────────────────
//
//	POST   /api/taxonomy/values                  {axisId, label, color?}
//	PUT    /api/taxonomy/values/{id}             {label, color?}
//	DELETE /api/taxonomy/values/{id}?moveTo=<id>  recipes with this value get moveTo

type valueRequest struct {
	AxisID uint    `json:"axisId"`
	Label  string  `json:"label"`
	Color  *string `json:"color"`
}

func labelTaken(tx *gorm.DB, axisID uint, label string, exceptID uint) bool {
	var n int64
	tx.Model(&entities.TaxonomyValue{}).
		Where("axis_id = ? AND lower(label) = lower(?) AND id <> ?", axisID, label, exceptID).
		Count(&n)
	return n > 0
}

func ValueCreate(w http.ResponseWriter, r *http.Request) {
	var req valueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	label := strings.TrimSpace(req.Label)
	if label == "" || req.AxisID == 0 {
		resp.Error(w, http.StatusBadRequest, "axisId and label are required")
		return
	}
	conn := db.DB()

	var axis entities.TaxonomyAxis
	if err := conn.First(&axis, req.AxisID).Error; err != nil {
		resp.Error(w, http.StatusBadRequest, "unknown axis")
		return
	}
	if labelTaken(conn, axis.ID, label, 0) {
		resp.Error(w, http.StatusConflict, "value already exists")
		return
	}

	var maxSort int
	conn.Model(&entities.TaxonomyValue{}).Where("axis_id = ?", axis.ID).
		Select("COALESCE(MAX(sort_order), 0)").Scan(&maxSort)

	val := entities.TaxonomyValue{AxisID: axis.ID, Label: label, SortOrder: maxSort + 1}
	if axis.HasColor {
		val.Color = req.Color
	}
	if err := conn.Omit("Axis").Create(&val).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.JSON(w, http.StatusCreated, mappers.MapTaxonomyValueToDto(val))
}

func ValueUpdate(w http.ResponseWriter, r *http.Request) {
	var req valueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	label := strings.TrimSpace(req.Label)
	if label == "" {
		resp.Error(w, http.StatusBadRequest, "label is required")
		return
	}
	conn := db.DB()

	var val entities.TaxonomyValue
	if err := conn.Preload("Axis").First(&val, pathID(r)).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "value not found")
		return
	}
	if labelTaken(conn, val.AxisID, label, val.ID) {
		resp.Error(w, http.StatusConflict, "value already exists")
		return
	}

	updates := map[string]any{"label": label}
	if val.Axis.HasColor {
		updates["color"] = req.Color
	}
	if err := conn.Model(&entities.TaxonomyValue{}).Where("id = ?", val.ID).Updates(updates).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	val.Label = label
	if val.Axis.HasColor {
		val.Color = req.Color
	}
	resp.JSON(w, http.StatusOK, mappers.MapTaxonomyValueToDto(val))
}

func ValueDelete(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()

	var val entities.TaxonomyValue
	if err := conn.Preload("Axis").First(&val, pathID(r)).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "value not found")
		return
	}

	var moveTo uint
	if s := r.URL.Query().Get("moveTo"); s != "" {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			resp.Error(w, http.StatusBadRequest, "invalid moveTo")
			return
		}
		var target entities.TaxonomyValue
		if err := conn.First(&target, n).Error; err != nil || target.AxisID != val.AxisID || target.ID == val.ID {
			resp.Error(w, http.StatusBadRequest, "moveTo must be another value of the same axis")
			return
		}
		moveTo = target.ID
	}

	var used int64
	conn.Model(&entities.RecipeTaxonomy{}).Where("value_id = ?", val.ID).Count(&used)
	if used > 0 && val.Axis.Required && moveTo == 0 {
		resp.JSON(w, http.StatusConflict, map[string]any{
			"error":  fmt.Sprintf("value used in %d recipes; axis is required, pass moveTo", used),
			"usedIn": used,
		})
		return
	}

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if moveTo != 0 {
			// recipes that already have the target value don't get it twice
			if err := tx.Exec(
				"INSERT INTO recipe_taxonomy (recipe_id, value_id) "+
					"SELECT recipe_id, ? FROM recipe_taxonomy WHERE value_id = ? ON CONFLICT DO NOTHING",
				moveTo, val.ID,
			).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("value_id = ?", val.ID).Delete(&entities.RecipeTaxonomy{}).Error; err != nil {
			return err
		}
		return tx.Delete(&entities.TaxonomyValue{}, val.ID).Error
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusOK, map[string]any{"deleted": val.ID, "moved": used})
}
