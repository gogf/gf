// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"context"
	"time"
)

// Manager for sessions.
type Manager struct {
	ttl     time.Duration // TTL for sessions.
	storage Storage       // Storage interface for session storage.
}

// New creates and returns a new session manager.
func New(ttl time.Duration, storage ...Storage) *Manager {
	m := &Manager{
		ttl: ttl,
	}
	if len(storage) > 0 && storage[0] != nil {
		m.storage = storage[0]
	} else {
		// It uses StorageFile in default.
		m.storage = NewStorageFile(DefaultStorageFilePath, ttl)
	}
	return m
}

// New creates or fetches the session for given session id.
// The parameter `sessionId` is optional, it creates a new one if not it's passed
// depending on Storage.New.
func (m *Manager) New(ctx context.Context, sessionId ...string) *Session {
	var id string
	if len(sessionId) > 0 && sessionId[0] != "" {
		id = sessionId[0]
	}
	return &Session{
		id:      id,
		ctx:     ctx,
		manager: m,
	}
}

// SetStorage sets the session storage for manager.
func (m *Manager) SetStorage(storage Storage) {
	m.storage = storage
}

// GetStorage returns the session storage of current manager.
func (m *Manager) GetStorage() Storage {
	return m.storage
}

// SetTTL the TTL for the session manager.
func (m *Manager) SetTTL(ttl time.Duration) {
	m.ttl = ttl
}

// GetTTL returns the TTL of the session manager.
func (m *Manager) GetTTL() time.Duration {
	return m.ttl
}
