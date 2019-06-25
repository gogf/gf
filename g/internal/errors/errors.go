// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package errors provides simple functions to manipulate errors.
//
// This package can be scalable due to https://go.googlesource.com/proposal/+/master/design/go2draft.md.
package errors

import "github.com/gogf/gf/g/util/gconv"

// errorWrapper is a simple wrapper for errors.
type errorWrapper struct {
	s string
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
	return &errorWrapper{
		s: text,
	}
}

// Wrap wraps error with text.
func Wrap(err error, text string) error {
	if err == nil {
		return nil
	}
	return NewText(text + ": " + err.Error())
}

// Error implements interface Error.
func (e *errorWrapper) Error() string {
	return e.s
}
