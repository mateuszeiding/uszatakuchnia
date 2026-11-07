package migrations

import (
	"uszatakuchnia/services"

	"gorm.io/gorm"
)

// 0003: przelicza cache recipes.nutrition — nowe pole incomplete („bez: …”) i przeliczniki l ↔ ml.
var recomputeNutrition = Migration{
	ID: "0003_recompute_nutrition",
	Up: func(tx *gorm.DB) error {
		return services.RecomputeRecipes(tx) // świeża baza: brak przepisów, nic do zrobienia
	},
}
