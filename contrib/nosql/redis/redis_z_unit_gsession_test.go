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
)

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
