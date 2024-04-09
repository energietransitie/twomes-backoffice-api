package repositories

import (
	"errors"
	"fmt"

	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcelist"
	"github.com/energietransitie/twomes-backoffice-api/twomes/datasourcetype"
	"gorm.io/gorm"
)

type DataSourceListRepository struct {
	db *gorm.DB
}

func NewDataSourceListRepository(db *gorm.DB) *DataSourceListRepository {
	return &DataSourceListRepository{
		db: db,
	}
}

// Database representation of a [datasourcelist.DataSourceList].
type DataSourceListModel struct {
	gorm.Model
	Items    []DataSourceTypeModel `gorm:"many2many:data_source_list_items;"`
	Campaign []CampaignModel       `gorm:"foreignKey:DataSourceListID"`
	Name     string
}

// Set the name of the table in the database.
func (DataSourceListModel) TableName() string {
	return "data_source_list"
}

// Create a new DataSourceListModel from a [datasourcelist.DataSourceList]
func MakeDataSourceListModel(dataSourceList datasourcelist.DataSourceList) DataSourceListModel {
	var dataSourceTypeModels []DataSourceTypeModel
	for _, item := range dataSourceList.Items {
		dataSourceTypeModels = append(dataSourceTypeModels, MakeDataSourceTypeModel(item))
	}

	return DataSourceListModel{
		Model: gorm.Model{ID: dataSourceList.ID},
		Name:  dataSourceList.Name,
		Items: dataSourceTypeModels,
	}
}

// Create a [datasourcelist.DataSourceList] from a DataSourceListModel.
func (m *DataSourceListModel) fromModel() datasourcelist.DataSourceList {
	var items []datasourcetype.DataSourceType

	for _, datasourceListItemModel := range m.Items {
		items = append(items, datasourceListItemModel.fromModel())
	}

	return datasourcelist.DataSourceList{
		ID:    m.Model.ID,
		Name:  m.Name,
		Items: items,
	}
}

func (r *DataSourceListRepository) Create(dataSourceList datasourcelist.DataSourceList) (datasourcelist.DataSourceList, error) {
	dataSourceListModel := MakeDataSourceListModel(dataSourceList)

	// Check for duplicate orders
	orderMap := make(map[uint]bool)
	orderMap[0] = true
	for _, item := range dataSourceList.Items {
		if orderMap[item.Order] && item.Order != 0 {
			return datasourcelist.DataSourceList{}, fmt.Errorf("duplicate order found: %d", item.Order)
		}
		orderMap[item.Order] = true
	}

	tx := r.db.Begin()
	if err := tx.Create(&dataSourceListModel).Error; err != nil {
		tx.Rollback()
		return datasourcelist.DataSourceList{}, fmt.Errorf("failed to create DataSourceListModel: %w", err)
	}

	// Update order in join table
	for _, item := range dataSourceList.Items {
		// Find existing DataSourceTypeModel by ID
		var dataSourceTypeModel DataSourceTypeModel
		if err := tx.First(&dataSourceTypeModel, item.ID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return datasourcelist.DataSourceList{}, fmt.Errorf("DataSourceTypeModel with ID %d not found: %w", item.ID, err)
			}
			tx.Rollback()
			return datasourcelist.DataSourceList{}, fmt.Errorf("failed to find DataSourceTypeModel: %w", err)
		}

		orderNumber, _ := findMaxKey(orderMap)
		//Autoincrement if order is set to 0
		if item.Order == 0 {
			orderNumber = orderNumber + 1
			orderMap[orderNumber] = true
			item.Order = orderNumber
		} else {
			orderNumber = item.Order
		}

		//Update order
		var existingDataSourceListItem DataSourceListItems
		if err := tx.Where("data_source_list_model_id = ? AND data_source_type_model_id = ?", dataSourceListModel.ID, dataSourceTypeModel.ID).First(&existingDataSourceListItem).Update("order", orderNumber).Error; err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				tx.Rollback()
				return datasourcelist.DataSourceList{}, fmt.Errorf("failed to update existing DataSourceListItems: %w", err)
			}
		}

		// Update order directly in dataSourceListModel.Items so we can return it in the response
		for idx := range dataSourceListModel.Items {
			if dataSourceListModel.Items[idx].ID == dataSourceTypeModel.ID {
				dataSourceListModel.Items[idx].Order = orderNumber
				break
			}
		}
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return datasourcelist.DataSourceList{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return dataSourceListModel.fromModel(), nil
}

func (r *DataSourceListRepository) Delete(datasourcelist datasourcelist.DataSourceList) error {
	datasourceListModel := MakeDataSourceListModel(datasourcelist)
	return r.db.Create(&datasourceListModel).Error
}

func (r *DataSourceListRepository) Find(datasourceList datasourcelist.DataSourceList) (datasourcelist.DataSourceList, error) {
	datasourceListModel := MakeDataSourceListModel(datasourceList)
	err := r.db.Where(&datasourceListModel).First(&datasourceListModel).Error
	return datasourceListModel.fromModel(), err
}

func (r *DataSourceListRepository) GetAll() ([]datasourcelist.DataSourceList, error) {
	var datasourceLists []datasourcelist.DataSourceList

	var datasourceListsModels []DataSourceListModel
	err := r.db.Find(&datasourceListsModels).Error
	if err != nil {
		return nil, err
	}

	for _, datasourceListModel := range datasourceListsModels {
		datasourceLists = append(datasourceLists, datasourceListModel.fromModel())
	}

	return datasourceLists, nil
}

func findMaxKey(orderMap map[uint]bool) (uint, error) {
	if len(orderMap) == 0 {
		return 0, fmt.Errorf("map is empty")
	}

	var maxKey uint

	for key := range orderMap {
		if key > maxKey {
			maxKey = key
		}
	}

	return maxKey, nil
}
