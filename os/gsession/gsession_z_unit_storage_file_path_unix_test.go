//go:build unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file checks the session storage boundary for symbolic links on Unix.

package gsession_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
)

// Test_StorageFilePathSymlinkEscape rejects external and unresolved symbolic links.
func Test_StorageFilePathSymlinkEscape(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			root    = t.TempDir()
			path    = filepath.Join(root, "storage")
			outside = filepath.Join(root, "outside")
			ctx     = context.Background()
		)
		gt.AssertNil(os.Mkdir(path, 0o700))
		gt.AssertNil(os.Mkdir(outside, 0o700))
		gt.AssertNil(os.Mkdir(filepath.Join(outside, "nested"), 0o700))
		gt.AssertNil(gsession.NewStorageFile(outside, time.Minute).SetSession(ctx, "original", gmap.NewStrAnyMapFrom(map[string]any{"value": "original"}, true), time.Minute))
		outsidePath := filepath.Join(outside, "original.session")
		before, err := os.ReadFile(outsidePath)
		gt.AssertNil(err)
		gt.AssertNil(os.Symlink(outsidePath, filepath.Join(path, "leaf.session")))
		gt.AssertNil(os.Symlink(outside, filepath.Join(path, "parent")))
		gt.AssertNil(os.Symlink(filepath.Join(outside, "nested"), filepath.Join(path, "back")))
		gt.AssertNil(os.Symlink(filepath.Join(outside, "missing.session"), filepath.Join(path, "dangling.session")))
		gt.AssertNil(os.Symlink(filepath.Join(outside, "missing-parent"), filepath.Join(path, "dangling-parent")))
		storage := gsession.NewStorageFile(path, time.Minute)
		for _, id := range []string{"leaf", "parent/original", "parent/missing", "back/../original", "back/../missing", "dangling", "dangling-parent/missing"} {
			assertStorageFilePathRejected(t, storage, id)
		}
		after, err := os.ReadFile(outsidePath)
		gt.AssertNil(err)
		gt.Assert(after, before)
		_, err = os.Stat(filepath.Join(outside, "missing.session"))
		gt.Assert(os.IsNotExist(err), true)
		_, err = os.Stat(filepath.Join(outside, "missing-parent"))
		gt.Assert(os.IsNotExist(err), true)
	})
}

// Test_StorageFilePathSymlinkCompatibility accepts configured roots and links inside the root.
func Test_StorageFilePathSymlinkCompatibility(t *testing.T) {
	gtest.C(t, func(gt *gtest.T) {
		var (
			ctx  = context.Background()
			root = t.TempDir()
			path = filepath.Join(root, "storage")
		)
		gt.AssertNil(os.Mkdir(path, 0o700))
		gt.AssertNil(os.MkdirAll(filepath.Join(path, "nested", "deeper"), 0o700))
		alias := filepath.Join(root, "alias")
		gt.AssertNil(os.Symlink(path, alias))
		gt.AssertNil(os.Symlink(filepath.Join(path, "nested", "deeper"), filepath.Join(path, "parent")))
		storage := gsession.NewStorageFile(alias, time.Minute)
		gt.AssertNil(storage.SetSession(ctx, "original", gmap.NewStrAnyMapFrom(map[string]any{"value": "original"}, true), time.Minute))
		gt.AssertNil(os.Symlink(filepath.Join(path, "original.session"), filepath.Join(path, "leaf.session")))
		for _, id := range []string{"leaf", "parent/custom", "parent/../custom"} {
			gt.AssertNil(storage.SetSession(ctx, id, gmap.NewStrAnyMapFrom(map[string]any{"value": "changed"}, true), time.Minute))
			if id == "parent/../custom" {
				// The filesystem resolves .. after the link, rather than cleaning it to the root.
				_, err := os.Stat(filepath.Join(path, "nested", "custom.session"))
				gt.AssertNil(err)
				_, err = os.Stat(filepath.Join(path, "custom.session"))
				gt.Assert(os.IsNotExist(err), true)
			}
			data, err := storage.GetSession(ctx, id, time.Minute)
			gt.AssertNil(err)
			gt.Assert(data.Get("value"), "changed")
			gt.AssertNil(storage.UpdateTTL(ctx, id, time.Second))
			gt.AssertNil(storage.RemoveAll(ctx, id))
		}
		data, err := storage.GetSession(ctx, "original", time.Minute)
		gt.AssertNil(err)
		gt.Assert(data.Get("value"), "original")
	})
}
