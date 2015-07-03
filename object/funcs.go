package object

import (
	"time"
	"reflect"
)

func is_value_basic(value interface{}) bool {
	switch value_type_kind(value) {
	case reflect.Chan:
		fallthrough
	case reflect.Func:
		fallthrough
	case reflect.Interface:
		fallthrough
	case reflect.Map:
		fallthrough
	case reflect.Slice:
		switch t := value.(type) {
		case []byte:
			return true
		default:
			_ = t
		}
		return false
	case reflect.Struct:
		switch t := value.(type) {
		case time.Duration:
			return true
		case time.Time:
			return true
		default:
			_ = t
		}
		return false		
	}
	return true
}

func value_type_kind(value interface{}) reflect.Kind {	
	vt := reflect.TypeOf(value)
	kindOfValue := vt.Kind()	
	for {
		if kindOfValue == reflect.Ptr {
			vt = vt.Elem()
			kindOfValue = vt.Kind()
			continue
		} 
		break
	}
	return kindOfValue
}