// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// GroupGeneric provides generic functions of redis.
type GroupGeneric struct {
	redis *Redis
}

// GroupGeneric creates and returns GroupGeneric.
func (r *Redis) GroupGeneric() gredis.IGroupGeneric {
	return GroupGeneric{
		redis: r,
	}
}

// Copy copies the value stored at the source key to the destination key.
//
// By default, the destination key is created in the logical database used by the connection.
// The DB option allows specifying an alternative logical database index for the destination key.
//
// The command returns an error when the destination key already exists.
//
// It returns:
// - 1 if source was copied.
// - 0 if source was not copied.
//
// https://redis.io/commands/copy/
func (r GroupGeneric) Copy(ctx context.Context, source, destination string, option ...gredis.CopyOption) (int64, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "Copy", mustMergeOptionToArgs(
		[]interface{}{source, destination}, usedOption,
	)...)
	return v.Int64(), err
}

// Exists returns if key exists.
// The user should be aware that if the same existing key is mentioned in the arguments multiple times,
// it will be counted multiple times.
// So if some key exists, EXISTS some key will return 2.
//
// https://redis.io/commands/exists/
func (r GroupGeneric) Exists(ctx context.Context, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "Exists", gconv.Interfaces(keys)...)
	return v.Int64(), err
}

// Type returns the string representation of the type of the value stored at key.
// The different types that can be returned are: string, list, set, zset, hash and stream.
//
// It returns type of key, or none when key does not exist.
//
// https://redis.io/commands/type/
func (r GroupGeneric) Type(ctx context.Context, key string) (string, error) {
	v, err := r.redis.Do(ctx, "Type", key)
	return v.String(), err
}

// Unlink is very similar to DEL: it removes the specified keys. Just like DEL a key is ignored if it does not exist.
// However, the command performs the actual memory reclaiming in a different thread, so it is not blocking, while DEL is.
// This is where the command name comes from: the command just unlinks the keys from the keyspace.
// The actual removal will happen later asynchronously.
//
// It returns the number of keys that were unlinked.
//
// https://redis.io/commands/unlink/
func (r GroupGeneric) Unlink(ctx context.Context, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "Unlink", gconv.Interfaces(keys)...)
	return v.Int64(), err
}

// Rename renames key to newKey. It returns an error when key does not exist.
// If newKey already exists it is overwritten, when this happens RENAME executes an implicit DEL operation,
// so if the deleted key contains a very big value it may cause high latency even if RENAME itself is usually a constant-time operation.
//
// In Cluster mode, both key and newKey must be in the same hash slot,
// meaning that in practice only keys that have the same hashtag can be reliably renamed in cluster.
//
// https://redis.io/commands/rename/
func (r GroupGeneric) Rename(ctx context.Context, key, newKey string) error {
	_, err := r.redis.Do(ctx, "Rename", key, newKey)
	return err
}

// RenameNX renames key to newKey if newKey does not yet exist.
// It returns an error when key does not exist.
// In Cluster mode, both key and newKey must be in the same hash slot,
// meaning that in practice only keys that have the same hashtag can be reliably renamed in cluster.
//
// It returns:
// - 1 if key was renamed to newKey.
// - 0 if newKey already exists.
//
// https://redis.io/commands/renamenx/
func (r GroupGeneric) RenameNX(ctx context.Context, key, newKey string) (int64, error) {
	v, err := r.redis.Do(ctx, "RenameNX", key, newKey)
	return v.Int64(), err
}

// Move moves key from the currently selected database (see SELECT) to the specified destination database.
// When key already exists in the destination database, or it does not exist in the source database,
// it does nothing.
// It is possible to use MOVE as a locking primitive because of this.
//
// It returns:
// - 1 if key was moved.
// - 0 if key was not moved.
//
// https://redis.io/commands/move/
func (r GroupGeneric) Move(ctx context.Context, key string, db int) (int64, error) {
	v, err := r.redis.Do(ctx, "Move", key, db)
	return v.Int64(), err
}

// Del removes the specified keys.
// a key is ignored if it does not exist.
//
// It returns the number of keys that were removed.
//
// https://redis.io/commands/del/
func (r GroupGeneric) Del(ctx context.Context, keys ...string) (int64, error) {
	v, err := r.redis.Do(ctx, "Del", keys)
	return v.Int64(), err
}

// RandomKey return a random key from the currently selected database.
//
// It returns the random key, or nil when the database is empty.
//
// https://redis.io/commands/randomkey/
func (r GroupGeneric) RandomKey(ctx context.Context) (string, error) {
	v, err := r.redis.Do(ctx, "RandomKey")
	return v.String(), err
}

// DBSize return the number of keys in the currently-selected database.
//
// https://redis.io/commands/dbsize/
func (r GroupGeneric) DBSize(ctx context.Context) (int64, error) {
	v, err := r.redis.Do(ctx, "DBSize")
	return v.Int64(), err
}

// Keys return all keys matching pattern.
//
// While the time complexity for this operation is O(N), the constant times are fairly low.
// For example, Redis running on an entry level laptop can scan a 1 million key database in 40 milliseconds.
//
// https://redis.io/commands/keys/
func (r GroupGeneric) Keys(ctx context.Context, pattern string) ([]string, error) {
	v, err := r.redis.Do(ctx, "Keys", pattern)
	return v.Strings(), err
}

// FlushDB delete all the keys of the currently selected DB. This command never fails.
//
// https://redis.io/commands/flushdb/
func (r GroupGeneric) FlushDB(ctx context.Context, option ...gredis.FlushOp) error {
	_, err := r.redis.Do(ctx, "FlushDB", gconv.Interfaces(option)...)
	return err
}

// FlushAll delete all the keys of all the existing databases, not just the currently selected one.
// This command never fails.
// By default, FlushAll will synchronously flush all the databases.
//
// It is possible to use one of the following modifiers to dictate the flushing mode explicitly:
// ASYNC: flushes the databases asynchronously
// SYNC: flushes the databases synchronously
//
// Note: an asynchronous FlushAll command only deletes keys that were present at the time the command was invoked.
// Keys created during an asynchronous flush will be unaffected.
//
// https://redis.io/commands/flushall/
func (r GroupGeneric) FlushAll(ctx context.Context, option ...gredis.FlushOp) error {
	_, err := r.redis.Do(ctx, "FlushAll", gconv.Interfaces(option)...)
	return err
}

// Expire sets a timeout on key.
// After the timeout has expired, the key will automatically be deleted.
//
// It returns:
// - 1 if the timeout was set.
// - 0 if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments.
//
// https://redis.io/commands/expire/
func (r GroupGeneric) Expire(ctx context.Context, key string, seconds time.Duration, option ...gredis.ExpireOption) (int64, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "Expire", mustMergeOptionToArgs(
		[]interface{}{key, seconds.Seconds()}, usedOption,
	)...)
	return v.Int64(), err
}

// ExpireAt has the same effect and semantic as EXPIRE, but instead of specifying the number of
// seconds representing the TTL (time to live), it takes an absolute Unix timestamp (seconds since
// January 1, 1970).
// A timestamp in the past will delete the key immediately.
//
// It returns:
// - 1 if the timeout was set.
// - 0 if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments.
//
// https://redis.io/commands/expireat/
func (r GroupGeneric) ExpireAt(ctx context.Context, key string, time time.Time, option ...gredis.ExpireOption) (int64, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "ExpireAt", mustMergeOptionToArgs(
		[]interface{}{key, gtime.New(time).Timestamp()}, usedOption,
	)...)
	return v.Int64(), err
}

// ExpireTime returns the absolute time at which the given key will expire.
//
// It returns:
// - -1 if the key exists but has no associated expiration time.
// - -2 if the key does not exist.
//
// https://redis.io/commands/expiretime/
func (r GroupGeneric) ExpireTime(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "ExpireTime", key)
}

// TTL returns the remaining time to live of a key that has a timeout.
// This introspection capability allows a Redis client to check how many seconds a given key
// will continue to be part of the dataset.
// In Redis 2.6 or older the command returns -1 if the key does not exist or if the key exist but has
// no associated expire.
//
// Starting with Redis 2.8 the return value in case of error changed:
//
// The command returns -2 if the key does not exist.
// The command returns -1 if the key exists but has no associated expire.
// See also the PTTL command that returns the same information with milliseconds resolution
// (Only available in Redis 2.6 or greater).
//
// It returns TTL in seconds, or a negative value in order to signal an error (see the description above).
//
// https://redis.io/commands/ttl/
func (r GroupGeneric) TTL(ctx context.Context, key string) (time.Duration, error) {
	v, err := r.redis.Do(ctx, "TTL", key)
	return v.Duration(), err
}

// Persist removes the existing timeout on key, turning the key from volatile (a key with an expire set)
// to persistent (a key that will never expire as no timeout is associated).
//
// It returns:
// - 1 if the timeout was removed.
// - 0 if key does not exist or does not have an associated timeout.
//
// https://redis.io/commands/persist/
func (r GroupGeneric) Persist(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "Persist", key)
	return v.Int64(), err
}

// PExpire works exactly like EXPIRE but the time to live of the key is specified in milliseconds
// instead of seconds.
//
// It returns:
// - 1 if the timeout was set.
// - 0 if the timeout was not set. e.g. key doesn't exist, or operation skipped due to the provided arguments.
//
// https://redis.io/commands/pexpire/
func (r GroupGeneric) PExpire(ctx context.Context, key string, milliseconds time.Duration, option ...gredis.ExpireOption) (int64, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "PExpire", mustMergeOptionToArgs(
		[]interface{}{key, milliseconds.Milliseconds()}, usedOption,
	)...)
	return v.Int64(), err
}

// PExpireAt has the same effect and semantic as ExpireAt, but the Unix time at which the key will
// expire is specified in milliseconds instead of seconds.
//
// https://redis.io/commands/pexpireat/
func (r GroupGeneric) PExpireAt(ctx context.Context, key string, time time.Time, option ...gredis.ExpireOption) (int64, error) {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	v, err := r.redis.Do(ctx, "PExpireAt", mustMergeOptionToArgs(
		[]interface{}{key, gtime.New(time).TimestampMilli()}, usedOption,
	)...)
	return v.Int64(), err
}

// PExpireTime returns the expiration time of given `key`.
//
// It returns:
// - -1 if the key exists but has no associated expiration time.
// - -2 if the key does not exist.
//
// https://redis.io/commands/pexpiretime/
func (r GroupGeneric) PExpireTime(ctx context.Context, key string) (*gvar.Var, error) {
	return r.redis.Do(ctx, "PExpireTime", key)
}

// PTTL like TTL this command returns the remaining time to live of a key that has an expired set,
// with the sole difference that TTL returns the amount of remaining time in seconds while PTTL
// returns it in milliseconds.
//
// In Redis 2.6 or older the command returns -1 if the key does not exist or if the key exist but has
// no associated expire.
//
// It returns TTL in milliseconds, or a negative value in order to signal an error (see the description above).
//
//  https://redis.io/commands/pttl/
func (r GroupGeneric) PTTL(ctx context.Context, key string) (int64, error) {
	v, err := r.redis.Do(ctx, "PTTL", key)
	return v.Int64(), err
}
