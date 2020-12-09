// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"github.com/gogf/gf/container/gmap"
	"time"

	"github.com/gogf/gf/os/gcache"
)

// Manager for sessions.
type Manager struct {
	ttl         time.Duration // TTL for sessions.
	storage     Storage       // Storage interface for session storage.
	sessionData *gcache.Cache // Session data cache for session TTL.
}

// New creates and returns a new session manager.
func New(ttl time.Duration, storage ...Storage) *Manager {
	m := &Manager{
		ttl:         ttl,
		sessionData: gcache.New(),
	}
	if len(storage) > 0 && storage[0] != nil {
		m.storage = storage[0]
	} else {
		m.storage = NewStorageFile()
	}
	return m
}

// New creates or fetches the session for given session id.
// The parameter <sessionId> is optional, it creates a new one if not it's passed
// depending on Storage.New.
func (m *Manager) New(sessionId ...string) *Session {
	var id string
	if len(sessionId) > 0 && sessionId[0] != "" {
		id = sessionId[0]
	}
	return &Session{
		id:      id,
		manager: m,
	}
}

// SetStorage sets the session storage for manager.
func (m *Manager) SetStorage(storage Storage) {
	m.storage = storage
}

// SetTTL the TTL for the session manager.
func (m *Manager) SetTTL(ttl time.Duration) {
	m.ttl = ttl
}

// TTL returns the TTL of the session manager.
func (m *Manager) TTL() time.Duration {
	return m.ttl
}

// UpdateSessionTTL updates the ttl for given session.
func (m *Manager) UpdateSessionTTL(sessionId string, data *gmap.StrAnyMap) {
	m.sessionData.Set(sessionId, data, m.ttl)
}
