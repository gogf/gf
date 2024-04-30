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

func getFloatConvertFunc(fieldType reflect.Type) fieldScanFunc {

	var convert func(originType, elemTyp reflect.Type) fieldScanFunc

	convert = func(originType, elemTyp reflect.Type) fieldScanFunc {

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
			panic(fmt.Errorf("不支持从float类型转换到%v", fieldType))
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
	parseFunc func(val string, base int, bitSize int) (T, error)) fieldScanFunc {
	//
	convert := func(originType, elemType reflect.Type, parseFunc func(val string, base int, bitSize int) (T, error)) fieldScanFunc {
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
			panic(fmt.Errorf("不支持从int类型转换到%v", fieldType))
		}
	}

	switch fieldType.Kind() {
	case reflect.Pointer:
		return convert(fieldType, fieldType.Elem(), parseFunc)
	default:
		return convert(fieldType, fieldType, parseFunc)
	}

}

func getStringConvertFunc(fieldType reflect.Type) fieldScanFunc {

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
		// 检查长度是否相同，如果不同，则不转换
	case reflect.Ptr:
		if fieldType.Elem().Kind() == reflect.String {
			return convertString
		}
	default:
		// convert to int
	}
	// 支持将字符串转换到结构体，切片 map，符合json格式的数据即可
	convertJson := getJsonConvertFunc(fieldType, false)
	if convertJson != nil {
		return convertJson
	}

	return nil
}

func getBoolConvertFunc(fieldType reflect.Type) fieldScanFunc {

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
	panic(fmt.Errorf("不支持从bool类型转换到%v", fieldType))
}

func getTimeConvertFunc(fieldType reflect.Type) fieldScanFunc {
	// 可能不同数据库存储的时间格式不一样
	switch fieldType.String() {
	case "*time.Time":
		return func(val any, dst reflect.Value) (err error) {
			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				// mysql time(10:01:01)类型会被转换为[]byte
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			ptr := dst.Interface().(*time.Time)
			*ptr = t
			// dst.Set(reflect.ValueOf(&t))
			return nil
		}
	case "time.Time":
		return func(val any, dst reflect.Value) (err error) {

			var t time.Time
			switch v := val.(type) {
			case time.Time:
				t = v
			case string:
				// mysql time(10:01:01)类型会被转换为[]byte
				t, err = parseTime(v)
				if err != nil {
					return err
				}
			}
			ptr := dst.Addr().Interface().(*time.Time)
			*ptr = t
			// dst.Set(reflect.ValueOf(t))
			return nil
		}
	default:
		if fieldType.Kind() == reflect.String {
			// TODO 可以格式化一下
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
	panic(fmt.Errorf("不支持从time类型转换到%v", fieldType))
}

func getDecimalConvertFunc(fieldType reflect.Type) fieldScanFunc {
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
		// todo 是否需要支持第三方库的decimal 类型
		panic(fmt.Errorf("不支持从decimal类型转换到%v", fieldType))
	}
}

// bit类型可以转换到任意整数类型，bool类型，string，[]byte,
func getBitConvertFunc(fieldType reflect.Type) fieldScanFunc {
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
	// case reflect.Float32, reflect.Float64: // todo 是否支持bit转换到float?

	case reflect.String: // todo 如果是字符串，是否需要将 1 转为'1' ？
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
		panic(fmt.Errorf("不支持从bit类型转换到%v", fieldType))
	case reflect.Bool:
		return func(val any, dst reflect.Value) error {
			v, _ := val.(string)
			if len(v) > 0 {
				dst.SetBool(v[len(v)-1] == 0)
			}
			return nil
		}
	default:
		panic(fmt.Errorf("不支持从bit类型转换到%v", fieldType))
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

// json 可以转换到以下类型之一,调用标准库的json.Marshal
// struct *struct
// map[string]any
// []int, []string 所有基础类型的切片
// []struct []*struct
func getJsonConvertFunc(fieldType reflect.Type, errPanic bool) fieldScanFunc {

	convertJson := func(val any, dst reflect.Value) error {
		v, _ := val.(string)
		if v == "" {
			return nil
		}
		if dst.Kind() == reflect.Ptr {
			return json.Unmarshal([]byte(v), dst.Interface())
		}
		return json.Unmarshal([]byte(v), dst.Addr().Interface())
	}

	check := func(typ reflect.Type) fieldScanFunc {
		switch typ.Kind() {
		case reflect.Map:
			return convertJson
		case reflect.Struct:
			return convertJson
		case reflect.Slice:
			return convertJson
		default:
			if errPanic {
				panic(fmt.Errorf("不支持从json类型转换到%v", fieldType))
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
