package services

import (
	"errors"
	"ppdb-be/models"
	"ppdb-be/repositories"
)

type SchoolProfileService interface {
	CreateSchoolProfile(req models.SchoolProfileRequest, userID int64) error
	GetSchoolProfiles() ([]models.SchoolProfile, error)
	GetSchoolProfileByID(id int64) (*models.SchoolProfile, error)
	UpdateSchoolProfile(id int64, req models.SchoolProfileRequest, userID int64) error
	DeleteSchoolProfile(id int64) error
}

type schoolProfileService struct {
	repo repositories.SchoolProfileRepository
}

func NewSchoolProfileService(repo repositories.SchoolProfileRepository) SchoolProfileService {
	return &schoolProfileService{repo: repo}
}

func (s *schoolProfileService) CreateSchoolProfile(req models.SchoolProfileRequest, userID int64) error {
	profile := &models.SchoolProfile{
		NamaSekolah:      req.NamaSekolah,
		TingkatanSekolah: req.TingkatanSekolah,
		Alamat:           req.Alamat,
		TahunBerdiri:     req.TahunBerdiri,
		Logo:             req.Logo,
		Mail:             req.Mail,
		Phone:            req.Phone,
		InformasiLainnya: req.InformasiLainnya,
		SystemName:       req.SystemName,
	}
	if profile.SystemName == "" {
		profile.SystemName = "PPDB"
	}
	profile.CreatedBy = &userID
	return s.repo.Create(profile)
}

func (s *schoolProfileService) GetSchoolProfiles() ([]models.SchoolProfile, error) {
	return s.repo.FindAll()
}

func (s *schoolProfileService) GetSchoolProfileByID(id int64) (*models.SchoolProfile, error) {
	return s.repo.FindByID(id)
}

func (s *schoolProfileService) UpdateSchoolProfile(id int64, req models.SchoolProfileRequest, userID int64) error {
	profile, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("school profile not found")
	}

	profile.NamaSekolah = req.NamaSekolah
	profile.TingkatanSekolah = req.TingkatanSekolah
	profile.Alamat = req.Alamat
	profile.TahunBerdiri = req.TahunBerdiri
	profile.Logo = req.Logo
	profile.Mail = req.Mail
	profile.Phone = req.Phone
	profile.InformasiLainnya = req.InformasiLainnya
	profile.SystemName = req.SystemName
	if profile.SystemName == "" {
		profile.SystemName = "PPDB"
	}

	profile.LastModifiedBy = &userID

	return s.repo.Update(profile)
}

func (s *schoolProfileService) DeleteSchoolProfile(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("school profile not found")
	}
	return s.repo.Delete(id)
}
