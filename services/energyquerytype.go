package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquerytype"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvariety"
)

type EnergyQueryTypeService struct {
	repository energyquerytype.EnergyQueryTypeRepository
}

// Create a new EnergyQueryTypeService.
func NewEnergyQueryTypeService(repository energyquerytype.EnergyQueryTypeRepository) *EnergyQueryTypeService {
	return &EnergyQueryTypeService{
		repository: repository,
	}
}

func (s *EnergyQueryTypeService) Create(queryVariety energyqueryvariety.EnergyQueryVariety, formula string) (energyquerytype.EnergyQueryType, error) {
	EnergyQueryType := energyquerytype.MakeEnergyQueryType(queryVariety, formula)
	return s.repository.Create(EnergyQueryType)
}

func (s *EnergyQueryTypeService) Find(energyQueryType energyquerytype.EnergyQueryType) (energyquerytype.EnergyQueryType, error) {
	return s.repository.Find(energyQueryType)
}

func (s *EnergyQueryTypeService) FindById(id uint) (energyquerytype.EnergyQueryType, error) {
	return s.repository.Find(energyquerytype.EnergyQueryType{ID: id})
}

func (s *EnergyQueryTypeService) GetAll() ([]energyquerytype.EnergyQueryType, error) {
	return s.repository.GetAll()
}

func (s *EnergyQueryTypeService) Delete(energyQueryType energyquerytype.EnergyQueryType) error {
	return s.repository.Delete(energyQueryType)
}

func (s *EnergyQueryTypeService) GetByIDForDataSourceType(id uint) (interface{}, error) {
	return s.repository.Find(energyquerytype.EnergyQueryType{ID: id})
}

func (s *EnergyQueryTypeService) GetTableName() string {
	return "energy_query_type"
}
