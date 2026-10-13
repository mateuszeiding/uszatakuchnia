package entities

import "time"

// Photo — zdjęcie wgrane przez admina (kadrowane po stronie klienta do 1200×1200 JPEG).
// Serwowane przez GET /api/photos/:id. RecipePhoto.ImageURL wskazuje na ten adres.
type Photo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ContentType string    `gorm:"type:varchar(32);not null"`
	Data        []byte    `gorm:"type:bytea;not null"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (Photo) TableName() string { return "photos" }
