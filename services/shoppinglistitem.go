package services

import (
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

func (s *ShoppingListItemService) Create(sourceID uint, schedule []string, itemType shoppinglistitemtype.ShoppingListItemType) (shoppinglistitem.ShoppingListItem, error) {
	listitemType, err := s.shoppingListItemTypeService.Find(itemType)
	if err != nil {
		return shoppinglistitem.ShoppingListItem{}, err
	}
	shoppingListItem := shoppinglistitem.MakeShoppingListItem(sourceID, schedule, listitemType)
	return shoppingListItem, nil
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
