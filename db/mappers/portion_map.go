package mappers

import (
	"sort"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

type portionFields struct {
	mode            string
	unit            *entities.Unit
	sizes           []entities.RecipePieceSize
	shape           *string
	d, a, b, pieces *int
}

func mapPortion(p portionFields) dtos.PortionDto {
	out := dtos.PortionDto{Mode: p.mode, Sizes: []dtos.PieceSizeDto{}}
	if out.Mode == "" {
		out.Mode = entities.PortionServings
	}
	if p.unit != nil {
		u := MapUnitToDto(*p.unit)
		out.PieceUnit = &u
	}
	sizes := append([]entities.RecipePieceSize{}, p.sizes...)
	sort.SliceStable(sizes, func(i, j int) bool { return sizes[i].SortOrder < sizes[j].SortOrder })
	for _, s := range sizes {
		out.Sizes = append(out.Sizes, dtos.PieceSizeDto{DiameterCm: s.DiameterCm, DoughG: s.DoughG, IsBase: s.IsBase})
	}
	if p.shape != nil {
		out.Form = &dtos.FormDto{Shape: *p.shape, DCm: p.d, ACm: p.a, BCm: p.b, Pieces: p.pieces}
	}
	return out
}

// MapRecipePortion needs the PieceUnit and PieceSizes preloads.
func MapRecipePortion(e entities.Recipe) dtos.PortionDto {
	return mapPortion(portionFields{e.PortionMode, e.PieceUnit, e.PieceSizes, e.FormShape, e.FormDCm, e.FormACm, e.FormBCm, e.FormPieces})
}

func mapRecipeBasePortion(e entities.RecipeBase) dtos.PortionDto {
	return mapPortion(portionFields{e.PortionMode, e.PieceUnit, e.PieceSizes, e.FormShape, e.FormDCm, e.FormACm, e.FormBCm, e.FormPieces})
}
