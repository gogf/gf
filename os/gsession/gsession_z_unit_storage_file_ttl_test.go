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

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_StorageFile_updateSessionTTl covers updateSessionTTl: it writes a fresh
// timestamp into an existing session file, and returns an error for a missing one.
func Test_StorageFile_updateSessionTTl(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		dir := t.TempDir()
		ttl := time.Hour
		s := NewStorageFile(dir, ttl)

		data := gmap.NewStrAnyMapFrom(map[string]any{"k": "v"}, true)
		t.AssertNil(s.SetSession(ctx, "s1", data, ttl))

		// Happy path: existing session gets its timestamp refreshed.
		t.AssertNil(s.updateSessionTTl(ctx, "s1"))

		// Missing session: OpenWithFlag fails.
		t.AssertNE(s.updateSessionTTl(ctx, "missing"), nil)

		// Path-escaping id: sessionFilePath returns an error.
		t.AssertNE(s.updateSessionTTl(ctx, "../outside"), nil)
	})
}

// Test_StorageFile_GetSession_Expired covers the branch where the stored timestamp
// is older than the TTL, so GetSession returns nil without an error.
func Test_StorageFile_GetSession_Expired(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		dir := t.TempDir()
		ttl := time.Millisecond // extremely short TTL
		s := NewStorageFile(dir, ttl)

		data := gmap.NewStrAnyMapFrom(map[string]any{"k": "v"}, true)
		t.AssertNil(s.SetSession(ctx, "exp", data, ttl))

		// Wait for the session to expire, then read it.
		time.Sleep(10 * time.Millisecond)
		got, err := s.GetSession(ctx, "exp", ttl)
		t.AssertNil(err)
		t.Assert(got, nil)
	})
}

// Test_StorageFile_CryptoGetSet covers SetSession and GetSession with encryption
// enabled, exercising the AES encrypt/decrypt branches.
func Test_StorageFile_CryptoGetSet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		dir := t.TempDir()
		ttl := time.Hour
		s := NewStorageFile(dir, ttl)
		s.SetCryptoEnabled(true)

		data := gmap.NewStrAnyMapFrom(map[string]any{"secret": "value"}, true)
		t.AssertNil(s.SetSession(ctx, "enc", data, ttl))

		got, err := s.GetSession(ctx, "enc", ttl)
		t.AssertNil(err)
		t.AssertNE(got, nil)
		t.Assert(got.Get("secret"), "value")
	})
}
