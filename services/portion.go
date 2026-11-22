package services

import (
	"fmt"
	"sort"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// effectivePortion is the portion the saved recipe will have: the request's, or (nil = unchanged) the stored one.
func effectivePortion(tx *gorm.DB, req dtos.UpsertRecipeRequest, selfID uint) (dtos.UpsertPortionRequest, error) {
	if req.Portion != nil {
		return *req.Portion, nil
	}
	if selfID == 0 {
		return dtos.UpsertPortionRequest{Mode: entities.PortionServings}, nil
	}
	var e entities.Recipe
	if err := tx.Preload("PieceSizes").Take(&e, selfID).Error; err != nil {
		return dtos.UpsertPortionRequest{}, err
	}
	return StoredPortion(e), nil
}

// StoredPortion: the recipe's portion as a request (PieceSizes preloaded).
func StoredPortion(e entities.Recipe) dtos.UpsertPortionRequest {
	p := dtos.UpsertPortionRequest{Mode: e.PortionMode, PieceUnitId: e.PieceUnitID, Sizes: []dtos.PieceSizeDto{}}
	if p.Mode == "" {
		p.Mode = entities.PortionServings
	}
	for _, s := range e.PieceSizes {
		p.Sizes = append(p.Sizes, dtos.PieceSizeDto{DiameterCm: s.DiameterCm, DoughG: s.DoughG, IsBase: s.IsBase})
	}
	if e.FormShape != nil {
		p.Form = &dtos.FormDto{Shape: *e.FormShape, DCm: e.FormDCm, ACm: e.FormACm, BCm: e.FormBCm, Pieces: e.FormPieces}
	}
	return p
}

func cmInRange(v *int) bool {
	return v != nil && *v >= entities.PortionMinCm && *v <= entities.PortionMaxCm
}

// validatePortion: data of every mode must be valid even in a draft (it's stored whatever the mode);
// publishing needs what the chosen mode uses.
func validatePortion(p dtos.UpsertPortionRequest, servings int, publish bool, units map[uint]entities.Unit, add func(field, msg string)) {
	rangeMsg := fmt.Sprintf("od %d do %d cm", entities.PortionMinCm, entities.PortionMaxCm)
	switch p.Mode {
	case entities.PortionServings, entities.PortionSized, entities.PortionForm:
	default:
		add("portion.mode", "Nieznany rodzaj porcji")
	}

	if p.PieceUnitId != nil {
		u, ok := units[*p.PieceUnitId]
		if !ok {
			add("portion.pieceUnitId", "Nieznana jednostka")
		} else if u.FormOne == nil || u.FormFew == nil || u.FormMany == nil || u.NoQty {
			add("portion.pieceUnitId", fmt.Sprintf("Jednostka „%s” nie ma odmiany (np. placek / placki / placków)", u.Name))
		}
	}

	seen := map[int]bool{}
	bases := 0
	for i, s := range p.Sizes {
		f := fmt.Sprintf("portion.sizes[%d]", i)
		d := s.DiameterCm
		if !cmInRange(&d) {
			add(f+".diameterCm", "Średnica "+rangeMsg)
		} else if seen[d] {
			add(f+".diameterCm", fmt.Sprintf("Rozmiar %d cm jest już na liście", d))
		}
		seen[d] = true
		if s.DoughG <= 0 {
			add(f+".doughG", "Podaj masę ciasta na sztukę")
		}
		if s.IsBase {
			bases++
		}
	}
	if bases > 1 {
		add("portion.sizes", "Zaznacz jeden rozmiar „w przepisie”")
	}

	if fm := p.Form; fm != nil {
		switch fm.Shape {
		case entities.FormRound:
			if fm.DCm != nil && !cmInRange(fm.DCm) {
				add("portion.form.dCm", "Średnica "+rangeMsg)
			}
		case entities.FormRect:
			if fm.ACm != nil && !cmInRange(fm.ACm) {
				add("portion.form.aCm", "Długość "+rangeMsg)
			}
			if fm.BCm != nil && !cmInRange(fm.BCm) {
				add("portion.form.bCm", "Szerokość "+rangeMsg)
			}
		default:
			add("portion.form.shape", "Wybierz kształt formy")
		}
		if fm.Pieces != nil && (*fm.Pieces < 1 || *fm.Pieces > 100) {
			add("portion.form.pieces", "Liczba kawałków od 1 do 100")
		}
	}

	if !publish {
		return
	}
	switch p.Mode {
	case entities.PortionServings:
		if servings <= 0 {
			add("servings", "Podaj liczbę porcji")
		}
	case entities.PortionSized:
		if servings <= 0 {
			add("servings", "Podaj liczbę sztuk w przepisie")
		}
		if p.PieceUnitId == nil {
			add("portion.pieceUnitId", "Wybierz nazwę sztuki")
		}
		if len(p.Sizes) == 0 {
			add("portion.sizes", "Dodaj co najmniej jeden rozmiar")
		} else if bases == 0 {
			add("portion.sizes", "Zaznacz rozmiar „w przepisie”")
		}
	case entities.PortionForm:
		fm := p.Form
		switch {
		case fm == nil:
			add("portion.form.shape", "Wybierz kształt formy")
		case fm.Shape == entities.FormRound && fm.DCm == nil:
			add("portion.form.dCm", "Podaj średnicę formy")
		case fm.Shape == entities.FormRect && (fm.ACm == nil || fm.BCm == nil):
			add("portion.form.aCm", "Podaj długość i szerokość formy")
		}
	}
}

// applyPortion copies the portion onto the recipe (sizes go in replacePieceSizes).
func applyPortion(e *entities.Recipe, p *dtos.UpsertPortionRequest) {
	if p == nil {
		if e.PortionMode == "" {
			e.PortionMode = entities.PortionServings
		}
		return
	}
	e.PortionMode = p.Mode
	e.PieceUnitID = p.PieceUnitId
	e.PieceUnit = nil
	e.FormShape, e.FormDCm, e.FormACm, e.FormBCm, e.FormPieces = nil, nil, nil, nil, nil
	if fm := p.Form; fm != nil {
		shape := fm.Shape
		e.FormShape = &shape
		e.FormPieces = fm.Pieces
		// only the dimensions of the chosen shape: a stale diameter must not linger next to a rect form
		if shape == entities.FormRound {
			e.FormDCm = fm.DCm
		} else {
			e.FormACm, e.FormBCm = fm.ACm, fm.BCm
		}
	}
}

func replacePieceSizes(tx *gorm.DB, recipeID uint, p *dtos.UpsertPortionRequest) error {
	if p == nil {
		return nil
	}
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipePieceSize{}).Error; err != nil {
		return err
	}
	// chips on the page go from the smallest size, whatever order the author typed them in
	sizes := append([]dtos.PieceSizeDto{}, p.Sizes...)
	sort.SliceStable(sizes, func(i, j int) bool { return sizes[i].DiameterCm < sizes[j].DiameterCm })
	for i, s := range sizes {
		row := entities.RecipePieceSize{RecipeID: recipeID, DiameterCm: s.DiameterCm, DoughG: s.DoughG, IsBase: s.IsBase, SortOrder: i + 1}
		if err := tx.Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}
