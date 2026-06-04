package entities

type RecipeTag struct {
	RecipeID  uint   `gorm:"not null;primaryKey"`
	Tag       string `gorm:"type:varchar(128);not null;primaryKey"`
	GroupName string `gorm:"type:varchar(32);not null"` // 'diet' or 'practical'

	Recipe Recipe `gorm:"foreignKey:RecipeID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

func (RecipeTag) TableName() string { return "recipe_tags" }
