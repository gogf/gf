//go:build unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file verifies internal session file maintenance and early-write cleanup on Unix.

package gsession

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// Test_StorageFilePathInternalMaintenance keeps external links outside TTL and expiry writes.
func Test_StorageFilePathInternalMaintenance(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx         = context.Background()
			path        = t.TempDir()
			outside     = t.TempDir()
			outsidePath = filepath.Join(outside, "sentinel.session")
			linkPath    = filepath.Join(path, "linked.session")
			expiredPath = filepath.Join(path, "expired.session")
			storage     = NewStorageFile(path, time.Second)
			sentinel    = append(gbinary.EncodeInt64(0), []byte(`{"value":"outside"}`)...)
		)
		gt.AssertNil(os.WriteFile(outsidePath, sentinel, 0o600))
		gt.AssertNil(os.Symlink(outsidePath, linkPath))

		// An actual TTL write must reject the link before changing the external timestamp.
		err := storage.updateSessionTTl(ctx, "linked")
		gt.AssertNE(err, nil)
		gt.Assert(gerror.Code(err).Code(), gcode.CodeInvalidParameter.Code())
		after, err := os.ReadFile(outsidePath)
		gt.AssertNil(err)
		gt.AssertEQ(after, sentinel)

		// The direct expiry path must reject the external link instead of deleting it.
		err = storage.checkAndClearSessionFile(ctx, linkPath)
		gt.AssertNE(err, nil)
		gt.Assert(gerror.Code(err).Code(), gcode.CodeInvalidParameter.Code())
		after, err = os.ReadFile(outsidePath)
		gt.AssertNil(err)
		gt.AssertEQ(after, sentinel)

		// A real scan removes a local expired file while preserving the rejected link.
		gt.AssertNil(os.WriteFile(expiredPath, sentinel, 0o600))
		storage.timelyClearExpiredSessionFile(ctx)
		_, err = os.Stat(expiredPath)
		gt.Assert(os.IsNotExist(err), true)
		info, err := os.Lstat(linkPath)
		gt.AssertNil(err)
		gt.Assert(info.Mode()&os.ModeSymlink != 0, true)
		after, err = os.ReadFile(outsidePath)
		gt.AssertNil(err)
		gt.AssertEQ(after, sentinel)
		entries, err := os.ReadDir(path)
		gt.AssertNil(err)
		gt.Assert(len(entries), 1)
		gt.Assert(entries[0].Name(), "linked.session")
	})
}

// Test_StorageFilePathInternalUpdateTTLQueue rejects invalid ids before adding them to the queue.
func Test_StorageFilePathInternalUpdateTTLQueue(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx     = context.Background()
			path    = t.TempDir()
			storage = NewStorageFile(path, time.Minute)
		)
		for _, ttl := range []time.Duration{time.Millisecond, time.Minute} {
			err := storage.UpdateTTL(ctx, "../outside", ttl)
			gt.AssertNE(err, nil)
			gt.Assert(gerror.Code(err).Code(), gcode.CodeInvalidParameter.Code())
			gt.Assert(storage.updatingIdSet.Size(), 0)
			gt.Assert(storage.updatingIdSet.Contains("../outside"), false)
		}
		// A valid id still enters the queue, so validation does not disable TTL updates.
		gt.AssertNil(storage.UpdateTTL(ctx, "valid", time.Minute))
		gt.Assert(storage.updatingIdSet.Size(), 1)
		gt.Assert(storage.updatingIdSet.Contains("valid"), true)
	})
}

// Test_StorageFilePathInternalWriteCleanup exercises Stat failure after creating a private file.
func Test_StorageFilePathInternalWriteCleanup(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx      = context.Background()
			path     = t.TempDir()
			loopPath = filepath.Join(path, "loop.session")
			storage  = NewStorageFile(path, time.Minute)
		)
		gt.AssertNil(os.Symlink(loopPath, loopPath))
		// Calling the private writer bypasses the public path guard to test its own cleanup.
		err := storage.setSessionFile(ctx, loopPath, []byte(`{"value":"test"}`))
		gt.AssertNE(err, nil)
		gt.Assert(errors.Is(err, syscall.ELOOP), true)
		entries, err := os.ReadDir(path)
		gt.AssertNil(err)
		gt.Assert(len(entries), 1)
		gt.Assert(entries[0].Name(), "loop.session")
		gt.Assert(entries[0].Type()&os.ModeSymlink != 0, true)
	})
}

// Test_StorageFilePathInternalFilesystemRoot preserves the root while resolving a missing file.
func Test_StorageFilePathInternalFilesystemRoot(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			root    = filepath.VolumeName(os.TempDir()) + string(filepath.Separator)
			id      = "gsession-boundary-" + guid.S()
			probe   = filepath.Join(root, id+".session")
			storage = &StorageFile{path: root}
		)
		// Use only metadata checks against a random missing name; never create a root file.
		_, err := os.Lstat(probe)
		gt.Assert(os.IsNotExist(err), true)
		gt.AssertNil(storage.checkSessionFilePath(probe))
		path, err := storage.sessionFilePath(id)
		gt.AssertNil(err)
		gt.AssertEQ(filepath.Clean(path), probe)
		_, err = os.Lstat(path)
		gt.Assert(os.IsNotExist(err), true)
	})
}

// Test_StorageFilePathInternalRelativeDefault clears absolute scan results under a relative root.
func Test_StorageFilePathInternalRelativeDefault(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx          = context.Background()
			absolute     = t.TempDir()
			expiredPath  = filepath.Join(absolute, "expired.session")
			savedDefault = DefaultStorageFilePath
		)
		cwd, err := os.Getwd()
		gt.AssertNil(err)
		relative, err := filepath.Rel(cwd, absolute)
		gt.AssertNil(err)
		// Configure a relative default without changing the working directory.
		DefaultStorageFilePath = relative
		defer func() {
			DefaultStorageFilePath = savedDefault
		}()
		storage := NewStorageFile("", time.Second)
		gt.AssertEQ(storage.path, relative)
		gt.AssertNil(os.WriteFile(expiredPath, append(gbinary.EncodeInt64(0), []byte(`{"value":"expired"}`)...), 0o600))
		files, err := gfile.ScanDirFile(storage.path, "*.session", false)
		gt.AssertNil(err)
		gt.Assert(len(files), 1)
		gt.Assert(filepath.IsAbs(files[0]), true)
		// The same scanner result used by periodic maintenance must pass the guard.
		gt.AssertNil(storage.checkAndClearSessionFile(ctx, files[0]))
		_, err = os.Stat(expiredPath)
		gt.Assert(os.IsNotExist(err), true)
	})
}

// Test_StorageFilePathInternalWriteMissingParent handles CreateTemp failure without leaving files.
func Test_StorageFilePathInternalWriteMissingParent(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx     = context.Background()
			path    = t.TempDir()
			storage = NewStorageFile(path, time.Minute)
		)
		// Bypass public validation to exercise the writer's real missing-directory error.
		err := storage.setSessionFile(ctx, filepath.Join(path, "missing", "session.session"), []byte(`{"value":"test"}`))
		gt.AssertNE(err, nil)
		gt.Assert(errors.Is(err, os.ErrNotExist), true)
		entries, err := os.ReadDir(path)
		gt.AssertNil(err)
		gt.Assert(len(entries), 0)
	})
}
