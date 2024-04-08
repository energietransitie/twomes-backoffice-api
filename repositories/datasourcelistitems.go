package repositories

// Custom Many2Many table to support duplicates and list order
type DataSourceListItems struct {
	DataSourceListID uint `gorm:"primaryKey"`
	DataSourceTypeID uint `gorm:"primaryKey"`
	Order            uint `gorm:"primaryKey"`
}

func (DataSourceListItems) TableName() string {
	return "data_source_list_items"
}
