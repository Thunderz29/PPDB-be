package repositories

import (
	"ppdb-be/models"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	// SettingConfig
	CreateSettingConfig(config *models.SettingConfig) error
	FindAllSettingConfig() ([]models.SettingConfig, error)
	FindSettingConfigByID(id int64) (*models.SettingConfig, error)
	UpdateSettingConfig(config *models.SettingConfig) error
	DeleteSettingConfig(id int64) error

	// Parameter
	CreateParameter(param *models.Parameter) error
	FindAllParameter() ([]models.Parameter, error)
	FindParameterByID(id int64) (*models.Parameter, error)
	UpdateParameter(param *models.Parameter) error
	DeleteParameter(id int64) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

// --- Setting Config ---
func (r *configRepository) CreateSettingConfig(config *models.SettingConfig) error {
	return r.db.Create(config).Error
}

func (r *configRepository) FindAllSettingConfig() ([]models.SettingConfig, error) {
	var configs []models.SettingConfig
	err := r.db.Find(&configs).Error
	return configs, err
}

func (r *configRepository) FindSettingConfigByID(id int64) (*models.SettingConfig, error) {
	var config models.SettingConfig
	err := r.db.Where("config_id = ?", id).First(&config).Error
	return &config, err
}

func (r *configRepository) UpdateSettingConfig(config *models.SettingConfig) error {
	return r.db.Save(config).Error
}

func (r *configRepository) DeleteSettingConfig(id int64) error {
	return r.db.Where("config_id = ?", id).Delete(&models.SettingConfig{}).Error
}

// --- Parameter ---
func (r *configRepository) CreateParameter(param *models.Parameter) error {
	return r.db.Create(param).Error
}

func (r *configRepository) FindAllParameter() ([]models.Parameter, error) {
	var params []models.Parameter
	err := r.db.Find(&params).Error
	return params, err
}

func (r *configRepository) FindParameterByID(id int64) (*models.Parameter, error) {
	var param models.Parameter
	err := r.db.Where("param_id = ?", id).First(&param).Error
	return &param, err
}

func (r *configRepository) UpdateParameter(param *models.Parameter) error {
	return r.db.Save(param).Error
}

func (r *configRepository) DeleteParameter(id int64) error {
	return r.db.Where("param_id = ?", id).Delete(&models.Parameter{}).Error
}
