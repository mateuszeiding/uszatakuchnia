package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/handlers"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

type unitRequest struct {
	Name  string   `json:"name"`
	Forms []string `json:"forms"` // [1, 2–4, 5+] albo puste = nieodmienna
	NoQty bool     `json:"noQty"`
}

// Units — słownik jednostek (panel → Klasyfikacja → Jednostki).
//
//	POST   /api/taxonomy/units
//	PUT    /api/taxonomy/units/:id
//	DELETE /api/taxonomy/units/:id?moveTo=<id>  składniki w przepisach dostają moveTo (bez = bez jednostki)
//
// Ta sama funkcja obsługuje (limit 12 funkcji Vercela, vercel.json → ?kind=…):
//   - słownik sprzętu: /api/taxonomy/equipment[/:id] → ?kind=equipment,
//   - wyzwania: /api/challenges[/list|/:id] → ?kind=challenges (handlers.Challenges),
//   - formaty wyzwań: /api/taxonomy/challenge-formats[/:id] → ?kind=challenge-formats.
func Units(w http.ResponseWriter, r *http.Request) {
	if resp.Begin(w, r) {
		return
	}
	switch r.URL.Query().Get("kind") {
	case "equipment":
		equipmentHandler(w, r)
		return
	case "challenges":
		handlers.Challenges(w, r)
		return
	case "challenge-formats":
		handlers.ChallengeFormats(w, r)
		return
	}
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		saveUnit(w, r)
	case http.MethodDelete:
		deleteUnit(w, r)
	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

func saveUnit(w http.ResponseWriter, r *http.Request) {
	var req unitRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	name := strings.TrimSpace(req.Name)
	if name == "" {
		resp.Error(w, http.StatusBadRequest, "name is required")
		return
	}
	forms := make([]string, 0, 3)
	for _, f := range req.Forms {
		if f = strings.TrimSpace(f); f != "" {
			forms = append(forms, f)
		}
	}
	if len(forms) != 0 && len(forms) != 3 {
		resp.Error(w, http.StatusBadRequest, "forms must have 3 entries (1, 2–4, 5+) or be empty")
		return
	}

	conn := db.DB()
	var unit entities.Unit
	if r.Method == http.MethodPut {
		if err := conn.First(&unit, r.URL.Query().Get("id")).Error; err != nil {
			resp.Error(w, http.StatusNotFound, "unit not found")
			return
		}
	} else {
		conn.Model(&entities.Unit{}).Select("COALESCE(MAX(sort_order), 0) + 1").Scan(&unit.SortOrder)
	}

	var dup int64
	conn.Model(&entities.Unit{}).Where("lower(name) = lower(?) AND id <> ?", name, unit.ID).Count(&dup)
	if dup > 0 {
		resp.Error(w, http.StatusConflict, "unit already exists")
		return
	}

	unit.Name = name
	unit.NoQty = req.NoQty
	unit.FormOne, unit.FormFew, unit.FormMany = nil, nil, nil
	if len(forms) == 3 {
		unit.FormOne, unit.FormFew, unit.FormMany = &forms[0], &forms[1], &forms[2]
	}
	if err := conn.Save(&unit).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	status := http.StatusOK
	if r.Method == http.MethodPost {
		status = http.StatusCreated
	}
	resp.JSON(w, status, mappers.MapUnitToDto(unit))
}

func deleteUnit(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	var unit entities.Unit
	if err := conn.First(&unit, r.URL.Query().Get("id")).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "unit not found")
		return
	}

	var moveTo *uint
	if s := r.URL.Query().Get("moveTo"); s != "" {
		n, err := strconv.ParseUint(s, 10, 64)
		if err != nil || uint(n) == unit.ID {
			resp.Error(w, http.StatusBadRequest, "invalid moveTo")
			return
		}
		var target entities.Unit
		if err := conn.First(&target, n).Error; err != nil {
			resp.Error(w, http.StatusBadRequest, "unknown moveTo unit")
			return
		}
		moveTo = &target.ID
	}

	var used int64
	conn.Model(&entities.RecipeIngredient{}).Where("unit_id = ?", unit.ID).Count(&used)

	if err := conn.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&entities.RecipeIngredient{}).Where("unit_id = ?", unit.ID).
			Update("unit_id", moveTo).Error; err != nil {
			return err
		}
		return tx.Delete(&entities.Unit{}, unit.ID).Error
	}); err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.JSON(w, http.StatusOK, map[string]any{"deleted": unit.ID, "moved": used})
}

// ─── Sprzęt (Panel → Klasyfikacja → Sprzęt) ─────────────────
//
//	POST   /api/taxonomy/equipment                      {name}
//	PUT    /api/taxonomy/equipment/:id                  {name}
//	DELETE /api/taxonomy/equipment/:id?replaceWith=<id>  przepisy dostają inny sprzęt (doprecyzowania zostają)
//	DELETE /api/taxonomy/equipment/:id?removeFromRecipes=true
func equipmentHandler(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	q := r.URL.Query()
	id, _ := strconv.ParseUint(q.Get("id"), 10, 64)

	switch r.Method {
	case http.MethodPost, http.MethodPut:
		if r.Method == http.MethodPut && id == 0 {
			resp.Error(w, http.StatusBadRequest, "id is required")
			return
		}
		if r.Method == http.MethodPost {
			id = 0
		}
		var req struct {
			Name string `json:"name"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			resp.Error(w, http.StatusBadRequest, "invalid body")
			return
		}
		e, err := services.SaveEquipment(conn, uint(id), req.Name)
		switch {
		case errors.Is(err, services.ErrEquipmentNotFound):
			resp.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrEquipmentTaken):
			resp.Error(w, http.StatusConflict, err.Error())
		case err != nil:
			resp.Error(w, http.StatusBadRequest, err.Error())
		default:
			status := http.StatusOK
			if r.Method == http.MethodPost {
				status = http.StatusCreated
			}
			resp.JSON(w, status, mappers.MapEquipmentToDto(e))
		}
	case http.MethodDelete:
		replaceWith, _ := strconv.ParseUint(q.Get("replaceWith"), 10, 64)
		used, err := services.DeleteEquipment(conn, uint(id), uint(replaceWith), q.Get("removeFromRecipes") == "true")
		switch {
		case errors.Is(err, services.ErrEquipmentNotFound):
			resp.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrEquipmentInUse):
			resp.JSON(w, http.StatusConflict, map[string]any{"error": err.Error(), "usedIn": used})
		case errors.Is(err, services.ErrEquipmentReplace):
			resp.Error(w, http.StatusBadRequest, err.Error())
		case err != nil:
			resp.Error(w, http.StatusInternalServerError, err.Error())
		default:
			resp.JSON(w, http.StatusOK, map[string]any{"deleted": id, "recipes": used})
		}
	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
