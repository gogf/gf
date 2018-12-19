// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
// @author john, zseeker

package glog

import (
    "fmt"
    "gitee.com/johng/gf/g/container/gtype"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfpool"
    "gitee.com/johng/gf/g/os/gmlock"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gregex"
    "io"
    "os"
    "runtime"
    "strings"
    "sync"
    "time"
)

type Logger struct {
    mu           sync.RWMutex
    pr           *Logger             // 父级Logger
    io           io.Writer           // 日志内容写入的IO接口
    path         *gtype.String       // 日志写入的目录路径
    file         *gtype.String       // 日志文件名称格式
    level        *gtype.Int          // 日志输出等级
    btSkip       *gtype.Int          // 错误产生时的backtrace回调信息skip条数
    btStatus     *gtype.Int          // 是否当打印错误时同时开启backtrace打印(默认-1，表示默认打印逻辑 - 错误才打印)
    printHeader  *gtype.Bool         // 是否不打印前缀信息(时间，级别等)
    alsoStdPrint *gtype.Bool         // 控制台打印开关，当输出到文件/自定义输出时也同时打印到终端
}

const (
    gDEFAULT_FILE_FORMAT     = `{Y-m-d}.log`
    gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE|os.O_WRONLY|os.O_APPEND
    gDEFAULT_FPOOL_PERM      = os.FileMode(0666)
    gDEFAULT_FPOOL_EXPIRE    = 60000
)

var (
    // 默认的日志换行符
    ln    = "\n"
    // 标准输出互斥锁，防止标准输出串日志
    stdMu = sync.RWMutex{}
)

// 初始化日志换行符
func init() {
    if runtime.GOOS == "windows" {
        ln = "\r\n"
    }
}

// 新建自定义的日志操作对象
func New() *Logger {
    return &Logger {
        io           : nil,
        path         : gtype.NewString(),
        file         : gtype.NewString(gDEFAULT_FILE_FORMAT),
        level        : gtype.NewInt(defaultLevel.Val()),
        btSkip       : gtype.NewInt(),
        btStatus     : gtype.NewInt(-1),
        printHeader  : gtype.NewBool(true),
        alsoStdPrint : gtype.NewBool(true),
    }
}

// Logger深拷贝
func (l *Logger) Clone() *Logger {
    return &Logger {
        pr           : l,
        io           : l.GetWriter(),
        path         : l.path.Clone(),
        file         : l.file.Clone(),
        level        : l.level.Clone(),
        btSkip       : l.btSkip.Clone(),
        btStatus    : l.btStatus.Clone(),
        printHeader  : l.printHeader.Clone(),
        alsoStdPrint : l.alsoStdPrint.Clone(),
    }
}

// 设置日志记录等级
func (l *Logger) SetLevel(level int) {
    l.level.Set(level)
}

// 获取日志记录等级
func (l *Logger) GetLevel() int {
    return l.level.Val()
}

// 快捷方法，打开或关闭DEBU日志信息
func (l *Logger) SetDebug(debug bool) {
    if debug {
        l.level.Set(l.level.Val() | LEVEL_DEBU)
    } else {
        l.level.Set(l.level.Val() & ^LEVEL_DEBU)
    }
}

func (l *Logger) SetBacktrace(enabled bool) {
    if enabled {
        l.btStatus.Set(1)
    } else {
        l.btStatus.Set(0)
    }

}

// 设置BacktraceSkip
func (l *Logger) SetBacktraceSkip(skip int) {
    l.btSkip.Set(skip)
}

// 可自定义IO接口，IO可以是文件输出、标准输出、网络输出
func (l *Logger) SetWriter(writer io.Writer) {
    l.mu.Lock()
    l.io = writer
    l.mu.Unlock()
}

// 返回自定义的IO，默认为nil
func (l *Logger) GetWriter() io.Writer {
    l.mu.RLock()
    r := l.io
    l.mu.RUnlock()
    return r
}

// 获取默认的文件IO
func (l *Logger) getFilePointer() *gfpool.File {
    if path := l.path.Val(); path != "" {
        // 文件名称中使用"{}"包含的内容使用gtime格式化
        file, _ := gregex.ReplaceStringFunc(`{.+?}`, l.file.Val(), func(s string) string {
            return gtime.Now().Format(strings.Trim(s, "{}"))
        })
        // 如果日志目录不存在则创建目录路径
        if !gfile.Exists(path) {
            if err := gfile.Mkdir(path); err != nil {
                fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
                return nil
            }
        }
        fpath   := path + gfile.Separator + file
        if fp, err := gfpool.Open(fpath, gDEFAULT_FILE_POOL_FLAGS, gDEFAULT_FPOOL_PERM, gDEFAULT_FPOOL_EXPIRE); err == nil {
            return fp
        } else {
            fmt.Fprintln(os.Stderr, err)
        }
    }
    return nil
}

// 设置日志文件的存储目录路径
func (l *Logger) SetPath(path string) error {
    // 如果目录不存在，则递归创建
    if !gfile.Exists(path) {
       if err := gfile.Mkdir(path); err != nil {
           fmt.Fprintln(os.Stderr, fmt.Sprintf(`[glog] mkdir "%s" failed: %s`, path, err.Error()))
           return err
       }
    }
    l.path.Set(strings.TrimRight(path, gfile.Separator))
    return nil
}

// 获取设置的日志目录路径
func (l *Logger) GetPath() string {
    return l.path.Val()
}

// 日志文件名称
func (l *Logger) SetFile(file string) {
    l.file.Set(file)
}

// 设置写日志时开启or关闭控制台打印，默认是关闭的
func (l *Logger) SetStdPrint(enabled bool) {
    l.alsoStdPrint.Set(enabled)
}

// 这里的写锁保证统一时刻只会写入一行日志，防止串日志的情况
func (l *Logger) print(std io.Writer, s string) {
    // 优先使用自定义的IO输出
    if l.printHeader.Val() {
        s = l.format(s)
    }
    writer := l.GetWriter()
    if writer == nil {
        // 如果设置的writer为空，那么其次判断是否有文件输出设置
        // 内部使用了内存锁，保证在glog中对同一个日志文件的并发写入不会串日志(并发安全)
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
        // 当没有设置writer时，需要判断是否允许输出到标准输出
        if l.alsoStdPrint.Val() {
            l.doStdLockPrint(std, s)
        }
    } else {
        l.doStdLockPrint(writer, s)
    }
}

// 并发安全打印到标准输出
func (l *Logger) doStdLockPrint(std io.Writer, s string) {
    stdMu.Lock()
    if _, err := std.Write([]byte(s)); err != nil {
        fmt.Fprintln(os.Stderr, err.Error())
    }
    stdMu.Unlock()
}

// 核心打印数据方法(标准输出)
func (l *Logger) stdPrint(s string) {
    l.print(os.Stdout, s)
}

// 核心打印数据方法(标准错误)
func (l *Logger) errPrint(s string) {
    // 记录调用回溯信息
    status := l.btStatus.Val()
    if status == -1 || status == 1 {
        s = l.appendBacktrace(s)
    }
    // 防止串日志情况，这里不使用stderr，而是使用stdout
    l.print(os.Stdout, s)
}

// 输出内容中添加回溯信息
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

// 直接打印回溯信息，参数skip表示调用端往上多少级开始回溯
func (l *Logger) PrintBacktrace(skip...int) {
    l.Println(l.appendBacktrace("", skip...))
}

// 获取文件调用回溯字符串，参数skip表示调用端往上多少级开始回溯
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

func (l *Logger) Fatal(v ...interface{}) {
    l.errPrint("[FATA] " + fmt.Sprintln(v...))
    os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.errPrint("[FATA] " + fmt.Sprintf(format, v...))
    os.Exit(1)
}

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

// 判断给定level是否满足
func (l *Logger) checkLevel(level int) bool {
    return l.level.Val() & level > 0
}