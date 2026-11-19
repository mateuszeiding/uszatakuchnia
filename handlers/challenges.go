package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"uszatakuchnia/db"
	"uszatakuchnia/db/mappers"
	"uszatakuchnia/dtos"
	resp "uszatakuchnia/http"
	"uszatakuchnia/services"
)

// ─── Wyzwania ─────────────────
//
//	GET    /api/challenges/list   lista, najnowsze na górze (też GET /api/challenges)
//	GET    /api/challenges/{id}   szczegóły (404, gdy nie ma)
//	POST   /api/challenges        nowe wyzwanie (od razu publiczne)
//	PUT    /api/challenges/{id}   zapis
//	DELETE /api/challenges/{id}

// ChallengeList — GET /api/challenges/list.
func ChallengeList(w http.ResponseWriter, r *http.Request) {
	list, err := services.LoadChallenges(db.DB())
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	out := make([]dtos.ChallengeDto, 0, len(list))
	for _, c := range list {
		out = append(out, mappers.MapChallengeToDto(c))
	}
	resp.JSON(w, http.StatusOK, out)
}

// ChallengeGet — GET /api/challenges/{id}.
func ChallengeGet(w http.ResponseWriter, r *http.Request) {
	c, err := services.LoadChallenge(db.DB(), pathID(r))
	switch {
	case errors.Is(err, services.ErrChallengeNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case err != nil:
		resp.Error(w, http.StatusInternalServerError, err.Error())
	default:
		resp.JSON(w, http.StatusOK, mappers.MapChallengeToDto(c))
	}
}

// ChallengeCreate — POST /api/challenges.
func ChallengeCreate(w http.ResponseWriter, r *http.Request) { saveChallenge(w, r, 0) }

// ChallengeUpdate — PUT /api/challenges/{id}.
func ChallengeUpdate(w http.ResponseWriter, r *http.Request) { saveChallenge(w, r, pathID(r)) }

func saveChallenge(w http.ResponseWriter, r *http.Request, id uint) {
	conn := db.DB()
	var req dtos.UpsertChallengeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}
	req = services.NormalizeChallengeRequest(req)
	fieldErrs, err := services.ValidateChallenge(conn, &req)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	if len(fieldErrs) > 0 {
		resp.JSON(w, http.StatusUnprocessableEntity, map[string]any{"error": "validation failed", "fields": fieldErrs})
		return
	}
	savedID, err := services.SaveChallenge(conn, id, req)
	if errors.Is(err, services.ErrChallengeNotFound) {
		resp.Error(w, http.StatusNotFound, err.Error())
		return
	}
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	c, err := services.LoadChallenge(conn, savedID)
	if err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}
	status := http.StatusOK
	if id == 0 {
		status = http.StatusCreated
	}
	resp.JSON(w, status, mappers.MapChallengeToDto(c))
}

// ChallengeDelete — DELETE /api/challenges/{id}.
func ChallengeDelete(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	err := services.DeleteChallenge(db.DB(), id)
	switch {
	case errors.Is(err, services.ErrChallengeNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case err != nil:
		resp.Error(w, http.StatusInternalServerError, err.Error())
	default:
		resp.JSON(w, http.StatusOK, map[string]any{"deleted": id})
	}
}

// ─── Formaty wyzwań (panel → Klasyfikacja; lista z liczbą użyć jest w GET /api/taxonomy/list) ───
//
//	POST   /api/taxonomy/challenge-formats                  {name}
//	PUT    /api/taxonomy/challenge-formats/{id}             {name}
//	DELETE /api/taxonomy/challenge-formats/{id}?moveTo=<id>  używany: wyzwania przechodzą do moveTo

// ChallengeFormatCreate — POST /api/taxonomy/challenge-formats.
func ChallengeFormatCreate(w http.ResponseWriter, r *http.Request) { saveChallengeFormat(w, r, 0) }

// ChallengeFormatUpdate — PUT /api/taxonomy/challenge-formats/{id}.
func ChallengeFormatUpdate(w http.ResponseWriter, r *http.Request) {
	saveChallengeFormat(w, r, pathID(r))
}

func saveChallengeFormat(w http.ResponseWriter, r *http.Request, id uint) {
	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid body")
		return
	}
	f, err := services.SaveChallengeFormat(db.DB(), id, req.Name)
	switch {
	case errors.Is(err, services.ErrChallengeFormatNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrChallengeFormatTaken):
		resp.Error(w, http.StatusConflict, err.Error())
	case err != nil:
		resp.Error(w, http.StatusBadRequest, err.Error())
	default:
		status := http.StatusOK
		if id == 0 {
			status = http.StatusCreated
		}
		resp.JSON(w, status, mappers.MapChallengeFormatToDto(f))
	}
}

// ChallengeFormatDelete — DELETE /api/taxonomy/challenge-formats/{id}?moveTo=<id>.
func ChallengeFormatDelete(w http.ResponseWriter, r *http.Request) {
	id := pathID(r)
	moveTo, _ := strconv.ParseUint(r.URL.Query().Get("moveTo"), 10, 64)
	used, err := services.DeleteChallengeFormat(db.DB(), id, uint(moveTo))
	switch {
	case errors.Is(err, services.ErrChallengeFormatNotFound):
		resp.Error(w, http.StatusNotFound, err.Error())
	case errors.Is(err, services.ErrChallengeFormatInUse):
		resp.JSON(w, http.StatusConflict, map[string]any{"error": err.Error(), "usedIn": used})
	case errors.Is(err, services.ErrChallengeFormatMoveTo):
		resp.Error(w, http.StatusBadRequest, err.Error())
	case err != nil:
		resp.Error(w, http.StatusInternalServerError, err.Error())
	default:
		resp.JSON(w, http.StatusOK, map[string]any{"deleted": id, "moved": used})
	}
}
