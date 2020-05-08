// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package glog

import (
	"context"
	"io"
)

// Expose returns the default logger of glog.
func Expose() *Logger {
	return logger
}

// Ctx is a chaining function,
// which sets the context for current logging.
// The parameter <keys> specifies the context keys for retrieving values.
func Ctx(ctx context.Context, keys ...interface{}) *Logger {
	return logger.Ctx(ctx, keys...)
}

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

// LevelStr is a chaining function,
// which sets logging level for the current logging content output using level string.
func LevelStr(levelStr string) *Logger {
	return logger.LevelStr(levelStr)
}

// Skip is a chaining function,
// which sets stack skip for the current logging content output.
// It also affects the caller file path checks when line number printing enabled.
func Skip(skip int) *Logger {
	return logger.Skip(skip)
}

// Stack is a chaining function,
// which sets stack options for the current logging content output .
func Stack(enabled bool, skip ...int) *Logger {
	return logger.Stack(enabled, skip...)
}

// StackWithFilter is a chaining function,
// which sets stack filter for the current logging content output .
func StackWithFilter(filter string) *Logger {
	return logger.StackWithFilter(filter)
}

// StdPrint is a chaining function,
// which enables/disables stdout for the current logging content output.
// It's enabled in default.
func Stdout(enabled ...bool) *Logger {
	return logger.Stdout(enabled...)
}

// Header is a chaining function,
// which enables/disables log header for the current logging content output.
// It's enabled in default.
func Header(enabled ...bool) *Logger {
	return logger.Header(enabled...)
}

// Line is a chaining function,
// which enables/disables printing its caller file along with its line number.
// The parameter <long> specified whether print the long absolute file path, eg: /a/b/c/d.go:23.
func Line(long ...bool) *Logger {
	return logger.Line(long...)
}

// Async is a chaining function,
// which enables/disables async logging output feature.
func Async(enabled ...bool) *Logger {
	return logger.Async(enabled...)
}
