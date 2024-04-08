package services

import (
	"fmt"

	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcetype"
)

type DataSourceTypeService struct {
	repository datasourcetype.DataSourceTypeRepository

	//Service for setting item types
	deviceTypeService      *DeviceTypeService
	cloudFeedTypeService   *CloudFeedTypeService
	energyQueryTypeService *EnergyQueryTypeService
}

// Create a new DataSourceTypeService.
func NewDataSourceTypeService(
	repository datasourcetype.DataSourceTypeRepository,
	deviceTypeService *DeviceTypeService,
	cloudFeedTypeService *CloudFeedTypeService,
	energyQueryTypeService *EnergyQueryTypeService,
) *DataSourceTypeService {
	return &DataSourceTypeService{
		repository:             repository,
		deviceTypeService:      deviceTypeService,
		cloudFeedTypeService:   cloudFeedTypeService,
		energyQueryTypeService: energyQueryTypeService,
	}
}

// Used so we do not have to hardcode the check as much
type Source interface {
	GetByIDForDataSourceType(id uint) (interface{}, error)
	GetTableName() string
}

func (s *DataSourceTypeService) Create(
	typeSourceID uint,
	itemType datasourcetype.CategoryType,
	precedes []datasourcetype.DataSourceType,
	installationManualUrl string,
	infoUrl string,
	uploadSchedule []string,
	measurementSchedule []string,
	notificationThreshold string,
) (datasourcetype.DataSourceType, error) {

	//Ensures that the source associated with a given sourceID matches the expected item type
	source, err := s.GetSourceByID(typeSourceID)
	if err != nil {
		return datasourcetype.DataSourceType{}, fmt.Errorf("error retrieving source: %w", err)
	}

	if source.GetTableName() != string(itemType) {
		return datasourcetype.DataSourceType{}, fmt.Errorf("sourceID %s does not match itemType %s", source.GetTableName(), itemType)
	}
	//

	dataSourceType := datasourcetype.MakeDataSourceType(
		typeSourceID,
		itemType,
		installationManualUrl,
		infoUrl,
		precedes,
		uploadSchedule,
		measurementSchedule,
		notificationThreshold,
	)

	return s.repository.Create(dataSourceType)
}

func (s *DataSourceTypeService) Find(dataSourceType datasourcetype.DataSourceType) (datasourcetype.DataSourceType, error) {
	return s.repository.Find(dataSourceType)
}

func (s *DataSourceTypeService) GetAll() ([]datasourcetype.DataSourceType, error) {
	return s.repository.GetAll()
}

func (s *DataSourceTypeService) Delete(dataSourceType datasourcetype.DataSourceType) error {
	return s.repository.Delete(dataSourceType)
}

func (s *DataSourceTypeService) GetSourceByID(sourceID uint) (Source, error) {
	sources := []Source{
		s.deviceTypeService,
		s.cloudFeedTypeService,
		s.energyQueryTypeService,
	}

	for _, src := range sources {
		_, err := src.GetByIDForDataSourceType(sourceID)
		if err == nil {
			return src, nil
		}
	}

	return nil, fmt.Errorf("sourceID not found")
}

func (s *DataSourceTypeService) GetSourceByIDAndTable(sourceID uint, table string) (interface{}, error) {
	sources := []Source{
		s.deviceTypeService,
		s.cloudFeedTypeService,
		s.energyQueryTypeService,
	}

	var selectedSource Source

	switch table {
	case "device_type":
		selectedSource = sources[0]
	case "cloud_feed_type":
		selectedSource = sources[1]
	case "energy_query_type":
		selectedSource = sources[2]
	default:
		return nil, fmt.Errorf("unsupported table type: %s", table)
	}

	item, err := selectedSource.GetByIDForDataSourceType(sourceID)
	if err == nil {
		return item, nil
	}

	return nil, fmt.Errorf("sourceID not found")
}
