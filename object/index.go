package object

type Indexes interface{
	CreateIndex(namespace string) error
	RemoveIndex(namespace string) error
	Index(uri string, value interface{}) error
	IndexIfNotExist(uri string, value interface{}) error
	Query(namespace string, query interface{}, offset int64, size int64) (int64, []string, error)
}

