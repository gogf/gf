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

// StorageBase is a base implement for Session Storage.
type StorageBase struct{}

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageBase) New(ctx context.Context, ttl time.Duration) (id string, err error) {
	return "", ErrorDisabled
}

// Get retrieves certain session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageBase) Get(ctx context.Context, sessionId string, key string) (value interface{}, err error) {
	return nil, ErrorDisabled
}

// Data retrieves all key-value pairs as map from storage.
func (s *StorageBase) Data(ctx context.Context, sessionId string) (sessionData map[string]interface{}, err error) {
	return nil, ErrorDisabled
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageBase) GetSize(ctx context.Context, sessionId string) (size int, err error) {
	return 0, ErrorDisabled
}

// Set sets key-value session pair to the storage.
// The parameter `ttl` specifies the TTL for the session id (not for the key-value pair).
func (s *StorageBase) Set(ctx context.Context, sessionId string, key string, value interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter `ttl` specifies the TTL for the session id(not for the key-value pair).
func (s *StorageBase) SetMap(ctx context.Context, sessionId string, mapData map[string]interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// Remove deletes key with its value from storage.
func (s *StorageBase) Remove(ctx context.Context, sessionId string, key string) error {
	return ErrorDisabled
}

// RemoveAll deletes session from storage.
func (s *StorageBase) RemoveAll(ctx context.Context, sessionId string) error {
	return ErrorDisabled
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageBase) GetSession(ctx context.Context, sessionId string, ttl time.Duration) (*gmap.StrAnyMap, error) {
	return nil, ErrorDisabled
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageBase) SetSession(ctx context.Context, sessionId string, sessionData *gmap.StrAnyMap, ttl time.Duration) error {
	return ErrorDisabled
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageBase) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	return ErrorDisabled
}
