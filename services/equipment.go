package services

import (
	"errors"
	"fmt"
	"strings"
	"unicode/utf8"

	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"

	"gorm.io/gorm"
)

// MaxEquipmentSpecLen is the longest equipment detail (characters), matching varchar(160).
const MaxEquipmentSpecLen = 160

// MaxEquipmentNameLen is the longest equipment name in the dictionary (characters), matching varchar(120).
const MaxEquipmentNameLen = 120

func validateEquipment(tx *gorm.DB, req dtos.UpsertRecipeRequest, add func(field, msg string)) error {
	if req.Equipment == nil {
		return nil
	}
	ids := []uint{}
	for _, e := range *req.Equipment {
		ids = append(ids, e.EquipmentId)
	}
	known := map[uint]bool{}
	if len(ids) > 0 {
		var found []uint
		if err := tx.Model(&entities.Equipment{}).Where("id IN ?", uniqueUints(ids)).Pluck("id", &found).Error; err != nil {
			return err
		}
		for _, id := range found {
			known[id] = true
		}
	}
	for i, e := range *req.Equipment {
		f := fmt.Sprintf("equipment[%d]", i)
		if !known[e.EquipmentId] {
			add(f+".equipmentId", "Wybierz sprzęt ze słownika")
		}
		if e.Spec != nil && utf8.RuneCountInString(strings.TrimSpace(*e.Spec)) > MaxEquipmentSpecLen {
			add(f+".spec", fmt.Sprintf("Doprecyzowanie może mieć najwyżej %d znaków", MaxEquipmentSpecLen))
		}
	}
	return nil
}

// replaceRecipeEquipment replaces the recipe's equipment (nil in the request = unchanged).
func replaceRecipeEquipment(tx *gorm.DB, recipeID uint, req dtos.UpsertRecipeRequest) error {
	if req.Equipment == nil {
		return nil
	}
	if err := tx.Where("recipe_id = ?", recipeID).Delete(&entities.RecipeEquipment{}).Error; err != nil {
		return err
	}
	for i, e := range *req.Equipment {
		var spec *string
		if e.Spec != nil && strings.TrimSpace(*e.Spec) != "" {
			s := strings.TrimSpace(*e.Spec)
			spec = &s
		}
		row := entities.RecipeEquipment{
			RecipeID: recipeID, SortOrder: i + 1, EquipmentID: e.EquipmentId, Spec: spec, Optional: e.Optional,
		}
		if err := tx.Omit("Equipment").Create(&row).Error; err != nil {
			return err
		}
	}
	return nil
}

// ─── Dictionary (admin → Classification → Equipment) ─────────────────

var (
	ErrEquipmentName     = errors.New("name is required")
	ErrEquipmentTaken    = errors.New("equipment already exists")
	ErrEquipmentNotFound = errors.New("equipment not found")
	ErrEquipmentInUse    = errors.New("equipment is used in recipes; pass replaceWith or removeFromRecipes")
	ErrEquipmentReplace  = errors.New("replaceWith must be another equipment")
)

func cleanEquipmentName(name string) (string, error) {
	name = strings.Join(strings.Fields(name), " ")
	if name == "" {
		return "", ErrEquipmentName
	}
	if utf8.RuneCountInString(name) > MaxEquipmentNameLen {
		return "", fmt.Errorf("name is longer than %d characters", MaxEquipmentNameLen)
	}
	return name, nil
}

func equipmentNameTaken(tx *gorm.DB, name string, exceptID uint) (bool, error) {
	var n int64
	err := tx.Model(&entities.Equipment{}).Where("lower(name) = lower(?) AND id <> ?", name, exceptID).Count(&n).Error
	return n > 0, err
}

// SaveEquipment creates (id 0) or renames a dictionary item.
func SaveEquipment(tx *gorm.DB, id uint, rawName string) (entities.Equipment, error) {
	name, err := cleanEquipmentName(rawName)
	if err != nil {
		return entities.Equipment{}, err
	}
	var e entities.Equipment
	if id != 0 {
		if err := tx.First(&e, id).Error; err != nil {
			return e, ErrEquipmentNotFound
		}
	}
	if taken, err := equipmentNameTaken(tx, name, e.ID); err != nil {
		return e, err
	} else if taken {
		return e, ErrEquipmentTaken
	}
	e.Name = name
	return e, tx.Save(&e).Error
}

func EquipmentUsage(tx *gorm.DB) (map[uint]int, error) {
	var rows []struct {
		ID uint
		N  int
	}
	if err := tx.Model(&entities.RecipeEquipment{}).
		Select("equipment_id AS id, COUNT(DISTINCT recipe_id) AS n").
		Group("equipment_id").Scan(&rows).Error; err != nil {
		return nil, err
	}
	out := map[uint]int{}
	for _, r := range rows {
		out[r.ID] = r.N
	}
	return out, nil
}

// DeleteEquipment deletes a dictionary item. If it's in use: replaceWith (≠ 0) moves recipes to other equipment
// (details stay), removeFromRecipes removes it from recipes. Returns the number of recipes affected.
func DeleteEquipment(tx *gorm.DB, id, replaceWith uint, removeFromRecipes bool) (int64, error) {
	var e entities.Equipment
	if err := tx.First(&e, id).Error; err != nil {
		return 0, ErrEquipmentNotFound
	}
	var used int64
	if err := tx.Model(&entities.RecipeEquipment{}).Where("equipment_id = ?", id).
		Distinct("recipe_id").Count(&used).Error; err != nil {
		return 0, err
	}
	if replaceWith != 0 {
		var target entities.Equipment
		if replaceWith == id || tx.First(&target, replaceWith).Error != nil {
			return 0, ErrEquipmentReplace
		}
	}
	if used > 0 && replaceWith == 0 && !removeFromRecipes {
		return used, ErrEquipmentInUse
	}

	err := tx.Transaction(func(tx *gorm.DB) error {
		if used > 0 {
			if replaceWith != 0 {
				if err := tx.Model(&entities.RecipeEquipment{}).Where("equipment_id = ?", id).
					Update("equipment_id", replaceWith).Error; err != nil {
					return err
				}
			} else {
				recipeIDs := []uint{}
				if err := tx.Model(&entities.RecipeEquipment{}).Where("equipment_id = ?", id).
					Distinct().Pluck("recipe_id", &recipeIDs).Error; err != nil {
					return err
				}
				if err := tx.Where("equipment_id = ?", id).Delete(&entities.RecipeEquipment{}).Error; err != nil {
					return err
				}
				if err := renumberRecipeEquipment(tx, recipeIDs); err != nil {
					return err
				}
			}
		}
		return tx.Delete(&entities.Equipment{}, id).Error
	})
	return used, err
}

// renumberRecipeEquipment: after a delete, the order has no gaps (1, 2, 3…).
func renumberRecipeEquipment(tx *gorm.DB, recipeIDs []uint) error {
	if len(recipeIDs) == 0 {
		return nil
	}
	// two steps, so we don't collide with the (recipe_id, sort_order) key
	if err := tx.Exec(`UPDATE recipe_equipment SET sort_order = -sort_order WHERE recipe_id IN ?`, recipeIDs).Error; err != nil {
		return err
	}
	return tx.Exec(`
		UPDATE recipe_equipment re SET sort_order = n.rn
		FROM (SELECT recipe_id, sort_order, ROW_NUMBER() OVER (PARTITION BY recipe_id ORDER BY sort_order DESC) AS rn
		      FROM recipe_equipment WHERE recipe_id IN ?) n
		WHERE re.recipe_id = n.recipe_id AND re.sort_order = n.sort_order`, recipeIDs).Error
}
