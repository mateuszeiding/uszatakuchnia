package enum

type AromaType int

//go:generate stringer -type=AromaType
const (
	CITRUS AromaType = iota
	FRUITY
	FLORAL
	HERBAL
	GREEN
	SPICY
	SULFUROUS
	EARTHY
	WOODY
	NUTTY
	SMOKY
	CREAMY
	FERMENTED
	OTHER
)
