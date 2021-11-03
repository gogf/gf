// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gutil"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	config = &gredis.Config{
		Address: `:6379`,
		Db:      1,
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
			Address:     "192.111.0.2:6379",
			Db:          1,
			DialTimeout: time.Second,
		}
		redis, err := gredis.New(config1)
		t.AssertNil(err)
		t.AssertNE(redis, nil)
		defer redis.Close(ctx)

		_, err = redis.Do(ctx, "info")
		t.AssertNE(err, nil)

		config1 = &gredis.Config{
			Address: "127.0.0.1:6379",
			Db:      100,
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
		t.AssertNil(err)

		_, err = conn.Do(ctx, "Subscribe", "gf")
		t.AssertNil(err)

		_, err = redis.Do(ctx, "PUBLISH", "gf", "test")
		t.AssertNil(err)

		v, err = conn.Receive(ctx)
		t.AssertNil(err)
		t.Assert(v.Val().(*gredis.Subscription).Channel, "gf")

		v, err = conn.Receive(ctx)
		t.AssertNil(err)
		t.Assert(v.Val().(*gredis.Message).Channel, "gf")
		t.Assert(v.Val().(*gredis.Message).Payload, "test")

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
