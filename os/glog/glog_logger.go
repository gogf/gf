// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/os/gtimer"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gogf/gf/debug/gdebug"

	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gfpool"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/util/gconv"
)

// Logger is the struct for logging management.
type Logger struct {
	mu     sync.Mutex // Mutex is not for common logging, but for file rotation feature.
	parent *Logger    // Parent logger.
	config Config     // Logger configuration.
}

const (
	gDEFAULT_FILE_FORMAT     = `{Y-m-d}.log`
	gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	gDEFAULT_FPOOL_PERM      = os.FileMode(0666)
	gDEFAULT_FPOOL_EXPIRE    = time.Minute
	gPATH_FILTER_KEY         = "/os/glog/glog"
)

const (
	F_ASYNC      = 1 << iota // Print logging content asynchronously。
	F_FILE_LONG              // Print full file name and line number: /a/b/c/d.go:23.
	F_FILE_SHORT             // Print final file name element and line number: d.go:23. overrides F_FILE_LONG.
	F_TIME_DATE              // Print the date in the local time zone: 2009-01-23.
	F_TIME_TIME              // Print the time in the local time zone: 01:23:23.
	F_TIME_MILLI             // Print the time with milliseconds in the local time zone: 01:23:23.675.
	F_TIME_STD   = F_TIME_DATE | F_TIME_MILLI
)

// New creates and returns a custom logger.
func New() *Logger {
	logger := &Logger{
		config: DefaultConfig(),
	}
	gtimer.AddOnce(time.Second, logger.rotateChecks)
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
	logger := Logger{}
	logger = *l
	logger.parent = l
	return &logger
}

// getFilePointer returns the file pinter for file logging.
// It returns nil if file logging is disabled, or file opening fails.
func (l *Logger) getFilePointer(now time.Time) *gfpool.File {
	if path := l.config.Path; path != "" {
		// Create path if it does not exist.
		if !gfile.Exists(path) {
			if err := gfile.Mkdir(path); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
				return nil
			}
		}
		if fp, err := gfpool.Open(
			l.getFilePath(now),
			gDEFAULT_FILE_POOL_FLAGS,
			gDEFAULT_FPOOL_PERM,
			gDEFAULT_FPOOL_EXPIRE); err == nil {
			return fp
		} else {
			fmt.Fprintln(os.Stderr, err)
		}
	}
	return nil
}

// getFilePath returns the logging file path.
func (l *Logger) getFilePath(now time.Time) string {
	// Content containing "{}" in the file name is formatted using gtime.
	file, _ := gregex.ReplaceStringFunc(`{.+?}`, l.config.File, func(s string) string {
		return gtime.New(now).Format(strings.Trim(s, "{}"))
	})
	return gfile.Join(l.config.Path, file)
}

// print prints <s> to defined writer, logging file or passed <std>.
func (l *Logger) print(std io.Writer, lead string, value ...interface{}) {
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
			if len(value) > 0 {
				buffer.WriteByte(' ')
			}
		}
		// Caller path.
		callerPath := ""
		if l.config.Flags&F_FILE_LONG > 0 {
			_, path, line := gdebug.CallerWithFilter(gPATH_FILTER_KEY, l.config.StSkip)
			callerPath = fmt.Sprintf(`%s:%d: `, path, line)
		}
		if l.config.Flags&F_FILE_SHORT > 0 {
			_, path, line := gdebug.CallerWithFilter(gPATH_FILTER_KEY, l.config.StSkip)
			callerPath = fmt.Sprintf(`%s:%d: `, gfile.Basename(path), line)
		}
		if len(callerPath) > 0 {
			buffer.WriteString(callerPath)
		}
		// Prefix.
		if len(l.config.Prefix) > 0 {
			buffer.WriteString(l.config.Prefix + " ")
		}
	}
	// Convert value to string.
	tempStr := ""
	valueStr := ""
	for _, v := range value {
		if err, ok := v.(error); ok {
			tempStr = fmt.Sprintf("%+v", err)
		} else {
			tempStr = gconv.String(v)
		}
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
		if f := l.getFilePointer(now); f != nil {
			defer f.Close()
			// Rotation file size checks.
			if l.config.RotateSize > 0 {
				state, err := f.Stat()
				if err != nil {
					panic(err)
				}
				if state.Size() > l.config.RotateSize {
					l.rotateFile(now)
					l.printToWriter(now, std, buffer)
					return
				}
			}
			if _, err := io.WriteString(f, buffer.String()); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
		// Allow output to stdout?
		if l.config.StdoutPrint {
			if _, err := std.Write(buffer.Bytes()); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	} else {
		if _, err := l.config.Writer.Write(buffer.Bytes()); err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
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
	filters := []string{gPATH_FILTER_KEY}
	if l.config.StFilter != "" {
		filters = append(filters, l.config.StFilter)
	}
	return gdebug.StackWithFilters(filters, stackSkip)
}
