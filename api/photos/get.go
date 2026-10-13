package api

import (
	"net/http"
	"strconv"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	resp "uszatakuchnia/http"
)

// Get — GET /api/photos/:id. Zdjęcia są niezmienne (nowe zdjęcie = nowe id), więc cache na rok.
func Get(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	id, err := strconv.ParseUint(r.URL.Query().Get("id"), 10, 64)
	if err != nil {
		resp.Error(w, http.StatusBadRequest, "invalid id")
		return
	}

	var photo entities.Photo
	if err := db.DB().First(&photo, id).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "photo not found")
		return
	}

	w.Header().Set("Content-Type", photo.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(photo.Data)
}
