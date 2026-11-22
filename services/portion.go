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

// IsPieceUnit: a unit that can name pieces — inflected (1 / 2–4 / 5+) and with a quantity.
func IsPieceUnit(u entities.Unit) bool {
	return u.FormOne != nil && u.FormFew != nil && u.FormMany != nil && !u.NoQty
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

	// Fields of the other modes are saved too, so their messages name the mode (and the size row): the author may
	// not see those fields right now.
	const sizedName, formName = "Sztuki w rozmiarach", "Forma lub blacha"

	if p.PieceUnitId != nil {
		u, ok := units[*p.PieceUnitId]
		if !ok {
			add("portion.pieceUnitId", sizedName+": nieznana jednostka sztuki")
		} else if !IsPieceUnit(u) {
			add("portion.pieceUnitId", fmt.Sprintf("%s: jednostka „%s” nie ma odmiany (np. placek / placki / placków)", sizedName, u.Name))
		}
	}

	seen := map[int]bool{}
	bases := 0
	for i, s := range p.Sizes {
		f := fmt.Sprintf("portion.sizes[%d]", i)
		d := s.DiameterCm
		row := fmt.Sprintf("%s, rozmiar %d", sizedName, i+1)
		if cmInRange(&d) {
			row = fmt.Sprintf("%s, rozmiar %d cm", sizedName, d)
		}
		if !cmInRange(&d) {
			add(f+".diameterCm", row+": średnica "+rangeMsg)
		} else if seen[d] {
			add(f+".diameterCm", fmt.Sprintf("%s: rozmiar %d cm jest już na liście", sizedName, d))
		}
		seen[d] = true
		if s.DoughG <= 0 {
			add(f+".doughG", row+": podaj masę ciasta na sztukę")
		}
		if s.IsBase {
			bases++
		}
	}
	if bases > 1 {
		add("portion.sizes", sizedName+": zaznacz jeden rozmiar „w przepisie”")
	}

	if fm := p.Form; fm != nil {
		switch fm.Shape {
		case entities.FormRound:
			if fm.DCm != nil && !cmInRange(fm.DCm) {
				add("portion.form.dCm", formName+": średnica "+rangeMsg)
			}
		case entities.FormRect:
			if fm.ACm != nil && !cmInRange(fm.ACm) {
				add("portion.form.aCm", formName+": długość "+rangeMsg)
			}
			if fm.BCm != nil && !cmInRange(fm.BCm) {
				add("portion.form.bCm", formName+": szerokość "+rangeMsg)
			}
		default:
			add("portion.form.shape", formName+": wybierz kształt formy")
		}
		if fm.Pieces != nil && (*fm.Pieces < 1 || *fm.Pieces > 100) {
			add("portion.form.pieces", formName+": liczba kawałków od 1 do 100")
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
