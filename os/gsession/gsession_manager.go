// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"time"

	"github.com/gogf/gf/os/gcache"
)

// Manager for sessions.
type Manager struct {
	ttl      time.Duration // TTL for sessions.
	storage  Storage       // Storage interface for session storage Set/Get.
	sessions *gcache.Cache // Session cache for session TTL.
}

// New creates and returns a new session manager.
func New(ttl time.Duration, storage ...Storage) *Manager {
	m := &Manager{
		ttl:      ttl,
		storage:  NewStorageFile(),
		sessions: gcache.New(),
	}
	if len(storage) > 0 && storage[0] != nil {
		m.storage = storage[0]
	}
	return m
}

// New creates or fetches the session for given session id.
func (m *Manager) New(sessionId ...string) *Session {
	var id string
	if len(sessionId) > 0 && sessionId[0] != "" {
		id = sessionId[0]
	}
	// NOTE:
	// We CANNOT creates and stores it directly to manager
	// as it might be a fake and invalid session id
	// which would consumes your memory as much as possible.
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
// If this session is dirty, it also exports it to storage.
func (m *Manager) UpdateSessionTTL(id string, session *Session) {
	m.sessions.Set(id, session, m.ttl)
}
