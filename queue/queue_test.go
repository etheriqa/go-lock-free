package queue

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestQueue(t *testing.T) {
	assert := assert.New(t)

	{
		q := NewQueue()
		_, ok := q.Dequeue()
		assert.False(ok)
	}

	{
		q := NewQueue()
		q.Enqueue(1)

		{
			n, ok := q.Dequeue()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			_, ok := q.Dequeue()
			assert.False(ok)
		}
	}

	{
		q := NewQueue()
		q.Enqueue(1)
		q.Enqueue(2)

		{
			n, ok := q.Dequeue()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			n, ok := q.Dequeue()
			assert.Equal(2, n)
			assert.True(ok)
		}

		{
			_, ok := q.Dequeue()
			assert.False(ok)
		}
	}
}

func BenchmarkNaiveQueue(b *testing.B) {
	type node struct {
		value interface{}
		next  *node
	}
	sentinel := &node{}
	head := sentinel
	tail := sentinel
	mtxHead := sync.Mutex{}
	mtxTail := sync.Mutex{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			func() {
				mtxTail.Lock()
				defer mtxTail.Unlock()

				tail.next = &node{}
				tail = tail.next
			}()

			func() {
				mtxHead.Lock()
				defer mtxHead.Unlock()

				head = head.next
			}()
		}
	})
}

func BenchmarkChannel(b *testing.B) {
	q := make(chan interface{}, 1000)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q <- nil
			<-q
		}
	})
}

func BenchmarkLockFreeQueue(b *testing.B) {
	q := NewQueue()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.Enqueue(nil)
			q.Dequeue()
		}
	})

	if _, ok := q.Dequeue(); ok {
		b.Fail()
	}
}
