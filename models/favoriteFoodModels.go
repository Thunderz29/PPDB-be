package models

import (
	"time"
)

type FavoriteFood struct {
	UserID        int       `gorm:"primaryKey;column:user_id" json:"user_id"`
	RecipeID      int       `gorm:"primaryKey;column:recipe_id" json:"recipe_id"`
	IsFavorite    bool       `gorm:"column:is_favorite" json:"is_favorite"`
	CreatedBy     string     `gorm:"column:created_by" json:"created_by"`
	CreatedTime   time.Time  `gorm:"column:created_time" json:"created_time"`
	ModifiedBy    string     `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime  time.Time  `gorm:"column:modified_time" json:"modified_time"`
	User          User       `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Recipe        Recipe     `gorm:"foreignKey:RecipeID;references:RecipeID" json:"recipe"`
}

// TableName specifies the table name for the model
func (FavoriteFood) TableName() string {
	return "favorite_foods"
}