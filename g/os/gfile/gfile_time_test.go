//have test 100%
package gfile

import (
	"testing"
	"github.com/gogf/gf/g/test/gtest"

)

func TestMTime(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(MTime("./testfile/dirfiles/t1.txt"),1554883732)
		gtest.Assert(MTime(""),0)
	})
}

func TestMTimeMillisecond(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(MTimeMillisecond("./testfile/dirfiles/t1.txt"),129)
		gtest.Assert(MTimeMillisecond(""),0)
	})
}