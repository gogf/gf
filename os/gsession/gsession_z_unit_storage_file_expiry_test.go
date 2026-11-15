// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

//go:build !windows

// This file tests session file expiry: timelyClearExpiredSessionFile and
// checkAndClearSessionFile. It is skipped on Windows because Windows does not
// allow removing a file that is still open, which checkAndClearSessionFile does
// after opening the file to read its timestamp.

package gsession

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func writeSessionFileRaw(t *gtest.T, path string, timestampMilli int64) {
	t.Helper()
	body := []byte(`{"x":1}`)
	buf := make([]byte, 8+len(body))
	copy(buf[:8], gbinary.EncodeInt64(timestampMilli))
	copy(buf[8:], body)
	t.AssertNil(os.WriteFile(path, buf, os.ModePerm))
}

// Test_StorageFile_checkAndClearSessionFile exercises checkAndClearSessionFile
// directly: expired files are removed and fresh files are kept.
func Test_StorageFile_checkAndClearSessionFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		dir := t.TempDir()
		ttl := 100 * time.Millisecond
		s := NewStorageFile(dir, ttl)

		// Expired file: timestamp so old it has long since expired.
		expiredPath := filepath.Join(dir, "expired.session")
		writeSessionFileRaw(t, expiredPath, gtime.TimestampMilli()-int64(time.Hour/time.Millisecond))
		err := s.checkAndClearSessionFile(ctx, expiredPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(expiredPath), false)

		// Fresh file: timestamp is now, should not be removed.
		freshPath := filepath.Join(dir, "fresh.session")
		writeSessionFileRaw(t, freshPath, gtime.TimestampMilli())
		err = s.checkAndClearSessionFile(ctx, freshPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(freshPath), true)
	})
}

// Test_StorageFile_timelyClearExpiredSessionFile exercises the scan-and-clear loop:
// expired session files are deleted and non-expired ones survive.
func Test_StorageFile_timelyClearExpiredSessionFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		dir := t.TempDir()
		ttl := 100 * time.Millisecond
		s := NewStorageFile(dir, ttl)

		expiredPath := filepath.Join(dir, "old.session")
		freshPath := filepath.Join(dir, "new.session")
		writeSessionFileRaw(t, expiredPath, gtime.TimestampMilli()-int64(time.Hour/time.Millisecond))
		writeSessionFileRaw(t, freshPath, gtime.TimestampMilli())

		s.timelyClearExpiredSessionFile(ctx)

		t.Assert(gfile.Exists(expiredPath), false)
		t.Assert(gfile.Exists(freshPath), true)
	})
}
