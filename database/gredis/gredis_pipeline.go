// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT License was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file defines the Pipeliner, PipelinerOperation, PipelinerGroup interfaces,
// the six pipeline command group interfaces, and the Tx interface used for
// pipeline and transaction support in the Redis adapter.

package gredis

import (
	"context"
	"time"
)

// Pipeliner queues Redis commands for batch execution in a single network round-trip.
// Commands are buffered locally and sent to the server when Exec is called.
// A Pipeliner is obtained via Redis.Pipeline() or Redis.TxPipeline().
//
// Usage example:
//
//	pipe := redis.Pipeline(ctx)
//	cmd1 := pipe.PipelineGroupHash().HSet(ctx, "key", map[string]any{"field": "value"})
//	cmd2 := pipe.PipelineGroupString().Get(ctx, "key2")
//	pipe.Exec(ctx)
//	val1, _ := cmd1.Result()
//	val2, _ := cmd2.Result()
type Pipeliner interface {
	PipelinerOperation
	PipelinerGroup
}

// PipelinerOperation defines the core pipeline operations.
type PipelinerOperation interface {
	// Do queues a raw Redis command and returns its future Cmd.
	// The command is not sent to the server until Exec is called.
	Do(ctx context.Context, command string, args ...any) *Cmd

	// Exec sends all queued commands to the server in a single batch,
	// populates each queued Cmd's result, and returns all Cmds.
	// After Exec returns, all Cmd objects returned by queued commands are populated.
	Exec(ctx context.Context) ([]*Cmd, error)

	// Discard discards all queued commands without sending them to the server.
	Discard()
}

// PipelinerGroup provides typed command access for the Pipeliner,
// mirroring the existing Redis group interfaces but with *Cmd return types.
type PipelinerGroup interface {
	PipelineGroupGeneric() IPipelineGroupGeneric
	PipelineGroupHash() IPipelineGroupHash
	PipelineGroupString() IPipelineGroupString
	PipelineGroupList() IPipelineGroupList
	PipelineGroupSet() IPipelineGroupSet
	PipelineGroupSortedSet() IPipelineGroupSortedSet
}

// IPipelineGroupGeneric manages generic Redis operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
// Results are populated after Pipeliner.Exec() is called.
type IPipelineGroupGeneric interface {
	Copy(ctx context.Context, source, destination string, option ...CopyOption) *Cmd
	Exists(ctx context.Context, keys ...string) *Cmd
	Type(ctx context.Context, key string) *Cmd
	Unlink(ctx context.Context, keys ...string) *Cmd
	Rename(ctx context.Context, key, newKey string) *Cmd
	RenameNX(ctx context.Context, key, newKey string) *Cmd
	Move(ctx context.Context, key string, db int) *Cmd
	Del(ctx context.Context, keys ...string) *Cmd
	RandomKey(ctx context.Context) *Cmd
	DBSize(ctx context.Context) *Cmd
	Keys(ctx context.Context, pattern string) *Cmd
	Scan(ctx context.Context, cursor uint64, option ...ScanOption) *Cmd
	FlushDB(ctx context.Context, option ...FlushOp) *Cmd
	FlushAll(ctx context.Context, option ...FlushOp) *Cmd
	Expire(ctx context.Context, key string, seconds int64, option ...ExpireOption) *Cmd
	ExpireAt(ctx context.Context, key string, time time.Time, option ...ExpireOption) *Cmd
	ExpireTime(ctx context.Context, key string) *Cmd
	TTL(ctx context.Context, key string) *Cmd
	Persist(ctx context.Context, key string) *Cmd
	PExpire(ctx context.Context, key string, milliseconds int64, option ...ExpireOption) *Cmd
	PExpireAt(ctx context.Context, key string, time time.Time, option ...ExpireOption) *Cmd
	PExpireTime(ctx context.Context, key string) *Cmd
	PTTL(ctx context.Context, key string) *Cmd
}

// IPipelineGroupHash manages Redis hash operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
type IPipelineGroupHash interface {
	HSet(ctx context.Context, key string, fields map[string]any) *Cmd
	HSetNX(ctx context.Context, key, field string, value any) *Cmd
	HGet(ctx context.Context, key, field string) *Cmd
	HStrLen(ctx context.Context, key, field string) *Cmd
	HExists(ctx context.Context, key, field string) *Cmd
	HDel(ctx context.Context, key string, fields ...string) *Cmd
	HLen(ctx context.Context, key string) *Cmd
	HIncrBy(ctx context.Context, key, field string, increment int64) *Cmd
	HIncrByFloat(ctx context.Context, key, field string, increment float64) *Cmd
	HMSet(ctx context.Context, key string, fields map[string]any) *Cmd
	HMGet(ctx context.Context, key string, fields ...string) *Cmd
	HKeys(ctx context.Context, key string) *Cmd
	HVals(ctx context.Context, key string) *Cmd
	HGetAll(ctx context.Context, key string) *Cmd
}

// IPipelineGroupString manages Redis string operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
type IPipelineGroupString interface {
	Set(ctx context.Context, key string, value any, option ...SetOption) *Cmd
	SetNX(ctx context.Context, key string, value any) *Cmd
	SetEX(ctx context.Context, key string, value any, ttlInSeconds int64) *Cmd
	Get(ctx context.Context, key string) *Cmd
	GetDel(ctx context.Context, key string) *Cmd
	GetEX(ctx context.Context, key string, option ...GetEXOption) *Cmd
	GetSet(ctx context.Context, key string, value any) *Cmd
	StrLen(ctx context.Context, key string) *Cmd
	Append(ctx context.Context, key string, value string) *Cmd
	SetRange(ctx context.Context, key string, offset int64, value string) *Cmd
	GetRange(ctx context.Context, key string, start, end int64) *Cmd
	Incr(ctx context.Context, key string) *Cmd
	IncrBy(ctx context.Context, key string, increment int64) *Cmd
	IncrByFloat(ctx context.Context, key string, increment float64) *Cmd
	Decr(ctx context.Context, key string) *Cmd
	DecrBy(ctx context.Context, key string, decrement int64) *Cmd
	MSet(ctx context.Context, keyValueMap map[string]any) *Cmd
	MSetNX(ctx context.Context, keyValueMap map[string]any) *Cmd
	MGet(ctx context.Context, keys ...string) *Cmd
}

// IPipelineGroupList manages Redis list operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
type IPipelineGroupList interface {
	LPush(ctx context.Context, key string, values ...any) *Cmd
	LPushX(ctx context.Context, key string, element any, elements ...any) *Cmd
	RPush(ctx context.Context, key string, values ...any) *Cmd
	RPushX(ctx context.Context, key string, value any) *Cmd
	LPop(ctx context.Context, key string, count ...int) *Cmd
	RPop(ctx context.Context, key string, count ...int) *Cmd
	LRem(ctx context.Context, key string, count int64, value any) *Cmd
	LLen(ctx context.Context, key string) *Cmd
	LIndex(ctx context.Context, key string, index int64) *Cmd
	LInsert(ctx context.Context, key string, op LInsertOp, pivot, value any) *Cmd
	LSet(ctx context.Context, key string, index int64, value any) *Cmd
	LRange(ctx context.Context, key string, start, stop int64) *Cmd
	LTrim(ctx context.Context, key string, start, stop int64) *Cmd
	RPopLPush(ctx context.Context, source, destination string) *Cmd
}

// IPipelineGroupSet manages Redis set operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
type IPipelineGroupSet interface {
	SAdd(ctx context.Context, key string, member any, members ...any) *Cmd
	SIsMember(ctx context.Context, key string, member any) *Cmd
	SPop(ctx context.Context, key string, count ...int) *Cmd
	SRandMember(ctx context.Context, key string, count ...int) *Cmd
	SRem(ctx context.Context, key string, member any, members ...any) *Cmd
	SMove(ctx context.Context, source, destination string, member any) *Cmd
	SCard(ctx context.Context, key string) *Cmd
	SMembers(ctx context.Context, key string) *Cmd
	SMIsMember(ctx context.Context, key, member any, members ...any) *Cmd
	SInter(ctx context.Context, key string, keys ...string) *Cmd
	SInterStore(ctx context.Context, destination string, key string, keys ...string) *Cmd
	SUnion(ctx context.Context, key string, keys ...string) *Cmd
	SUnionStore(ctx context.Context, destination, key string, keys ...string) *Cmd
	SDiff(ctx context.Context, key string, keys ...string) *Cmd
	SDiffStore(ctx context.Context, destination string, key string, keys ...string) *Cmd
}

// IPipelineGroupSortedSet manages Redis sorted set operations in pipeline mode.
// Each method queues the command and returns a *Cmd future.
type IPipelineGroupSortedSet interface {
	ZAdd(ctx context.Context, key string, option *ZAddOption, member ZAddMember, members ...ZAddMember) *Cmd
	ZScore(ctx context.Context, key string, member any) *Cmd
	ZIncrBy(ctx context.Context, key string, increment float64, member any) *Cmd
	ZCard(ctx context.Context, key string) *Cmd
	ZCount(ctx context.Context, key string, min, max string) *Cmd
	ZRange(ctx context.Context, key string, start, stop int64, option ...ZRangeOption) *Cmd
	ZRevRange(ctx context.Context, key string, start, stop int64, option ...ZRevRangeOption) *Cmd
	ZRank(ctx context.Context, key string, member any) *Cmd
	ZRevRank(ctx context.Context, key string, member any) *Cmd
	ZRem(ctx context.Context, key string, member any, members ...any) *Cmd
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) *Cmd
	ZRemRangeByScore(ctx context.Context, key string, min, max string) *Cmd
	ZRemRangeByLex(ctx context.Context, key string, min, max string) *Cmd
	ZLexCount(ctx context.Context, key, min, max string) *Cmd
}

// Tx represents a Redis transaction context used with Watch.
// It embeds Pipeliner so commands can be queued within the transaction callback.
// Commands queued on Tx are executed atomically via MULTI/EXEC.
type Tx interface {
	Pipeliner
}
