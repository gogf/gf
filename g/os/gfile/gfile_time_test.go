package gfile_test

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/test/gtest"
	"os"
	"testing"
)

func TestMTime(t *testing.T) {
	gtest.Case(t, func() {

		var (
			file1   string = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		gtest.Assert(err, nil)

		gtest.Assert(gfile.MTime(testpath()+file1), fileobj.ModTime().Unix())
		gtest.Assert(gfile.MTime(""), 0)
	})
}

func TestMTimeMillisecond(t *testing.T) {
	gtest.Case(t, func() {
		var (
			file1   string = "/testfile_t1.txt"
			err     error
			fileobj os.FileInfo
		)

		createTestFile(file1, "")
		defer delTestFiles(file1)
		fileobj, err = os.Stat(testpath() + file1)
		gtest.Assert(err, nil)

		gtest.AssertGTE(gfile.MTimeMillisecond(testpath()+file1), fileobj.ModTime().Nanosecond()/1000000)
		gtest.Assert(gfile.MTimeMillisecond(""), 0)
	})
}
