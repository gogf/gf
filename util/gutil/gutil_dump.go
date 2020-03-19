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
	"os"
	"reflect"
	"strings"

	"github.com/gogf/gf/internal/empty"
	"github.com/gogf/gf/util/gconv"
)

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
		if b, ok := v.([]byte); ok {
			buffer.Write(b)
		} else {
			rv := reflect.ValueOf(value)
			kind := rv.Kind()
			// If it is a pointer, we should find its real data type.
			if kind == reflect.Ptr {
				rv = rv.Elem()
				kind = rv.Kind()
			}
			switch kind {
			// If <value> is type of array, it converts the value of even number index as its key and
			// the value of odd number index as its corresponding value.
			// Eg:
			// []string{"k1","v1","k2","v2"} => map[string]interface{}{"k1":"v1", "k2":"v2"}
			// []string{"k1","v1","k2"} => map[string]interface{}{"k1":"v1", "k2":nil}
			case reflect.Slice, reflect.Array:
				length := rv.Len()
				for i := 0; i < length; i += 2 {
					if i+1 < length {
						m[String(rv.Index(i).Interface())] = rv.Index(i + 1).Interface()
					} else {
						m[String(rv.Index(i).Interface())] = nil
					}
				}
			case reflect.Map:
				ks := rv.MapKeys()
				for _, k := range ks {
					m[String(k.Interface())] = rv.MapIndex(k).Interface()
				}
			case reflect.Struct:
				// Map converting interface check.
				if v, ok := value.(apiMapStrAny); ok {
					return v.MapStrAny()
				}
				rt := rv.Type()
				name := ""
				tagArray := structTagPriority
				switch len(tags) {
				case 0:
					// No need handle.
				case 1:
					tagArray = append(strings.Split(tags[0], ","), structTagPriority...)
				default:
					tagArray = append(tags, structTagPriority...)
				}
				var rtField reflect.StructField
				var rvField reflect.Value
				var rvKind reflect.Kind
				for i := 0; i < rv.NumField(); i++ {
					rtField = rt.Field(i)
					rvField = rv.Field(i)
					// Only convert the public attributes.
					fieldName := rtField.Name
					if !utils.IsLetterUpper(fieldName[0]) {
						continue
					}
					name = ""
					fieldTag := rtField.Tag
					for _, tag := range tagArray {
						if name = fieldTag.Get(tag); name != "" {
							break
						}
					}
					if name == "" {
						name = strings.TrimSpace(fieldName)
					} else {
						// Support json tag feature: -, omitempty
						name = strings.TrimSpace(name)
						if name == "-" {
							continue
						}
						array := strings.Split(name, ",")
						if len(array) > 1 {
							switch strings.TrimSpace(array[1]) {
							case "omitempty":
								if empty.IsEmpty(rvField.Interface()) {
									continue
								} else {
									name = strings.TrimSpace(array[0])
								}
							default:
								name = strings.TrimSpace(array[0])
							}
						}
					}
					if recursive {
						rvKind = rvField.Kind()
						if rvKind == reflect.Ptr {
							rvField = rvField.Elem()
							rvKind = rvField.Kind()
						}
						if rvKind == reflect.Struct {
							for k, v := range doMapConvert(rvField.Interface(), recursive, tags...) {
								m[k] = v
							}
						} else {
							m[name] = rvField.Interface()
						}
					} else {
						m[name] = rvField.Interface()
					}
				}
			default:
				return nil
			}
			if m := gconv.Map(v); m != nil && len(m) > 0 {
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

// Throw throws out an exception, which can be caught be TryCatch or recover.
func Throw(exception interface{}) {
	panic(exception)
}

// TryCatch implements try...catch... logistics.
func TryCatch(try func(), catch ...func(exception interface{})) {
	if len(catch) > 0 {
		// If <catch> is given, it's used to handle the exception.
		defer func() {
			if e := recover(); e != nil {
				catch[0](e)
			}
		}()
	} else {
		// If no <catch> function passed, it filters the exception.
		defer func() {
			recover()
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
