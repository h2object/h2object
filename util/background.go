package util

import (
	"sync"
)

type Background struct{
	sync.WaitGroup
}

func (b *Background) Work(fn func()) {
	b.Add(1)
	go func() {
		fn()
		b.Done()
	}()	
}
