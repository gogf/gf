// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/guid"
	"testing"
	"time"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/test/gtest"
	redis2 "github.com/gomodule/redigo/redis"
)

var (
	config = gredis.Config{
		Host: "127.0.0.1",
		Port: 6379,
		Db:   1,
	}
)

func Test_NewClose(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		t.AssertNE(redis, nil)
		err := redis.Close()
		t.Assert(err, nil)
	})
}

func Test_Do(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		defer redis.Close()
		_, err := redis.Do("SET", "k", "v")
		t.Assert(err, nil)

		r, err := redis.Do("GET", "k")
		t.Assert(err, nil)
		t.Assert(r, []byte("v"))

		_, err = redis.Do("DEL", "k")
		t.Assert(err, nil)
		r, err = redis.Do("GET", "k")
		t.Assert(err, nil)
		t.Assert(r, nil)
	})
}

func Test_Stats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		defer redis.Close()
		redis.SetMaxIdle(2)
		redis.SetMaxActive(100)
		redis.SetIdleTimeout(500 * time.Millisecond)
		redis.SetMaxConnLifetime(500 * time.Millisecond)

		array := make([]*gredis.Conn, 0)
		for i := 0; i < 10; i++ {
			array = append(array, redis.Conn())
		}
		stats := redis.Stats()
		t.Assert(stats.ActiveCount, 10)
		t.Assert(stats.IdleCount, 0)
		for i := 0; i < 10; i++ {
			array[i].Close()
		}
		stats = redis.Stats()
		t.Assert(stats.ActiveCount, 2)
		t.Assert(stats.IdleCount, 2)
		//time.Sleep(3000*time.Millisecond)
		//stats  = redis.Stats()
		//fmt.Println(stats)
		//t.Assert(stats.ActiveCount,  0)
		//t.Assert(stats.IdleCount,    0)
	})
}

func Test_Conn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		defer redis.Close()
		conn := redis.Conn()
		defer conn.Close()

		key := gconv.String(gtime.TimestampNano())
		value := []byte("v")
		r, err := conn.Do("SET", key, value)
		t.Assert(err, nil)

		r, err = conn.Do("GET", key)
		t.Assert(err, nil)
		t.Assert(r, value)

		_, err = conn.Do("DEL", key)
		t.Assert(err, nil)
		r, err = conn.Do("GET", key)
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
		defer redis.Close()

		conn := redis.Conn()
		defer conn.Close()

		_, err := conn.Do("SET", "k", "v")
		t.Assert(err, nil)

		r, err := conn.Do("GET", "k")
		t.Assert(err, nil)
		t.Assert(r, []byte("v"))

		_, err = conn.Do("DEL", "k")
		t.Assert(err, nil)
		r, err = conn.Do("GET", "k")
		t.Assert(err, nil)
		t.Assert(r, nil)
	})
}

func Test_Error(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config1 := gredis.Config{
			Host:           "127.0.0.2",
			Port:           6379,
			Db:             1,
			ConnectTimeout: time.Second,
		}
		redis := gredis.New(config1)
		_, err := redis.Do("info")
		t.AssertNE(err, nil)

		config1 = gredis.Config{
			Host: "127.0.0.1",
			Port: 6379,
			Db:   1,
			Pass: "666666",
		}
		redis = gredis.New(config1)
		_, err = redis.Do("info")
		t.AssertNE(err, nil)

		config1 = gredis.Config{
			Host: "127.0.0.1",
			Port: 6379,
			Db:   100,
		}
		redis = gredis.New(config1)
		_, err = redis.Do("info")
		t.AssertNE(err, nil)

		redis = gredis.Instance("gf")
		t.Assert(redis == nil, true)
		gredis.ClearConfig()

		redis = gredis.New(config)
		defer redis.Close()
		_, err = redis.DoVar("SET", "k", "v")
		t.Assert(err, nil)

		v, err := redis.DoVar("GET", "k")
		t.Assert(err, nil)
		t.Assert(v.String(), "v")

		conn := redis.GetConn()
		_, err = conn.DoVar("SET", "k", "v")
		t.Assert(err, nil)

		//v, err = conn.ReceiveVar()
		//t.Assert(err, nil)
		//t.Assert(v.String(), "v")

		psc := redis2.PubSubConn{Conn: conn}
		psc.Subscribe("gf")
		redis.DoVar("PUBLISH", "gf", "gf test")
		go func() {
			for {
				v, _ := conn.ReceiveVar()
				switch obj := v.Val().(type) {
				case redis2.Message:
					t.Assert(string(obj.Data), "gf test")
				case redis2.Subscription:

				}
			}
		}()

		time.Sleep(time.Second)
	})
}

func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		defer func() {
			redis.Do("DEL", "key-true")
			redis.Do("DEL", "key-false")
		}()

		_, err := redis.Do("SET", "key-true", true)
		t.Assert(err, nil)

		_, err = redis.Do("SET", "key-false", false)
		t.Assert(err, nil)

		r, err := redis.DoVar("GET", "key-true")
		t.Assert(err, nil)
		t.Assert(r.Bool(), true)

		r, err = redis.DoVar("GET", "key-false")
		t.Assert(err, nil)
		t.Assert(r.Bool(), false)
	})
}

func Test_Int(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		key := guid.S()
		defer redis.Do("DEL", key)

		_, err := redis.Do("SET", key, 1)
		t.Assert(err, nil)

		r, err := redis.DoVar("GET", key)
		t.Assert(err, nil)
		t.Assert(r.Int(), 1)
	})
}

func Test_HSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		redis := gredis.New(config)
		key := guid.S()
		defer redis.Do("DEL", key)

		_, err := redis.Do("HSET", key, "name", "john")
		t.Assert(err, nil)

		r, err := redis.DoVar("HGETALL", key)
		t.Assert(err, nil)
		t.Assert(r.Strings(), g.ArrayStr{"name", "john"})
	})
}

func Test_HGetAll1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var err error
		redis := gredis.New(config)
		key := guid.S()
		defer redis.Do("DEL", key)

		_, err = redis.Do("HSET", key, "id", 100)
		t.Assert(err, nil)
		_, err = redis.Do("HSET", key, "name", "john")
		t.Assert(err, nil)

		r, err := redis.DoVar("HGETALL", key)
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
			err   error
			key   = guid.S()
			redis = gredis.New(config)
		)
		defer redis.Do("DEL", key)

		_, err = redis.Do("HSET", key, "id", 100)
		t.Assert(err, nil)
		_, err = redis.Do("HSET", key, "name", "john")
		t.Assert(err, nil)

		result, err := redis.DoVar("HGETALL", key)
		t.Assert(err, nil)

		t.Assert(gconv.Uint(result.MapStrVar()["id"]), 100)
		t.Assert(result.MapStrVar()["id"].Uint(), 100)
	})
}

func Test_Auto_Marshal(t *testing.T) {
	var (
		err   error
		redis = gredis.New(config)
		key   = guid.S()
	)
	defer redis.Do("DEL", key)

	type User struct {
		Id   int
		Name string
	}

	gtest.C(t, func(t *gtest.T) {
		user := &User{
			Id:   10000,
			Name: "john",
		}

		_, err = redis.Do("SET", key, user)
		t.Assert(err, nil)

		r, err := redis.DoVar("GET", key)
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
