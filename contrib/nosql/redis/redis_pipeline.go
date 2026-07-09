// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file implements the gredis.Pipeliner interface by wrapping go-redis's
// redis.Pipeliner. It provides Pipeline, TxPipeline, and Watch operations,
// along with typed command groups that queue commands and return *gredis.Cmd
// futures. Results are populated after Exec is called.

package redis

import (
	"context"
	"errors"
	"reflect"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// cmdEntry pairs a gredis.Cmd future with a populate closure that fills its
// result after the underlying go-redis command has been executed by Exec.
type cmdEntry struct {
	// cmd is the gredis.Cmd future returned to the caller.
	cmd *gredis.Cmd

	// populate reads the result from the underlying go-redis command and
	// converts it to a *gvar.Var. It is called once during Exec.
	populate func() (*gvar.Var, error)
}

// redisPipeliner implements gredis.Pipeliner by wrapping go-redis's
// redis.Pipeliner. Each queued command creates a cmdEntry holding a
// gredis.Cmd future and a populate closure. Exec sends all queued
// commands in a single round-trip, then iterates cmdEntries to populate
// each Cmd's result.
type redisPipeliner struct {
	// pipe is the underlying go-redis pipeliner (regular or transactional).
	pipe redis.Pipeliner

	// cmdList holds all queued commands in insertion order.
	cmdList []*cmdEntry
}

// newRedisPipeliner creates a redisPipeliner wrapping the given go-redis pipeliner.
func newRedisPipeliner(pipe redis.Pipeliner) *redisPipeliner {
	return &redisPipeliner{
		pipe: pipe,
	}
}

// Do queues a raw Redis command and returns its future Cmd.
// The command is not sent to the server until Exec is called.
// Struct, map, and slice values (except []byte) are JSON-serialized
// to match the behavior of Conn.Do.
func (p *redisPipeliner) Do(ctx context.Context, command string, args ...any) *gredis.Cmd {
	// Serialize struct/map/slice values to JSON, matching Conn.Do behavior.
	for k, v := range args {
		reflectInfo := gutil.OriginTypeAndKind(v)
		switch reflectInfo.OriginKind {
		case reflect.Struct, reflect.Map, reflect.Slice, reflect.Array:
			if _, ok := v.([]byte); !ok {
				marshaled, err := gjson.Marshal(v)
				if err != nil {
					cmd := &gredis.Cmd{}
					cmd.SetErr(err)
					return cmd
				}
				args[k] = marshaled
			}
		}
	}

	goCmd := p.pipe.Do(ctx, append([]any{command}, args...)...)
	cmd := &gredis.Cmd{}
	p.cmdList = append(p.cmdList, &cmdEntry{
		cmd: cmd,
		populate: func() (*gvar.Var, error) {
			return pipelineResultToVar(goCmd.Result())
		},
	})
	return cmd
}

// Exec sends all queued commands to the server in a single batch,
// populates each queued Cmd's result, and returns all Cmds.
// After Exec returns, all Cmd objects returned by queued commands are populated.
func (p *redisPipeliner) Exec(ctx context.Context) ([]*gredis.Cmd, error) {
	_, err := p.pipe.Exec(ctx)
	// Allow redis.Nil as non-error for pipeline results.
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	for _, entry := range p.cmdList {
		val, populateErr := entry.populate()
		entry.cmd.SetVal(val)
		if populateErr != nil && !errors.Is(populateErr, redis.Nil) {
			entry.cmd.SetErr(populateErr)
		}
	}

	result := make([]*gredis.Cmd, len(p.cmdList))
	for i, entry := range p.cmdList {
		result[i] = entry.cmd
	}
	return result, nil
}

// Discard discards all queued commands without sending them to the server.
func (p *redisPipeliner) Discard() {
	p.pipe.Discard()
	p.cmdList = p.cmdList[:0]
}

// PipelineGroupGeneric returns the generic command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupGeneric() gredis.IPipelineGroupGeneric {
	return pipelineGroupGeneric{pipeliner: p}
}

// PipelineGroupHash returns the hash command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupHash() gredis.IPipelineGroupHash {
	return pipelineGroupHash{pipeliner: p}
}

// PipelineGroupString returns the string command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupString() gredis.IPipelineGroupString {
	return pipelineGroupString{pipeliner: p}
}

// PipelineGroupList returns the list command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupList() gredis.IPipelineGroupList {
	return pipelineGroupList{pipeliner: p}
}

// PipelineGroupSet returns the set command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupSet() gredis.IPipelineGroupSet {
	return pipelineGroupSet{pipeliner: p}
}

// PipelineGroupSortedSet returns the sorted set command group for pipeline mode.
func (p *redisPipeliner) PipelineGroupSortedSet() gredis.IPipelineGroupSortedSet {
	return pipelineGroupSortedSet{pipeliner: p}
}

// Pipeline returns a Pipeliner for batching multiple commands into a single
// network round-trip. Commands are buffered locally and sent when Exec is called.
func (r *Redis) Pipeline(_ context.Context) gredis.Pipeliner {
	return newRedisPipeliner(r.client.Pipeline())
}

// TxPipeline returns a Pipeliner that wraps commands in a MULTI/EXEC transaction.
// All queued commands are executed atomically by the Redis server.
func (r *Redis) TxPipeline(_ context.Context) gredis.Pipeliner {
	return newRedisPipeliner(r.client.TxPipeline())
}

// Watch watches the given keys for modifications and executes fn in a transaction.
// If any watched key is modified by another client before the transaction executes,
// the transaction is aborted and Watch returns a transaction-abort error.
func (r *Redis) Watch(ctx context.Context, fn func(gredis.Tx) error, keys ...string) error {
	return r.client.Watch(ctx, func(tx *redis.Tx) error {
		pipeliner := newRedisPipeliner(tx.TxPipeline())
		return fn(pipeliner)
	}, keys...)
}

// pipelineResultToVar converts a go-redis command result to a *gvar.Var,
// handling redis.Nil and common type conversions. This mirrors the logic
// in Conn.resultToVar.
func pipelineResultToVar(result any, err error) (*gvar.Var, error) {
	if err == redis.Nil {
		err = nil
	}
	if err == nil {
		switch v := result.(type) {
		case []byte:
			return gvar.New(string(v)), err

		case []any:
			return gvar.New(gconv.Strings(v)), err

		case *redis.Message:
			result = &gredis.Message{
				Channel:      v.Channel,
				Pattern:      v.Pattern,
				Payload:      v.Payload,
				PayloadSlice: v.PayloadSlice,
			}

		case *redis.Subscription:
			result = &gredis.Subscription{
				Kind:    v.Kind,
				Channel: v.Channel,
				Count:   v.Count,
			}
		}
	}
	return gvar.New(result), err
}

// pipelineGroupGeneric implements gredis.IPipelineGroupGeneric.
type pipelineGroupGeneric struct {
	pipeliner *redisPipeliner
}

// Copy queues a Copy command and returns its future Cmd.
// https://redis.io/commands/copy/
func (g pipelineGroupGeneric) Copy(ctx context.Context, source, destination string, option ...gredis.CopyOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "Copy", mustMergeOptionToArgs(
		[]any{source, destination}, usedOption,
	)...)
}

// Exists queues an Exists command and returns its future Cmd.
// https://redis.io/commands/exists/
func (g pipelineGroupGeneric) Exists(ctx context.Context, keys ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Exists", gconv.Interfaces(keys)...)
}

// Type queues a Type command and returns its future Cmd.
// https://redis.io/commands/type/
func (g pipelineGroupGeneric) Type(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Type", key)
}

// Unlink queues an Unlink command and returns its future Cmd.
// https://redis.io/commands/unlink/
func (g pipelineGroupGeneric) Unlink(ctx context.Context, keys ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Unlink", gconv.Interfaces(keys)...)
}

// Rename queues a Rename command and returns its future Cmd.
// https://redis.io/commands/rename/
func (g pipelineGroupGeneric) Rename(ctx context.Context, key, newKey string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Rename", key, newKey)
}

// RenameNX queues a RenameNX command and returns its future Cmd.
// https://redis.io/commands/renamenx/
func (g pipelineGroupGeneric) RenameNX(ctx context.Context, key, newKey string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "RenameNX", key, newKey)
}

// Move queues a Move command and returns its future Cmd.
// https://redis.io/commands/move/
func (g pipelineGroupGeneric) Move(ctx context.Context, key string, db int) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Move", key, db)
}

// Del queues a Del command and returns its future Cmd.
// https://redis.io/commands/del/
func (g pipelineGroupGeneric) Del(ctx context.Context, keys ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Del", gconv.Interfaces(keys)...)
}

// RandomKey queues a RandomKey command and returns its future Cmd.
// https://redis.io/commands/randomkey/
func (g pipelineGroupGeneric) RandomKey(ctx context.Context) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "RandomKey")
}

// DBSize queues a DBSize command and returns its future Cmd.
// https://redis.io/commands/dbsize/
func (g pipelineGroupGeneric) DBSize(ctx context.Context) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "DBSize")
}

// Keys queues a Keys command and returns its future Cmd.
// https://redis.io/commands/keys/
func (g pipelineGroupGeneric) Keys(ctx context.Context, pattern string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Keys", pattern)
}

// Scan queues a Scan command and returns its future Cmd.
// https://redis.io/commands/scan/
func (g pipelineGroupGeneric) Scan(ctx context.Context, cursor uint64, option ...gredis.ScanOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0].ToUsedOption()
	}
	return g.pipeliner.Do(ctx, "Scan", mustMergeOptionToArgs(
		[]any{cursor}, usedOption,
	)...)
}

// FlushDB queues a FlushDB command and returns its future Cmd.
// https://redis.io/commands/flushdb/
func (g pipelineGroupGeneric) FlushDB(ctx context.Context, option ...gredis.FlushOp) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "FlushDB", gconv.Interfaces(option)...)
}

// FlushAll queues a FlushAll command and returns its future Cmd.
// https://redis.io/commands/flushall/
func (g pipelineGroupGeneric) FlushAll(ctx context.Context, option ...gredis.FlushOp) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "FlushAll", gconv.Interfaces(option)...)
}

// Expire queues an Expire command and returns its future Cmd.
// https://redis.io/commands/expire/
func (g pipelineGroupGeneric) Expire(ctx context.Context, key string, seconds int64, option ...gredis.ExpireOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "Expire", mustMergeOptionToArgs(
		[]any{key, seconds}, usedOption,
	)...)
}

// ExpireAt queues an ExpireAt command and returns its future Cmd.
// https://redis.io/commands/expireat/
func (g pipelineGroupGeneric) ExpireAt(ctx context.Context, key string, time time.Time, option ...gredis.ExpireOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "ExpireAt", mustMergeOptionToArgs(
		[]any{key, gtime.New(time).Timestamp()}, usedOption,
	)...)
}

// ExpireTime queues an ExpireTime command and returns its future Cmd.
// https://redis.io/commands/expiretime/
func (g pipelineGroupGeneric) ExpireTime(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ExpireTime", key)
}

// TTL queues a TTL command and returns its future Cmd.
// https://redis.io/commands/ttl/
func (g pipelineGroupGeneric) TTL(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "TTL", key)
}

// Persist queues a Persist command and returns its future Cmd.
// https://redis.io/commands/persist/
func (g pipelineGroupGeneric) Persist(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Persist", key)
}

// PExpire queues a PExpire command and returns its future Cmd.
// https://redis.io/commands/pexpire/
func (g pipelineGroupGeneric) PExpire(ctx context.Context, key string, milliseconds int64, option ...gredis.ExpireOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "PExpire", mustMergeOptionToArgs(
		[]any{key, milliseconds}, usedOption,
	)...)
}

// PExpireAt queues a PExpireAt command and returns its future Cmd.
// https://redis.io/commands/pexpireat/
func (g pipelineGroupGeneric) PExpireAt(ctx context.Context, key string, time time.Time, option ...gredis.ExpireOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "PExpireAt", mustMergeOptionToArgs(
		[]any{key, gtime.New(time).TimestampMilli()}, usedOption,
	)...)
}

// PExpireTime queues a PExpireTime command and returns its future Cmd.
// https://redis.io/commands/pexpiretime/
func (g pipelineGroupGeneric) PExpireTime(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "PExpireTime", key)
}

// PTTL queues a PTTL command and returns its future Cmd.
// https://redis.io/commands/pttl/
func (g pipelineGroupGeneric) PTTL(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "PTTL", key)
}

// pipelineGroupHash implements gredis.IPipelineGroupHash.
type pipelineGroupHash struct {
	pipeliner *redisPipeliner
}

// HSet queues an HSet command and returns its future Cmd.
// https://redis.io/commands/hset/
func (g pipelineGroupHash) HSet(ctx context.Context, key string, fields map[string]any) *gredis.Cmd {
	s := []any{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	return g.pipeliner.Do(ctx, "HSet", s...)
}

// HSetNX queues an HSetNX command and returns its future Cmd.
// https://redis.io/commands/hsetnx/
func (g pipelineGroupHash) HSetNX(ctx context.Context, key, field string, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HSetNX", key, field, value)
}

// HGet queues an HGet command and returns its future Cmd.
// https://redis.io/commands/hget/
func (g pipelineGroupHash) HGet(ctx context.Context, key, field string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HGet", key, field)
}

// HStrLen queues an HStrLen command and returns its future Cmd.
// https://redis.io/commands/hstrlen/
func (g pipelineGroupHash) HStrLen(ctx context.Context, key, field string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HSTRLEN", key, field)
}

// HExists queues an HExists command and returns its future Cmd.
// https://redis.io/commands/hexists/
func (g pipelineGroupHash) HExists(ctx context.Context, key, field string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HExists", key, field)
}

// HDel queues an HDel command and returns its future Cmd.
// https://redis.io/commands/hdel/
func (g pipelineGroupHash) HDel(ctx context.Context, key string, fields ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HDel", append([]any{key}, gconv.Interfaces(fields)...)...)
}

// HLen queues an HLen command and returns its future Cmd.
// https://redis.io/commands/hlen/
func (g pipelineGroupHash) HLen(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HLen", key)
}

// HIncrBy queues an HIncrBy command and returns its future Cmd.
// https://redis.io/commands/hincrby/
func (g pipelineGroupHash) HIncrBy(ctx context.Context, key, field string, increment int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HIncrBy", key, field, increment)
}

// HIncrByFloat queues an HIncrByFloat command and returns its future Cmd.
// https://redis.io/commands/hincrbyfloat/
func (g pipelineGroupHash) HIncrByFloat(ctx context.Context, key, field string, increment float64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HIncrByFloat", key, field, increment)
}

// HMSet queues an HMSet command and returns its future Cmd.
// https://redis.io/commands/hmset/
func (g pipelineGroupHash) HMSet(ctx context.Context, key string, fields map[string]any) *gredis.Cmd {
	s := []any{key}
	for k, v := range fields {
		s = append(s, k, v)
	}
	return g.pipeliner.Do(ctx, "HMSet", s...)
}

// HMGet queues an HMGet command and returns its future Cmd.
// https://redis.io/commands/hmget/
func (g pipelineGroupHash) HMGet(ctx context.Context, key string, fields ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HMGet", append([]any{key}, gconv.Interfaces(fields)...)...)
}

// HKeys queues an HKeys command and returns its future Cmd.
// https://redis.io/commands/hkeys/
func (g pipelineGroupHash) HKeys(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HKeys", key)
}

// HVals queues an HVals command and returns its future Cmd.
// https://redis.io/commands/hvals/
func (g pipelineGroupHash) HVals(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HVals", key)
}

// HGetAll queues an HGetAll command and returns its future Cmd.
// https://redis.io/commands/hgetall/
func (g pipelineGroupHash) HGetAll(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "HGetAll", key)
}

// pipelineGroupString implements gredis.IPipelineGroupString.
type pipelineGroupString struct {
	pipeliner *redisPipeliner
}

// Set queues a Set command and returns its future Cmd.
// https://redis.io/commands/set/
func (g pipelineGroupString) Set(ctx context.Context, key string, value any, option ...gredis.SetOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "Set", mustMergeOptionToArgs(
		[]any{key, value}, usedOption,
	)...)
}

// SetNX queues a SetNX command and returns its future Cmd.
// https://redis.io/commands/setnx/
func (g pipelineGroupString) SetNX(ctx context.Context, key string, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SetNX", key, value)
}

// SetEX queues a SetEX command and returns its future Cmd.
// https://redis.io/commands/setex/
func (g pipelineGroupString) SetEX(ctx context.Context, key string, value any, ttlInSeconds int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SetEX", key, ttlInSeconds, value)
}

// Get queues a Get command and returns its future Cmd.
// https://redis.io/commands/get/
func (g pipelineGroupString) Get(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Get", key)
}

// GetDel queues a GetDel command and returns its future Cmd.
// https://redis.io/commands/getdel/
func (g pipelineGroupString) GetDel(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "GetDel", key)
}

// GetEX queues a GetEX command and returns its future Cmd.
// https://redis.io/commands/getex/
func (g pipelineGroupString) GetEX(ctx context.Context, key string, option ...gredis.GetEXOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "GetEX", mustMergeOptionToArgs(
		[]any{key}, usedOption,
	)...)
}

// GetSet queues a GetSet command and returns its future Cmd.
// https://redis.io/commands/getset/
func (g pipelineGroupString) GetSet(ctx context.Context, key string, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "GetSet", key, value)
}

// StrLen queues a StrLen command and returns its future Cmd.
// https://redis.io/commands/strlen/
func (g pipelineGroupString) StrLen(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "StrLen", key)
}

// Append queues an Append command and returns its future Cmd.
// https://redis.io/commands/append/
func (g pipelineGroupString) Append(ctx context.Context, key string, value string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Append", key, value)
}

// SetRange queues a SetRange command and returns its future Cmd.
// https://redis.io/commands/setrange/
func (g pipelineGroupString) SetRange(ctx context.Context, key string, offset int64, value string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SetRange", key, offset, value)
}

// GetRange queues a GetRange command and returns its future Cmd.
// https://redis.io/commands/getrange/
func (g pipelineGroupString) GetRange(ctx context.Context, key string, start, end int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "GetRange", key, start, end)
}

// Incr queues an Incr command and returns its future Cmd.
// https://redis.io/commands/incr/
func (g pipelineGroupString) Incr(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Incr", key)
}

// IncrBy queues an IncrBy command and returns its future Cmd.
// https://redis.io/commands/incrby/
func (g pipelineGroupString) IncrBy(ctx context.Context, key string, increment int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "IncrBy", key, increment)
}

// IncrByFloat queues an IncrByFloat command and returns its future Cmd.
// https://redis.io/commands/incrbyfloat/
func (g pipelineGroupString) IncrByFloat(ctx context.Context, key string, increment float64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "IncrByFloat", key, increment)
}

// Decr queues a Decr command and returns its future Cmd.
// https://redis.io/commands/decr/
func (g pipelineGroupString) Decr(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "Decr", key)
}

// DecrBy queues a DecrBy command and returns its future Cmd.
// https://redis.io/commands/decrby/
func (g pipelineGroupString) DecrBy(ctx context.Context, key string, decrement int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "DecrBy", key, decrement)
}

// MSet queues an MSet command and returns its future Cmd.
// https://redis.io/commands/mset/
func (g pipelineGroupString) MSet(ctx context.Context, keyValueMap map[string]any) *gredis.Cmd {
	var args []any
	for k, v := range keyValueMap {
		args = append(args, k, v)
	}
	return g.pipeliner.Do(ctx, "MSet", args...)
}

// MSetNX queues an MSetNX command and returns its future Cmd.
// https://redis.io/commands/msetnx/
func (g pipelineGroupString) MSetNX(ctx context.Context, keyValueMap map[string]any) *gredis.Cmd {
	var args []any
	for k, v := range keyValueMap {
		args = append(args, k, v)
	}
	return g.pipeliner.Do(ctx, "MSetNX", args...)
}

// MGet queues an MGet command and returns its future Cmd.
// https://redis.io/commands/mget/
func (g pipelineGroupString) MGet(ctx context.Context, keys ...string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "MGet", gconv.Interfaces(keys)...)
}

// pipelineGroupList implements gredis.IPipelineGroupList.
type pipelineGroupList struct {
	pipeliner *redisPipeliner
}

// LPush queues an LPush command and returns its future Cmd.
// https://redis.io/commands/lpush/
func (g pipelineGroupList) LPush(ctx context.Context, key string, values ...any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LPush", append([]any{key}, values...)...)
}

// LPushX queues an LPushX command and returns its future Cmd.
// https://redis.io/commands/lpushx/
func (g pipelineGroupList) LPushX(ctx context.Context, key string, element any, elements ...any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LPushX", append([]any{key, element}, elements...)...)
}

// RPush queues an RPush command and returns its future Cmd.
// https://redis.io/commands/rpush/
func (g pipelineGroupList) RPush(ctx context.Context, key string, values ...any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "RPush", append([]any{key}, values...)...)
}

// RPushX queues an RPushX command and returns its future Cmd.
// https://redis.io/commands/rpushx/
func (g pipelineGroupList) RPushX(ctx context.Context, key string, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "RPushX", key, value)
}

// LPop queues an LPop command and returns its future Cmd.
// https://redis.io/commands/lpop/
func (g pipelineGroupList) LPop(ctx context.Context, key string, count ...int) *gredis.Cmd {
	if len(count) > 0 {
		return g.pipeliner.Do(ctx, "LPop", key, count[0])
	}
	return g.pipeliner.Do(ctx, "LPop", key)
}

// RPop queues an RPop command and returns its future Cmd.
// https://redis.io/commands/rpop/
func (g pipelineGroupList) RPop(ctx context.Context, key string, count ...int) *gredis.Cmd {
	if len(count) > 0 {
		return g.pipeliner.Do(ctx, "RPop", key, count[0])
	}
	return g.pipeliner.Do(ctx, "RPop", key)
}

// LRem queues an LRem command and returns its future Cmd.
// https://redis.io/commands/lrem/
func (g pipelineGroupList) LRem(ctx context.Context, key string, count int64, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LRem", key, count, value)
}

// LLen queues an LLen command and returns its future Cmd.
// https://redis.io/commands/llen/
func (g pipelineGroupList) LLen(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LLen", key)
}

// LIndex queues an LIndex command and returns its future Cmd.
// https://redis.io/commands/lindex/
func (g pipelineGroupList) LIndex(ctx context.Context, key string, index int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LIndex", key, index)
}

// LInsert queues an LInsert command and returns its future Cmd.
// https://redis.io/commands/linsert/
func (g pipelineGroupList) LInsert(ctx context.Context, key string, op gredis.LInsertOp, pivot, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LInsert", key, string(op), pivot, value)
}

// LSet queues an LSet command and returns its future Cmd.
// https://redis.io/commands/lset/
func (g pipelineGroupList) LSet(ctx context.Context, key string, index int64, value any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LSet", key, index, value)
}

// LRange queues an LRange command and returns its future Cmd.
// https://redis.io/commands/lrange/
func (g pipelineGroupList) LRange(ctx context.Context, key string, start, stop int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LRange", key, start, stop)
}

// LTrim queues an LTrim command and returns its future Cmd.
// https://redis.io/commands/ltrim/
func (g pipelineGroupList) LTrim(ctx context.Context, key string, start, stop int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "LTrim", key, start, stop)
}

// RPopLPush queues an RPopLPush command and returns its future Cmd.
// https://redis.io/commands/rpoplpush/
func (g pipelineGroupList) RPopLPush(ctx context.Context, source, destination string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "RPopLPush", source, destination)
}

// pipelineGroupSet implements gredis.IPipelineGroupSet.
type pipelineGroupSet struct {
	pipeliner *redisPipeliner
}

// SAdd queues an SAdd command and returns its future Cmd.
// https://redis.io/commands/sadd/
func (g pipelineGroupSet) SAdd(ctx context.Context, key string, member any, members ...any) *gredis.Cmd {
	s := []any{key}
	s = append(s, member)
	s = append(s, members...)
	return g.pipeliner.Do(ctx, "SAdd", s...)
}

// SIsMember queues an SIsMember command and returns its future Cmd.
// https://redis.io/commands/sismember/
func (g pipelineGroupSet) SIsMember(ctx context.Context, key string, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SIsMember", key, member)
}

// SPop queues an SPop command and returns its future Cmd.
// https://redis.io/commands/spop/
func (g pipelineGroupSet) SPop(ctx context.Context, key string, count ...int) *gredis.Cmd {
	s := []any{key}
	s = append(s, gconv.Interfaces(count)...)
	return g.pipeliner.Do(ctx, "SPop", s...)
}

// SRandMember queues an SRandMember command and returns its future Cmd.
// https://redis.io/commands/srandmember/
func (g pipelineGroupSet) SRandMember(ctx context.Context, key string, count ...int) *gredis.Cmd {
	s := []any{key}
	s = append(s, gconv.Interfaces(count)...)
	return g.pipeliner.Do(ctx, "SRandMember", s...)
}

// SRem queues an SRem command and returns its future Cmd.
// https://redis.io/commands/srem/
func (g pipelineGroupSet) SRem(ctx context.Context, key string, member any, members ...any) *gredis.Cmd {
	s := []any{key}
	s = append(s, member)
	s = append(s, members...)
	return g.pipeliner.Do(ctx, "SRem", s...)
}

// SMove queues an SMove command and returns its future Cmd.
// https://redis.io/commands/smove/
func (g pipelineGroupSet) SMove(ctx context.Context, source, destination string, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SMove", source, destination, member)
}

// SCard queues an SCard command and returns its future Cmd.
// https://redis.io/commands/scard/
func (g pipelineGroupSet) SCard(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SCard", key)
}

// SMembers queues an SMembers command and returns its future Cmd.
// https://redis.io/commands/smembers/
func (g pipelineGroupSet) SMembers(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "SMembers", key)
}

// SMIsMember queues an SMIsMember command and returns its future Cmd.
// https://redis.io/commands/smismember/
func (g pipelineGroupSet) SMIsMember(ctx context.Context, key, member any, members ...any) *gredis.Cmd {
	s := []any{key, member}
	s = append(s, members...)
	return g.pipeliner.Do(ctx, "SMIsMember", s...)
}

// SInter queues an SInter command and returns its future Cmd.
// https://redis.io/commands/sinter/
func (g pipelineGroupSet) SInter(ctx context.Context, key string, keys ...string) *gredis.Cmd {
	s := []any{key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SInter", s...)
}

// SInterStore queues an SInterStore command and returns its future Cmd.
// https://redis.io/commands/sinterstore/
func (g pipelineGroupSet) SInterStore(ctx context.Context, destination string, key string, keys ...string) *gredis.Cmd {
	s := []any{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SInterStore", s...)
}

// SUnion queues an SUnion command and returns its future Cmd.
// https://redis.io/commands/sunion/
func (g pipelineGroupSet) SUnion(ctx context.Context, key string, keys ...string) *gredis.Cmd {
	s := []any{key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SUnion", s...)
}

// SUnionStore queues an SUnionStore command and returns its future Cmd.
// https://redis.io/commands/sunionstore/
func (g pipelineGroupSet) SUnionStore(ctx context.Context, destination, key string, keys ...string) *gredis.Cmd {
	s := []any{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SUnionStore", s...)
}

// SDiff queues an SDiff command and returns its future Cmd.
// https://redis.io/commands/sdiff/
func (g pipelineGroupSet) SDiff(ctx context.Context, key string, keys ...string) *gredis.Cmd {
	s := []any{key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SDiff", s...)
}

// SDiffStore queues an SDiffStore command and returns its future Cmd.
// https://redis.io/commands/sdiffstore/
func (g pipelineGroupSet) SDiffStore(ctx context.Context, destination string, key string, keys ...string) *gredis.Cmd {
	s := []any{destination, key}
	s = append(s, gconv.Interfaces(keys)...)
	return g.pipeliner.Do(ctx, "SDiffStore", s...)
}

// pipelineGroupSortedSet implements gredis.IPipelineGroupSortedSet.
type pipelineGroupSortedSet struct {
	pipeliner *redisPipeliner
}

// ZAdd queues a ZAdd command and returns its future Cmd.
// https://redis.io/commands/zadd/
func (g pipelineGroupSortedSet) ZAdd(
	ctx context.Context, key string, option *gredis.ZAddOption, member gredis.ZAddMember, members ...gredis.ZAddMember,
) *gredis.Cmd {
	s := mustMergeOptionToArgs(
		[]any{key}, option,
	)
	s = append(s, member.Score, member.Member)
	for _, item := range members {
		s = append(s, item.Score, item.Member)
	}
	return g.pipeliner.Do(ctx, "ZAdd", s...)
}

// ZScore queues a ZScore command and returns its future Cmd.
// https://redis.io/commands/zscore/
func (g pipelineGroupSortedSet) ZScore(ctx context.Context, key string, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZScore", key, member)
}

// ZIncrBy queues a ZIncrBy command and returns its future Cmd.
// https://redis.io/commands/zincrby/
func (g pipelineGroupSortedSet) ZIncrBy(ctx context.Context, key string, increment float64, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZIncrBy", key, increment, member)
}

// ZCard queues a ZCard command and returns its future Cmd.
// https://redis.io/commands/zcard/
func (g pipelineGroupSortedSet) ZCard(ctx context.Context, key string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZCard", key)
}

// ZCount queues a ZCount command and returns its future Cmd.
// https://redis.io/commands/zcount/
func (g pipelineGroupSortedSet) ZCount(ctx context.Context, key string, min, max string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZCount", key, min, max)
}

// ZRange queues a ZRange command and returns its future Cmd.
// https://redis.io/commands/zrange/
func (g pipelineGroupSortedSet) ZRange(ctx context.Context, key string, start, stop int64, option ...gredis.ZRangeOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "ZRange", mustMergeOptionToArgs(
		[]any{key, start, stop}, usedOption,
	)...)
}

// ZRevRange queues a ZRevRange command and returns its future Cmd.
// https://redis.io/commands/zrevrange/
func (g pipelineGroupSortedSet) ZRevRange(ctx context.Context, key string, start, stop int64, option ...gredis.ZRevRangeOption) *gredis.Cmd {
	var usedOption any
	if len(option) > 0 {
		usedOption = option[0]
	}
	return g.pipeliner.Do(ctx, "ZRevRange", mustMergeOptionToArgs(
		[]any{key, start, stop}, usedOption,
	)...)
}

// ZRank queues a ZRank command and returns its future Cmd.
// https://redis.io/commands/zrank/
func (g pipelineGroupSortedSet) ZRank(ctx context.Context, key string, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZRank", key, member)
}

// ZRevRank queues a ZRevRank command and returns its future Cmd.
// https://redis.io/commands/zrevrank/
func (g pipelineGroupSortedSet) ZRevRank(ctx context.Context, key string, member any) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZRevRank", key, member)
}

// ZRem queues a ZRem command and returns its future Cmd.
// https://redis.io/commands/zrem/
func (g pipelineGroupSortedSet) ZRem(ctx context.Context, key string, member any, members ...any) *gredis.Cmd {
	s := []any{key}
	s = append(s, member)
	s = append(s, members...)
	return g.pipeliner.Do(ctx, "ZRem", s...)
}

// ZRemRangeByRank queues a ZRemRangeByRank command and returns its future Cmd.
// https://redis.io/commands/zremrangebyrank/
func (g pipelineGroupSortedSet) ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZRemRangeByRank", key, start, stop)
}

// ZRemRangeByScore queues a ZRemRangeByScore command and returns its future Cmd.
// https://redis.io/commands/zremrangebyscore/
func (g pipelineGroupSortedSet) ZRemRangeByScore(ctx context.Context, key string, min, max string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZRemRangeByScore", key, min, max)
}

// ZRemRangeByLex queues a ZRemRangeByLex command and returns its future Cmd.
// https://redis.io/commands/zremrangebylex/
func (g pipelineGroupSortedSet) ZRemRangeByLex(ctx context.Context, key string, min, max string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZRemRangeByLex", key, min, max)
}

// ZLexCount queues a ZLexCount command and returns its future Cmd.
// https://redis.io/commands/zlexcount/
func (g pipelineGroupSortedSet) ZLexCount(ctx context.Context, key, min, max string) *gredis.Cmd {
	return g.pipeliner.Do(ctx, "ZLexCount", key, min, max)
}
