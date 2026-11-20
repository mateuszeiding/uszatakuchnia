package entities

// Unit is a unit of measure. Inflected forms for numerals: 1 / 2–4 / 5+ (e.g. ząbek / ząbki / ząbków).
// No forms = the unit doesn't inflect (g, ml).
type Unit struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	Name      string  `gorm:"type:varchar(32);not null;uniqueIndex:ux_units_name"`
	FormOne   *string `gorm:"type:varchar(32)"`
	FormFew   *string `gorm:"type:varchar(32)"`
	FormMany  *string `gorm:"type:varchar(32)"`
	NoQty     bool    `gorm:"not null;default:false"` // "do smaku" (to taste): no quantity
	SortOrder int     `gorm:"not null;default:0"`
}

func (Unit) TableName() string { return "units" }
