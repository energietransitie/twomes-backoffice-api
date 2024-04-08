package services

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcelist"
	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcetype"
)

type DataSourceListService struct {
	repository datasourcelist.DataSourceListRepository

	// Service used for the items
	shoppingListItemService *DataSourceTypeService
}

// Create a new DataSourceListService.
func NewDataSourceListService(repository datasourcelist.DataSourceListRepository, shoppingListItemService *DataSourceTypeService) *DataSourceListService {
	return &DataSourceListService{
		repository:              repository,
		shoppingListItemService: shoppingListItemService,
	}
}

func (s *DataSourceListService) Create(name string, items []datasourcetype.DataSourceType) (datasourcelist.DataSourceList, error) {
	for i, item := range items {
		listItem, err := s.shoppingListItemService.Find(item)
		if err != nil {
			return datasourcelist.DataSourceList{}, err
		}
		items[i] = listItem
	}

	datasourcelist := datasourcelist.MakeDataSourceList(items, name)
	return s.repository.Create(datasourcelist)
}

func (s *DataSourceListService) Find(shoppingList datasourcelist.DataSourceList) (datasourcelist.DataSourceList, error) {
	return s.repository.Find(shoppingList)
}

func (s *DataSourceListService) GetAll() ([]datasourcelist.DataSourceList, error) {
	return s.repository.GetAll()
}

func (s *DataSourceListService) Delete(shoppingList datasourcelist.DataSourceList) error {
	return s.repository.Delete(shoppingList)
}
