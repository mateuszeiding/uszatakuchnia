package entities

type TasteType struct {
	ID   uint   `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"type:varchar(64);not null;uniqueIndex:ux_taste_type_code"`
	Name string `gorm:"type:varchar(128);not null"`
}

func (TasteType) TableName() string { return "taste_types" }
