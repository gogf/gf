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
// The param <v> can be multiple variables.
func (l *Logger) Print(v ...interface{}) {
    l.printStd(fmt.Sprintln(v...))
}

// Printf prints <v> with format <format> using fmt.Sprintf.
// The param <v> can be multiple variables.
func (l *Logger) Printf(format string, v ...interface{}) {
    l.printStd(fmt.Sprintf(format + ln, v...))
}

// See Print.
func (l *Logger) Println(v ...interface{}) {
    l.Print(v...)
}

// Deprecated.
// Use Printf instead.
func (l *Logger) Printfln(format string, v ...interface{}) {
    l.printStd(fmt.Sprintf(format + ln, v...))
}

// Fatal prints the logging content with [FATA] header and newline, then exit the current process.
func (l *Logger) Fatal(v ...interface{}) {
    l.printErr("[FATA] " + fmt.Sprintln(v...))
    os.Exit(1)
}

// Fatalf prints the logging content with [FATA] header, custom format and newline, then exit the current process.
func (l *Logger) Fatalf(format string, v ...interface{}) {
    l.printErr("[FATA] " + fmt.Sprintf(format + ln, v...))
    os.Exit(1)
}

// Deprecated.
// Use Fatalf instead.
func (l *Logger) Fatalfln(format string, v ...interface{}) {
    l.printErr("[FATA] " + fmt.Sprintf(format + ln, v...))
    os.Exit(1)
}

// Panic prints the logging content with [PANI] header and newline, then panics.
func (l *Logger) Panic(v ...interface{}) {
    s := fmt.Sprintln(v...)
    l.printErr("[PANI] " + s)
    panic(s)
}

// Panicf prints the logging content with [PANI] header, custom format and newline, then panics.
func (l *Logger) Panicf(format string, v ...interface{}) {
    s := fmt.Sprintf(format + ln, v...)
    l.printErr("[PANI] " + s)
    panic(s)
}

// Deprecated.
// Use Panicf instead.
func (l *Logger) Panicfln(format string, v ...interface{}) {
    s := fmt.Sprintf(format + ln, v...)
    l.printErr("[PANI] " + s)
    panic(s)
}

// Info prints the logging content with [INFO] header and newline.
func (l *Logger) Info(v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.printStd("[INFO] " + fmt.Sprintln(v...))
    }
}

// Infof prints the logging content with [INFO] header, custom format and newline.
func (l *Logger) Infof(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.printStd("[INFO] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Infof instead.
func (l *Logger) Infofln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_INFO) {
        l.printStd("[INFO] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// Debug prints the logging content with [DEBU] header and newline.
func (l *Logger) Debug(v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.printStd("[DEBU] " + fmt.Sprintln(v...))
    }
}

// Debugf prints the logging content with [DEBU] header, custom format and newline.
func (l *Logger) Debugf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.printStd("[DEBU] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Debugf instead.
func (l *Logger) Debugfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_DEBU) {
        l.printStd("[DEBU] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// Notice prints the logging content with [NOTI] header and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Notice(v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.printErr("[NOTI] " + fmt.Sprintln(v...))
    }
}

// Noticef prints the logging content with [NOTI] header, custom format and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Noticef(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.printErr("[NOTI] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Noticef instead.
func (l *Logger) Noticefln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_NOTI) {
        l.printErr("[NOTI] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// Warning prints the logging content with [WARN] header and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Warning(v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.printErr("[WARN] " + fmt.Sprintln(v...))
    }
}

// Warningf prints the logging content with [WARN] header, custom format and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Warningf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.printErr("[WARN] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Warningf instead.
func (l *Logger) Warningfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_WARN) {
        l.printErr("[WARN] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// Error prints the logging content with [ERRO] header and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Error(v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.printErr("[ERRO] " + fmt.Sprintln(v...))
    }
}

// Errorf prints the logging content with [ERRO] header, custom format and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Errorf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.printErr("[ERRO] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Errorf instead.
func (l *Logger) Errorfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_ERRO) {
        l.printErr("[ERRO] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// Critical prints the logging content with [CRIT] header and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Critical(v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.printErr("[CRIT] " + fmt.Sprintln(v...))
    }
}

// Criticalf prints the logging content with [CRIT] header, custom format and newline.
// It also prints caller backtrace info if backtrace feature is enabled.
func (l *Logger) Criticalf(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.printErr("[CRIT] " + fmt.Sprintf(format + ln, v...))
    }
}

// Deprecated.
// Use Criticalf instead.
func (l *Logger) Criticalfln(format string, v ...interface{}) {
    if l.checkLevel(LEVEL_CRIT) {
        l.printErr("[CRIT] " + fmt.Sprintf(format + ln, v...) + ln)
    }
}

// checkLevel checks whether the given <level> could be output.
func (l *Logger) checkLevel(level int) bool {
    return l.level & level > 0
}