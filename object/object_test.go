package object

import (
	"log"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

type Person struct{
	Name string `object:"name"`
	Age  int64	`object:"age"`
	Good bool	`object:"good"`
}

func TestOBJECT(t *testing.T) {
	log.Println("======OBJECT test start========")

	obj, err := NewObject("/","name")
	assert.Nil(t, err)

	assert.Nil(t, obj.SetValue("liujianping"))
	
	assert.Equal(t, obj.Value(), "liujianping")

	var name string

	assert.Nil(t, obj.GetValue("/name", &name))
	assert.Equal(t, name, "liujianping")

	obj2, err := NewObject("/", "age")
	assert.Nil(t, err)

	assert.Nil(t, obj2.SetValue(35))
	assert.Equal(t, obj2.Value(), 35)

	var age int
	assert.Nil(t, obj2.GetValue("/age", &age))
	assert.Equal(t, age, 35)

	ns1, err := NewNamespace("/user")
	assert.Nil(t, err)

	ns2, err := NewNamespace("/user/liujianping")
	assert.Nil(t, err)

	assert.Nil(t, ns2.SetValue(map[string]interface{}{
			"name":"liujianping",
			"age":35,
			"good": true,
		}))

	assert.Nil(t, ns1.AddObject(ns2))

	ns3, err := CreateOBJECTByValue("/user/zhangsan", Person{
		Name: "zhangsan",
		Age:23,
		Good: false,
		})
	assert.Nil(t, err)
	assert.Nil(t, ns1.AddObject(ns3))

	log.Println("ns1 value:", ns1.Value())
	log.Println("ns2 value:", ns2.Value())
	log.Println("ns3 value:", ns3.Value())

	var name2 string
	assert.Nil(t, ns2.GetValue("/user/liujianping/name", &name2))
	assert.Equal(t, name2, "liujianping")

	var name3 string
	assert.Nil(t, ns3.GetValue("/user/zhangsan/name", &name3))
	assert.Equal(t, name3, "zhangsan")

	var p1 Person
	assert.Nil(t, ns3.GetValue("/user/zhangsan", &p1))
	log.Println("p1=", p1)

 	ps := []Person{
 		Person{
 			Name: "zhangsan1",
			Age:23,
			Good: false,
 		},
 		Person{
 			Name: "zhangsan2",
			Age:24,
			Good: true,
 		},
 		Person{
 			Name: "zhangsan3",
			Age:25,
			Good: false,
 		},
 	}
 	ns4, err := CreateOBJECTByValue("/friends", ps)
	assert.Nil(t, err)
	log.Println("ns4 value:", ns4.Value())

	var rs []Person
	assert.Nil(t, ns4.GetValue("/friends", &rs))
	log.Println("ns4 get value:", rs)

	log.Println("======OBJECT test end========")
}


func TestByteCoder(t *testing.T) {
	log.Println("======BYTECODER test start========")
	var bcode BoltCoder

	str, err := NewObject("/object", "string")
	assert.Nil(t, err)
	assert.Nil(t, str.SetValue("abcd"))

	byts, err := str.GetBytes(bcode)
	assert.Nil(t, err)

	str2, err := NewObject("/object", "string")
	assert.Nil(t, err)
	assert.Nil(t, str2.SetBytes(bcode, byts))

	assert.Equal(t, str2.Value(), "abcd")

	i, err := NewObject("/object", "int")
	assert.Nil(t, err)
	assert.Nil(t, i.SetValue(23))

	byts2, err := i.GetBytes(bcode)
	assert.Nil(t, err)

	i2, err := NewObject("/object", "int")
	assert.Nil(t, err)
	assert.Nil(t, i2.SetBytes(bcode, byts2))

	assert.Equal(t, i2.Value(), 23)


	f, err := NewObject("/object", "float")
	assert.Nil(t, err)
	assert.Nil(t, f.SetValue(2.3))

	byts3, err := f.GetBytes(bcode)
	assert.Nil(t, err)

	f2, err := NewObject("/object", "int")
	assert.Nil(t, err)
	assert.Nil(t, f2.SetBytes(bcode, byts3))

	assert.Equal(t, f2.Value(), 2.3)

	now := time.Now()

	tm, err := NewObject("/object", "time")
	assert.Nil(t, err)
	assert.Nil(t, tm.SetValue(now))

	byts4, err := tm.GetBytes(bcode)
	assert.Nil(t, err)

	tm2, err := NewObject("/object", "time")
	assert.Nil(t, err)
	assert.Nil(t, tm2.SetBytes(bcode, byts4))

	assert.Equal(t, tm2.Value().(time.Time).Equal(now), true)


	log.Println("======BYTECODER test end========")
}

func TestStore(t *testing.T) {
	log.Println("======BoltStore test start========")
	var coder BoltCoder

	store := NewBoltStore("./testdata","bolt.db", coder)
	assert.Nil(t, store.Load())
	ns, err := CreateOBJECTByValue("/user/zhangsan", Person{
		Name: "zhangsan",
		Age:23,
		Good: false,
		})
	assert.Nil(t, err)

	assert.Nil(t, store.PutOBJECT(ns))

	obj, err := store.GetOBJECT(ns.URI())
	assert.Nil(t, err)
	log.Println("get object:", obj.Value())

	name, err := store.GetOBJECT("/user/zhangsan/name")
	assert.Nil(t, err)
	log.Println("get name:", name.Value())

	age, err := store.GetOBJECT("/user/zhangsan/age")
	assert.Nil(t, err)
	log.Println("get age:", age.Value())

	good, err := store.GetOBJECT("/user/zhangsan/good")
	assert.Nil(t, err)
	log.Println("get good:", good.Value())

	sz, err := store.Size("/")
	assert.Nil(t, err)
	log.Println("size of (/):", sz)

	sz1, err := store.Size("/user")
	assert.Nil(t, err)
	log.Println("size of (/user):", sz1)
	
	sz2, err := store.Size("/user/zhangsan")
	assert.Nil(t, err)
	log.Println("size of (/user/zhangsan):", sz2)

	sz3, err := store.Size("/user/zhangsan/name")
	assert.Nil(t, err)
	log.Println("size of (/user/zhangsan/name):", sz3)

	log.Println("======BoltStore test end========")
}

func TestIndexes(t *testing.T) {
	log.Println("======BleveIndexes test start========")
	
	log.Println("======BleveIndexes test end========")
}

