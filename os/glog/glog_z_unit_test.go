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
	"strings"
	"sync"
	"testing"
	"time"

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

func TestSetFile(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetFile("test.log")
	})
}

func TestSetTimeFormat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)

		l.SetTimeFormat("2006-01-02T15:04:05.000Z07:00")
		l.Debug(ctx, "test")

		t.AssertGE(len(strings.Split(w.String(), "[DEBU]")), 1)
		datetime := strings.Trim(strings.Split(w.String(), "[DEBU]")[0], " ")

		_, err := time.Parse("2006-01-02T15:04:05.000Z07:00", datetime)
		t.AssertNil(err)
		_, err = time.Parse("2006-01-02 15:04:05.000", datetime)
		t.AssertNE(err, nil)
		_, err = time.Parse("Mon, 02 Jan 2006 15:04:05 -0700", datetime)
		t.AssertNE(err, nil)
	})
}

func TestSetLevel(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevel(glog.LEVEL_ALL)
		t.Assert(glog.GetLevel()&glog.LEVEL_ALL, glog.LEVEL_ALL)
	})
}

func TestSetAsync(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetAsync(false)
	})
}

func TestSetStdoutPrint(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetStdoutPrint(false)
	})
}

func TestSetHeaderPrint(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetHeaderPrint(false)
	})
}

func TestSetPrefix(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetPrefix("log_prefix")
	})
}

func TestSetConfigWithMap(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetConfigWithMap(map[string]interface{}{
			"level": "all",
		}), nil)
	})
}

func TestSetPath(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetPath("/var/log"), nil)
		t.Assert(glog.GetPath(), "/var/log")
	})
}

func TestSetWriter(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetWriter(os.Stdout)
		t.Assert(glog.GetWriter(), os.Stdout)
	})
}

func TestSetFlags(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetFlags(glog.F_ASYNC)
		t.Assert(glog.GetFlags(), glog.F_ASYNC)
	})
}

func TestSetCtxKeys(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetCtxKeys("SpanId", "TraceId")
		t.Assert(glog.GetCtxKeys(), []string{"SpanId", "TraceId"})
	})
}

func TestPrintStack(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.PrintStack(ctx, 1)
	})
}

func TestSetStack(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetStack(true)
		t.Assert(glog.GetStack(1), "")
	})
}

func TestSetLevelStr(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.Assert(glog.SetLevelStr("all"), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		t.AssertNE(l.SetLevelStr("test"), nil)
	})
}

func TestSetLevelPrefix(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevelPrefix(glog.LEVEL_ALL, "LevelPrefix")
		t.Assert(glog.GetLevelPrefix(glog.LEVEL_ALL), "LevelPrefix")
	})
}

func TestSetLevelPrefixes(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetLevelPrefixes(map[int]string{
			glog.LEVEL_ALL: "ALL_Prefix",
		})
	})
}

func TestSetHandlers(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetHandlers(func(ctx context.Context, in *glog.HandlerInput) {
		})
	})
}

func TestSetWriterColorEnable(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		glog.SetWriterColorEnable(true)
	})
}

func TestInstance(t *testing.T) {
	defaultLog := glog.DefaultLogger().Clone()
	defer glog.SetDefaultLogger(defaultLog)
	gtest.C(t, func(t *gtest.T) {
		t.AssertNE(glog.Instance("gf"), nil)
	})
}

func TestGetConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		config := glog.DefaultLogger().GetConfig()
		t.Assert(config.Path, "")
		t.Assert(config.StdoutPrint, true)
	})
}

func TestWrite(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		len, err := l.Write([]byte("GoFrame"))
		t.AssertNil(err)
		t.Assert(len, 7)
	})
}

func TestChainingTo(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.DefaultLogger().Clone()
		logTo := l.To(os.Stdout)
		t.AssertNE(logTo, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logTo := l.To(os.Stdout)
		t.AssertNE(logTo, nil)
	})
}

func TestChainingPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.DefaultLogger().Clone()
		logPath := l.Path("./")
		t.AssertNE(logPath, nil)
	})
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logPath := l.Path("./")
		t.AssertNE(logPath, nil)
	})
}

func TestChainingCat(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logCat := l.Cat(".gf")
		t.AssertNE(logCat, nil)
	})
}

func TestChainingLevel(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logLevel := l.Level(glog.LEVEL_ALL)
		t.AssertNE(logLevel, nil)
	})
}

func TestChainingLevelStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logLevelStr := l.LevelStr("all")
		t.AssertNE(logLevelStr, nil)
	})
}

func TestChainingSkip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logSkip := l.Skip(1)
		t.AssertNE(logSkip, nil)
	})
}

func TestChainingStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logStack := l.Stack(true)
		t.AssertNE(logStack, nil)
	})
}

func TestChainingStackWithFilter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logStackWithFilter := l.StackWithFilter("gtest")
		t.AssertNE(logStackWithFilter, nil)
	})
}

func TestChainingStdout(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logStdout := l.Stdout(true)
		t.AssertNE(logStdout, nil)
	})
}

func TestChainingHeader(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logHeader := l.Header(true)
		t.AssertNE(logHeader, nil)
	})
}

func TestChainingLine(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logLine := l.Line(true)
		t.AssertNE(logLine, nil)
	})
}

func TestChainingAsync(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		logAsync := l.Async(true)
		t.AssertNE(logAsync, nil)
	})
}

func TestConfigSetDebug(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		l.SetDebug(false)
	})
}

func TestConfigAppendCtxKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		l.AppendCtxKeys("Trace-Id", "Span-Id", "Test")
		l.AppendCtxKeys("Trace-Id-New", "Span-Id-New", "Test")
	})
}

func TestConfigSetPath(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		t.AssertNE(l.SetPath(""), nil)
	})
}

func TestConfigSetStdoutColorDisabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		l := glog.New()
		l.SetStdoutColorDisabled(false)
	})
}

func TestCtx(t *testing.T) {
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

func TestCtxConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		w := bytes.NewBuffer(nil)
		l := glog.NewWithWriter(w)
		m := map[string]interface{}{
			"CtxKeys": g.SliceStr{"Trace-Id", "Span-Id", "Test"},
		}
		var nilMap map[string]interface{}

		err := l.SetConfigWithMap(m)
		t.AssertNil(err)
		err = l.SetConfigWithMap(nilMap)
		t.AssertNE(err, nil)

		ctx := context.WithValue(context.Background(), "Trace-Id", "1234567890")
		ctx = context.WithValue(ctx, "Span-Id", "abcdefg")

		l.Print(ctx, 1, 2, 3)
		t.Assert(gstr.Count(w.String(), "1234567890"), 1)
		t.Assert(gstr.Count(w.String(), "abcdefg"), 1)
		t.Assert(gstr.Count(w.String(), "1 2 3"), 1)
	})
}

func TestConcurrent(t *testing.T) {
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
