//have test 100%
package gfile

import (
	"testing"
	"github.com/gogf/gf/g/test/gtest"

)

func TestMTime(t *testing.T) {
	gtest.Case(t, func() {
		//拷贝到其它地方，再测试的时候，这个文件的修改值会变，所以用大于来断言
		gtest.AssertGT(MTime("./testfile/dirfiles/t1.txt"),1454883732)
		gtest.Assert(MTime(""),0)
	})
}

func TestMTimeMillisecond(t *testing.T) {
	gtest.Case(t, func() {
		gtest.AssertGT(MTimeMillisecond("./testfile/dirfiles/t1.txt"),0) //这里有值
		gtest.Assert(MTimeMillisecond(""),0)
	})
}