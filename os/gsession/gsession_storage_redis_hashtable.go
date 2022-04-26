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
	"github.com/gogf/gf/v2/util/gconv"
)

// StorageRedisHashTable implements the Session Storage interface with redis hash table.
type StorageRedisHashTable struct {
	redis  *gredis.Redis // Redis client for session storage.
	prefix string        // Redis sessionIdToRedisKey prefix for session id.
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

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageRedisHashTable) New(ctx context.Context, ttl time.Duration) (id string, err error) {
	return "", ErrorDisabled
}

// Get retrieves session value with given sessionIdToRedisKey.
// It returns nil if the sessionIdToRedisKey does not exist in the session.
func (s *StorageRedisHashTable) Get(ctx context.Context, sessionId string, key string) (value interface{}, err error) {
	v, err := s.redis.Do(ctx, "HGET", s.sessionIdToRedisKey(sessionId), key)
	if err != nil {
		return nil, err
	}
	if v.IsNil() {
		return nil, nil
	}
	return v.String(), nil
}

// Data retrieves all sessionIdToRedisKey-value pairs as map from storage.
func (s *StorageRedisHashTable) Data(ctx context.Context, sessionId string) (data map[string]interface{}, err error) {
	v, err := s.redis.Do(ctx, "HGETALL", s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	data = make(map[string]interface{})
	array := v.Interfaces()
	for i := 0; i < len(array); i += 2 {
		if array[i+1] != nil {
			data[gconv.String(array[i])] = gconv.String(array[i+1])
		} else {
			data[gconv.String(array[i])] = array[i+1]
		}
	}
	return data, nil
}

// GetSize retrieves the size of sessionIdToRedisKey-value pairs from storage.
func (s *StorageRedisHashTable) GetSize(ctx context.Context, sessionId string) (size int, err error) {
	r, err := s.redis.Do(ctx, "HLEN", s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return -1, err
	}
	return r.Int(), nil
}

// Set sets sessionIdToRedisKey-value session pair to the storage.
// The parameter `ttl` specifies the TTL for the session id (not for the sessionIdToRedisKey-value pair).
func (s *StorageRedisHashTable) Set(ctx context.Context, sessionId string, key string, value interface{}, ttl time.Duration) error {
	_, err := s.redis.Do(ctx, "HSET", s.sessionIdToRedisKey(sessionId), key, value)
	return err
}

// SetMap batch sets sessionIdToRedisKey-value session pairs with map to the storage.
// The parameter `ttl` specifies the TTL for the session id(not for the sessionIdToRedisKey-value pair).
func (s *StorageRedisHashTable) SetMap(ctx context.Context, sessionId string, data map[string]interface{}, ttl time.Duration) error {
	array := make([]interface{}, len(data)*2+1)
	array[0] = s.sessionIdToRedisKey(sessionId)

	index := 1
	for k, v := range data {
		array[index] = k
		array[index+1] = v
		index += 2
	}
	_, err := s.redis.Do(ctx, "HMSET", array...)
	return err
}

// Remove deletes sessionIdToRedisKey with its value from storage.
func (s *StorageRedisHashTable) Remove(ctx context.Context, sessionId string, key string) error {
	_, err := s.redis.Do(ctx, "HDEL", s.sessionIdToRedisKey(sessionId), key)
	return err
}

// RemoveAll deletes all sessionIdToRedisKey-value pairs from storage.
func (s *StorageRedisHashTable) RemoveAll(ctx context.Context, sessionId string) error {
	_, err := s.redis.Do(ctx, "DEL", s.sessionIdToRedisKey(sessionId))
	return err
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter `ttl` specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter `data` is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageRedisHashTable) GetSession(ctx context.Context, sessionId string, ttl time.Duration, data *gmap.StrAnyMap) (*gmap.StrAnyMap, error) {
	intlog.Printf(ctx, "StorageRedisHashTable.GetSession: %s, %v", sessionId, ttl)
	r, err := s.redis.Do(ctx, "EXISTS", s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	if r.Bool() {
		return gmap.NewStrAnyMap(true), nil
	}
	return nil, nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedisHashTable) SetSession(ctx context.Context, sessionId string, data *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedisHashTable.SetSession: %s, %v", sessionId, ttl)
	_, err := s.redis.Do(ctx, "EXPIRE", s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedisHashTable) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedisHashTable.UpdateTTL: %s, %v", sessionId, ttl)
	_, err := s.redis.Do(ctx, "EXPIRE", s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()))
	return err
}

func (s *StorageRedisHashTable) sessionIdToRedisKey(sessionId string) string {
	return s.prefix + sessionId
}
