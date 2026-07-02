package models

type User struct {
	UserID       int64   `gorm:"primaryKey;column:user_id;autoIncrement" json:"user_id"`
	RoleID       *int64  `gorm:"column:role_id" json:"role_id"`
	UserFullName string  `gorm:"column:user_full_name;type:varchar(50);not null" json:"user_full_name"`
	UserName     string  `gorm:"column:user_name;type:varchar(20);not null" json:"user_name"`
	UserPhone    *string `gorm:"column:user_phone;type:varchar(15)" json:"user_phone"`
	UserEmail    string  `gorm:"column:user_email;type:varchar(50);not null" json:"user_email"`
	UserPassword string  `gorm:"column:user_password;type:text" json:"-"`
	BaseModel
}

func (User) TableName() string {
	return "users"
}
