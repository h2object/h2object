package object

type Store interface{
	PutOBJECT(obj OBJECT) error
	GetOBJECT(uri string) (OBJECT, error)
	DelOBJECT(obj OBJECT) error

	Put(uri string, val interface{}) error
	Get(uri string, nested bool) (interface{}, error)
	Del(uri string) error

	MultiGet(uri string, suffix string, nested bool) ([]interface{}, error)
	
	Size(uri string) (int64, error)
}

