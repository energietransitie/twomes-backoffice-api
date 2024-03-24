package repositories

import (
	"errors"
	"strings"

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
	SourceID              uint
	TypeID                uint
	Type                  ShoppingListItemTypeModel `gorm:"foreignkey:TypeID"`
	Precedes              []ShoppingListItemModel   `gorm:"many2many:Shopping_List_Item_Precedes;"`
	UploadSchedule        string                    `gorm:"type:text"` //It will crash if we do not specify text here!
	MeasurementSchedule   string                    `gorm:"type:text"`
	NotificationThreshold string
}

// Set the name of the table in the database.
func (ShoppingListItemModel) TableName() string {
	return "shopping_list_item"
}

// Create a new ShoppingListItemModel from a [shoppinglistitem.ShoppinglistItem]
func MakeShoppingListItemModel(shoppinglistitem shoppinglistitem.ShoppingListItem) ShoppingListItemModel {
	var shoppingListItemModels []ShoppingListItemModel
	for _, item := range shoppinglistitem.Precedes {
		shoppingListItemModels = append(shoppingListItemModels, MakeShoppingListItemModel(item))
	}

	return ShoppingListItemModel{
		Model:                 gorm.Model{ID: shoppinglistitem.ID},
		SourceID:              shoppinglistitem.SourceID,
		TypeID:                shoppinglistitem.Type.ID,
		Type:                  MakeShoppingListItemTypeModel(shoppinglistitem.Type),
		Precedes:              shoppingListItemModels,
		UploadSchedule:        ConvertScheduleToModel(shoppinglistitem.UploadSchedule),
		MeasurementSchedule:   ConvertScheduleToModel(shoppinglistitem.UploadSchedule),
		NotificationThreshold: shoppinglistitem.NotificationThreshold,
	}
}

// Create a [shoppinglistitem.ShoppingListItem] from a ShoppingListItemModel
func (m *ShoppingListItemModel) fromModel() shoppinglistitem.ShoppingListItem {
	var items []shoppinglistitem.ShoppingListItem
	for _, shoppingListItemModel := range m.Precedes {
		items = append(items, shoppingListItemModel.fromModel())
	}

	return shoppinglistitem.ShoppingListItem{
		ID:                    m.Model.ID,
		SourceID:              m.SourceID,
		Type:                  m.Type.fromModel(),
		Precedes:              items,
		UploadSchedule:        ConvertScheduleToArray(m.UploadSchedule),
		MeasurementSchedule:   ConvertScheduleToArray(m.MeasurementSchedule),
		NotificationThreshold: m.NotificationThreshold,
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

func ConvertScheduleToModel(jsonArray []string) string {
	var convertedString string

	for _, item := range jsonArray {
		convertedString += item + ","
	}
	convertedString = strings.TrimSuffix(convertedString, ",")

	return convertedString
}

func ConvertScheduleToArray(modelString string) []string {
	convertedArray := strings.Split(modelString, ",")
	return convertedArray
}

// Check if we did not make a loop that can softlock the app
func (s *ShoppingListItemModel) AfterSave(tx *gorm.DB) (err error) {
	var emptySlice []uint
	if s.CheckforCircular(s, emptySlice) {
		if err := tx.Rollback().Error; err != nil {
			return err
		}
		return errors.New("circular reference detected, transaction rolled back")
	}
	return nil
}

func (s *ShoppingListItemModel) CheckforCircular(item *ShoppingListItemModel, previousIDs []uint) bool {
	previousIDs = append(previousIDs, item.ID)
	for _, elem := range item.Precedes {
		for _, ID := range previousIDs {
			if elem.ID == ID || s.CheckforCircular(&elem, previousIDs) {
				return true
			}
		}
	}
	return false
}
