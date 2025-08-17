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
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/internal/intlog"
)

// StorageRedisHashTable implements the Session Storage interface with redis hash table.
type StorageRedisHashTable struct {
	StorageBase
	redis  *gredis.Redis // Redis client for session storage.
	prefix string        // Redis key prefix for session id.
}

// NewStorageRedisHashTable creates and returns a redis hash table storage object for session.
func NewStorageRedisHashTable(redis *gredis.Redis, prefix ...string) *StorageRedisHashTable {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedisHashTable{
		redis: redis,
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	return s
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageRedisHashTable) Get(ctx context.Context, sessionId string, key string) (value interface{}, err error) {
	v, err := s.redis.HGet(ctx, s.sessionIdToRedisKey(sessionId), key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	return v.String(), nil
}

// Data retrieves all key-value pairs as map from storage.
func (s *StorageRedisHashTable) Data(ctx context.Context, sessionId string) (data map[string]interface{}, err error) {
	m, err := s.redis.HGetAll(ctx, s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	return m.Map(), nil
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageRedisHashTable) GetSize(ctx context.Context, sessionId string) (size int, err error) {
	v, err := s.redis.HLen(ctx, s.sessionIdToRedisKey(sessionId))
	return int(v), err
}

// Set sets key-value session pair to the storage.
// The parameter `ttl` specifies the TTL for the session id (not for the key-value pair).
func (s *StorageRedisHashTable) Set(ctx context.Context, sessionId string, key string, value interface{}, ttl time.Duration) error {
	_, err := s.redis.HSet(ctx, s.sessionIdToRedisKey(sessionId), map[string]interface{}{
		key: value,
	})
	return err
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter `ttl` specifies the TTL for the session id(not for the key-value pair).
func (s *StorageRedisHashTable) SetMap(ctx context.Context, sessionId string, data map[string]interface{}, ttl time.Duration) error {
	err := s.redis.HMSet(ctx, s.sessionIdToRedisKey(sessionId), data)
	return err
}

// Remove deletes key with its value from storage.
func (s *StorageRedisHashTable) Remove(ctx context.Context, sessionId string, key string) error {
	_, err := s.redis.HDel(ctx, s.sessionIdToRedisKey(sessionId), key)
	return err
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedisHashTable) RemoveAll(ctx context.Context, sessionId string) error {
	_, err := s.redis.Del(ctx, s.sessionIdToRedisKey(sessionId))
	return err
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageRedisHashTable) GetSession(ctx context.Context, sessionId string, ttl time.Duration) (*gmap.StrAnyMap, error) {
	intlog.Printf(ctx, "StorageRedisHashTable.GetSession: %s, %v", sessionId, ttl)
	v, err := s.redis.Exists(ctx, s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	if v > 0 {
		// It does not store the session data in memory, it so returns an empty map.
		// It retrieves session data items directly through redis server each time.
		return gmap.NewStrAnyMap(true), nil
	}
	return nil, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedisHashTable) SetSession(ctx context.Context, sessionId string, sessionData *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedisHashTable.SetSession: %s, %v", sessionId, ttl)
	_, err := s.redis.Expire(ctx, s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedisHashTable) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedisHashTable.UpdateTTL: %s, %v", sessionId, ttl)
	_, err := s.redis.Expire(ctx, s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

// sessionIdToRedisKey converts and returns the redis key for given session id.
func (s *StorageRedisHashTable) sessionIdToRedisKey(sessionId string) string {
	return s.prefix + sessionId
}
