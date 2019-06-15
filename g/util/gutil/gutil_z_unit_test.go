package gutil_test

import (
	"testing"

	"github.com/gogf/gf/g/test/gtest"
	"github.com/gogf/gf/g/util/gutil"
)

func Test_Dump(t *testing.T) {
	gtest.Case(t, func() {
		gutil.Dump(map[int]int{
			100: 100,
		})
		gtest.Assert("", "")
	})

	gtest.Case(t, func() {
		gutil.Dump(map[string]interface{}{"": func() {}})
		gtest.Assert("", "")
	})

	gtest.Case(t, func() {
		gutil.Dump([]byte("gutil Dump test"))
		gtest.Assert("", "")
	})
}

func Test_PrintBacktrace(t *testing.T) {
	gtest.Case(t, func() {
		gutil.PrintBacktrace()
		gtest.Assert("", "")
	})
}

func Test_TryCatch(t *testing.T) {

	gutil.TryCatch(func() {
	}, func(err interface{}) {
	})
	gtest.Assert("", "")

	gutil.TryCatch(func() {
	})
	gtest.Assert("", "")

	gutil.TryCatch(func() {
		panic("gutil TryCatch test")
	}, func(err interface{}) {
	})
	gtest.Assert("", "")
}

func Test_IsEmpty(t *testing.T) {
	gtest.Assert(gutil.IsEmpty(1), false)
}

func Test_Throw(t *testing.T) {
	gtest.Case(t, func() {
		defer func() {
			if e := recover(); e != nil {
				gtest.Assert(e, "gutil Throw test")
			}
		}()

		gutil.Throw("gutil Throw test")
	})
}
