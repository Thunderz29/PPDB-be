package models

type SchoolProfile struct {
	ID               int64  `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	NamaSekolah      string `gorm:"column:nama_sekolah;size:255;not null" json:"nama_sekolah"`
	TingkatanSekolah string `gorm:"column:tingkatan_sekolah;size:255;not null" json:"tingkatan_sekolah"`
	Alamat           string `gorm:"column:alamat;not null" json:"alamat"`
	TahunBerdiri     string `gorm:"column:tahun_berdiri;size:255;not null" json:"tahun_berdiri"`
	Logo             string `gorm:"column:logo" json:"logo"`
	Mail             string `gorm:"column:mail;size:255;not null" json:"mail"`
	Phone            string `gorm:"column:phone;size:255;not null" json:"phone"`
	InformasiLainnya string `gorm:"column:informasi_lainnya" json:"informasi_lainnya"`
	SystemName       string `gorm:"column:system_name;size:255;default:'PPDB'" json:"system_name"`
	BaseModel
	IsDeleted *bool `gorm:"-" json:"-"`
}

func (SchoolProfile) TableName() string {
	return "profil_sekolah"
}

type SchoolProfileRequest struct {
	NamaSekolah      string `json:"nama_sekolah" binding:"required,max=255"`
	TingkatanSekolah string `json:"tingkatan_sekolah" binding:"required,max=255"`
	Alamat           string `json:"alamat" binding:"required"`
	TahunBerdiri     string `json:"tahun_berdiri" binding:"required,max=255"`
	Logo             string `json:"logo"`
	Mail             string `json:"mail" binding:"required,email,max=255"`
	Phone            string `json:"phone" binding:"required,max=255"`
	InformasiLainnya string `json:"informasi_lainnya"`
	SystemName       string `json:"system_name" binding:"max=255"`
}
