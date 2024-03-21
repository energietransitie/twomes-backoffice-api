package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"
)

type ShoppingListItemTypeService struct {
	repository shoppinglistitemtype.ShoppingListItemTypeRepository
}

// Create a new ShoppingListItemTypeService.
func NewShoppingListItemTypeService(repository shoppinglistitemtype.ShoppingListItemTypeRepository) *ShoppingListItemTypeService {
	return &ShoppingListItemTypeService{
		repository: repository,
	}
}

func (s *ShoppingListItemTypeService) Create(tableName string) (shoppinglistitemtype.ShoppingListItemType, error) {
	shoppingListItemType := shoppinglistitemtype.MakeShoppingListItemType(tableName)
	return s.repository.Create(shoppingListItemType)
}

func (s *ShoppingListItemTypeService) Find(shoppingListItemType shoppinglistitemtype.ShoppingListItemType) (shoppinglistitemtype.ShoppingListItemType, error) {
	return s.repository.Find(shoppingListItemType)
}

func (s *ShoppingListItemTypeService) GetAll() ([]shoppinglistitemtype.ShoppingListItemType, error) {
	return s.repository.GetAll()
}

func (s *ShoppingListItemTypeService) Delete(shoppingListItemType shoppinglistitemtype.ShoppingListItemType) error {
	return s.repository.Delete(shoppingListItemType)
}
