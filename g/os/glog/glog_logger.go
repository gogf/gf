// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// @author john, zseeker

package glog

import (
    "errors"
    "fmt"
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/os/gfile"
    "github.com/gogf/gf/g/os/gfpool"
    "github.com/gogf/gf/g/os/gmlock"
    "github.com/gogf/gf/g/os/gtime"
    "github.com/gogf/gf/g/text/gregex"
    "io"
    "os"
    "runtime"
    "strings"
    "sync"
    "time"
)

type Logger struct {
    mu           sync.RWMutex
    pr           *Logger             // Parent logger.
	writer       io.Writer           // Customized io.Writer.
    path         *gtype.String       // Logging directory path.
    file         *gtype.String       // Format for logging file.
    level        *gtype.Int          // Output level.
    btSkip       *gtype.Int          // Skip count for backtrace.
    btStatus     *gtype.Int          // Backtrace status(1: enabled - default; 0: disabled)
    printHeader  *gtype.Bool         // Print header or not(true in default).
    alsoStdPrint *gtype.Bool         // Output to stdout or not(true in default).
}

const (
    gDEFAULT_FILE_FORMAT     = `{Y-m-d}.log`
    gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE|os.O_WRONLY|os.O_APPEND
    gDEFAULT_FPOOL_PERM      = os.FileMode(0666)
    gDEFAULT_FPOOL_EXPIRE    = 60000
)

var (
    // Default line break.
    ln    = "\n"
    // Mutex to ensure log output sequence.
    stdMu = sync.RWMutex{}
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
        path         : gtype.NewString(),
        file         : gtype.NewString(gDEFAULT_FILE_FORMAT),
        level        : gtype.NewInt(defaultLevel.Val()),
        btSkip       : gtype.NewInt(),
        btStatus     : gtype.NewInt(1),
        printHeader  : gtype.NewBool(true),
        alsoStdPrint : gtype.NewBool(true),
    }
    logger.writer = &Writer {
	    logger : logger,
    }
    return logger
}

// Clone returns a new logger, which is the clone the current logger.
func (l *Logger) Clone() *Logger {
	logger := &Logger {
        pr           : l,
        path         : l.path.Clone(),
        file         : l.file.Clone(),
        level        : l.level.Clone(),
        btSkip       : l.btSkip.Clone(),
        btStatus     : l.btStatus.Clone(),
        printHeader  : l.printHeader.Clone(),
        alsoStdPrint : l.alsoStdPrint.Clone(),
    }
	logger.writer = &Writer {
		logger : logger,
	}
	return logger
}

// SetLevel sets the logging level.
func (l *Logger) SetLevel(level int) {
    l.level.Set(level)
}

// GetLevel returns the logging level value.
func (l *Logger) GetLevel() int {
    return l.level.Val()
}

// SetDebug enables/disables the debug level for logger.
// The debug level is enabled in default.
func (l *Logger) SetDebug(debug bool) {
    if debug {
        l.level.Set(l.level.Val() | LEVEL_DEBU)
    } else {
        l.level.Set(l.level.Val() & ^LEVEL_DEBU)
    }
}

// SetBacktrace enables/disables the backtrace feature in failure logging outputs.
func (l *Logger) SetBacktrace(enabled bool) {
    if enabled {
        l.btStatus.Set(1)
    } else {
        l.btStatus.Set(0)
    }
}

// SetBacktraceSkip sets the backtrace offset from the end point.
func (l *Logger) SetBacktraceSkip(skip int) {
    l.btSkip.Set(skip)
}

// SetWriter sets the customized logging <writer> for logging.
// The <writer> object should implements the io.Writer interface.
// Developer can use customized logging <writer> to redirect logging output to another service,
// eg: kafka, mysql, mongodb, etc.
func (l *Logger) SetWriter(writer io.Writer) {
    l.mu.Lock()
    l.writer = writer
    l.mu.Unlock()
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns a default writer if no customized writer set.
func (l *Logger) GetWriter() io.Writer {
    l.mu.RLock()
    r := l.writer
    l.mu.RUnlock()
    return r
}

// getFilePointer returns the file pinter for file logging.
// It returns nil if file logging is disabled, or file opening fails.
func (l *Logger) getFilePointer() *gfpool.File {
    if path := l.path.Val(); path != "" {
        // Content containing "{}" in the file name is formatted using gtime
        file, _ := gregex.ReplaceStringFunc(`{.+?}`, l.file.Val(), func(s string) string {
            return gtime.Now().Format(strings.Trim(s, "{}"))
        })
        // Create path if it does not exist。
        if !gfile.Exists(path) {
            if err := gfile.Mkdir(path); err != nil {
                fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
                return nil
            }
        }
        fpath := path + gfile.Separator + file
        if fp, err := gfpool.Open(fpath, gDEFAULT_FILE_POOL_FLAGS, gDEFAULT_FPOOL_PERM, gDEFAULT_FPOOL_EXPIRE); err == nil {
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
    l.path.Set(strings.TrimRight(path, gfile.Separator))
    return nil
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func (l *Logger) GetPath() string {
    return l.path.Val()
}

// SetFile sets the file name <pattern> for file logging.
// Datetime pattern can be used in <pattern>, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func (l *Logger) SetFile(pattern string) {
    l.file.Set(pattern)
}

// SetStdPrint sets whether output the logging contents to stdout, which is false in default.
func (l *Logger) SetStdPrint(enabled bool) {
    l.alsoStdPrint.Set(enabled)
}

// print prints <s> to defined writer, logging file or passed <std>.
// It internally uses memory lock for file logging to ensure logging sequence.
func (l *Logger) print(std io.Writer, s string) {
    // Customized writer has the most high priority.
    if l.printHeader.Val() {
        s = l.format(s)
    }
    writer := l.GetWriter()
    if _, ok := writer.(*Writer); ok {
        if f := l.getFilePointer(); f != nil {
            defer f.Close()
            key := l.path.Val()
            gmlock.Lock(key)
            _, err := io.WriteString(f, s)
            gmlock.Unlock(key)
            if err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
            }
        }
        // Also output to stdout?
        if l.alsoStdPrint.Val() {
            l.doStdLockPrint(std, s)
        }
    } else {
        l.doStdLockPrint(writer, s)
    }
}

// doStdLockPrint prints <s> to <std> concurrent-safely.
func (l *Logger) doStdLockPrint(std io.Writer, s string) {
    stdMu.Lock()
    if _, err := std.Write([]byte(s)); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
    }
    stdMu.Unlock()
}

// stdPrint prints content <s> without backtrace.
func (l *Logger) stdPrint(s string) {
    l.print(os.Stdout, s)
}

// stdPrint prints content <s> with backtrace check.
func (l *Logger) errPrint(s string) {
    if l.btStatus.Val() == 1 {
        s = l.appendBacktrace(s)
    }
    // In matter of sequence, do not use stderr here, but use the same stdout.
    l.print(os.Stdout, s)
}

// appendBacktrace appends backtrace to the <s>.
func (l *Logger) appendBacktrace(s string, skip...int) string {
    trace := l.GetBacktrace(skip...)
    if trace != "" {
        backtrace := "Backtrace:" + ln + trace
        if len(s) > 0 {
            if s[len(s)-1] == byte('\n') {
                s = s + backtrace + ln
            } else {
                s = s + ln + backtrace + ln
            }
        } else {
            s = backtrace
        }
    }
    return s
}

// PrintBacktrace prints the caller backtrace, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func (l *Logger) PrintBacktrace(skip...int) {
    l.Println(l.appendBacktrace("", skip...))
}

// GetBacktrace returns the caller backtrace content, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func (l *Logger) GetBacktrace(skip...int) string {
    customSkip := 0
    if len(skip) > 0 {
        customSkip = skip[0]
    }
    backtrace := ""
    index     := 1
    from      := 0
    // 首先定位业务文件开始位置
    for i := 0; i < 10; i++ {
        if _, file, _, ok := runtime.Caller(i); ok {
            if !gregex.IsMatchString("/g/os/glog/glog.+$", file) {
                from = i
                break
            }
        }
    }
    // 从业务文件开始位置根据自定义的skip开始backtrace
    goRoot := runtime.GOROOT()
    for i := from + customSkip + l.btSkip.Val(); i < 10000; i++ {
        if _, file, cline, ok := runtime.Caller(i); ok && file != "" {
            // 不打印出go源码路径及glog包文件路径，日志打印必须从业务源码文件开始，且从glog包文件开始检索
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

func (l *Logger) format(s string) string {
    return time.Now().Format("2006-01-02 15:04:05.000 ") + s
}

func (l *Logger) Print(v ...interface{}) {
    l.stdPrint(fmt.Sprintln(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
    l.stdPrint(fmt.Sprintln(v...))
}

func (l *Logger) Printfln(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format + ln, v...))
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func (l *Logger) Fatal(v ...interface{}) {
    l.errPrint("[FATA] " + fmt.Sprintln(v...))
    os.Exit(1)
}

// Fatalf prints the logging content with [FATA] header and custom format, then exit the current process.
func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.errPrint("[FATA] " + fmt.Sprintf(format, v...))
    os.Exit(1)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func (l *Logger) Fatalfln(format string, v ...interface{}) {
    l.errPrint("[FATA] " + fmt.Sprintf(format + ln, v...))
    os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
    s := fmt.Sprintln(v...)
    l.errPrint("[PANI] " + s)
    panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    l.errPrint("[PANI] " + s)
    panic(s)
}

func (l *Logger) Panicfln(format string, v ...interface{}) {
    s := fmt.Sprintf(format + ln, v...)
    l.errPrint("[PANI] " + s)
    panic(s)
}

func (l *Logger) Info(v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.stdPrint("[INFO] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Infof(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.stdPrint("[INFO] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Infofln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.stdPrint("[INFO] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Debug(v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.stdPrint("[DEBU] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Debugf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Debugfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Notice(v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.errPrint("[NOTI] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Noticef(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.errPrint("[NOTI] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Noticefln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.errPrint("[NOTI] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Warning(v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.errPrint("[WARN] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Warningf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.errPrint("[WARN] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Warningfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.errPrint("[WARN] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Error(v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.errPrint("[ERRO] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.errPrint("[ERRO] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Errorfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.errPrint("[ERRO] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Critical(v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.errPrint("[CRIT] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.errPrint("[CRIT] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Criticalfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.errPrint("[CRIT] " + fmt.Sprintf(format, v...) + ln)
    }
}

// checkLevel checks whether the given <level> could be output.
func (l *Logger) checkLevel(level int) bool {
    return l.level.Val() & level > 0
}