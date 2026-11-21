package services

import (
	"errors"
	"fmt"
	"sort"
	"strings"
	"unicode/utf8"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

const (
	MaxTermNameLen        = 60
	MaxTermExplanationLen = 280
)

var (
	ErrTermNotFound = errors.New("term not found")
	ErrTermInvalid  = errors.New("validation failed")
	ErrTermTaken    = errors.New("term already exists")
	ErrTermInUse    = errors.New("term is used in recipes; pass replaceWith or removeMarks")
	ErrTermReplace  = errors.New("replaceWith must be another term")
)

func oneLine(s string) string { return strings.Join(strings.Fields(s), " ") }

// SaveGlossaryTerm creates (id 0) or updates a term. ErrTermInvalid / ErrTermTaken come with field errors.
func SaveGlossaryTerm(tx *gorm.DB, id uint, req dtos.UpsertGlossaryTermRequest) (entities.GlossaryTerm, []dtos.FieldError, error) {
	var t entities.GlossaryTerm
	if id != 0 {
		if err := tx.First(&t, id).Error; err != nil {
			return t, nil, ErrTermNotFound
		}
	}
	name, expl := oneLine(req.Name), oneLine(req.Explanation)
	errs := []dtos.FieldError{}
	add := func(f, m string) { errs = append(errs, dtos.FieldError{Field: f, Message: m}) }
	switch {
	case name == "":
		add("name", "Podaj nazwę pojęcia")
	case utf8.RuneCountInString(name) > MaxTermNameLen:
		add("name", fmt.Sprintf("Nazwa może mieć najwyżej %d znaków", MaxTermNameLen))
	}
	switch {
	case expl == "":
		add("explanation", "Napisz wyjaśnienie")
	case utf8.RuneCountInString(expl) > MaxTermExplanationLen:
		add("explanation", fmt.Sprintf("Wyjaśnienie może mieć najwyżej %d znaków", MaxTermExplanationLen))
	}
	if len(errs) > 0 {
		return t, errs, ErrTermInvalid
	}
	var n int64
	if err := tx.Model(&entities.GlossaryTerm{}).Where("lower(name) = lower(?) AND id <> ?", name, t.ID).Count(&n).Error; err != nil {
		return t, nil, err
	}
	if n > 0 {
		return t, []dtos.FieldError{{Field: "name", Message: "Takie pojęcie już jest w słowniczku"}}, ErrTermTaken
	}
	t.Name, t.Explanation = name, expl
	return t, nil, tx.Save(&t).Error
}

// stepsWithTerms: steps whose text or hints contain a glossary ref (to `id`, or any with id 0).
func stepsWithTerms(tx *gorm.DB, id uint) ([]entities.RecipeStep, error) {
	like := "%{{term:%"
	if id != 0 {
		like = fmt.Sprintf("%%{{term:%d|%%", id)
	}
	var steps []entities.RecipeStep
	err := tx.Where("text LIKE ? OR hints::text LIKE ?", like, like).Order("recipe_id, step_no").Find(&steps).Error
	return steps, err
}

// TermUsage: term id → ids of the recipes (drafts too) that mark it at least once.
func TermUsage(tx *gorm.DB) (map[uint]map[uint]bool, error) {
	steps, err := stepsWithTerms(tx, 0)
	if err != nil {
		return nil, err
	}
	out := map[uint]map[uint]bool{}
	for _, s := range steps {
		for _, id := range TermRefIDs(append([]string{s.Text}, s.Hints...)...) {
			if out[id] == nil {
				out[id] = map[uint]bool{}
			}
			out[id][s.RecipeID] = true
		}
	}
	return out, nil
}

// TermRecipes lists the recipes that use the term, A–Z (the delete dialog names them).
func TermRecipes(tx *gorm.DB, id uint) ([]dtos.GlossaryUsageDto, error) {
	usage, err := TermUsage(tx)
	if err != nil {
		return nil, err
	}
	ids := []uint{}
	for rid := range usage[id] {
		ids = append(ids, rid)
	}
	out := []dtos.GlossaryUsageDto{}
	if len(ids) == 0 {
		return out, nil
	}
	err = tx.Model(&entities.Recipe{}).Select("id, name, status").Where("id IN ?", ids).
		Order("lower(name), id").Scan(&out).Error
	return out, err
}

// DeleteGlossaryTerm deletes a term. If recipes mark it: replaceWith (≠ 0) points the marks at another term,
// removeMarks leaves the words as plain text. Returns the number of recipes rewritten.
func DeleteGlossaryTerm(tx *gorm.DB, id, replaceWith uint, removeMarks bool) (int, error) {
	if replaceWith == id {
		return 0, ErrTermReplace
	}
	forUpdate := clause.Locking{Strength: "UPDATE"}
	recipes := 0
	err := tx.Transaction(func(tx *gorm.DB) error {
		// Both terms stay locked until commit: a concurrent delete of the replacement waits, so the rewritten
		// marks never point at a term that's gone. Locked in id order, so two crossing deletes can't deadlock.
		ids := []uint{id}
		if replaceWith != 0 {
			ids = append(ids, replaceWith)
		}
		var locked []entities.GlossaryTerm
		if err := tx.Clauses(forUpdate).Where("id IN ?", ids).Order("id").Find(&locked).Error; err != nil {
			return err
		}
		found := map[uint]bool{}
		for _, t := range locked {
			found[t.ID] = true
		}
		if !found[id] {
			return ErrTermNotFound
		}
		if replaceWith != 0 && !found[replaceWith] {
			return ErrTermReplace
		}
		// the steps are locked too: two deletes touching one step can't overwrite each other's rewrite
		steps, err := stepsWithTerms(tx.Clauses(forUpdate), id)
		if err != nil {
			return err
		}
		seen := map[uint]bool{}
		for _, s := range steps {
			if !seen[s.RecipeID] {
				seen[s.RecipeID] = true
				recipes++
			}
		}
		if recipes > 0 && replaceWith == 0 && !removeMarks {
			return ErrTermInUse
		}
		for _, s := range steps {
			hints := make(entities.StringList, len(s.Hints))
			for i, h := range s.Hints {
				hints[i] = ReplaceTermRefs(h, id, replaceWith)
			}
			if err := tx.Model(&entities.RecipeStep{}).
				Where("recipe_id = ? AND step_no = ?", s.RecipeID, s.StepNo).
				Updates(map[string]any{"text": ReplaceTermRefs(s.Text, id, replaceWith), "hints": hints}).Error; err != nil {
				return err
			}
		}
		return tx.Delete(&entities.GlossaryTerm{}, id).Error
	})
	return recipes, err
}

// RecipeTerms: the terms marked in these steps (text and hints), A–Z; missing ids are skipped.
func RecipeTerms(tx *gorm.DB, steps []entities.RecipeStep) ([]dtos.GlossaryTermDto, error) {
	texts := []string{}
	for _, s := range steps {
		texts = append(texts, s.Text)
		texts = append(texts, s.Hints...)
	}
	out := []dtos.GlossaryTermDto{}
	ids := TermRefIDs(texts...)
	if len(ids) == 0 {
		return out, nil
	}
	var terms []entities.GlossaryTerm
	if err := tx.Where("id IN ?", ids).Find(&terms).Error; err != nil {
		return nil, err
	}
	for _, t := range terms {
		out = append(out, dtos.GlossaryTermDto{ID: t.ID, Name: t.Name, Explanation: t.Explanation})
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Name) < strings.ToLower(out[j].Name) })
	return out, nil
}
