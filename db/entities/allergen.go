package entities

type Allergen struct {
	Code  string
	Label string
}

// Allergens — 14 alergenów UE (rozporządzenie 1169/2011). Lista stała, kolejność wyświetlania.
var Allergens = []Allergen{
	{"gluten", "gluten"},
	{"milk", "mleko"},
	{"eggs", "jaja"},
	{"fish", "ryby"},
	{"crustaceans", "skorupiaki"},
	{"molluscs", "mięczaki"},
	{"peanuts", "orzeszki ziemne"},
	{"nuts", "orzechy"},
	{"soy", "soja"},
	{"celery", "seler"},
	{"mustard", "gorczyca"},
	{"sesame", "sezam"},
	{"sulphites", "siarczyny"},
	{"lupin", "łubin"},
}

func IsAllergenCode(code string) bool {
	for _, a := range Allergens {
		if a.Code == code {
			return true
		}
	}
	return false
}

type IngredientAllergen struct {
	IngredientID uint   `gorm:"not null;primaryKey"`
	Code         string `gorm:"type:varchar(16);not null;primaryKey"`
}

func (IngredientAllergen) TableName() string { return "ingredient_allergens" }
