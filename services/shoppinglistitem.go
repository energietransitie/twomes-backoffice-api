package services

import (
	"fmt"

	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"
)

type ShoppingListItemService struct {
	repository shoppinglistitem.ShoppingListItemRepository

	//Service for setting item types
	shoppingListItemTypeService *ShoppingListItemTypeService
	deviceTypeService           *DeviceTypeService
	cloudFeedService            *CloudFeedService
}

// Create a new ShoppingListItemService.
func NewShoppingListItemService(
	repository shoppinglistitem.ShoppingListItemRepository,
	shoppingListItemTypeService *ShoppingListItemTypeService,
	deviceTypeService *DeviceTypeService,
	cloudFeedService *CloudFeedService,
) *ShoppingListItemService {
	return &ShoppingListItemService{
		repository:                  repository,
		shoppingListItemTypeService: shoppingListItemTypeService,
		deviceTypeService:           deviceTypeService,
		cloudFeedService:            cloudFeedService,
	}
}

// Used so we do not have to hardcode the check as much
type Source interface {
	GetByIDForShoppingList(id uint) (interface{}, error)
	GetTableName() string
}

func (s *ShoppingListItemService) Create(
	sourceID uint,
	itemType shoppinglistitemtype.ShoppingListItemType,
	precedes []shoppinglistitem.ShoppingListItem,
	uploadSchedule []string,
	measurementSchedule []string,
	notificationThreshold string,
) (shoppinglistitem.ShoppingListItem, error) {

	//Ensures that the source associated with a given sourceID matches the expected item type
	foundType, err := s.shoppingListItemTypeService.Find(itemType)
	if err != nil {
		return shoppinglistitem.ShoppingListItem{}, err
	}

	source, err := s.GetSourceByID(sourceID)
	if err != nil {
		return shoppinglistitem.ShoppingListItem{}, fmt.Errorf("error retrieving source: %w", err)
	}

	if source.GetTableName() != foundType.Name {
		return shoppinglistitem.ShoppingListItem{}, fmt.Errorf("sourceID %s does not match itemType %s", source.GetTableName(), foundType.Name)
	}
	//

	shoppingListItem := shoppinglistitem.MakeShoppingListItem(sourceID, itemType, precedes, uploadSchedule, measurementSchedule, notificationThreshold)

	return s.repository.Create(shoppingListItem)
}

func (s *ShoppingListItemService) Find(shoppingListItem shoppinglistitem.ShoppingListItem) (shoppinglistitem.ShoppingListItem, error) {
	return s.repository.Find(shoppingListItem)
}

func (s *ShoppingListItemService) GetAll() ([]shoppinglistitem.ShoppingListItem, error) {
	return s.repository.GetAll()
}

func (s *ShoppingListItemService) Delete(shoppingListItem shoppinglistitem.ShoppingListItem) error {
	return s.repository.Delete(shoppingListItem)
}

func (s *ShoppingListItemService) GetSourceByID(sourceID uint) (Source, error) {
	sources := []Source{
		s.deviceTypeService,
		s.cloudFeedService,
		//&EnergyQueryService{},
	}

	for _, src := range sources {
		_, err := src.GetByIDForShoppingList(sourceID)
		if err == nil {
			return src, nil
		}
	}

	return nil, fmt.Errorf("sourceID not found")
}
