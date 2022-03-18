// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/gogf/gf/v2/internal/utils"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/debug/gdebug"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gfpool"
	"github.com/gogf/gf/v2/os/gmlock"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
)

// Logger is the struct for logging management.
type Logger struct {
	init   *gtype.Bool // Initialized.
	parent *Logger     // Parent logger, if it is not empty, it means the logger is used in chaining function.
	config Config      // Logger configuration.
}

const (
	defaultFileFormat                 = `{Y-m-d}.log`
	defaultFileFlags                  = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	defaultFilePerm                   = os.FileMode(0666)
	defaultFileExpire                 = time.Minute
	pathFilterKey                     = "/os/glog/glog"
	memoryLockPrefixForPrintingToFile = "glog.printToFile:"
)

const (
	F_ASYNC      = 1 << iota // Print logging content asynchronously。
	F_FILE_LONG              // Print full file name and line number: /a/b/c/d.go:23.
	F_FILE_SHORT             // Print final file name element and line number: d.go:23. overrides F_FILE_LONG.
	F_TIME_DATE              // Print the date in the local time zone: 2009-01-23.
	F_TIME_TIME              // Print the time in the local time zone: 01:23:23.
	F_TIME_MILLI             // Print the time with milliseconds in the local time zone: 01:23:23.675.
	F_CALLER_FN              // Print Caller function name and package: main.main
	F_TIME_STD   = F_TIME_DATE | F_TIME_MILLI
)

// New creates and returns a custom logger.
func New() *Logger {
	logger := &Logger{
		init:   gtype.NewBool(),
		config: DefaultConfig(),
	}
	return logger
}

// NewWithWriter creates and returns a custom logger with io.Writer.
func NewWithWriter(writer io.Writer) *Logger {
	l := New()
	l.SetWriter(writer)
	return l
}

// Clone returns a new logger, which is the clone the current logger.
// It's commonly used for chaining operations.
func (l *Logger) Clone() *Logger {
	newLogger := New()
	newLogger.config = l.config
	newLogger.parent = l
	return newLogger
}

// getFilePath returns the logging file path.
// The logging file name must have extension name of "log".
func (l *Logger) getFilePath(now time.Time) string {
	// Content containing "{}" in the file name is formatted using gtime.
	file, _ := gregex.ReplaceStringFunc(`{.+?}`, l.config.File, func(s string) string {
		return gtime.New(now).Format(strings.Trim(s, "{}"))
	})
	file = gfile.Join(l.config.Path, file)
	return file
}

// print prints `s` to defined writer, logging file or passed `std`.
func (l *Logger) print(ctx context.Context, level int, values ...interface{}) {
	// Lazy initialize for rotation feature.
	// It uses atomic reading operation to enhance the performance checking.
	// It here uses CAP for performance and concurrent safety.
	p := l
	if p.parent != nil {
		p = p.parent
	}
	// It just initializes once for each logger.
	if p.config.RotateSize > 0 || p.config.RotateExpire > 0 {
		if !p.init.Val() && p.init.Cas(false, true) {
			gtimer.AddOnce(context.Background(), p.config.RotateCheckInterval, p.rotateChecksTimely)
			intlog.Printf(ctx, "logger rotation initialized: every %s", p.config.RotateCheckInterval.String())
		}
	}

	var (
		now   = time.Now()
		input = &HandlerInput{
			Logger:       l,
			Buffer:       bytes.NewBuffer(nil),
			Time:         now,
			Color:        defaultLevelColor[level],
			Level:        level,
			handlerIndex: -1,
		}
	)
	if l.config.HeaderPrint {
		// Time.
		timeFormat := ""
		if l.config.Flags&F_TIME_DATE > 0 {
			timeFormat += "2006-01-02"
		}
		if l.config.Flags&F_TIME_TIME > 0 {
			if timeFormat != "" {
				timeFormat += " "
			}
			timeFormat += "15:04:05"
		}
		if l.config.Flags&F_TIME_MILLI > 0 {
			if timeFormat != "" {
				timeFormat += " "
			}
			timeFormat += "15:04:05.000"
		}
		if len(timeFormat) > 0 {
			input.TimeFormat = now.Format(timeFormat)
		}

		// Level string.
		input.LevelFormat = l.getLevelPrefixWithBrackets(level)

		// Caller path and Fn name.
		if l.config.Flags&(F_FILE_LONG|F_FILE_SHORT|F_CALLER_FN) > 0 {
			callerFnName, path, line := gdebug.CallerWithFilter(
				[]string{utils.StackFilterKeyForGoFrame},
				l.config.StSkip,
			)
			if l.config.Flags&F_CALLER_FN > 0 {
				if len(callerFnName) > 2 {
					input.CallerFunc = fmt.Sprintf(`[%s]`, callerFnName)
				}
			}
			if line >= 0 && len(path) > 1 {
				if l.config.Flags&F_FILE_LONG > 0 {
					input.CallerPath = fmt.Sprintf(`%s:%d:`, path, line)
				}
				if l.config.Flags&F_FILE_SHORT > 0 {
					input.CallerPath = fmt.Sprintf(`%s:%d:`, gfile.Basename(path), line)
				}
			}
		}
		// Prefix.
		if len(l.config.Prefix) > 0 {
			input.Prefix = l.config.Prefix
		}
	}
	// Convert value to string.
	if ctx != nil {
		// Tracing values.
		spanCtx := trace.SpanContextFromContext(ctx)
		if traceId := spanCtx.TraceID(); traceId.IsValid() {
			input.CtxStr = traceId.String()
		}
		// Context values.
		if len(l.config.CtxKeys) > 0 {
			for _, ctxKey := range l.config.CtxKeys {
				var ctxValue interface{}
				if ctxValue = ctx.Value(ctxKey); ctxValue == nil {
					ctxValue = ctx.Value(gctx.StrKey(gconv.String(ctxKey)))
				}
				if ctxValue != nil {
					if input.CtxStr != "" {
						input.CtxStr += ", "
					}
					input.CtxStr += gconv.String(ctxValue)
				}
			}
		}
		if input.CtxStr != "" {
			input.CtxStr = "{" + input.CtxStr + "}"
		}
	}
	var tempStr string
	for _, v := range values {
		tempStr = gconv.String(v)
		if len(input.Content) > 0 {
			if input.Content[len(input.Content)-1] == '\n' {
				// Remove one blank line(\n\n).
				if len(tempStr) > 0 && tempStr[0] == '\n' {
					input.Content += tempStr[1:]
				} else {
					input.Content += tempStr
				}
			} else {
				input.Content += " " + tempStr
			}
		} else {
			input.Content = tempStr
		}
	}
	if l.config.Flags&F_ASYNC > 0 {
		input.IsAsync = true
		err := asyncPool.Add(ctx, func(ctx context.Context) {
			input.Next(ctx)
		})
		if err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	} else {
		input.Next(ctx)
	}
}

// doDefaultPrint outputs the logging content according configuration.
func (l *Logger) doDefaultPrint(ctx context.Context, input *HandlerInput) *bytes.Buffer {
	var buffer *bytes.Buffer
	if l.config.Writer == nil {
		// Allow output to stdout?
		if l.config.StdoutPrint {
			if buf := l.printToStdout(ctx, input); buf != nil {
				buffer = buf
			}
		}

		// Output content to disk file.
		if l.config.Path != "" {
			if buf := l.printToFile(ctx, input.Time, input); buf != nil {
				buffer = buf
			}
		}
	} else {
		// Output to custom writer.
		if buf := l.printToWriter(ctx, input); buf != nil {
			buffer = buf
		}
	}
	return buffer
}

// printToWriter writes buffer to writer.
func (l *Logger) printToWriter(ctx context.Context, input *HandlerInput) *bytes.Buffer {
	if l.config.Writer != nil {
		var (
			buffer = input.getRealBuffer(l.config.WriterColorEnable)
		)
		if _, err := l.config.Writer.Write(buffer.Bytes()); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
		return buffer
	}
	return nil
}

// printToStdout outputs logging content to stdout.
func (l *Logger) printToStdout(ctx context.Context, input *HandlerInput) *bytes.Buffer {
	if l.config.StdoutPrint {
		var (
			err    error
			buffer = input.getRealBuffer(!l.config.StdoutColorDisabled)
		)
		// This will lose color in Windows os system.
		// if _, err := os.Stdout.Write(input.getRealBuffer(true).Bytes()); err != nil {

		// This will print color in Windows os system.
		if _, err = fmt.Fprint(color.Output, buffer.String()); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
		return buffer
	}
	return nil
}

// printToFile outputs logging content to disk file.
func (l *Logger) printToFile(ctx context.Context, t time.Time, in *HandlerInput) *bytes.Buffer {
	var (
		buffer        = in.getRealBuffer(l.config.WriterColorEnable)
		logFilePath   = l.getFilePath(t)
		memoryLockKey = memoryLockPrefixForPrintingToFile + logFilePath
	)
	gmlock.Lock(memoryLockKey)
	defer gmlock.Unlock(memoryLockKey)

	// Rotation file size checks.
	if l.config.RotateSize > 0 {
		if gfile.Size(logFilePath) > l.config.RotateSize {
			l.rotateFileBySize(ctx, t)
		}
	}
	// Logging content outputting to disk file.
	if file := l.getFilePointer(ctx, logFilePath); file == nil {
		intlog.Errorf(ctx, `got nil file pointer for: %s`, logFilePath)
	} else {
		if _, err := file.Write(buffer.Bytes()); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
		if err := file.Close(); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}
	return buffer
}

// getFilePointer retrieves and returns a file pointer from file pool.
func (l *Logger) getFilePointer(ctx context.Context, path string) *gfpool.File {
	file, err := gfpool.Open(
		path,
		defaultFileFlags,
		defaultFilePerm,
		defaultFileExpire,
	)
	if err != nil {
		// panic(err)
		intlog.Errorf(ctx, `%+v`, err)
	}
	return file
}

// printStd prints content `s` without stack.
func (l *Logger) printStd(ctx context.Context, level int, value ...interface{}) {
	l.print(ctx, level, value...)
}

// printStd prints content `s` with stack check.
func (l *Logger) printErr(ctx context.Context, level int, value ...interface{}) {
	if l.config.StStatus == 1 {
		if s := l.GetStack(); s != "" {
			value = append(value, "\nStack:\n"+s)
		}
	}
	// In matter of sequence, do not use stderr here, but use the same stdout.
	l.print(ctx, level, value...)
}

// format formats `values` using fmt.Sprintf.
func (l *Logger) format(format string, value ...interface{}) string {
	return fmt.Sprintf(format, value...)
}

// PrintStack prints the caller stack,
// the optional parameter `skip` specify the skipped stack offset from the end point.
func (l *Logger) PrintStack(ctx context.Context, skip ...int) {
	if s := l.GetStack(skip...); s != "" {
		l.Print(ctx, "Stack:\n"+s)
	} else {
		l.Print(ctx)
	}
}

// GetStack returns the caller stack content,
// the optional parameter `skip` specify the skipped stack offset from the end point.
func (l *Logger) GetStack(skip ...int) string {
	stackSkip := l.config.StSkip
	if len(skip) > 0 {
		stackSkip += skip[0]
	}
	filters := []string{pathFilterKey}
	if l.config.StFilter != "" {
		filters = append(filters, l.config.StFilter)
	}
	return gdebug.StackWithFilters(filters, stackSkip)
}

// GetConfig returns the configuration of current Logger.
func (l *Logger) GetConfig() Config {
	return l.config
}
