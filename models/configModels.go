package models

import "time"

type SettingConfig struct {
	ConfigID int64  `gorm:"primaryKey;autoIncrement;column:config_id" json:"config_id"`
	Activity string `gorm:"column:activity;size:255" json:"activity"`
	Message  string `gorm:"column:message" json:"message"`
	IsActive *bool  `gorm:"column:is_active" json:"is_active"`
	BaseModel
	IsDeleted *bool `gorm:"-" json:"-"`
}

func (SettingConfig) TableName() string {
	return "setting_config"
}

type SettingConfigRequest struct {
	Activity string `json:"activity" binding:"required,max=255"`
	Message  string `json:"message" binding:"required"`
	IsActive *bool  `json:"is_active" binding:"required"`
}

type Parameter struct {
	ParamID          int64      `gorm:"primaryKey;autoIncrement;column:param_id" json:"param_id"`
	TahunAjar        string     `gorm:"column:tahun_ajar;size:9" json:"tahun_ajar"`
	BiayaPendaftaran float64    `gorm:"column:biaya_pendaftaran" json:"biaya_pendaftaran"`
	Ups              float64    `gorm:"column:ups" json:"ups"`
	Usm              float64    `gorm:"column:usm" json:"usm"`
	TglParam         *time.Time `gorm:"column:tgl_param" json:"tgl_param"`
	MatpelUjian1     *int64     `gorm:"column:matpel_ujian1" json:"matpel_ujian1"`
	MatpelUjian2     *int64     `gorm:"column:matpel_ujian2" json:"matpel_ujian2"`
	MatpelUjian3     *int64     `gorm:"column:matpel_ujian3" json:"matpel_ujian3"`
	MatpelUjian4     *int64     `gorm:"column:matpel_ujian4" json:"matpel_ujian4"`
	BaseModel
	IsDeleted *bool `gorm:"-" json:"-"`
}

func (Parameter) TableName() string {
	return "parameter"
}

type ParameterRequest struct {
	TahunAjar        string     `json:"tahun_ajar" binding:"required,max=9"`
	BiayaPendaftaran float64    `json:"biaya_pendaftaran"`
	Ups              float64    `json:"ups"`
	Usm              float64    `json:"usm"`
	TglParam         *time.Time `json:"tgl_param"`
	MatpelUjian1     *int64     `json:"matpel_ujian1"`
	MatpelUjian2     *int64     `json:"matpel_ujian2"`
	MatpelUjian3     *int64     `json:"matpel_ujian3"`
	MatpelUjian4     *int64     `json:"matpel_ujian4"`
}
