package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquery"
)

type EnergyQueryService struct {
	repository energyquery.EnergyQueryRepository
}

// Create a new EnergyQueryService.
func NewEnergyQueryService(repository energyquery.EnergyQueryRepository) *EnergyQueryService {
	return &EnergyQueryService{
		repository: repository,
	}
}

func (s *EnergyQueryService) Create(name string, formula string) (energyquery.EnergyQuery, error) {
	EnergyQuery := energyquery.MakeEnergyQuery(name, formula)
	return s.repository.Create(EnergyQuery)
}

func (s *EnergyQueryService) Find(energyQuery energyquery.EnergyQuery) (energyquery.EnergyQuery, error) {
	return s.repository.Find(energyQuery)
}

func (s *EnergyQueryService) FindById(id uint) (energyquery.EnergyQuery, error) {
	return s.repository.Find(energyquery.EnergyQuery{ID: id})
}

func (s *EnergyQueryService) GetAll() ([]energyquery.EnergyQuery, error) {
	return s.repository.GetAll()
}

func (s *EnergyQueryService) Delete(energyQuery energyquery.EnergyQuery) error {
	return s.repository.Delete(energyQuery)
}

func (s *EnergyQueryService) GetByIDForShoppingList(id uint) (interface{}, error) {
	return s.repository.Find(energyquery.EnergyQuery{ID: id})
}

func (s *EnergyQueryService) GetTableName() string {
	return "energy_query"
}
