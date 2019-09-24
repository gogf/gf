// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"testing"
	"time"

	"github.com/gogf/gf/test/gtest"
)

func Test_Manager_Basic(t *testing.T) {
	ttl := time.Second
	storage := NewStorageFile()
	manager := New(ttl, storage)
	sessionId := ""
	gtest.Case(t, func() {
		session := manager.New()
		defer session.Close()
		session.Set("k1", "v1")
		session.Set("k2", "v2")
		gtest.Assert(session.Get("k1"), "v1")
		gtest.Assert(session.Get("k2"), "v2")
		sessionId = session.Id()
	})

	time.Sleep(500 * time.Millisecond)
	gtest.Case(t, func() {
		gtest.AssertNE(sessionId, "")
		gtest.Assert(manager.New(sessionId).Get("k1"), "v1")
		gtest.Assert(manager.New(sessionId).Get("k2"), "v2")
	})

	time.Sleep(1000 * time.Millisecond)
	gtest.Case(t, func() {
		gtest.AssertNE(sessionId, "")
		gtest.Assert(manager.New(sessionId).Get("k1"), nil)
		gtest.Assert(manager.New(sessionId).Get("k2"), nil)
	})
}
