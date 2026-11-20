package mappers

import (
	"uszatakuchnia/db/entities"
	"uszatakuchnia/dtos"
)

// MapRecipeEquipmentToDto needs Equipment.Equipment preloaded (ordered by sort_order).
func MapRecipeEquipmentToDto(a []entities.RecipeEquipment) []dtos.RecipeEquipmentDto {
	out := make([]dtos.RecipeEquipmentDto, 0, len(a))
	for _, e := range a {
		out = append(out, dtos.RecipeEquipmentDto{
			EquipmentID: e.EquipmentID, Name: e.Equipment.Name, Spec: e.Spec, Optional: e.Optional,
		})
	}
	return out
}

func MapEquipmentToDto(e entities.Equipment) dtos.EquipmentDto {
	return dtos.EquipmentDto{ID: e.ID, Name: e.Name}
}
