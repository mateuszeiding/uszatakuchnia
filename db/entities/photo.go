package entities

import "time"

// Photo uploaded by the admin (cropped client-side to a 1200×1200 JPEG).
// Served by GET /api/photos/{id}; RecipePhoto.ImageURL points there.
type Photo struct {
	ID          uint      `gorm:"primaryKey;autoIncrement"`
	ContentType string    `gorm:"type:varchar(32);not null"`
	Data        []byte    `gorm:"type:bytea;not null"`
	CreatedAt   time.Time `gorm:"not null"`
}

func (Photo) TableName() string { return "photos" }
