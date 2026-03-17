package aring

import "sync"

type Ring[T any] struct {
	mu    sync.RWMutex
	buf   []T
	size  int
	head  int
	count int
}

func New[T any](size int) *Ring[T] {
	return &Ring[T]{buf: make([]T, size), size: size}
}

func (r *Ring[T]) Push(v T) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.buf[r.head] = v
	r.head = (r.head + 1) % r.size
	if r.count < r.size {
		r.count++
	}
}

func (r *Ring[T]) Len() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return r.count
}

// Snapshot returns oldest->newest
func (r *Ring[T]) Snapshot() []T {
	r.mu.RLock()
	defer r.mu.RUnlock()

	out := make([]T, r.count)
	if r.count == 0 {
		return out
	}
	start := (r.head - r.count + r.size) % r.size
	for i := 0; i < r.count; i++ {
		out[i] = r.buf[(start+i)%r.size]
	}
	return out
}
