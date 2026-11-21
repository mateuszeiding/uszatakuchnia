package dtos

type GlossaryTermDto struct {
	ID          uint   `json:"id"`
	Name        string `json:"name"`
	Explanation string `json:"explanation"`
	// Recipes (drafts too) that mark the term; owner only, nil for everyone else.
	UsedIn *int `json:"usedIn,omitempty"`
}

type UpsertGlossaryTermRequest struct {
	Name        string `json:"name"`
	Explanation string `json:"explanation"`
}

type GlossaryUsageDto struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Status string `json:"status"`
}
