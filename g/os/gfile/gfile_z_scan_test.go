package gfile_test

import (
	"testing"

	"github.com/gogf/gf/g/os/gfile"

	"github.com/gogf/gf/g/test/gtest"
)

func Test_Scan(t *testing.T) {
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
		gtest.AssertNE(teatPath+gfile.Separator+"dir1", files)
		gtest.AssertNE(teatPath+gfile.Separator+"dir2", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir1"+gfile.Separator+"file1", files)
		gtest.AssertIN(teatPath+gfile.Separator+"dir2"+gfile.Separator+"file2", files)
	})
}
