package gfile_test

import (
	"testing"

	"github.com/gogf/gf/g/os/gfile"

	"github.com/gogf/gf/g/test/gtest"
)

func Test_ScanDir(t *testing.T) {
	teatPath := gfile.Dir(gfile.SourcePath()) + gfile.Separator + "testdata"
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
	teatPath := gfile.Dir(gfile.SourcePath()) + gfile.Separator + "testdata"
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
