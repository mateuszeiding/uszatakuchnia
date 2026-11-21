package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"
)

// ─── Glossary (admin → Słowniczek; terms marked in steps and hints) ─────────────────
//
//	GET    /api/glossary                          public; the owner also gets usedIn
//	POST   /api/glossary                          {name, explanation}
//	PUT    /api/glossary/{id}                     {name, explanation}
//	GET    /api/glossary/{id}/usage               owner: recipes that mark the term
//	DELETE /api/glossary/{id}?replaceWith=<id>    marks point at another term
//	DELETE /api/glossary/{id}?removeMarks=true    the words stay as plain text

func GlossaryList(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	var terms []entities.GlossaryTerm
	if err := conn.Order("lower(name), id").Find(&terms).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	// usage counts drafts, so only the owner sees it
	var usage map[uint]map[uint]bool
	if resp.IsOwner(r) {
		var err error
		if usage, err = services.TermUsage(conn); err != nil {
			resp.Error(w, http.StatusInternalServerError, err.Error())
			return
		}
	}
	out := make([]dtos.GlossaryTermDto, 0, len(terms))
	for _, t := range terms {
		d := dtos.GlossaryTermDto{ID: t.ID, Name: t.Name, Explanation: t.Explanation}
		if usage != nil {
			n := len(usage[t.ID])
			d.UsedIn = &n
		}
		out = append(out, d)
	}
	resp.JSON(w, http.StatusOK, out)
}

func GlossaryCreate(w http.ResponseWriter, r *http.Request) { saveGlossaryTerm(w, r, 0) }

func GlossaryUpdate(w http.ResponseWriter, r *http.Request) { saveGlossaryTerm(w, r, pathID(r)) }

func saveGlossaryTerm(w http.ResponseWriter, r *http.Request, id uint) {
	var req dtos.UpsertGlossaryTermRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	t, fields, err := services.SaveGlossaryTerm(db.DB(), id, req)
	switch {
	case errors.Is(err, services.ErrTermNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrTermInvalid):
		resp.JSON(w, http.StatusUnprocessableEntity, map[string]any{"error": err.Error(), "fields": fields})
	case errors.Is(err, services.ErrTermTaken):
		resp.JSON(w, http.StatusConflict, map[string]any{"error": err.Error(), "fields": fields})
	case err != nil:
		resp.Error(w, http.StatusInternalServerError, err.Error())
	default:
		status := http.StatusOK
		if id == 0 {
			status = http.StatusCreated
		}
		resp.JSON(w, status, dtos.GlossaryTermDto{ID: t.ID, Name: t.Name, Explanation: t.Explanation})
	}
}

// GlossaryUsage lists recipe names, drafts included, so it's for the owner only.
func GlossaryUsage(w http.ResponseWriter, r *http.Request) {
	if !resp.IsOwner(r) {
		resp.Error(w, http.StatusUnauthorized, "unauthorized")
		return
	}
	conn := db.DB()
	var t entities.GlossaryTerm
	if err := conn.First(&t, pathID(r)).Error; err != nil {
		resp.Error(w, http.StatusNotFound, services.ErrTermNotFound.Error())
		return
	}
	list, err := services.TermRecipes(conn, t.ID)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	resp.JSON(w, http.StatusOK, list)
}

func GlossaryDelete(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	id := pathID(r)
	replaceWith, _ := strconv.ParseUint(q.Get("replaceWith"), 10, 32)
	used, err := services.DeleteGlossaryTerm(db.DB(), id, uint(replaceWith), q.Get("removeMarks") == "true")
	switch {
	case errors.Is(err, services.ErrTermNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrTermInUse):
		resp.JSON(w, http.StatusConflict, map[string]any{"error": err.Error(), "usedIn": used})
	case errors.Is(err, services.ErrTermReplace):
		resp.Error(w, http.StatusBadRequest, err.Error())
	case err != nil:
		resp.Error(w, http.StatusInternalServerError, err.Error())
	default:
		resp.JSON(w, http.StatusOK, map[string]any{"deleted": id, "recipes": used})
	}
}
