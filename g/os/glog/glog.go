// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// 日志模块.
// 直接文件/输出操作，没有异步逻辑，没有使用缓存或者通道
package glog

import (
    "sync"
    "os"
    "io"
    "time"
    "fmt"
    "errors"
    "strings"
    "path/filepath"
    "gitee.com/johng/gf/g/os/gfile"
    "gitee.com/johng/gf/g/os/gfilepool"
    "runtime"
    "strconv"
)

type Logger struct {
    mutex        sync.RWMutex
    logio        io.Writer
    debug        bool         // 是否允许输出DEBUG信息
    logpath      string       // 日志写入的目录路径
    lastlogdate  string       // 上一次写入日志的日期，例如: 2006-01-02
}

// 默认的日志对象
var logger = New()

// 新建自定义的日志操作对象
func New() *Logger {
    return &Logger{
        debug : true,
    }
}

// 日志日志目录绝对路径
func SetLogPath(path string) {
    logger.SetLogPath(path)
}

// 设置是否允许输出DEBUG信息
func SetDebug(debug bool) {
    logger.SetDebug(debug)
}

// 获取日志目录绝对路径
func GetLogPath() string {
    return logger.GetLogPath()
}

func Print(v ...interface{}) {
    logger.Print(v ...)
}

func Printf(format string, v ...interface{}) {
    logger.Printf(format, v ...)
}

func Println(v ...interface{}) {
    logger.Println(v ...)
}

func Printfln(format string, v ...interface{}) {
    logger.Printfln(format, v ...)
}

func Fatal(v ...interface{}) {
    logger.Fatal(v ...)
}

func Fatalf(format string, v ...interface{}) {
    logger.Fatalf(format, v ...)
}

func Fatalln(v ...interface{}) {
    logger.Fatalln(v ...)
}

func Fatalfln(format string, v ...interface{}) {
    logger.Fatalfln(format, v ...)
}

func Panic(v ...interface{}) {
    logger.Panic(v ...)
}

func Panicf(format string, v ...interface{}) {
    logger.Panicf(format, v ...)
}

func Panicln(v ...interface{}) {
    logger.Panicln(v ...)
}

func Panicfln(format string, v ...interface{}) {
    logger.Panicfln(format, v ...)
}

func Info(v ...interface{}) {
    logger.Info(v...)
}

func Debug(v ...interface{}) {
    logger.Debug(v...)
}

func Notice(v ...interface{}) {
    logger.Notice(v...)
}

func Warning(v ...interface{}) {
    logger.Warning(v...)
}

func Error(v ...interface{}) {
    logger.Error(v...)
}

func Critical(v ...interface{}) {
    logger.Critical(v...)
}

func Infof(format string, v ...interface{}) {
    logger.Infof(format, v...)
}

func Debugf(format string, v ...interface{}) {
    logger.Debugf(format, v...)
}

func Noticef(format string, v ...interface{}) {
    logger.Noticef(format, v...)
}

func Warningf(format string, v ...interface{}) {
    logger.Warningf(format, v...)
}

func Errorf(format string, v ...interface{}) {
    logger.Errorf(format, v...)
}

func Criticalf(format string, v ...interface{}) {
    logger.Criticalf(format, v...)
}

func Infofln(format string, v ...interface{}) {
    logger.Infofln(format, v...)
}

func Debugfln(format string, v ...interface{}) {
    logger.Debugfln(format, v...)
}

func Noticefln(format string, v ...interface{}) {
    logger.Noticefln(format, v...)
}

func Warningfln(format string, v ...interface{}) {
    logger.Warningfln(format, v...)
}

func Errorfln(format string, v ...interface{}) {
    logger.Errorfln(format, v...)
}

func Criticalfln(format string, v ...interface{}) {
    logger.Criticalfln(format, v...)
}

func (l *Logger) GetLogIO() io.Writer {
    l.mutex.RLock()
    r := l.logio
    l.mutex.RUnlock()
    return r
}

func (l *Logger) GetDebug() bool {
    l.mutex.RLock()
    r := l.debug
    l.mutex.RUnlock()
    return r
}

func (l *Logger) GetLogPath() string {
    l.mutex.RLock()
    r := l.logpath
    l.mutex.RUnlock()
    return r
}

func (l *Logger) GetLastLogDate() string {
    l.mutex.RLock()
    r := l.lastlogdate
    l.mutex.RUnlock()
    return r
}

func (l *Logger) SetLogIO(w io.Writer) {
    l.mutex.RLock()
    l.logio = w
    l.mutex.RUnlock()
}

func (l *Logger) SetDebug(debug bool) {
    l.mutex.Lock()
    l.debug = debug
    l.mutex.Unlock()
}

// 设置日志文件的存储目录路径
func (l *Logger) SetLogPath(path string) error {
    // 检测目录权限
    if !gfile.Exists(path) {
        if err := gfile.Mkdir(path); err != nil {
            fmt.Fprintln(os.Stderr, err)
            return err
        }
    }
    if !gfile.IsWritable(path) {
        errstr := path + " is no writable for current user"
        fmt.Fprintln(os.Stderr, errstr)
        return errors.New(errstr)
    }
    l.mutex.Lock()
    l.logpath = strings.TrimRight(path, string(filepath.Separator))
    l.mutex.Unlock()
    // 重新检查日志io对象
    l.checkLogIO()
    return nil
}

// 检查文件名称是否已经过期，如果过期那么需要新建一个日志文件(默认按照日期分隔)
func (l *Logger) checkLogIO() {
    date := time.Now().Format("2006-01-02")
    if date != l.GetLastLogDate() {
        if path := l.GetLogPath(); path != "" {
            fname := date + ".log"
            fpath := path + string(filepath.Separator) + fname
            if fp, err := gfilepool.OpenWithPool(fpath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 600); err == nil {
                l.SetLogIO(fp.File())
            } else {
                fmt.Fprintln(os.Stderr, err)
            }
        }
    }
}

// 这里的互斥锁保证统一时刻只会写入一行日志，防止串日志的情况
func (l *Logger) print(logio io.Writer, s string) {
    l.mutex.Lock()
    fmt.Fprint(logio, l.format(s))
    l.mutex.Unlock()
}

// 核心打印数据方法(标准输出)
func (l *Logger) stdPrint(s string) {
    l.checkLogIO()
    logio := l.GetLogIO()
    if logio == nil {
        logio = os.Stdout
    }
    l.print(logio, s)
}

// 核心打印数据方法(标准错误)
func (l *Logger) errPrint(s string) {
    l.checkLogIO()
    logio := l.GetLogIO()
    if logio == nil {
        logio = os.Stderr
    }
    // 记录调用回溯信息
    backtrace := l.backtrace()
    if s[len(s) - 1] == byte('\n') {
        s = s + backtrace + "\n"
    } else {
        s = s + "\n" + backtrace + "\n"
    }
    l.print(logio, s)
}

// 调用回溯字符串
func (l *Logger) backtrace() string {
    backtrace := "Trace:\n"
    for i := 1; i < 10000; i++ {
        if _, cfile, cline, ok := runtime.Caller(i + 3); ok {
            backtrace += strconv.Itoa(i) + ". " + cfile + ":" + strconv.Itoa(cline) + "\n"
        } else {
            break
        }
    }
    return backtrace
}

func (l *Logger) format(s string) string {
    return time.Now().Format("2006-01-02 15:04:05 ") + s
}

func (l *Logger) Print(v ...interface{}) {
    l.stdPrint(fmt.Sprint(v...))
}

func (l *Logger) Printf(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format, v...))
}

func (l *Logger) Println(v ...interface{}) {
    l.stdPrint(fmt.Sprintln(v...))
}

func (l *Logger) Printfln(format string, v ...interface{}) {
    l.stdPrint(fmt.Sprintf(format + "\n", v...))
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
    l.errPrint(fmt.Sprintln(v...))
    os.Exit(1)
}

func (l *Logger) Fatalfln(format string, v ...interface{}) {
    l.errPrint(fmt.Sprintf(format + "\n", v...))
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
    s := fmt.Sprintln(v...)
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Panicfln(format string, v ...interface{}) {
    s := fmt.Sprintf(format + "\n", v...)
    l.errPrint(s)
    panic(s)
}

func (l *Logger) Info(v ...interface{}) {
    l.stdPrint("[INFO] " + fmt.Sprintln(v...))
}

func (l *Logger) Debug(v ...interface{}) {
    if l.GetDebug() {
        l.stdPrint("[DEBU] " + fmt.Sprintln(v...))
    }
}

func (l *Logger) Notice(v ...interface{}) {
    l.errPrint("[NOTI] " + fmt.Sprintln(v...))
}

func (l *Logger) Warning(v ...interface{}) {
    l.errPrint("[WARN] " + fmt.Sprintln(v...))
}

func (l *Logger) Error(v ...interface{}) {
    l.errPrint("[ERRO] " + fmt.Sprintln(v...))
}

func (l *Logger) Critical(v ...interface{}) {
    l.errPrint("[CRIT] " + fmt.Sprintln(v...))
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
    l.stdPrint("[INFO] " + fmt.Sprintf(format, v...) + "\n")
}

func (l *Logger) Debugfln(format string, v ...interface{}) {
    if l.GetDebug() {
        l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...) + "\n")
    }
}

func (l *Logger) Noticefln(format string, v ...interface{}) {
    l.errPrint("[NOTI] " + fmt.Sprintf(format, v...) + "\n")
}

func (l *Logger) Warningfln(format string, v ...interface{}) {
    l.errPrint("[WARN] " + fmt.Sprintf(format, v...) + "\n")
}

func (l *Logger) Errorfln(format string, v ...interface{}) {
    l.errPrint("[ERRO] " + fmt.Sprintf(format, v...) + "\n")
}

func (l *Logger) Criticalfln(format string, v ...interface{}) {
    l.errPrint("[CRIT] " + fmt.Sprintf(format, v...) + "\n")
}