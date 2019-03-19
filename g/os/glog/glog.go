// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
// @author john, zseeker

// Package glog implements powerful and easy-to-use levelled logging functionality.
// 
// 日志模块, 直接文件/输出操作，没有异步逻辑，没有使用缓存或者通道
package glog

import (
    "github.com/gogf/gf/g/container/gtype"
    "github.com/gogf/gf/g/internal/cmdenv"
    "io"
)

const (
    LEVEL_ALL  = LEVEL_DEBU | LEVEL_INFO | LEVEL_NOTI | LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT
    LEVEL_DEV  = LEVEL_ALL
    LEVEL_PROD = LEVEL_WARN | LEVEL_ERRO | LEVEL_CRIT
    LEVEL_DEBU = 1 << iota
    LEVEL_INFO
    LEVEL_NOTI
    LEVEL_WARN
    LEVEL_ERRO
    LEVEL_CRIT
)

var (
    // default level for log
    defaultLevel = gtype.NewInt(LEVEL_ALL)

    // default logger object, for package method usage
    logger = New()
)

func init() {
    SetDebug(cmdenv.Get("gf.glog.debug", true).Bool())
}

// SetPath sets the directory path for file logging.
//
// 日志日志目录绝对路径.
func SetPath(path string) {
    logger.SetPath(path)
}

// SetFile sets the file name <pattern> for file logging.
// Datetime pattern can be used in <pattern>, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
//
// 日志文件名称.
func SetFile(pattern string) {
    logger.SetFile(pattern)
}

// SetLevel sets the default logging level.
//
// 设置全局的日志记录等级.
func SetLevel(level int) {
    logger.SetLevel(level)
    defaultLevel.Set(level)
}

// SetWriter sets the customed logging <writer> for logging.
// The <writer> object should implements the io.Writer interface.
// Developer can use customed logging <writer> to redirect logging output to another service, 
// eg: kafka, mysql, mongodb, etc.
//
// 可自定义IO接口，IO可以是文件输出、标准输出、网络输出.
func SetWriter(writer io.Writer) {
    logger.SetWriter(writer)
}

// GetWriter returns the customed writer object, which implements the io.Writer interface.
// It returns nil if no customed writer set.
// 
// 返回自定义的IO，默认为nil.
func GetWriter() io.Writer {
    return logger.GetWriter()
}

// GetLevel returns the default logging level value.
// 
// 获取全局的日志记录等级.
func GetLevel() int {
    return defaultLevel.Val()
}

// SetDebug enables/disables the debug level for default logger.
// The debug level is enbaled in default.
// 
// 设置是否允许输出DEBUG信息.
func SetDebug(debug bool) {
    logger.SetDebug(debug)
}

// SetStdPrint sets whether ouptput the logging contents to stdout, which is false indefault.
// 
// 设置写日志的同时开启or关闭控制台打印，默认是关闭的
func SetStdPrint(open bool) {
    logger.SetStdPrint(open)
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
//
// 获取日志目录绝对路径
func GetPath() string {
    return logger.GetPath()
}

// PrintBacktrace prints the caller backtrace, 
// the optional parameter <skip> specify the skipped backtraces offset from the end point.
//
// 打印文件调用回溯信息
func PrintBacktrace(skip...int) {
    logger.PrintBacktrace(skip...)
}

// GetBacktrace returns the caller backtrace content, 
// the optional parameter <skip> specify the skipped backtraces offset from the end point.
//
// 获取文件调用回溯信息.
func GetBacktrace(skip...int) string {
    return logger.GetBacktrace(skip...)
}

// SetBacktrace enables/disables the backtrace feature in failure logging outputs.
//
// 是否关闭全局的backtrace信息
func SetBacktrace(enabled bool) {
    logger.SetBacktrace(enabled)
}

// To is a chaining function, 
// which redirects current logging content output to the sepecified <writer>.
// 
// 链式操作，设置下一次写入日志内容的Writer
func To(writer io.Writer) *Logger {
    return logger.To(writer)
}

// Path is a chaining function,
// which sets the directory path to <path> for current logging content output.
//
// 链式操作，设置下一次输出的日志路径。
func Path(path string) *Logger {
    return logger.Path(path)
}

// Cat is a chaining function, 
// which sets the category to <category> for current logging content output.
// 
// 设置下一次输出的分类，支持多级分类设置.
func Cat(category string) *Logger {
    return logger.Cat(category)
}

// File is a chaining function, 
// which sets file name <pattern> for the current logging content output.
//
// 设置日志输出文件名称格式
func File(pattern string) *Logger {
    return logger.File(pattern)
}

// Level is a chaining function, 
// which sets logging level for the current logging content output.
//
// 设置日志打印等级.
func Level(level int) *Logger {
    return logger.Level(level)
}

// Backtrace is a chaining function, 
// which sets backtrace options for the current logging content output .
//
// 设置文件调用回溯信息.
func Backtrace(enabled bool, skip...int) *Logger {
    return logger.Backtrace(enabled, skip...)
}

// StdPrint is a chaining function, 
// which enables/disables stdout for the current logging content output.
//
// 是否允许在设置输出文件时同时也输出到终端
func StdPrint(enabled bool) *Logger {
    return logger.StdPrint(enabled)
}

// Header is a chaining function, 
// which enables/disables log header for the current logging content output.
//
// 是否打印每行日志头信息(默认开启)
func Header(enabled bool) *Logger {
    return logger.Header(enabled)
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

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(v ...interface{}) {
    logger.Fatal(v ...)
}

// Fatalf prints the logging content with [FATA] header and custom format, then exit the current process.
func Fatalf(format string, v ...interface{}) {
    logger.Fatalf(format, v ...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalfln(format string, v ...interface{}) {
    logger.Fatalfln(format, v ...)
}

func Panic(v ...interface{}) {
    logger.Panic(v ...)
}

func Panicf(format string, v ...interface{}) {
    logger.Panicf(format, v ...)
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
