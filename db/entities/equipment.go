package entities

// Equipment is the equipment dictionary (stand mixer, springform pan, thermometer…). Name is unique
// case-insensitively (index ux_equipment_name_lower, migration 0007).
type Equipment struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Name string `gorm:"type:varchar(120);not null"`
}

func (Equipment) TableName() string { return "equipment" }

// RecipeEquipment is an equipment item in a recipe. Deleting the recipe deletes its items; equipment
// used in a recipe can't be deleted from the dictionary without replacing or removing it there (API).
type RecipeEquipment struct {
	RecipeID  uint `gorm:"not null;primaryKey"`
	SortOrder int  `gorm:"not null;primaryKey"`

	EquipmentID uint      `gorm:"not null;index:ix_receq_equipment"`
	Equipment   Equipment `gorm:"foreignKey:EquipmentID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`

	// Details, e.g. "3–4 l, grube dno", "albo ręcznie — ok. 10–12 minut".
	Spec     *string `gorm:"type:varchar(160)"`
	Optional bool    `gorm:"not null;default:false"`
}

func (RecipeEquipment) TableName() string { return "recipe_equipment" }
