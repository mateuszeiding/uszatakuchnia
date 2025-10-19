package entities

import "uszatakuchnia/dtos"

type Image struct {
	ID      uint    `gorm:"primaryKey"`
	RefType RefType `gorm:"type:varchar(32);not null;uniqueIndex:ux_images_ref,priority:1"`
	RefID   uint    `gorm:"not null;uniqueIndex:ux_images_ref,priority:2"`

	Provider ImageProvider `gorm:"type:varchar(32);not null"`
	PhotoID  string        `gorm:"not null"`

	ImageURLSmall   string `gorm:"not null"`
	ImageURLRegular string `gorm:"not null"`
	ImageURLRaw     string `gorm:"not null"`

	AuthorName       string `gorm:"not null"`
	AuthorUsername   string `gorm:"not null"`
	AuthorProfileURL string `gorm:"not null"`
	PhotoPageURL     string `gorm:"not null"`
}
type ImageProvider string

const (
	Unsplash ImageProvider = "unsplash"
)

type RefType string

const (
	IngredientRef RefType = "ingredient"
)

func (i *Image) ToDto() *dtos.ImageDto {
	if i == nil {
		return nil
	}

	dto := &dtos.ImageDto{
		URLs: dtos.ImageURLs{
			Small:   i.ImageURLSmall,
			Regular: i.ImageURLRegular,
			Raw:     i.ImageURLRaw,
		},
		Author: dtos.ImageAuthor{
			Name:       i.AuthorName,
			Username:   i.AuthorUsername,
			ProfileURL: i.AuthorProfileURL,
		},
	}

	return dto
}
