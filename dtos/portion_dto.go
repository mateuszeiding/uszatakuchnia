package dtos

// PortionDto: what one serving of the recipe is. mode "servings" ignores the rest; "sized" uses pieceUnit +
// sizes (pieces in the recipe = servings); "form" uses form.
type PortionDto struct {
	Mode      string         `json:"mode"`
	PieceUnit *UnitDto       `json:"pieceUnit"`
	Sizes     []PieceSizeDto `json:"sizes"`
	Form      *FormDto       `json:"form"`
}

type PieceSizeDto struct {
	DiameterCm int  `json:"diameterCm"`
	DoughG     int  `json:"doughG"`
	IsBase     bool `json:"isBase"`
}

// FormDto: round uses dCm, rect uses aCm × bCm. Pieces = pieces cut from the form (for per-piece macros); nil = not given.
type FormDto struct {
	Shape  string `json:"shape"`
	DCm    *int   `json:"dCm"`
	ACm    *int   `json:"aCm"`
	BCm    *int   `json:"bCm"`
	Pieces *int   `json:"pieces"`
}

// UpsertPortionRequest replaces the whole portion (fields of every mode, so switching modes in the form loses nothing).
type UpsertPortionRequest struct {
	Mode        string         `json:"mode"`
	PieceUnitId *uint          `json:"pieceUnitId"`
	Sizes       []PieceSizeDto `json:"sizes"`
	Form        *FormDto       `json:"form"`
}
