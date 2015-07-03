package nodes

import (
	"sync"
	"github.com/h2object/h2object/util"
)

type node struct{
	sync.RWMutex
	parent *node
	name string
	path string
	sons map[string]*node
	bind interface{}
}

type Nodes struct{
	root *node
}

func NewNodes() *Nodes {
	return &Nodes{
		root: &node{
			parent: nil,
			name: "/",
			path: "/",
			sons: make(map[string]*node),
			bind: nil,
		},
	}
}

func (ns *Nodes) Node(path string) *node {
	folders := util.PathFolders(path)
	if folders[0] != "/" || len(folders) == 0 {
		return nil
	}
	
	var n *node = ns.root
	if len(folders) > 1 {
		for _, v := range folders[1:] {
			n = n.ChildMust(v)
		}
	}
	return n
}

func (n *node) Path() string {
	return n.path
}
func (n *node) ChildPath(name string) string {
	if n.path == "/" {
		return n.path + name
	}
	return n.path + "/" + name
}

func (n *node) Parent() *node {
	return n.parent
}

func (n *node) ChildMust(name string) *node {
	n.RLock()
	defer n.RUnlock()

	if son, ok := n.sons[name]; ok {
		return son
	}

	son := &node{
		parent: n,
		name: name,
		path: n.ChildPath(name),
		sons: make(map[string]*node),
		bind: nil,
	}
	n.sons[name] = son
	return son
}

func (n *node) SetBind(b interface{}) {
	n.Lock()
	defer n.Unlock()
	n.bind = b
}

func (n *node) GetBind() interface{} {
	n.RLock()
	defer n.RUnlock()
	return n.bind
}
