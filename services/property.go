package services

import (
	"github.com/energietransitie/twomes-backoffice-api/ports"
	"github.com/energietransitie/twomes-backoffice-api/twomes"
)

type PropertyService struct {
	repository ports.PropertyRepository
}

// Create a new PropertyService.
func NewPropertyService(repository ports.PropertyRepository) *PropertyService {
	return &PropertyService{
		repository: repository,
	}
}

func (s *PropertyService) Create(name string) (twomes.Property, error) {
	property := twomes.MakeProperty(name)
	return s.repository.Create(property)
}

func (s *PropertyService) Find(property twomes.Property) (twomes.Property, error) {
	return s.repository.Find(property)
}

func (s *PropertyService) GetByID(id uint) (twomes.Property, error) {
	return s.repository.Find(twomes.Property{ID: id})
}

func (s *PropertyService) GetByName(name string) (twomes.Property, error) {
	return s.repository.Find(twomes.Property{Name: name})
}
