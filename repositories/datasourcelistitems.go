package repositories

// Custom Many2Many table to support duplicates and list order
//Yes, it is unfortunate that we have to do it manual since Gorm will throw a tantrum and make duplicate entries among other things
type DataSourceListItems struct {
	ID                    uint
	DataSourceListModelID uint
	DataSourceTypeModelID uint
	Order                 uint
}

func (DataSourceListItems) TableName() string {
	return "data_source_list_items"
}
