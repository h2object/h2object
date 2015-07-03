package object

import (
	"os"
	"fmt"
	"sync"
	"path"
	"github.com/blevesearch/bleve"
	"github.com/h2object/h2object/nodes"
	"github.com/h2object/h2object/util"
)

var IndexParent bool = true
var IndexOffset int = 0
var IndexSize int = 20

type BleveIndexes struct{
	sync.RWMutex
	root 	string
	name 	string
	nodes *nodes.Nodes
}

func NewBleveIndexes(indexRoot string, indexName string) *BleveIndexes {
	return &BleveIndexes{
		root: indexRoot,
		name: indexName,
		nodes: nodes.NewNodes(),
	}
}

func (indexes *BleveIndexes) index(namespace string, createIfNotExist bool) (bleve.Index, error) {
	dir := path.Join(indexes.root, namespace) 
	fn := path.Join(indexes.root, namespace, indexes.name)
	
	if _, err := os.Stat(fn); err == nil {
		if idx, err := bleve.Open(fn); err == nil {
			return idx, nil
		}
	}

	if !createIfNotExist {
		return nil, nil
	}

	os.RemoveAll(fn)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	mapping := bleve.NewIndexMapping()
	idx, err := bleve.New(fn, mapping)
	return idx, err
}

func (indexes *BleveIndexes) CreateIndex(namespace string) error {
	ns, err := util.Namespace(namespace)
	if err != nil {
		return err
	}

	indexes.Lock()
	defer indexes.Unlock()

	node := indexes.nodes.Node(ns)
	if node == nil {
		return fmt.Errorf("none node for namespace: %s", ns)
	}

	if node.GetBind() == nil {
		index, err := indexes.index(ns, true)
		if err != nil {
			return err
		}
		node.SetBind(index)
	}
	return nil
}

func (indexes *BleveIndexes) RemoveIndex(namespace string) error {
	ns, err := util.Namespace(namespace)
	if err != nil {
		return err
	}

	indexes.Lock()
	defer indexes.Unlock()

	node := indexes.nodes.Node(ns)
	if node == nil {
		return fmt.Errorf("none node for namespace: %s", ns)
	}

	if node.GetBind() != nil {
		node.SetBind(nil)
	}

	// remove the index from fs
	os.RemoveAll(path.Join(indexes.root, namespace, indexes.name))
	return nil	
}

func (indexes *BleveIndexes) IndexIfNotExist(uri string, value interface{}) error {
	indexes.Lock()
	defer indexes.Unlock()

	node := indexes.nodes.Node(path.Dir(uri))
	if node == nil {
		return fmt.Errorf("none node for namespace: %s", path.Dir(uri))
	}

	for {
		if bind := node.GetBind(); bind != nil {
			index := bind.(bleve.Index)
			if err := index.Index(uri, value); err != nil {
				return err
			}
		} else {
			index, err := indexes.index(node.Path(), true)
			if err != nil {
				return err
			}
			if index != nil {
				node.SetBind(index)
				index.Index(uri, value)		
			}			
		}

		if !IndexParent {
			break
		}
		node = node.Parent()
		if node == nil {
			return nil
		}
	}

	return nil
}

func (indexes *BleveIndexes) Index(uri string, value interface{}) error {
	indexes.Lock()
	defer indexes.Unlock()

	node := indexes.nodes.Node(path.Dir(uri))
	if node == nil {
		return fmt.Errorf("none node for namespace: %s", path.Dir(uri))
	}

	for {
		if bind := node.GetBind(); bind != nil {
			index := bind.(bleve.Index)
			if err := index.Index(uri, value); err != nil {
				return err
			}
		} else {
			index, err := indexes.index(node.Path(), false)
			if err != nil {
				return err
			}
			if index != nil {
				node.SetBind(index)
				index.Index(uri, value)		
			}			
		}

		if !IndexParent {
			break
		}
		node = node.Parent()
		if node == nil {
			return nil
		}
	}

	return nil
}

func (indexes *BleveIndexes) Query(namespace string, query interface{}, offset int64, size int64) (int64, []string, error) {
	var total int64 = 0
	var uris []string
	ns, err := util.Namespace(namespace)
	if err != nil {
		return total, uris, err
	}

	indexes.RLock()
	defer indexes.RUnlock()
	
	node := indexes.nodes.Node(ns)
	if node == nil {
		return total, uris, fmt.Errorf("none node for namespace: %s", ns)
	}

	var idx bleve.Index
	if node.GetBind() == nil {
		index, err := indexes.index(ns, false)
		if err != nil {
			return total, uris, err
		}
		if index == nil {
			return total, uris, nil
		}

		node.SetBind(index)
		idx = index
	} else {
		idx = node.GetBind().(bleve.Index)
	}

	q, ok := query.(bleve.Query)
	if !ok {
		return total, uris, fmt.Errorf("query type convert failed")
	} 
	request := bleve.NewSearchRequestOptions(q, IndexSize, IndexOffset, false)		
	response, err := idx.Search(request);
	if err != nil {
		return total, uris, err
	}
	total = int64(response.Total)
	for _, doc := range response.Hits {
		uris = append(uris, doc.ID)
	}
	return total, uris, nil
}