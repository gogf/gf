package glog

import (
    "sync"
    "os"
    "io"
    "strings"
    "reflect"
    "path/filepath"
    "time"
    "fmt"
)

type Logger struct {
    mutex        sync.RWMutex
    logio        io.Writer
    logpath      string       // 日志写入的目录路径
    lastlogdate  string       // 上一次写入日志的日期，例如: 2006-01-02
}

// 默认的日志对象
var logger = New()

// 新建自定义的日志操作对象
func New() *Logger {
    return &Logger{ }
}

func SetLogPath(path string) {
    logger.SetLogPath(path)
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

func Fatal(v ...interface{}) {
    logger.Fatal(v ...)
}

func Fatalf(format string, v ...interface{}) {
    logger.Fatalf(format, v ...)
}

func Fatalln(v ...interface{}) {
    logger.Fatalln(v ...)
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

func (l *Logger) GetLogIO() io.Writer {
    l.mutex.RLock()
    r := l.logio
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

// 设置日志文件的存储目录路径
func (l *Logger) SetLogPath(path string) {
    l.mutex.Lock()
    l.logpath  = strings.TrimRight(path, string(filepath.Separator))
    l.mutex.Unlock()
    // 重新检查日志io对象
    l.checkLogIO()
}

// 检查文件名称是否已经过期
func (l *Logger) checkLogIO() {
    date := time.Now().Format("2006-01-02")
    if date != l.GetLastLogDate() {
        path := l.GetLogPath()
        if path != "" {
            if !exists(path) {
                mkdir(path)
            }

            l.mutex.Lock()
            fname     := date + ".log"
            fpath     := l.logpath + string(filepath.Separator) + fname
            fio, err  := os.OpenFile(fpath, os.O_WRONLY|os.O_APPEND, 0755)
            if err == nil && fio != nil {
                if l.logio != nil && reflect.TypeOf(l.logio).String() == "*os.File" {
                    l.logio.(*os.File).Close()
                }
                l.logio = fio
            } else {
                fmt.Fprintln(os.Stderr, err)
            }
            l.mutex.Unlock()
        }
    }
}

// 核心打印数据方法(标准输出)
func (l *Logger) stdPrint(s string) {
    l.checkLogIO()
    l.mutex.Lock()
    if l.logio == nil {
        fmt.Fprint(os.Stdout, l.format(s))
    } else {
        fmt.Fprint(l.logio, l.format(s))
    }
    l.mutex.Unlock()
}

// 核心打印数据方法(标准错误)
func (l *Logger) errPrint(s string) {
    l.checkLogIO()
    l.mutex.Lock()
    if l.logio == nil {
        fmt.Fprint(os.Stderr, l.format(s))
    } else {
        fmt.Fprint(l.logio, l.format(s))
    }
    l.mutex.Unlock()
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

func (l *Logger) Info(v ...interface{}) {
    l.stdPrint("[INFO] " + fmt.Sprintln(v...))
}

func (l *Logger) Debug(v ...interface{}) {
    l.stdPrint("[DEBU] " + fmt.Sprintln(v...))
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
    l.stdPrint("[DEBU] " + fmt.Sprintf(format, v...))
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

// 给定文件的绝对路径创建文件
func mkdir(path string) error {
    err  := os.MkdirAll(path, os.ModePerm)
    if err != nil {
        return err
    }
    return nil
}

// 判断所给路径文件/文件夹是否存在
func exists(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}