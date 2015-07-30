package queue

import (
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
