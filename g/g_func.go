// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package g

import (
    "github.com/gogf/gf/g/container/gvar"
    "github.com/gogf/gf/g/internal/empty"
    "github.com/gogf/gf/g/net/ghttp"
    "github.com/gogf/gf/g/util/gutil"
)

// NewVar returns a *gvar.Var.
func NewVar(i interface{}, unsafe...bool) *Var {
    return gvar.New(i, unsafe...)
}

// Wait blocks until all the web servers shutdown.
func Wait() {
    ghttp.Wait()
}

// Dump dumps a variable to stdout with more manually readable.
func Dump(i...interface{}) {
    gutil.Dump(i...)
}

// Export exports a variable to string with more manually readable.
func Export(i...interface{}) string {
    return gutil.Export(i...)
}

// Throw throws a exception, which can be caught by TryCatch function.
// It always be used in TryCatch function.
func Throw(exception interface{}) {
    gutil.Throw(exception)
}

// TryCatch does the try...catch... logic.
func TryCatch(try func(), catch ... func(exception interface{})) {
    gutil.TryCatch(try, catch...)
}

// IsEmpty checks given value empty or not.
// It returns false if value is: integer(0), bool(false), slice/map(len=0), nil;
// or else true.
func IsEmpty(value interface{}) bool {
    return empty.IsEmpty(value)
}