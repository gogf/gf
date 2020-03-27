// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcompress_test

import (
	"bytes"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/encoding/gcompress"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_ZipPath(t *testing.T) {
	// file
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("zip", "path1", "1.txt")
		dstPath := gdebug.TestDataPath("zip", "zip.zip")

		t.Assert(gfile.Exists(dstPath), false)
		err := gcompress.ZipPath(srcPath, dstPath)
		t.Assert(err, nil)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		err = gcompress.UnZipFile(dstPath, tempDirPath)
		t.Assert(err, nil)
		defer gfile.Remove(tempDirPath)

		t.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
	})
	// directory
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("zip")
		dstPath := gdebug.TestDataPath("zip", "zip.zip")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.Assert(err, nil)

		t.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath, dstPath)
		t.Assert(err, nil)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		err = gcompress.UnZipFile(dstPath, tempDirPath)
		t.Assert(err, nil)
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
	// multiple paths joined using char ','
	gtest.C(t, func(t *gtest.T) {
		srcPath := gdebug.TestDataPath("zip")
		srcPath1 := gdebug.TestDataPath("zip", "path1")
		srcPath2 := gdebug.TestDataPath("zip", "path2")
		dstPath := gdebug.TestDataPath("zip", "zip.zip")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.Assert(err, nil)

		t.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath1+", "+srcPath2, dstPath)
		t.Assert(err, nil)
		t.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		zipContent := gfile.GetBytes(dstPath)
		t.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		t.Assert(err, nil)
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
		srcPath := gdebug.TestDataPath("zip")
		srcPath1 := gdebug.TestDataPath("zip", "path1")
		srcPath2 := gdebug.TestDataPath("zip", "path2")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		t.Assert(err, nil)

		writer := bytes.NewBuffer(nil)
		t.Assert(writer.Len(), 0)
		err = gcompress.ZipPathWriter(srcPath1+", "+srcPath2, writer)
		t.Assert(err, nil)
		t.AssertGT(writer.Len(), 0)

		tempDirPath := gfile.TempDir(gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		t.Assert(err, nil)

		zipContent := writer.Bytes()
		t.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		t.Assert(err, nil)
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
