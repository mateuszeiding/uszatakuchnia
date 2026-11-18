package entities

// Equipment — słownik sprzętu (robot planetarny, tortownica, termometr…). Nazwa unikalna bez wielkości liter
// (indeks ux_equipment_name_lower, migracja 0007).
type Equipment struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(120);not null"`
}

func (Equipment) TableName() string { return "equipment" }

// RecipeEquipment — pozycja sprzętu w przepisie. Usunięcie przepisu usuwa pozycje; sprzętu używanego
// w przepisie nie da się usunąć ze słownika bez zamiany albo usunięcia z przepisów (API).
type RecipeEquipment struct {
	RecipeID  uint `gorm:"not null;primaryKey"`
	SortOrder int  `gorm:"not null;primaryKey"`

	EquipmentID uint      `gorm:"not null;index:ix_receq_equipment"`
	Equipment   Equipment `gorm:"foreignKey:EquipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Doprecyzowanie: „3–4 l, grube dno”, „albo ręcznie — ok. 10–12 minut”.
	Spec     *string `gorm:"type:varchar(160)"`
	Optional bool    `gorm:"not null;default:false"`
}

func (RecipeEquipment) TableName() string { return "recipe_equipment" }
