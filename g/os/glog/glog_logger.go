// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package glog

import (
    "os"
    "io"
    "time"
    "fmt"
    "errors"
    "strings"
    "runtime"
    "strconv"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/util/gregex"
    "gitee.com/johng/gf/g/os/gfilepool"
    "gitee.com/johng/gf/g/container/gtype"
)

const (
    gDEFAULT_FILE_POOL_FLAGS = os.O_CREATE|os.O_WRONLY|os.O_APPEND
)

// 默认的日志换行符
var ln = "\n"

// 初始化日志换行符
// @author zseeker
// @date   2018-05-23
func init() {
    if runtime.GOOS == "windows" {
        ln = "\r\n"
    }
}

// Logger浅拷贝，除了path，其他都是父对象属性的指针
func (l *Logger) Clone() *Logger {
    return &Logger {
        pr       : l,
        io       : l.GetIO(),
        path     : gtype.NewString(l.path.Val()),
        debug    : l.debug,
        btSkip   : l.btSkip,
        stdprint : l.stdprint,
    }
}

// 设置下一次输出的日志分类(可以按照文件目录层级设置)，在当前logpath或者当前工作目录下创建category目录，
// 这是一个链式操作，可以设置多个分类，将会创建层级的日志分类目录。
func (l *Logger) Cat(category string) *Logger {
    path := l.path.Val()
    if path == "" {
        path = gfile.Pwd()
    }
    path += gfile.Separator + category
    if v := loggerMap.Get(path); v != nil {
        return v.(*Logger)
    } else {
        logger := l.Clone()
        logger.SetPath(path)
        loggerMap.Set(path, logger)
        return logger
    }
}

// 设置BacktraceSkip
func (l *Logger) SetBacktraceSkip(skip int) {
    l.btSkip.Set(skip)
}

// 可自定义IO接口，IO可以是文件输出、标准输出、网络输出
func (l *Logger) SetIO(w io.Writer) {
    l.mu.Lock()
    l.io = w
    l.mu.Unlock()
}

// 返回自定义的IO，默认为nil
func (l *Logger) GetIO() io.Writer {
    l.mu.RLock()
    r := l.io
    l.mu.RUnlock()
    if r == nil && l.pr != nil {
        return l.pr.GetIO()
    }
    return r
}

// 获取默认的文件IO
func (l *Logger) getFileByPool() *gfilepool.File {
    if path := l.path.Val(); path != "" {
        fpath := path + gfile.Separator + time.Now().Format("2006-01-02.log")
        if fp, err := gfilepool.OpenWithPool(fpath, gDEFAULT_FILE_POOL_FLAGS, 0666, 86400000); err == nil {
            return fp
        } else {
            fmt.Fprintln(os.Stderr, err)
        }
    }
    return nil
}

func (l *Logger) GetDebug() bool {
    return l.debug.Val()
}

func (l *Logger) SetDebug(debug bool) {
    l.debug.Set(debug)
}

// 设置日志文件的存储目录路径
func (l *Logger) SetPath(path string) error {
    // 检测目录权限
    if !gfile.Exists(path) {
        if err := gfile.Mkdir(path); err != nil {
            fmt.Fprintln(os.Stderr, "glog mkdir failed: " + err.Error())
            return err
        }
    }
    if !gfile.IsWritable(path) {
        errstr := path + " is no writable for current user"
        fmt.Fprintln(os.Stderr, errstr)
        return errors.New(errstr)
    }
    l.path.Set(strings.TrimRight(path, gfile.Separator))
    return nil
}

// 设置写日志时开启or关闭控制台打印，默认是关闭的
// @author zseeker
// @date   2018-05-24
func (l *Logger) SetStdPrint(open bool) {
    l.stdprint.Set(open)
}

// 这里的写锁保证统一时刻只会写入一行日志，防止串日志的情况
func (l *Logger) print(def io.Writer, s string) {
    // 优先使用自定义的IO输出
    writer := l.GetIO()
    if writer == nil {
        // 如果设置的IO为空，那么其次判断是否有文件输出设置
        if f := l.getFileByPool(); f != nil {
            writer = f
            // 如果有文件设置那么需要判断是否同时输出到文件和终端
            // @author zseeker
            if l.stdprint.Val() {
                writer = io.MultiWriter(writer, def)
            }
            defer f.Close()
        } else {
            writer = def
        }
    }
    l.mu.Lock()
    fmt.Fprint(writer, l.format(s))
    l.mu.Unlock()
}

// 核心打印数据方法(标准输出)
func (l *Logger) stdPrint(s string) {
    l.print(os.Stdout, s)
}

// 核心打印数据方法(标准错误)
func (l *Logger) errPrint(s string) {
    // 记录调用回溯信息
    backtrace := l.GetBacktrace(3)
    if s[len(s) - 1] == byte('\n') {
        s = s + backtrace + ln
    } else {
        s = s + ln + backtrace + ln
    }
    l.print(os.Stderr, s)
}

// 直接打印回溯信息，参数skip表示调用端往上多少级开始回溯
func (l *Logger) PrintBacktrace(skip...int) {
    customSkip := 1
    if len(skip) > 0 {
        customSkip += skip[0]
    }
    l.Println(l.GetBacktrace(customSkip))
}

// 获取文件调用回溯字符串，参数skip表示调用端往上多少级开始回溯
func (l *Logger) GetBacktrace(skip...int) string {
    customSkip := 0
    if len(skip) > 0 {
        customSkip += skip[0]
    }
    backtrace := ""
    index     := 1
    for i := 1; i < 10000; i++ {
        if _, cfile, cline, ok := runtime.Caller(customSkip + i + l.btSkip.Val()); ok {
            // 不打印出go源码路径
            if !gregex.IsMatchString("^" + gfile.GoRootOfBuild(), cfile) {
                backtrace += strconv.Itoa(index) + ".\t" + cfile + ":" + strconv.Itoa(cline) + ln
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
    l.stdPrint(fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
    l.stdPrint(fmt.Sprint(v...) + ln)
}

func (l *Logger) Printfln(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format + ln, v...))
}

func (l *Logger) Fatal(v ...interface{}) {
    l.errPrint(fmt.Sprint(v...))
    os.Exit(1)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.errPrint(fmt.Sprintf(format, v...))
    os.Exit(1)
}

func (l *Logger) Fatalln(v ...interface{}) {
    l.errPrint(fmt.Sprint(v...) + ln)
    os.Exit(1)
}

func (l *Logger) Fatalfln(format string, v ...interface{}) {
    l.errPrint(fmt.Sprintf(format + ln, v...))
    os.Exit(1)
}

func (l *Logger) Panic(v ...interface{}) {
    s := fmt.Sprint(v...)
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format, v...)
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Panicln(v ...interface{}) {
    s := fmt.Sprint(v...) + ln
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Panicfln(format string, v ...interface{}) {
    s := fmt.Sprintf(format + ln, v...)
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Info(v ...interface{}) {
    l.stdPrint("[INFO] " + fmt.Sprint(v...) + ln)
}

func (l *Logger) Debug(v ...interface{}) {
    if l.GetDebug() {
        l.stdPrint("[DEBU] " + fmt.Sprint(v...) + ln)
    }
}

func (l *Logger) Notice(v ...interface{}) {
    l.errPrint("[NOTI] " + fmt.Sprint(v...) + ln)
}

func (l *Logger) Warning(v ...interface{}) {
    l.errPrint("[WARN] " + fmt.Sprint(v...) + ln)
}

func (l *Logger) Error(v ...interface{}) {
    l.errPrint("[ERRO] " + fmt.Sprint(v...) + ln)
}

func (l *Logger) Critical(v ...interface{}) {
    l.errPrint("[CRIT] " + fmt.Sprint(v...) + ln)
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.stdPrint("[INFO] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Debugf(format string, v ...interface{}) {
    if l.GetDebug() {
        l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...))
    }
}

func (l *Logger) Noticef(format string, v ...interface{}) {
    l.errPrint("[NOTI] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Warningf(format string, v ...interface{}) {
    l.errPrint("[WARN] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.errPrint("[ERRO] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Criticalf(format string, v ...interface{}) {
    l.errPrint("[CRIT] " + fmt.Sprintf(format, v...))
}

func (l *Logger) Infofln(format string, v ...interface{}) {
    l.stdPrint("[INFO] " + fmt.Sprintf(format, v...) + ln)
}

func (l *Logger) Debugfln(format string, v ...interface{}) {
    if l.GetDebug() {
        l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...) + ln)
    }
}

func (l *Logger) Noticefln(format string, v ...interface{}) {
    l.errPrint("[NOTI] " + fmt.Sprintf(format, v...) + ln)
}

func (l *Logger) Warningfln(format string, v ...interface{}) {
    l.errPrint("[WARN] " + fmt.Sprintf(format, v...) + ln)
}

func (l *Logger) Errorfln(format string, v ...interface{}) {
    l.errPrint("[ERRO] " + fmt.Sprintf(format, v...) + ln)
}

func (l *Logger) Criticalfln(format string, v ...interface{}) {
    l.errPrint("[CRIT] " + fmt.Sprintf(format, v...) + ln)
}