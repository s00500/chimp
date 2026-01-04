package store

// SessionBackend defines the interface for session storage backends.
// Implementations can store sessions in memory, Redis, or other backends.
type SessionBackend[T Initializable[T]] interface {
	// Get retrieves a session by ID. Returns the session and true if found,
	// or a new initialized session and false if not found.
	Get(id string) (*Session[T], bool)

	// Set stores a session by ID.
	Set(id string, session *Session[T])

	// Delete removes a session by ID.
	Delete(id string)
}
