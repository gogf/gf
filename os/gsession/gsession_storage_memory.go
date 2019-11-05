// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"github.com/gogf/gf/container/gmap"
	"time"
)

// StorageMemory implements the Session Storage interface with memory.
type StorageMemory struct{}

// NewStorageMemory creates and returns a file storage object for session.
func NewStorageMemory() *StorageMemory {
	return &StorageMemory{}
}

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageMemory) New(ttl time.Duration) (id string) {
	return ""
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageMemory) Get(id string, key string) interface{} {
	return nil
}

// GetMap retrieves all key-value pairs as map from storage.
func (s *StorageMemory) GetMap(id string) map[string]interface{} {
	return nil
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageMemory) GetSize(id string) int {
	return -1
}

// Set sets key-value session pair to the storage.
// The parameter <ttl> specifies the TTL for the session id (not for the key-value pair).
func (s *StorageMemory) Set(id string, key string, value interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter <ttl> specifies the TTL for the session id(not for the key-value pair).
func (s *StorageMemory) SetMap(id string, data map[string]interface{}, ttl time.Duration) error {
	return ErrorDisabled
}

// Remove deletes key with its value from storage.
func (s *StorageMemory) Remove(id string, key string) error {
	return ErrorDisabled
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageMemory) RemoveAll(id string) error {
	return ErrorDisabled
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter <ttl> specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter <data> is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageMemory) GetSession(id string, ttl time.Duration, data *gmap.StrAnyMap) (*gmap.StrAnyMap, error) {
	return data, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageMemory) SetSession(id string, data *gmap.StrAnyMap, ttl time.Duration) error {
	return nil
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageMemory) UpdateTTL(id string, ttl time.Duration) error {
	return nil
}

// doUpdateTTL updates the TTL for session id.
func (s *StorageMemory) doUpdateTTL(id string) error {
	return nil
}
