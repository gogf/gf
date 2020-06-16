// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/util/gconv"
)

// StorageRedisHashTable implements the Session Storage interface with redis hash table.
type StorageRedisHashTable struct {
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

// New creates a session id.
// This function can be used for custom session creation.
func (s *StorageRedisHashTable) New(ttl time.Duration) (id string) {
	return ""
}

// Get retrieves session value with given key.
// It returns nil if the key does not exist in the session.
func (s *StorageRedisHashTable) Get(id string, key string) interface{} {
	r, _ := s.redis.Do("HGET", s.key(id), key)
	if r != nil {
		return gconv.String(r)
	}
	return r
}

// GetMap retrieves all key-value pairs as map from storage.
func (s *StorageRedisHashTable) GetMap(id string) map[string]interface{} {
	r, err := s.redis.DoVar("HGETALL", s.key(id))
	if err != nil {
		return nil
	}
	array := r.Interfaces()
	m := make(map[string]interface{})
	for i := 0; i < len(array); i += 2 {
		if array[i+1] != nil {
			m[gconv.String(array[i])] = gconv.String(array[i+1])
		} else {
			m[gconv.String(array[i])] = array[i+1]
		}
	}
	return m
}

// GetSize retrieves the size of key-value pairs from storage.
func (s *StorageRedisHashTable) GetSize(id string) int {
	r, _ := s.redis.DoVar("HLEN", s.key(id))
	return r.Int()
}

// Set sets key-value session pair to the storage.
// The parameter <ttl> specifies the TTL for the session id (not for the key-value pair).
func (s *StorageRedisHashTable) Set(id string, key string, value interface{}, ttl time.Duration) error {
	_, err := s.redis.Do("HSET", s.key(id), key, value)
	return err
}

// SetMap batch sets key-value session pairs with map to the storage.
// The parameter <ttl> specifies the TTL for the session id(not for the key-value pair).
func (s *StorageRedisHashTable) SetMap(id string, data map[string]interface{}, ttl time.Duration) error {
	array := make([]interface{}, len(data)*2+1)
	array[0] = s.key(id)

	index := 1
	for k, v := range data {
		array[index] = k
		array[index+1] = v
		index += 2
	}
	_, err := s.redis.Do("HMSET", array...)
	return err
}

// Remove deletes key with its value from storage.
func (s *StorageRedisHashTable) Remove(id string, key string) error {
	_, err := s.redis.Do("HDEL", s.key(id), key)
	return err
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedisHashTable) RemoveAll(id string) error {
	_, err := s.redis.Do("DEL", s.key(id))
	return err
}

// GetSession returns the session data as *gmap.StrAnyMap for given session id from storage.
//
// The parameter <ttl> specifies the TTL for this session, and it returns nil if the TTL is exceeded.
// The parameter <data> is the current old session data stored in memory,
// and for some storage it might be nil if memory storage is disabled.
//
// This function is called ever when session starts.
func (s *StorageRedisHashTable) GetSession(id string, ttl time.Duration, data *gmap.StrAnyMap) (*gmap.StrAnyMap, error) {
	intlog.Printf("StorageRedisHashTable.GetSession: %s, %v", id, ttl)
	r, err := s.redis.DoVar("EXISTS", s.key(id))
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
func (s *StorageRedisHashTable) SetSession(id string, data *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf("StorageRedisHashTable.SetSession: %s, %v", id, ttl)
	_, err := s.redis.Do("EXPIRE", s.key(id), int64(ttl.Seconds()))
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedisHashTable) UpdateTTL(id string, ttl time.Duration) error {
	intlog.Printf("StorageRedisHashTable.UpdateTTL: %s, %v", id, ttl)
	_, err := s.redis.Do("EXPIRE", s.key(id), int64(ttl.Seconds()))
	return err
}

func (s *StorageRedisHashTable) key(id string) string {
	return s.prefix + id
}
