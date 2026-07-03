package services

import (
	"errors"
	"ppdb-be/models"
	"ppdb-be/repositories"
)

type MasterSchoolService interface {
	CreateSchool(req models.MasterSchoolRequest, userID int64) error
	GetSchools() ([]models.MasterSchool, error)
	GetSchoolByID(id int64) (*models.MasterSchool, error)
	UpdateSchool(id int64, req models.MasterSchoolRequest, userID int64) error
	DeleteSchool(id int64) error
}

type masterSchoolService struct {
	repo repositories.MasterSchoolRepository
}

func NewMasterSchoolService(repo repositories.MasterSchoolRepository) MasterSchoolService {
	return &masterSchoolService{repo: repo}
}

func (s *masterSchoolService) CreateSchool(req models.MasterSchoolRequest, userID int64) error {
	school := &models.MasterSchool{
		NamaSmp:   req.NamaSmp,
		AlamatSmp: req.AlamatSmp,
	}
	school.CreatedBy = &userID
	return s.repo.Create(school)
}

func (s *masterSchoolService) GetSchools() ([]models.MasterSchool, error) {
	return s.repo.FindAll()
}

func (s *masterSchoolService) GetSchoolByID(id int64) (*models.MasterSchool, error) {
	return s.repo.FindByID(id)
}

func (s *masterSchoolService) UpdateSchool(id int64, req models.MasterSchoolRequest, userID int64) error {
	school, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("master school not found")
	}

	school.NamaSmp = req.NamaSmp
	school.AlamatSmp = req.AlamatSmp
	school.LastModifiedBy = &userID

	return s.repo.Update(school)
}

func (s *masterSchoolService) DeleteSchool(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("master school not found")
	}
	return s.repo.Delete(id)
}
