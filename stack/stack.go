package stack

import (
	"math/rand"
	"sync/atomic"
	"unsafe"
)

const nExchangers = 4 // TODO

const (
	stateFree = iota
	statePushing
	statePoping
	stateExchanging
)

type node struct {
	value interface{}
	next  unsafe.Pointer
}

type exchanger struct {
	value interface{}
	state uint
}

type Stack struct {
	top        unsafe.Pointer
	exchangers []unsafe.Pointer
}

func newNode(v interface{}) *node {
	return &node{
		value: v,
		next:  nil,
	}
}

func newExchanger() *exchanger {
	return &exchanger{
		value: nil,
		state: stateFree,
	}
}

func NewStack() *Stack {
	exchangers := make([]unsafe.Pointer, nExchangers)
	for i := 0; i < len(exchangers); i++ {
		exchangers[i] = unsafe.Pointer(newExchanger())
	}
	return &Stack{
		top:        nil,
		exchangers: exchangers,
	}
}

func (s *Stack) Push(v interface{}) {
	for {
		if s.tryPushTop(v) {
			return
		}
		if s.tryPushExchanger(v) {
			return
		}
	}
}

func (s *Stack) Pop() (v interface{}, ok bool) {
	for {
		if top, ok := s.tryPopTop(); ok {
			if top == nil {
				return nil, false
			} else {
				return top.value, true
			}
		}
		if v, ok := s.tryPopExchanger(); ok {
			return v, true
		}
	}
}

func (s *Stack) tryPushTop(v interface{}) bool {
	oldTop := atomic.LoadPointer(&s.top)
	newTop := unsafe.Pointer(newNode(v))
	(*node)(newTop).next = oldTop
	return atomic.CompareAndSwapPointer(&s.top, oldTop, newTop)
}

func (s *Stack) tryPopTop() (*node, bool) {
	oldTop := atomic.LoadPointer(&s.top)
	if oldTop == nil {
		return nil, true
	}
	newTop := (*node)(oldTop).next
	if ok := atomic.CompareAndSwapPointer(&s.top, oldTop, newTop); ok {
		return (*node)(oldTop), true
	}
	return nil, false
}

func (s *Stack) tryPushExchanger(v interface{}) bool {
	// TODO time out
	slotptr := &s.exchangers[rand.Intn(nExchangers)]
	for i := 0; i < 100; i++ {
		oldSlot := atomic.LoadPointer(slotptr)
		switch (*exchanger)(oldSlot).state {
		case stateFree:
			newSlot := unsafe.Pointer(&exchanger{
				value: v,
				state: statePushing,
			})
			if !atomic.CompareAndSwapPointer(slotptr, oldSlot, newSlot) {
				break
			}
			for j := 0; j < 100; j++ {
				currentSlot := atomic.LoadPointer(slotptr)
				if (*exchanger)(currentSlot).state != stateExchanging {
					continue
				}
				atomic.StorePointer(slotptr, unsafe.Pointer(newExchanger()))
				return true
			}
			if !atomic.CompareAndSwapPointer(slotptr, newSlot, oldSlot) {
				atomic.StorePointer(slotptr, unsafe.Pointer(newExchanger()))
				return true
			}
			return false
		case statePushing:
			return false
		case statePoping:
			newSlot := unsafe.Pointer(&exchanger{
				value: v,
				state: stateExchanging,
			})
			if !atomic.CompareAndSwapPointer(slotptr, oldSlot, newSlot) {
				break
			}
			return true
		case stateExchanging:
			return false
		}
	}
	return false
}

func (s *Stack) tryPopExchanger() (interface{}, bool) {
	// TODO time out
	slotptr := &s.exchangers[rand.Intn(nExchangers)]
	for i := 0; i < 100; i++ {
		oldSlot := atomic.LoadPointer(slotptr)
		switch (*exchanger)(oldSlot).state {
		case stateFree:
			newSlot := unsafe.Pointer(&exchanger{
				state: statePoping,
			})
			if !atomic.CompareAndSwapPointer(slotptr, oldSlot, newSlot) {
				break
			}
			for j := 0; j < 100; j++ {
				currentSlot := atomic.LoadPointer(slotptr)
				if (*exchanger)(currentSlot).state != stateExchanging {
					continue
				}
				atomic.StorePointer(slotptr, unsafe.Pointer(newExchanger()))
				return (*exchanger)(currentSlot).value, true
			}
			if !atomic.CompareAndSwapPointer(slotptr, newSlot, oldSlot) {
				currentSlot := atomic.LoadPointer(slotptr)
				atomic.StorePointer(slotptr, unsafe.Pointer(newExchanger()))
				return (*exchanger)(currentSlot).value, true
			}
			return nil, false
		case statePushing:
			newSlot := unsafe.Pointer(&exchanger{
				state: stateExchanging,
			})
			if !atomic.CompareAndSwapPointer(slotptr, oldSlot, newSlot) {
				break
			}
			return (*exchanger)(oldSlot).value, true
		case statePoping:
			return nil, false
		case stateExchanging:
			return nil, false
		}
	}
	return nil, false
}
