package store

import (
	"context"
	"crypto/rand"
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
	dirty    atomic.Bool // Track if session needs to be saved
}

func generateRandomKey() []byte {
	key := make([]byte, 32) // 32 bytes = 256-bit key
	_, err := rand.Read(key)
	if err != nil {
		panic(err)
	}
	return key
}

func CookieStore(sessionKey []byte) *sessions.CookieStore {
	if len(sessionKey) == 0 {
		sessionKey = generateRandomKey()
	}
	sessionstore := sessions.NewCookieStore(sessionKey)
	sessionstore.Options.SameSite = http.SameSiteDefaultMode
	sessionstore.Options.Secure = false // needed as we dont use https
	return sessionstore
}

// CreateStaticStore: Rarely we need to use a fixed static store, EG: Wails fails to support cookies...
func CreateStaticStore[T Initializable[T]]() (middleWare func(next http.Handler) http.Handler) {
	// Wails fails to support cookies... this sucks...
	// We need a global state now... For now we create it inside of this function.
	var globalState Lockable[T]
	globalState.MutateOnly(func(s *T) {
		newValue := (*s).Initialize()
		*s = newValue
	})

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), "session", &globalState)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CreateSessionStore[T Initializable[T]](sessionname string, gorillaStore *sessions.CookieStore, backend SessionBackend[T]) (middleWare func(next http.Handler) http.Handler) {
	middleWare = func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := gorillaStore.Get(r, sessionname)

			id, ok := session.Values["id"].(string)
			if !ok {
				session.Values["id"] = uuid.NewString()
				id = session.Values["id"].(string)
			}

			s, exists := backend.Get(id)

			if !exists {
				// New session, save cookie
				session.Save(r, w)
				s.dirty.Store(true) // Mark as dirty to ensure first save
			}

			// Set up dirty flag callback so mutations mark session for saving
			s.State.markDirty = func() {
				s.dirty.Store(true)
			}

			// Update last interaction time in memory
			lastInteraction := s.lastInteraction.Load()
			now := time.Now()

			// Only update if more than 1 minute has passed (reduce unnecessary writes)
			if now.Sub(lastInteraction) > time.Minute {
				s.lastInteraction.Store(now)
				s.dirty.Store(true)
			}

			ctx := context.WithValue(r.Context(), "session", &s.State)
			next.ServeHTTP(w, r.WithContext(ctx))

			// Only persist if session was modified (dirty flag set by mutations or last interaction update)
			if s.dirty.Load() {
				backend.Set(id, s)
				s.dirty.Store(false)
			}
		})
	}
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
