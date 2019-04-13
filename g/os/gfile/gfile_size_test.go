//test 100%

package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestSize(t *testing.T) {
	gtest.Case(t, func() {
		var (
			paths1 string = "/testfile_t1.txt"
			sizes  int64
		)

		CreateTestFile(paths1, "abcdefghijklmn")
		defer DelTestFiles(paths1)

		sizes = Size(Testpath() + paths1)
		gtest.Assert(sizes, 14)

		sizes = Size("")
		gtest.Assert(sizes, 0)

	})
}

func TestFormatSize(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(FormatSize(0), "0.00B")
		gtest.Assert(FormatSize(16), "16.00B")

		gtest.Assert(FormatSize(1024), "1.00K")

		gtest.Assert(FormatSize(16000000), "15.26M")

		gtest.Assert(FormatSize(1600000000), "1.49G")

		gtest.Assert(FormatSize(9600000000000), "8.73T")
		gtest.Assert(FormatSize(9600000000000000), "8.53P")

		gtest.Assert(FormatSize(9600000000000000000), "TooLarge")

	})
}

func TestReadableSize(t *testing.T) {
	gtest.Case(t, func() {

		var (
			paths1 string = "/testfile_t1.txt"
		)
		CreateTestFile(paths1, "abcdefghijklmn")
		defer DelTestFiles(paths1)
		gtest.Assert(ReadableSize(Testpath()+paths1), "14.00B")
		gtest.Assert(ReadableSize(""), "0.00B")

	})
}
