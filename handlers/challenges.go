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

// ─── Challenges ─────────────────
//
//	GET    /api/challenges/list   newest first (also GET /api/challenges)
//	GET    /api/challenges/{id}   404 when missing
//	POST   /api/challenges        public right away (no draft)

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

func ChallengeCreate(w http.ResponseWriter, r *http.Request) { saveChallenge(w, r, 0) }

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

// ─── Challenge formats (admin → Classification; the list with usage counts is in GET /api/taxonomy/list) ───
//
//	POST   /api/taxonomy/challenge-formats                  {name}
//	PUT    /api/taxonomy/challenge-formats/{id}             {name}
//	DELETE /api/taxonomy/challenge-formats/{id}?moveTo=<id>  in use: its challenges move to moveTo

func ChallengeFormatCreate(w http.ResponseWriter, r *http.Request) { saveChallengeFormat(w, r, 0) }

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
