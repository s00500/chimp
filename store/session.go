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

	State Lockable[T]

	notfresh bool
}

func CookieStore(sessionKey []byte) *sessions.CookieStore {
	sessionstore := sessions.NewCookieStore(sessionKey)
	sessionstore.Options.SameSite = http.SameSiteDefaultMode
	sessionstore.Options.Secure = false // needed as we dont use https
	return sessionstore
}

func CreateSessionStore[T Initializable[T]](sessionname string, gorillaStore *sessions.CookieStore) (middleWare func(next http.Handler) http.Handler) {
	var globalStore map[string]*Session[T] = map[string]*Session[T]{}
	var globalStoreMu sync.RWMutex

	GetSessionByID := func(id string) (s *Session[T], isFresh bool) {
		globalStoreMu.RLock()
		if s, ok := globalStore[id]; ok {
			globalStoreMu.RUnlock()
			return s, false
		}
		globalStoreMu.RUnlock()

		// Create New
		s = &Session[T]{State: Lockable[T]{}}
		s.State.MutateOnly(func(v *T) {
			*v = (*v).Initialize()
		})

		s.lastInteraction.Store(time.Now())
		globalStoreMu.Lock()
		globalStore[id] = s
		globalStoreMu.Unlock()

		return s, true
	}

	middleWare = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := gorillaStore.Get(r, sessionname)

			id, ok := session.Values["id"].(string)
			if !ok {
				session.Values["id"] = uuid.NewString()
				id = session.Values["id"].(string)
			}

			s, isFresh := GetSessionByID(id)

			if isFresh {
				session.Save(r, w)
			}

			s.lastInteraction.Store(time.Now()) // Sessions timeout if they are not interactied with for session timeout

			ctx := context.WithValue(r.Context(), "session", &s.State)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

	go func() {
		t := time.NewTicker(time.Second * 15)
		for range t.C {
			globalStoreMu.RLock()
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
			globalStoreMu.RUnlock()
		}
	}()
	return middleWare
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

// TODO:
func (s *Session[T]) Clear() {

	s.l.Lock()
	s.notfresh = false
	s.l.Unlock()

	s.State.MutateOnly(func(v *T) {
		*v = (*v).Initialize()
	})
}

// func GetSession[T Initializable[T]](r *http.Request) *Session[T] {
// 	return r.Context().Value("session").(*Session[T])
// }

/*
// TODO: make this a state geter helper
func State(ctx context.Context) SessionState {
	val := ctx.Value("session")
	if val == nil {
		return SessionState{}
	}
	s := ctx.Value("session").(*store.Session[SessionState])
	if s == nil {
		return SessionState{}
	}

	ref, drop := s.ReadRef()
	defer drop()
	return ref
}

func getSessionData(ctx context.Context) *store.Session[sessionstate.SessionState] {
	return ctx.Value("session").(*store.Session[sessionstate.SessionState])
}

*/
