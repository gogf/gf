// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package glog implements powerful and easy-to-use levelled logging functionality.
package glog

import (
	"io"

	"github.com/gogf/gf/internal/cmdenv"
	"github.com/gogf/gf/os/grpool"
)

var (
	// Default logger object, for package method usage
	logger = New()
	// Goroutine pool for async logging output.
	// It uses only one asynchronize worker to ensure log sequence.
	asyncPool = grpool.New(1)
	// defaultDebug enables debug level or not in default,
	// which can be configured using command option or system environment.
	defaultDebug = true
)

func init() {
	defaultDebug = cmdenv.Get("gf.glog.debug", true).Bool()
	SetDebug(defaultDebug)
}

// Default returns the default logger.
func DefaultLogger() *Logger {
	return logger
}

// SetDefaultLogger sets the default logger for package glog.
// Note that there might be concurrent safety issue if calls this function
// in different goroutines.
func SetDefaultLogger(l *Logger) {
	logger = l
}

// SetPath sets the directory path for file logging.
func SetPath(path string) error {
	return logger.SetPath(path)
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

// SetAsync enables/disables async logging output feature for default logger.
func SetAsync(enabled bool) {
	logger.SetAsync(enabled)
}

// SetStdoutPrint sets whether ouptput the logging contents to stdout, which is true in default.
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

// PrintStack prints the caller stack,
// the optional parameter <skip> specify the skipped stack offset from the end point.
func PrintStack(skip ...int) {
	logger.PrintStack(skip...)
}

// GetStack returns the caller stack content,
// the optional parameter <skip> specify the skipped stack offset from the end point.
func GetStack(skip ...int) string {
	return logger.GetStack(skip...)
}

// SetStack enables/disables the stack feature in failure logging outputs.
func SetStack(enabled bool) {
	logger.SetStack(enabled)
}
