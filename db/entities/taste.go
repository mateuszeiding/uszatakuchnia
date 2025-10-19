package entities

type Taste struct {
	ID     uint      `gorm:"primaryKey;autoIncrement"`
	TypeID uint      `gorm:"not null;uniqueIndex:ux_taste_type_id"`
	Type   TasteType `gorm:"foreignKey:TypeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT"`
}

func (Taste) TableName() string { return "taste" }
