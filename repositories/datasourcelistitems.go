package repositories

// Custom Many2Many table to support duplicates and list order
type DataSourceListItems struct {
	ID                    uint
	DataSourceListModelID uint
	DataSourceTypeModelID uint
	Order                 uint
}

func (DataSourceListItems) TableName() string {
	return "data_source_list_items"
}
