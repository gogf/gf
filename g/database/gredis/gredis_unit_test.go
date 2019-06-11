// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
    "github.com/gogf/gf/g/database/gredis"
    "github.com/gogf/gf/g/test/gtest"
    "testing"
    "time"
)

var (
    config = gredis.Config{
        Host : "127.0.0.1",
        Port : 6379,
        Db   : 1,
    }
)

func Test_NewClose(t *testing.T) {
    gtest.Case(t, func() {
        redis := gredis.New(config)
        gtest.AssertNE(redis, nil)
        err := redis.Close()
        gtest.Assert(err, nil)
    })
}

func Test_Do(t *testing.T) {
    gtest.Case(t, func() {
        redis := gredis.New(config)
        defer redis.Close()
        _, err := redis.Do("SET", "k", "v")
        gtest.Assert(err, nil)

        r, err := redis.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r, []byte("v"))

        _, err  = redis.Do("DEL", "k")
        gtest.Assert(err, nil)
        r, err  = redis.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r, nil)
    })
}

func Test_Send(t *testing.T) {
    gtest.Case(t, func() {
        redis := gredis.New(config)
        defer redis.Close()
        err := redis.Send("SET", "k", "v")
        gtest.Assert(err, nil)

        r, err := redis.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r,   []byte("v"))
    })
}

func Test_Stats(t *testing.T) {
    gtest.Case(t, func() {
        redis := gredis.New(config)
        defer redis.Close()
        redis.SetMaxIdle(2)
        redis.SetMaxActive(100)
        redis.SetIdleTimeout(500*time.Millisecond)
        redis.SetMaxConnLifetime(500*time.Millisecond)

        array := make([]*gredis.Conn, 0)
        for i := 0; i < 10; i++ {
            array = append(array, redis.Conn())
        }
        stats := redis.Stats()
        gtest.Assert(stats.ActiveCount, 10)
        gtest.Assert(stats.IdleCount,    0)
        for i := 0; i < 10; i++ {
            array[i].Close()
        }
        stats  = redis.Stats()
        gtest.Assert(stats.ActiveCount,  2)
        gtest.Assert(stats.IdleCount,    2)
        //time.Sleep(3000*time.Millisecond)
        //stats  = redis.Stats()
        //fmt.Println(stats)
        //gtest.Assert(stats.ActiveCount,  0)
        //gtest.Assert(stats.IdleCount,    0)
    })
}

func Test_Conn(t *testing.T) {
    gtest.Case(t, func() {
        redis := gredis.New(config)
        defer redis.Close()
        conn := redis.Conn()
        defer conn.Close()


        r, err := conn.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r,   []byte("v"))

        _, err  = conn.Do("DEL", "k")
        gtest.Assert(err, nil)
        r, err  = conn.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r, nil)
    })
}

func Test_Instance(t *testing.T) {
    gtest.Case(t, func() {
        group := "my-test"
        gredis.SetConfig(config, group)
        defer gredis.RemoveConfig(group)
        redis := gredis.Instance(group)
        defer redis.Close()

        conn := redis.Conn()
        defer conn.Close()

        _, err := conn.Do("SET", "k", "v")
        gtest.Assert(err, nil)

        r, err := conn.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r,   []byte("v"))

        _, err  = conn.Do("DEL", "k")
        gtest.Assert(err, nil)
        r, err  = conn.Do("GET", "k")
        gtest.Assert(err, nil)
        gtest.Assert(r, nil)
    })
}
