package gfile_test

import (
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestSize(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_t1.txt"
			sizes  int64
		)

		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)

		sizes = gfile.Size(testpath() + paths1)
		gtest.Assert(sizes, 14)

		sizes = gfile.Size("")
		gtest.Assert(sizes, 0)

	})
}

func TestFormatSize(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(gfile.FormatSize(0), "0.00B")
		gtest.Assert(gfile.FormatSize(16), "16.00B")

		gtest.Assert(gfile.FormatSize(1024), "1.00K")

		gtest.Assert(gfile.FormatSize(16000000), "15.26M")

		gtest.Assert(gfile.FormatSize(1600000000), "1.49G")

		gtest.Assert(gfile.FormatSize(9600000000000), "8.73T")
		gtest.Assert(gfile.FormatSize(9600000000000000), "8.53P")

		gtest.Assert(gfile.FormatSize(9600000000000000000), "TooLarge")

	})
}

func TestReadableSize(t *testing.T) {
	gtest.Case(t, func() {

		var (
			paths1 string = "/testfile_t1.txt"
		)
		createTestFile(paths1, "abcdefghijklmn")
		defer delTestFiles(paths1)
		gtest.Assert(gfile.ReadableSize(testpath()+paths1), "14.00B")
		gtest.Assert(gfile.ReadableSize(""), "0.00B")

	})
}
