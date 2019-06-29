// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package errors provides simple functions to manipulate errors.
package gerror

import (
	"bytes"
	"fmt"
	"runtime"
	"strings"

	"github.com/gogf/gf/g/util/gconv"
	"github.com/pkg/errors"
)

type stacker interface {
	StackTrace() errors.StackTrace
}

type causer interface {
	Cause() error
}

const (
	gFILTER_KEY = "/g/errors/gerror/gerror.go"
)

var (
	// goRootForFilter is used for stack filtering purpose.
	goRootForFilter = runtime.GOROOT()
)

func init() {
	if goRootForFilter != "" {
		goRootForFilter = strings.Replace(goRootForFilter, "\\", "/", -1)
	}
}

// New returns an error that formats as the given value.
func New(value interface{}) error {
	if value == nil {
		return nil
	}
	return NewText(gconv.String(value))
}

// NewText returns an error that formats as the given text.
func NewText(text string) error {
	if text == "" {
		return nil
	}
	return errors.New(text)
}

// Wrap wraps error with text.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return errors.Wrap(err, text)
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// If err is nil, Wrapf returns nil.
func Wrapf(err error, format string, args ...interface{}) error {
	return errors.Wrapf(err, format, args...)
}

// Cause returns the underlying cause of the error, if possible.
// An error value has a cause if it implements the following
// interface:
//
//     type causer interface {
//            Cause() error
//     }
//
// If the error does not implement Cause, the original error will
// be returned. If the error is nil, nil will be returned without further
// investigation.
func Cause(err error) error {
	return errors.Cause(err)
}

// Stack returns the stack callers as string.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if _, ok := err.(causer); !ok {
		return ""
	}
	index := 1
	buffer := bytes.NewBuffer(nil)
	for err != nil {
		cause, ok := err.(causer)
		if !ok {
			if err, ok := err.(stacker); ok {
				buffer.WriteString(fmt.Sprintf("%d.\t%s\n", index, err))
				index++
				formatSubStack(err, buffer)
			}
			break
		}
		if err, ok := err.(stacker); ok {
			buffer.WriteString(fmt.Sprintf("%d.\t%s\n", index, err))
			index++
			formatSubStack(err, buffer)
		}
		err = cause.Cause()
	}
	return buffer.String()
}

// formatSubStack formats the stack for error.
func formatSubStack(err stacker, buffer *bytes.Buffer) {
	index := 1
	for _, f := range err.StackTrace() {
		if fn := runtime.FuncForPC(uintptr(f) - 1); fn != nil {
			file, line := fn.FileLine(uintptr(f) - 1)
			if strings.Contains(file, gFILTER_KEY) {
				continue
			}
			if goRootForFilter != "" && len(file) >= len(goRootForFilter) && file[0:len(goRootForFilter)] == goRootForFilter {
				continue
			}
			buffer.WriteString(fmt.Sprintf("\t%d).\t%s\n\t\t%s:%d\n", index, fn.Name(), file, line))
			index++
		}
	}
}
