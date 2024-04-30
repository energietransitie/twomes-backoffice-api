package datasourcetype

type DataSourceTypeRepository interface {
	Find(dataSourceType DataSourceType) (DataSourceType, error)
	GetAll() ([]DataSourceType, error)
	Create(DataSourceType) (DataSourceType, error)
	Delete(DataSourceType) error
}
