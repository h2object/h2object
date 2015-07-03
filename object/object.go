package object

import (
	"fmt"
	"path"
	"errors"
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"github.com/h2object/cast"
	"github.com/h2object/h2object/util"
)

type OBJECT_TYPE int 

const (
	OBJECT_T = iota
	NAMESPACE_T = iota
)

var CaseSensitive bool = false

type OBJECT interface{
	Namespace() string
	Key() string
	URI() string
	Type() OBJECT_TYPE
	Value() interface{}
	GetBytes(bc ByteCoder) ([]byte, error)
	SetBytes(bc ByteCoder, bytes []byte) error
	SetValue(value interface{}) error
	GetValue(uri string, ret interface{}) error
}


func CreateOBJECTByType(uri string, atype OBJECT_TYPE) (OBJECT, error) {
	switch atype {
	case OBJECT_T:
		return NewObjectByURI(uri)
	case NAMESPACE_T:
		return NewNamespace(uri)
	}
	return nil, errors.New("unknown type")
}

func CreateOBJECTByValue(uri string, value interface{}) (OBJECT, error) {
	if is_value_basic(value) {
		obj, err := NewObjectByURI(uri)
		if err != nil {
			return nil, err
		}
		if err := obj.SetValue(value); err != nil {
			return nil, err
		}
		return obj, nil
	}

	ns, err := NewNamespace(uri)
	if err != nil {
		return nil, err
	}		

	if err := ns.SetValue(value); err != nil {
		return nil, err
	}
	return ns, nil
}

type Object struct{
	namespace 	string
	key 	  	string
	value 	  	interface{}
}

func NewObject(namespace, key string) (*Object, error) {
	if !strings.HasPrefix(namespace, "/") {
		return nil, errors.New("namespace need absolute path:" + namespace)
	}
	if key == "" {
		return nil, errors.New("object key absent")
	}
	var ns string = namespace
	var k string = key
	if !CaseSensitive {
		ns = strings.ToLower(namespace)
		k  = strings.ToLower(key)	
	}	
	if namespace != "/" {
		ns = strings.TrimSuffix(ns, "/")
	}
	return &Object{
		namespace: ns,
		key: k,
	}, nil
}

func NewObjectByURI(uri string) (*Object, error) {
	return NewObject(path.Split(uri))
}

func (obj *Object) Namespace() string {
	return obj.namespace
}

func (obj *Object) Key() string {
	return obj.key
}

func (obj *Object) URI() string {
	return path.Join(obj.namespace, obj.key)
}

func (obj *Object) Type() OBJECT_TYPE {
	return OBJECT_T
}

func (obj *Object) GetBytes(bc ByteCoder) ([]byte, error) {
	if bc == nil {
		return nil, errors.New("byte coder absent")
	}
	return bc.Encode(obj.value)
}

func (obj *Object) SetBytes(bc ByteCoder, byts []byte) error {
	if bc == nil {
		return errors.New("byte coder absent")
	}
	val, err := bc.Decode(byts)
	if err != nil {
		return err
	}
	return obj.SetValue(val)
}

func (obj *Object) Value() interface{} {
	return obj.value
}

// implement json marshal only object value
func (obj *Object) MarshalJSON() ([]byte, error) {
	return json.Marshal(obj.value)
}

func (obj *Object) SetValue(value interface{}) error {
	if !is_value_basic(value) {
		return errors.New("value is not basic type")	
	}
	obj.value = value
	return nil
}

func (obj *Object) GetValue(uri string, ret interface{}) error {
	if uri == path.Join(obj.namespace, obj.key) {				
		return util.Convert(obj.value, ret)
	}
	return errors.New(uri + " wrong")
}

type Namespace struct{
	namespace string
	objects   map[string]OBJECT
}

func NewNamespace(ns string) (*Namespace, error) {
	if !strings.HasPrefix(ns, "/") {
		return nil, errors.New("namespace need absolute path:" + ns)
	}
	var n string = ns
	if !CaseSensitive {
		n = strings.ToLower(ns)
	}
	if n != "/" {
		n = strings.TrimSuffix(n, "/")
	}
	return &Namespace{
		namespace: n,
		objects: make(map[string]OBJECT),
	},nil
}

func (ns *Namespace) Namespace() string {
	return ns.namespace
}

func (ns *Namespace) Key() string {
	return ""
}

func (ns *Namespace) URI() string {
	return ns.namespace
}

func (ns *Namespace) GetBytes(bc ByteCoder) ([]byte, error) {	
	if bc == nil {
		return nil, errors.New("byte coder absent")
	}
	return bc.Encode(ns.Value())
}

func (ns *Namespace) SetBytes(bc ByteCoder, byts []byte) error {	
	if bc == nil {
		return errors.New("byte coder absent")
	}
	val, err := bc.Decode(byts)
	if err != nil {
		return err
	}
	return ns.SetValue(val)
}

// implement json marshal only object value
func (ns *Namespace) MarshalJSON() ([]byte, error) {
	return json.Marshal(ns.Value())
}


func (ns *Namespace) Type() OBJECT_TYPE {
	return NAMESPACE_T
}

func (ns *Namespace) Value() interface{} {
	result := make(map[string]interface{})
	for k, v := range ns.objects {
		result[k] = v.Value()
	}
	return result
}

func (ns *Namespace) SetValue(value interface{}) error {
	if is_value_basic(value) {
		return errors.New("value is basic type")
	}
	
	switch value_type_kind(value) {
	case reflect.Slice:
		vs := reflect.ValueOf(value)
		for i := 0; i <  vs.Len(); i++ {
			if err := ns.Set(strconv.Itoa(i), vs.Index(i).Interface()); err != nil {
				return err
			}
		}
		return nil
	case reflect.Map:
		objs, err := cast.ToStringMapE(value)
		if err != nil {
			return err
		}

		for k, v := range objs {
			if err := ns.Set(k, v); err != nil {
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
				if err := ns.Set(tag, vv.Field(i).Interface()); err != nil {
					return err
				}		
			}
		}	
		return nil
	}
	return errors.New("unknow error")
}

func (ns *Namespace) GetValue(uri string, ret interface{}) error {
	if !strings.HasPrefix(uri, ns.namespace) {
		return errors.New(uri + " wrong")
	}
	folders := strings.Split(strings.Trim(strings.TrimPrefix(uri, ns.namespace), "/"), "/")
	var nsloop OBJECT = ns	
	for _, folder := range folders {
		if len(folder) == 0 {
			continue
		}
		obj, ok := nsloop.(*Namespace).objects[folder]
		if !ok {
			return errors.New(uri + " wrong")
		}

		if obj.Type() == OBJECT_T {
			return obj.GetValue(uri, ret)
		}

		nsloop = obj
	}

	return util.Convert(nsloop.Value(), ret)
}


func (ns *Namespace) AddObject(obj OBJECT) error {
	switch obj.Type() {
	case OBJECT_T:
		if obj.Namespace() != ns.namespace {
			return fmt.Errorf("namespace not equal (%s) (%s)", obj.Namespace(), ns.namespace)
		}
		ns.objects[obj.Key()] = obj
	case NAMESPACE_T:
		parent, key := path.Split(obj.Namespace())
		if strings.TrimSuffix(parent, "/") != ns.namespace {
			return fmt.Errorf("namespace not equal (%s) (%s)", parent, ns.namespace)
		}
		ns.objects[key] = obj	
	}
	
	return nil
}

func (ns *Namespace) DelObject(obj OBJECT) {
	switch obj.Type() {
	case OBJECT_T:
		if obj.Namespace() == ns.namespace {
			delete(ns.objects, obj.Key())
		}
	case NAMESPACE_T:
		parent, key := path.Split(obj.Namespace())
		if strings.TrimSuffix(parent, "/") == ns.namespace {
			delete(ns.objects, key)
		}
	}
}

func (ns *Namespace) Set(key string, val interface{}) error {
	uri := path.Join(ns.namespace, key)
	obj, err := CreateOBJECTByValue(uri, val)
	if err != nil {
		return err
	}
	return ns.AddObject(obj)
}

func (ns *Namespace) Del(key string) {
	delete(ns.objects, key)
}


