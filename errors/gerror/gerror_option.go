// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

// Option is option for creating error.
type Option struct {
	Error error  // Wrapped error if any.
	Stack bool   // Whether recording stack information into error.
	Text  string // Error text, which is created by New* functions.
	Code  int    // Error code if necessary.
}

// NewOption creates and returns an error with Option.
// It is the senior usage for creating error, which is often used internally in framework.
func NewOption(option Option) error {
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
