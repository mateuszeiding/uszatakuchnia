package db

import (
	"fmt"
	"os"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/services"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func runSeed(gdb *gorm.DB) error {
	if err := seedIngredientTypes(gdb); err != nil {
		return err
	}
	if err := seedTaxonomy(gdb); err != nil {
		return err
	}
	if err := seedUnits(gdb); err != nil {
		return err
	}
	return SeedAll(gdb)
}

// --------------------
// Słowniki
// --------------------

func seedIngredientTypes(tx *gorm.DB) error {
	types := []entities.IngredientType{
		{Code: "VEGETABLE", Name: "warzywa i grzyby"},
		{Code: "FRUIT", Name: "owoce"},
		{Code: "MEAT", Name: "mięso i drób"},
		{Code: "FISH", Name: "ryby i owoce morza"},
		{Code: "DAIRY", Name: "nabiał i jaja"},
		{Code: "GRAIN", Name: "zboża i makarony"},
		{Code: "LEGUME", Name: "rośliny strączkowe"},
		{Code: "SPICE", Name: "przyprawy i zioła"},
		{Code: "FAT", Name: "oleje i tłuszcze"},
		{Code: "OTHER", Name: "inne"},
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&types).Error
}

type taxonomyValueSeed struct {
	Label string
	Color string // klucz tokenu, "" = brak
}

type taxonomyAxisSeed struct {
	Axis   entities.TaxonomyAxis
	Values []taxonomyValueSeed
}

func labels(ls ...string) []taxonomyValueSeed {
	out := make([]taxonomyValueSeed, 0, len(ls))
	for _, l := range ls {
		out = append(out, taxonomyValueSeed{Label: l})
	}
	return out
}

func seedTaxonomy(tx *gorm.DB) error {
	axes := []taxonomyAxisSeed{
		{
			Axis: entities.TaxonomyAxis{Code: "dish", Label: "Typ dania", Multi: true, Required: true, HasColor: true, FilterMode: "or", SortOrder: 1},
			Values: []taxonomyValueSeed{
				{"śniadanie", "coral"},
				{"zupa", "tomato"},
				{"danie główne", "terracotta"},
				{"sałatka", "lime"},
				{"przekąska", "plum"},
				{"deser", "pink"},
				{"pieczywo i wypieki", "brown"},
				{"sos / dodatek", "mustard"},
				{"napój", "teal"},
			},
		},
		{
			Axis:   entities.TaxonomyAxis{Code: "base", Label: "Baza", Multi: true, FilterMode: "or", SortOrder: 2},
			Values: labels("mięso", "drób", "ryby", "owoce morza", "jajka", "nabiał", "warzywa", "strączki", "grzyby", "makaron", "ryż i kasze"),
		},
		{
			Axis:   entities.TaxonomyAxis{Code: "cuisine", Label: "Kuchnia", Multi: false, FilterMode: "or", SortOrder: 3},
			Values: labels("polska", "włoska", "francuska", "hiszpańska", "grecka", "bliskowschodnia", "indyjska", "azjatycka", "meksykańska", "amerykańska", "skandynawska"),
		},
		{
			Axis:   entities.TaxonomyAxis{Code: "diet", Label: "Diety", Multi: true, FilterMode: "and", SortOrder: 4},
			Values: labels("wegetariańskie", "wegańskie", "bez glutenu", "bez laktozy", "bez cukru", "wysokobiałkowe", "niskokaloryczne", "keto", "low fodmap"),
		},
		{
			Axis:   entities.TaxonomyAxis{Code: "practical", Label: "Praktyczne", Multi: true, FilterMode: "or", SortOrder: 5},
			Values: labels("1 patelnia", "1 garnek", "do 5 składników", "można mrozić", "na imprezę", "dla dzieci", "meal prep"),
		},
	}

	for _, a := range axes {
		axis := a.Axis
		if err := tx.Where(entities.TaxonomyAxis{Code: axis.Code}).FirstOrCreate(&axis).Error; err != nil {
			return err
		}
		for i, v := range a.Values {
			val := entities.TaxonomyValue{AxisID: axis.ID, Label: v.Label, SortOrder: i + 1}
			if v.Color != "" {
				c := v.Color
				val.Color = &c
			}
			if err := tx.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "axis_id"}, {Name: "label"}},
				DoNothing: true,
			}).Create(&val).Error; err != nil {
				return err
			}
		}
	}
	return nil
}

func seedUnits(tx *gorm.DB) error {
	forms := func(one, few, many string) (*string, *string, *string) { return &one, &few, &many }
	unit := func(name string, f ...string) entities.Unit {
		u := entities.Unit{Name: name}
		if len(f) == 3 {
			u.FormOne, u.FormFew, u.FormMany = forms(f[0], f[1], f[2])
		}
		return u
	}

	units := []entities.Unit{
		unit("g"),
		unit("kg"),
		unit("ml"),
		unit("l"),
		unit("łyżka", "łyżka", "łyżki", "łyżek"),
		unit("łyżeczka", "łyżeczka", "łyżeczki", "łyżeczek"),
		unit("szklanka", "szklanka", "szklanki", "szklanek"),
		unit("szczypta", "szczypta", "szczypty", "szczypt"),
		unit("ząbek", "ząbek", "ząbki", "ząbków"),
		unit("sztuka", "sztuka", "sztuki", "sztuk"),
		unit("filet", "filet", "filety", "filetów"),
		unit("plaster", "plaster", "plastry", "plastrów"),
		unit("pęczek", "pęczek", "pęczki", "pęczków"),
		unit("garść", "garść", "garście", "garści"),
		unit("listek", "listek", "listki", "listków"),
		{Name: "do smaku", NoQty: true},
	}
	for i := range units {
		units[i].SortOrder = i + 1
	}
	return tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "name"}},
		DoNothing: true,
	}).Create(&units).Error
}

// --------------------
// Ingredients
// --------------------

type IngredientImageSeed struct {
	PostID string // Unsplash photo id
}

type IngredientSeed struct {
	Name       string
	TypeCode   string
	Allergens  []string // kody z entities.Allergens
	ParentName *string
	Image      *IngredientImageSeed

	Nutrition   map[string]entities.MacroRange // na 100 g, klucze z entities.NutrientKeys
	Conversions map[string]entities.MacroRange // nazwa jednostki → gramy
}

// rg — zakres do seeda.
func rg(min, max float64) entities.MacroRange {
	r := entities.MacroRange{Min: min}
	if max > min {
		r.Max = &max
	}
	return r
}

// one — pojedyncza wartość do seeda.
func one(v float64) entities.MacroRange { return entities.MacroRange{Min: v} }

// nut — makro na 100 g w kolejności: kcal, białko, tłuszcz, węglowodany, błonnik.
func nut(kcal, protein, fat, carbs, fiber entities.MacroRange) map[string]entities.MacroRange {
	return map[string]entities.MacroRange{"kcal": kcal, "protein": protein, "fat": fat, "carbs": carbs, "fiber": fiber}
}

// conv — przeliczniki: pary nazwa jednostki, gramy.
func conv(pairs ...any) map[string]entities.MacroRange {
	out := map[string]entities.MacroRange{}
	for i := 0; i+1 < len(pairs); i += 2 {
		out[pairs[i].(string)] = pairs[i+1].(entities.MacroRange)
	}
	return out
}

func loadIngredientTypeMap(tx *gorm.DB) (map[string]uint, error) {
	out := map[string]uint{}
	var rows []entities.IngredientType
	if err := tx.Find(&rows).Error; err != nil {
		return nil, err
	}
	for _, r := range rows {
		out[r.Code] = r.ID
	}
	return out, nil
}

func findRootIngredientID(tx *gorm.DB, name string) (*uint, error) {
	var ing entities.Ingredient
	if err := tx.Where("parent_id IS NULL AND name = ?", name).First(&ing).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("parent not found (root ingredient): %s", name)
		}
		return nil, err
	}
	return &ing.ID, nil
}

func SeedOneIngredient(tx *gorm.DB, seed IngredientSeed) (uint, error) {
	typeMap, err := loadIngredientTypeMap(tx)
	if err != nil {
		return 0, err
	}

	typeID, ok := typeMap[seed.TypeCode]
	if !ok {
		return 0, fmt.Errorf("unknown ingredient type code: %s", seed.TypeCode)
	}

	var parentID *uint
	if seed.ParentName != nil {
		pid, err := findRootIngredientID(tx, *seed.ParentName)
		if err != nil {
			return 0, err
		}
		parentID = pid
	}

	var ing entities.Ingredient
	if err := tx.Where(entities.Ingredient{Name: seed.Name, ParentID: parentID}).
		Assign(entities.Ingredient{TypeID: typeID}).
		FirstOrCreate(&ing).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create ingredient: %w", err)
	}

	if err := tx.Where("ingredient_id = ?", ing.ID).Delete(&entities.IngredientAllergen{}).Error; err != nil {
		return 0, err
	}
	for _, code := range seed.Allergens {
		if !entities.IsAllergenCode(code) {
			return 0, fmt.Errorf("unknown allergen %q for %s", code, seed.Name)
		}
		if err := tx.Create(&entities.IngredientAllergen{IngredientID: ing.ID, Code: code}).Error; err != nil {
			return 0, err
		}
	}

	for _, k := range entities.NutrientKeys {
		if r, ok := seed.Nutrition[k]; ok {
			ing.SetMacro(k, &r)
		} else {
			ing.SetMacro(k, nil)
		}
	}
	if err := tx.Omit("Type", "Parent", "Image", "Allergens", "Conversions").Save(&ing).Error; err != nil {
		return 0, err
	}
	if err := tx.Where("ingredient_id = ?", ing.ID).Delete(&entities.IngredientUnitConversion{}).Error; err != nil {
		return 0, err
	}
	for unitName, grams := range seed.Conversions {
		unitID, err := findUnitID(tx, unitName)
		if err != nil {
			return 0, err
		}
		c := entities.IngredientUnitConversion{IngredientID: ing.ID, UnitID: *unitID, GramsMin: grams.Min, GramsMax: grams.Max}
		if err := tx.Create(&c).Error; err != nil {
			return 0, err
		}
	}

	if seed.Image != nil && os.Getenv("UNSPLASH_ACCESS_KEY") != "" {
		if err := services.UpsertImageFromUnsplash(tx, ing.ID, entities.IngredientRef, seed.Image.PostID); err != nil {
			return 0, err
		}
	}

	return ing.ID, nil
}

// --------------------
// Recipes
// --------------------

type RecipePhotoSeed struct {
	URL string
}

type RecipeSeed struct {
	Name        string
	Servings    int
	Tagline     *string
	Description *string
	TimeMin     int
	TimeMax     *int
	Difficulty  *int
	Status      string // "draft" | "published"
	Photo       *RecipePhotoSeed

	// Klasyfikacja: kod osi → etykiety wartości
	Taxonomy map[string][]string

	Steps []RecipeStepSeed
	Items []RecipeIngredientSeed
}

type RecipeStepSeed struct {
	Section  *string
	Text     string
	Prep     bool
	PrepLead *string
}

type RecipeIngredientSeed struct {
	Section *string

	IngredientName       string
	IngredientParentName *string

	Quantity string // "2", "½", "2–3"; "" = brak
	Unit     string // nazwa jednostki ze słownika; "" = brak
	Note     *string
}

func findIngredientID(tx *gorm.DB, name string, parentName *string) (uint, error) {
	var parentID *uint
	if parentName != nil {
		pid, err := findRootIngredientID(tx, *parentName)
		if err != nil {
			return 0, err
		}
		parentID = pid
	}

	var ing entities.Ingredient
	q := tx.Where("name = ?", name)
	if parentID == nil {
		q = q.Where("parent_id IS NULL")
	} else {
		q = q.Where("parent_id = ?", *parentID)
	}

	if err := q.First(&ing).Error; err != nil {
		return 0, fmt.Errorf("ingredient not found: %s (parent=%v): %w", name, parentName, err)
	}

	return ing.ID, nil
}

func findUnitID(tx *gorm.DB, name string) (*uint, error) {
	if name == "" {
		return nil, nil
	}
	var u entities.Unit
	if err := tx.Where("name = ?", name).First(&u).Error; err != nil {
		return nil, fmt.Errorf("unit not found: %s: %w", name, err)
	}
	return &u.ID, nil
}

func findTaxonomyValueID(tx *gorm.DB, axisCode, label string) (uint, error) {
	var v entities.TaxonomyValue
	if err := tx.Joins("Axis").
		Where(`"Axis".code = ? AND taxonomy_values.label = ?`, axisCode, label).
		First(&v).Error; err != nil {
		return 0, fmt.Errorf("taxonomy value not found: %s/%s: %w", axisCode, label, err)
	}
	return v.ID, nil
}

func upsertRecipePhoto(tx *gorm.DB, recipeID uint, url string) error {
	photo := entities.RecipePhoto{RecipeID: recipeID, ImageURL: url}
	return tx.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "recipe_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"image_url": photo.ImageURL,
		}),
	}).Create(&photo).Error
}

func SeedOneRecipe(tx *gorm.DB, seed RecipeSeed) (uint, error) {
	if seed.Servings <= 0 {
		return 0, fmt.Errorf("recipe servings must be > 0: %s", seed.Name)
	}

	status := seed.Status
	if status == "" {
		status = entities.StatusPublished
	}

	var r entities.Recipe
	if err := tx.Where(entities.Recipe{Name: seed.Name}).
		Assign(entities.Recipe{
			Servings:    seed.Servings,
			Description: seed.Description,
			Tagline:     seed.Tagline,
			TimeMin:     seed.TimeMin,
			TimeMax:     seed.TimeMax,
			Difficulty:  seed.Difficulty,
			Status:      status,
		}).
		FirstOrCreate(&r).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create recipe: %w", err)
	}

	if seed.Photo != nil && seed.Photo.URL != "" {
		if err := upsertRecipePhoto(tx, r.ID, seed.Photo.URL); err != nil {
			return 0, err
		}
	}

	// Kroki, składniki i klasyfikacja — usuń stare, wstaw nowe
	for _, m := range []any{&entities.RecipeStep{}, &entities.RecipeIngredient{}, &entities.RecipeTaxonomy{}} {
		if err := tx.Where("recipe_id = ?", r.ID).Delete(m).Error; err != nil {
			return 0, err
		}
	}

	for i, s := range seed.Steps {
		step := entities.RecipeStep{
			RecipeID: r.ID, StepNo: i + 1, Section: s.Section, Text: s.Text, Prep: s.Prep, PrepLead: s.PrepLead,
		}
		if err := tx.Create(&step).Error; err != nil {
			return 0, err
		}
	}

	for i, it := range seed.Items {
		ingID, err := findIngredientID(tx, it.IngredientName, it.IngredientParentName)
		if err != nil {
			return 0, err
		}
		unitID, err := findUnitID(tx, it.Unit)
		if err != nil {
			return 0, err
		}
		var qty *string
		if it.Quantity != "" {
			if _, ok := services.ParseQuantity(it.Quantity); !ok {
				return 0, fmt.Errorf("invalid quantity %q in %s", it.Quantity, seed.Name)
			}
			qty = strPtr(it.Quantity)
		}

		item := entities.RecipeIngredient{
			RecipeID:     r.ID,
			SortOrder:    i + 1,
			Section:      it.Section,
			IngredientID: ingID,
			Quantity:     qty,
			UnitID:       unitID,
			Note:         it.Note,
		}
		if err := tx.Create(&item).Error; err != nil {
			return 0, err
		}
	}

	for axisCode, ls := range seed.Taxonomy {
		for _, l := range ls {
			vid, err := findTaxonomyValueID(tx, axisCode, l)
			if err != nil {
				return 0, err
			}
			if err := tx.Create(&entities.RecipeTaxonomy{RecipeID: r.ID, ValueID: vid}).Error; err != nil {
				return 0, err
			}
		}
	}

	return r.ID, nil
}

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// ── Składniki ────────────────────────────────────────────
		ings := []IngredientSeed{
			// Warzywa i grzyby
			{Name: "Cebula", TypeCode: "VEGETABLE", Nutrition: nut(one(40), one(1.1), one(0.1), one(9.3), one(1.7)), Conversions: conv("sztuka", rg(100, 150)), Image: &IngredientImageSeed{PostID: "bC1fXU1v98U"}},
			{Name: "Czosnek", TypeCode: "VEGETABLE", Nutrition: nut(one(149), one(6.4), one(0.5), one(33), one(2.1)), Conversions: conv("ząbek", rg(4, 6)), Image: &IngredientImageSeed{PostID: "_b8n1KTZ8rw"}},
			{Name: "Pomidor", TypeCode: "VEGETABLE", Nutrition: nut(rg(18, 24), rg(0.9, 1.2), one(0.2), rg(3.5, 4.5), rg(1, 1.2)), Conversions: conv("sztuka", rg(100, 150)), Image: &IngredientImageSeed{PostID: "Jy7e7RjOFVo"}},
			{Name: "Papryczka chili", TypeCode: "VEGETABLE", Nutrition: nut(one(40), one(1.9), one(0.4), one(8.8), one(1.5)), Conversions: conv("sztuka", rg(10, 20)), Image: &IngredientImageSeed{PostID: "iBsmi-wCXNE"}},
			{Name: "Marchew", TypeCode: "VEGETABLE", Nutrition: nut(one(41), one(0.9), one(0.2), one(9.6), one(2.8)), Conversions: conv("sztuka", rg(60, 100)), Image: &IngredientImageSeed{PostID: "0bPFoBML0qs"}},
			{Name: "Ziemniaki", TypeCode: "VEGETABLE", Nutrition: nut(one(77), one(2), one(0.1), one(17), one(2.2)), Conversions: conv("sztuka", rg(100, 200)), Image: &IngredientImageSeed{PostID: "rB7-iPJADpI"}},
			{Name: "Szpinak", TypeCode: "VEGETABLE", Nutrition: nut(one(23), one(2.9), one(0.4), one(3.6), one(2.2)), Conversions: conv("garść", rg(25, 35)), Image: &IngredientImageSeed{PostID: "GI39ElkDQFo"}},
			{Name: "Cukinia", TypeCode: "VEGETABLE", Nutrition: nut(one(17), one(1.2), one(0.3), one(3.1), one(1)), Conversions: conv("sztuka", rg(200, 300)), Image: &IngredientImageSeed{PostID: "SbzMoWL0xvw"}},
			{Name: "Papryka czerwona", TypeCode: "VEGETABLE", Nutrition: nut(one(31), one(1), one(0.3), one(6), one(2.1)), Conversions: conv("sztuka", rg(150, 200)), Image: &IngredientImageSeed{PostID: "iKKOhzHXR9E"}},
			// Owoce
			{Name: "Cytryna", TypeCode: "FRUIT", Nutrition: nut(one(29), one(1.1), one(0.3), one(9.3), one(2.8)), Conversions: conv("sztuka", rg(100, 130)), Image: &IngredientImageSeed{PostID: "7WAGthfGJ9w"}},
			{Name: "Limonka", TypeCode: "FRUIT", Nutrition: nut(one(30), one(0.7), one(0.2), one(10.5), one(2.8)), Conversions: conv("sztuka", rg(50, 70))},
			// Mięso i drób
			{Name: "Kurczak", TypeCode: "MEAT", Nutrition: nut(rg(110, 180), rg(18, 23), rg(1.5, 10), one(0), one(0)), Conversions: conv("filet", rg(150, 250)), Image: &IngredientImageSeed{PostID: "pNcFMdEe09Q"}},
			{Name: "Wołowina", TypeCode: "MEAT", Nutrition: nut(rg(150, 250), rg(18, 21), rg(7, 18), one(0), one(0)), Image: &IngredientImageSeed{PostID: "u2QnkcLNsBY"}},
			{Name: "Boczek", TypeCode: "MEAT", Nutrition: nut(rg(400, 540), rg(12, 16), rg(37, 50), rg(0, 1.5), one(0)), Conversions: conv("plaster", rg(15, 25))},
			// Ryby i owoce morza
			{Name: "Łosoś", TypeCode: "FISH", Allergens: []string{"fish"}, Nutrition: nut(rg(180, 210), rg(19, 20.5), rg(11, 14), one(0), one(0)), Conversions: conv("filet", rg(120, 180))},
			{Name: "Dorsz", TypeCode: "FISH", Allergens: []string{"fish"}, Nutrition: nut(rg(78, 82), rg(17.5, 18.5), rg(0.6, 0.8), one(0), one(0)), Conversions: conv("filet", rg(120, 180))},
			// Nabiał i jaja
			{Name: "Masło", TypeCode: "DAIRY", Allergens: []string{"milk"}, Nutrition: nut(rg(735, 745), one(0.7), rg(81, 83), one(0.6), one(0)), Conversions: conv("łyżka", one(15), "łyżeczka", one(5))},
			{Name: "Parmezan", TypeCode: "DAIRY", Allergens: []string{"milk"}, Nutrition: nut(rg(390, 431), rg(32, 38), rg(26, 29), rg(0, 4), one(0)), Conversions: conv("łyżka", rg(5, 7))},
			{Name: "Jajka", TypeCode: "DAIRY", Allergens: []string{"eggs"}, Nutrition: nut(rg(140, 143), rg(12.5, 13), rg(9.5, 10), one(0.7), one(0)), Conversions: conv("sztuka", rg(50, 65))},
			{Name: "Śmietana 30%", TypeCode: "DAIRY", Allergens: []string{"milk"}, Nutrition: nut(rg(290, 300), rg(2.2, 2.5), one(30), rg(3, 3.5), one(0)), Conversions: conv("łyżka", rg(12, 15))},
			// Zboża i makarony
			{Name: "Makaron spaghetti", TypeCode: "GRAIN", Allergens: []string{"gluten"}, Nutrition: nut(rg(350, 360), rg(12, 13), rg(1.5, 2), rg(70, 72), rg(3, 3.5))},
			// Oleje i tłuszcze
			{Name: "Oliwa z oliwek", TypeCode: "FAT", Nutrition: nut(one(884), one(0), one(100), one(0), one(0)), Conversions: conv("łyżka", one(10), "łyżeczka", rg(4, 5), "ml", one(0.92))},
			// Inne
			{Name: "Bulion warzywny", TypeCode: "OTHER", Allergens: []string{"celery"}, Nutrition: nut(rg(5, 15), rg(0.2, 0.5), rg(0, 0.5), rg(0.5, 2), one(0))},
			// Przyprawy i zioła
			{Name: "Bazylia", TypeCode: "SPICE", Nutrition: nut(one(23), one(3.2), one(0.6), one(2.7), one(1.6)), Conversions: conv("garść", rg(5, 10), "listek", rg(0.3, 0.6)), Image: &IngredientImageSeed{PostID: "DDKdfRe1GxA"}},
			{Name: "Pietruszka (nać)", TypeCode: "SPICE", Nutrition: nut(one(36), one(3), one(0.8), one(6.3), one(3.3)), Conversions: conv("garść", rg(10, 15), "pęczek", rg(30, 50)), Image: &IngredientImageSeed{PostID: "r6BUzN_jTHg"}},
			{Name: "Tymianek", TypeCode: "SPICE", Nutrition: nut(one(101), one(5.6), one(1.7), one(24), one(14)), Conversions: conv("pęczek", rg(20, 30))},
			{Name: "Rozmaryn", TypeCode: "SPICE", Nutrition: nut(one(131), one(3.3), one(5.9), one(21), one(14))},
			{Name: "Oregano", TypeCode: "SPICE"},
			{Name: "Sól", TypeCode: "SPICE"},
			{Name: "Pieprz czarny", TypeCode: "SPICE", Nutrition: nut(one(251), one(10), one(3.3), one(64), one(25)), Conversions: conv("szczypta", rg(0.3, 0.5))},
			{Name: "Kurkuma", TypeCode: "SPICE", Nutrition: nut(one(312), one(9.7), one(3.3), one(67), one(22.7)), Conversions: conv("łyżeczka", rg(2.5, 3))},
			{Name: "Kmin rzymski", TypeCode: "SPICE", Nutrition: nut(one(375), one(17.8), one(22), one(44), one(10.5)), Conversions: conv("łyżeczka", rg(2, 2.5))},
			{Name: "Papryka słodka", TypeCode: "SPICE", Nutrition: nut(one(282), one(14), one(13), one(54), one(35)), Conversions: conv("łyżeczka", rg(2.3, 3))},
			{Name: "Papryka wędzona", TypeCode: "SPICE", Nutrition: nut(one(282), one(14), one(13), one(54), one(35)), Conversions: conv("łyżeczka", rg(2.3, 3))},
		}

		for _, s := range ings {
			if _, err := SeedOneIngredient(tx, s); err != nil {
				return err
			}
		}

		// ── Przepisy ─────────────────────────────────────────────

		secSos := strPtr("Sos")
		secMakaron := strPtr("Makaron")
		secDressing := strPtr("Dressing")
		secKurczak := strPtr("Kurczak")

		recipes := []RecipeSeed{
			{
				Name:       "Spaghetti aglio e olio",
				Servings:   2,
				Tagline:    strPtr("Trzy składniki, gotowe w 15 minut. Klasyk za który nie wypada przepraszać."),
				TimeMin:    15,
				TimeMax:    intPtr(20),
				Difficulty: intPtr(1),
				Status:     "published",
				Taxonomy: map[string][]string{
					"dish":    {"danie główne"},
					"base":    {"makaron"},
					"cuisine": {"włoska"},
					"diet":    {"wegetariańskie"},
				},
				Steps: []RecipeStepSeed{
					{Section: secMakaron, Text: "Ugotuj makaron al dente w mocno osolonej wodzie. Zachowaj szklankę wody z gotowania."},
					{Section: secSos, Text: "Pokrój czosnek w cienkie plasterki. Na dużej patelni rozgrzej oliwę na średnim ogniu — czosnek ma powoli się złocić, nie smażyć. Dodaj chili."},
					{Section: secSos, Text: "Przełóż makaron na patelnię, dolej kilka łyżek wody z gotowania. Mieszaj aż sos będzie jedwabisty i połączy się z makaronem."},
					{Text: "Posyp pietruszką, opcjonalnie startym parmezanem. Podawaj od razu."},
				},
				Items: []RecipeIngredientSeed{
					{Section: secMakaron, IngredientName: "Makaron spaghetti", Quantity: "200", Unit: "g"},
					{Section: secSos, IngredientName: "Czosnek", Quantity: "4", Unit: "ząbek"},
					{Section: secSos, IngredientName: "Oliwa z oliwek", Quantity: "4", Unit: "łyżka"},
					{Section: secSos, IngredientName: "Papryczka chili", Quantity: "1", Unit: "sztuka", Note: strPtr("lub szczypta płatków chili")},
					{IngredientName: "Pietruszka (nać)", Quantity: "1", Unit: "garść"},
					{IngredientName: "Parmezan", Quantity: "30", Unit: "g", Note: strPtr("opcjonalnie")},
				},
			},
			{
				Name:       "Pieczony łosoś z cytryną i tymiankiem",
				Servings:   2,
				Tagline:    strPtr("Maślany, soczysty łosoś gotowy w 20 minut. Jedna blacha, zero bałaganu."),
				TimeMin:    25,
				Difficulty: intPtr(1),
				Status:     "published",
				Taxonomy: map[string][]string{
					"dish": {"danie główne"},
					"base": {"ryby"},
					"diet": {"bez glutenu"},
				},
				Steps: []RecipeStepSeed{
					{Text: "Rozgrzej piekarnik do 200°C. Blachę wyłóż papierem do pieczenia."},
					{Text: "Filety połóż skórą w dół. Posmaruj roztopionym masłem, posól i popieprz. Na wierzch połóż gałązki tymianku i plasterki cytryny."},
					{Text: "Piecz 15–18 minut — łosoś jest gotowy gdy mięso łatwo rozdziela się na płaty. Jeśli filety są cienkie, skróć do 12 minut."},
					{Text: "Skrop sokiem z reszty cytryny i podawaj od razu."},
				},
				Items: []RecipeIngredientSeed{
					{IngredientName: "Łosoś", Quantity: "2", Unit: "filet", Note: strPtr("ok. 150 g każdy")},
					{IngredientName: "Masło", Quantity: "2", Unit: "łyżka"},
					{IngredientName: "Cytryna", Quantity: "1", Unit: "sztuka"},
					{IngredientName: "Tymianek", Quantity: "½", Unit: "pęczek", Note: strPtr("świeży")},
					{IngredientName: "Sól", Unit: "do smaku"},
					{IngredientName: "Pieprz czarny", Unit: "do smaku"},
				},
			},
			{
				Name:       "Kurczak tikka masala",
				Servings:   4,
				Tagline:    strPtr("Kremowy, aromatyczny sos z kurczakiem. Nie potrzebujesz żadnego tandooru."),
				TimeMin:    50,
				Difficulty: intPtr(2),
				Status:     "published",
				Taxonomy: map[string][]string{
					"dish":    {"danie główne"},
					"base":    {"drób"},
					"cuisine": {"indyjska"},
					"diet":    {"bez glutenu"},
				},
				Steps: []RecipeStepSeed{
					{Section: secKurczak, Prep: true, PrepLead: strPtr("12 h"), Text: "Kurczaka pokrój w kawałki. Wymieszaj z kurkumą, kminem, papryką wędzoną i oliwą. Marynuj minimum 30 minut, najlepiej całą noc."},
					{Section: secKurczak, Text: "Obsmaż kurczaka na mocno rozgrzanej patelni z każdej strony — partiami, żeby się nie dusił. Odłóż."},
					{Section: secSos, Text: "Na tej samej patelni podsmaż cebulę do zezłocenia. Dodaj czosnek i chili, smaż 1 minutę."},
					{Section: secSos, Text: "Wrzuć pomidory, dolej szklankę bulionu. Gotuj na małym ogniu 15 minut aż sos zgęstnieje."},
					{Section: secSos, Text: "Dodaj śmietanę i kurczaka z powrotem. Gotuj 5 minut. Dopraw do smaku."},
					{Text: "Posyp kolendrą lub pietruszką. Podawaj z ryżem lub chlebem naan."},
				},
				Items: []RecipeIngredientSeed{
					{Section: secKurczak, IngredientName: "Kurczak", Quantity: "600", Unit: "g", Note: strPtr("piersi lub udka bez kości")},
					{Section: secKurczak, IngredientName: "Kurkuma", Quantity: "1", Unit: "łyżeczka"},
					{Section: secKurczak, IngredientName: "Kmin rzymski", Quantity: "1", Unit: "łyżeczka"},
					{Section: secKurczak, IngredientName: "Papryka wędzona", Quantity: "1", Unit: "łyżeczka"},
					{Section: secKurczak, IngredientName: "Oliwa z oliwek", Quantity: "2", Unit: "łyżka"},
					{Section: secSos, IngredientName: "Cebula", Quantity: "1", Unit: "sztuka"},
					{Section: secSos, IngredientName: "Czosnek", Quantity: "3–4", Unit: "ząbek"},
					{Section: secSos, IngredientName: "Papryczka chili", Quantity: "1", Unit: "sztuka"},
					{Section: secSos, IngredientName: "Pomidor", Quantity: "400", Unit: "g", Note: strPtr("z puszki lub świeże")},
					{Section: secSos, IngredientName: "Bulion warzywny", Quantity: "250", Unit: "ml"},
					{Section: secSos, IngredientName: "Śmietana 30%", Quantity: "150", Unit: "ml"},
					{Section: secSos, IngredientName: "Papryka słodka", Quantity: "2", Unit: "łyżeczka"},
				},
			},
			{
				Name:       "Sałatka z pieczonym boczkiem i szpinakiem",
				Servings:   2,
				Tagline:    strPtr("Ciepły dressing, chrupiący boczek i świeży szpinak. Gotowe w 15 minut."),
				TimeMin:    15,
				Difficulty: intPtr(1),
				Status:     "published",
				Taxonomy: map[string][]string{
					"dish":      {"sałatka"},
					"base":      {"mięso", "warzywa"},
					"diet":      {"bez glutenu"},
					"practical": {"1 patelnia", "do 5 składników"},
				},
				Steps: []RecipeStepSeed{
					{Section: secDressing, Text: "Na patelni usmaż boczek na chrupiąco. Zostaw 2 łyżki tłuszczu na patelni."},
					{Section: secDressing, Text: "Do tłuszczu dodaj czosnek, smaż 30 sekund. Dodaj sok z cytryny i szczyptę pieprzu — dressing gotowy."},
					{Text: "Szpinak ułóż w misce. Polej ciepłym dressingiem — liście lekko zwiotczeją. Dodaj boczek i od razu podawaj."},
				},
				Items: []RecipeIngredientSeed{
					{IngredientName: "Szpinak", Quantity: "150", Unit: "g", Note: strPtr("baby spinach")},
					{IngredientName: "Boczek", Quantity: "100", Unit: "g"},
					{Section: secDressing, IngredientName: "Czosnek", Quantity: "2", Unit: "ząbek"},
					{Section: secDressing, IngredientName: "Cytryna", Quantity: "½", Unit: "sztuka", Note: strPtr("sok")},
					{Section: secDressing, IngredientName: "Pieprz czarny", Quantity: "1", Unit: "szczypta"},
				},
			},
		}

		for _, r := range recipes {
			if _, err := SeedOneRecipe(tx, r); err != nil {
				return err
			}
		}

		// Makro przepisów liczone ze składników (recipes.nutrition).
		if err := services.RecomputeRecipes(tx); err != nil {
			return err
		}

		return nil
	})
}

func strPtr(s string) *string { return &s }
func intPtr(i int) *int       { return &i }
