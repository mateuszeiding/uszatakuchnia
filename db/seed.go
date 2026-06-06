package db

import (
	"fmt"
	"uszatakuchnia/db/entities"
	"uszatakuchnia/services"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func runSeed(gdb *gorm.DB) error {
	ingredientTypes := []entities.IngredientType{
		{Code: "VEGETABLE", Name: "Warzywo"},
		{Code: "FRUIT", Name: "Owoc"},
		{Code: "HERB", Name: "Zioło"},
		{Code: "SPICE", Name: "Przyprawa"},
		{Code: "MEAT", Name: "Mięso"},
		{Code: "FISH", Name: "Ryba"},
		{Code: "OTHER", Name: "Inne"},
	}

	if err := gdb.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "code"}},
		DoNothing: true,
	}).Create(&ingredientTypes).Error; err != nil {
		return err
	}

	if err := SeedAll(gdb); err != nil {
		return err
	}

	return nil
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
	IsAllergen bool
	ParentName *string
	Image      *IngredientImageSeed
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
		Assign(entities.Ingredient{TypeID: typeID, IsAllergen: seed.IsAllergen}).
		FirstOrCreate(&ing).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create ingredient: %w", err)
	}

	if seed.Image != nil {
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
	Region      *string
	TimeMinutes *int
	Difficulty  *int
	KcalPerServing *int
	Status      string // "draft" | "published"
	NeedsPrep   bool
	Photo       *RecipePhotoSeed

	Categories    []string
	DietTags      []string
	PracticalTags []string

	Steps []RecipeStepSeed
	Items []RecipeIngredientSeed
}

type RecipeStepSeed struct {
	StepNo  int
	Section *string
	Text    string
}

type RecipeIngredientSeed struct {
	SortOrder int
	Section   *string

	IngredientName       string
	IngredientParentName *string

	Amount *float64
	Unit   *string
	Note   *string
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
		status = "published"
	}

	var r entities.Recipe
	if err := tx.Where(entities.Recipe{Name: seed.Name}).
		Assign(entities.Recipe{
			Servings:       seed.Servings,
			Description:    seed.Description,
			Tagline:        seed.Tagline,
			Region:         seed.Region,
			TimeMinutes:    seed.TimeMinutes,
			Difficulty:     seed.Difficulty,
			KcalPerServing: seed.KcalPerServing,
			Status:         status,
			NeedsPrep:      seed.NeedsPrep,
		}).
		FirstOrCreate(&r).Error; err != nil {
		return 0, fmt.Errorf("failed to find or create recipe: %w", err)
	}

	if seed.Photo != nil && seed.Photo.URL != "" {
		if err := upsertRecipePhoto(tx, r.ID, seed.Photo.URL); err != nil {
			return 0, err
		}
	}

	// Steps
	for _, s := range seed.Steps {
		step := entities.RecipeStep{
			RecipeID: r.ID, StepNo: s.StepNo, Section: s.Section, Text: s.Text,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "recipe_id"}, {Name: "step_no"}},
			DoUpdates: clause.Assignments(map[string]any{
				"section": step.Section,
				"text":    step.Text,
			}),
		}).Create(&step).Error; err != nil {
			return 0, err
		}
	}

	// Ingredients
	for _, it := range seed.Items {
		ingID, err := findIngredientID(tx, it.IngredientName, it.IngredientParentName)
		if err != nil {
			return 0, err
		}

		item := entities.RecipeIngredient{
			RecipeID:     r.ID,
			SortOrder:    it.SortOrder,
			Section:      it.Section,
			IngredientID: ingID,
			Amount:       it.Amount,
			Unit:         it.Unit,
			Note:         it.Note,
		}
		if err := tx.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "recipe_id"}, {Name: "sort_order"}},
			DoUpdates: clause.Assignments(map[string]any{
				"section":       item.Section,
				"ingredient_id": item.IngredientID,
				"amount":        item.Amount,
				"unit":          item.Unit,
				"note":          item.Note,
			}),
		}).Create(&item).Error; err != nil {
			return 0, err
		}
	}

	// Tags (categories + diet + practical) — usuń stare, wstaw nowe
	tx.Where("recipe_id = ?", r.ID).Delete(&entities.RecipeTag{})
	for _, tag := range seed.Categories {
		tx.Create(&entities.RecipeTag{RecipeID: r.ID, Tag: tag, GroupName: "category"})
	}
	for _, tag := range seed.DietTags {
		tx.Create(&entities.RecipeTag{RecipeID: r.ID, Tag: tag, GroupName: "diet"})
	}
	for _, tag := range seed.PracticalTags {
		tx.Create(&entities.RecipeTag{RecipeID: r.ID, Tag: tag, GroupName: "practical"})
	}

	return r.ID, nil
}

func SeedAll(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {

		// ── Składniki ────────────────────────────────────────────
		ings := []IngredientSeed{
			// Warzywa
			{Name: "Cebula", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "bC1fXU1v98U"}},
			{Name: "Czosnek", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "_b8n1KTZ8rw"}},
			{Name: "Pomidor", TypeCode: "FRUIT", Image: &IngredientImageSeed{PostID: "Jy7e7RjOFVo"}},
			{Name: "Papryczka chili", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "iBsmi-wCXNE"}},
			{Name: "Marchew", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "0bPFoBML0qs"}},
			{Name: "Ziemniaki", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "rB7-iPJADpI"}},
			{Name: "Szpinak", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "GI39ElkDQFo"}},
			{Name: "Cukinia", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "SbzMoWL0xvw"}},
			{Name: "Papryka czerwona", TypeCode: "VEGETABLE", Image: &IngredientImageSeed{PostID: "iKKOhzHXR9E"}},
			// Owoce
			{Name: "Cytryna", TypeCode: "FRUIT", Image: &IngredientImageSeed{PostID: "7WAGthfGJ9w"}},
			{Name: "Limonka", TypeCode: "FRUIT"},
			// Zioła
			{Name: "Bazylia", TypeCode: "HERB", Image: &IngredientImageSeed{PostID: "DDKdfRe1GxA"}},
			{Name: "Pietruszka (nać)", TypeCode: "HERB", Image: &IngredientImageSeed{PostID: "r6BUzN_jTHg"}},
			{Name: "Tymianek", TypeCode: "HERB"},
			{Name: "Rozmaryn", TypeCode: "HERB"},
			{Name: "Oregano", TypeCode: "HERB"},
			// Mięso
			{Name: "Kurczak", TypeCode: "MEAT", Image: &IngredientImageSeed{PostID: "pNcFMdEe09Q"}},
			{Name: "Wołowina", TypeCode: "MEAT", Image: &IngredientImageSeed{PostID: "u2QnkcLNsBY"}},
			{Name: "Boczek", TypeCode: "MEAT"},
			// Ryby
			{Name: "Łosoś", TypeCode: "FISH"},
			{Name: "Dorsz", TypeCode: "FISH"},
			// Inne
			{Name: "Masło", TypeCode: "OTHER", IsAllergen: true},
			{Name: "Oliwa z oliwek", TypeCode: "OTHER"},
			{Name: "Makaron spaghetti", TypeCode: "OTHER"},
			{Name: "Parmezan", TypeCode: "OTHER", IsAllergen: true},
			{Name: "Jajka", TypeCode: "OTHER", IsAllergen: true},
			{Name: "Śmietana 30%", TypeCode: "OTHER", IsAllergen: true},
			{Name: "Bulion warzywny", TypeCode: "OTHER"},
			// Przyprawy
			{Name: "Sól", TypeCode: "SPICE"},
			{Name: "Pieprz czarny", TypeCode: "SPICE"},
			{Name: "Kurkuma", TypeCode: "SPICE"},
			{Name: "Kmin rzymski", TypeCode: "SPICE"},
			{Name: "Papryka słodka", TypeCode: "SPICE"},
			{Name: "Papryka wędzona", TypeCode: "SPICE"},
		}

		for _, s := range ings {
			if _, err := SeedOneIngredient(tx, s); err != nil {
				return err
			}
		}

		// ── Przepisy ─────────────────────────────────────────────

		// helpers
		secSos := strPtr("Sos")
		secMakaron := strPtr("Makaron")
		secDressing := strPtr("Dressing")
		secKurczak := strPtr("Kurczak")

		wloska := strPtr("włoska")
		azjatycka := strPtr("azjatycka")
		srodziemnomorska := strPtr("śródziemnomorska")

		recipes := []RecipeSeed{
			{
				Name:        "Spaghetti aglio e olio",
				Servings:    2,
				Tagline:     strPtr("Trzy składniki, gotowe w 15 minut. Klasyk za który nie wypada przepraszać."),
				Region:      wloska,
				TimeMinutes: intPtr(20),
				Difficulty:  intPtr(1),
				Status:      "published",
				Categories:  []string{"makarony"},
				PracticalTags: []string{"30 min", "1 patelnia", "do 5 składników"},
				Steps: []RecipeStepSeed{
					{StepNo: 1, Section: secMakaron, Text: "Ugotuj makaron al dente w mocno osolonej wodzie. Zachowaj szklankę wody z gotowania."},
					{StepNo: 2, Section: secSos, Text: "Pokrój czosnek w cienkie plasterki. Na dużej patelni rozgrzej oliwę na średnim ogniu — czosnek ma powoli się złocić, nie smażyć. Dodaj chili."},
					{StepNo: 3, Section: secSos, Text: "Przełóż makaron na patelnię, dolej kilka łyżek wody z gotowania. Mieszaj aż sos będzie jedwabisty i połączy się z makaronem."},
					{StepNo: 4, Text: "Posyp pietruszką, opcjonalnie startym parmezanem. Podawaj od razu."},
				},
				Items: []RecipeIngredientSeed{
					{SortOrder: 1, Section: secMakaron, IngredientName: "Makaron spaghetti", Amount: f64Ptr(200), Unit: strPtr("g")},
					{SortOrder: 2, Section: secSos, IngredientName: "Czosnek", Amount: f64Ptr(4), Unit: strPtr("ząbek")},
					{SortOrder: 3, Section: secSos, IngredientName: "Oliwa z oliwek", Amount: f64Ptr(4), Unit: strPtr("łyżka")},
					{SortOrder: 4, Section: secSos, IngredientName: "Papryczka chili", Amount: f64Ptr(1), Unit: strPtr("sztuka"), Note: strPtr("lub szczypta płatków chili")},
					{SortOrder: 5, IngredientName: "Pietruszka (nać)", Amount: f64Ptr(1), Unit: strPtr("garść")},
					{SortOrder: 6, IngredientName: "Parmezan", Amount: f64Ptr(30), Unit: strPtr("g"), Note: strPtr("opcjonalnie")},
				},
			},
			{
				Name:          "Pieczony łosoś z cytryną i tymiankiem",
				Servings:      2,
				Tagline:       strPtr("Maślany, soczysty łosoś gotowy w 20 minut. Jedna blacha, zero bałaganu."),
				Region:        srodziemnomorska,
				TimeMinutes:   intPtr(25),
				Difficulty:    intPtr(1),
				KcalPerServing: intPtr(420),
				Status:        "published",
				Categories:    []string{"ryby"},
				DietTags:      []string{"bez glutenu"},
				PracticalTags: []string{"30 min", "1 patelnia", "do 5 składników"},
				Steps: []RecipeStepSeed{
					{StepNo: 1, Text: "Rozgrzej piekarnik do 200°C. Blachę wyłóż papierem do pieczenia."},
					{StepNo: 2, Text: "Filety połóż skórą w dół. Posmaruj roztopionym masłem, posól i popieprzysz. Na wierzch połóż gałązki tymianku i plasterki cytryny."},
					{StepNo: 3, Text: "Piecz 15–18 minut — łosoś jest gotowy gdy mięso łatwo rozdziela się na płaty. Jeśli filety są cienkie, skróć do 12 minut."},
					{StepNo: 4, Text: "Skrop sokiem z reszty cytryny i podawaj od razu."},
				},
				Items: []RecipeIngredientSeed{
					{SortOrder: 1, IngredientName: "Łosoś", Amount: f64Ptr(2), Unit: strPtr("filety"), Note: strPtr("ok. 150g każdy")},
					{SortOrder: 2, IngredientName: "Masło", Amount: f64Ptr(2), Unit: strPtr("łyżka")},
					{SortOrder: 3, IngredientName: "Cytryna", Amount: f64Ptr(1), Unit: strPtr("sztuka")},
					{SortOrder: 4, IngredientName: "Tymianek", Amount: f64Ptr(3), Unit: strPtr("garść"), Note: strPtr("świeży")},
					{SortOrder: 5, IngredientName: "Sól", Amount: f64Ptr(1), Unit: strPtr("szczypta")},
					{SortOrder: 6, IngredientName: "Pieprz czarny", Amount: f64Ptr(1), Unit: strPtr("szczypta")},
				},
			},
			{
				Name:          "Kurczak tikka masala",
				Servings:      4,
				Tagline:       strPtr("Kremowy, aromatyczny sos z kurczakiem. Nie potrzebujesz żadnego tandooru."),
				Region:        azjatycka,
				TimeMinutes:   intPtr(50),
				Difficulty:    intPtr(2),
				KcalPerServing: intPtr(510),
				Status:        "published",
				NeedsPrep:     true,
				Categories:    []string{"mieso"},
				DietTags:      []string{"bez glutenu"},
				Steps: []RecipeStepSeed{
					{StepNo: 1, Section: secKurczak, Text: "Kurczaka pokrój w kawałki. Wymieszaj z kurkumą, kminem, papryką wędzoną i oliwą. Marynuj minimum 30 minut (najlepiej noc)."},
					{StepNo: 2, Section: secKurczak, Text: "Obsmaż kurczaka na mocno rozgrzanej patelni z każdej strony — partiami, żeby się nie dusił. Odłóż."},
					{StepNo: 3, Section: secSos, Text: "Na tej samej patelni podsmaż cebulę do zezłocenia. Dodaj czosnek i chili, smaż 1 minutę."},
					{StepNo: 4, Section: secSos, Text: "Wrzuć pomidory, dolej szklankę bulionu. Gotuj na małym ogniu 15 minut aż sos zgęstnieje."},
					{StepNo: 5, Section: secSos, Text: "Dodaj śmietanę i kurczaka z powrotem. Gotuj 5 minut. Dopraw do smaku."},
					{StepNo: 6, Text: "Posyp kolendrą lub pietruszką. Podawaj z ryżem lub chlebem naan."},
				},
				Items: []RecipeIngredientSeed{
					{SortOrder: 1, Section: secKurczak, IngredientName: "Kurczak", Amount: f64Ptr(600), Unit: strPtr("g"), Note: strPtr("piersi lub udka bez kości")},
					{SortOrder: 2, Section: secKurczak, IngredientName: "Kurkuma", Amount: f64Ptr(1), Unit: strPtr("łyżeczka")},
					{SortOrder: 3, Section: secKurczak, IngredientName: "Kmin rzymski", Amount: f64Ptr(1), Unit: strPtr("łyżeczka")},
					{SortOrder: 4, Section: secKurczak, IngredientName: "Papryka wędzona", Amount: f64Ptr(1), Unit: strPtr("łyżeczka")},
					{SortOrder: 5, Section: secKurczak, IngredientName: "Oliwa z oliwek", Amount: f64Ptr(2), Unit: strPtr("łyżka")},
					{SortOrder: 6, Section: secSos, IngredientName: "Cebula", Amount: f64Ptr(1), Unit: strPtr("sztuka")},
					{SortOrder: 7, Section: secSos, IngredientName: "Czosnek", Amount: f64Ptr(4), Unit: strPtr("ząbek")},
					{SortOrder: 8, Section: secSos, IngredientName: "Papryczka chili", Amount: f64Ptr(1), Unit: strPtr("sztuka")},
					{SortOrder: 9, Section: secSos, IngredientName: "Pomidor", Amount: f64Ptr(400), Unit: strPtr("g"), Note: strPtr("z puszki lub świeże")},
					{SortOrder: 10, Section: secSos, IngredientName: "Bulion warzywny", Amount: f64Ptr(250), Unit: strPtr("ml")},
					{SortOrder: 11, Section: secSos, IngredientName: "Śmietana 30%", Amount: f64Ptr(150), Unit: strPtr("ml")},
					{SortOrder: 12, Section: secSos, IngredientName: "Papryka słodka", Amount: f64Ptr(2), Unit: strPtr("łyżeczka")},
				},
			},
			{
				Name:          "Sałatka z pieczonym boczkiem i szpinakiem",
				Servings:      2,
				Tagline:       strPtr("Ciepły dressing, chrupiący boczek i świeży szpinak. Gotowe w 15 minut."),
				TimeMinutes:   intPtr(15),
				Difficulty:    intPtr(1),
				Status:        "published",
				Categories:    []string{"salatki"},
				PracticalTags: []string{"30 min", "1 patelnia"},
				Steps: []RecipeStepSeed{
					{StepNo: 1, Section: secDressing, Text: "Na patelni usmaż boczek na chrupiąco. Zostaw 2 łyżki tłuszczu na patelni."},
					{StepNo: 2, Section: secDressing, Text: "Do tłuszczu dodaj czosnek, smaż 30 sekund. Dodaj sok z cytryny i szczyptę pieprzu — dressing gotowy."},
					{StepNo: 3, Text: "Szpinak ułóż w misce. Polej ciepłym dressingiem — liście lekko zwiotczeją. Dodaj boczek i od razu podawaj."},
				},
				Items: []RecipeIngredientSeed{
					{SortOrder: 1, IngredientName: "Szpinak", Amount: f64Ptr(150), Unit: strPtr("g"), Note: strPtr("baby spinach")},
					{SortOrder: 2, IngredientName: "Boczek", Amount: f64Ptr(100), Unit: strPtr("g")},
					{SortOrder: 3, Section: secDressing, IngredientName: "Czosnek", Amount: f64Ptr(2), Unit: strPtr("ząbek")},
					{SortOrder: 4, Section: secDressing, IngredientName: "Cytryna", Amount: f64Ptr(0.5), Unit: strPtr("sztuka"), Note: strPtr("sok")},
					{SortOrder: 5, Section: secDressing, IngredientName: "Pieprz czarny", Amount: f64Ptr(1), Unit: strPtr("szczypta")},
				},
			},
		}

		for _, r := range recipes {
			if _, err := SeedOneRecipe(tx, r); err != nil {
				return err
			}
		}

		return nil
	})
}

func strPtr(s string) *string  { return &s }
func intPtr(i int) *int        { return &i }
func f64Ptr(f float64) *float64 { return &f }
