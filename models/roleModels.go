package models

type Role struct {
	RoleID   int64  `gorm:"primaryKey;autoIncrement;column:role_id" json:"role_id"`
	RoleName string `gorm:"column:role_name;size:50" json:"role_name"`
	BaseModel
	IsDeleted *bool `gorm:"-" json:"-"`
}

func (Role) TableName() string {
	return "role"
}

type RoleRequest struct {
	RoleName string `json:"role_name" binding:"required,max=50"`
}
