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
	gtest.Case(t, func() {
		srcPath := gfile.Join(gdebug.TestDataPath(), "zip", "path1", "1.txt")
		dstPath := gfile.Join(gdebug.TestDataPath(), "zip", "zip.zip")

		gtest.Assert(gfile.Exists(dstPath), false)
		err := gcompress.ZipPath(srcPath, dstPath)
		gtest.Assert(err, nil)
		gtest.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		gtest.Assert(err, nil)

		err = gcompress.UnZipFile(dstPath, tempDirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(tempDirPath)

		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
	})
	// directory
	gtest.Case(t, func() {
		srcPath := gfile.Join(gdebug.TestDataPath(), "zip")
		dstPath := gfile.Join(gdebug.TestDataPath(), "zip", "zip.zip")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		gtest.Assert(err, nil)

		gtest.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath, dstPath)
		gtest.Assert(err, nil)
		gtest.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		gtest.Assert(err, nil)

		err = gcompress.UnZipFile(dstPath, tempDirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(tempDirPath)

		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "zip", "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "zip", "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
	// multiple paths joined using char ','
	gtest.Case(t, func() {
		srcPath := gfile.Join(gdebug.TestDataPath(), "zip")
		srcPath1 := gfile.Join(gdebug.TestDataPath(), "zip", "path1")
		srcPath2 := gfile.Join(gdebug.TestDataPath(), "zip", "path2")
		dstPath := gfile.Join(gdebug.TestDataPath(), "zip", "zip.zip")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		gtest.Assert(err, nil)

		gtest.Assert(gfile.Exists(dstPath), false)
		err = gcompress.ZipPath(srcPath1+", "+srcPath2, dstPath)
		gtest.Assert(err, nil)
		gtest.Assert(gfile.Exists(dstPath), true)
		defer gfile.Remove(dstPath)

		tempDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		gtest.Assert(err, nil)

		zipContent := gfile.GetBytes(dstPath)
		gtest.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(tempDirPath)

		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}

func Test_ZipPathWriter(t *testing.T) {
	gtest.Case(t, func() {
		srcPath := gfile.Join(gdebug.TestDataPath(), "zip")
		srcPath1 := gfile.Join(gdebug.TestDataPath(), "zip", "path1")
		srcPath2 := gfile.Join(gdebug.TestDataPath(), "zip", "path2")

		pwd := gfile.Pwd()
		err := gfile.Chdir(srcPath)
		defer gfile.Chdir(pwd)
		gtest.Assert(err, nil)

		writer := bytes.NewBuffer(nil)
		gtest.Assert(writer.Len(), 0)
		err = gcompress.ZipPathWriter(srcPath1+", "+srcPath2, writer)
		gtest.Assert(err, nil)
		gtest.AssertGT(writer.Len(), 0)

		tempDirPath := gfile.Join(gfile.TempDir(), gtime.TimestampNanoStr())
		err = gfile.Mkdir(tempDirPath)
		gtest.Assert(err, nil)

		zipContent := writer.Bytes()
		gtest.AssertGT(len(zipContent), 0)
		err = gcompress.UnZipContent(zipContent, tempDirPath)
		gtest.Assert(err, nil)
		defer gfile.Remove(tempDirPath)

		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path1", "1.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path1", "1.txt")),
		)
		gtest.Assert(
			gfile.GetContents(gfile.Join(tempDirPath, "path2", "2.txt")),
			gfile.GetContents(gfile.Join(srcPath, "path2", "2.txt")),
		)
	})
}
