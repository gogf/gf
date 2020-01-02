// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_Copy(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths  string = "/testfile_copyfile1.txt"
			topath string = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		gtest.Assert(gfile.Copy(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		gtest.Assert(gfile.IsFile(testpath()+topath), true)
		gtest.AssertNE(gfile.Copy("", ""), nil)
	})
}

func Test_CopyFile(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths  string = "/testfile_copyfile1.txt"
			topath string = "/testfile_copyfile2.txt"
		)

		createTestFile(paths, "")
		defer delTestFiles(paths)

		gtest.Assert(gfile.CopyFile(testpath()+paths, testpath()+topath), nil)
		defer delTestFiles(topath)

		gtest.Assert(gfile.IsFile(testpath()+topath), true)
		gtest.AssertNE(gfile.CopyFile("", ""), nil)
	})
	// Content replacement.
	gtest.Case(t, func() {
		src := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		dst := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		gfile.PutContents(src, srcContent)
		gfile.PutContents(dst, dstContent)
		gtest.Assert(gfile.GetContents(src), srcContent)
		gtest.Assert(gfile.GetContents(dst), dstContent)

		err := gfile.CopyFile(src, dst)
		gtest.Assert(err, nil)
		gtest.Assert(gfile.GetContents(src), srcContent)
		gtest.Assert(gfile.GetContents(dst), srcContent)
	})
}

func Test_CopyDir(t *testing.T) {
	gtest.Case(t, func() {
		var (
			dirpath1 string = "/testcopydir1"
			dirpath2 string = "/testcopydir2"
		)

		havelist1 := []string{
			"t1.txt",
			"t2.txt",
		}

		createDir(dirpath1)
		for _, v := range havelist1 {
			createTestFile(dirpath1+"/"+v, "")
		}
		defer delTestFiles(dirpath1)

		yfolder := testpath() + dirpath1
		tofolder := testpath() + dirpath2

		if gfile.IsDir(tofolder) {
			gtest.Assert(gfile.Remove(tofolder), nil)
			gtest.Assert(gfile.Remove(""), nil)
		}

		gtest.Assert(gfile.CopyDir(yfolder, tofolder), nil)
		defer delTestFiles(tofolder)

		// 检查复制后的旧文件夹是否真实存在
		gtest.Assert(gfile.IsDir(yfolder), true)

		// 检查复制后的旧文件夹中的文件是否真实存在
		for _, v := range havelist1 {
			gtest.Assert(gfile.IsFile(yfolder+"/"+v), true)
		}

		// 检查复制后的新文件夹是否真实存在
		gtest.Assert(gfile.IsDir(tofolder), true)

		// 检查复制后的新文件夹中的文件是否真实存在
		for _, v := range havelist1 {
			gtest.Assert(gfile.IsFile(tofolder+"/"+v), true)
		}

		gtest.Assert(gfile.Remove(tofolder), nil)
		gtest.Assert(gfile.Remove(""), nil)
	})
	// Content replacement.
	gtest.Case(t, func() {
		src := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		dst := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr(), gtime.TimestampNanoStr())
		srcContent := "1"
		dstContent := "1"
		gfile.PutContents(src, srcContent)
		gfile.PutContents(dst, dstContent)
		gtest.Assert(gfile.GetContents(src), srcContent)
		gtest.Assert(gfile.GetContents(dst), dstContent)

		err := gfile.CopyDir(gfile.Dir(src), gfile.Dir(dst))
		gtest.Assert(err, nil)
		gtest.Assert(gfile.GetContents(src), srcContent)
		gtest.Assert(gfile.GetContents(dst), srcContent)
	})
}
