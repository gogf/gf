// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres_test

import (
	_ "github.com/gogf/gf/os/gres/testdata/data"

	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"

	"github.com/gogf/gf/os/gres"
)

func Test_Basic(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gres.Get("none"), nil)
		t.Assert(gres.Contains("none"), false)
		t.Assert(gres.Contains("dir1"), true)
	})

	gtest.C(t, func(t *gtest.T) {
		path := "dir1/test1"
		file := gres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)

		info := file.FileInfo()
		t.AssertNE(info, nil)
		t.Assert(info.IsDir(), false)
		t.Assert(info.Name(), "test1")

		rc, err := file.Open()
		t.Assert(err, nil)
		defer rc.Close()

		b := make([]byte, 5)
		n, err := rc.Read(b)
		t.Assert(n, 5)
		t.Assert(err, nil)
		t.Assert(string(b), "test1")

		t.Assert(file.Content(), "test1 content")
	})

	gtest.C(t, func(t *gtest.T) {
		path := "dir2"
		file := gres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)

		info := file.FileInfo()
		t.AssertNE(info, nil)
		t.Assert(info.IsDir(), true)
		t.Assert(info.Name(), "dir2")

		rc, err := file.Open()
		t.Assert(err, nil)
		defer rc.Close()

		t.Assert(file.Content(), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		path := "dir2/test2"
		file := gres.Get(path)
		t.AssertNE(file, nil)
		t.Assert(file.Name(), path)
		t.Assert(file.Content(), "test2 content")
	})
}

func Test_Get(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(gres.Get("dir1/test1"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		file := gres.GetWithIndex("dir1", g.SliceStr{"test1"})
		t.AssertNE(file, nil)
		t.Assert(file.Name(), "dir1/test1")
	})
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gres.GetContent("dir1"), "")
		t.Assert(gres.GetContent("dir1/test1"), "test1 content")
	})
}

func Test_ScanDir(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		path := "dir1"
		files := gres.ScanDir(path, "*", false)
		t.AssertNE(files, nil)
		t.Assert(len(files), 2)
	})
	gtest.C(t, func(t *gtest.T) {
		path := "dir1"
		files := gres.ScanDir(path, "*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 3)
	})

	gtest.C(t, func(t *gtest.T) {
		path := "dir1"
		files := gres.ScanDir(path, "*.*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
		t.Assert(files[0].Name(), "dir1/sub/sub-test1.txt")
		t.Assert(files[0].Content(), "sub-test1 content")
	})
}

func Test_ScanDirFile(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		path := "dir2"
		files := gres.ScanDirFile(path, "*", false)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
	})
	gtest.C(t, func(t *gtest.T) {
		path := "dir2"
		files := gres.ScanDirFile(path, "*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 2)
	})

	gtest.C(t, func(t *gtest.T) {
		path := "dir2"
		files := gres.ScanDirFile(path, "*.*", true)
		t.AssertNE(files, nil)
		t.Assert(len(files), 1)
		t.Assert(files[0].Name(), "dir2/sub/sub-test2.txt")
		t.Assert(files[0].Content(), "sub-test2 content")
	})
}
