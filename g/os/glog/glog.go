// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// @author john, zseeker

// Package glog implements powerful and easy-to-use levelled logging functionality.
package glog

import (
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
    // Default logger object, for package method usage
    logger = New()
)

func init() {
    SetDebug(cmdenv.Get("gf.glog.debug", true).Bool())
}

// SetPath sets the directory path for file logging.
func SetPath(path string) {
    logger.SetPath(path)
}

// GetPath returns the logging directory path for file logging.
// It returns empty string if no directory path set.
func GetPath() string {
	return logger.GetPath()
}

// SetFile sets the file name <pattern> for file logging.
// Datetime pattern can be used in <pattern>, eg: access-{Ymd}.log.
// The default file name pattern is: Y-m-d.log, eg: 2018-01-01.log
func SetFile(pattern string) {
    logger.SetFile(pattern)
}

// SetLevel sets the default logging level.
func SetLevel(level int) {
    logger.SetLevel(level)
}

// GetLevel returns the default logging level value.
func GetLevel() int {
	return logger.GetLevel()
}

// SetWriter sets the customized logging <writer> for logging.
// The <writer> object should implements the io.Writer interface.
// Developer can use customized logging <writer> to redirect logging output to another service, 
// eg: kafka, mysql, mongodb, etc.
func SetWriter(writer io.Writer) {
    logger.SetWriter(writer)
}

// GetWriter returns the customized writer object, which implements the io.Writer interface.
// It returns nil if no customized writer set.
func GetWriter() io.Writer {
    return logger.GetWriter()
}

// SetDebug enables/disables the debug level for default logger.
// The debug level is enbaled in default.
func SetDebug(debug bool) {
    logger.SetDebug(debug)
}

// SetStdoutPrint sets whether ouptput the logging contents to stdout, which is false in default.
func SetStdoutPrint(enabled bool) {
    logger.SetStdoutPrint(enabled)
}

// SetHeaderPrint sets whether output header of the logging contents, which is true in default.
func SetHeaderPrint(enabled bool) {
	logger.SetHeaderPrint(enabled)
}

// SetPrefix sets prefix string for every logging content.
// Prefix is part of header, which means if header output is shut, no prefix will be output.
func SetPrefix(prefix string) {
	logger.SetPrefix(prefix)
}

// SetFlags sets extra flags for logging output features.
func SetFlags(flags int) {
	logger.SetFlags(flags)
}

// GetFlags returns the flags of logger.
func GetFlags() int {
	return logger.GetFlags()
}

// PrintBacktrace prints the caller backtrace, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func PrintBacktrace(skip...int) {
    logger.PrintBacktrace(skip...)
}

// GetBacktrace returns the caller backtrace content, 
// the optional parameter <skip> specify the skipped backtrace offset from the end point.
func GetBacktrace(skip...int) string {
    return logger.GetBacktrace(skip...)
}

// SetBacktrace enables/disables the backtrace feature in failure logging outputs.
func SetBacktrace(enabled bool) {
    logger.SetBacktrace(enabled)
}
