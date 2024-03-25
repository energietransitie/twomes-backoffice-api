package services

import (
	"errors"

	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryupload"
	"github.com/energietransitie/twomes-backoffice-api/twomes/energyqueryvalue"
)

var (
	ErrEmptyQueryUpload = errors.New("no values in upload")
)

type EnergyQueryUploadService struct {
	repository energyqueryupload.EnergyQueryUploadRepository
}

// Create a new EnergyQueryUploadService.
func NewEnergyQueryUploadService(repository energyqueryupload.EnergyQueryUploadRepository) *EnergyQueryUploadService {
	return &EnergyQueryUploadService{
		repository: repository,
	}
}

func (s *EnergyQueryUploadService) Create(queryID uint, buildingID uint, energyQueryValues []energyqueryvalue.EnergyQueryValue) (energyqueryupload.EnergyQueryUpload, error) {
	if len(energyQueryValues) <= 0 {
		return energyqueryupload.EnergyQueryUpload{}, ErrEmptyUpload
	}

	queryUpload := energyqueryupload.MakeEnergyQueryUpload(queryID, buildingID, energyQueryValues)
	queryUpload, err := s.repository.Create(queryUpload)

	return queryUpload, err
}
