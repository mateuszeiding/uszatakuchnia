package router

import (
	"net/http"

	h "uszatakuchnia/handlers"
)

// Routes is the whole API (+ link-preview bot pages at the end). A new endpoint = a new line here (vercel.json
// sends all of /api/* to one function, so nothing else is needed). Writes (anything but GET) require the owner:
// resp.Begin in api/index.go enforces it, and router_test.go checks every route in this table.
var Routes = []Route{
	// Recipes
	{http.MethodGet, "/api/recipes/list", h.RecipeList}, // ?admin=true also returns drafts (owner)
	{http.MethodPost, "/api/recipes", h.RecipeCreate},
	{http.MethodGet, "/api/recipes/{id}", h.RecipeGet}, // draft: 404 for non-owners
	{http.MethodPut, "/api/recipes/{id}", h.RecipeUpdate},
	{http.MethodDelete, "/api/recipes/{id}", h.RecipeDelete},

	// Ingredients
	{http.MethodGet, "/api/ingredients/list", h.IngredientList},
	{http.MethodPost, "/api/ingredients", h.IngredientCreate},
	{http.MethodPut, "/api/ingredients/{id}", h.IngredientUpdate},
	{http.MethodDelete, "/api/ingredients/{id}", h.IngredientDelete},
	{http.MethodGet, "/api/enums/ingredient", h.IngredientTypes},

	// Dictionaries (admin → Classification)
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

	// Challenges
	{http.MethodGet, "/api/challenges/list", h.ChallengeList},
	{http.MethodGet, "/api/challenges", h.ChallengeList}, // GET without an id = list
	{http.MethodPost, "/api/challenges", h.ChallengeCreate},
	{http.MethodGet, "/api/challenges/{id}", h.ChallengeGet},
	{http.MethodPut, "/api/challenges/{id}", h.ChallengeUpdate},
	{http.MethodDelete, "/api/challenges/{id}", h.ChallengeDelete},

	// Glossary (terms marked in steps and hints)
	{http.MethodGet, "/api/glossary", h.GlossaryList},
	{http.MethodPost, "/api/glossary", h.GlossaryCreate},
	{http.MethodPut, "/api/glossary/{id}", h.GlossaryUpdate},
	{http.MethodGet, "/api/glossary/{id}/usage", h.GlossaryUsage}, // owner
	{http.MethodDelete, "/api/glossary/{id}", h.GlossaryDelete}, // ?replaceWith= | ?removeMarks=true

	// Photos
	{http.MethodPost, "/api/photos", h.PhotoUpload},
	{http.MethodGet, "/api/photos/{id}", h.PhotoGet},

	// Pages for link-preview bots (Messenger, WhatsApp…): index.html with the recipe meta.
	// The only routes outside /api/; vercel.json sends only their user agents here (tested in router_test.go).
	{http.MethodGet, "/recipes/{id}", h.RecipeSharePage},
	{http.MethodGet, "/recipes/{id}/zakupy", h.ShoppingSharePage},
}
