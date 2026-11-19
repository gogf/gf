// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog_test

import (
	"bytes"
	"context"
	"log/slog"
	"strings"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

type slogLogValuer struct {
	value string
}

func (v slogLogValuer) LogValue() slog.Value {
	return slog.StringValue(v.value)
}

type synchronizedBuffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (b *synchronizedBuffer) Write(p []byte) (int, error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.Write(p)
}

func (b *synchronizedBuffer) String() string {
	b.mu.Lock()
	defer b.mu.Unlock()
	return b.buffer.String()
}

func Test_SlogHandler_Levels(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)

		slogger.DebugContext(ctx, "debug message")
		slogger.InfoContext(ctx, "info message")
		slogger.Log(ctx, slog.LevelInfo+2, "notice message")
		slogger.WarnContext(ctx, "warn message")
		slogger.Log(ctx, slog.LevelWarn+2, "warn plus message")
		slogger.ErrorContext(ctx, "error message")

		output := buffer.String()
		t.Assert(strings.Contains(output, "[DEBU] debug message"), true)
		t.Assert(strings.Contains(output, "[INFO] info message"), true)
		t.Assert(strings.Contains(output, "[NOTI] notice message"), true)
		t.Assert(strings.Contains(output, "[WARN] warn message"), true)
		t.Assert(strings.Contains(output, "[WARN] warn plus message"), true)
		t.Assert(strings.Contains(output, "[ERRO] error message"), true)
	})
}

func Test_SlogHandler_AttrsAndGroups(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)

		slogger.InfoContext(ctx, "plain attrs",
			slog.String("key", "value"),
			slog.Group("user", slog.String("name", "john")),
		)
		slogger.With("app", "demo").
			WithGroup("req").
			With("id", 7).
			InfoContext(ctx, "grouped attrs", slog.String("path", "/x"))

		output := buffer.String()
		t.Assert(strings.Contains(output, "plain attrs key value user.name john"), true)
		t.Assert(strings.Contains(output, "grouped attrs app demo req.id 7 req.path /x"), true)
	})
}

func Test_SlogHandler_StructureHandler(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)
		logger.SetHandlers(glog.HandlerStructure)

		slogger.InfoContext(ctx, "msg",
			slog.String("k", "a b"),
			slog.Group("user", slog.Int("id", 1)),
			slog.String("nl", "x\n2026-01-01 [ERRO] forged"),
		)

		output := buffer.String()
		t.Assert(strings.Contains(output, `Content=msg k="a b" user.id=1`), true)
		t.Assert(strings.Contains(output, `nl="x\n2026-01-01 [ERRO] forged"`), true)
		t.Assert(strings.Contains(output, "\n2026-01-01"), false)
	})
}

func Test_SlogHandler_LogValuer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)

		slogger.InfoContext(ctx, "logvaluer", slog.Any("lv", slogLogValuer{value: "resolved"}))

		output := buffer.String()
		t.Assert(strings.Contains(output, "logvaluer lv resolved"), true)
	})
}

func Test_SlogHandler_EmptyAttr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = context.Background()
			buffer = bytes.NewBuffer(nil)
			logger = glog.NewWithWriter(buffer)
			// Record.AddAttrs omits empty Attr, so feed the handler directly.
			handler = glog.NewSlogHandler(logger).WithAttrs([]slog.Attr{
				{}, slog.String("k", "v"),
			})
			slogger = slog.New(handler)
		)
		logger.SetLevel(glog.LEVEL_ALL)

		slogger.InfoContext(ctx, "empty attr")

		output := buffer.String()
		t.Assert(strings.Contains(output, "empty attr k v"), true)
		t.Assert(strings.Contains(output, "<nil>"), false)
	})
}

func Test_SlogHandler_Caller(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)
		logger.SetFlags(glog.F_FILE_SHORT | glog.F_CALLER_FN)

		slogger.InfoContext(ctx, "caller check")

		output := buffer.String()
		t.Assert(strings.Contains(output, "glog_test.Test_SlogHandler_Caller"), true)
		t.Assert(strings.Contains(output, "glog_z_unit_slog_test.go:"), true)
	})
}

// Test_SlogHandler_ErrorStack pins the -trimpath behavior: without anchoring the
// stack at the record PC, log/slog frames are not filtered because runtime.GOROOT()
// is empty in such builds.
func Test_SlogHandler_ErrorStack(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			slogger = slog.New(glog.NewSlogHandler(logger))
		)
		logger.SetLevel(glog.LEVEL_ALL)
		logger.SetFlags(glog.F_CALLER_FN)

		slogger.ErrorContext(ctx, "error check")

		output := buffer.String()
		t.Assert(strings.Contains(output, "glog_test.Test_SlogHandler_ErrorStack"), true)
		t.Assert(strings.Contains(output, "log/slog."), false)
	})
}

func Test_SlogHandler_WithAttrsImmutable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			buffer  = bytes.NewBuffer(nil)
			logger  = glog.NewWithWriter(buffer)
			handler = glog.NewSlogHandler(logger)
		)
		logger.SetLevel(glog.LEVEL_ALL)

		// This chain leaves spare capacity in the attributes slice, so a derived
		// handler without a cloned slice would overwrite the sibling's attributes.
		base := handler.
			WithAttrs([]slog.Attr{slog.String("a", "1")}).
			WithAttrs([]slog.Attr{slog.String("b", "2")}).
			WithAttrs([]slog.Attr{slog.String("c", "3")})
		x := base.WithAttrs([]slog.Attr{slog.String("x", "9")})
		y := base.WithAttrs([]slog.Attr{slog.String("y", "8")})

		slog.New(x).InfoContext(ctx, "from x")
		slog.New(y).InfoContext(ctx, "from y")

		output := buffer.String()
		t.Assert(strings.Contains(output, "from x a 1 b 2 c 3 x 9"), true)
		t.Assert(strings.Contains(output, "from y a 1 b 2 c 3 y 8"), true)
	})
}

func Test_SlogHandler_Concurrent(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx    = context.Background()
			buffer = &synchronizedBuffer{}
			logger = glog.NewWithWriter(buffer)
			base   = slog.New(glog.NewSlogHandler(logger)).With("app", "demo")
			wg     sync.WaitGroup
		)
		logger.SetLevel(glog.LEVEL_INFO)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				slogger := base.WithGroup("req").With("id", i)
				for j := 0; j < 20; j++ {
					slogger.InfoContext(ctx, "concurrent", slog.Int("j", j))
				}
			}(i)
		}
		wg.Wait()
		t.Assert(strings.Count(buffer.String(), "concurrent app demo req.id "), 200)
	})
}

func Test_SlogHandler_Enabled(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx     = context.Background()
			logger  = glog.NewWithWriter(bytes.NewBuffer(nil))
			handler = glog.NewSlogHandler(logger)
		)

		logger.SetLevel(glog.LEVEL_WARN | glog.LEVEL_ERRO | glog.LEVEL_CRIT)

		t.Assert(handler.Enabled(ctx, slog.LevelDebug), false)
		t.Assert(handler.Enabled(ctx, slog.LevelInfo), false)
		t.Assert(handler.Enabled(ctx, slog.LevelInfo+2), false)
		t.Assert(handler.Enabled(ctx, slog.LevelWarn), true)
		t.Assert(handler.Enabled(ctx, slog.LevelWarn+2), true)
		t.Assert(handler.Enabled(ctx, slog.LevelError), true)
		t.Assert(handler.Enabled(ctx, slog.LevelError+2), true)

		logger.SetLevel(glog.LEVEL_NOTI)

		t.Assert(handler.Enabled(ctx, slog.LevelInfo+2), true)
	})
}
