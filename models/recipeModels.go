package models

import (
	"time"
)

type Recipe struct {
	RecipeID      int           `gorm:"primaryKey;column:recipe_id" json:"recipe_id"`
	CategoryID    int           `gorm:"column:category_id" json:"category_id"`
	UserID        int           `gorm:"column:user_id" json:"user_id"`
	LevelID       int           `gorm:"column:level_id" json:"level_id"`
	RecipeName    string         `gorm:"column:recipe_name" json:"recipe_name"`
	ImageFilename string         `gorm:"column:image_filename" json:"image_filename"`
	TimeCook      *int           `gorm:"column:time_cook" json:"time_cook"`
	Ingridient    string         `gorm:"type:text;column:ingridient" json:"ingridient"`
	HowToCook     string         `gorm:"type:text;column:how_to_cook" json:"how_to_cook"`
	IsDeleted     bool           `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedBy     string         `gorm:"column:created_by" json:"created_by"`
	CreatedTime   time.Time      `gorm:"column:created_time" json:"created_time"`
	ModifiedBy    string         `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime  time.Time      `gorm:"column:modified_time" json:"modified_time"`
	Category      Category `gorm:"foreignKey:CategoryID;references:CategoryID" json:"category"`
	User          User     `gorm:"foreignKey:UserID;references:UserID" json:"user"`
	Level         Level    `gorm:"foreignKey:LevelID;references:LevelID" json:"level"`
}

// TableName specifies the table name for the model
func (Recipe) TableName() string {
	return "recipes"
}
