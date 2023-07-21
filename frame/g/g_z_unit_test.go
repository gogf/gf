package g_test

import (
	"context"
	"github.com/gogf/gf/v2/errors/gerror"
	"os"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
)

var (
	ctx = context.TODO()
)

func Test_NewVar(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.NewVar(1).Int(), 1)
		t.Assert(g.NewVar(1, true).Int(), 1)
	})
}

func Test_Dump(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.Dump("GoFrame")
	})
}

func Test_DumpTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpTo(os.Stdout, "GoFrame", gutil.DumpOption{})
	})
}

func Test_DumpWithType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpWithType("GoFrame", 123)
	})
}

func Test_DumpWithOption(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.DumpWithOption("GoFrame", gutil.DumpOption{})
	})
}

func Test_Try(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.Try(ctx, func(ctx context.Context) {
			g.Dump("GoFrame")
		})
	})
}

func Test_TryFunc(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		err := g.TryFunc(ctx, func(ctx context.Context) error {
			return nil
		})
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.TryFunc(ctx, func(ctx context.Context) error {
			return gerror.New("GoFrame")
		})
		t.Assert(err.Error(), "GoFrame")
	})
	gtest.C(t, func(t *gtest.T) {
		err := g.TryFunc(ctx, func(ctx context.Context) error {
			g.Throw("GoFrame")
			return nil
		})
		t.Assert(err.Error(), "GoFrame")
	})
}

func Test_TryCatch(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.TryCatch(ctx, func(ctx context.Context) {
			g.Dump("GoFrame")
		}, func(ctx context.Context, exception error) {
			g.Dump(exception)
		})
	})
	gtest.C(t, func(t *gtest.T) {
		g.TryCatch(ctx, func(ctx context.Context) {
			g.Throw("GoFrame")
		}, func(ctx context.Context, exception error) {
			t.Assert(exception.Error(), "GoFrame")
		})
	})
}

func Test_IsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.IsNil(nil), true)
		t.Assert(g.IsNil(0), false)
		t.Assert(g.IsNil("GoFrame"), false)
	})
}

func Test_IsEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(g.IsEmpty(nil), true)
		t.Assert(g.IsEmpty(0), true)
		t.Assert(g.IsEmpty("GoFrame"), false)
	})
}

func Test_SetDebug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		g.SetDebug(true)
	})
}

func Test_Object(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(g.Client(), nil)
		t.AssertNE(g.Server(), nil)
		t.AssertNE(g.TCPServer(), nil)
		t.AssertNE(g.UDPServer(), nil)
		t.AssertNE(g.View(), nil)
		t.AssertNE(g.Config(), nil)
		t.AssertNE(g.Cfg(), nil)
		t.AssertNE(g.Resource(), nil)
		t.AssertNE(g.I18n(), nil)
		t.AssertNE(g.Res(), nil)
		t.AssertNE(g.Log(), nil)
		t.AssertNE(g.Validator(), nil)
	})
}
