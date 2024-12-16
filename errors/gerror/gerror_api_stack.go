// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gerror

import (
	"errors"
	"runtime"
)

// stack represents a stack of program counters.
type stack []uintptr

const (
	// maxStackDepth marks the max stack depth for error back traces.
	maxStackDepth = 64
)

// Cause returns the root cause error of `err`.
func Cause(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(ICause); ok {
		return e.Cause()
	}
	if e, ok := err.(IUnwrap); ok {
		return Cause(e.Unwrap())
	}
	return err
}

// Stack returns the stack callers as string.
// It returns the error string directly if the `err` does not support stacks.
func Stack(err error) string {
	if err == nil {
		return ""
	}
	if e, ok := err.(IStack); ok {
		return e.Stack()
	}
	return err.Error()
}

// Current creates and returns the current level error.
// It returns nil if current level error is nil.
func Current(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(ICurrent); ok {
		return e.Current()
	}
	return err
}

// Unwrap returns the next level error.
// It returns nil if current level error or the next level error is nil.
func Unwrap(err error) error {
	if err == nil {
		return nil
	}
	if e, ok := err.(IUnwrap); ok {
		return e.Unwrap()
	}
	return nil
}

// HasStack checks and reports whether `err` implemented interface `gerror.IStack`.
func HasStack(err error) bool {
	_, ok := err.(IStack)
	return ok
}

// Equal reports whether current error `err` equals to error `target`.
// Please note that, in default comparison logic for `Error`,
// the errors are considered the same if both the `code` and `text` of them are the same.
func Equal(err, target error) bool {
	if err == target {
		return true
	}
	if e, ok := err.(IEqual); ok {
		return e.Equal(target)
	}
	if e, ok := target.(IEqual); ok {
		return e.Equal(err)
	}
	return false
}

// Is reports whether current error `err` has error `target` in its chaining errors.
// There's similar function HasError which is designed and implemented early before errors.Is of go stdlib.
// It is now alias of errors.Is of go stdlib, to guarantee the same performance as go stdlib.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func As(err error, target any) bool {
	return errors.As(err, target)
}

// HasError performs as Is.
// This function is designed and implemented early before errors.Is of go stdlib.
// Deprecated: use Is instead.
func HasError(err, target error) bool {
	return errors.Is(err, target)
}

// callers returns the stack callers.
// Note that it here just retrieves the caller memory address array not the caller information.
func callers(skip ...int) stack {
	var (
		pcs [maxStackDepth]uintptr
		n   = 3
	)
	if len(skip) > 0 {
		n += skip[0]
	}
	return pcs[:runtime.Callers(n, pcs[:])]
}
