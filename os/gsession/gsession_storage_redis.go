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
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtimer"
)

// StorageRedis implements the Session Storage interface with redis.
type StorageRedis struct {
	StorageBase
	redis         *gredis.Redis   // Redis client for session storage.
	prefix        string          // Redis key prefix for session id.
	updatingIdMap *gmap.StrIntMap // Updating TTL set for session id.
}

const (
	// DefaultStorageRedisLoopInterval is the interval updating TTL for session ids
	// in last duration.
	DefaultStorageRedisLoopInterval = 10 * time.Second
)

// NewStorageRedis creates and returns a redis storage object for session.
func NewStorageRedis(redis *gredis.Redis, prefix ...string) *StorageRedis {
	if redis == nil {
		panic("redis instance for storage cannot be empty")
		return nil
	}
	s := &StorageRedis{
		redis:         redis,
		updatingIdMap: gmap.NewStrIntMap(true),
	}
	if len(prefix) > 0 && prefix[0] != "" {
		s.prefix = prefix[0]
	}
	// Batch updates the TTL for session ids timely.
	gtimer.AddSingleton(context.Background(), DefaultStorageRedisLoopInterval, func(ctx context.Context) {
		intlog.Print(context.TODO(), "StorageRedis.timer start")
		var (
			err        error
			sessionId  string
			ttlSeconds int
		)
		for {
			if sessionId, ttlSeconds = s.updatingIdMap.Pop(); sessionId == "" {
				break
			} else {
				if err = s.doUpdateExpireForSession(context.TODO(), sessionId, ttlSeconds); err != nil {
					intlog.Errorf(context.TODO(), `%+v`, err)
				}
			}
		}
		intlog.Print(context.TODO(), "StorageRedis.timer end")
	})
	return s
}

// RemoveAll deletes all key-value pairs from storage.
func (s *StorageRedis) RemoveAll(ctx context.Context, sessionId string) error {
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
func (s *StorageRedis) GetSession(ctx context.Context, sessionId string, ttl time.Duration) (*gmap.StrAnyMap, error) {
	intlog.Printf(ctx, "StorageRedis.GetSession: %s, %v", sessionId, ttl)
	r, err := s.redis.Do(ctx, "GET", s.sessionIdToRedisKey(sessionId))
	if err != nil {
		return nil, err
	}
	content := r.Bytes()
	if len(content) == 0 {
		return nil, nil
	}
	var m map[string]interface{}
	if err = json.UnmarshalUseNumber(content, &m); err != nil {
		return nil, err
	}
	if m == nil {
		return nil, nil
	}
	return gmap.NewStrAnyMapFrom(m, true), nil
}

// SetSession updates the data map for specified session id.
// This function is called ever after session, which is changed dirty, is closed.
// This copy all session data map from memory to storage.
func (s *StorageRedis) SetSession(ctx context.Context, sessionId string, sessionData *gmap.StrAnyMap, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedis.SetSession: %s, %v, %v", sessionId, sessionData, ttl)
	content, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}
	_, err = s.redis.Do(ctx, "SETEX", s.sessionIdToRedisKey(sessionId), int64(ttl.Seconds()), content)
	return err
}

// UpdateTTL updates the TTL for specified session id.
// This function is called ever after session, which is not dirty, is closed.
// It just adds the session id to the async handling queue.
func (s *StorageRedis) UpdateTTL(ctx context.Context, sessionId string, ttl time.Duration) error {
	intlog.Printf(ctx, "StorageRedis.UpdateTTL: %s, %v", sessionId, ttl)
	if ttl >= DefaultStorageRedisLoopInterval {
		s.updatingIdMap.Set(sessionId, int(ttl.Seconds()))
	}
	return nil
}

// doUpdateTTL updates the TTL for session id.
func (s *StorageRedis) doUpdateExpireForSession(ctx context.Context, sessionId string, ttlSeconds int) error {
	intlog.Printf(ctx, "StorageRedis.doUpdateTTL: %s, %d", sessionId, ttlSeconds)
	_, err := s.redis.Do(ctx, "EXPIRE", s.sessionIdToRedisKey(sessionId), ttlSeconds)
	return err
}

// sessionIdToRedisKey converts and returns the redis key for given session id.
func (s *StorageRedis) sessionIdToRedisKey(sessionId string) string {
	return s.prefix + sessionId
}
