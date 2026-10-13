package entities

// Unit — jednostka miary. Formy odmiany dla liczebników: 1 / 2–4 / 5+ (np. ząbek / ząbki / ząbków).
// Brak form = jednostka nieodmienna (g, ml).
type Unit struct {
	ID        uint    `gorm:"primaryKey;autoIncrement"`
	Name      string  `gorm:"type:varchar(32);not null;uniqueIndex:ux_units_name"`
	FormOne   *string `gorm:"type:varchar(32)"`
	FormFew   *string `gorm:"type:varchar(32)"`
	FormMany  *string `gorm:"type:varchar(32)"`
	NoQty     bool    `gorm:"not null;default:false"` // "do smaku" — bez ilości
	SortOrder int     `gorm:"not null;default:0"`
}

func (Unit) TableName() string { return "units" }
