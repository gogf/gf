// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"bytes"
	"testing"

	"github.com/gogf/gf/v2/encoding/gcompress"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ZipPath(t *testing.T) {
	// file
	gtest.C(t, func(t *gtest.T) {
		srcPath := gtest.DataPath("zip", "path1", "1.txt")
		dstPath := gtest.DataPath("zip", "zip.zip")

		t.Assert(gfile.Exists(dstPath), false)
		t.Assert(gcompress.ZipPath(srcPath, dstPath), nil)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		// unzip to temporary dir.
		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		t.Assert(gfile.Mkdir(tempDirPath), nil)
		t.Assert(gcompress.UnZipFile(dstPath, tempDirPath), nil)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "1.txt")),
			gfile.GetContents(srcPath),
		)
	})
	// multiple files
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath1 = gtest.DataPath("zip", "path1", "1.txt")
			srcPath2 = gtest.DataPath("zip", "path2", "2.txt")
			dstPath  = gfile.Temp(gtime.TimestampNanoStr(), "zip.zip")
		)
		if p := gfile.Dir(dstPath); !gfile.Exists(p) {
			t.Assert(gfile.Mkdir(p), nil)
		}

		t.Assert(gfile.Exists(dstPath), false)
		err := gcompress.ZipPath(srcPath1+","+srcPath2, dstPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		// unzip to another temporary dir.
		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		t.Assert(gfile.Mkdir(tempDirPath), nil)
		err = gcompress.UnZipFile(dstPath, tempDirPath)
		t.AssertNil(err)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "1.txt")),
			gfile.GetContents(srcPath1),
		)
		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "2.txt")),
			gfile.GetContents(srcPath2),
		)
	})
	// one dir and one file.
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath1 = gtest.DataPath("zip", "path1")
			srcPath2 = gtest.DataPath("zip", "path2", "2.txt")
			dstPath  = gfile.Temp(gtime.TimestampNanoStr(), "zip.zip")
		)
		if p := gfile.Dir(dstPath); !gfile.Exists(p) {
			t.Assert(gfile.Mkdir(p), nil)
		}

		t.Assert(gfile.Exists(dstPath), false)
		err := gcompress.ZipPath(srcPath1+","+srcPath2, dstPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		// unzip to another temporary dir.
		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		t.Assert(gfile.Mkdir(tempDirPath), nil)
		err = gcompress.UnZipFile(dstPath, tempDirPath)
		t.AssertNil(err)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath1, "1.txt")),
		)
		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "2.txt")),
			gfile.GetContents(srcPath2),
		)
	})
	// directory.
	gtest.C(t, func(t *gtest.T) {
		srcPath := gtest.DataPath("zip")
		dstPath := gtest.DataPath("zip", "zip.zip")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.AssertNil(err)

		t.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath, dstPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.AssertNil(err)

		err = gcompress.UnZipFile(dstPath, tempDirPath)
		t.AssertNil(err)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "zip", "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "zip", "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
	// multiple directory paths joined using char ','.
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath  = gtest.DataPath("zip")
			srcPath1 = gtest.DataPath("zip", "path1")
			srcPath2 = gtest.DataPath("zip", "path2")
			dstPath  = gtest.DataPath("zip", "zip.zip")
		)
		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.AssertNil(err)

		t.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath1+", "+srcPath2, dstPath)
		t.AssertNil(err)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.AssertNil(err)

		zipContent := gfile.GetBytes(dstPath)
		t.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		t.AssertNil(err)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}

func Test_ZipPathWriter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			srcPath  = gtest.DataPath("zip")
			srcPath1 = gtest.DataPath("zip", "path1")
			srcPath2 = gtest.DataPath("zip", "path2")
		)
		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.AssertNil(err)

		writer := bytes.NewBuffer(nil)
		t.Assert(writer.Len(), 0)
		err = gcompress.ZipPathWriter(srcPath1+", "+srcPath2, writer)
		t.AssertNil(err)
		t.AssertGT(writer.Len(), 0)

		tempDirPath := gfile.Temp(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.AssertNil(err)

		zipContent := writer.Bytes()
		t.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		t.AssertNil(err)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}
