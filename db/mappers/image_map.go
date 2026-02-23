package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

func MapImageToDto(e *entities.Image) *dtos.ImageDto {
	if e == nil {
		return nil
	}

	dto := &dtos.ImageDto{
		URLs: dtos.ImageURLs{
			Small:   e.ImageURLSmall,
			Regular: e.ImageURLRegular,
			Raw:     e.ImageURLRaw,
		},
		Author: dtos.ImageAuthor{
			Name:       e.AuthorName,
			Username:   e.AuthorUsername,
			ProfileURL: e.AuthorProfileURL,
		},
	}

	return dto
}
