package entities

type Taste struct {
	ID     int       `gorm:"primaryKey;autoIncrement"`
	TypeID int       `gorm:"not null;uniqueIndex:ux_taste_type_id"`
	Type   TasteType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Taste) TableName() string { return "taste" }
