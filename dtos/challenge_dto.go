package dtos

import "time"

// ChallengeIngredientDto — pozycja tekstowa; Required = „Musisz użyć”.
type ChallengeIngredientDto struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

// ChallengeDto — wyzwanie (lista i szczegóły mają ten sam kształt; losowanie jest po stronie klienta).
type ChallengeDto struct {
	ID           uint                     `json:"id"`
	Title        string                   `json:"title"`
	FormatID     uint                     `json:"formatId"`
	Format       string                   `json:"format"`
	Description  string                   `json:"description"`
	TimeLimitMin *int                     `json:"timeLimitMin"`
	Difficulty   int                      `json:"difficulty"`
	Source       *string                  `json:"source"`
	Ingredients  []ChallengeIngredientDto `json:"ingredients"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

// UpsertChallengeRequest — POST /api/challenges, PUT /api/challenges/:id.
// Format: formatId albo nazwa w `format` (kopie w content/challenges/ trzymają nazwę, nie ID).
type UpsertChallengeRequest struct {
	Title        string                   `json:"title"`
	FormatID     uint                     `json:"formatId"`
	Format       string                   `json:"format"`
	Description  string                   `json:"description"`
	TimeLimitMin *int                     `json:"timeLimitMin"`
	Difficulty   int                      `json:"difficulty"`
	Source       *string                  `json:"source"`
	Ingredients  []ChallengeIngredientDto `json:"ingredients"`
}

// ChallengeFormatDto — pozycja słownika formatów; UsedIn = w ilu wyzwaniach.
type ChallengeFormatDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	UsedIn int    `json:"usedIn"`
}
