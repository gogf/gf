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
	"testing"

	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/test/gtest"
)

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
		slogger.ErrorContext(ctx, "error message")

		output := buffer.String()
		t.Assert(strings.Contains(output, "[DEBU]"), true)
		t.Assert(strings.Contains(output, "[INFO]"), true)
		t.Assert(strings.Contains(output, "[NOTI]"), true)
		t.Assert(strings.Contains(output, "[WARN]"), true)
		t.Assert(strings.Contains(output, "[ERRO]"), true)
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
		t.Assert(strings.Contains(output, "plain attrs key=value user.name=john"), true)
		t.Assert(strings.Contains(output, "grouped attrs app=demo req.id=7 req.path=/x"), true)
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
		t.Assert(handler.Enabled(ctx, slog.LevelWarn), true)
		t.Assert(handler.Enabled(ctx, slog.LevelError), true)
	})
}
