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
		_, ok := q.dequeue()
		assert.False(ok)
	}

	{
		q := NewQueue()
		q.enqueue(1)

		{
			n, ok := q.dequeue()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			_, ok := q.dequeue()
			assert.False(ok)
		}
	}

	{
		q := NewQueue()
		q.enqueue(1)
		q.enqueue(2)

		{
			n, ok := q.dequeue()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			n, ok := q.dequeue()
			assert.Equal(2, n)
			assert.True(ok)
		}

		{
			_, ok := q.dequeue()
			assert.False(ok)
		}
	}
}

func BenchmarkNaiveQueue(b *testing.B) {
	q := make([]interface{}, 1000)
	head := 0
	tail := 0
	size := 0
	mtx := sync.Mutex{}

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			for func() bool {
				mtx.Lock()
				defer mtx.Unlock()

				if size == 1000 {
					return true
				}
				q[tail] = nil
				tail = (tail + 1) % 1000
				size++
				return false
			}() {
			}

			for func() bool {
				mtx.Lock()
				defer mtx.Unlock()

				if size == 0 {
					return true
				}
				_ = q[head]
				head = (head + 1) % 1000
				size--
				return false
			}() {
			}
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

func BenchmarkQueue(b *testing.B) {
	q := NewQueue()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			q.enqueue(nil)
			q.dequeue()
		}
	})

	if _, ok := q.dequeue(); ok {
		b.Fail()
	}
}
