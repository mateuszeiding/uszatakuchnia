package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	resp "uszatakuchnia/http"

	"gorm.io/gorm"
)

type valueRequest struct {
	AxisID uint    `json:"axisId"`
	Label  string  `json:"label"`
	Color  *string `json:"color"`
}

// Values — wartości osi klasyfikacji (panel → Klasyfikacja).
//
//	POST   /api/taxonomy/values                  {axisId, label, color?}
//	PUT    /api/taxonomy/values/:id              {label, color?}
//	DELETE /api/taxonomy/values/:id?moveTo=<id>  przepisy z tą wartością dostają moveTo
func Values(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	switch r.Method {
	case http.MethodPost:
		createValue(w, r)
	case http.MethodPut:
		updateValue(w, r)
	case http.MethodDelete:
		deleteValue(w, r)
	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func labelTaken(tx *gorm.DB, axisID uint, label string, exceptID uint) bool {
	var n int64
	tx.Model(&entities.TaxonomyValue{}).
		Where("axis_id = ? AND lower(label) = lower(?) AND id <> ?", axisID, label, exceptID).
		Count(&n)
	return n > 0
}

func createValue(w http.ResponseWriter, r *http.Request) {
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

func updateValue(w http.ResponseWriter, r *http.Request) {
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
	if err := conn.Preload("Axis").First(&val, r.URL.Query().Get("id")).Error; err != nil {
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

func deleteValue(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()

	var val entities.TaxonomyValue
	if err := conn.Preload("Axis").First(&val, r.URL.Query().Get("id")).Error; err != nil {
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
			// przepisy, które mają już wartość docelową, nie dostają jej drugi raz
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
