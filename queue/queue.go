package queue

import (
	"sync/atomic"
	"unsafe"
)

type node struct {
	value interface{}
	next  unsafe.Pointer
}

type Queue struct {
	head unsafe.Pointer
	tail unsafe.Pointer
}

func newNode(v interface{}) *node {
	return &node{
		value: v,
		next:  nil,
	}
}

func NewQueue() *Queue {
	sentinel := unsafe.Pointer(newNode(nil))
	return &Queue{
		head: sentinel,
		tail: sentinel,
	}
}

func (q *Queue) Enqueue(v interface{}) {
	n := unsafe.Pointer(newNode(v))
	for {
		last := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*node)(last).next)
		if last != atomic.LoadPointer(&q.tail) {
			continue
		}
		if next != nil {
			atomic.CompareAndSwapPointer(&q.tail, last, next)
			continue
		}
		if atomic.CompareAndSwapPointer(&(*node)(last).next, next, n) {
			atomic.CompareAndSwapPointer(&q.tail, last, n)
			return
		}

	}
}

func (q *Queue) Dequeue() (v interface{}, ok bool) {
	for {
		first := atomic.LoadPointer(&q.head)
		last := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*node)(first).next)
		if first != atomic.LoadPointer(&q.head) {
			continue
		}
		if first == last {
			if next == nil {
				return nil, false
			}
			atomic.CompareAndSwapPointer(&q.tail, last, next)
		} else {
			v := (*node)(next).value
			if atomic.CompareAndSwapPointer(&q.head, first, next) {
				return v, true
			}
		}
	}
}
