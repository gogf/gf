// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import "github.com/gogf/gf/v2/errors/gcode"

// Option is option for creating error.
type Option struct {
	Error error      // Wrapped error if any.
	Stack bool       // Whether recording stack information into error.
	Text  string     // Error text, which is created by New* functions.
	Code  gcode.Code // Error code if necessary.
}

// NewWithOption creates and returns a custom error with Option.
// It is the senior usage for creating error, which is often used internally in framework.
func NewWithOption(option Option) error {
	err := &Error{
		error: option.Error,
		text:  option.Text,
		code:  option.Code,
	}
	if option.Stack {
		err.stack = callers()
	}
	return err
}

// NewOption creates and returns a custom error with Option.
// Deprecated: use NewWithOption instead.
func NewOption(option Option) error {
	return NewWithOption(option)
}
