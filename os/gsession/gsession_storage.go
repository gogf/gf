// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
)

// Storage is the interface definition for session storage.
type Storage interface {
	// New creates a custom session id.
	// This function can be used for custom session creation.
	New(ctx context.Context, ttl time.Duration) (sessionId string, err error)

	// Get retrieves and returns certain session value with given key.
	// It returns nil if the key does not exist in the session.
	Get(ctx context.Context, sessionId string, key string) (value interface{}, err error)

	// GetSize retrieves and returns the size of key-value pairs from storage.
	GetSize(ctx context.Context, sessionId string) (size int, err error)

	// Data retrieves all key-value pairs as map from storage.
	Data(ctx context.Context, sessionId string) (sessionData map[string]interface{}, err error)

	// Set sets one key-value session pair to the storage.
	// The parameter `ttl` specifies the TTL for the session id.
	Set(ctx context.Context, sessionId string, key string, value interface{}, ttl time.Duration) error

	// SetMap batch sets key-value session pairs as map to the storage.
	// The parameter `ttl` specifies the TTL for the session id.
	SetMap(ctx context.Context, sessionId string, mapData map[string]interface{}, ttl time.Duration) error

	// Remove deletes key-value pair from specified session from storage.
	Remove(ctx context.Context, sessionId string, key string) error

	// RemoveAll deletes session from storage.
	RemoveAll(ctx context.Context, sessionId string) error

	// GetSession returns the session data as `*gmap.StrAnyMap` for given session from storage.
	//
	// The parameter `ttl` specifies the TTL for this session.
	// The parameter `data` is the current old session data stored in memory,
	// and for some storage it might be nil if memory storage is disabled.
	//
	// This function is called ever when session starts.
	// It returns nil if the session does not exist or its TTL is expired.
	GetSession(ctx context.Context, sessionId string, ttl time.Duration) (*gmap.StrAnyMap, error)

	// SetSession updates the data for specified session id.
	// This function is called ever after session, which is changed dirty, is closed.
	// This copy all session data map from memory to storage.
	SetSession(ctx context.Context, sessionId string, sessionData *gmap.StrAnyMap, ttl time.Duration) error

	// UpdateTTL updates the TTL for specified session id.
	// This function is called ever after session, which is not dirty, is closed.
	UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error
}
