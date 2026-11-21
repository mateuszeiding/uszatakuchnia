package handlers

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"uszatakuchnia/db"
	"uszatakuchnia/db/entities"
	resp "uszatakuchnia/http"
)

const maxPhotoBytes = 4 << 20 // Vercel's body limit is 4.5 MB; the client sends a ~200–500 KB JPEG

var allowedPhotoTypes = map[string]bool{"image/jpeg": true, "image/png": true, "image/webp": true}

// PhotoGet: photos are immutable (new photo = new id), so they're cached for a year.
func PhotoGet(w http.ResponseWriter, r *http.Request) {
	var photo entities.Photo
	if err := db.DB().First(&photo, pathID(r)).Error; err != nil {
		resp.Error(w, http.StatusNotFound, "photo not found")
		return
	}

	w.Header().Set("Content-Type", photo.ContentType)
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(photo.Data)
}

// PhotoUpload: body = raw image bytes, Content-Type = the image type.
func PhotoUpload(w http.ResponseWriter, r *http.Request) {
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
