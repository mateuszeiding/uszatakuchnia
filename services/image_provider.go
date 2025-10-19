package services

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"uszatakuchnia/db/entities"
	unsplash "uszatakuchnia/db/helpers"
)

func UpsertImageFromUnsplash(tx *gorm.DB, refID uint, refType entities.RefType, photoID string) error {
	meta, err := unsplash.FetchPhotoMeta(photoID)
	if err != nil {
		return err
	}

	im := entities.Image{
		RefType:          refType,
		RefID:            refID,
		Provider:         entities.Unsplash,
		PhotoID:          photoID,
		ImageURLSmall:    meta.Small,
		ImageURLRegular:  meta.Regular,
		ImageURLRaw:      meta.Raw,
		AuthorName:       meta.AuthorName,
		AuthorUsername:   meta.AuthorUsername,
		AuthorProfileURL: meta.AuthorProfileURL,
		PhotoPageURL:     meta.PhotoPageURL,
	}

	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{
			{Name: "ref_type"},
			{Name: "ref_id"},
		},
		DoUpdates: clause.Assignments(map[string]any{
			"provider":           im.Provider,
			"photo_id":           im.PhotoID,
			"image_url_small":    im.ImageURLSmall,
			"image_url_regular":  im.ImageURLRegular,
			"image_url_raw":      im.ImageURLRaw,
			"author_name":        im.AuthorName,
			"author_username":    im.AuthorUsername,
			"author_profile_url": im.AuthorProfileURL,
			"photo_page_url":     im.PhotoPageURL,
		}),
	}).Create(&im).Error
}
