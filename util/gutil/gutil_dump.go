// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gutil

import (
	"bytes"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"github.com/gogf/gf/util/gconv"
	"os"
	"reflect"
)

<<<<<<< HEAD
=======
// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
// Dump prints variables <i...> to stdout with more manually readable.
func Dump(i ...interface{}) {
	s := Export(i...)
	if s != "" {
		fmt.Println(s)
	}
}

// Export returns variables <i...> as a string with more manually readable.
func Export(i ...interface{}) string {
	buffer := bytes.NewBuffer(nil)
	for _, v := range i {
<<<<<<< HEAD
		if b, ok := v.([]byte); ok {
			buffer.Write(b)
		} else {
			rv := reflect.ValueOf(v)
			kind := rv.Kind()
=======
		switch r := v.(type) {
		case []byte:
			buffer.Write(r)
		case string:
			buffer.WriteString(r)
		default:
			var (
				rv   = reflect.ValueOf(v)
				kind = rv.Kind()
			)
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			case reflect.Slice, reflect.Array:
				v = gconv.Interfaces(v)
<<<<<<< HEAD
			case reflect.Map, reflect.Struct:
				v = gconv.Map(v)
=======
			case reflect.Map:
				v = gconv.Map(v)
			case reflect.Struct:
				if r, ok := v.(apiMapStrAny); ok {
					v = r.MapStrAny()
				} else if r, ok := v.(apiString); ok {
					v = r.String()
				}
>>>>>>> 4ae89dc9f62ced2aaf3c7eeb2eaf438c65c1521c
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
