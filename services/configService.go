package services

import (
	"errors"
	"ppdb-be/models"
	"ppdb-be/repositories"
)

type ConfigService interface {
	// SettingConfig
	CreateSettingConfig(req models.SettingConfigRequest, userID int64) error
	GetSettingConfigs() ([]models.SettingConfig, error)
	GetSettingConfigByID(id int64) (*models.SettingConfig, error)
	UpdateSettingConfig(id int64, req models.SettingConfigRequest, userID int64) error
	DeleteSettingConfig(id int64) error

	// Parameter
	CreateParameter(req models.ParameterRequest, userID int64) error
	GetParameters() ([]models.Parameter, error)
	GetParameterByID(id int64) (*models.Parameter, error)
	UpdateParameter(id int64, req models.ParameterRequest, userID int64) error
	DeleteParameter(id int64) error
}

type configService struct {
	repo repositories.ConfigRepository
}

func NewConfigService(repo repositories.ConfigRepository) ConfigService {
	return &configService{repo: repo}
}

// --- Setting Config ---
func (s *configService) CreateSettingConfig(req models.SettingConfigRequest, userID int64) error {
	config := &models.SettingConfig{
		Activity: req.Activity,
		Message:  req.Message,
		IsActive: req.IsActive,
	}
	config.CreatedBy = &userID
	return s.repo.CreateSettingConfig(config)
}

func (s *configService) GetSettingConfigs() ([]models.SettingConfig, error) {
	return s.repo.FindAllSettingConfig()
}

func (s *configService) GetSettingConfigByID(id int64) (*models.SettingConfig, error) {
	return s.repo.FindSettingConfigByID(id)
}

func (s *configService) UpdateSettingConfig(id int64, req models.SettingConfigRequest, userID int64) error {
	config, err := s.repo.FindSettingConfigByID(id)
	if err != nil {
		return errors.New("setting config not found")
	}

	config.Activity = req.Activity
	config.Message = req.Message
	config.IsActive = req.IsActive
	config.LastModifiedBy = &userID

	return s.repo.UpdateSettingConfig(config)
}

func (s *configService) DeleteSettingConfig(id int64) error {
	_, err := s.repo.FindSettingConfigByID(id)
	if err != nil {
		return errors.New("setting config not found")
	}
	return s.repo.DeleteSettingConfig(id)
}

// --- Parameter ---
func (s *configService) CreateParameter(req models.ParameterRequest, userID int64) error {
	param := &models.Parameter{
		TahunAjar:        req.TahunAjar,
		BiayaPendaftaran: req.BiayaPendaftaran,
		Ups:              req.Ups,
		Usm:              req.Usm,
		TglParam:         req.TglParam,
		MatpelUjian1:     req.MatpelUjian1,
		MatpelUjian2:     req.MatpelUjian2,
		MatpelUjian3:     req.MatpelUjian3,
		MatpelUjian4:     req.MatpelUjian4,
	}
	param.CreatedBy = &userID
	return s.repo.CreateParameter(param)
}

func (s *configService) GetParameters() ([]models.Parameter, error) {
	return s.repo.FindAllParameter()
}

func (s *configService) GetParameterByID(id int64) (*models.Parameter, error) {
	return s.repo.FindParameterByID(id)
}

func (s *configService) UpdateParameter(id int64, req models.ParameterRequest, userID int64) error {
	param, err := s.repo.FindParameterByID(id)
	if err != nil {
		return errors.New("parameter not found")
	}

	param.TahunAjar = req.TahunAjar
	param.BiayaPendaftaran = req.BiayaPendaftaran
	param.Ups = req.Ups
	param.Usm = req.Usm
	param.TglParam = req.TglParam
	param.MatpelUjian1 = req.MatpelUjian1
	param.MatpelUjian2 = req.MatpelUjian2
	param.MatpelUjian3 = req.MatpelUjian3
	param.MatpelUjian4 = req.MatpelUjian4

	param.LastModifiedBy = &userID

	return s.repo.UpdateParameter(param)
}

func (s *configService) DeleteParameter(id int64) error {
	_, err := s.repo.FindParameterByID(id)
	if err != nil {
		return errors.New("parameter not found")
	}
	return s.repo.DeleteParameter(id)
}
