package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// MapChallengeToDto wymaga preloadu Format i Ingredients (kolejność sort_order).
func MapChallengeToDto(c entities.Challenge) dtos.ChallengeDto {
	ings := make([]dtos.ChallengeIngredientDto, 0, len(c.Ingredients))
	for _, i := range c.Ingredients {
		ings = append(ings, dtos.ChallengeIngredientDto{Name: i.Name, Required: i.Required})
	}
	return dtos.ChallengeDto{
		ID:           c.ID,
		Title:        c.Title,
		FormatID:     c.FormatID,
		Format:       c.Format.Name,
		Description:  c.Description,
		TimeLimitMin: c.TimeLimitMin,
		Difficulty:   c.Difficulty,
		Source:       c.Source,
		Ingredients:  ings,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func MapChallengeFormatToDto(f entities.ChallengeFormat) dtos.ChallengeFormatDto {
	return dtos.ChallengeFormatDto{ID: f.ID, Name: f.Name}
}
