package repositories

// Custom Many2Many table to support duplicates and list order
type DataSourceListItems struct {
	DataSourceListModelID uint `gorm:"primaryKey;autoIncrement:false;"`
	DataSourceTypeModelID uint `gorm:"primaryKey;autoIncrement:false;"`
	Order                 uint `gorm:"primaryKey;autoIncrement:false;default:0;"`
}

func (DataSourceListItems) TableName() string {
	return "data_source_list_items"
}
