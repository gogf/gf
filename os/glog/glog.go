// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glog implements powerful and easy-to-use leveled logging functionality.
package glog

import (
	"context"

	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/os/grpool"
	"github.com/gogf/gf/v2/util/gconv"
)

// ILogger is the API interface for logger.
type ILogger interface {
	Print(ctx context.Context, v ...any)
	Printf(ctx context.Context, format string, v ...any)
	Debug(ctx context.Context, v ...any)
	Debugf(ctx context.Context, format string, v ...any)
	Info(ctx context.Context, v ...any)
	Infof(ctx context.Context, format string, v ...any)
	Notice(ctx context.Context, v ...any)
	Noticef(ctx context.Context, format string, v ...any)
	Warning(ctx context.Context, v ...any)
	Warningf(ctx context.Context, format string, v ...any)
	Error(ctx context.Context, v ...any)
	Errorf(ctx context.Context, format string, v ...any)
	Critical(ctx context.Context, v ...any)
	Criticalf(ctx context.Context, format string, v ...any)
	Panic(ctx context.Context, v ...any)
	Panicf(ctx context.Context, format string, v ...any)
	Fatal(ctx context.Context, v ...any)
	Fatalf(ctx context.Context, format string, v ...any)
}

const (
	commandEnvKeyForDebug = "gf.glog.debug"
)

var (
	// Ensure Logger implements ILogger.
	_ ILogger = &Logger{}

	// Default logger object, for package method usage.
	defaultLogger = New()

	// Goroutine pool for async logging output.
	// It uses only one asynchronous worker to ensure log sequence.
	asyncPool = grpool.New(1)

	// defaultDebug enables debug level or not in default,
	// which can be configured using command option or system environment.
	defaultDebug = true
)

func init() {
	defaultDebug = gconv.Bool(command.GetOptWithEnv(commandEnvKeyForDebug, "true"))
	SetDebug(defaultDebug)
}

// DefaultLogger returns the default logger.
func DefaultLogger() *Logger {
	return defaultLogger
}

// SetDefaultLogger sets the default logger for package glog.
// Note that there might be concurrent safety issue if calls this function
// in different goroutines.
func SetDefaultLogger(l *Logger) {
	defaultLogger = l
}
