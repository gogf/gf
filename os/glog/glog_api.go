// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

// Print prints <v> with newline using fmt.Sprintln.
// The parameter <v> can be multiple variables.
func Print(v ...interface{}) {
	logger.Print(v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func Printf(format string, v ...interface{}) {
	logger.Printf(format, v...)
}

// See Print.
func Println(v ...interface{}) {
	logger.Println(v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(v ...interface{}) {
	logger.Fatal(v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(format string, v ...interface{}) {
	logger.Fatalf(format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(v ...interface{}) {
	logger.Panic(v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func Panicf(format string, v ...interface{}) {
	logger.Panicf(format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(v ...interface{}) {
	logger.Info(v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(format string, v ...interface{}) {
	logger.Infof(format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(v ...interface{}) {
	logger.Debug(v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(format string, v ...interface{}) {
	logger.Debugf(format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(v ...interface{}) {
	logger.Notice(v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(format string, v ...interface{}) {
	logger.Noticef(format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(v ...interface{}) {
	logger.Warning(v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(format string, v ...interface{}) {
	logger.Warningf(format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(v ...interface{}) {
	logger.Error(v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(format string, v ...interface{}) {
	logger.Errorf(format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(v ...interface{}) {
	logger.Critical(v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(format string, v ...interface{}) {
	logger.Criticalf(format, v...)
}
