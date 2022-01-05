// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gres_test

import (
	"strings"
	"testing"

	_ "github.com/gogf/gf/v2/os/gres/testdata/data"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gres"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_PackToGoFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath    = gdebug.TestDataPath("files")
			goFilePath = gfile.TempDir(gtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gres.PackToGoFile(srcPath, goFilePath, pkgName)
		)
		t.Assert(err, nil)
		_ = gfile.Remove(goFilePath)
	})
}

func Test_Pack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath   = gdebug.TestDataPath("files")
			data, err = gres.Pack(srcPath)
		)
		t.Assert(err, nil)

		r := gres.New()
		err = r.Add(string(data))
		t.Assert(err, nil)
		t.Assert(r.Contains("files/"), true)
	})
}

func Test_PackToFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath = gdebug.TestDataPath("files")
			dstPath = gfile.TempDir(gtime.TimestampNanoStr())
			err     = gres.PackToFile(srcPath, dstPath)
		)
		t.Assert(err, nil)

		defer gfile.Remove(dstPath)

		r := gres.New()
		err = r.Load(dstPath)
		t.Assert(err, nil)
		t.Assert(r.Contains("files"), true)
	})
}

func Test_PackMulti(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath    = gdebug.TestDataPath("files")
			goFilePath = gfile.TempDir(gtime.TimestampNanoStr(), "data.go")
			pkgName    = "data"
			array, err = gfile.ScanDir(srcPath, "*", false)
		)
		t.Assert(err, nil)
		err = gres.PackToGoFile(strings.Join(array, ","), goFilePath, pkgName)
		t.Assert(err, nil)
		_ = gfile.Remove(goFilePath)
	})
}

func Test_PackWithPrefix1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath    = gdebug.TestDataPath("files")
			goFilePath = gfile.TempDir(gtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gres.PackToGoFile(srcPath, goFilePath, pkgName, "www/gf-site/test")
		)
		t.AssertNil(err)
		_ = gfile.Remove(goFilePath)
	})
}

func Test_PackWithPrefix2(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath    = gdebug.TestDataPath("files")
			goFilePath = gfile.TempDir(gtime.TimestampNanoStr(), "testdata.go")
			pkgName    = "testdata"
			err        = gres.PackToGoFile(srcPath, goFilePath, pkgName, "/var/www/gf-site/test")
		)
		t.Assert(err, nil)
		_ = gfile.Remove(goFilePath)
	})
}

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

func Test_Export(t *testing.T) {
	gres.Dump()
	gtest.C(t, func(t *gtest.T) {
		var (
			src = `template`
			dst = gfile.TempDir(gtime.TimestampNanoStr())
			err = gres.Export(src, dst)
		)
		defer gfile.Remove(dst)
		t.AssertNil(err)
		files, err := gfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 14)

		name := `template/index.html`
		t.Assert(gfile.GetContents(gfile.Join(dst, name)), gres.GetContent(name))
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			src = `template`
			dst = gfile.TempDir(gtime.TimestampNanoStr())
			err = gres.Export(src, dst, gres.ExportOption{
				RemovePrefix: `template`,
			})
		)
		defer gfile.Remove(dst)
		t.AssertNil(err)
		files, err := gfile.ScanDir(dst, "*", true)
		t.AssertNil(err)
		t.Assert(len(files), 13)

		nameInRes := `template/index.html`
		nameInSys := `index.html`
		t.Assert(gfile.GetContents(gfile.Join(dst, nameInSys)), gres.GetContent(nameInRes))
	})
}
