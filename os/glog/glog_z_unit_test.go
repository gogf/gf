// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"bytes"
	"context"
	"os"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

func TestCase(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(glog.Instance(), nil)
	})
}

func TestDefaultLogger(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(defaultLog, nil)
		log := glog.New()
		glog.SetDefaultLogger(log)
		t.AssertEQ(glog.DefaultLogger(), defaultLog)
		t.AssertEQ(glog.Expose(), defaultLog)
	})
}

func TestAPI(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		glog.Print(ctx, "Print")
		glog.Printf(ctx, "%s", "Printf")
		glog.Info(ctx, "Info")
		glog.Infof(ctx, "%s", "Infof")
		glog.Debug(ctx, "Debug")
		glog.Debugf(ctx, "%s", "Debugf")
		glog.Notice(ctx, "Notice")
		glog.Noticef(ctx, "%s", "Noticef")
		glog.Warning(ctx, "Warning")
		glog.Warningf(ctx, "%s", "Warningf")
		glog.Error(ctx, "Error")
		glog.Errorf(ctx, "%s", "Errorf")
		glog.Critical(ctx, "Critical")
		glog.Criticalf(ctx, "%s", "Criticalf")
	})
}

func TestChaining(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)

	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(glog.Cat("module"), nil)
		t.AssertNE(glog.File("test.log"), nil)
		t.AssertNE(glog.Level(glog.LEVEL_ALL), nil)
		t.AssertNE(glog.LevelStr("all"), nil)
		t.AssertNE(glog.Skip(1), nil)
		t.AssertNE(glog.Stack(false), nil)
		t.AssertNE(glog.StackWithFilter("none"), nil)
		t.AssertNE(glog.Stdout(false), nil)
		t.AssertNE(glog.Header(false), nil)
		t.AssertNE(glog.Line(false), nil)
		t.AssertNE(glog.Async(false), nil)
	})
}

func Test_SetFile(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetFile("test.log")
	})
}

func Test_SetLevel(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevel(glog.LEVEL_ALL)
		t.Assert(glog.GetLevel(), glog.LEVEL_ALL)
	})
}

func Test_SetAsync(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetAsync(false)
	})
}

func Test_SetStdoutPrint(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetStdoutPrint(false)
	})
}

func Test_SetHeaderPrint(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetHeaderPrint(false)
	})
}

func Test_SetPrefix(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetPrefix("log_prefix")
	})
}

func Test_SetConfigWithMap(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetConfigWithMap(map[string]interface{}{
			"level": "all",
		}), nil)
	})
}

func Test_SetPath(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetPath("/var/log"), nil)
		t.Assert(glog.GetPath(), "/var/log")
	})
}

func Test_SetWriter(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetWriter(os.Stdout)
		t.Assert(glog.GetWriter(), os.Stdout)
	})
}

func Test_SetFlags(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetFlags(glog.F_ASYNC)
		t.Assert(glog.GetFlags(), glog.F_ASYNC)
	})
}

func Test_SetCtxKeys(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetCtxKeys("SpanId", "TraceId")
		t.Assert(glog.GetCtxKeys(), []string{"SpanId", "TraceId"})
	})
}

func Test_PrintStack(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.PrintStack(ctx, 1)
	})
}

func Test_SetStack(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetStack(true)
		t.Assert(glog.GetStack(1), "")
	})
}

func Test_SetLevelStr(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetLevelStr("all"), nil)
	})
}

func Test_SetLevelPrefix(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevelPrefix(glog.LEVEL_ALL, "LevelPrefix")
		t.Assert(glog.GetLevelPrefix(glog.LEVEL_ALL), "LevelPrefix")
	})
}

func Test_SetLevelPrefixes(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevelPrefixes(map[int]string{
			glog.LEVEL_ALL: "ALL_Prefix",
		})
	})
}

func Test_SetHandlers(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetHandlers(func(ctx context.Context, in *glog.HandlerInput) {
		})
	})
}

func Test_SetWriterColorEnable(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetWriterColorEnable(true)
	})
}

func Test_Ctx(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		l.SetCtxKeys("Trace-Id", "Span-Id", "Test")
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)
	})
}

func Test_Ctx_Config(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		m := map[string]interface{}{
			"CtxKeys": g.SliceStr{"Trace-Id", "Span-Id", "Test"},
		}
		err := l.SetConfigWithMap(m)
		t.AssertNil(err)
		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)
	})
}

func Test_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c := 1000
		l := glog.New()
		s := "@1234567890#"
		f := "test.log"
		p := gfile.Temp(gtime.TimestampNanoStr())
		t.Assert(l.SetPath(p), nil)
		defer gfile.Remove(p)
		wg := sync.WaitGroup{}
		ch := make(chan struct{})
		for i := 0; i < c; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				<-ch
				l.File(f).Stdout(false).Print(ctx, s)
			}()
		}
		close(ch)
		wg.Wait()
		content := gfile.GetContents(gfile.Join(p, f))
		t.Assert(gstr.Count(content, s), c)
	})
}
