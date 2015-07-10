package object

import (
	"fmt"
	"time"
	"path"
	"errors"
	"strings"
	"reflect"
	"strconv"
	"github.com/boltdb/bolt"
	"github.com/h2object/cast"
	"github.com/h2object/h2object/util"
)

type BoltStore struct{
	db *bolt.DB
	root  string
	name  string
	coder ByteCoder
}

func NewBoltStore(root string, name string, coder ByteCoder) *BoltStore {
	return &BoltStore{
		root: root,
		name: name,
		coder: coder,
	}
}

func (store *BoltStore) Load() error {
	if store.db == nil {
		db, err := bolt.Open(path.Join(store.root, store.name), 0600, 
			&bolt.Options{Timeout: 1 * time.Second})
		if err != nil {
			return err
		}
		store.db = db
	}
	return nil
}

func (store *BoltStore) PutOBJECT(object OBJECT) error {
	if object == nil {
		return nil
	}
	return store.Put(object.URI(), object.Value())
}

func (store *BoltStore) GetOBJECT(uri string) (OBJECT, error) {
	value, err := store.Get(uri, true)
	if err != nil {
		return nil, err
	}
	return CreateOBJECTByValue(uri, value)		
}

func (store *BoltStore) DelOBJECT(object OBJECT) error {
	if object == nil {
		return nil
	}
	return store.Del(object.URI())
}


func (store *BoltStore) Put(uri string, value interface{}) error {
	namespace, key := util.NamespaceKey(uri)
	bucketkeys := util.PathFolders(namespace)
	var bucket *bolt.Bucket = nil
	var err error
	e := store.db.Update(func(tx *bolt.Tx) error {		
		for _, bucketkey := range bucketkeys {
			if bucket == nil {
				bucket, err = tx.CreateBucketIfNotExists([]byte(bucketkey))
				if err != nil {
					return err
				}
			} else {
				bucket, err = bucket.CreateBucketIfNotExists([]byte(bucketkey))
				if err != nil {
					return err
				}
			}
		}
		return store.put(bucket, key, value)
	})
	return e
}
func (store *BoltStore) Get(uri string, nested bool) (interface{}, error) {
	namespace, key := util.NamespaceKey(uri) 
	bucketkeys := util.PathFolders(namespace)
	var bucket *bolt.Bucket = nil
	var result interface{}
	e := store.db.View(func(tx *bolt.Tx) error {	
		for _, bucketkey := range bucketkeys {
			if bucket == nil {
				bucket = tx.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			} else {
				bucket = bucket.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			}
		}
		if bucket == nil {
			return errors.New(uri + " not exists")
		}

		val, err := store.get(bucket, key, nested)
		if err != nil {
			return err
		}
		result = val
		return nil
	})	
	return result, e
}

func (store *BoltStore) MultiGet(uri string, suffix string, nested bool) ([]interface{}, error) {
	results := []interface{}{}
	bucketkeys := util.PathFolders(uri)
	var bucket *bolt.Bucket = nil
	e := store.db.View(func(tx *bolt.Tx) error {	
		for _, bucketkey := range bucketkeys {
			if bucket == nil {
				bucket = tx.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			} else {
				bucket = bucket.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			}
		}
		if bucket == nil {
			return errors.New(uri + " not exists")
		}

		vals, err := store.multi_get(bucket, suffix, nested)
		if err != nil {
			return err
		}
		results = append(results, vals...)	
		return nil
	})	
	return results, e
}

func (store *BoltStore) Del(uri string) error {
	namespace, key := util.NamespaceKey(uri)
	bucketkeys := util.PathFolders(namespace)
	var bucket *bolt.Bucket = nil
	var err error
	e := store.db.Update(func(tx *bolt.Tx) error {		
		for _, bucketkey := range bucketkeys {
			if bucket == nil {
				bucket, err = tx.CreateBucketIfNotExists([]byte(bucketkey))
				if err != nil {
					return err
				}
			} else {
				bucket, err = bucket.CreateBucketIfNotExists([]byte(bucketkey))
				if err != nil {
					return err
				}
			}
		}
		bucket.DeleteBucket([]byte(key))
		bucket.Delete([]byte(key))
		return nil
	})
	return e
}

func (store *BoltStore) Size(uri string) (int64, error) {
	namespace, key := util.NamespaceKey(uri) 
	bucketkeys := util.PathFolders(namespace)
	var bucket *bolt.Bucket = nil
	var result int64 = 0
	e := store.db.View(func(tx *bolt.Tx) error {	
		for _, bucketkey := range bucketkeys {
			if bucket == nil {
				bucket = tx.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			} else {
				bucket = bucket.Bucket([]byte(bucketkey))
				if bucket == nil {
					return errors.New(uri + " not exists")
				}
			}
		}
		if bucket == nil {
			return errors.New(uri + " not exists")
		}

		result = store.size(bucket, key, nil)
		return nil
	})	
	return result, e
}

func (store *BoltStore) size(bucket *bolt.Bucket, key string, size *AtomicInt64) int64 {
	var statics *AtomicInt64 = size
	if statics == nil {
		statics = NewAtomicInt64(0)	
	}

	if key != "" {
		b := bucket.Get([]byte(key))
		if b != nil {
			return statics.Caculate(int64(len(key) + len(b)))
		}

		bucket = bucket.Bucket([]byte(key))
	}

	if bucket != nil {
		c := bucket.Cursor()
		// to do staticsing
		for k, v := c.First(); k != nil; k, v = c.Next() {			
			if v != nil {
				statics.Caculate(int64(len(k) + len(v)))
			} else {				
				store.size(bucket, string(k), statics)
			}
		}
		statics.Caculate(int64(len(key)))
		return statics.Value()
	}

	//! null 
	return statics.Value()
}

func (store *BoltStore) put(bucket *bolt.Bucket, key string, value interface{}) error {
	var err error
	if is_value_basic(value) {
		if key == "" {
			return fmt.Errorf("key is absent")		
		}
		byts, err := store.coder.Encode(value)
		if err != nil {
			return err
		}
		return bucket.Put([]byte(key), byts)
	}

	if key != "" {
		bucket, err = bucket.CreateBucketIfNotExists([]byte(key))
		if err != nil {
			return err
		}	
	}	

	switch value_type_kind(value) {
	case reflect.Slice:
		vs := reflect.ValueOf(value)
		for i := 0; i <  vs.Len(); i++ {
			if err := store.put(bucket, strconv.Itoa(i), vs.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		vs, err := cast.ToStringMapE(value)
		if err != nil {
			return err
		}
		for k, v := range vs {
			if err := store.put(bucket, k, v); err != nil {
				return err
			}
		}
		return nil
	case reflect.Struct:
		vt := reflect.TypeOf(value)
		vv := reflect.ValueOf(value)
		for i := 0; i < vt.NumField(); i ++ {
			vfield := vt.Field(i)
			if tag := vfield.Tag.Get("object"); tag != "" {
				if err := store.put(bucket, tag, vv.Field(i).Interface()); err != nil {
					return err
				}
			}
		}	
		return nil
	}
	return fmt.Errorf("unsupport store value type")
}

func (store *BoltStore) get(bucket *bolt.Bucket, key string, nested bool) (interface{}, error) {
	if key != "" {
		b := bucket.Get([]byte(key))
		if b != nil {
			return store.coder.Decode(b)
		}
		bucket = bucket.Bucket([]byte(key))
	}

	if bucket != nil {
		m := make(map[string]interface{})
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if v != nil {
				vv, err := store.coder.Decode(v)
				if err != nil {
					return nil, err
				}
				m[string(k)] = vv
			} 

			if v == nil && nested == true {
				vv, err := store.get(bucket, string(k), nested)
				if err != nil {
					return nil, err
				}
				m[string(k)] = vv
			}
		}
		return m, nil
	}

	return nil, fmt.Errorf("bucket not exist")
}

func (store *BoltStore) multi_get(bucket *bolt.Bucket, suffix string, nested bool) ([]interface{}, error) {	
	results := []interface{}{}
	if bucket != nil {
		c := bucket.Cursor()
		for k, v := c.First(); k != nil; k, v = c.Next() {
			if suffix == "" {
				if v != nil {
					vv, err := store.coder.Decode(v)
					if err != nil {
						return nil, err
					}
					results = append(results, vv)
					continue
				}

				if nested {
					bucket_nested := bucket.Bucket(k)
					if bucket_nested != nil {
						nests , err := store.multi_get(bucket_nested, suffix, nested)
						if err != nil {
							return results, err
						}
						results = append(results, nests...)
					}
				}
			} else {
				if strings.HasSuffix(string(k), suffix) {
					vv, err := store.get(bucket, string(k), true)
					if err != nil {
						return results, err
					}
					results = append(results, vv)
					continue
				} else {
					if nested && v == nil {
						bucket_nested := bucket.Bucket(k)
						if bucket_nested != nil {
							nests , err := store.multi_get(bucket_nested, suffix, nested)
							if err != nil {
								return results, err
							}
							results = append(results, nests...)
						}
					}
				}
			}
		}
	}
	return results, nil
}
