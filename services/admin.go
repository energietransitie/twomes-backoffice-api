package services

import (
	"time"

	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type AdminService struct {
	repository ports.AdminRepository

	// Services used when creating an admin.
	authService ports.AuthorizationService
}

// Create a new AdminService.
func NewAdminService(repository ports.AdminRepository, authService ports.AuthorizationService) *AdminService {
	return &AdminService{
		repository:  repository,
		authService: authService,
	}
}

func (s *AdminService) Create(name string, expiry time.Time) (twomes.Admin, error) {
	admin := twomes.MakeAdmin(name, expiry)
	admin, err := s.repository.Create(admin)
	if err != nil {
		return twomes.Admin{}, err
	}

	admin.AuthorizationToken, err = s.authService.CreateToken(twomes.AdminToken, admin.ID, admin.Expiry)
	if err != nil {
		return twomes.Admin{}, err
	}

	return admin, nil
}

func (s *AdminService) Find(admin twomes.Admin) (twomes.Admin, error) {
	return s.repository.Find(admin)
}

func (s *AdminService) GetAll() ([]twomes.Admin, error) {
	return s.repository.GetAll()
}

func (s *AdminService) Delete(admin twomes.Admin) error {
	return s.repository.Delete(admin)
}

func (s *AdminService) Reactivate(admin twomes.Admin) (twomes.Admin, error) {
	admin, err := s.repository.Find(admin)
	if err != nil {
		return twomes.Admin{}, err
	}

	admin.Reactivate()

	admin, err = s.repository.Update(admin)
	if err != nil {
		return twomes.Admin{}, err
	}

	admin.AuthorizationToken, err = s.authService.CreateToken(twomes.AdminToken, admin.ID, admin.Expiry)
	if err != nil {
		return twomes.Admin{}, err
	}

	return admin, nil
}

func (s *AdminService) SetExpiry(admin twomes.Admin, expiry time.Time) (twomes.Admin, error) {
	admin, err := s.repository.Find(admin)
	if err != nil {
		return twomes.Admin{}, err
	}

	admin.SetExpiry(expiry)
	return s.repository.Update(admin)
}
