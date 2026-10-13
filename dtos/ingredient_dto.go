package dtos

type IngredientDto struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	ParentID  *uint     `json:"parentId"`
	Allergens []string  `json:"allergens"` // kody z entities.Allergens
	UsedIn    int       `json:"usedIn"`    // liczba przepisów (tylko lista)
	Type      TypeDto   `json:"type"`
	Image     *ImageDto `json:"image"`
}
