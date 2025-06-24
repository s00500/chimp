package store

import (
	"context"
	"sync"
)

type Initializable[T any] interface {
	Initialize() T
}

type Lockable[T Initializable[T]] struct {
	mu sync.RWMutex
	v  T
}

func (l *Lockable[T]) Use() (ref T, mutate func(func(state *T)), drop context.CancelFunc) {
	l.mu.RLock()

	drop = context.CancelFunc(func() {
		l.mu.Unlock()
	})

	mutate = func(f func(state *T)) {
		l.mu.RUnlock()
		l.mu.Lock()
		f(&l.v)
		l.mu.Unlock()
		l.mu.RLock()
	}

	return l.v, mutate, drop // Returning as cancelfunc to ensure it gets called and linted
}

/*
func (l *Lockable[T]) Read(read func(state T)) {
	l.mu.RLock()
	read(l.v)
	l.mu.RUnlock()
}
*/

func (l *Lockable[T]) MutateOnly(mutate func(state *T)) {
	l.mu.Lock()
	mutate(&l.v)
	l.mu.Unlock()
}
