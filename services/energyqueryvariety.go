package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvariety"
)

type EnergyQueryVarietyService struct {
	repository energyqueryvariety.EnergyQueryVarietyRepository
}

// Create a new EnergyQueryVarietyService.
func NewEnergyQueryVarietyService(repository energyqueryvariety.EnergyQueryVarietyRepository) *EnergyQueryVarietyService {
	return &EnergyQueryVarietyService{
		repository: repository,
	}
}

func (s *EnergyQueryVarietyService) Create(name string) (energyqueryvariety.EnergyQueryVariety, error) {
	EnergyQueryVariety := energyqueryvariety.MakeEnergyQueryVariety(name)
	return s.repository.Create(EnergyQueryVariety)
}

func (s *EnergyQueryVarietyService) Find(energyQueryVariety energyqueryvariety.EnergyQueryVariety) (energyqueryvariety.EnergyQueryVariety, error) {
	return s.repository.Find(energyQueryVariety)
}

func (s *EnergyQueryVarietyService) FindById(id uint) (energyqueryvariety.EnergyQueryVariety, error) {
	return s.repository.Find(energyqueryvariety.EnergyQueryVariety{ID: id})
}

func (s *EnergyQueryVarietyService) GetAll() ([]energyqueryvariety.EnergyQueryVariety, error) {
	return s.repository.GetAll()
}

func (s *EnergyQueryVarietyService) Delete(energyQueryVariety energyqueryvariety.EnergyQueryVariety) error {
	return s.repository.Delete(energyQueryVariety)
}
