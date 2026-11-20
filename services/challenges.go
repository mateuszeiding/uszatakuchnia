package services

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// Limits match the database columns (characters).
const (
	MaxChallengeTitleLen      = 255
	MaxChallengeSourceLen     = 255
	MaxChallengeIngredientLen = 120
	MaxChallengeFormatNameLen = 64
	MaxChallengeTimeLimitMin  = 7 * 24 * 60
	MinChallengeDifficulty    = 1
	MaxChallengeDifficulty    = 5
)

var (
	ErrChallengeNotFound       = errors.New("challenge not found")
	ErrChallengeFormatName     = errors.New("name is required")
	ErrChallengeFormatTaken    = errors.New("format already exists")
	ErrChallengeFormatNotFound = errors.New("format not found")
	ErrChallengeFormatInUse    = errors.New("format is used in challenges; pass moveTo")
	ErrChallengeFormatMoveTo   = errors.New("moveTo must be another format")
)

var manyBlankLines = regexp.MustCompile(`\n[ \t]*\n(?:[ \t]*\n)+`)

// NormalizeChallengeRequest trims fields, normalizes description paragraphs (one blank line between them),
// drops empty ingredient items; an empty source = none.
func NormalizeChallengeRequest(req dtos.UpsertChallengeRequest) dtos.UpsertChallengeRequest {
	req.Title = strings.Join(strings.Fields(req.Title), " ")
	req.Format = strings.Join(strings.Fields(req.Format), " ")

	desc := strings.ReplaceAll(req.Description, "\r\n", "\n")
	lines := strings.Split(desc, "\n")
	for i, l := range lines {
		lines[i] = strings.TrimRight(l, " \t")
	}
	desc = strings.TrimSpace(strings.Join(lines, "\n"))
	req.Description = manyBlankLines.ReplaceAllString(desc, "\n\n")

	if req.Source != nil {
		s := strings.Join(strings.Fields(*req.Source), " ")
		if s == "" {
			req.Source = nil
		} else {
			req.Source = &s
		}
	}

	ings := make([]dtos.ChallengeIngredientDto, 0, len(req.Ingredients))
	for _, i := range req.Ingredients {
		name := strings.Join(strings.Fields(i.Name), " ")
		if name == "" {
			continue
		}
		ings = append(ings, dtos.ChallengeIngredientDto{Name: name, Required: i.Required})
	}
	req.Ingredients = ings
	if req.Cautions != nil {
		c := CleanCautions(*req.Cautions)
		req.Cautions = &c
	}
	return req
}

// validateChallengeFields validates without the database (ValidateChallenge checks the format).
func validateChallengeFields(req dtos.UpsertChallengeRequest) []dtos.FieldError {
	errs := []dtos.FieldError{}
	add := func(field, msg string) { errs = append(errs, dtos.FieldError{Field: field, Message: msg}) }

	switch n := utf8.RuneCountInString(req.Title); {
	case n == 0:
		add("title", "Wpisz tytuł")
	case n > MaxChallengeTitleLen:
		add("title", fmt.Sprintf("Tytuł może mieć najwyżej %d znaków", MaxChallengeTitleLen))
	}
	if req.Description == "" {
		add("description", "Opisz zadanie")
	}
	if req.Difficulty < MinChallengeDifficulty || req.Difficulty > MaxChallengeDifficulty {
		add("difficulty", "Wybierz trudność od 1 do 5")
	}
	if req.TimeLimitMin != nil {
		if *req.TimeLimitMin <= 0 {
			add("timeLimitMin", "Limit czasu musi być większy od zera (puste = bez limitu)")
		} else if *req.TimeLimitMin > MaxChallengeTimeLimitMin {
			add("timeLimitMin", fmt.Sprintf("Limit czasu może mieć najwyżej %d minut", MaxChallengeTimeLimitMin))
		}
	}
	if req.Source != nil && utf8.RuneCountInString(*req.Source) > MaxChallengeSourceLen {
		add("source", fmt.Sprintf("Źródło może mieć najwyżej %d znaków", MaxChallengeSourceLen))
	}
	for i, ing := range req.Ingredients {
		if utf8.RuneCountInString(ing.Name) > MaxChallengeIngredientLen {
			add(fmt.Sprintf("ingredients[%d].name", i),
				fmt.Sprintf("Składnik może mieć najwyżej %d znaków", MaxChallengeIngredientLen))
		}
	}
	if req.Cautions != nil {
		validateCautions(*req.Cautions, add)
	}
	if req.FormatID == 0 && req.Format == "" {
		add("formatId", "Wybierz format")
	}
	return errs
}

// ValidateChallenge expects req after NormalizeChallengeRequest. A format given by name is turned into formatId.
func ValidateChallenge(tx *gorm.DB, req *dtos.UpsertChallengeRequest) ([]dtos.FieldError, error) {
	errs := validateChallengeFields(*req)
	if req.FormatID == 0 && req.Format == "" {
		return errs, nil
	}
	var f entities.ChallengeFormat
	q := tx.Model(&entities.ChallengeFormat{})
	if req.FormatID != 0 {
		q = q.Where("id = ?", req.FormatID)
	} else {
		q = q.Where("lower(name) = lower(?)", req.Format)
	}
	res := q.Limit(1).Find(&f)
	if res.Error != nil {
		return nil, res.Error
	}
	if res.RowsAffected == 0 {
		errs = append(errs, dtos.FieldError{Field: "formatId", Message: "Nie ma takiego formatu"})
	} else {
		req.FormatID = f.ID
	}
	return errs, nil
}

// LoadChallenges returns all challenges, newest first.
func LoadChallenges(tx *gorm.DB) ([]entities.Challenge, error) {
	var list []entities.Challenge
	err := tx.Preload("Format").
		Preload("Ingredients", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order") }).
		Order("created_at DESC, id DESC").
		Find(&list).Error
	return list, err
}

func LoadChallenge(tx *gorm.DB, id uint) (entities.Challenge, error) {
	var c entities.Challenge
	res := tx.Preload("Format").
		Preload("Ingredients", func(tx *gorm.DB) *gorm.DB { return tx.Order("sort_order") }).
		Limit(1).Find(&c, id)
	if res.Error != nil {
		return c, res.Error
	}
	if res.RowsAffected == 0 {
		return c, ErrChallengeNotFound
	}
	return c, nil
}

// SaveChallenge creates (id 0) or overwrites a challenge along with its ingredient list.
// req must have gone through NormalizeChallengeRequest and ValidateChallenge (formatId set).
func SaveChallenge(conn *gorm.DB, id uint, req dtos.UpsertChallengeRequest) (uint, error) {
	err := conn.Transaction(func(tx *gorm.DB) error {
		c := entities.Challenge{
			Title:        req.Title,
			FormatID:     req.FormatID,
			Description:  req.Description,
			TimeLimitMin: req.TimeLimitMin,
			Difficulty:   req.Difficulty,
			Source:       req.Source,
		}
		if req.Cautions != nil {
			c.Cautions = *req.Cautions
		}
		if id == 0 {
			if err := tx.Omit("Format", "Ingredients").Create(&c).Error; err != nil {
				return err
			}
			id = c.ID
		} else {
			cols := []any{"format_id", "description", "time_limit_min", "difficulty", "source", "updated_at"}
			if req.Cautions != nil {
				cols = append(cols, "cautions")
			}
			res := tx.Model(&entities.Challenge{}).Where("id = ?", id).
				Select("title", cols...).
				Updates(&c)
			if res.Error != nil {
				return res.Error
			}
			if res.RowsAffected == 0 {
				return ErrChallengeNotFound
			}
			if err := tx.Where("challenge_id = ?", id).Delete(&entities.ChallengeIngredient{}).Error; err != nil {
				return err
			}
		}
		for i, ing := range req.Ingredients {
			row := entities.ChallengeIngredient{ChallengeID: id, SortOrder: i + 1, Name: ing.Name, Required: ing.Required}
			if err := tx.Create(&row).Error; err != nil {
				return err
			}
		}
		return nil
	})
	return id, err
}

func DeleteChallenge(tx *gorm.DB, id uint) error {
	res := tx.Delete(&entities.Challenge{}, id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return ErrChallengeNotFound
	}
	return nil
}

// ─── Formats (admin → Classification → Challenge formats) ─────────

func cleanChallengeFormatName(name string) (string, error) {
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return "", ErrChallengeFormatName
	}
	if utf8.RuneCountInString(name) > MaxChallengeFormatNameLen {
		return "", fmt.Errorf("name is longer than %d characters", MaxChallengeFormatNameLen)
	}
	return name, nil
}

// SaveChallengeFormat creates (id 0, last in order) or renames a format.
func SaveChallengeFormat(tx *gorm.DB, id uint, rawName string) (entities.ChallengeFormat, error) {
	name, err := cleanChallengeFormatName(rawName)
	if err != nil {
		return entities.ChallengeFormat{}, err
	}
	var f entities.ChallengeFormat
	if id != 0 {
		if err := tx.First(&f, id).Error; err != nil {
			return f, ErrChallengeFormatNotFound
		}
	} else {
		tx.Model(&entities.ChallengeFormat{}).Select("COALESCE(MAX(sort_order), 0) + 1").Scan(&f.SortOrder)
	}
	var n int64
	if err := tx.Model(&entities.ChallengeFormat{}).
		Where("lower(name) = lower(?) AND id <> ?", name, f.ID).Count(&n).Error; err != nil {
		return f, err
	}
	if n > 0 {
		return f, ErrChallengeFormatTaken
	}
	f.Name = name
	return f, tx.Save(&f).Error
}

// ChallengeFormatUsage: number of challenges using each format.
func ChallengeFormatUsage(tx *gorm.DB) (map[uint]int, error) {
	var rows []struct {
		ID uint
		N  int
	}
	if err := tx.Model(&entities.Challenge{}).Select("format_id AS id, COUNT(*) AS n").
		Group("format_id").Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := map[uint]int{}
	for _, r := range rows {
		out[r.ID] = r.N
	}
	return out, nil
}

// DeleteChallengeFormat deletes a format. One in use needs moveTo (a challenge requires a format):
// its challenges move to the other format. Returns the number of challenges moved.
func DeleteChallengeFormat(conn *gorm.DB, id, moveTo uint) (int64, error) {
	var f entities.ChallengeFormat
	if err := conn.First(&f, id).Error; err != nil {
		return 0, ErrChallengeFormatNotFound
	}
	if moveTo != 0 {
		var target entities.ChallengeFormat
		if moveTo == id || conn.First(&target, moveTo).Error != nil {
			return 0, ErrChallengeFormatMoveTo
		}
	}
	var used int64
	if err := conn.Model(&entities.Challenge{}).Where("format_id = ?", id).Count(&used).Error; err != nil {
		return 0, err
	}
	if used > 0 && moveTo == 0 {
		return used, ErrChallengeFormatInUse
	}
	err := conn.Transaction(func(tx *gorm.DB) error {
		if used > 0 {
			// no UpdatedAt: a dictionary change isn't a challenge edit
			if err := tx.Exec(`UPDATE challenges SET format_id = ? WHERE format_id = ?`, moveTo, id).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&entities.ChallengeFormat{}, id).Error
	})
	return used, err
}
