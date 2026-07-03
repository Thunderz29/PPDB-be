package services

import (
	"errors"
	"ppdb-be/models"
	"ppdb-be/repositories"
)

type RoleService interface {
	CreateRole(req models.RoleRequest, userID int64) error
	GetRoles() ([]models.Role, error)
	GetRoleByID(id int64) (*models.Role, error)
	UpdateRole(id int64, req models.RoleRequest, userID int64) error
	DeleteRole(id int64) error
}

type roleService struct {
	repo repositories.RoleRepository
}

func NewRoleService(repo repositories.RoleRepository) RoleService {
	return &roleService{repo: repo}
}

func (s *roleService) CreateRole(req models.RoleRequest, userID int64) error {
	role := &models.Role{
		RoleName: req.RoleName,
	}
	role.CreatedBy = &userID
	return s.repo.Create(role)
}

func (s *roleService) GetRoles() ([]models.Role, error) {
	return s.repo.FindAll()
}

func (s *roleService) GetRoleByID(id int64) (*models.Role, error) {
	return s.repo.FindByID(id)
}

func (s *roleService) UpdateRole(id int64, req models.RoleRequest, userID int64) error {
	role, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}

	role.RoleName = req.RoleName
	role.LastModifiedBy = &userID

	return s.repo.Update(role)
}

func (s *roleService) DeleteRole(id int64) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("role not found")
	}
	return s.repo.Delete(id)
}
