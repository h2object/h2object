package object

import (
	"fmt"
	"time"
	"bytes"
	"reflect"
	"github.com/h2object/cast"
)

const (
	value_nil = iota
	value_bool = iota
	value_number = iota
	value_duration = iota
	value_datetime = iota
	value_string = iota	
	value_bytes = iota
)

// JSON basic type: null, bool, number, string
// JSON composite type: slice, map 

type BoltCoder struct{}

func (bc BoltCoder) Decode(value []byte) (interface{}, error) {
	if len(value) == 0 {
		return nil, fmt.Errorf("bytes len zero")
	}
	switch value[0] {
	case value_nil:
		return nil, nil
	case value_bool:
		if value[1] == 1 {
			return true, nil
		}
		return false, nil
	case value_number:
		f, err := cast.ToFloat64E(string(value[1:]))
		if err != nil {
			return nil, err
		}
		return f, nil		
	case value_duration:
		f, err := cast.ToFloat64E(string(value[1:]))
		if err != nil {
			return nil, err
		}
		return time.Duration(f), nil
	case value_datetime:
		tm := &time.Time{}
		if err := tm.UnmarshalBinary(value[1:]); err != nil {
			return nil, err
		}
		return *tm, nil
	case value_string:
		return string(value[1:]), nil
	case value_bytes:
		return value[1:], nil
	}
	return nil, nil
}

func number(val interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	buffer.WriteByte(value_number)
	s, err := cast.ToStringE(cast.ToFloat64(val))
	if err != nil {
		return nil, err
	}
	buffer.WriteString(s)
	return buffer.Bytes(), nil
}


func (bc BoltCoder) Encode(value interface{}) ([]byte, error) {
	var buffer bytes.Buffer
	if value == nil {
		buffer.WriteByte(value_nil)
		return buffer.Bytes(), nil
	}
	
	switch t := value.(type) {
	case bool:
		buffer.WriteByte(value_bool)
		if value.(bool) {
			buffer.WriteByte(1)
		} else {
			buffer.WriteByte(0)
		}
		return buffer.Bytes(), nil
	case int:		
		return number(value)
	case int8:
		return number(value)
	case int16:
		return number(value)
	case int32:
		return number(value)
	case int64:
		return number(value)
	case float32:
		return number(value)
	case float64:
		return number(value)	
	case string:
		buffer.WriteByte(value_string)
		buffer.WriteString(value.(string))
		return buffer.Bytes(), nil	
	case []byte:
		buffer.WriteByte(value_bytes)
		buffer.Write(value.([]byte))
		return buffer.Bytes(), nil	
	case time.Duration:
		buffer.WriteByte(value_duration)
		s, err := cast.ToStringE(int(value.(time.Duration)))
		if err != nil {
			return nil, err
		}
		buffer.WriteString(s)
		return buffer.Bytes(), nil	
	case time.Time:
		buffer.WriteByte(value_datetime)
		binary, err := value.(time.Time).MarshalBinary()
		if err != nil {
			return nil, err
		}
		buffer.Write(binary)
		return buffer.Bytes(), nil	
	default:
		_ = t
	}

	return nil, fmt.Errorf("unsupport type (%s) encode", reflect.TypeOf(value).Kind().String())
}

