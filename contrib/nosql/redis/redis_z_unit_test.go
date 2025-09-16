// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"
	"time"

	goredis "github.com/redis/go-redis/v9"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_NewClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)

		err = redis.Close(ctx)
		t.AssertNil(err)
	})
}

func Test_Client(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		// Test getting the client
		client := redis.Client()
		t.AssertNE(client, nil)

		// Test type assertion to goredis.UniversalClient
		universalClient, ok := client.(goredis.UniversalClient)
		t.Assert(ok, true)
		t.AssertNE(universalClient, nil)

		// Test that we can use the client directly for redis operations
		// This demonstrates that the returned client is properly configured
		result := universalClient.Set(ctx, "test-client-key", "test-value", 0)
		t.AssertNil(result.Err())

		getResult := universalClient.Get(ctx, "test-client-key")
		t.AssertNil(getResult.Err())
		t.Assert(getResult.Val(), "test-value")

		// Clean up
		delResult := universalClient.Del(ctx, "test-client-key")
		t.AssertNil(delResult.Err())

		// Test Pipeline functionality
		pipe := universalClient.Pipeline()
		t.AssertNE(pipe, nil)

		// Add multiple commands to the pipeline
		pipe.Set(ctx, "pipeline-key1", "value1", 0)
		pipe.Set(ctx, "pipeline-key2", "value2", 0)
		pipe.Set(ctx, "pipeline-key3", "value3", 0)
		pipe.Get(ctx, "pipeline-key1")
		pipe.Get(ctx, "pipeline-key2")
		pipe.Get(ctx, "pipeline-key3")

		// Execute the pipeline
		results, err := pipe.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 6) // 3 SET commands + 3 GET commands

		// Verify the SET results
		for i := range 3 {
			t.AssertNil(results[i].Err())
		}

		// Verify the GET results
		getResults := results[3:]
		t.Assert(getResults[0].(*goredis.StringCmd).Val(), "value1")
		t.Assert(getResults[1].(*goredis.StringCmd).Val(), "value2")
		t.Assert(getResults[2].(*goredis.StringCmd).Val(), "value3")

		// Clean up pipeline test keys
		cleanupPipe := universalClient.Pipeline()
		cleanupPipe.Del(ctx, "pipeline-key1")
		cleanupPipe.Del(ctx, "pipeline-key2")
		cleanupPipe.Del(ctx, "pipeline-key3")
		_, err = cleanupPipe.Exec(ctx)
		t.AssertNil(err)
	})
}

func Test_Do(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		_, err := redis.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = redis.Do(ctx, "DEL", "k")
		t.AssertNil(err)
		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Conn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		key := gconv.String(gtime.TimestampNano())
		value := []byte("v")
		r, err := conn.Do(ctx, "SET", key, value)
		t.AssertNil(err)

		r, err = conn.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r, value)

		_, err = conn.Do(ctx, "DEL", key)
		t.AssertNil(err)
		r, err = conn.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Instance(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		group := "my-test"
		gredis.SetConfig(config, group)
		defer gredis.RemoveConfig(group)

		redis := gredis.Instance(group)
		defer redis.Close(ctx)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err := conn.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = conn.Do(ctx, "DEL", "k")
		t.AssertNil(err)
		r, err = conn.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config1 := &gredis.Config{
			Address:     "192.111.0.2:6379",
			Db:          1,
			DialTimeout: time.Second,
		}
		r, err := gredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "info")
		t.AssertNE(err, nil)

		config1 = &gredis.Config{
			Address: "127.0.0.1:6379",
			Db:      100,
		}
		r, err = gredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "info")
		t.AssertNE(err, nil)

		r = gredis.Instance("gf")
		t.Assert(r == nil, true)
		gredis.ClearConfig()

		r, err = gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(r, nil)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		v, err := r.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(v.String(), "v")

		conn, err := r.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)
	})
}

func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() {
			redis.Do(ctx, "DEL", "key-true")
			redis.Do(ctx, "DEL", "key-false")
		}()

		_, err := redis.Do(ctx, "SET", "key-true", true)
		t.AssertNil(err)

		_, err = redis.Do(ctx, "SET", "key-false", false)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", "key-true")
		t.AssertNil(err)
		t.Assert(r.Bool(), true)

		r, err = redis.Do(ctx, "GET", "key-false")
		t.AssertNil(err)
		t.Assert(r.Bool(), false)
	})
}

func Test_Int(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		key := guid.S()
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "SET", key, 1)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r.Int(), 1)
	})
}

func Test_HSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		key := guid.S()
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)
		t.Assert(r.MapStrStr(), g.MapStrStr{"name": "john"})
	})
}

func Test_HGetAll1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			key = guid.S()
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "id", 100)
		t.AssertNil(err)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)
		t.Assert(r.Map(), g.MapStrAny{
			"id":   100,
			"name": "john",
		})
	})
}

func Test_HGetAll2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			key = guid.S()
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HSET", key, "id", 100)
		t.AssertNil(err)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.AssertNil(err)

		result, err := redis.Do(ctx, "HGETALL", key)
		t.AssertNil(err)

		t.Assert(gconv.Uint(result.MapStrVar()["id"]), 100)
		t.Assert(result.MapStrVar()["id"].Uint(), 100)
	})
}

func Test_HMSet(t *testing.T) {
	// map
	gtest.C(t, func(t *gtest.T) {
		var (
			key  = guid.S()
			data = g.Map{
				"name":  "gf",
				"sex":   0,
				"score": 100,
			}
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HMSET", append(g.Slice{key}, gutil.MapToSlice(data)...)...)
		t.AssertNil(err)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.AssertNil(err)
		t.Assert(v.Slice(), g.Slice{data["name"]})
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name  string `json:"name"`
			Sex   int    `json:"sex"`
			Score int    `json:"score"`
		}
		var (
			key  = guid.S()
			data = &User{
				Name:  "gf",
				Sex:   0,
				Score: 100,
			}
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)
		defer redis.Do(ctx, "DEL", key)

		_, err = redis.Do(ctx, "HMSET", append(g.Slice{key}, gutil.StructToSlice(data)...)...)
		t.AssertNil(err)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.AssertNil(err)
		t.Assert(v.Slice(), g.Slice{data.Name})
	})
}

func Test_Auto_Marshal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			key = guid.S()
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		defer redis.Do(ctx, "DEL", key)

		type User struct {
			Id   int
			Name string
		}

		user := &User{
			Id:   10000,
			Name: "john",
		}

		_, err = redis.Do(ctx, "SET", key, user)
		t.AssertNil(err)

		r, err := redis.Do(ctx, "GET", key)
		t.AssertNil(err)
		t.Assert(r.Map(), g.MapStrAny{
			"Id":   user.Id,
			"Name": user.Name,
		})

		var user2 *User
		t.Assert(r.Struct(&user2), nil)
		t.Assert(user2.Id, user.Id)
		t.Assert(user2.Name, user.Name)
	})
}

func Test_Auto_MarshalSlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			key = "user-slice"
		)
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Do(ctx, "DEL", key)
		type User struct {
			Id   int
			Name string
		}
		var (
			result *gvar.Var
			users1 = []User{
				{
					Id:   1,
					Name: "john1",
				},
				{
					Id:   2,
					Name: "john2",
				},
			}
		)

		_, err = redis.Do(ctx, "SET", key, users1)
		t.AssertNil(err)

		result, err = redis.Do(ctx, "GET", key)
		t.AssertNil(err)

		var users2 []User
		err = result.Structs(&users2)
		t.AssertNil(err)
		t.Assert(users2, users1)
	})
}
