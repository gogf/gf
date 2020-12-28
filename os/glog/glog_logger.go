// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"context"
	"fmt"
	"github.com/gogf/gf/container/gtype"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gfpool"
	"github.com/gogf/gf/os/gmlock"
	"github.com/gogf/gf/os/gtimer"
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
	defaultFileFormat = `{Y-m-d}.log`
	defaultFileFlags  = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	defaultFilePerm   = os.FileMode(0666)
	defaultFileExpire = time.Minute
	pathFilterKey     = "/os/glog/glog"
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
func (l *Logger) print(std io.Writer, lead string, values ...interface{}) {
	// Lazy initialize for rotation feature.
	// It uses atomic reading operation to enhance the performance checking.
	// It here uses CAP for performance and concurrent safety.
	p := l
	if p.parent != nil {
		p = p.parent
	}
	if !p.init.Val() && p.init.Cas(false, true) {
		// It just initializes once for each logger.
		if p.config.RotateSize > 0 || p.config.RotateExpire > 0 {
			gtimer.AddOnce(p.config.RotateCheckInterval, p.rotateChecksTimely)
			intlog.Printf("logger rotation initialized: every %s", p.config.RotateCheckInterval.String())
		}
	}

	var (
		now    = time.Now()
		buffer = bytes.NewBuffer(nil)
	)
	if l.config.HeaderPrint {
		// Time.
		timeFormat := ""
		if l.config.Flags&F_TIME_DATE > 0 {
			timeFormat += "2006-01-02 "
		}
		if l.config.Flags&F_TIME_TIME > 0 {
			timeFormat += "15:04:05 "
		}
		if l.config.Flags&F_TIME_MILLI > 0 {
			timeFormat += "15:04:05.000 "
		}
		if len(timeFormat) > 0 {
			buffer.WriteString(now.Format(timeFormat))
		}
		// Lead string.
		if len(lead) > 0 {
			buffer.WriteString(lead)
			if len(values) > 0 {
				buffer.WriteByte(' ')
			}
		}
		// Caller path and Fn name.
		if l.config.Flags&(F_FILE_LONG|F_FILE_SHORT|F_CALLER_FN) > 0 {
			callerPath := ""
			callerFnName, path, line := gdebug.CallerWithFilter(pathFilterKey, l.config.StSkip)
			if l.config.Flags&F_CALLER_FN > 0 {
				buffer.WriteString(fmt.Sprintf(`[%s] `, callerFnName))
			}
			if l.config.Flags&F_FILE_LONG > 0 {
				callerPath = fmt.Sprintf(`%s:%d: `, path, line)
			}
			if l.config.Flags&F_FILE_SHORT > 0 {
				callerPath = fmt.Sprintf(`%s:%d: `, gfile.Basename(path), line)
			}
			buffer.WriteString(callerPath)

		}
		// Prefix.
		if len(l.config.Prefix) > 0 {
			buffer.WriteString(l.config.Prefix + " ")
		}
	}
	// Convert value to string.
	var (
		tempStr  = ""
		valueStr = ""
	)
	// Context values.
	if l.ctx != nil && len(l.config.CtxKeys) > 0 {
		ctxStr := ""
		for _, key := range l.config.CtxKeys {
			if v := l.ctx.Value(key); v != nil {
				if ctxStr != "" {
					ctxStr += ", "
				}
				ctxStr += fmt.Sprintf("%s: %+v", key, v)
			}
		}
		if ctxStr != "" {
			buffer.WriteString(fmt.Sprintf("{%s} ", ctxStr))
		}
	}
	for _, v := range values {
		tempStr = gconv.String(v)
		if len(valueStr) > 0 {
			if valueStr[len(valueStr)-1] == '\n' {
				// Remove one blank line(\n\n).
				if tempStr[0] == '\n' {
					valueStr += tempStr[1:]
				} else {
					valueStr += tempStr
				}
			} else {
				valueStr += " " + tempStr
			}
		} else {
			valueStr = tempStr
		}
	}
	buffer.WriteString(valueStr + "\n")
	if l.config.Flags&F_ASYNC > 0 {
		err := asyncPool.Add(func() {
			l.printToWriter(now, std, buffer)
		})
		if err != nil {
			intlog.Error(err)
		}
	} else {
		l.printToWriter(now, std, buffer)
	}
}

// printToWriter writes buffer to writer.
func (l *Logger) printToWriter(now time.Time, std io.Writer, buffer *bytes.Buffer) {
	if l.config.Writer == nil {
		// Output content to disk file.
		if l.config.Path != "" {
			l.printToFile(now, buffer)
		}
		// Allow output to stdout?
		if l.config.StdoutPrint {
			if _, err := std.Write(buffer.Bytes()); err != nil {
				intlog.Error(err)
			}
		}
	} else {
		if _, err := l.config.Writer.Write(buffer.Bytes()); err != nil {
			// panic(err)
			intlog.Error(err)
		}
	}
}

// printToFile outputs logging content to disk file.
func (l *Logger) printToFile(now time.Time, buffer *bytes.Buffer) {
	var (
		logFilePath   = l.getFilePath(now)
		memoryLockKey = "glog.printToFile:" + logFilePath
	)
	gmlock.Lock(memoryLockKey)
	defer gmlock.Unlock(memoryLockKey)

	// Rotation file size checks.
	if l.config.RotateSize > 0 {
		if gfile.Size(logFilePath) > l.config.RotateSize {
			l.rotateFileBySize(now)
		}
	}
	// Logging content outputting to disk file.
	if file := l.getFilePointer(logFilePath); file == nil {
		intlog.Errorf(`got nil file pointer for: %s`, logFilePath)
	} else {
		if _, err := file.Write(buffer.Bytes()); err != nil {
			intlog.Error(err)
		}
		if err := file.Close(); err != nil {
			intlog.Error(err)
		}
	}
}

// getFilePointer retrieves and returns a file pointer from file pool.
func (l *Logger) getFilePointer(path string) *gfpool.File {
	file, err := gfpool.Open(
		path,
		defaultFileFlags,
		defaultFilePerm,
		defaultFileExpire,
	)
	if err != nil {
		// panic(err)
		intlog.Error(err)
	}
	return file
}

// printStd prints content <s> without stack.
func (l *Logger) printStd(lead string, value ...interface{}) {
	l.print(os.Stdout, lead, value...)
}

// printStd prints content <s> with stack check.
func (l *Logger) printErr(lead string, value ...interface{}) {
	if l.config.StStatus == 1 {
		if s := l.GetStack(); s != "" {
			value = append(value, "\nStack:\n"+s)
		}
	}
	// In matter of sequence, do not use stderr here, but use the same stdout.
	l.print(os.Stdout, lead, value...)
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
