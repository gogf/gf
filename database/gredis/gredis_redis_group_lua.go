// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"
)

type RedisGroupLua struct {
	redis *Redis
}

func (r *Redis) Lua() *RedisGroupLua {
	return &RedisGroupLua{
		redis: r,
	}
}

// Eval Invoke the execution of a server-side Lua script.
//
//The first argument is the script's source code. Scripts are written in Lua and executed by the embedded Lua 5.1 interpreter in Redis.
//
//The second argument is the number of input key name arguments, followed by all the keys accessed by the script. These names of input keys are available to the script as the KEYS global runtime variable Any additional input arguments should not represent names of keys.
//
//Important: to ensure the correct execution of scripts, both in standalone and clustered deployments, all names of keys that a script accesses must be explicitly provided as input key arguments. The script should only access keys whose names are given as input arguments. Scripts should never access keys with programmatically-generated names or based on the contents of data structures stored in the database.
//
//Please refer to the Redis Programmability and Introduction to Eval Scripts for more information about Lua scripts.
//
// https://redis.io/commands/eval/
func (r *RedisGroupLua) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	v, err := r.redis.Do(ctx, "EVAL", script, keys, args)
	return v.Interface(), err
}

// EvalSha Evaluate a script from the server's cache by its SHA1 digest.
//
//The server caches scripts by using the SCRIPT LOAD command. The command is otherwise identical to EVAL.
//
//Please refer to the Redis Programmability and Introduction to Eval Scripts for more information about Lua scripts.
//
// https://redis.io/commands/evalsha/
func (r *RedisGroupLua) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	v, err := r.redis.Do(ctx, "EVALSHA", sha1, keys, args)
	return v.Interface(), err
}

// ScriptLoad  Load a script into the scripts cache, without executing it. After the specified command is loaded into the script cache it will be callable using EVALSHA with the correct SHA1 digest of the script, exactly like after the first successful invocation of EVAL.
//
//The script is guaranteed to stay in the script cache forever (unless SCRIPT FLUSH is called).
//
//The command works in the same way even if the script was already present in the script cache.
//
// https://redis.io/commands/script-load/
func (r *RedisGroupLua) ScriptLoad(ctx context.Context, script string) (string, error) {
	v, err := r.redis.Do(ctx, "SCRIPT LOAD", script)
	return v.String(), err
}

// ScriptExists Returns information about the existence of the scripts in the script cache.
//
//This command accepts one or more SHA1 digests and returns a list of ones or zeros to signal if the scripts are already defined or not inside the script cache. This can be useful before a pipelining operation to ensure that scripts are loaded (and if not, to load them using SCRIPT LOAD) so that the pipelining operation can be performed solely using EVALSHA instead of EVAL to save bandwidth.
//
// https://redis.io/commands/script-exists/
func (r *RedisGroupLua) ScriptExists(ctx context.Context, sha1s ...string) ([]interface{}, error) {
	v, err := r.redis.Do(ctx, "SCRIPT EXISTS", sha1s)
	return v.Interfaces(), err
}

// ScriptFlush Flush the Lua scripts cache.
//
//By default, SCRIPT FLUSH will synchronously flush the cache. Starting with Redis 6.2, setting the lazyfree-lazy-user-flush configuration directive to "yes" changes the default flush mode to asynchronous.
//
// https://redis.io/commands/script-flush/
func (r *RedisGroupLua) ScriptFlush(ctx context.Context, options string) (string, error) {
	v, err := r.redis.Do(ctx, "SCRIPT FLUSH", options)
	return v.String(), err
}

// ScriptKill Kills the currently executing EVAL script, assuming no write operation was yet performed by the script.
//
//This command is mainly useful to kill a script that is running for too much time(for instance, because it entered an infinite loop because of a bug). The script will be killed, and the client currently blocked into EVAL will see the command returning with an error.
//
//If the script has already performed write operations, it can not be killed in this way because it would violate Lua's script atomicity contract. In such a case, only SHUTDOWN NOSAVE can kill the script, killing the Redis process in a hard way and preventing it from persisting with half-written information.
//
// https://redis.io/commands/script-kill/
func (r *RedisGroupLua) ScriptKill(ctx context.Context) (string, error) {
	v, err := r.redis.Do(ctx, "SCRIPT KILL")
	return v.String(), err
}
