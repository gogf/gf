// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"os"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_Copy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(gfile.Copy(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(gfile.IsFile(testpath()+topath), true)
		t.AssertNE(gfile.Copy(paths, ""), nil)
		t.AssertNE(gfile.Copy("", topath), nil)
	})
}

func Test_Copy_File_To_Dir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			src = gtest.DataPath("dir1", "file1")
			dst = gfile.Temp(guid.S(), "dir2")
		)
		err := gfile.Mkdir(dst)
		t.AssertNil(err)
		defer gfile.Remove(dst)

		err = gfile.Copy(src, dst)
		t.AssertNil(err)

		expectPath := gfile.Join(dst, "file1")
		t.Assert(gfile.GetContents(expectPath), gfile.GetContents(src))
	})
}

func Test_Copy_Dir_To_File(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			src = gtest.DataPath("dir1")
			dst = gfile.Temp(guid.S(), "file2")
		)
		f, err := gfile.Create(dst)
		t.AssertNil(err)
		defer f.Close()
		defer gfile.Remove(dst)

		err = gfile.Copy(src, dst)
		t.AssertNE(err, nil)
	})
}

func Test_CopyFile(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			paths  = "/testfile_copyfile1.txt"
			topath = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		t.Assert(gfile.CopyFile(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		t.Assert(gfile.IsFile(testpath()+topath), true)
		t.AssertNE(gfile.CopyFile(paths, ""), nil)
		t.AssertNE(gfile.CopyFile("", topath), nil)
	})
	// Content replacement.
	gtest.C(t, func(t *gtest.T) {
		src := gfile.Temp(gtime.TimestampNanoStr())
		dst := gfile.Temp(gtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		t.Assert(gfile.PutContents(src, srcContent), nil)
		t.Assert(gfile.PutContents(dst, dstContent), nil)
		t.Assert(gfile.GetContents(src), srcContent)
		t.Assert(gfile.GetContents(dst), dstContent)

		t.Assert(gfile.CopyFile(src, dst), nil)
		t.Assert(gfile.GetContents(src), srcContent)
		t.Assert(gfile.GetContents(dst), srcContent)
	})
	// Set mode
	gtest.C(t, func(t *gtest.T) {
		var (
			src     = "/testfile_copyfile1.txt"
			dst     = "/testfile_copyfile2.txt"
			dstMode = os.FileMode(0600)
		)
		t.AssertNil(createTestFile(src, ""))
		defer delTestFiles(src)

		t.Assert(gfile.CopyFile(testpath()+src, testpath()+dst, gfile.CopyOption{Mode: dstMode}), nil)
		defer delTestFiles(dst)

		dstStat, err := gfile.Stat(testpath() + dst)
		t.AssertNil(err)
		t.Assert(dstStat.Mode().Perm(), dstMode)
	})
	// Preserve src file's mode
	gtest.C(t, func(t *gtest.T) {
		var (
			src = "/testfile_copyfile1.txt"
			dst = "/testfile_copyfile2.txt"
		)
		t.AssertNil(createTestFile(src, ""))
		defer delTestFiles(src)

		t.Assert(gfile.CopyFile(testpath()+src, testpath()+dst, gfile.CopyOption{PreserveMode: true}), nil)
		defer delTestFiles(dst)

		srcStat, err := gfile.Stat(testpath() + src)
		t.AssertNil(err)
		dstStat, err := gfile.Stat(testpath() + dst)
		t.AssertNil(err)
		t.Assert(srcStat.Mode().Perm(), dstStat.Mode().Perm())
	})
}

func Test_CopyDir(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			dirPath1 = "/test-copy-dir1"
			dirPath2 = "/test-copy-dir2"
		)
		haveList := []string{
			"t1.txt",
			"t2.txt",
		}
		createDir(dirPath1)
		for _, v := range haveList {
			t.Assert(createTestFile(dirPath1+"/"+v, ""), nil)
		}
		defer delTestFiles(dirPath1)

		var (
			yfolder  = testpath() + dirPath1
			tofolder = testpath() + dirPath2
		)

		if gfile.IsDir(tofolder) {
			t.Assert(gfile.Remove(tofolder), nil)
			t.Assert(gfile.Remove(""), nil)
		}

		t.Assert(gfile.CopyDir(yfolder, tofolder), nil)
		defer delTestFiles(tofolder)

		t.Assert(gfile.IsDir(yfolder), true)

		for _, v := range haveList {
			t.Assert(gfile.IsFile(yfolder+"/"+v), true)
		}

		t.Assert(gfile.IsDir(tofolder), true)

		for _, v := range haveList {
			t.Assert(gfile.IsFile(tofolder+"/"+v), true)
		}

		t.Assert(gfile.Remove(tofolder), nil)
		t.Assert(gfile.Remove(""), nil)
	})
	// Content replacement.
	gtest.C(t, func(t *gtest.T) {
		src := gfile.Temp(gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		dst := gfile.Temp(gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		defer func() {
			gfile.Remove(src)
			gfile.Remove(dst)
		}()
		srcContent := "1"
		dstContent := "1"
		t.Assert(gfile.PutContents(src, srcContent), nil)
		t.Assert(gfile.PutContents(dst, dstContent), nil)
		t.Assert(gfile.GetContents(src), srcContent)
		t.Assert(gfile.GetContents(dst), dstContent)

		err := gfile.CopyDir(gfile.Dir(src), gfile.Dir(dst))
		t.AssertNil(err)
		t.Assert(gfile.GetContents(src), srcContent)
		t.Assert(gfile.GetContents(dst), srcContent)

		t.AssertNE(gfile.CopyDir(gfile.Dir(src), ""), nil)
		t.AssertNE(gfile.CopyDir("", gfile.Dir(dst)), nil)
	})
}
