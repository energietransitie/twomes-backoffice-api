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
	SourceID uint                      `gorm:"polymorphic:Type;polymorphicValue:device_type,cloud_feed"`
	Schedule []string                  `gorm:"type:json"` //It will crash if we do not specify json here!
	Type     ShoppingListItemTypeModel `gorm:"foreignKey:ID"`
}

// Set the name of the table in the database.
func (ShoppingListItemModel) TableName() string {
	return "shopping_list_item"
}

// Create a new ShoppingListItemModel from a [shoppinglistitem.ShoppinglistItem]
func MakeShoppingListItemModel(shoppinglistitem shoppinglistitem.ShoppingListItem) ShoppingListItemModel {
	return ShoppingListItemModel{
		Model:    gorm.Model{ID: shoppinglistitem.ID},
		SourceID: shoppinglistitem.SourceID,
		Schedule: shoppinglistitem.Schedule,
		Type:     MakeShoppingListItemTypeModel(shoppinglistitem.Type),
	}
}

// Create a [shoppinglistitem.ShoppingListItem] from a ShoppingListItemModel
func (m *ShoppingListItemModel) fromModel() shoppinglistitem.ShoppingListItem {
	return shoppinglistitem.ShoppingListItem{
		ID:       m.Model.ID,
		SourceID: m.SourceID,
		Schedule: m.Schedule,
		Type:     m.Type.fromModel(),
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

func (r *ShoppingListItemRepository) Find(shoppingListItem shoppinglistitem.ShoppingListItem) (shoppinglistitem.ShoppingListItem, error) {
	shoppingListItemModel := MakeShoppingListItemModel(shoppingListItem)
	err := r.db.Where(&shoppingListItemModel).First(&shoppingListItemModel).Error
	return shoppingListItemModel.fromModel(), err
}

func (r *ShoppingListItemRepository) GetAll() ([]shoppinglistitem.ShoppingListItem, error) {
	var shoppingListItems []shoppinglistitem.ShoppingListItem

	var shoppingListItemModels []ShoppingListItemModel
	err := r.db.Find(&shoppingListItemModels).Error
	if err != nil {
		return nil, err
	}

	for _, shoppingListItemModel := range shoppingListItemModels {
		shoppingListItems = append(shoppingListItems, shoppingListItemModel.fromModel())
	}

	return shoppingListItems, nil
}
