package services

import (
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type AppService struct {
	repository ports.AppRepository
}

// Create a new AppService.
func NewAppService(repository ports.AppRepository) *AppService {
	return &AppService{
		repository: repository,
	}
}

// Create a new app.
func (s *AppService) Create(name string, provisioningURLTemplate string) (twomes.App, error) {
	app := twomes.MakeApp(name, provisioningURLTemplate)
	return s.repository.Create(app)
}

func (s *AppService) Find(app twomes.App) (twomes.App, error) {
	return s.repository.Find(app)
}

func (s *AppService) GetAll() ([]twomes.App, error) {
	return s.repository.GetAll()
}

// Get an app by its ID.
func (s *AppService) GetByID(id uint) (twomes.App, error) {
	return s.repository.Find(twomes.App{ID: id})
}
