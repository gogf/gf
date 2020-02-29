// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession_test

import (
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gsession"
	"testing"
	"time"

	"github.com/gogf/gf/test/gtest"
)

func Test_StorageRedisHashTable(t *testing.T) {
	redis, err := gredis.NewFromStr("127.0.0.1:6379,0")
	gtest.Assert(err, nil)

	storage := gsession.NewStorageRedisHashTable(redis)
	manager := gsession.New(time.Second, storage)
	sessionId := ""
	gtest.Case(t, func() {
		s := manager.New()
		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		s.Sets(g.Map{
			"k3": "v3",
			"k4": "v4",
		})
		gtest.Assert(s.IsDirty(), true)
		sessionId = s.Id()
	})
	gtest.Case(t, func() {
		s := manager.New(sessionId)
		gtest.Assert(s.Get("k1"), "v1")
		gtest.Assert(s.Get("k2"), "v2")
		gtest.Assert(s.Get("k3"), "v3")
		gtest.Assert(s.Get("k4"), "v4")
		gtest.Assert(len(s.Map()), 4)
		gtest.Assert(s.Map()["k1"], "v1")
		gtest.Assert(s.Map()["k4"], "v4")
		gtest.Assert(s.Id(), sessionId)
		gtest.Assert(s.Size(), 4)
		gtest.Assert(s.Contains("k1"), true)
		gtest.Assert(s.Contains("k3"), true)
		gtest.Assert(s.Contains("k5"), false)
		s.Remove("k4")
		gtest.Assert(s.Size(), 3)
		gtest.Assert(s.Contains("k3"), true)
		gtest.Assert(s.Contains("k4"), false)
		s.RemoveAll()
		gtest.Assert(s.Size(), 0)
		gtest.Assert(s.Contains("k1"), false)
		gtest.Assert(s.Contains("k2"), false)
		s.Sets(g.Map{
			"k5": "v5",
			"k6": "v6",
		})
		gtest.Assert(s.Size(), 2)
		gtest.Assert(s.Contains("k5"), true)
		gtest.Assert(s.Contains("k6"), true)
		s.Close()
	})

	time.Sleep(1500 * time.Millisecond)
	gtest.Case(t, func() {
		s := manager.New(sessionId)
		gtest.Assert(s.Size(), 0)
		gtest.Assert(s.Get("k5"), nil)
		gtest.Assert(s.Get("k6"), nil)
	})
}

func Test_StorageRedisHashTablePrefix(t *testing.T) {
	redis, err := gredis.NewFromStr("127.0.0.1:6379,0")
	gtest.Assert(err, nil)

	prefix := "s_"
	storage := gsession.NewStorageRedisHashTable(redis, prefix)
	manager := gsession.New(time.Second, storage)
	sessionId := ""
	gtest.Case(t, func() {
		s := manager.New()
		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		s.Sets(g.Map{
			"k3": "v3",
			"k4": "v4",
		})
		gtest.Assert(s.IsDirty(), true)
		sessionId = s.Id()
	})
	gtest.Case(t, func() {
		s := manager.New(sessionId)
		gtest.Assert(s.Get("k1"), "v1")
		gtest.Assert(s.Get("k2"), "v2")
		gtest.Assert(s.Get("k3"), "v3")
		gtest.Assert(s.Get("k4"), "v4")
		gtest.Assert(len(s.Map()), 4)
		gtest.Assert(s.Map()["k1"], "v1")
		gtest.Assert(s.Map()["k4"], "v4")
		gtest.Assert(s.Id(), sessionId)
		gtest.Assert(s.Size(), 4)
		gtest.Assert(s.Contains("k1"), true)
		gtest.Assert(s.Contains("k3"), true)
		gtest.Assert(s.Contains("k5"), false)
		s.Remove("k4")
		gtest.Assert(s.Size(), 3)
		gtest.Assert(s.Contains("k3"), true)
		gtest.Assert(s.Contains("k4"), false)
		s.RemoveAll()
		gtest.Assert(s.Size(), 0)
		gtest.Assert(s.Contains("k1"), false)
		gtest.Assert(s.Contains("k2"), false)
		s.Sets(g.Map{
			"k5": "v5",
			"k6": "v6",
		})
		gtest.Assert(s.Size(), 2)
		gtest.Assert(s.Contains("k5"), true)
		gtest.Assert(s.Contains("k6"), true)
		s.Close()
	})

	time.Sleep(1500 * time.Millisecond)
	gtest.Case(t, func() {
		s := manager.New(sessionId)
		gtest.Assert(s.Size(), 0)
		gtest.Assert(s.Get("k5"), nil)
		gtest.Assert(s.Get("k6"), nil)
	})
}
