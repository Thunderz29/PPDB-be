package models

import "time"

type User struct {
	UserID       uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username     string    `gorm:"column:username" json:"username"`
	Fullname     string    `gorm:"column:fullname" json:"fullname"`
	Password     string    `gorm:"column:password" json:"password"`
	Role         string    `gorm:"column:role" json:"role"`
	IsDeleted    bool      `gorm:"column:is_deleted" json:"is_deleted"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedBy   string    `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// TableName specifies the table name for the model
func (User) TableName() string {
	return "users"
}
