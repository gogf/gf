// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (r *Redis) Do(ctx context.Context, command string, args ...any) (*gvar.Var, error) {
	conn, err := r.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close(ctx)
	}()
	return conn.Do(ctx, command, args...)
}

// Close closes the redis connection pool, which will release all connections reserved by this pool.
// It is commonly not necessary to call Close manually.
func (r *Redis) Close(ctx context.Context) (err error) {
	if err = r.client.Close(); err != nil {
		err = gerror.Wrap(err, `Operation Client Close failed`)
	}
	return
}

// Conn retrieves and returns a connection object for continuous operations.
// Note that you should call Close function manually if you do not use this connection any further.
func (r *Redis) Conn(ctx context.Context) (gredis.Conn, error) {
	return &Conn{
		redis: r,
	}, nil
}

// Client returns the underlying redis client instance.
// This method provides access to the raw redis client for advanced operations
// that are not covered by the standard Redis interface.
//
// Example usage with type assertion:
//
//	import goredis "github.com/redis/go-redis/v9"
//
//	func ExampleUsage(ctx context.Context, redis *Redis) error {
//		client := redis.Client()
//		universalClient, ok := client.(goredis.UniversalClient)
//		if !ok {
//			return errors.New("failed to assert to UniversalClient")
//		}
//
//		// Use universalClient for advanced operations like Pipeline
//		pipe := universalClient.Pipeline()
//		pipe.Set(ctx, "key1", "value1", 0)
//		pipe.Set(ctx, "key2", "value2", 0)
//		results, err := pipe.Exec(ctx)
//		if err != nil {
//			return err
//		}
//		// ... handle results
//		return nil
//	}
func (r *Redis) Client() gredis.RedisRawClient {
	return r.client
}
