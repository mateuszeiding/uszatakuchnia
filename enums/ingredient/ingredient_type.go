package ingredient

type IngredientType int

//go:generate stringer -type=IngredientType
const (
	VEGETABLE IngredientType = iota
	FRUIT
	HERB
	SPICE
	GRAIN
	LEGUME
	NUT
	SEED
	DAIRY
	MEAT
	FISH
	FUNGUS
	OTHER
)
