// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"fmt"
	"os"
)

// Print prints <v> with newline using fmt.Sprintln.
// The parameter <v> can be multiple variables.
func (l *Logger) Print(v ...interface{}) {
	l.printStd("", v...)
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The parameter <v> can be multiple variables.
func (l *Logger) Printf(format string, v ...interface{}) {
	l.printStd("", l.format(format, v...))
}

// Println is alias of Print.
// See Print.
func (l *Logger) Println(v ...interface{}) {
	l.Print(v...)
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func (l *Logger) Fatal(v ...interface{}) {
	l.printErr(l.getLevelPrefixWithBrackets(LEVEL_FATA), v...)
	os.Exit(1)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.printErr(l.getLevelPrefixWithBrackets(LEVEL_FATA), l.format(format, v...))
	os.Exit(1)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func (l *Logger) Panic(v ...interface{}) {
	l.printErr(l.getLevelPrefixWithBrackets(LEVEL_PANI), v...)
	panic(fmt.Sprint(v...))
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.printErr(l.getLevelPrefixWithBrackets(LEVEL_PANI), l.format(format, v...))
	panic(l.format(format, v...))
}

// Info prints the logging content with [INFO] header and newline.
func (l *Logger) Info(v ...interface{}) {
	if l.checkLevel(LEVEL_INFO) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_INFO), v...)
	}
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func (l *Logger) Infof(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_INFO) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_INFO), l.format(format, v...))
	}
}

// Debug prints the logging content with [DEBU] header and newline.
func (l *Logger) Debug(v ...interface{}) {
	if l.checkLevel(LEVEL_DEBU) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_DEBU), v...)
	}
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func (l *Logger) Debugf(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_DEBU) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_DEBU), l.format(format, v...))
	}
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Notice(v ...interface{}) {
	if l.checkLevel(LEVEL_NOTI) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_NOTI), v...)
	}
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Noticef(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_NOTI) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_NOTI), l.format(format, v...))
	}
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warning(v ...interface{}) {
	if l.checkLevel(LEVEL_WARN) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_WARN), v...)
	}
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Warningf(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_WARN) {
		l.printStd(l.getLevelPrefixWithBrackets(LEVEL_WARN), l.format(format, v...))
	}
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Error(v ...interface{}) {
	if l.checkLevel(LEVEL_ERRO) {
		l.printErr(l.getLevelPrefixWithBrackets(LEVEL_ERRO), v...)
	}
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Errorf(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_ERRO) {
		l.printErr(l.getLevelPrefixWithBrackets(LEVEL_ERRO), l.format(format, v...))
	}
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Critical(v ...interface{}) {
	if l.checkLevel(LEVEL_CRIT) {
		l.printErr(l.getLevelPrefixWithBrackets(LEVEL_CRIT), v...)
	}
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller stack info if stack feature is enabled.
func (l *Logger) Criticalf(format string, v ...interface{}) {
	if l.checkLevel(LEVEL_CRIT) {
		l.printErr(l.getLevelPrefixWithBrackets(LEVEL_CRIT), l.format(format, v...))
	}
}

// checkLevel checks whether the given <level> could be output.
func (l *Logger) checkLevel(level int) bool {
	return l.config.Level&level > 0
}
