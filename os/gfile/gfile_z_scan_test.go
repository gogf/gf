// Copyright 2017-2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"github.com/gogf/gf/debug/gdebug"
	"testing"

	"github.com/gogf/gf/os/gfile"

	"github.com/gogf/gf/test/gtest"
)

func Test_ScanDir(t *testing.T) {
	teatPath := gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata"
	gtest.Case(t, func() {
		files, err := gfile.ScanDir(teatPath, "*", false)
		gtest.Assert(err, nil)
		gtest.AssertIN(teatPath+gfile.Separator+"dir1", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir2", files)
		gtest.AssertNE(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
	})
	gtest.Case(t, func() {
		files, err := gfile.ScanDir(teatPath, "*", true)
		gtest.Assert(err, nil)
		gtest.AssertIN(teatPath+gfile.Separator+"dir1", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir2", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir2"+gfile.Separator+"file2", files)
	})
}

func Test_ScanDirFile(t *testing.T) {
	teatPath := gfile.Dir(gdebug.CallerFilePath()) + gfile.Separator + "testdata"
	gtest.Case(t, func() {
		files, err := gfile.ScanDirFile(teatPath, "*", false)
		gtest.Assert(err, nil)
		gtest.Assert(len(files), 0)
	})
	gtest.Case(t, func() {
		files, err := gfile.ScanDirFile(teatPath, "*", true)
		gtest.Assert(err, nil)
		gtest.AssertNI(teatPath+gfile.Separator+"dir1", files)
		gtest.AssertNI(teatPath+gfile.Separator+"dir2", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir2"+gfile.Separator+"file2", files)
	})
}
