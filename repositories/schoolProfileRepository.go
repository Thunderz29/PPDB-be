package repositories

import (
	"ppdb-be/models"

	"gorm.io/gorm"
)

type SchoolProfileRepository interface {
	Create(profile *models.SchoolProfile) error
	FindAll() ([]models.SchoolProfile, error)
	FindByID(id int64) (*models.SchoolProfile, error)
	Update(profile *models.SchoolProfile) error
	Delete(id int64) error
}

type schoolProfileRepository struct {
	db *gorm.DB
}

func NewSchoolProfileRepository(db *gorm.DB) SchoolProfileRepository {
	return &schoolProfileRepository{db: db}
}

func (r *schoolProfileRepository) Create(profile *models.SchoolProfile) error {
	return r.db.Create(profile).Error
}

func (r *schoolProfileRepository) FindAll() ([]models.SchoolProfile, error) {
	var profiles []models.SchoolProfile
	err := r.db.Find(&profiles).Error
	return profiles, err
}

func (r *schoolProfileRepository) FindByID(id int64) (*models.SchoolProfile, error) {
	var profile models.SchoolProfile
	err := r.db.Where("id = ?", id).First(&profile).Error
	return &profile, err
}

func (r *schoolProfileRepository) Update(profile *models.SchoolProfile) error {
	return r.db.Save(profile).Error
}

func (r *schoolProfileRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&models.SchoolProfile{}).Error
}
