// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/gogf/gf/g/internal/debug"

	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfpool"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/util/gconv"
)

type Logger struct {
	parent      *Logger   // Parent logger.
	writer      io.Writer // Customized io.Writer.
	flags       int       // Extra flags for logging output features.
	path        string    // Logging directory path.
	file        string    // Format for logging file.
	level       int       // Output level.
	prefix      string    // Prefix string for every logging content.
	stSkip      int       // Skip count for stack.
	stStatus    int       // Stack status(1: enabled - default; 0: disabled)
	headerPrint bool      // Print header or not(true in default).
	stdoutPrint bool      // Output to stdout or not(true in default).
}

const (
	gDEFAULT_FILE_FORMAT     = `{Y-m-d}.log`
	gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE | os.O_WRONLY | os.O_APPEND
	gDEFAULT_FPOOL_PERM      = os.FileMode(0666)
	gDEFAULT_FPOOL_EXPIRE    = 60000
	gPATH_FILTER_KEY         = "/g/os/glog/glog"
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
		file:        gDEFAULT_FILE_FORMAT,
		flags:       F_TIME_STD,
		level:       LEVEL_ALL,
		stStatus:    1,
		headerPrint: true,
		stdoutPrint: true,
	}
	return logger
}

// Clone returns a new logger, which is the clone the current logger.
func (l *Logger) Clone() *Logger {
	logger := Logger{}
	logger = *l
	logger.parent = l
	return &logger
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level int) {
	l.level = level
}

// GetLevel returns the logging level value.
func (l *Logger) GetLevel() int {
	return l.level
}

// SetDebug enables/disables the debug level for logger.
// The debug level is enabled in default.
func (l *Logger) SetDebug(debug bool) {
	if debug {
		l.level = l.level | LEVEL_DEBU
	} else {
		l.level = l.level & ^LEVEL_DEBU
	}
}

// SetAsync enables/disables async logging output feature.
func (l *Logger) SetAsync(enabled bool) {
	if enabled {
		l.flags = l.flags | F_ASYNC
	} else {
		l.flags = l.flags & ^F_ASYNC
	}
}

// SetFlags sets extra flags for logging output features.
func (l *Logger) SetFlags(flags int) {
	l.flags = flags
}

// GetFlags returns the flags of logger.
func (l *Logger) GetFlags() int {
	return l.flags
}

// SetStack enables/disables the stack feature in failure logging outputs.
func (l *Logger) SetStack(enabled bool) {
	if enabled {
		l.stStatus = 1
	} else {
		l.stStatus = 0
	}
}

// SetStackSkip sets the stack offset from the end point.
func (l *Logger) SetStackSkip(skip int) {
	l.stSkip = skip
}

// SetWriter sets the customized logging <writer> for logging.
// The <writer> object should implements the io.Writer interface.
// Developer can use customized logging <writer> to redirect logging output to another service,
// eg: kafka, mysql, mongodb, etc.
func (l *Logger) SetWriter(writer io.Writer) {
	l.writer = writer
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns nil if no writer previously set.
func (l *Logger) GetWriter() io.Writer {
	return l.writer
}

// getFilePointer returns the file pinter for file logging.
// It returns nil if file logging is disabled, or file opening fails.
func (l *Logger) getFilePointer() *gfpool.File {
	if path := l.path; path != "" {
		// Content containing "{}" in the file name is formatted using gtime
		file, _ := gregex.ReplaceStringFunc(`{.+?}`, l.file, func(s string) string {
			return gtime.Now().Format(strings.Trim(s, "{}"))
		})
		// Create path if it does not exist。
		if !gfile.Exists(path) {
			if err := gfile.Mkdir(path); err != nil {
				fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
				return nil
			}
		}
		if fp, err := gfpool.Open(
			path+gfile.Separator+file,
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

// SetPath sets the directory path for file logging.
func (l *Logger) SetPath(path string) error {
	if path == "" {
		return errors.New("path is empty")
	}
	if !gfile.Exists(path) {
		if err := gfile.Mkdir(path); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
			return err
		}
	}
	l.path = strings.TrimRight(path, gfile.Separator)
	return nil
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func (l *Logger) GetPath() string {
	return l.path
}

// SetFile sets the file name <pattern> for file logging.
// Datetime pattern can be used in <pattern>, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func (l *Logger) SetFile(pattern string) {
	l.file = pattern
}

// SetStdoutPrint sets whether output the logging contents to stdout, which is true in default.
func (l *Logger) SetStdoutPrint(enabled bool) {
	l.stdoutPrint = enabled
}

// SetHeaderPrint sets whether output header of the logging contents, which is true in default.
func (l *Logger) SetHeaderPrint(enabled bool) {
	l.headerPrint = enabled
}

// SetPrefix sets prefix string for every logging content.
// Prefix is part of header, which means if header output is shut, no prefix will be output.
func (l *Logger) SetPrefix(prefix string) {
	l.prefix = prefix
}

// print prints <s> to defined writer, logging file or passed <std>.
func (l *Logger) print(std io.Writer, lead string, value ...interface{}) {
	buffer := bytes.NewBuffer(nil)
	if l.headerPrint {
		// Time.
		timeFormat := ""
		if l.flags&F_TIME_DATE > 0 {
			timeFormat += "2006-01-02 "
		}
		if l.flags&F_TIME_TIME > 0 {
			timeFormat += "15:04:05 "
		}
		if l.flags&F_TIME_MILLI > 0 {
			timeFormat += "15:04:05.000 "
		}
		if len(timeFormat) > 0 {
			buffer.WriteString(time.Now().Format(timeFormat))
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
		if l.flags&F_FILE_LONG > 0 {
			callerPath = debug.CallerWithFilter(gPATH_FILTER_KEY, l.stSkip) + ": "
		}
		if l.flags&F_FILE_SHORT > 0 {
			callerPath = gfile.Basename(debug.CallerWithFilter(gPATH_FILTER_KEY, l.stSkip)) + ": "
		}
		if len(callerPath) > 0 {
			buffer.WriteString(callerPath)
		}
		// Prefix.
		if len(l.prefix) > 0 {
			buffer.WriteString(l.prefix + " ")
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
	if l.flags&F_ASYNC > 0 {
		asyncPool.Add(func() {
			l.printToWriter(std, buffer)
		})
	} else {
		l.printToWriter(std, buffer)
	}
}

// printToWriter writes buffer to writer.
func (l *Logger) printToWriter(std io.Writer, buffer *bytes.Buffer) {
	if l.writer == nil {
		if f := l.getFilePointer(); f != nil {
			defer f.Close()
			if _, err := io.WriteString(f, buffer.String()); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
		// Allow output to stdout?
		if l.stdoutPrint {
			if _, err := std.Write(buffer.Bytes()); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	} else {
		if _, err := l.writer.Write(buffer.Bytes()); err != nil {
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
	if l.stStatus == 1 {
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
	number := 1
	if len(skip) > 0 {
		number = skip[0] + 1
	}
	return debug.StackWithFilter(gPATH_FILTER_KEY, number)
}
