package util

import (
	"sync/atomic"
)

type AtomicInt64 struct{
	val *int64
}

func NewAtomicInt64(init int64) *AtomicInt64 {
	var i int64 = init
	return &AtomicInt64{
		val: &i,
	}
}

func (ai *AtomicInt64) Caculate(delta int64) int64 {
	if delta == 0 {
		return atomic.LoadInt64(ai.val)
	}
	return atomic.AddInt64(ai.val, delta)
}

func (ai *AtomicInt64) Value() int64 {
	return atomic.LoadInt64(ai.val)
}