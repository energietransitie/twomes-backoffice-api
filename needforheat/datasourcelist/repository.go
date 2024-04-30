package datasourcelist

type DataSourceListRepository interface {
	Find(dataSourceList DataSourceList) (DataSourceList, error)
	GetAll() ([]DataSourceList, error)
	Create(DataSourceList) (DataSourceList, error)
	Delete(DataSourceList) error
}
