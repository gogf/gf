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

// apiVal is used for type assert api for Val().
type apiVal interface {
	Val() interface{}
}

// apiString is used for type assert api for String().
type apiString interface {
	String() string
}

// apiMapStrAny is the interface support for converting struct parameter to map.
type apiMapStrAny interface {
	MapStrAny() map[string]interface{}
}

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
	for _, value := range i {
		switch r := value.(type) {
		case []byte:
			buffer.Write(r)
		case string:
			buffer.WriteString(r)
		default:
			var (
				reflectValue = reflect.ValueOf(value)
				reflectKind  = reflectValue.Kind()
			)
			for reflectKind == reflect.Ptr {
				reflectValue = reflectValue.Elem()
				reflectKind = reflectValue.Kind()
			}
			switch reflectKind {
			case reflect.Slice, reflect.Array:
				value = gconv.Interfaces(value)
			case reflect.Map:
				value = gconv.Map(value)
			case reflect.Struct:
				converted := false
				if r, ok := value.(apiVal); ok {
					if result := r.Val(); result != nil {
						value = result
						converted = true
					}
				}
				if !converted {
					if r, ok := value.(apiMapStrAny); ok {
						if result := r.MapStrAny(); result != nil {
							value = result
							converted = true
						}
					}
				}
				if !converted {
					if r, ok := value.(apiString); ok {
						value = r.String()
					}
				}
			}
			encoder := json.NewEncoder(buffer)
			encoder.SetEscapeHTML(false)
			encoder.SetIndent("", "\t")
			if err := encoder.Encode(value); err != nil {
				fmt.Fprintln(os.Stderr, err.Error())
			}
		}
	}
	return buffer.String()
}
