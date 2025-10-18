package taste

type TasteType int

//go:generate stringer -type=TasteType
const (
	SWEET TasteType = iota
	SOUR
	SALTY
	BITTER
	UMAMI
)
