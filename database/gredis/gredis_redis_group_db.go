// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
	"github.com/gogf/gf/v2/util/gconv"
)

type RedisGroupDB struct {
	redis *Redis
}

func (r *Redis) DB() *RedisGroupDB {
	return &RedisGroupDB{
		redis: r,
	}
}

// Exists return if key exists.
// The user should be aware that if the same existing key is mentioned in the arguments multiple times,
// it will be counted multiple times.
// So if somekey exists, EXISTS somekey somekey will return 2.
//
// https://redis.io/commands/exists/
func (r *RedisGroupDB) Exists(ctx context.Context, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "EXISTS", keys)
	return v.Int64(), err
}

// Type return the string representation of the type of the value stored at key.
// The different types that can be returned are: string, list, set, zset, hash and stream.
//
// https://redis.io/commands/type/
func (r *RedisGroupDB) Type(ctx context.Context, key string) (string, error) {
	v, err := r.redis.Do(ctx, "TYPE", key)
	return v.String(), err
}

// Rename key to newkey. It return an error when key does not exist.
// If newkey already exists it is overwritten, when this happens RENAME executes an implicit DEL operation,
// so if the deleted key contains a very big value it may cause high latency even if RENAME itself is usually a constant-time operation.
//
// In Cluster mode, both key and newkey must be in the same hash slot,
// meaning that in practice only keys that have the same hash tag can be reliably renamed in cluster.
//
// https://redis.io/commands/rename/
func (r *RedisGroupDB) Rename(ctx context.Context, key, newKey string) (string, error) {
	v, err := r.redis.Do(ctx, "RENAME", key, newKey)
	return v.String(), err
}

// RenameNX renames key to newkey if newkey does not yet exist.
// It return an error when key does not exist.
// In Cluster mode, both key and newkey must be in the same hash slot,
// meaning that in practice only keys that have the same hash tag can be reliably renamed in cluster.
//
// https://redis.io/commands/renamenx/
func (r *RedisGroupDB) RenameNX(ctx context.Context, key, newKey string) (bool, error) {
	v, err := r.redis.Do(ctx, "RENAME", key, newKey)
	return v.Bool(), err
}

// Move  key from the currently selected database (see SELECT) to the specified destination database.
// When key already exists in the destination database, or it does not exist in the source database,
// it does nothing.
// It is possible to use MOVE as a locking primitive because of this.
//
// https://redis.io/commands/move/
func (r *RedisGroupDB) Move(ctx context.Context, key, db string) (bool, error) {
	v, err := r.redis.Do(ctx, "MOVE", key, db)
	return v.Bool(), err
}

// Del removes the specified keys.
// a key is ignored if it does not exist.
//
// https://redis.io/commands/del/
func (r *RedisGroupDB) Del(ctx context.Context, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "DEL", keys)
	return v.Int64(), err
}

// RandomKey return a random key from the currently selected database.
//
// https://redis.io/commands/randomkey/
func (r *RedisGroupDB) RandomKey(ctx context.Context) (string, error) {
	v, err := r.redis.Do(ctx, "RANDOMKEY")
	return v.String(), err
}

// DBSize return the number of keys in the currently-selected database.
//
// https://redis.io/commands/dbsize/
func (r *RedisGroupDB) DBSize(ctx context.Context) (int64, error) {
	v, err := r.redis.Do(ctx, "DBSIZE")
	return v.Int64(), err
}

// Keys return all keys matching pattern.
// While the time complexity for this operation is O(N), the constant times are fairly low.
// For example, Redis running on an entry level laptop can scan a 1 million key database in 40 milliseconds.
// consider KEYS as a command that should only be used in production environments with extreme care.
// It may ruin performance when it is executed against large databases.
// This command is intended for debugging and special operations, such as changing your keyspace layout.
// Don't use KEYS in your regular application code.
// If you're looking for a way to find keys in a subset of your keyspace, consider using SCAN or sets.
//
// https://redis.io/commands/keys/
func (r *RedisGroupDB) Keys(ctx context.Context, pattern string) ([]string, error) {
	v, err := r.redis.Do(ctx, "KEYS", pattern)
	return gconv.SliceStr(v), err
}

// FlushDB delete all the keys of the currently selected DB. This command never fails.
//
// By default, FLUSHDB will synchronously flush all keys from the database.
// Starting with Redis 6.2, setting the lazyfree-lazy-user-flush configuration directive to "yes"
// changes the default flush mode to asynchronous.
// It is possible to use one of the following modifiers to dictate the flushing mode explicitly:
// ASYNC: flushes the database asynchronously
// SYNC: flushes the database synchronously
// Note: an asynchronous FLUSHDB command only deletes keys that were present at the time the command was invoked.
// Keys created during an asynchronous flush will be unaffected.
//
// https://redis.io/commands/flushdb/
func (r *RedisGroupDB) FlushDB(ctx context.Context, options string) error {
	_, err := r.redis.Do(ctx, "FLUSHDB", options)
	return err
}

// FlushAll delete all the keys of all the existing databases, not just the currently selected one.
// This command never fails.
// By default, FLUSHALL will synchronously flush all the databases.
// Starting with Redis 6.2, setting the lazyfree-lazy-user-flush configuration directive to "yes" changes the
// default flush mode to asynchronous.
// It is possible to use one of the following modifiers to dictate the flushing mode explicitly:
// ASYNC: flushes the databases asynchronously
// SYNC: flushes the databases synchronously
// Note: an asynchronous FLUSHALL command only deletes keys that were present at the time the command was invoked.
// Keys created during an asynchronous flush will be unaffected.
//
// https://redis.io/commands/flushall/
func (r *RedisGroupDB) FlushAll(ctx context.Context, options string) error {
	_, err := r.redis.Do(ctx, "FLUSHALL", options)
	return err
}
