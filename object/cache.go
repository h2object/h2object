package object

import "time"

type Cache interface{
	Set(string, interface{}, time.Duration) 
	Add(string, interface{}, time.Duration) error
	Get(string) (interface{}, bool)
	Delete(string)
	DeleteExpired()
}