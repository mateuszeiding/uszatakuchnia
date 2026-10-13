package api

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	resp "uszatakuchnia/http"
)

const maxPhotoBytes = 4 << 20 // limit body Vercel to 4,5 MB; klient wysyła ~200–500 KB JPEG

var allowedPhotoTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}

// Upload — POST /api/photos, body = surowe bajty obrazu, Content-Type = typ obrazu.
func Upload(w http.ResponseWriter, r *http.Request) {
	if resp.HandleCORS(w, r) {
		return
	}
	if r.Method != http.MethodPost {
		resp.Error(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if !allowedPhotoTypes[ct] {
		resp.Error(w, http.StatusUnsupportedMediaType, "unsupported image type")
		return
	}

	data, err := io.ReadAll(http.MaxBytesReader(w, r.Body, maxPhotoBytes))
	if err != nil {
		resp.Error(w, http.StatusRequestEntityTooLarge, "image too large")
		return
	}
	if len(data) == 0 {
		resp.Error(w, http.StatusBadRequest, "empty body")
		return
	}

	photo := entities.Photo{ContentType: ct, Data: data, CreatedAt: time.Now()}
	if err := db.DB().Create(&photo).Error; err != nil {
		resp.Error(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp.JSON(w, http.StatusCreated, map[string]any{
		"id":  photo.ID,
		"url": fmt.Sprintf("/api/photos/%d", photo.ID),
	})
}
