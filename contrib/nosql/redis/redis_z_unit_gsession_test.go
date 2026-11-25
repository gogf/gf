// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// ─────────────────────── StorageRedis ────────────────────────────

func Test_StorageRedis_NilPanics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() { t.AssertNE(recover(), nil) }()
		gsession.NewStorageRedis(nil)
	})
}

func Test_StorageRedis_GetSession_New(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedis(redis, "gstest:")
		id := guid.S()
		m, err := storage.GetSession(ctx, id, time.Second*10)
		t.AssertNil(err)
		t.Assert(m, nil) // new session returns nil
	})
}

func Test_StorageRedis_SetSession_GetSession(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedis(redis, "gstest:")
		id := guid.S()
		defer redis.Del(ctx, "gstest:"+id)

		manager := gsession.New(time.Second*30, storage)
		s := manager.New(ctx, id)
		t.AssertNil(s.Set("foo", "bar"))
		t.AssertNil(s.Set("num", 42))
		t.AssertNil(s.Close())

		s2 := manager.New(ctx, id)
		v, err := s2.Get("foo")
		t.AssertNil(err)
		t.Assert(v, "bar")
		v2, err := s2.Get("num")
		t.AssertNil(err)
		t.Assert(v2, 42)
		t.AssertNil(s2.Close())
	})
}

func Test_StorageRedis_RemoveAll(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedis(redis, "gstest:")
		id := guid.S()

		manager := gsession.New(time.Second*30, storage)
		s := manager.New(ctx, id)
		t.AssertNil(s.Set("k", "v"))
		t.AssertNil(s.Close())

		t.AssertNil(storage.RemoveAll(ctx, id))

		m, err := storage.GetSession(ctx, id, time.Second*30)
		t.AssertNil(err)
		t.Assert(m, nil)
	})
}

func Test_StorageRedis_UpdateTTL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedis(redis, "gstest:")
		id := guid.S()
		defer redis.Del(ctx, "gstest:"+id)

		manager := gsession.New(time.Second*30, storage)
		s := manager.New(ctx, id)
		t.AssertNil(s.Set("k", "v"))
		t.AssertNil(s.Close())

		// UpdateTTL must not return an error for a known session.
		t.AssertNil(storage.UpdateTTL(ctx, id, time.Second*60))
	})
}

// ─────────────────── StorageRedisHashTable ────────────────────────

func Test_StorageRedisHashTable_NilPanics(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func() { t.AssertNE(recover(), nil) }()
		gsession.NewStorageRedisHashTable(nil)
	})
}

func Test_StorageRedisHashTable_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedisHashTable(redis, "gshttest:")
		id := guid.S()
		defer storage.RemoveAll(ctx, id)

		t.AssertNil(storage.Set(ctx, id, "field1", "value1", time.Second*30))
		v, err := storage.Get(ctx, id, "field1")
		t.AssertNil(err)
		t.Assert(v, "value1")
	})
}

func Test_StorageRedisHashTable_SetMap_Data(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedisHashTable(redis, "gshttest:")
		id := guid.S()
		defer storage.RemoveAll(ctx, id)

		t.AssertNil(storage.SetMap(ctx, id, map[string]any{"a": 1, "b": "two"}, time.Second*30))
		data, err := storage.Data(ctx, id)
		t.AssertNil(err)
		t.Assert(data["a"], 1)
		t.Assert(data["b"], "two")
	})
}

func Test_StorageRedisHashTable_GetSize(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedisHashTable(redis, "gshttest:")
		id := guid.S()
		defer storage.RemoveAll(ctx, id)

		t.AssertNil(storage.Set(ctx, id, "x", 1, time.Second*30))
		t.AssertNil(storage.Set(ctx, id, "y", 2, time.Second*30))
		size, err := storage.GetSize(ctx, id)
		t.AssertNil(err)
		t.Assert(size, 2)
	})
}

func Test_StorageRedisHashTable_Remove(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := gsession.NewStorageRedisHashTable(redis, "gshttest:")
		id := guid.S()
		defer storage.RemoveAll(ctx, id)

		t.AssertNil(storage.Set(ctx, id, "del", "me", time.Second*30))
		t.AssertNil(storage.Remove(ctx, id, "del"))
		v, err := storage.Get(ctx, id, "del")
		t.AssertNil(err)
		t.Assert(v, nil)
	})
}

func Test_StorageRedisHashTable_GetSetSession(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		manager := gsession.New(time.Second*30, gsession.NewStorageRedisHashTable(redis, "gshttest:"))
		s := manager.New(ctx)
		t.AssertNil(s.Set("hello", "world"))
		id, err := s.Id()
		t.AssertNil(err)
		t.AssertNil(s.Close())

		s2 := manager.New(ctx, id)
		v, err := s2.Get("hello")
		t.AssertNil(err)
		t.Assert(v, "world")
		t.AssertNil(s2.Close())
	})
}

func Test_StorageRedisHashTable_ReadRefreshesTTL(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ttl     = 3 * time.Second
			manager = gsession.New(ttl, gsession.NewStorageRedisHashTable(redis, "test_session_ttl:"))
		)
		s := manager.New(ctx)
		t.AssertNil(s.Set("k", "v"))
		id, err := s.Id()
		t.AssertNil(err)
		t.AssertNil(s.Close())
		defer redis.Del(ctx, "test_session_ttl:"+id)

		// Requests that only read the session must keep it alive.
		for range 3 {
			time.Sleep(time.Second + 500*time.Millisecond)
			s = manager.New(ctx, id)
			v, err := s.Get("k")
			t.AssertNil(err)
			t.Assert(v, "v")
			t.AssertNil(s.Close())
		}
	})
}
