// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"io"
)

// To is a chaining function, 
// which redirects current logging content output to the sepecified <writer>.
func To(writer io.Writer) *Logger {
    return logger.To(writer)
}

// Path is a chaining function,
// which sets the directory path to <path> for current logging content output.
func Path(path string) *Logger {
    return logger.Path(path)
}

// Cat is a chaining function, 
// which sets the category to <category> for current logging content output.
func Cat(category string) *Logger {
    return logger.Cat(category)
}

// File is a chaining function, 
// which sets file name <pattern> for the current logging content output.
func File(pattern string) *Logger {
    return logger.File(pattern)
}

// Level is a chaining function, 
// which sets logging level for the current logging content output.
func Level(level int) *Logger {
    return logger.Level(level)
}

// Skip is a chaining function,
// which sets backtrace skip for the current logging content output.
// It also affects the caller file path checks when line number printing enabled.
func Skip(skip int) *Logger {
	return logger.Skip(skip)
}

// Backtrace is a chaining function, 
// which sets backtrace options for the current logging content output .
func Backtrace(enabled bool, skip...int) *Logger {
    return logger.Backtrace(enabled, skip...)
}

// StdPrint is a chaining function, 
// which enables/disables stdout for the current logging content output.
// It's enabled in default.
func Stdout(enabled...bool) *Logger {
    return logger.Stdout(enabled...)
}

// Header is a chaining function, 
// which enables/disables log header for the current logging content output.
// It's enabled in default.
func Header(enabled...bool) *Logger {
    return logger.Header(enabled...)
}

// Line is a chaining function,
// which enables/disables printing its caller file along with its line number.
// The param <long> specified whether print the long absolute file path, eg: /a/b/c/d.go:23.
func Line(long...bool) *Logger {
	return logger.Line(long...)
}
