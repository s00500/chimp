package store

import (
	"sync"
	"time"
)

type memoryBackend[T Initializable[T]] struct {
	store map[string]*Session[T]
	mu    sync.RWMutex
}

// MemoryStore creates an in-memory session backend.
// Sessions are stored in a map and automatically cleaned up when expired.
func MemoryStore[T Initializable[T]]() SessionBackend[T] {
	backend := &memoryBackend[T]{
		store: make(map[string]*Session[T]),
	}

	// Start background cleanup goroutine
	go func() {
		ticker := time.NewTicker(15 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			backend.cleanup()
		}
	}()

	return backend
}

func (m *memoryBackend[T]) Get(id string) (*Session[T], bool) {
	m.mu.RLock()
	session, exists := m.store[id]
	m.mu.RUnlock()

	if exists {
		return session, true
	}

	// Create new session
	session = &Session[T]{State: Lockable[T]{}}
	session.State.MutateOnly(func(v *T) {
		*v = (*v).Initialize()
	})
	session.lastInteraction.Store(time.Now())

	return session, false
}

func (m *memoryBackend[T]) Set(id string, session *Session[T]) {
	m.mu.Lock()
	m.store[id] = session
	m.mu.Unlock()
}

func (m *memoryBackend[T]) Delete(id string) {
	m.mu.Lock()
	delete(m.store, id)
	m.mu.Unlock()
}

func (m *memoryBackend[T]) cleanup() {
	m.mu.RLock()
	expiredKeys := make([]string, 0)

	for key, session := range m.store {
		if time.Since(session.lastInteraction.Load()) > sessionExpireryTime {
			expiredKeys = append(expiredKeys, key)
		}
	}
	m.mu.RUnlock()

	if len(expiredKeys) > 0 {
		m.mu.Lock()
		for _, key := range expiredKeys {
			// Double-check expiration while holding write lock
			if session, exists := m.store[key]; exists {
				if time.Since(session.lastInteraction.Load()) > sessionExpireryTime {
					delete(m.store, key)
				}
			}
		}
		m.mu.Unlock()
	}
}
