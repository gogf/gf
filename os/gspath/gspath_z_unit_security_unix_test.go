//go:build unix

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gspath_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSPath_Search_SymlinkEscape(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			base    = t.TempDir()
			root    = filepath.Join(base, "root")
			outside = filepath.Join(base, "outside")
		)
		t.AssertNil(os.Mkdir(root, 0o755))
		t.AssertNil(os.Mkdir(outside, 0o755))
		t.AssertNil(os.WriteFile(filepath.Join(root, "inside.txt"), []byte("inside"), 0o644))
		t.AssertNil(os.WriteFile(filepath.Join(outside, "secret.txt"), []byte("secret"), 0o644))
		t.AssertNil(os.Symlink(filepath.Join(root, "inside.txt"), filepath.Join(root, "inside-link.txt")))
		t.AssertNil(os.Symlink(filepath.Join(outside, "secret.txt"), filepath.Join(root, "escape-file.txt")))
		t.AssertNil(os.Symlink(outside, filepath.Join(root, "escape-dir")))

		// A symlink resolving inside the root is still searchable.
		path, _ := gspath.Search(root, "inside-link.txt")
		t.Assert(path, filepath.Join(root, "inside-link.txt"))

		// Symlinks resolving outside the root are not returned.
		path, _ = gspath.Search(root, "escape-file.txt")
		t.Assert(path, "")
		path, _ = gspath.Search(root, "escape-dir/secret.txt")
		t.Assert(path, "")

		// Cache entries are validated by their resolved paths as well.
		path, _ = gspath.SearchWithCache(root, "inside-link.txt")
		t.Assert(path, filepath.Join(root, "inside-link.txt"))
		path, _ = gspath.SearchWithCache(root, "escape-file.txt")
		t.Assert(path, "")
		path, _ = gspath.SearchWithCache(root, "escape-dir/secret.txt")
		t.Assert(path, "")
	})
}

func TestSPath_RootSlash(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sp := gspath.New("/", false)
		t.Assert(sp.Paths(), []string{"/"})

		filePath := filepath.Join(t.TempDir(), "gf.txt")
		t.AssertNil(os.WriteFile(filePath, []byte("test"), 0o644))
		path, isDir := gspath.Search("/", strings.TrimPrefix(filePath, "/"))
		t.Assert(path, filePath)
		t.Assert(isDir, false)
	})
}
