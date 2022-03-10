// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gfile_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ScanDir(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := gfile.ScanDir(teatPath, "*", false)
		t.AssertNil(err)
		t.AssertIN(teatPath+gfile.Separator+"dir1", files)
		t.AssertIN(teatPath+gfile.Separator+"dir2", files)
		t.AssertNE(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
	})
	gtest.C(t, func(t *gtest.T) {
		files, err := gfile.ScanDir(teatPath, "*", true)
		t.AssertNil(err)
		t.AssertIN(teatPath+gfile.Separator+"dir1", files)
		t.AssertIN(teatPath+gfile.Separator+"dir2", files)
		t.AssertIN(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
		t.AssertIN(teatPath+gfile.Separator+"dir2"+gfile.Separator+"file2", files)
	})
}

func Test_ScanDirFunc(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := gfile.ScanDirFunc(teatPath, "*", true, func(path string) string {
			if gfile.Name(path) != "file1" {
				return ""
			}
			return path
		})
		t.AssertNil(err)
		t.Assert(len(files), 1)
		t.Assert(gfile.Name(files[0]), "file1")
	})
}

func Test_ScanDirFile(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		files, err := gfile.ScanDirFile(teatPath, "*", false)
		t.AssertNil(err)
		t.Assert(len(files), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		files, err := gfile.ScanDirFile(teatPath, "*", true)
		t.AssertNil(err)
		t.AssertNI(teatPath+gfile.Separator+"dir1", files)
		t.AssertNI(teatPath+gfile.Separator+"dir2", files)
		t.AssertIN(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
		t.AssertIN(teatPath+gfile.Separator+"dir2"+gfile.Separator+"file2", files)
	})
}

func Test_ScanDirFileFunc(t *testing.T) {
	teatPath := gdebug.TestDataPath()
	gtest.C(t, func(t *gtest.T) {
		array := garray.New()
		files, err := gfile.ScanDirFileFunc(teatPath, "*", false, func(path string) string {
			array.Append(1)
			return path
		})
		t.AssertNil(err)
		t.Assert(len(files), 0)
		t.Assert(array.Len(), 0)
	})
	gtest.C(t, func(t *gtest.T) {
		array := garray.New()
		files, err := gfile.ScanDirFileFunc(teatPath, "*", true, func(path string) string {
			array.Append(1)
			if gfile.Basename(path) == "file1" {
				return path
			}
			return ""
		})
		t.AssertNil(err)
		t.Assert(len(files), 1)
		t.Assert(array.Len(), 3)
	})
}
