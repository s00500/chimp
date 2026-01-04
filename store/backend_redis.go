package store

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisBackend[T Initializable[T]] struct {
	client  redis.UniversalClient
	prefix  string
	ttlFunc func(*Session[T]) time.Duration
}

// sessionData is a serializable representation of Session[T]
type sessionData[T any] struct {
	LastInteraction time.Time
	StateValue      T
	NotFresh        bool
}

// RedisStore creates a Redis-backed session backend.
// Uses redis.UniversalClient for compatibility with Redis, Valkey, DragonflyDB, KeyDB, clusters, etc.
// Sessions are automatically expired using Redis TTL (sliding expiration on each access).
// Optional ttlFunc parameter allows dynamic TTL calculation per session (e.g., based on "keep me logged in").
// If ttlFunc is nil, defaults to sessionExpireryTime (1 hour).
func RedisStore[T Initializable[T]](client redis.UniversalClient, ttlFunc ...func(*Session[T]) time.Duration) SessionBackend[T] {
	backend := &redisBackend[T]{
		client: client,
		prefix: "chimp:session:",
	}
	if len(ttlFunc) > 0 && ttlFunc[0] != nil {
		backend.ttlFunc = ttlFunc[0]
	}
	return backend
}

func (r *redisBackend[T]) Get(id string) (*Session[T], bool) {
	ctx := context.Background()
	key := r.prefix + id

	// Try to get from Redis
	data, err := r.client.Get(ctx, key).Bytes()
	if err == redis.Nil {
		// Session doesn't exist, create new
		//fmt.Printf("[Redis] Session not found, creating new: %s\n", id)
		session := &Session[T]{State: Lockable[T]{}}
		session.State.MutateOnly(func(v *T) {
			*v = (*v).Initialize()
		})
		session.lastInteraction.Store(time.Now())
		return session, false
	}

	if err != nil {
		// Redis error, create new session as fallback
		//fmt.Printf("[Redis] Error getting session %s: %v, creating new\n", id, err)
		session := &Session[T]{State: Lockable[T]{}}
		session.State.MutateOnly(func(v *T) {
			*v = (*v).Initialize()
		})
		session.lastInteraction.Store(time.Now())
		return session, false
	}

	// Deserialize session
	var sd sessionData[T]
	buf := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buf)
	if err := decoder.Decode(&sd); err != nil {
		// Deserialization failed, create new session
		fmt.Printf("[Redis] Error deserializing session %s: %v, creating new\n", id, err)
		session := &Session[T]{State: Lockable[T]{}}
		session.State.MutateOnly(func(v *T) {
			*v = (*v).Initialize()
		})
		session.lastInteraction.Store(time.Now())
		return session, false
	}

	// Reconstruct Session
	session := &Session[T]{State: Lockable[T]{}, notfresh: sd.NotFresh}
	session.State.MutateOnly(func(v *T) {
		*v = sd.StateValue
	})
	session.lastInteraction.Store(sd.LastInteraction)

	// Refresh TTL (sliding expiration)
	ttl := sessionExpireryTime
	if r.ttlFunc != nil {
		ttl = r.ttlFunc(session)
	}
	r.client.Expire(ctx, key, ttl)

	fmt.Printf("[Redis] Session loaded: %s (TTL: %v)\n", id, ttl)
	return session, true
}

func (r *redisBackend[T]) Set(id string, session *Session[T]) {
	ctx := context.Background()
	key := r.prefix + id

	// Serialize session
	ref, _, drop := session.State.Use()
	defer drop()

	sd := sessionData[T]{
		LastInteraction: session.lastInteraction.Load(),
		StateValue:      ref,
		NotFresh:        session.notfresh,
	}

	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(sd); err != nil {
		// Log error but don't panic - session will be recreated on next access
		fmt.Printf("Error encoding session %s: %v\n", id, err)
		return
	}

	// Store in Redis with TTL
	ttl := sessionExpireryTime
	if r.ttlFunc != nil {
		ttl = r.ttlFunc(session)
	}
	if err := r.client.Set(ctx, key, buf.Bytes(), ttl).Err(); err != nil {
		fmt.Printf("[Redis] Error storing session %s: %v\n", id, err)
		return
	}
	fmt.Printf("[Redis] Session stored: %s (TTL: %v, size: %d bytes)\n", id, ttl, len(buf.Bytes()))
}

func (r *redisBackend[T]) Delete(id string) {
	ctx := context.Background()
	key := r.prefix + id
	r.client.Del(ctx, key)
}
