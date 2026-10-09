// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// This file checks the session file path boundary and custom id compatibility.

package gsession_test

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gsession"
	"github.com/gogf/gf/v2/test/gtest"
)

// assertStorageFilePathRejected checks every session-id file operation rejects the path.
func assertStorageFilePathRejected(t *testing.T, storage *gsession.StorageFile, id string) {
	t.Helper()
	ctx := context.Background()
	t.Run("GetSession", func(t *testing.T) {
		gtest.C(t, func(gt *gtest.T) {
			data, err := storage.GetSession(ctx, id, time.Minute)
			gt.AssertNil(data)
			gt.AssertNE(err, nil)
		})
	})
	operations := []struct {
		name string
		run  func() error
	}{
		{"SetSession", func() error {
			return storage.SetSession(ctx, id, gmap.NewStrAnyMapFrom(map[string]any{"value": "changed"}, true), time.Minute)
		}},
		{"RemoveAll", func() error { return storage.RemoveAll(ctx, id) }},
		{"UpdateTTLShort", func() error { return storage.UpdateTTL(ctx, id, time.Second) }},
		{"UpdateTTL", func() error { return storage.UpdateTTL(ctx, id, time.Minute) }},
		{"SessionSize", func() error { _, err := gsession.New(time.Minute, storage).New(ctx, id).Size(); return err }},
	}
	for _, operation := range operations {
		t.Run(operation.name, func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) { gt.AssertNE(operation.run(), nil) })
		})
	}
}

// Test_StorageFilePathTraversal checks ids cannot read, replace, delete or refresh an external file.
func Test_StorageFilePathTraversal(t *testing.T) {
	for _, id := range []string{filepath.Join("..", "outside"), strings.Join([]string{"nested", "..", "..", "outside"}, string(os.PathSeparator))} {
		t.Run(id, func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) {
				root := t.TempDir()
				path := filepath.Join(root, "storage")
				gt.AssertNil(os.Mkdir(path, 0o700))
				gt.AssertNil(os.Mkdir(filepath.Join(path, "nested"), 0o700))
				outside := gsession.NewStorageFile(root, time.Minute)
				gt.AssertNil(outside.SetSession(context.Background(), "outside", gmap.NewStrAnyMapFrom(map[string]any{"value": "original"}, true), time.Minute))
				outsidePath := filepath.Join(root, "outside.session")
				before, err := os.ReadFile(outsidePath)
				gt.AssertNil(err)
				assertStorageFilePathRejected(t, gsession.NewStorageFile(path, time.Minute), id)
				after, err := os.ReadFile(outsidePath)
				gt.AssertNil(err)
				gt.Assert(after, before)
			})
		})
	}
	gtest.C(t, func(gt *gtest.T) {
		path := t.TempDir()
		assertStorageFilePathRejected(t, gsession.NewStorageFile(path, time.Minute), filepath.Join(path, "absolute"))
		entries, err := os.ReadDir(path)
		gt.AssertNil(err)
		gt.Assert(len(entries), 0)
	})
}

// Test_StorageFilePathCustomIds keeps the existing suffix and safe nested directory behavior.
func Test_StorageFilePathCustomIds(t *testing.T) {
	ids := []string{
		gsession.NewSessionId(), "custom-id", "", ".", "..", ".." + string(os.PathSeparator),
		strings.Repeat("a", 247), filepath.Join("nested", "custom"),
		"nested" + string(os.PathSeparator), strings.Join([]string{"nested", "..", "custom"}, string(os.PathSeparator)),
	}
	for i, id := range ids {
		t.Run(fmt.Sprintf("id-%d", i), func(t *testing.T) {
			gtest.C(t, func(gt *gtest.T) {
				var (
					ctx     = context.Background()
					path    = t.TempDir()
					storage = gsession.NewStorageFile(path, time.Minute)
				)
				gt.AssertNil(os.Mkdir(filepath.Join(path, "nested"), 0o700))
				data, err := storage.GetSession(ctx, id, time.Minute)
				gt.AssertNil(err)
				gt.AssertNil(data)
				gt.AssertNil(storage.UpdateTTL(ctx, id, time.Second))
				gt.AssertNil(storage.SetSession(ctx, id, gmap.NewStrAnyMapFrom(map[string]any{"value": "original"}, true), time.Minute))
				_, err = os.Stat(gfile.Join(path, id) + ".session")
				gt.AssertNil(err)
				data, err = storage.GetSession(ctx, id, time.Minute)
				gt.AssertNil(err)
				gt.Assert(data.Get("value"), "original")
				gt.AssertNil(storage.RemoveAll(ctx, id))
				data, err = storage.GetSession(ctx, id, time.Minute)
				gt.AssertNil(err)
				gt.AssertNil(data)
			})
		})
	}
	gtest.C(t, func(gt *gtest.T) {
		storage := gsession.NewStorageFile(t.TempDir(), time.Minute)
		data, err := storage.GetSession(context.Background(), filepath.Join("missing", "parent", "file"), time.Minute)
		gt.AssertNil(err)
		gt.AssertNil(data)
	})
}
