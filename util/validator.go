package util

import (
	"fmt"
	"reflect"
	"regexp"
	"time"
	"strings"
	"errors"
)

type Validator interface {
	Key() string 
	Value() interface{}
	IsSatisfied(interface{}) bool
	DefaultMessage() string
}

func NumberInt(v interface{}) (int, error) {
	switch v.(type) {
	case int:
		return int(reflect.ValueOf(v).Int()), nil
	case int8:
		return int(reflect.ValueOf(v).Int()), nil
	case int32:
		return int(reflect.ValueOf(v).Int()), nil
	case int64:
		return int(reflect.ValueOf(v).Int()), nil
	case float32:
		return int(reflect.ValueOf(v).Float()), nil
	case float64:
		return int(reflect.ValueOf(v).Float()), nil
	}
	return 0, errors.New("not number type")
}

func NewValidator(key string, value interface{}) (Validator, error) {
	switch strings.ToLower(key) {
	case "required":
		if b, ok := value.(bool); ok {
			if b {
				return ValidRequired(), nil
			} else {
				return nil, errors.New("required need be true")
			}
		}
		return nil, errors.New("required need be bool type")
	case "email":
		if b, ok := value.(bool); ok {
			if b {
				return ValidEmail(), nil
			} else {
				return nil, errors.New("email need be true")
			}
		}
		return nil, errors.New("email need be bool type")
	case "min":
		i, err := NumberInt(value)
		if err != nil {
			return nil, err
		}
		return ValidMin(i), nil
	case "max":
		i, err := NumberInt(value)
		if err != nil {
			return nil, err
		}
		return ValidMax(i), nil
	case "minsize":
		i, err := NumberInt(value)
		if err != nil {
			return nil, err
		}
		return ValidMinSize(i), nil
	case "maxsize":
		i, err := NumberInt(value)
		if err != nil {
			return nil, err
		}
		return ValidMaxSize(i), nil
	case "length":
		i, err := NumberInt(value)
		if err != nil {
			return nil, err
		}
		return ValidLength(i), nil
	case "match":
		if i, ok := value.(string); ok {
			return ValidMatch(i), nil
		}
		return nil, errors.New("min need be float64 type")
	}
	return nil, errors.New("unknown validator format")
}

type Required struct{}

func ValidRequired() *Required {
	return &Required{}
}

func (r *Required) Key() string {
	return "required"
}

func (r *Required) Value() interface{} {
	return true
}

func (r *Required) IsSatisfied(obj interface{}) bool {
	if obj == nil {
		return false
	}

	if str, ok := obj.(string); ok {
		return len(str) > 0
	}
	if b, ok := obj.(bool); ok {
		return b
	}
	if i, ok := obj.(int); ok {
		return i != 0
	}
	if t, ok := obj.(time.Time); ok {
		return !t.IsZero()
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() > 0
	}
	return true
}

func (r *Required) DefaultMessage() string {
	return "Required"
}

type Min struct {
	Min int
}

func ValidMin(min int) *Min {
	return &Min{min}
}

func (m *Min) Key() string {
	return "min"
}

func (m *Min) Value() interface{} {
	return m.Min
}

func (m *Min) IsSatisfied(obj interface{}) bool {
	
	num, err := NumberInt(obj)
	if err == nil {
		return num >= m.Min
	}
	return false
}

func (m *Min) DefaultMessage() string {
	return fmt.Sprintf("Minimum is %d", m.Min)
}

type Max struct {
	Max int
}

func ValidMax(max int) *Max {
	return &Max{max}
}

func (m *Max) Key() string {
	return "max"
}

func (m *Max) Value() interface{} {
	return m.Max
}

func (m *Max) IsSatisfied(obj interface{}) bool {
	num, err := NumberInt(obj)
	if err == nil {
		return num <= m.Max
	}
	return false
}

func (m *Max) DefaultMessage() string {
	return fmt.Sprintf("Maximum is %d", m.Max)
}

// Requires an array or string to be at least a given length.
type MinSize struct {
	Min int
}

func ValidMinSize(min int) *MinSize {
	return &MinSize{min}
}

func (m *MinSize) Key() string {
	return "minsize"
}

func (m *MinSize) Value() interface{} {
	return m.Min
}

func (m *MinSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) >= m.Min
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() >= m.Min
	}
	return false
}

func (m *MinSize) DefaultMessage() string {
	return fmt.Sprintf("Minimum size is %d", m.Min)
}

// Requires an array or string to be at most a given length.
type MaxSize struct {
	Max int
}

func ValidMaxSize(max int) *MaxSize {
	return &MaxSize{max}
}

func (m *MaxSize) Key() string {
	return "maxsize"
}

func (m *MaxSize) Value() interface{} {
	return m.Max
}

func (m *MaxSize) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) <= m.Max
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() <= m.Max
	}
	return false
}

func (m *MaxSize) DefaultMessage() string {
	return fmt.Sprintf("Maximum size is %d", m.Max)
}

// Requires an array or string to be exactly a given length.
type Length struct {
	N int
}

func ValidLength(n int) *Length {
	return &Length{n}
}

func (l *Length) Key() string {
	return "length"
}

func (l *Length) Value() interface{} {
	return l.N
}

func (s *Length) IsSatisfied(obj interface{}) bool {
	if str, ok := obj.(string); ok {
		return len(str) == s.N
	}
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Slice {
		return v.Len() == s.N
	}
	return false
}

func (s *Length) DefaultMessage() string {
	return fmt.Sprintf("Required length is %d", s.N)
}

// Requires a string to match a given regex.
type Match struct {
	str 	string
	Regexp *regexp.Regexp
}

func ValidMatch(str string) *Match {
	regex := regexp.MustCompile(str)
	return &Match{
		str: str,
		Regexp: regex,
	}
}

func (m *Match) Key() string {
	return "match"
}

func (m *Match) Value() interface{} {
	return m.str
}

func (m *Match) IsSatisfied(obj interface{}) bool {
	str := obj.(string)
	return m.Regexp.MatchString(str)
}

func (m *Match) DefaultMessage() string {
	return fmt.Sprintf("Must match %s", m.str)
}

var emailPattern = "^[\\w!#$%&'*+/=?^_`{|}~-]+(?:\\.[\\w!#$%&'*+/=?^_`{|}~-]+)*@(?:[\\w](?:[\\w-]*[\\w])?\\.)+[a-zA-Z0-9](?:[\\w-]*[\\w])?$"

type Email struct {
	*Match
}

func ValidEmail() Email {
	return Email{ValidMatch(emailPattern)}
}

func (e Email) Key() string {
	return "email"
}

func (e Email) Value() interface{} {
	return true
}

func (e Email) DefaultMessage() string {
	return fmt.Sprintf("Must be a valid email address")
}
