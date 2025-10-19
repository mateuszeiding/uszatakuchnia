package dtos

type IngredientDto struct {
	ID         uint       `json:"id"`
	Name       string     `json:"name"`
	ParentID   *uint      `json:"parentId"`
	Aromas     []AromaDto `json:"aromas"`
	IsAllergen bool       `json:"isAllergen"`
	Type       string     `json:"type"`
	Tastes     []TasteDto `json:"tastes"`
}
