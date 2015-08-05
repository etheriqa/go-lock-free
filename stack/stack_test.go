package stack

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStack(t *testing.T) {
	assert := assert.New(t)

	{
		s := NewStack()
		_, ok := s.Pop()
		assert.False(ok)
	}

	{
		s := NewStack()
		s.Push(1)

		{
			n, ok := s.Pop()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			_, ok := s.Pop()
			assert.False(ok)
		}
	}

	{
		s := NewStack()
		s.Push(1)
		s.Push(2)

		{
			n, ok := s.Pop()
			assert.Equal(2, n)
			assert.True(ok)
		}

		{
			n, ok := s.Pop()
			assert.Equal(1, n)
			assert.True(ok)
		}

		{
			_, ok := s.Pop()
			assert.False(ok)
		}
	}
}

func BenchmarkLockFreeStack(b *testing.B) {
	s := NewStack()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			s.Push(nil)
			s.Pop()
		}
	})

	if _, ok := s.Pop(); ok {
		b.Fail()
	}
}
