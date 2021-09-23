// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/guid"
	"github.com/gogf/gf/util/gutil"
	"testing"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/test/gtest"
)

var (
	config = &gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}
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

func Test_Do(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		_, err = redis.Do(ctx, "SET", "k", "v")
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(r, []byte("v"))

		_, err = redis.Do(ctx, "DEL", "k")
		t.Assert(err, nil)
		r, err = redis.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(r, nil)
	})
}

func Test_Stats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		array := make([]*gredis.RedisConn, 0)
		for i := 0; i < 10; i++ {
			conn, err := redis.Conn(ctx)
			t.AssertNil(err)
			array = append(array, conn)
		}
		stats, err := redis.Stats(ctx)
		t.AssertNil(err)
		t.Assert(stats.ActiveCount(), 10)
		t.Assert(stats.IdleCount(), 0)

		for i := 0; i < 10; i++ {
			t.AssertNil(array[i].Close(ctx))
		}

		stats, err = redis.Stats(ctx)
		t.AssertNil(err)
		t.Assert(stats.ActiveCount(), 10)
		t.Assert(stats.IdleCount(), 10)
		//time.Sleep(3000*time.Millisecond)
		//stats  = redis.Stats()
		//fmt.Println(stats)
		//t.Assert(stats.ActiveCount,  0)
		//t.Assert(stats.IdleCount,    0)
	})
}

func Test_Conn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)

		key := gconv.String(gtime.TimestampNano())
		value := []byte("v")
		r, err := conn.Do(ctx, "SET", key, value)
		t.Assert(err, nil)

		r, err = conn.Do(ctx, "GET", key)
		t.Assert(err, nil)
		t.Assert(r, value)

		_, err = conn.Do(ctx, "DEL", key)
		t.Assert(err, nil)
		r, err = conn.Do(ctx, "GET", key)
		t.Assert(err, nil)
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
		t.Assert(err, nil)

		r, err := conn.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(r, []byte("v"))

		_, err = conn.Do(ctx, "DEL", "k")
		t.Assert(err, nil)
		r, err = conn.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(r, nil)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config1 := &gredis.Config{
			Host:           "192.111.0.2",
			Port:           6379,
			Db:             1,
			ConnectTimeout: time.Second,
		}
		redis, err := gredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		_, err = redis.Do(ctx, "info")
		t.AssertNE(err, nil)

		config1 = &gredis.Config{
			Host: "127.0.0.1",
			Port: 6379,
			Db:   100,
		}
		redis, err = gredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		_, err = redis.Do(ctx, "info")
		t.AssertNE(err, nil)

		redis = gredis.Instance("gf")
		t.Assert(redis == nil, true)
		gredis.ClearConfig()

		redis, err = gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		_, err = redis.Do(ctx, "SET", "k", "v")
		t.Assert(err, nil)

		v, err := redis.Do(ctx, "GET", "k")
		t.Assert(err, nil)
		t.Assert(v.String(), "v")

		conn, err := redis.Conn(ctx)
		t.AssertNil(err)
		defer conn.Close(ctx)
		_, err = conn.Do(ctx, "SET", "k", "v")
		t.Assert(err, nil)

		_, err = conn.Do(ctx, "Subscribe", "gf")
		t.Assert(err, nil)

		_, err = redis.Do(ctx, "PUBLISH", "gf", "test")
		t.Assert(err, nil)

		v, _ = conn.Receive(ctx)
		t.Assert(len(v.Strings()), 3)
		t.Assert(v.Strings()[2], "test")

		time.Sleep(time.Second)
	})
}

func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis, err := gredis.New(config)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		defer func() {
			redis.Do(ctx, "DEL", "key-true")
			redis.Do(ctx, "DEL", "key-false")
		}()

		_, err = redis.Do(ctx, "SET", "key-true", true)
		t.Assert(err, nil)

		_, err = redis.Do(ctx, "SET", "key-false", false)
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "GET", "key-true")
		t.Assert(err, nil)
		t.Assert(r.Bool(), true)

		r, err = redis.Do(ctx, "GET", "key-false")
		t.Assert(err, nil)
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
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "GET", key)
		t.Assert(err, nil)
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
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.Assert(err, nil)
		t.Assert(r.Strings(), g.ArrayStr{"name", "john"})
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
		t.Assert(err, nil)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "HGETALL", key)
		t.Assert(err, nil)
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
		t.Assert(err, nil)
		_, err = redis.Do(ctx, "HSET", key, "name", "john")
		t.Assert(err, nil)

		result, err := redis.Do(ctx, "HGETALL", key)
		t.Assert(err, nil)

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
		t.Assert(err, nil)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.Assert(err, nil)
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
		t.Assert(err, nil)
		v, err := redis.Do(ctx, "HMGET", key, "name")
		t.Assert(err, nil)
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
		t.Assert(err, nil)

		r, err := redis.Do(ctx, "GET", key)
		t.Assert(err, nil)
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
		t.Assert(err, nil)

		result, err = redis.Do(ctx, "GET", key)
		t.Assert(err, nil)

		var users2 []User
		err = result.Structs(&users2)
		t.Assert(err, nil)
		t.Assert(users2, users1)
	})
}
