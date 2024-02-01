package models

import (
	"time"
)

type Level struct {
	LevelID      int      `gorm:"primaryKey;column:level_id" json:"level_id"`
	LevelName    string    `gorm:"column:level_name;not null" json:"level_name"`
	IsDeleted    bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedBy   string    `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// TableName specifies the table name for the model
func (Level) TableName() string {
	return "levels"
}
