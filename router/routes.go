package router

import (
	"net/http"

	h "uszatakuchnia/handlers"
)

// Routes — całe API. Nowy endpoint = nowa linijka tutaj (vercel.json kieruje całe /api/* do
// jednej funkcji, więc nic więcej nie trzeba). Zapis (wszystko poza GET) wymaga właściciela —
// pilnuje tego resp.Begin w api/index.go, a test w router_test.go sprawdza każdą trasę z tej tabeli.
var Routes = []Route{
	// Przepisy
	{http.MethodGet, "/api/recipes/list", h.RecipeList}, // ?admin=true — także szkice (właściciel)
	{http.MethodPost, "/api/recipes", h.RecipeCreate},
	{http.MethodGet, "/api/recipes/{id}", h.RecipeGet}, // szkic: 404 dla nie-właściciela
	{http.MethodPut, "/api/recipes/{id}", h.RecipeUpdate},
	{http.MethodDelete, "/api/recipes/{id}", h.RecipeDelete},

	// Składniki
	{http.MethodGet, "/api/ingredients/list", h.IngredientList},
	{http.MethodPost, "/api/ingredients", h.IngredientCreate},
	{http.MethodPut, "/api/ingredients/{id}", h.IngredientUpdate},
	{http.MethodDelete, "/api/ingredients/{id}", h.IngredientDelete},
	{http.MethodGet, "/api/enums/ingredient", h.IngredientTypes},

	// Słowniki (panel → Klasyfikacja)
	{http.MethodGet, "/api/taxonomy/list", h.Dictionaries},
	{http.MethodPost, "/api/taxonomy/values", h.ValueCreate},
	{http.MethodPut, "/api/taxonomy/values/{id}", h.ValueUpdate},
	{http.MethodDelete, "/api/taxonomy/values/{id}", h.ValueDelete}, // ?moveTo=
	{http.MethodPost, "/api/taxonomy/units", h.UnitCreate},
	{http.MethodPut, "/api/taxonomy/units/{id}", h.UnitUpdate},
	{http.MethodDelete, "/api/taxonomy/units/{id}", h.UnitDelete}, // ?moveTo=
	{http.MethodPost, "/api/taxonomy/equipment", h.EquipmentCreate},
	{http.MethodPut, "/api/taxonomy/equipment/{id}", h.EquipmentUpdate},
	{http.MethodDelete, "/api/taxonomy/equipment/{id}", h.EquipmentDelete}, // ?replaceWith= | ?removeFromRecipes=true
	{http.MethodPost, "/api/taxonomy/challenge-formats", h.ChallengeFormatCreate},
	{http.MethodPut, "/api/taxonomy/challenge-formats/{id}", h.ChallengeFormatUpdate},
	{http.MethodDelete, "/api/taxonomy/challenge-formats/{id}", h.ChallengeFormatDelete}, // ?moveTo=

	// Wyzwania
	{http.MethodGet, "/api/challenges/list", h.ChallengeList},
	{http.MethodGet, "/api/challenges", h.ChallengeList}, // jak dawniej: GET bez id = lista
	{http.MethodPost, "/api/challenges", h.ChallengeCreate},
	{http.MethodGet, "/api/challenges/{id}", h.ChallengeGet},
	{http.MethodPut, "/api/challenges/{id}", h.ChallengeUpdate},
	{http.MethodDelete, "/api/challenges/{id}", h.ChallengeDelete},

	// Zdjęcia
	{http.MethodPost, "/api/photos", h.PhotoUpload},
	{http.MethodGet, "/api/photos/{id}", h.PhotoGet},
}
