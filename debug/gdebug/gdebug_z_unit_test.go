package gdebug_test

import (
	"testing"

	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func TestCallerPackage(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.CallerPackage(), "github.com/gogf/gf/v2/test/gtest")
	})
}

func TestCallerFunction(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.CallerFunction(), "C")
	})
}

func TestCallerFilePath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFilePath(), "gtest_util.go"), true)
	})
}

func TestCallerDirectory(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerDirectory(), "gtest"), true)
	})
}

func TestCallerFileLine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFileLine(), "gtest_util.go:35"), true)
	})
}

func TestCallerFileLineShort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.CallerFileLineShort(), "gtest_util.go:35"), true)
	})
}

func TestFuncPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.FuncPath(Test_FuncPath), "github.com/gogf/gf/v2/debug/gdebug_test.Test_FuncPath")
	})
}

func TestFuncName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gdebug.FuncName(Test_FuncName), "gdebug_test.Test_FuncName")
	})
}

func TestPrintStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		gdebug.PrintStack()
	})
}

func TestGoroutineId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(gdebug.GoroutineId(), 0)
	})
}

func TestStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.Stack(), "gtest_util.go:35"), true)
	})
}

func TestStackWithFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gstr.Contains(gdebug.StackWithFilter([]string{"github.com"}), "gtest_util.go:35"), true)
	})
}

func TestBinVersion(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(len(gdebug.BinVersion()), 0)
	})
}

func TestBinVersionMd5(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertGT(len(gdebug.BinVersionMd5()), 0)
	})
}
