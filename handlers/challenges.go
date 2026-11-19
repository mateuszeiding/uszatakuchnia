// Package handlers — handlery HTTP podpięte pod istniejącą funkcję Vercela zamiast własnej
// (limit 12 funkcji na wdrożenie; w api/ jest już 12). Funkcja w api/, która je woła,
// zaczyna od resp.Begin (CORS + bramka zapisu), więc tu tego nie powtarzamy.
//
// Trasy (vercel.json → /api/taxonomy/units?kind=…):
//
//	GET    /api/challenges/list                               lista, najnowsze na górze
//	GET    /api/challenges/:id                                szczegóły (404, gdy nie ma)
//	POST   /api/challenges                                    nowe wyzwanie
//	PUT    /api/challenges/:id                                zapis
//	DELETE /api/challenges/:id
//	POST   /api/taxonomy/challenge-formats                    {name}
//	PUT    /api/taxonomy/challenge-formats/:id                {name}
//	DELETE /api/taxonomy/challenge-formats/:id?moveTo=<id>    używany: wyzwania przechodzą do moveTo
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

func queryID(r *http.Request) uint {
	id, _ := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	return uint(id)
}

// Challenges — lista, szczegóły i zapis wyzwań.
func Challenges(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := queryID(r)

	switch r.Method {
	case http.MethodGet:
		if id == 0 {
			list, err := services.LoadChallenges(conn)
			if err != nil {
				resp.Error(w, http.StatusInternalServerError, err.Error())
				return
			}
			out := make([]dtos.ChallengeDto, 0, len(list))
			for _, c := range list {
				out = append(out, mappers.MapChallengeToDto(c))
			}
			resp.JSON(w, http.StatusOK, out)
			return
		}
		c, err := services.LoadChallenge(conn, id)
		switch {
		case errors.Is(err, services.ErrChallengeNotFound):
			resp.Error(w, http.StatusNotFound, err.Error())
		case err != nil:
			resp.Error(w, http.StatusInternalServerError, err.Error())
		default:
			resp.JSON(w, http.StatusOK, mappers.MapChallengeToDto(c))
		}

	case http.MethodPost, http.MethodPut:
		if r.Method == http.MethodPut && id == 0 {
			resp.Error(w, http.StatusBadRequest, "id is required")
			return
		}
		if r.Method == http.MethodPost {
			id = 0
		}
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
		if r.Method == http.MethodPost {
			status = http.StatusCreated
		}
		resp.JSON(w, status, mappers.MapChallengeToDto(c))

	case http.MethodDelete:
		err := services.DeleteChallenge(conn, id)
		switch {
		case errors.Is(err, services.ErrChallengeNotFound):
			resp.Error(w, http.StatusNotFound, err.Error())
		case err != nil:
			resp.Error(w, http.StatusInternalServerError, err.Error())
		default:
			resp.JSON(w, http.StatusOK, map[string]any{"deleted": id})
		}

	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}

// ChallengeFormats — słownik formatów (lista z liczbą użyć jest w GET /api/taxonomy/list).
func ChallengeFormats(w http.ResponseWriter, r *http.Request) {
	conn := db.DB()
	id := queryID(r)

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
		f, err := services.SaveChallengeFormat(conn, id, req.Name)
		switch {
		case errors.Is(err, services.ErrChallengeFormatNotFound):
			resp.Error(w, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrChallengeFormatTaken):
			resp.Error(w, http.StatusConflict, err.Error())
		case err != nil:
			resp.Error(w, http.StatusBadRequest, err.Error())
		default:
			status := http.StatusOK
			if r.Method == http.MethodPost {
				status = http.StatusCreated
			}
			resp.JSON(w, status, mappers.MapChallengeFormatToDto(f))
		}

	case http.MethodDelete:
		moveTo, _ := strconv.ParseUint(r.URL.Query().Get("moveTo"), 10, 64)
		used, err := services.DeleteChallengeFormat(conn, id, uint(moveTo))
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

	default:
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	}
}
