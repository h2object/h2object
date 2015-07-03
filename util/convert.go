package util

import (
	"strconv"
	"reflect"
	"github.com/mitchellh/mapstructure"
)

func map2slicehook(f reflect.Kind, t reflect.Kind, data interface{}) (interface{}, error) {
	if f != reflect.Map || t != reflect.Slice {
		return data, nil
	}

	s := []interface{}{}
	dataVal := reflect.Indirect(reflect.ValueOf(data))
	for i := 0; i < dataVal.Len(); i++ {
		vidx := dataVal.MapIndex(reflect.ValueOf(strconv.Itoa(i))) 
		if !vidx.IsValid() {
			break
		} 
		s = append(s, vidx.Interface())
	}
	return s, nil
}

func Convert(src interface{}, dest interface{}) error {
	config := &mapstructure.DecoderConfig{
		DecodeHook: map2slicehook,
		Metadata:	nil,
		Result:		dest,
		TagName:	"object",		 
	}

	decoder, err := mapstructure.NewDecoder(config)
	if err != nil {
		return err
	}

	return decoder.Decode(src)
}
