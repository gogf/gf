// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

// Option is option for creating error.
type Option struct {
	Error error  // Wrapped error.
	Stack bool   // Record stack information into error.
	Text  string // Error text, which is created by New* functions.
	Code  int    // Error code if necessary.
}

// NewOption creates and returns an error with Option.
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
