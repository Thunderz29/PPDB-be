package repositories

import (
	"ppdb-be/models"

	"gorm.io/gorm"
)

type MasterSchoolRepository interface {
	Create(school *models.MasterSchool) error
	FindAll() ([]models.MasterSchool, error)
	FindByID(id int64) (*models.MasterSchool, error)
	Update(school *models.MasterSchool) error
	Delete(id int64) error
}

type masterSchoolRepository struct {
	db *gorm.DB
}

func NewMasterSchoolRepository(db *gorm.DB) MasterSchoolRepository {
	return &masterSchoolRepository{db: db}
}

func (r *masterSchoolRepository) Create(school *models.MasterSchool) error {
	return r.db.Create(school).Error
}

func (r *masterSchoolRepository) FindAll() ([]models.MasterSchool, error) {
	var schools []models.MasterSchool
	err := r.db.Find(&schools).Error
	return schools, err
}

func (r *masterSchoolRepository) FindByID(id int64) (*models.MasterSchool, error) {
	var school models.MasterSchool
	err := r.db.Where("smp_id = ?", id).First(&school).Error
	return &school, err
}

func (r *masterSchoolRepository) Update(school *models.MasterSchool) error {
	return r.db.Save(school).Error
}

func (r *masterSchoolRepository) Delete(id int64) error {
	return r.db.Where("smp_id = ?", id).Delete(&models.MasterSchool{}).Error
}
