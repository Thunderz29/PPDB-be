package models

import "time"

type User struct {
	UserID         int64      `gorm:"primaryKey;column:user_id;autoIncrement" json:"user_id"`
	RoleID         *int64     `gorm:"column:role_id" json:"role_id"`
	UserFullName   string     `gorm:"column:user_full_name;type:varchar(50);not null" json:"user_full_name"`
	UserName       string     `gorm:"column:user_name;type:varchar(20);not null" json:"user_name"`
	UserPhone      *string    `gorm:"column:user_phone;type:varchar(15)" json:"user_phone"`
	UserEmail      string     `gorm:"column:user_email;type:varchar(50);not null" json:"user_email"`
	UserPassword   string     `gorm:"column:user_password;type:text" json:"-"`
	CreatedBy      *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedOn      *time.Time `gorm:"column:created_on;type:timestamp(6)" json:"created_on"`
	LastModifiedBy *int64     `gorm:"column:last_modified_by" json:"last_modified_by"`
	LastModifiedOn *time.Time `gorm:"column:last_modified_on;type:timestamp(6)" json:"last_modified_on"`
	IsDeleted      bool       `gorm:"column:is_deleted;default:false;not null" json:"is_deleted"`
}

func (User) TableName() string {
	return "users"
}
