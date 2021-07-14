// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession_test

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gsession"
	"testing"
	"time"

	"github.com/gogf/gf/test/gtest"
)

func Test_StorageMemory(t *testing.T) {
	storage := gsession.NewStorageMemory()
	manager := gsession.New(time.Second, storage)
	sessionId := ""
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO())
		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		s.Sets(g.Map{
			"k3": "v3",
			"k4": "v4",
		})
		t.Assert(s.IsDirty(), true)
		sessionId = s.Id()
	})

	time.Sleep(500 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO(), sessionId)
		t.Assert(s.Get("k1"), "v1")
		t.Assert(s.Get("k2"), "v2")
		t.Assert(s.Get("k3"), "v3")
		t.Assert(s.Get("k4"), "v4")
		t.Assert(len(s.Map()), 4)
		t.Assert(s.Map()["k1"], "v1")
		t.Assert(s.Map()["k4"], "v4")
		t.Assert(s.Id(), sessionId)
		t.Assert(s.Size(), 4)
		t.Assert(s.Contains("k1"), true)
		t.Assert(s.Contains("k3"), true)
		t.Assert(s.Contains("k5"), false)
		s.Remove("k4")
		t.Assert(s.Size(), 3)
		t.Assert(s.Contains("k3"), true)
		t.Assert(s.Contains("k4"), false)
		s.RemoveAll()
		t.Assert(s.Size(), 0)
		t.Assert(s.Contains("k1"), false)
		t.Assert(s.Contains("k2"), false)
		s.Sets(g.Map{
			"k5": "v5",
			"k6": "v6",
		})
		t.Assert(s.Size(), 2)
		t.Assert(s.Contains("k5"), true)
		t.Assert(s.Contains("k6"), true)
	})

	time.Sleep(1000 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO(), sessionId)
		t.Assert(s.Size(), 0)
		t.Assert(s.Get("k5"), nil)
		t.Assert(s.Get("k6"), nil)
	})
}
