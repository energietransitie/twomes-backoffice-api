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

func (s *ShoppingListItemService) Create(
	sourceID uint,
	itemType shoppinglistitemtype.ShoppingListItemType,
	precedes []shoppinglistitem.ShoppingListItem,
	uploadSchedule []string,
	measurementSchedule []string,
	notificationThreshold string,
) (shoppinglistitem.ShoppingListItem, error) {

	//Check if sourceID and itemType exists. SourceID can be deviceType or cloudfeed, itemType has a Name field with the table name.
	foundType, err := s.shoppingListItemTypeService.Find(itemType)
	if err != nil {
		return shoppinglistitem.ShoppingListItem{}, err
	}
	var sourceName string

	_, err = s.deviceTypeService.GetByID(sourceID)
	if err == nil {
		sourceName = "device_type"
	} else {
		_, err = s.cloudFeedService.GetByID(sourceID)
		if err != nil {
			return shoppinglistitem.ShoppingListItem{}, fmt.Errorf("sourceID not found")
		}
		sourceName = "cloud_feed"
	}

	if sourceName != foundType.Name {
		return shoppinglistitem.ShoppingListItem{}, fmt.Errorf("sourceID %s does not match itemType %s", sourceName, foundType.Name)
	}

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
