// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

var ctx = context.TODO()

func Test_NewSessionId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		id1 := NewSessionId()
		id2 := NewSessionId()
		t.AssertNE(id1, id2)
		t.Assert(len(id1), 32)
	})
}

func Test_Session_RegenerateId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		// 1. Test with memory storage
		storage := NewStorageMemory()
		manager := New(time.Hour, storage)
		session := manager.New(ctx)

		// Store some data
		err := session.Set("key1", "value1")
		t.AssertNil(err)
		err = session.Set("key2", "value2")
		t.AssertNil(err)

		// Get original session id
		oldId := session.MustId()

		// Test regenerate with deleteOld = true
		newId1, err := session.RegenerateId(true)
		t.AssertNil(err)
		t.AssertNE(oldId, newId1)

		// Verify data is preserved
		v1 := session.MustGet("key1")
		t.Assert(v1.String(), "value1")
		v2 := session.MustGet("key2")
		t.Assert(v2.String(), "value2")

		// Verify old session is deleted
		oldSession := manager.New(ctx)
		err = oldSession.SetId(oldId)
		t.AssertNil(err)
		v3 := oldSession.MustGet("key1")
		t.Assert(v3.IsNil(), true)

		// Test regenerate with deleteOld = false
		currentId := newId1
		newId2, err := session.RegenerateId(false)
		t.AssertNil(err)
		t.AssertNE(currentId, newId2)

		// Verify data is preserved in new session
		v4 := session.MustGet("key1")
		t.Assert(v4.String(), "value1")

		// Create another session instance with the previous id
		prevSession := manager.New(ctx)
		err = prevSession.SetId(currentId)
		t.AssertNil(err)
		// Data should still be accessible in previous session
		v5 := prevSession.MustGet("key1")
		t.Assert(v5.String(), "value1")
	})

	gtest.C(t, func(t *gtest.T) {
		// 2. Test with custom id function
		storage := NewStorageMemory()
		manager := New(time.Hour, storage)
		session := manager.New(ctx)

		customId := "custom_session_id"
		err := session.SetIdFunc(func(ttl time.Duration) string {
			return customId
		})
		t.AssertNil(err)

		newId, err := session.RegenerateId(true)
		t.AssertNil(err)
		t.Assert(newId, customId)
	})

	gtest.C(t, func(t *gtest.T) {
		// 3. Test with disabled storage
		storage := &StorageBase{} // implements Storage interface but all methods return ErrorDisabled
		manager := New(time.Hour, storage)
		session := manager.New(ctx)

		// Should still work even with disabled storage
		newId, err := session.RegenerateId(true)
		t.AssertNil(err)
		t.Assert(len(newId), 32)
	})
}

// Test MustRegenerateId
func Test_Session_MustRegenerateId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		storage := NewStorageMemory()
		manager := New(time.Hour, storage)
		session := manager.New(ctx)

		// Normal case should not panic
		t.AssertNil(session.Set("key", "value"))
		newId := session.MustRegenerateId(true)
		t.Assert(len(newId), 32)

		// Test with disabled storage (should not panic)
		storage2 := &StorageBase{}
		manager2 := New(time.Hour, storage2)
		session2 := manager2.New(ctx)
		newId2 := session2.MustRegenerateId(true)
		t.Assert(len(newId2), 32)
	})
}
