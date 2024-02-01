package models

import (
	"time"
	"github.com/go-playground/validator/v10"
)

type User struct {
	UserID       uint      `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username     string    `gorm:"column:username;not null" json:"username" validate:"required,max=100"`
	Fullname     string    `gorm:"column:fullname;not null" json:"fullname" validate:"required,max=255"`
	Password     string    `gorm:"column:password;not null" json:"password" validate:"required,min=6,max=50"`
	Role         string    `gorm:"column:role;not null" json:"role"`
	IsDeleted    bool      `gorm:"column:is_deleted;not null" json:"is_deleted"`
	CreatedBy    string    `gorm:"column:created_by" json:"created_by"`
	CreatedTime  time.Time `gorm:"column:created_time" json:"created_time"`
	ModifiedBy   string    `gorm:"column:modified_by" json:"modified_by"`
	ModifiedTime time.Time `gorm:"column:modified_time" json:"modified_time"`
}

// TableName specifies the table name for the model
func (User) TableName() string {
	return "users"
}

func (u *User) Validate() error {
    validate := validator.New()
    return validate.Struct(u)
}
