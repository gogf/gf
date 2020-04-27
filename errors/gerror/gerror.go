// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package errors provides simple functions to manipulate errors.
package gerror

import (
	"fmt"
)

// ApiStack is the interface for Stack feature.
type ApiStack interface {
	Stack() string
}

// ApiCause is the interface for Cause feature.
type ApiCause interface {
	Cause() error
}

// New creates and returns an error which is formatted from given text.
func New(text string) error {
	if text == "" {
		return nil
	}
	return &Error{
		stack: callers(),
		text:  text,
	}
}

// Newf returns an error that formats as the given format and args.
func Newf(format string, args ...interface{}) error {
	if format == "" {
		return nil
	}
	return &Error{
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
	}
}

// Wrap wraps error with text.
// It returns nil if given err is nil.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  text,
	}
}

// Wrapf returns an error annotating err with a stack trace
// at the point Wrapf is called, and the format specifier.
// It returns nil if given <err> is nil.
func Wrapf(err error, format string, args ...interface{}) error {
	if err == nil {
		return nil
	}
	return &Error{
		error: err,
		stack: callers(),
		text:  fmt.Sprintf(format, args...),
	}
}

// Cause returns the root cause error of <err>.
func Cause(err error) error {
	if err != nil {
		if e, ok := err.(ApiCause); ok {
			return e.Cause()
		}
	}
	return err
}

// Stack returns the stack callers as string.
// It returns an empty string if the <err> does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(ApiStack); ok {
		return e.Stack()
	}
	return ""
}
