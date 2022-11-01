// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

// RedisGroupScript provides script functions for redis.
type RedisGroupScript struct {
	redis *Redis
}

// GroupScript creates and returns RedisGroupScript.
func (r *Redis) GroupScript() RedisGroupScript {
	return RedisGroupScript{
		redis: r,
	}
}

// Eval invokes the execution of a server-side Lua script.
//
// https://redis.io/commands/eval/
func (r RedisGroupScript) Eval(ctx context.Context, script string, numKeys int64, keys []string, args []interface{}) (*gvar.Var, error) {
	var s = []interface{}{script, numKeys}
	s = append(s, gconv.Interfaces(keys)...)
	s = append(s, args...)
	v, err := r.redis.Do(ctx, "Eval", s...)
	return v, err
}

// EvalSha evaluates a script from the server's cache by its SHA1 digest.
//
// The server caches scripts by using the SCRIPT LOAD command.
// The command is otherwise identical to EVAL.
//
// https://redis.io/commands/evalsha/
func (r RedisGroupScript) EvalSha(ctx context.Context, sha1 string, numKeys int64, keys []string, args []interface{}) (*gvar.Var, error) {
	var s = []interface{}{sha1, numKeys}
	s = append(s, gconv.Interfaces(keys)...)
	s = append(s, args...)
	v, err := r.redis.Do(ctx, "EvalSha", s...)
	return v, err
}

// ScriptLoad loads a script into the scripts cache, without executing it.
//
// It returns the SHA1 digest of the script added into the script cache.
//
// https://redis.io/commands/script-load/
func (r RedisGroupScript) ScriptLoad(ctx context.Context, script string) (string, error) {
	v, err := r.redis.Do(ctx, "Script", "Load", script)
	return v.String(), err
}

// ScriptExists returns information about the existence of the scripts in the script cache.
//
// It returns an array of integers that correspond to the specified SHA1 digest arguments.
// For every corresponding SHA1 digest of a script that actually exists in the script cache,
// a 1 is returned, otherwise 0 is returned.
//
// https://redis.io/commands/script-exists/
func (r RedisGroupScript) ScriptExists(ctx context.Context, sha1 string, sha1s ...string) (map[string]bool, error) {
	var (
		s         []interface{}
		sha1Array = append([]interface{}{sha1}, gconv.Interfaces(sha1s)...)
	)
	s = append(s, "Exists")
	s = append(s, sha1Array...)
	result, err := r.redis.Do(ctx, "Script", s...)
	var (
		m           = make(map[string]bool)
		resultArray = result.Vars()
	)
	for i := 0; i < len(sha1Array); i++ {
		m[gconv.String(sha1Array[i])] = resultArray[i].Bool()
	}
	return m, err
}

// ScriptFlushOption provides options for function ScriptFlush.
type ScriptFlushOption struct {
	SYNC  bool // SYNC  flushes the cache synchronously.
	ASYNC bool // ASYNC flushes the cache asynchronously.
}

// ScriptFlush flush the Lua scripts cache.
//
// https://redis.io/commands/script-flush/
func (r RedisGroupScript) ScriptFlush(ctx context.Context, option ...ScriptFlushOption) error {
	var usedOption interface{}
	if len(option) > 0 {
		usedOption = option[0]
	}
	var s []interface{}
	s = append(s, "Flush")
	s = append(s, mustMergeOptionToArgs(
		[]interface{}{}, usedOption,
	)...)
	_, err := r.redis.Do(ctx, "Script", s...)
	return err
}

// ScriptKill kills the currently executing EVAL script, assuming no write operation was yet performed
// by the script.
//
// https://redis.io/commands/script-kill/
func (r RedisGroupScript) ScriptKill(ctx context.Context) error {
	_, err := r.redis.Do(ctx, "Script", "Kill")
	return err
}
