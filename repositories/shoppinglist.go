package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglist"
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"
	"gorm.io/gorm"
)

type ShoppingListRepository struct {
	db *gorm.DB
}

func NewShoppingListRepository(db *gorm.DB) *ShoppingListRepository {
	return &ShoppingListRepository{
		db: db,
	}
}

// Database representation of a [shoppinglist.ShoppingList].
type ShoppingListModel struct {
	gorm.Model
	Items        []ShoppingListItemModel    `gorm:"foreignKey:ID"`
	Dependencies [][2]ShoppingListItemModel `gorm:"foreignKey:ID"`
}

// Set the name of the table in the database.
func (ShoppingListModel) TableName() string {
	return "shopping_list"
}

// Create a new ShoppingListModel from a [shoppinglist.ShoppingList]
func MakeShoppingListModel(shoppinglist shoppinglist.ShoppingList) ShoppingListModel {
	var shoppingListItemModels []ShoppingListItemModel

	for _, item := range shoppinglist.Items {
		shoppingListItemModels = append(shoppingListItemModels, MakeShoppingListItemModel(item))
	}

	var dependencies [][2]ShoppingListItemModel

	for _, dependency := range shoppinglist.Dependencies {
		var dependencyModels [2]ShoppingListItemModel
		for i, item := range dependency {
			dependencyModels[i] = MakeShoppingListItemModel(item)
		}
		dependencies = append(dependencies, dependencyModels)
	}

	return ShoppingListModel{
		Model:        gorm.Model{ID: shoppinglist.ID},
		Items:        shoppingListItemModels,
		Dependencies: dependencies,
	}
}

// Create a [shoppinglist.ShoppingList] from a ShoppingListModel.
func (m *ShoppingListModel) fromModel() shoppinglist.ShoppingList {
	var items []shoppinglistitem.ShoppingListItem

	for _, shoppingListItemModel := range m.Items {
		items = append(items, shoppingListItemModel.fromModel())
	}

	var dependencies [][2]shoppinglistitem.ShoppingListItem

	for _, dependency := range m.Dependencies {
		var dependencyModels [2]shoppinglistitem.ShoppingListItem
		for i, item := range dependency {
			dependencyModels[i] = item.fromModel()
		}
		dependencies = append(dependencies, dependencyModels)
	}

	return shoppinglist.ShoppingList{
		ID:           m.Model.ID,
		Items:        items,
		Dependencies: dependencies,
	}
}

func (r *ShoppingListRepository) Create(shoppinglist shoppinglist.ShoppingList) (shoppinglist.ShoppingList, error) {
	shoppingListModel := MakeShoppingListModel(shoppinglist)
	err := r.db.Create(&shoppingListModel).Error
	return shoppingListModel.fromModel(), err
}

func (r *ShoppingListRepository) Delete(shoppinglist shoppinglist.ShoppingList) error {
	shoppingListModel := MakeShoppingListModel(shoppinglist)
	return r.db.Create(&shoppingListModel).Error
}

func (r *ShoppingListRepository) Find(shoppingList shoppinglist.ShoppingList) (shoppinglist.ShoppingList, error) {
	shoppingListModel := MakeShoppingListModel(shoppingList)
	err := r.db.Where(&shoppingListModel).First(&shoppingListModel).Error
	return shoppingListModel.fromModel(), err
}

func (r *ShoppingListRepository) GetAll() ([]shoppinglist.ShoppingList, error) {
	var shoppingLists []shoppinglist.ShoppingList

	var shoppingListsModels []ShoppingListModel
	err := r.db.Find(&shoppingListsModels).Error
	if err != nil {
		return nil, err
	}

	for _, shoppingListModel := range shoppingListsModels {
		shoppingLists = append(shoppingLists, shoppingListModel.fromModel())
	}

	return shoppingLists, nil
}
