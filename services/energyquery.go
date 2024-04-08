package services

import (
	"errors"

	"github.com/energietransitie/twomes-backoffice-api/twomes/account"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquery"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyquerytype"
	"github.com/energietransitie/twomes-backoffice-api/twomes/upload"
)

var (
	ErrQueryDoesNotBelongToAccount = errors.New("device does not belong to this account")
)

type EnergyQueryService struct {
	repository energyquery.EnergyQueryRepository

	// Services used when creating an energy query.
	accountService         *AccountService
	energyQueryTypeService *EnergyQueryTypeService

	// Services used when getting device info.
	uploadService *UploadService
}

// Create a new EnergyQueryService.
func NewEnergyQueryService(
	repository energyquery.EnergyQueryRepository,
	accountService *AccountService,
	energyQueryTypeService *EnergyQueryTypeService,
	uploadService *UploadService,
) *EnergyQueryService {
	return &EnergyQueryService{
		repository:             repository,
		accountService:         accountService,
		energyQueryTypeService: energyQueryTypeService,
		uploadService:          uploadService,
	}
}

func (s *EnergyQueryService) Create(queryType energyquerytype.EnergyQueryType, account account.Account, uploads []upload.Upload) (energyquery.EnergyQuery, error) {
	energyQuery := energyquery.MakeEnergyQuery(queryType, account, uploads)
	return s.repository.Create(energyQuery)
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
