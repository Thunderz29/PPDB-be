package models

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	CreatedBy      *int64     `gorm:"column:created_by" json:"created_by"`
	CreatedOn      *time.Time `gorm:"column:created_on;type:timestamp(6)" json:"created_on"`
	LastModifiedBy *int64     `gorm:"column:last_modified_by" json:"last_modified_by"`
	LastModifiedOn *time.Time `gorm:"column:last_modified_on;type:timestamp(6)" json:"last_modified_on"`
	IsDeleted      bool       `gorm:"column:is_deleted;default:false;not null" json:"is_deleted"`
}

func (b *BaseModel) BeforeCreate(tx *gorm.DB) error {
	now := time.Now()
	b.CreatedOn = &now
	b.IsDeleted = false

	if ctx := tx.Statement.Context; ctx != nil {
		if val, ok := ctx.Value("userID").(int64); ok {
			b.CreatedBy = &val
		}
	}
	return nil
}

func (b *BaseModel) BeforeUpdate(tx *gorm.DB) error {
	now := time.Now()
	b.LastModifiedOn = &now

	if ctx := tx.Statement.Context; ctx != nil {
		if val, ok := ctx.Value("userID").(int64); ok {
			b.LastModifiedBy = &val
		}
	}
	return nil
}
