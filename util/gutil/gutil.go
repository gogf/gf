// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides utility functions.
package gutil

import (
<<<<<<< HEAD
=======
	"fmt"
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
	"github.com/gogf/gf/internal/empty"
)

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
	panic(exception)
}

<<<<<<< HEAD
// TryCatch implements try...catch... logistics using internal panic...recover.
func TryCatch(try func(), catch ...func(exception interface{})) {
	defer func() {
		if e := recover(); e != nil && len(catch) > 0 {
			catch[0](e)
=======
// Try implements try... logistics using internal panic...recover.
// It returns error if any exception occurs, or else it returns nil.
func Try(try func()) (err error) {
	defer func() {
		if e := recover(); e != nil {
			err = fmt.Errorf(`%v`, e)
		}
	}()
	try()
	return
}

// TryCatch implements try...catch... logistics using internal panic...recover.
// It automatically calls function <catch> if any exception occurs ans passes the exception as an error.
func TryCatch(try func(), catch ...func(exception error)) {
	defer func() {
		if e := recover(); e != nil && len(catch) > 0 {
			catch[0](fmt.Errorf(`%v`, e))
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
		}
	}()
	try()
}

// IsEmpty checks given <value> empty or not.
// It returns false if <value> is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
	return empty.IsEmpty(value)
}
