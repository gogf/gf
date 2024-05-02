// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

func parseFloat[T int64 | uint64 | float64](typ reflect.Type, val string, dst reflect.Value, setFn func(T)) error {
	n, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return err
	}
	if dst.Kind() == reflect.Ptr {
		if dst.IsNil() {
			dst.Set(reflect.New(typ))
		}
		dst = dst.Elem()
	}
	setFn(T(n))
	return nil
}

func getFloatConvertFunc(fieldType reflect.Type) fieldConvertFunc {

	convert := func(originType, elemTyp reflect.Type) fieldConvertFunc {
		switch elemTyp.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return func(val any, dst reflect.Value) error {
				return parseFloat[int64](originType, val.(string), dst, dst.SetInt)
			}
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return func(val any, dst reflect.Value) error {
				return parseFloat[uint64](originType, val.(string), dst, dst.SetUint)
			}
		case reflect.Float32, reflect.Float64:
			return func(val any, dst reflect.Value) error {
				return parseFloat[float64](originType, val.(string), dst, dst.SetFloat)

			}
		case reflect.String:
			return func(val any, dst reflect.Value) error {
				dst.SetString(val.(string))
				return nil
			}
		default:
			return nil
		}
	}
	switch fieldType.Kind() {
	case reflect.Pointer:
		return convert(fieldType, fieldType.Elem())
	default:
		return convert(fieldType, fieldType)
	}
}

func getIntegerConvertFunc[T int64 | uint64](fieldType reflect.Type,
	parseFunc func(val string, base int, bitSize int) (T, error)) fieldConvertFunc {
	//
	convert := func(originType, elemType reflect.Type, parseFunc func(val string, base int, bitSize int) (T, error)) fieldConvertFunc {
		switch elemType.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return func(val any, dst reflect.Value) error {
				n, err := parseFunc(val.(string), 10, 64)
				if err != nil {
					return err
				}
				if dst.Kind() == reflect.Ptr {
					if dst.IsNil() {
						dst.Set(reflect.New(elemType))
					}
					dst = dst.Elem()
				}
				dst.SetInt(int64(n))
				return nil
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			return func(val any, dst reflect.Value) error {
				n, err := parseFunc(val.(string), 10, 64)
				if err != nil {
					return err
				}
				if dst.Kind() == reflect.Ptr {
					if dst.IsNil() {
						dst.Set(reflect.New(elemType))
					}
					dst = dst.Elem()
				}
				dst.SetUint(uint64(n))
				return nil
			}
		case reflect.Float32, reflect.Float64:
			return func(val any, dst reflect.Value) error {
				n, err := parseFunc(val.(string), 10, 64)
				if err != nil {
					return err
				}
				if dst.Kind() == reflect.Ptr {
					if dst.IsNil() {
						dst.Set(reflect.New(elemType))
					}
					dst = dst.Elem()
				}
				dst.SetFloat(float64(n))
				return nil
			}
		case reflect.String:
			return func(val any, dst reflect.Value) error {
				dst.SetString(val.(string))
				return nil
			}
		case reflect.Bool:
			return func(val any, dst reflect.Value) error {
				b, err := strconv.ParseBool(val.(string))
				if err != nil {
					return err
				}
				dst.SetBool(b)
				return nil
			}
		default:
			return nil
		}
	}

	switch fieldType.Kind() {
	case reflect.Pointer:
		return convert(fieldType, fieldType.Elem(), parseFunc)
	default:
		return convert(fieldType, fieldType, parseFunc)
	}

}

func getStringConvertFunc(fieldType reflect.Type) fieldConvertFunc {

	convertString := func(val any, dst reflect.Value) error {
		if dst.Kind() == reflect.Ptr {
			if dst.IsNil() {
				dst.Set(reflect.New(fieldType.Elem()))
			}
			dst = dst.Elem()
		}
		dst.SetString(val.(string))
		return nil
	}

	switch fieldType.Kind() {
	case reflect.String:
		return convertString
	case reflect.Slice:
		fieldType = fieldType.Elem()
		// []uint8 []byte
		if fieldType.Kind() == reflect.Uint8 {
			return func(val any, dst reflect.Value) error {
				v, _ := val.(string)
				dst.SetBytes([]byte(v))
				return nil
			}
		}
	case reflect.Array:
		// Check if the length is the same, if not, don't convert
	case reflect.Ptr:
		if fieldType.Elem().Kind() == reflect.String {
			return convertString
		}
	default:
		// convert to int
	}
	// Support converting strings to structs, slicing maps, and data in JSON format
	convertJson := getJsonConvertFunc(fieldType, false)
	if convertJson != nil {
		return convertJson
	}
	return nil
}

func getBoolConvertFunc(fieldType reflect.Type) fieldConvertFunc {

	convertBool := func(val any, dst reflect.Value) error {
		b, err := strconv.ParseBool(val.(string))
		if err != nil {
			return err
		}
		if dst.Kind() == reflect.Ptr {
			if dst.IsNil() {
				dst.Set(reflect.New(fieldType.Elem()))
			}
			dst = dst.Elem()
		}
		dst.SetBool(b)
		return nil
	}

	switch fieldType.Kind() {
	case reflect.Bool:
		return convertBool
	case reflect.String:
		return func(val any, dst reflect.Value) error {
			dst.SetString(val.(string))
			return nil
		}
	case reflect.Ptr:
		if fieldType.Elem().Kind() == reflect.Bool {
			return convertBool
		}
	}
	return nil
}

var (
	timeTimeType     = reflect.TypeOf(time.Time{})
	timeTimePtrType  = reflect.TypeOf(&time.Time{})
	gtimeTimeType    = reflect.TypeOf(gtime.Time{})
	gtimeTimePtrType = reflect.TypeOf(&gtime.Time{})
)

func getTimeConvertFunc(fieldType reflect.Type) fieldConvertFunc {
	// The time format may be different for different databases
	switch fieldType {
	case timeTimeType:
		return func(val any, dst reflect.Value) (err error) {
			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			ptr := dst.Addr().Interface().(*time.Time)
			*ptr = t
			return nil
		}
	case timeTimePtrType:
		return func(val any, dst reflect.Value) (err error) {
			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			if dst.IsNil() {
				dst.Set(reflect.New(timeTimeType))
			}
			ptr := dst.Interface().(*time.Time)
			*ptr = t
			return nil
		}
	case gtimeTimeType:
		return func(val any, dst reflect.Value) (err error) {
			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			ptr := dst.Addr().Interface().(*gtime.Time)
			ptr.Time = t
			return nil
		}
	case gtimeTimePtrType:
		return func(val any, dst reflect.Value) (err error) {
			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			if dst.IsNil() {
				dst.Set(reflect.New(gtimeTimeType))
			}
			ptr := dst.Interface().(*gtime.Time)
			ptr.Time = t
			return nil
		}
	default:
		// todo typ.ConvertibleTo(reflect.TypeOf(time.Time{}))？
		if fieldType.Kind() == reflect.String {
			// TODO Does formatting need to be done?
			return func(val any, dst reflect.Value) (err error) {
				switch v := val.(type) {
				case time.Time:
					dst.SetString(v.String())
				case string:
					dst.SetString(v)
				}
				return nil
			}
		}
	}
	return nil
}

func getDecimalConvertFunc(fieldType reflect.Type) fieldConvertFunc {
	switch fieldType.Kind() {
	case reflect.Float64, reflect.Float32:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			dst.SetFloat(n)
			return nil
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			dst.SetInt(int64(n))
			return nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			n, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return err
			}
			dst.SetUint(uint64(n))
			return nil
		}
	case reflect.String:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			dst.SetString(v)
			return nil
		}
	default:
		return nil
	}
}

// Bit types can be converted to arbitrary integer types, Boolean types, Stirling, [] bits,
func getBitConvertFunc(fieldType reflect.Type) fieldConvertFunc {
	switch fieldType.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(val any, dst reflect.Value) error {
			if dst.Kind() == reflect.Ptr {
				if dst.IsNil() {
					dst.Set(reflect.New(fieldType.Elem()))
				}
				dst = dst.Elem()
			}
			v, _ := val.(string)
			dst.SetInt(int64(bitArrayToUint64([]byte(v))))
			return nil
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(val any, dst reflect.Value) error {
			if dst.Kind() == reflect.Ptr {
				if dst.IsNil() {
					dst.Set(reflect.New(fieldType.Elem()))
				}
				dst = dst.Elem()
			}
			v, _ := val.(string)
			dst.SetUint(bitArrayToUint64([]byte(v)))
			return nil
		}
	// case reflect.Float32, reflect.Float64: // todo Does it support bit-to-float conversion?
	case reflect.String:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			dst.SetString(v)
			return nil
		}
	case reflect.Slice:
		// []byte
		if fieldType.Elem().Kind() == reflect.Uint8 {
			return func(val any, dst reflect.Value) error {
				v, _ := val.(string)
				dst.SetBytes([]byte(v))
				return nil
			}
		}
	case reflect.Bool:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			if len(v) > 0 {
				dst.SetBool(v[len(v)-1] == 0)
			}
			return nil
		}
	default:
		return nil
	}
	return nil
}

func bitArrayToUint64(b []byte) uint64 {
	var n uint64
	for _, v := range b {
		n = n<<8 | uint64(v)
	}
	return n
}

// json It can be converted to one of the following types, calling the standard library's json. Marshal
// struct *struct
// map[string]any
// []int, []string ... All base types of slices
// []struct []*struct
func getJsonConvertFunc(fieldType reflect.Type, errPanic bool) fieldConvertFunc {

	impl, fieldIsPtr := checkImplUnmarshalJSON(fieldType)

	convertJson := func(val any, dst reflect.Value) error {
		v, _ := val.(string)
		if v == "" {
			// todo Do I need to initialize a pointer field?
			return nil
		}
		// impl json.Unmarshaler
		if impl {
			if fieldIsPtr {
				if dst.IsNil() {
					dst.Set(reflect.New(fieldType.Elem()))
				}
				return dst.Interface().(json.Unmarshaler).UnmarshalJSON([]byte(v))
			}
			return dst.Addr().Interface().(json.Unmarshaler).UnmarshalJSON([]byte(v))
		}
		if dst.Kind() == reflect.Ptr {
			return json.Unmarshal([]byte(v), dst.Interface())
		}
		return json.Unmarshal([]byte(v), dst.Addr().Interface())
	}

	check := func(typ reflect.Type) fieldConvertFunc {
		switch typ.Kind() {
		case reflect.Map:
			return convertJson
		case reflect.Struct:
			return convertJson
		case reflect.Slice:
			return convertJson
		default:
			if errPanic {
				panic(fmt.Errorf("conversion from JSON type to %v is not supported", fieldType))
			}
		}
		return nil
	}
	switch fieldType.Kind() {
	case reflect.Pointer:
		return check(fieldType.Elem())
	default:
		return check(fieldType)
	}
}

func checkImplUnmarshalJSON(typ reflect.Type) (impl, isptr bool) {
	v := reflect.Value{}
	if typ.Kind() != reflect.Ptr {
		v = reflect.New(typ)
	} else {
		v = reflect.New(typ.Elem())
		isptr = true
	}
	switch v.Interface().(type) {
	case json.Unmarshaler:
		impl = true
	}
	return
}

const (
	dateFormat         = "2006-01-02"
	timeFormat         = "15:04:05.999999999"
	timetzFormat1      = "15:04:05.999999999-07:00:00"
	timetzFormat2      = "15:04:05.999999999-07:00"
	timetzFormat3      = "15:04:05.999999999-07"
	timestampFormat    = "2006-01-02 15:04:05.999999999"
	timestamptzFormat1 = "2006-01-02 15:04:05.999999999-07:00:00"
	timestamptzFormat2 = "2006-01-02 15:04:05.999999999-07:00"
	timestamptzFormat3 = "2006-01-02 15:04:05.999999999-07"
)

func parseTime(s string) (time.Time, error) {
	l := len(s)

	if l >= len("2006-01-02 15:04:05") {
		switch s[10] {
		case ' ':
			if c := s[l-6]; c == '+' || c == '-' {
				return time.Parse(timestamptzFormat2, s)
			}
			if c := s[l-3]; c == '+' || c == '-' {
				return time.Parse(timestamptzFormat3, s)
			}
			if c := s[l-9]; c == '+' || c == '-' {
				return time.Parse(timestamptzFormat1, s)
			}
			return time.ParseInLocation(timestampFormat, s, time.UTC)
		case 'T':
			return time.Parse(time.RFC3339Nano, s)
		}
	}

	if l >= len("15:04:05-07") {
		if c := s[l-6]; c == '+' || c == '-' {
			return time.Parse(timetzFormat2, s)
		}
		if c := s[l-3]; c == '+' || c == '-' {
			return time.Parse(timetzFormat3, s)
		}
		if c := s[l-9]; c == '+' || c == '-' {
			return time.Parse(timetzFormat1, s)
		}
	}

	if l < len("15:04:05") {
		return time.Time{}, fmt.Errorf("can't parse time=%q", s)
	}

	if s[2] == ':' {
		return time.ParseInLocation(timeFormat, s, time.UTC)
	}
	return time.ParseInLocation(dateFormat, s, time.UTC)
}
