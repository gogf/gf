// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres_test

import (
	"testing"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/test/gtest"

	"github.com/gogf/gf/os/gres"
	_ "github.com/gogf/gf/os/gres/testdata/data"
)

func Test_Basic(t *testing.T) {
	gres.Dump()
	gtest.Case(t, func() {
		gtest.Assert(gres.Get("none"), nil)
		gtest.Assert(gres.Contains("none"), false)
		gtest.Assert(gres.Contains("dir1"), true)
	})

	gtest.Case(t, func() {
		path := "dir1/test1"
		file := gres.Get(path)
		gtest.AssertNE(file, nil)
		gtest.Assert(file.Name(), path)

		info := file.FileInfo()
		gtest.AssertNE(info, nil)
		gtest.Assert(info.IsDir(), false)
		gtest.Assert(info.Name(), "test1")

		rc, err := file.Open()
		gtest.Assert(err, nil)
		defer rc.Close()

		b := make([]byte, 5)
		n, err := rc.Read(b)
		gtest.Assert(n, 5)
		gtest.Assert(err, nil)
		gtest.Assert(string(b), "test1")

		gtest.Assert(file.Content(), "test1 content")
	})

	gtest.Case(t, func() {
		path := "dir2"
		file := gres.Get(path)
		gtest.AssertNE(file, nil)
		gtest.Assert(file.Name(), path)

		info := file.FileInfo()
		gtest.AssertNE(info, nil)
		gtest.Assert(info.IsDir(), true)
		gtest.Assert(info.Name(), "dir2")

		rc, err := file.Open()
		gtest.Assert(err, nil)
		defer rc.Close()

		gtest.Assert(file.Content(), nil)
	})

	gtest.Case(t, func() {
		path := "dir2/test2"
		file := gres.Get(path)
		gtest.AssertNE(file, nil)
		gtest.Assert(file.Name(), path)
		gtest.Assert(file.Content(), "test2 content")
	})
}

func Test_Get(t *testing.T) {
	gres.Dump()
	gtest.Case(t, func() {
		gtest.AssertNE(gres.Get("dir1/test1"), nil)
	})
	gtest.Case(t, func() {
		file := gres.GetWithIndex("dir1", g.SliceStr{"test1"})
		gtest.AssertNE(file, nil)
		gtest.Assert(file.Name(), "dir1/test1")
	})
	gtest.Case(t, func() {
		gtest.Assert(gres.GetContent("dir1"), "")
		gtest.Assert(gres.GetContent("dir1/test1"), "test1 content")
	})
}

func Test_ScanDir(t *testing.T) {
	gres.Dump()
	gtest.Case(t, func() {
		path := "dir1"
		files := gres.ScanDir(path, "*", false)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 2)
	})
	gtest.Case(t, func() {
		path := "dir1"
		files := gres.ScanDir(path, "*", true)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 3)
	})

	gtest.Case(t, func() {
		path := "dir1"
		files := gres.ScanDir(path, "*.*", true)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 1)
		gtest.Assert(files[0].Name(), "dir1/sub/sub-test1.txt")
		gtest.Assert(files[0].Content(), "sub-test1 content")
	})
}

func Test_ScanDirFile(t *testing.T) {
	gres.Dump()
	gtest.Case(t, func() {
		path := "dir2"
		files := gres.ScanDirFile(path, "*", false)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 1)
	})
	gtest.Case(t, func() {
		path := "dir2"
		files := gres.ScanDirFile(path, "*", true)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 2)
	})

	gtest.Case(t, func() {
		path := "dir2"
		files := gres.ScanDirFile(path, "*.*", true)
		gtest.AssertNE(files, nil)
		gtest.Assert(len(files), 1)
		gtest.Assert(files[0].Name(), "dir2/sub/sub-test2.txt")
		gtest.Assert(files[0].Content(), "sub-test2 content")
	})
}
