// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_StorageFile(t *testing.T) {
	storage := gsession.NewStorageFile()
	manager := gsession.New(time.Second, storage)
	sessionId := ""
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO())
		defer s.Close()
		s.Set("k1", "v1")
		s.Set("k2", "v2")
		s.MustSet("k3", "v3")
		s.MustSet("k4", "v4")
		s.SetMap(g.Map{
			"kmap1": "kval1",
			"kmap2": "kval2",
		})
		s.MustSetMap(g.Map{
			"kmap3": "kval3",
			"kmap4": "kval4",
		})
		t.Assert(s.IsDirty(), true)
		sessionId = s.MustId()
	})

	time.Sleep(500 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO(), sessionId)
		t.Assert(s.MustGet("k1"), "v1")
		t.Assert(s.MustGet("k2"), "v2")
		t.Assert(s.MustGet("k3"), "v3")
		t.Assert(s.MustGet("k4"), "v4")
		t.Assert(len(s.MustData()), 8)
		t.Assert(s.MustData()["k1"], "v1")
		t.Assert(s.MustData()["k4"], "v4")
		t.Assert(s.MustId(), sessionId)
		t.Assert(s.MustSize(), 8)
		t.Assert(s.MustContains("k1"), true)
		t.Assert(s.MustContains("k3"), true)
		t.Assert(s.MustContains("k5"), false)
		s.Remove("k4")
		s.MustRemove("k4")
		t.Assert(s.MustSize(), 7)
		t.Assert(s.MustContains("k3"), true)
		t.Assert(s.MustContains("k4"), false)
		s.RemoveAll()
		t.Assert(s.MustSize(), 0)
		t.Assert(s.MustContains("k1"), false)
		t.Assert(s.MustContains("k2"), false)
		s.SetMap(g.Map{
			"k5": "v5",
			"k6": "v6",
		})
		t.Assert(s.MustSize(), 2)
		t.Assert(s.MustContains("k5"), true)
		t.Assert(s.MustContains("k6"), true)
	})

	time.Sleep(1000 * time.Millisecond)
	gtest.C(t, func(t *gtest.T) {
		s := manager.New(context.TODO(), sessionId)
		t.Assert(s.MustSize(), 0)
		t.Assert(s.MustGet("k5"), nil)
		t.Assert(s.MustGet("k6"), nil)
	})
}
