//go:build windows

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

func TestSPath_DriveRoot(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		root := filepath.VolumeName(gfile.Pwd()) + `\`
		sp := gspath.New(root, false)
		t.Assert(sp.Paths(), []string{root})
		filePath := filepath.Join(t.TempDir(), "gf.txt")
		t.AssertNil(os.WriteFile(filePath, []byte("test"), 0o644))
		rel, err := filepath.Rel(root, filePath)
		t.AssertNil(err)
		path, isDir := gspath.Search(root, rel)
		t.Assert(path, filePath)
		t.Assert(isDir, false)
	})
}
