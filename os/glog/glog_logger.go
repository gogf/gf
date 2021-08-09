// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"fmt"
	"github.com/fatih/color"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfpool"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/os/gtimer"
	"go.opentelemetry.io/otel/trace"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/debug/gdebug"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// Logger is the struct for logging management.
type Logger struct {
	ctx    context.Context // Context for logging.
	init   *gtype.Bool     // Initialized.
	parent *Logger         // Parent logger, if it is not empty, it means the logger is used in chaining function.
	config Config          // Logger configuration.
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
	logger := New()
	logger.ctx = l.ctx
	logger.config = l.config
	logger.parent = l
	return logger
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

// print prints <s> to defined writer, logging file or passed <std>.
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
			gtimer.AddOnce(p.config.RotateCheckInterval, p.rotateChecksTimely)
			intlog.Printf(ctx, "logger rotation initialized: every %s", p.config.RotateCheckInterval.String())
		}
	}

	var (
		now   = time.Now()
		input = &HandlerInput{
			logger: l,
			index:  -1,
			Ctx:    ctx,
			Time:   now,
			Color:  defaultLevelColor[level],
			Level:  level,
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
			callerFnName, path, line := gdebug.CallerWithFilter(pathFilterKey, l.config.StSkip)
			if l.config.Flags&F_CALLER_FN > 0 {
				input.CallerFunc = fmt.Sprintf(`[%s]`, callerFnName)
			}
			if l.config.Flags&F_FILE_LONG > 0 {
				input.CallerPath = fmt.Sprintf(`%s:%d:`, path, line)
			}
			if l.config.Flags&F_FILE_SHORT > 0 {
				input.CallerPath = fmt.Sprintf(`%s:%d:`, gfile.Basename(path), line)
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
			input.CtxStr = "{" + traceId.String() + "}"
		}
		// Context values.
		if len(l.config.CtxKeys) > 0 {
			ctxStr := ""
			for _, key := range l.config.CtxKeys {
				if v := ctx.Value(key); v != nil {
					if ctxStr != "" {
						ctxStr += ", "
					}
					ctxStr += fmt.Sprintf("%s: %+v", key, v)
				}
			}
			if ctxStr != "" {
				input.CtxStr += "{" + ctxStr + "}"
			}
		}
	}
	var tempStr string
	for _, v := range values {
		tempStr = gconv.String(v)
		if len(input.Content) > 0 {
			if input.Content[len(input.Content)-1] == '\n' {
				// Remove one blank line(\n\n).
				if tempStr[0] == '\n' {
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
		err := asyncPool.Add(func() {
			input.Next()
		})
		if err != nil {
			intlog.Error(ctx, err)
		}
	} else {
		input.Next()
	}
}

// doPrint outputs the logging content according configuration.
func (l *Logger) doPrint(ctx context.Context, input *HandlerInput) {
	if l.config.Writer == nil {
		// Output content to disk file.
		if l.config.Path != "" {
			l.printToFile(ctx, input.Time, input)
		}
		// Allow output to stdout?
		if l.config.StdoutPrint {
			l.printToStdout(ctx, input)
		}
	} else {
		// Output to custom writer.
		l.printToWriter(ctx, input)
	}
}

// printToWriter writes buffer to writer.
func (l *Logger) printToWriter(ctx context.Context, input *HandlerInput) {
	if l.config.Writer != nil {
		var (
			buffer = input.getBuffer(l.config.WriterColorEnable)
		)
		if _, err := l.config.Writer.Write(buffer.Bytes()); err != nil {
			intlog.Error(ctx, err)
		}
	}
}

// printToStdout outputs logging content to stdout.
func (l *Logger) printToStdout(ctx context.Context, input *HandlerInput) {
	if l.config.StdoutPrint {
		// This will lose color in Windows os system.
		// if _, err := os.Stdout.Write(input.getBuffer(true).Bytes()); err != nil {
		// This will print color in Windows os system.
		if _, err := fmt.Fprintf(color.Output, input.getBuffer(true).String()); err != nil {
			intlog.Error(ctx, err)
		}
	}
}

// printToFile outputs logging content to disk file.
func (l *Logger) printToFile(ctx context.Context, t time.Time, input *HandlerInput) {
	var (
		buffer        = input.getBuffer(l.config.WriterColorEnable)
		logFilePath   = l.getFilePath(t)
		memoryLockKey = memoryLockPrefixForPrintingToFile + logFilePath
	)
	gmlock.Lock(memoryLockKey)
	defer gmlock.Unlock(memoryLockKey)

	// Rotation file size checks.
	if l.config.RotateSize > 0 {
		if gfile.Size(logFilePath) > l.config.RotateSize {
			l.rotateFileBySize(t)
		}
	}
	// Logging content outputting to disk file.
	if file := l.getFilePointer(ctx, logFilePath); file == nil {
		intlog.Errorf(ctx, `got nil file pointer for: %s`, logFilePath)
	} else {
		if _, err := file.Write(buffer.Bytes()); err != nil {
			intlog.Error(ctx, err)
		}
		if err := file.Close(); err != nil {
			intlog.Error(ctx, err)
		}
	}
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
		intlog.Error(ctx, err)
	}
	return file
}

// getCtx returns the context which is set through chaining operations.
// It returns an empty context if no context set previously.
func (l *Logger) getCtx() context.Context {
	if l.ctx != nil {
		return l.ctx
	}
	return context.TODO()
}

// printStd prints content <s> without stack.
func (l *Logger) printStd(level int, value ...interface{}) {
	l.print(l.getCtx(), level, value...)
}

// printStd prints content <s> with stack check.
func (l *Logger) printErr(level int, value ...interface{}) {
	if l.config.StStatus == 1 {
		if s := l.GetStack(); s != "" {
			value = append(value, "\nStack:\n"+s)
		}
	}
	// In matter of sequence, do not use stderr here, but use the same stdout.
	l.print(l.getCtx(), level, value...)
}

// format formats <values> using fmt.Sprintf.
func (l *Logger) format(format string, value ...interface{}) string {
	return fmt.Sprintf(format, value...)
}

// PrintStack prints the caller stack,
// the optional parameter <skip> specify the skipped stack offset from the end point.
func (l *Logger) PrintStack(skip ...int) {
	if s := l.GetStack(skip...); s != "" {
		l.Println("Stack:\n" + s)
	} else {
		l.Println()
	}
}

// GetStack returns the caller stack content,
// the optional parameter <skip> specify the skipped stack offset from the end point.
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
