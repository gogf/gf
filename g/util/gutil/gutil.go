// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gutil provides utility functions.
package gutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogf/gf/g/internal/empty"
	"github.com/gogf/gf/g/util/gconv"
	"os"
	"runtime"
)

// Dump prints variables <i...> to stdout with more manually readable.
func Dump(i...interface{}) {
    s := Export(i...)
    if s != "" {
        fmt.Println(s)
    }
}

// Export returns variables <i...> as a string with more manually readable.
func Export(i...interface{}) string {
    buffer := bytes.NewBuffer(nil)
    for _, v := range i {
        if b, ok := v.([]byte); ok {
            buffer.Write(b)
        } else {
            if m := gconv.Map(v); m != nil {
            	v = m
            }
            encoder := json.NewEncoder(buffer)
            encoder.SetEscapeHTML(false)
            encoder.SetIndent("", "\t")
            if err := encoder.Encode(v); err != nil {
                fmt.Fprintln(os.Stderr, err.Error())
            }
        }
    }
    return buffer.String()
}

// PrintBacktrace prints the caller backtrace to stdout.
func PrintBacktrace() {
    index  := 1
    buffer := bytes.NewBuffer(nil)
    for i := 1; i < 10000; i++ {
        if _, path, line, ok := runtime.Caller(i); ok {
            buffer.WriteString(fmt.Sprintf(`%d. %s:%d%s`, index, path, line, "\n"))
            index++
        } else {
            break
        }
    }
    fmt.Print(buffer.String())
}

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
    panic(exception)
}

// TryCatch implements try...catch... logistics.
func TryCatch(try func(), catch ... func(exception interface{})) {
    if len(catch) > 0 {
        defer func() {
            if e := recover(); e != nil {
                catch[0](e)
            }
        }()
    }
    try()
}

// IsEmpty checks given <value> empty or not.
// It returns false if <value> is: integer(0), bool(false), slice/map(len=0), nil;
// or else returns true.
func IsEmpty(value interface{}) bool {
    return empty.IsEmpty(value)
}


