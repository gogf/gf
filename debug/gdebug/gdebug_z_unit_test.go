package gdebug_test

import (
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"testing"
)

func Test_Case(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.CallerPackage(), "github.com/gogf/gf/v2/test/gtest")
		t.Assert(gdebug.CallerFunction(), "C")
		t.Assert(gstr.Contains(gdebug.CallerFilePath(), "gtest_util.go"), true)
		t.Assert(gstr.Contains(gdebug.CallerDirectory(), "test\\gtest"), true)
		t.Assert(gstr.Contains(gdebug.CallerFileLine(), "gtest_util.go:36"), true)
		t.Assert(gstr.Contains(gdebug.CallerFileLineShort(), "gtest_util.go:36"), true)
		t.Assert(gdebug.FuncPath(Test_Case), "github.com/gogf/gf/v2/debug/gdebug_test.Test_Case")
		t.Assert(gdebug.FuncName(Test_Case), "gdebug_test.Test_Case")
		gdebug.PrintStack()
		t.AssertGT(gdebug.GoroutineId(), 0)
		t.Assert(gstr.Contains(gdebug.Stack(), "gtest_util.go:36"), true)
		t.Assert(gstr.Contains(gdebug.StackWithFilter([]string{"github.com"}), "gtest_util.go:36"), true)
		t.AssertGT(len(gdebug.BinVersion()), 0)
		t.AssertGT(len(gdebug.BinVersionMd5()), 0)
	})
}
