package repositories

import (
	"github.com/energietransitie/twomes-backoffice-api/twomes/shoppinglistitem"
	"gorm.io/gorm"
)

type ShoppingListItemRepository struct {
	db *gorm.DB
}

func NewShoppingListItemRepository(db *gorm.DB) *ShoppingListItemRepository {
	return &ShoppingListItemRepository{
		db: db,
	}
}

// Database representation of a [shoppinglistitem.ShoppingListItem].
type ShoppingListItemModel struct {
	gorm.Model
	SourceID uint                      `gorm:"polymorphic:Type;polymorphicValue:device,cloud_feed"`
	Type     ShoppingListItemTypeModel `gorm:"foreignKey:ID"`
	Schedule []string                  //Cronjob format
}

// Set the name of the table in the database.
func (ShoppingListItemModel) TableName() string {
	return "shoppinglist_item"
}

// Create a new ShoppingListItemModel from a [shoppinglistitem.ShoppinglistItem]
func MakeShoppingListItemModel(shoppinglistitem shoppinglistitem.ShoppingListItem) ShoppingListItemModel {
	return ShoppingListItemModel{
		Model:    gorm.Model{ID: shoppinglistitem.ID},
		SourceID: shoppinglistitem.SourceID,
		Type:     MakeShoppingListItemTypeModel(shoppinglistitem.Type),
		Schedule: shoppinglistitem.Schedule,
	}
}

// Create a [shoppinglistitem.ShoppingListItem] from a ShoppingListItemModel
func (m *ShoppingListItemModel) fromModel() shoppinglistitem.ShoppingListItem {
	return shoppinglistitem.ShoppingListItem{
		ID:       m.Model.ID,
		SourceID: m.SourceID,
		Type:     m.Type.fromModel(),
		Schedule: m.Schedule,
	}
}

func (r *ShoppingListItemRepository) Create(shoppinglistitem shoppinglistitem.ShoppingListItem) (shoppinglistitem.ShoppingListItem, error) {
	shoppingListItemModel := MakeShoppingListItemModel(shoppinglistitem)
	err := r.db.Create(&shoppingListItemModel).Error
	return shoppingListItemModel.fromModel(), err
}

func (r *ShoppingListItemRepository) Delete(shoppinglistitem shoppinglistitem.ShoppingListItem) error {
	shoppingListItemModel := MakeShoppingListItemModel(shoppinglistitem)
	return r.db.Create(&shoppingListItemModel).Error
}
