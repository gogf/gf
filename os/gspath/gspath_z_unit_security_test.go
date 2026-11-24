// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gspath_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gspath"
	"github.com/gogf/gf/v2/test/gtest"
)

func TestSPath_Search_Escape(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			base    = t.TempDir()
			root    = filepath.Join(base, "root")
			sibling = filepath.Join(base, "root-evil")
		)
		t.AssertNil(gfile.Mkdir(root))
		t.AssertNil(gfile.Mkdir(sibling))
		t.AssertNil(os.WriteFile(filepath.Join(root, "inside.txt"), []byte("inside"), 0o644))
		t.AssertNil(os.WriteFile(filepath.Join(sibling, "secret.txt"), []byte("secret"), 0o644))

		// The sibling directory shares the string prefix of the root.
		name := gfile.Join("..", "root-evil", "secret.txt")
		path, _ := gspath.Search(root, name)
		t.Assert(path, "")
		path, _ = gspath.SearchWithCache(root, name)
		t.Assert(path, "")

		// Normal files are still found in both modes.
		path, isDir := gspath.Search(root, "inside.txt")
		t.Assert(path, filepath.Join(root, "inside.txt"))
		t.Assert(isDir, false)
		path, isDir = gspath.SearchWithCache(root, "inside.txt")
		t.Assert(path, filepath.Join(root, "inside.txt"))
		t.Assert(isDir, false)
	})
}
