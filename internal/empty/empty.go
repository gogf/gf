// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package empty provides functions for checking empty/nil variables.
package empty

import (
	"reflect"
	"time"
)

// iString is used for type assert api for String().
type iString interface {
	String() string
}

// iInterfaces is used for type assert api for Interfaces.
type iInterfaces interface {
	Interfaces() []interface{}
}

// iMapStrAny is the interface support for converting struct parameter to map.
type iMapStrAny interface {
	MapStrAny() map[string]interface{}
}

type iTime interface {
	Date() (year int, month time.Month, day int)
	IsZero() bool
}

// IsEmpty checks whether given `value` empty.
// It returns true if `value` is in: 0, nil, false, "", len(slice/map/chan) == 0,
// or else it returns false.
func IsEmpty(value interface{}) bool {
	if value == nil {
		return true
	}
	// It firstly checks the variable as common types using assertion to enhance the performance,
	// and then using reflection.
	switch result := value.(type) {
	case int:
		return result == 0
	case int8:
		return result == 0
	case int16:
		return result == 0
	case int32:
		return result == 0
	case int64:
		return result == 0
	case uint:
		return result == 0
	case uint8:
		return result == 0
	case uint16:
		return result == 0
	case uint32:
		return result == 0
	case uint64:
		return result == 0
	case float32:
		return result == 0
	case float64:
		return result == 0
	case bool:
		return result == false
	case string:
		return result == ""
	case []byte:
		return len(result) == 0
	case []rune:
		return len(result) == 0
	case []int:
		return len(result) == 0
	case []string:
		return len(result) == 0
	case []float32:
		return len(result) == 0
	case []float64:
		return len(result) == 0
	case map[string]interface{}:
		return len(result) == 0

	default:
		// =========================
		// Common interfaces checks.
		// =========================
		if f, ok := value.(iTime); ok {
			if f == nil {
				return true
			}
			return f.IsZero()
		}
		if f, ok := value.(iString); ok {
			if f == nil {
				return true
			}
			return f.String() == ""
		}
		if f, ok := value.(iInterfaces); ok {
			if f == nil {
				return true
			}
			return len(f.Interfaces()) == 0
		}
		if f, ok := value.(iMapStrAny); ok {
			if f == nil {
				return true
			}
			return len(f.MapStrAny()) == 0
		}
		// Finally, using reflect.
		var rv reflect.Value
		if v, ok := value.(reflect.Value); ok {
			rv = v
		} else {
			rv = reflect.ValueOf(value)
		}

		switch rv.Kind() {
		case reflect.Bool:
			return !rv.Bool()

		case
			reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64:
			return rv.Int() == 0

		case
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64,
			reflect.Uintptr:
			return rv.Uint() == 0

		case
			reflect.Float32,
			reflect.Float64:
			return rv.Float() == 0

		case reflect.String:
			return rv.Len() == 0

		case reflect.Struct:
			for i := 0; i < rv.NumField(); i++ {
				if !IsEmpty(rv.Field(i).Interface()) {
					return false
				}
			}
			return true

		case
			reflect.Chan,
			reflect.Map,
			reflect.Slice,
			reflect.Array:
			return rv.Len() == 0

		case
			reflect.Func,
			reflect.Ptr,
			reflect.Interface,
			reflect.UnsafePointer:
			if rv.IsNil() {
				return true
			}
		}
	}
	return false
}

// IsEmptyLength checks whether given `value` is empty length.
// It returns true if `value` is in: nil, "", len(slice/map/chan) == 0,
// or else it returns false.
//func IsEmptyLength(value interface{}) bool {
//	if value == nil {
//		return true
//	}
//	// It firstly checks the variable as common types using assertion to enhance the performance,
//	// and then using reflection.
//	switch value := value.(type) {
//	case
//		int,
//		int8,
//		int16,
//		int32,
//		int64,
//		uint,
//		uint8,
//		uint16,
//		uint32,
//		uint64,
//		float32,
//		float64,
//		bool:
//		return false
//	case string:
//		return value == ""
//	case []byte:
//		return len(value) == 0
//	case []rune:
//		return len(value) == 0
//	case []int:
//		return len(value) == 0
//	case []string:
//		return len(value) == 0
//	case []float32:
//		return len(value) == 0
//	case []float64:
//		return len(value) == 0
//	case map[string]interface{}:
//		return len(value) == 0
//	default:
//		// =========================
//		// Common interfaces checks.
//		// =========================
//		if f, ok := value.(iTime); ok {
//			if f == nil {
//				return true
//			}
//			return f.IsZero()
//		}
//		if f, ok := value.(iString); ok {
//			if f == nil {
//				return true
//			}
//			return f.String() == ""
//		}
//		if f, ok := value.(iInterfaces); ok {
//			if f == nil {
//				return true
//			}
//			return len(f.Interfaces()) == 0
//		}
//		if f, ok := value.(iMapStrAny); ok {
//			if f == nil {
//				return true
//			}
//			return len(f.MapStrAny()) == 0
//		}
//		// Finally using reflect.
//		var rv reflect.Value
//		if v, ok := value.(reflect.Value); ok {
//			rv = v
//		} else {
//			rv = reflect.ValueOf(value)
//		}
//
//		switch rv.Kind() {
//		case
//			reflect.Int,
//			reflect.Int8,
//			reflect.Int16,
//			reflect.Int32,
//			reflect.Int64,
//			reflect.Uint,
//			reflect.Uint8,
//			reflect.Uint16,
//			reflect.Uint32,
//			reflect.Uint64,
//			reflect.Uintptr,
//			reflect.Float32,
//			reflect.Float64,
//			reflect.Bool:
//			return false
//		case reflect.String:
//			return rv.Len() == 0
//		case reflect.Struct:
//			for i := 0; i < rv.NumField(); i++ {
//				if !IsEmpty(rv) {
//					return false
//				}
//			}
//			return true
//		case reflect.Chan,
//			reflect.Map,
//			reflect.Slice,
//			reflect.Array:
//			return rv.Len() == 0
//		case reflect.Func,
//			reflect.Ptr,
//			reflect.Interface,
//			reflect.UnsafePointer:
//			if rv.IsNil() {
//				return true
//			}
//		}
//	}
//	return false
//}

// IsNil checks whether given `value` is nil.
// Parameter `traceSource` is used for tracing to the source variable if given `value` is type of pinter
// that also points to a pointer. It returns nil if the source is nil when `traceSource` is true.
// Note that it might use reflect feature which affects performance a little.
func IsNil(value interface{}, traceSource ...bool) bool {
	if value == nil {
		return true
	}
	var rv reflect.Value
	if v, ok := value.(reflect.Value); ok {
		rv = v
	} else {
		rv = reflect.ValueOf(value)
	}
	switch rv.Kind() {
	case reflect.Chan,
		reflect.Map,
		reflect.Slice,
		reflect.Func,
		reflect.Interface,
		reflect.UnsafePointer:
		return !rv.IsValid() || rv.IsNil()

	case reflect.Ptr:
		if len(traceSource) > 0 && traceSource[0] {
			for rv.Kind() == reflect.Ptr {
				rv = rv.Elem()
			}
			if !rv.IsValid() {
				return true
			}
			if rv.Kind() == reflect.Ptr {
				return rv.IsNil()
			}
		} else {
			return !rv.IsValid() || rv.IsNil()
		}
	}
	return false
}
