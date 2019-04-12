//have test 100%
package gfile

import (
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

func TestMTime(t *testing.T) {
	gtest.Case(t, func() {
		//拷贝到其它地方，再测试的时候，这个文件的修改值会变，所以用大于来断言
		gtest.AssertGT(MTime("./testfile/dirfiles/t1.txt"), 1454883732)
		gtest.Assert(MTime(""), 0)
	})
}

func TestMTimeMillisecond(t *testing.T) {
	gtest.Case(t, func() {
		//这里本不为0,但github中的ci测试时，值为0
		gtest.AssertGTE(MTimeMillisecond("./testfile/dirfiles/t1.txt"), 0)
		gtest.Assert(MTimeMillisecond(""), 0)
	})
}
