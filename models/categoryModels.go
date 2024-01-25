package models

import "time"

type Category struct {
	CategoryID   int      `gorm:"primaryKey;column:category_id" json:"category_id"`
	CategoryName string    `gorm:"column:category_name" json:"category_name"`
	IsDeleted    bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedBy   string    `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// TableName specifies the table name for the model
func (Category) TableName() string {
	return "categories"
}