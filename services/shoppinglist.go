package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglist"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"
)

type ShoppingListService struct {
	repository shoppinglist.ShoppingListRepository

	// Service used for the items
	shoppingListItemService *ShoppingListItemService
}

// Create a new ShoppingListService.
func NewShoppingListService(repository shoppinglist.ShoppingListRepository, shoppingListItemService *ShoppingListItemService) *ShoppingListService {
	return &ShoppingListService{
		repository:              repository,
		shoppingListItemService: shoppingListItemService,
	}
}

func (s *ShoppingListService) Create(description string, items []shoppinglistitem.ShoppingListItem) (shoppinglist.ShoppingList, error) {
	for i, item := range items {
		listItem, err := s.shoppingListItemService.Find(item)
		if err != nil {
			return shoppinglist.ShoppingList{}, err
		}
		items[i] = listItem
	}

	shoppinglist := shoppinglist.MakeShoppingList(items, description)
	return s.repository.Create(shoppinglist)
}

func (s *ShoppingListService) Find(shoppingList shoppinglist.ShoppingList) (shoppinglist.ShoppingList, error) {
	return s.repository.Find(shoppingList)
}

func (s *ShoppingListService) GetAll() ([]shoppinglist.ShoppingList, error) {
	return s.repository.GetAll()
}

func (s *ShoppingListService) Delete(shoppingList shoppinglist.ShoppingList) error {
	return s.repository.Delete(shoppingList)
}
