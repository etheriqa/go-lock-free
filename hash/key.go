package hash

import (
	"hash/fnv"
)

type Key interface {
	Hash() uint64
}

type Uint64 uint64
type Uintptr uintptr

func (k Uint64) Hash() uint64 {
	return hash(uint64(k))
}

func (k Uintptr) Hash() uint64 {
	return hash(uint64(k))
}

func hash(k uint64) uint64 {
	f := fnv.New64()
	f.Write([]byte{
		byte(k),
		byte(k >> 8),
		byte(k >> 16),
		byte(k >> 24),
		byte(k >> 32),
		byte(k >> 40),
		byte(k >> 48),
		byte(k >> 56),
	})
	return f.Sum64()
}
