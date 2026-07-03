package models

type MasterSchool struct {
	SmpID     int64  `gorm:"primaryKey;autoIncrement;column:smp_id" json:"smp_id"`
	NamaSmp   string `gorm:"column:nama_smp;size:255" json:"nama_smp"`
	AlamatSmp string `gorm:"column:alamat_smp;size:255" json:"alamat_smp"`
	BaseModel
	IsDeleted *bool `gorm:"-" json:"-"`
}

func (MasterSchool) TableName() string {
	return "master_smp"
}

type MasterSchoolRequest struct {
	NamaSmp   string `json:"nama_smp" binding:"required,max=255"`
	AlamatSmp string `json:"alamat_smp" binding:"required,max=255"`
}
