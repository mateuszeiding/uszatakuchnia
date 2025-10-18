package entities

type AromaType struct {
	ID   int    `gorm:"primaryKey;autoIncrement"`
	Code string `gorm:"type:varchar(64);not null;uniqueIndex:ux_aroma_type_code"`
	Name string `gorm:"type:varchar(128);not null"`
}

func (AromaType) TableName() string { return "aroma_types" }
