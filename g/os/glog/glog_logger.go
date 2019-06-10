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
	"github.com/gogf/gf/g/os/gfile"
	"github.com/gogf/gf/g/os/gfpool"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/text/gregex"
	"github.com/gogf/gf/g/util/gconv"
	"io"
	"os"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
    parent       *Logger      // Parent logger.
	writer       io.Writer    // Customized io.Writer.
	flags        int          // Extra flags for logging output features.
    path         string       // Logging directory path.
    file         string       // Format for logging file.
    level        int          // Output level.
    prefix       string       // Prefix string for every logging content.
    btSkip       int          // Skip count for backtrace.
    btStatus     int          // Backtrace status(1: enabled - default; 0: disabled)
    headerPrint  bool         // Print header or not(true in default).
    stdoutPrint  bool         // Output to stdout or not(true in default).
}

const (
    gDEFAULT_FILE_FORMAT     = `{Y-m-d}.log`
    gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE|os.O_WRONLY|os.O_APPEND
    gDEFAULT_FPOOL_PERM      = os.FileMode(0666)
    gDEFAULT_FPOOL_EXPIRE    = 60000
)

const (
	F_ASYNC      = 1 << iota // Print logging content asynchronously。
	F_FILE_LONG              // Print full file name and line number: /a/b/c/d.go:23.
	F_FILE_SHORT             // Print final file name element and line number: d.go:23. overrides F_FILE_LONG.
	F_TIME_DATE              // Print the date in the local time zone: 2009-01-23.
	F_TIME_TIME              // Print the time in the local time zone: 01:23:23.
	F_TIME_MILLI             // Print the time with milliseconds in the local time zone: 01:23:23.675.
	F_TIME_STD = F_TIME_DATE | F_TIME_MILLI
)

var (
    // Default line break.
    ln = "\n"
)

func init() {
	// Initialize log line breaks depending on underlying os.
    if runtime.GOOS == "windows" {
        ln = "\r\n"
    }
}

// New creates and returns a custom logger.
func New() *Logger {
    logger := &Logger {
        file         : gDEFAULT_FILE_FORMAT,
        flags        : F_TIME_STD,
        level        : LEVEL_ALL,
        btStatus     : 1,
        headerPrint  : true,
        stdoutPrint  : true,
    }
    return logger
}

// Clone returns a new logger, which is the clone the current logger.
func (l *Logger) Clone() *Logger {
	logger := Logger{}
	logger  = *l
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
        l.level = l.level  & ^LEVEL_DEBU
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

// SetBacktrace enables/disables the backtrace feature in failure logging outputs.
func (l *Logger) SetBacktrace(enabled bool) {
    if enabled {
        l.btStatus = 1
    } else {
        l.btStatus = 0
    }
}

// SetBacktraceSkip sets the backtrace offset from the end point.
func (l *Logger) SetBacktraceSkip(skip int) {
    l.btSkip = skip
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
        	path + gfile.Separator + file,
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
func (l *Logger) print(std io.Writer, lead string, value...interface{}) {
	buffer := bytes.NewBuffer(nil)
    if l.headerPrint {
	    // Time.
	    timeFormat := ""
	    if l.flags & F_TIME_DATE > 0 {
		    timeFormat += "2006-01-02 "
	    }
	    if l.flags & F_TIME_TIME > 0 {
		    timeFormat += "15:04:05 "
	    }
	    if l.flags & F_TIME_MILLI > 0 {
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
	    if l.flags & F_FILE_LONG > 0 {
		    callerPath = l.getLongFile() + ": "
	    }
	    if l.flags & F_FILE_SHORT > 0 {
		    callerPath = gfile.Basename(l.getLongFile()) + ": "
	    }
	    if len(callerPath) > 0 {
		    buffer.WriteString(callerPath)
	    }
	    // Prefix.
	    if len(l.prefix) > 0 {
		    buffer.WriteString(l.prefix + " ")
	    }
    }
	for k, v := range value {
		if k > 0 {
			buffer.WriteByte(' ')
		}
		buffer.WriteString(gconv.String(v))
	}
	buffer.WriteString(ln)
	if l.flags & F_ASYNC > 0 {
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

// printStd prints content <s> without backtrace.
func (l *Logger) printStd(lead string, value...interface{}) {
    l.print(os.Stdout, lead, value...)
}

// printStd prints content <s> with backtrace check.
func (l *Logger) printErr(lead string, value...interface{}) {
    if l.btStatus == 1 {
    	if s := l.GetBacktrace(); s != "" {
		    value = append(value, ln + "Backtrace:" + ln + s)
	    }
    }
    // In matter of sequence, do not use stderr here, but use the same stdout.
    l.print(os.Stdout, lead, value...)
}

// format formats <values> using fmt.Sprintf.
func (l *Logger) format(format string, value...interface{}) string {
	return fmt.Sprintf(format, value...)
}

// PrintBacktrace prints the caller backtrace, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func (l *Logger) PrintBacktrace(skip...int) {
	if s := l.GetBacktrace(skip...); s != "" {
		l.Println("Backtrace:" + ln + s)
	} else {
		l.Println()
	}
}

// GetBacktrace returns the caller backtrace content, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func (l *Logger) GetBacktrace(skip...int) string {
    customSkip := 0
    if len(skip) > 0 {
        customSkip = skip[0]
    }
    backtrace := ""
    from      := 0
    // Find the caller position exclusive of the glog file.
    for i := 0; i < 1000; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
            if !gregex.IsMatchString("/g/os/glog/glog.+$", file) {
                from = i
                break
            }
        }
    }
    // Find the true caller file path using custom skip.
	index  := 1
    goRoot := runtime.GOROOT()
    for i := from + customSkip + l.btSkip; i < 1000; i++ {
        if _, file, cline, ok := runtime.Caller(i); ok && len(file) > 2 {
            if (goRoot == "" || !gregex.IsMatchString("^" + goRoot, file)) && !gregex.IsMatchString(`<autogenerated>`, file) {
                backtrace += fmt.Sprintf(`%d. %s:%d%s`, index, file, cline, ln)
                index++
            }
        } else {
            break
        }
    }
    return backtrace
}

// getLongFile returns the absolute file path along with its line number of the caller.
func (l *Logger) getLongFile() string {
	from := 0
	// Find the caller position exclusive of the glog file.
	for i := 0; i < 1000; i++ {
		if _, file, _, ok := runtime.Caller(i); ok {
			if !gregex.IsMatchString("/g/os/glog/glog.+$", file) {
				from = i
				break
			}
		}
	}
	// Find the true caller file path using custom skip.
	goRoot := runtime.GOROOT()
	for i := from + l.btSkip; i < 1000; i++ {
		if _, file, line, ok := runtime.Caller(i); ok && len(file) > 2 {
			if (goRoot == "" || !gregex.IsMatchString("^" + goRoot, file)) && !gregex.IsMatchString(`<autogenerated>`, file) {
				return fmt.Sprintf(`%s:%d`, file, line)
			}
		} else {
			break
		}
	}
	return ""
}
