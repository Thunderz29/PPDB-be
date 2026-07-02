package models

import "time"

type UserSession struct {
	SessionID      int64      `gorm:"primaryKey;column:session_id;autoIncrement" json:"session_id"`
	UserID         int64      `gorm:"column:user_id;not null" json:"user_id"`
	AccessToken    string     `gorm:"column:access_token;type:text;not null" json:"access_token"`
	RefreshToken   *string    `gorm:"column:refresh_token;type:text" json:"refresh_token"`
	LoginAt        time.Time  `gorm:"column:login_at;type:timestamp(6);default:now();not null" json:"login_at"`
	LogoutAt       *time.Time `gorm:"column:logout_at;type:timestamp(6)" json:"logout_at"`
	ExpiredAt      time.Time  `gorm:"column:expired_at;type:timestamp(6);not null" json:"expired_at"`
	IPAddress      *string    `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	UserAgent      *string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	DeviceName     *string    `gorm:"column:device_name;type:varchar(100)" json:"device_name"`
	Platform       *string    `gorm:"column:platform;type:varchar(50)" json:"platform"`
	Browser        *string    `gorm:"column:browser;type:varchar(50)" json:"browser"`
	IsActive       bool       `gorm:"column:is_active;default:true;not null" json:"is_active"`
	IsRevoked      bool       `gorm:"column:is_revoked;default:false;not null" json:"is_revoked"`
	CreatedBy      *int64     `gorm:"column:created_by" json:"-"`
	CreatedOn      time.Time  `gorm:"column:created_on;type:timestamp(6);default:now();not null" json:"-"`
	LastModifiedBy *int64     `gorm:"column:last_modified_by" json:"-"`
	LastModifiedOn *time.Time `gorm:"column:last_modified_on;type:timestamp(6)" json:"-"`
}

func (UserSession) TableName() string {
	return "user_sessions"
}
