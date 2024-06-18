package store

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"

	"go.uber.org/atomic"
)

// Implement an in-memory, cookiebased session store. We can specify the type of the state in the beginning
// Security wise we fully depend on gorillas module.
// We only use this store to ease dealing with this session in a typed manner

const sessionExpireryTime = time.Hour

type Session[T Initializable[T]] struct {
	l               sync.RWMutex
	lastInteraction atomic.Time

	state Lockable[T]

	//values map[string]interface{} // other values
	notfresh bool
}

func CreateSessionStore[T Initializable[T]](gorillaStore *sessions.CookieStore) (getSession func(id string) *Session[T], middleWare func(next http.Handler) http.Handler) {
	var globalStore map[string]*Session[T] = map[string]*Session[T]{}
	var globalStoreMu sync.RWMutex

	GetSession := func(id string) *Session[T] {
		globalStoreMu.RLock()
		if s, ok := globalStore[id]; ok {
			globalStoreMu.RUnlock()
			return s
		}
		globalStoreMu.RUnlock()

		// Create New
		s := &Session[T]{state: Lockable[T]{}}
		s.state.v = s.state.v.Initialize()
		s.lastInteraction.Store(time.Now())
		globalStoreMu.Lock()
		globalStore[id] = s
		globalStoreMu.Unlock()

		return s
	}

	middleWare = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := gorillaStore.Get(r, "session-name")

			id, ok := session.Values["id"].(string)
			if !ok {
				session.Values["id"] = uuid.NewString()
				id = session.Values["id"].(string)
			}

			s := GetSession(id)

			session.Save(r, w)
			ctx := context.WithValue(r.Context(), "session", s)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	go func() {
		t := time.NewTicker(time.Second * 15)
		for {
			select {
			case <-t.C:
				globalStoreMu.RLock()
				defer globalStoreMu.RUnlock()

				for key, s := range globalStore {
					// Check expirery
					if time.Since(s.lastInteraction.Load()) > sessionExpireryTime {
						// reap it
						// grap the lock on the key first
						s.l.Lock()
						globalStoreMu.RUnlock()
						globalStoreMu.Lock()

						delete(globalStore, key)
						globalStoreMu.Unlock()
						globalStoreMu.RLock()
						s.l.Unlock()
					}
				}
			}
		}
	}()
	return GetSession, middleWare
}

func (s *Session[T]) ReadRef() (ref T, drop context.CancelFunc) {
	return s.state.ReadRef()
}
func (s *Session[T]) Read(read func(state T)) {
	s.state.Read()
	defer s.state.Drop()
	read(s.state.v)
}

func (s *Session[T]) Mutate(mutate func(state *T)) {
	s.state.Mutate(mutate)
	s.lastInteraction.Store(time.Now())
}

func (s *Session[T]) IsFresh() bool {
	s.l.RLock()
	if !s.notfresh {
		s.l.RUnlock()
		s.l.Lock()
		s.notfresh = true
		s.l.Unlock()
		return true
	}

	s.l.RUnlock()
	return false
}
