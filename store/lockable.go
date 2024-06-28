package store

import (
	"context"
	"sync"
)

type Initializable[T any] interface {
	Initialize() T
}

type Lockable[T Initializable[T]] struct {
	sync.RWMutex
	v     T
	state uint8 // 0 free, 1 read, 2 write
}

func (l *Lockable[T]) Read() {
	switch l.state {
	case 0:
		l.RLock()
	case 1:
		// do nothing
	case 2:
		l.Unlock()
		l.RLock()
	}
	l.state = 1
}

func (l *Lockable[T]) ReadRef() (ref T, drop context.CancelFunc) {
	switch l.state {
	case 0:
		l.RLock()
	case 1:
		// do nothing
	case 2:
		l.Unlock()
		l.RLock()
	}
	l.state = 1

	drop = context.CancelFunc(func() {
		l.Drop()
	})

	return l.v, drop // Returning as cancelfunc to ensure it gets called and linted
}

func (l *Lockable[T]) Write() {
	switch l.state {
	case 0:
		l.Lock()
	case 1:
		l.RUnlock()
		l.Lock()
	case 2:
		// do nothing
	}
	l.state = 2
}

func (l *Lockable[T]) Drop() {
	switch l.state {
	case 0:
	// do nothing
	case 1:
		l.RUnlock()
	case 2:
		l.Unlock()
	}
	l.state = 0
}

func (l *Lockable[T]) Mutate(mutate func(state *T)) {
	// save state
	oldState := l.state

	l.Write()
	mutate(&l.v)
	// restore old state
	switch oldState {
	case 0: // open
		l.Drop()
	case 1: // Read
		l.Drop()
		l.RLock()
	case 2:
		// was write locked before... not great though
	}
	l.state = oldState
}
