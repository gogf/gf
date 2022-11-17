package gdebug_test

import (
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
	"testing"
)

func Test_CallerPackage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.CallerPackage(), "github.com/gogf/gf/v2/test/gtest")
	})
}

func Test_CallerFunction(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.CallerFunction(), "C")
	})
}

func Test_CallerFilePath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFilePath(), "gtest_util.go"), true)
	})
}

func Test_CallerDirectory(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerDirectory(), "test\\gtest"), true)
	})
}

func Test_CallerFileLine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFileLine(), "gtest_util.go:36"), true)
	})
}

func Test_CallerFileLineShort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFileLineShort(), "gtest_util.go:36"), true)
	})
}

func Test_FuncPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.FuncPath(Test_FuncPath), "github.com/gogf/gf/v2/debug/gdebug_test.Test_FuncPath")
	})
}

func Test_FuncName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.FuncName(Test_FuncName), "gdebug_test.Test_FuncName")
	})
}

func Test_PrintStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gdebug.PrintStack()
	})
}

func Test_GoroutineId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(gdebug.GoroutineId(), 0)
	})
}

func Test_Stack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.Stack(), "gtest_util.go:36"), true)
	})
}

func Test_StackWithFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.StackWithFilter([]string{"github.com"}), "gtest_util.go:36"), true)
	})
}

func Test_BinVersion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(len(gdebug.BinVersion()), 0)
	})
}

func Test_BinVersionMd5(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(len(gdebug.BinVersionMd5()), 0)
	})
}
