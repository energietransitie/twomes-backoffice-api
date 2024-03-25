package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryproperty"
)

type EnergyQueryPropertyService struct {
	repository energyqueryproperty.EnergyQueryPropertyRepository
}

// Create a new EnergyQueryPropertyService.
func NewEnergyQueryPropertyService(repository energyqueryproperty.EnergyQueryPropertyRepository) *EnergyQueryPropertyService {
	return &EnergyQueryPropertyService{
		repository: repository,
	}
}

func (s *EnergyQueryPropertyService) Create(name string, unit string) (energyqueryproperty.EnergyQueryProperty, error) {
	queryProperty := energyqueryproperty.MakeEnergyQueryProperty(name, unit)
	return s.repository.Create(queryProperty)
}

func (s *EnergyQueryPropertyService) Find(queryProperty energyqueryproperty.EnergyQueryProperty) (energyqueryproperty.EnergyQueryProperty, error) {
	return s.repository.Find(queryProperty)
}

func (s *EnergyQueryPropertyService) GetByID(id uint) (energyqueryproperty.EnergyQueryProperty, error) {
	return s.repository.Find(energyqueryproperty.EnergyQueryProperty{ID: id})
}

func (s *EnergyQueryPropertyService) GetByName(name string) (energyqueryproperty.EnergyQueryProperty, error) {
	return s.repository.Find(energyqueryproperty.EnergyQueryProperty{Name: name})
}
