// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"io/fs"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

type testStaticFileInfo struct {
	name  string
	isDir bool
}

func (f testStaticFileInfo) Name() string       { return f.name }
func (f testStaticFileInfo) Size() int64        { return 0 }
func (f testStaticFileInfo) Mode() fs.FileMode  { return 0 }
func (f testStaticFileInfo) ModTime() time.Time { return time.Time{} }
func (f testStaticFileInfo) IsDir() bool        { return f.isDir }
func (f testStaticFileInfo) Sys() any           { return nil }

func Test_CompareStaticFiles(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		dirA := testStaticFileInfo{name: "a", isDir: true}
		dirB := testStaticFileInfo{name: "b", isDir: true}
		fileA := testStaticFileInfo{name: "a"}
		fileB := testStaticFileInfo{name: "b"}

		t.Assert(compareStaticFiles(dirA, fileA) < 0, true)
		t.Assert(compareStaticFiles(fileA, dirA) > 0, true)
		t.Assert(compareStaticFiles(dirA, dirB) < 0, true)
		t.Assert(compareStaticFiles(fileB, fileA) > 0, true)
		t.Assert(compareStaticFiles(fileA, fileA), 0)
	})
}
