// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package intlog provides internal logging for GoFrame development usage only.
package intlog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/internal/utils"
	"go.opentelemetry.io/otel/trace"
	"path/filepath"
	"time"
)

const (
	stackFilterKey = "/internal/intlog"
)

var (
	// isGFDebug marks whether printing GoFrame debug information.
	isGFDebug = false
)

func init() {
	isGFDebug = utils.IsDebugEnabled()
}

// SetEnabled enables/disables the internal logging manually.
// Note that this function is not concurrent safe, be aware of the DATA RACE.
func SetEnabled(enabled bool) {
	// If they're the same, it does not write the `isGFDebug` but only reading operation.
	if isGFDebug != enabled {
		isGFDebug = enabled
	}
}

// Print prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Print(ctx context.Context, v ...interface{}) {
	doPrint(ctx, fmt.Sprint(v...), false)
}

// Printf prints `v` with format `format` using fmt.Printf.
// The parameter `v` can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	doPrint(ctx, fmt.Sprintf(format, v...), false)
}

// Error prints `v` with newline using fmt.Println.
// The parameter `v` can be multiple variables.
func Error(ctx context.Context, v ...interface{}) {
	doPrint(ctx, fmt.Sprint(v...), true)
}

// Errorf prints `v` with format `format` using fmt.Printf.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	doPrint(ctx, fmt.Sprintf(format, v...), true)
}

func doPrint(ctx context.Context, content string, stack bool) {
	if !isGFDebug {
		return
	}
	buffer := bytes.NewBuffer(nil)
	buffer.WriteString(now())
	buffer.WriteString(" [INTE] ")
	buffer.WriteString(file())
	if s := traceIdStr(ctx); s != "" {
		buffer.WriteString(" " + s)
	}
	buffer.WriteString(content)
	buffer.WriteString("\n")
	if stack {
		buffer.WriteString(gdebug.StackWithFilter(stackFilterKey))
	}
	fmt.Print(buffer.String())
}

// traceIdStr retrieves and returns the trace id string for logging output.
func traceIdStr(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	spanCtx := trace.SpanContextFromContext(ctx)
	if traceId := spanCtx.TraceID(); traceId.IsValid() {
		return "{" + traceId.String() + "}"
	}
	return ""
}

// now returns current time string.
func now() string {
	return time.Now().Format("2006-01-02 15:04:05.000")
}

// file returns caller file name along with its line number.
func file() string {
	_, p, l := gdebug.CallerWithFilter(stackFilterKey)
	return fmt.Sprintf(`%s:%d`, filepath.Base(p), l)
}
