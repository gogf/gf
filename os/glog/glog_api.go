// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import "context"

// Print prints `v` with newline using fmt.Sprintln.
// The parameter `v` can be multiple variables.
func Print(ctx context.Context, v ...interface{}) {
	logger.Print(ctx, v...)
}

// Printf prints `v` with format `format` using fmt.Sprintf.
// The parameter `v` can be multiple variables.
func Printf(ctx context.Context, format string, v ...interface{}) {
	logger.Printf(ctx, format, v...)
}

// Println is alias of Print.
// See Print.
func Println(ctx context.Context, v ...interface{}) {
	logger.Println(ctx, v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func Fatal(ctx context.Context, v ...interface{}) {
	logger.Fatal(ctx, v...)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func Fatalf(ctx context.Context, format string, v ...interface{}) {
	logger.Fatalf(ctx, format, v...)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func Panic(ctx context.Context, v ...interface{}) {
	logger.Panic(ctx, v...)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func Panicf(ctx context.Context, format string, v ...interface{}) {
	logger.Panicf(ctx, format, v...)
}

// Info prints the logging content with [INFO] header and newline.
func Info(ctx context.Context, v ...interface{}) {
	logger.Info(ctx, v...)
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func Infof(ctx context.Context, format string, v ...interface{}) {
	logger.Infof(ctx, format, v...)
}

// Debug prints the logging content with [DEBU] header and newline.
func Debug(ctx context.Context, v ...interface{}) {
	logger.Debug(ctx, v...)
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func Debugf(ctx context.Context, format string, v ...interface{}) {
	logger.Debugf(ctx, format, v...)
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Notice(ctx context.Context, v ...interface{}) {
	logger.Notice(ctx, v...)
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Noticef(ctx context.Context, format string, v ...interface{}) {
	logger.Noticef(ctx, format, v...)
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Warning(ctx context.Context, v ...interface{}) {
	logger.Warning(ctx, v...)
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Warningf(ctx context.Context, format string, v ...interface{}) {
	logger.Warningf(ctx, format, v...)
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Error(ctx context.Context, v ...interface{}) {
	logger.Error(ctx, v...)
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Errorf(ctx context.Context, format string, v ...interface{}) {
	logger.Errorf(ctx, format, v...)
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func Critical(ctx context.Context, v ...interface{}) {
	logger.Critical(ctx, v...)
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func Criticalf(ctx context.Context, format string, v ...interface{}) {
	logger.Criticalf(ctx, format, v...)
}
