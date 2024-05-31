package manager

import (
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/alexedwards/scs/v2/memstore"
)

// SessionManager represents a wrapper around the scs.SessionManager (Alex Edwards) type.
// It provides full session support.
type SessionManager struct {
	*scs.SessionManager
}

// CreateSessionManager function creates new SessionManager that can be used for working with sessions.
// It is currently set to run without a database, that is, to store session data in memory.
func CreateSessionManager() *SessionManager {
	sessionManager := scs.New()
	sessionManager.Store = memstore.New()
	sessionManager.Lifetime = time.Hour

	return &SessionManager{sessionManager}
}
