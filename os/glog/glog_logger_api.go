// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"fmt"
	"os"
)

// Print prints `v` with newline using fmt.Sprintln.
// The parameter `v` can be multiple variables.
func (l *Logger) Print(ctx context.Context, v ...interface{}) {
	l.printStd(ctx, LEVEL_NONE, v...)
}

// Printf prints `v` with format `format` using fmt.Sprintf.
// The parameter `v` can be multiple variables.
func (l *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	l.printStd(ctx, LEVEL_NONE, l.format(format, v...))
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func (l *Logger) Fatal(ctx context.Context, v ...interface{}) {
	l.printErr(ctx, LEVEL_FATA, v...)
	os.Exit(1)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func (l *Logger) Fatalf(ctx context.Context, format string, v ...interface{}) {
	l.printErr(ctx, LEVEL_FATA, l.format(format, v...))
	os.Exit(1)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func (l *Logger) Panic(ctx context.Context, v ...interface{}) {
	l.printErr(ctx, LEVEL_PANI, v...)
	panic(fmt.Sprint(v...))
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func (l *Logger) Panicf(ctx context.Context, format string, v ...interface{}) {
	l.printErr(ctx, LEVEL_PANI, l.format(format, v...))
	panic(l.format(format, v...))
}

// Info prints the logging content with [INFO] header and newline.
func (l *Logger) Info(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_INFO) {
		l.printStd(ctx, LEVEL_INFO, v...)
	}
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_INFO) {
		l.printStd(ctx, LEVEL_INFO, l.format(format, v...))
	}
}

// Debug prints the logging content with [DEBU] header and newline.
func (l *Logger) Debug(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_DEBU) {
		l.printStd(ctx, LEVEL_DEBU, v...)
	}
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func (l *Logger) Debugf(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_DEBU) {
		l.printStd(ctx, LEVEL_DEBU, l.format(format, v...))
	}
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Notice(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_NOTI) {
		l.printStd(ctx, LEVEL_NOTI, v...)
	}
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Noticef(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_NOTI) {
		l.printStd(ctx, LEVEL_NOTI, l.format(format, v...))
	}
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warning(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_WARN) {
		l.printStd(ctx, LEVEL_WARN, v...)
	}
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warningf(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_WARN) {
		l.printStd(ctx, LEVEL_WARN, l.format(format, v...))
	}
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Error(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_ERRO) {
		l.printErr(ctx, LEVEL_ERRO, v...)
	}
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_ERRO) {
		l.printErr(ctx, LEVEL_ERRO, l.format(format, v...))
	}
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Critical(ctx context.Context, v ...interface{}) {
	if l.checkLevel(LEVEL_CRIT) {
		l.printErr(ctx, LEVEL_CRIT, v...)
	}
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Criticalf(ctx context.Context, format string, v ...interface{}) {
	if l.checkLevel(LEVEL_CRIT) {
		l.printErr(ctx, LEVEL_CRIT, l.format(format, v...))
	}
}

// checkLevel checks whether the given `level` could be output.
func (l *Logger) checkLevel(level int) bool {
	// nil logger, print nothing
	if l == nil {
		return false
	}
	return l.config.Level&level > 0
}
