package dtos

import "time"

type ChallengeIngredientDto struct {
	Name     string `json:"name"`
	Required bool   `json:"required"`
}

// The list and the details share one shape; the random pick happens client-side.
type ChallengeDto struct {
	ID           uint                     `json:"id"`
	Title        string                   `json:"title"`
	FormatID     uint                     `json:"formatId"`
	Format       string                   `json:"format"`
	Description  string                   `json:"description"`
	TimeLimitMin *int                     `json:"timeLimitMin"`
	Difficulty   int                      `json:"difficulty"`
	Source       *string                  `json:"source"`
	Cautions     []string                 `json:"cautions"`
	Ingredients  []ChallengeIngredientDto `json:"ingredients"`
	CreatedAt    time.Time                `json:"createdAt"`
	UpdatedAt    time.Time                `json:"updatedAt"`
}

// Format: formatId or a name in `format` (backups in content/challenges/ store the name, not the ID).
type UpsertChallengeRequest struct {
	Title        string                   `json:"title"`
	FormatID     uint                     `json:"formatId"`
	Format       string                   `json:"format"`
	Description  string                   `json:"description"`
	TimeLimitMin *int                     `json:"timeLimitMin"`
	Difficulty   int                      `json:"difficulty"`
	Source       *string                  `json:"source"`
	Cautions     *[]string                `json:"cautions"` // nil = unchanged (PUT); on create nil = []
	Ingredients  []ChallengeIngredientDto `json:"ingredients"`
}

type ChallengeFormatDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	UsedIn int    `json:"usedIn"`
}
