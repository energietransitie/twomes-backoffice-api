package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitemtype"
	"gorm.io/gorm"
)

type ShoppingListItemTypeRepository struct {
	db *gorm.DB
}

func NewShoppingListItemTypeRepository(db *gorm.DB) *ShoppingListItemTypeRepository {
	return &ShoppingListItemTypeRepository{
		db: db,
	}
}

// Database representation of a [shoppinglistitemtype.ShoppingListItemType].
type ShoppingListItemTypeModel struct {
	gorm.Model
	Name string
}

// Set the name of the table in the database.
func (ShoppingListItemTypeModel) TableName() string {
	return "shoppinglist_item_type"
}

// Create a new ShoppingListItemModel from a [shoppinglistitemtype.ShoppinglistItemType]
func MakeShoppingListItemTypeModel(shoppinglistitemtype shoppinglistitemtype.ShoppingListItemType) ShoppingListItemTypeModel {
	return ShoppingListItemTypeModel{
		Model: gorm.Model{ID: shoppinglistitemtype.ID},
		Name:  shoppinglistitemtype.Name,
	}
}

// Create a [shoppinglistitemType.ShoppingListItemType] from a ShoppingListItemTypeModel
func (m *ShoppingListItemTypeModel) fromModel() shoppinglistitemtype.ShoppingListItemType {
	return shoppinglistitemtype.ShoppingListItemType{
		ID:   m.Model.ID,
		Name: m.Name,
	}
}

func (r *ShoppingListItemTypeRepository) Create(shoppinglistitemtype shoppinglistitemtype.ShoppingListItemType) (shoppinglistitemtype.ShoppingListItemType, error) {
	shoppingListItemTypeModel := MakeShoppingListItemTypeModel(shoppinglistitemtype)
	err := r.db.Create(&shoppingListItemTypeModel).Error
	return shoppingListItemTypeModel.fromModel(), err
}

func (r *ShoppingListItemTypeRepository) Delete(shoppinglistitemtype shoppinglistitemtype.ShoppingListItemType) error {
	shoppingListItemTypeModel := MakeShoppingListItemTypeModel(shoppinglistitemtype)
	return r.db.Create(&shoppingListItemTypeModel).Error
}
