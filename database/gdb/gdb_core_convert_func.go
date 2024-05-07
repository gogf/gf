// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type iUnmarshalValue interface {
	UnmarshalValue(any) error
}

var convertFuncMap map[reflect.Kind]fieldConvertFunc

var (
	bytesType           = reflect.TypeOf((*[]byte)(nil)).Elem()
	timePtrType         = reflect.TypeOf((*time.Time)(nil))
	timeType            = timePtrType.Elem()
	jsonRawMessageType  = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	jsonUnmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	scannerType         = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	gtimePtrType        = reflect.TypeOf((*gtime.Time)(nil))
	gtimeType           = gtimePtrType.Elem()
)

func init() {
	convertFuncMap = map[reflect.Kind]fieldConvertFunc{
		reflect.Bool:       convertToBool,
		reflect.Int:        convertToInt64,
		reflect.Int8:       convertToInt64,
		reflect.Int16:      convertToInt64,
		reflect.Int32:      convertToInt64,
		reflect.Int64:      convertToInt64,
		reflect.Uint:       convertToUint64,
		reflect.Uint8:      convertToUint64,
		reflect.Uint16:     convertToUint64,
		reflect.Uint32:     convertToUint64,
		reflect.Uint64:     convertToUint64,
		reflect.Uintptr:    convertToUint64,
		reflect.Float32:    convertToFloat64,
		reflect.Float64:    convertToFloat64,
		reflect.Complex64:  nil,
		reflect.Complex128: nil,
		reflect.Array:      nil,
		// reflect.Interface:     convertToInterface,
		reflect.Ptr:           nil,
		reflect.Slice:         convertToJSON,
		reflect.Map:           convertToJSON,
		reflect.Struct:        convertToJSON,
		reflect.String:        convertToString,
		reflect.UnsafePointer: nil,
	}
}

var converterMap sync.Map

func getConverter(typ reflect.Type, n int) fieldConvertFunc {
	// 不支持二级指针赋值
	if n > 1 {
		return nil
	}

	if v, ok := converterMap.Load(typ); ok {
		return v.(fieldConvertFunc)
	}

	fn := getConvertFunc(typ, n)

	if v, ok := converterMap.LoadOrStore(typ, fn); ok {
		return v.(fieldConvertFunc)
	}
	return fn
}

func getBitConvertFunc(typ reflect.Type, n int) fieldConvertFunc {
	if n > 1 {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		fn := getBitConvertFunc(typ.Elem(), n+1)
		if fn != nil {
			return ptrConverter(fn)
		}
		return nil
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dst reflect.Value, src any) error {
			return bitsToNumber[int64](dst.SetInt, src)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dst reflect.Value, src any) error {
			return bitsToNumber[uint64](dst.SetUint, src)
		}
	case reflect.Float32, reflect.Float64:
		return func(dst reflect.Value, src any) error {
			return bitsToNumber[float64](dst.SetFloat, src)
		}
	default:
		return nil

	}
}

// [0xff,0xff] => 0xfffff
func bitsToNumber[T int64 | uint64 | float64](setFn func(T), src any) error {
	switch src := src.(type) {
	case nil:
		setFn(0)
		return nil
	case int64:
		setFn(T(src))
		return nil
	case uint64:
		setFn(T(src))
		return nil
	case []byte:
		toUint64 := bitArrayToUint64(src)
		setFn(T(toUint64))
		return nil
	case string:
		toUint64 := bitArrayToUint64(unsafeStringToBytes(src))
		setFn(T(toUint64))
		return nil
	default:
		return convertError()
	}
}

func bitArrayToUint64(b []byte) uint64 {
	var n uint64
	for _, v := range b {
		n = n<<8 | uint64(v)
	}
	return n
}

func getDecimalConvertFunc(typ reflect.Type, n int) fieldConvertFunc {
	if n > 1 {
		return nil
	}
	if typ.Kind() == reflect.Ptr {
		fn := getDecimalConvertFunc(typ.Elem(), n+1)
		if fn != nil {
			return ptrConverter(fn)
		}
		return nil
	}
	switch typ.Kind() {
	case reflect.Float32, reflect.Float64:
		return convertToFloat64
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dst reflect.Value, src any) error {
			return decimalConvertFunc[int64](dst.SetInt, src)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dst reflect.Value, src any) error {
			return decimalConvertFunc[uint64](dst.SetUint, src)
		}
	default:
		return nil
	}
}

func decimalConvertFunc[T int64 | uint64](set func(T), src any) error {
	switch sv := src.(type) {
	case []byte:
		val, err := strconv.ParseFloat(unsafeBytesToString(sv), 64)
		if err != nil {
			return err
		}
		set(T(val))
	case string:
		val, err := strconv.ParseFloat(sv, 64)
		if err != nil {
			return err
		}
		set(T(val))
	//case float32:
	//	dest.SetUint(uint64(sv))
	case float64:
		set(T(sv))
	default:
		return convertError()
	}
	return nil
}

func ptrConverter(fn fieldConvertFunc) fieldConvertFunc {
	return func(dest reflect.Value, src interface{}) error {
		// 如果从数据库查询出来的是null值
		if src == nil {
			//if !dest.CanAddr() {
			//	if dest.IsNil() {
			//		return nil
			//	}
			//	// 如果结构体字段不是nil
			//	return fn(dest.Elem(), src)
			//}
			if !dest.IsNil() {
				// 如果结构体字段不是nil，重新初始化
				dest.Set(reflect.New(dest.Type().Elem()))
			}
			return nil
		}

		if dest.IsNil() {
			dest.Set(reflect.New(dest.Type().Elem()))
		}

		if dest.Kind() == reflect.Map {
			return fn(dest, src)
		}
		return fn(dest.Elem(), src)
	}
}

func addrConverter(fn fieldConvertFunc) fieldConvertFunc {
	return func(dest reflect.Value, src interface{}) error {
		return fn(dest.Addr(), src)
	}
}

func getConvertFunc(typ reflect.Type, n int) fieldConvertFunc {
	kind := typ.Kind()

	if kind == reflect.Ptr {
		if fn := getConverter(typ.Elem(), n+1); fn != nil {
			return ptrConverter(fn)
		}
	}

	switch typ {
	case bytesType:
		return convertToBytes
	case timeType:
		return convertToTime
	case gtimeType:
		// 对于gtime.Time类型，不走Scan方法的逻辑
		return convertToGTime
	case jsonRawMessageType:
		return convertToBytes
	}

	fn := checkInterface(typ)
	if fn != nil {
		return fn
	}

	// 如果是底层是[]byte 类型的,比如 type MyBytes []byte
	if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8 {
		return convertToBytes
	}

	return convertFuncMap[kind]
}

func checkInterface(typ reflect.Type) fieldConvertFunc {
	isptr := false
	ptrType := typ
	if typ.Kind() == reflect.Ptr {
		isptr = true
	} else {
		ptrType = reflect.PointerTo(typ)
	}

	// 1.如果实现了sql.Scanner接口
	if ptrType.Implements(scannerType) {
		if isptr == false {
			return addrConverter(convertToScanner)
		}
		return convertToScanner
	}
	// 2.如果实现了json.Unmarshal接口
	if typ.Implements(jsonUnmarshalerType) {
		if isptr == false {
			return addrConverter(convertToJsonUnmarshal)
		}
		return convertToJsonUnmarshal
	}
	return nil
}

func convertToBool(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case nil:
		dest.SetBool(false)
		return nil
	case bool:
		dest.SetBool(src)
		return nil
	case int64:
		dest.SetBool(src != 0)
		return nil
	case []byte:
		f, err := strconv.ParseBool(unsafeBytesToString(src))
		if err != nil {
			return err
		}
		dest.SetBool(f)
		return nil
	case string:
		f, err := strconv.ParseBool(src)
		if err != nil {
			return err
		}
		dest.SetBool(f)
		return nil
	default:
		return convertError()
	}
}

func convertToInt64(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		dest.SetInt(0)
		return nil
	case int64:
		dest.SetInt(sv)
		return nil
	case uint64:
		dest.SetInt(int64(sv))
		return nil
	case []byte:
		n, err := strconv.ParseInt(unsafeBytesToString(sv), 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(n)
		return nil
	case string:
		n, err := strconv.ParseInt(sv, 10, 64)
		if err != nil {
			return err
		}
		dest.SetInt(n)
		return nil
	case int8: // dm tinyint
		dest.SetInt(int64(sv))
		return nil
	case int32: // dm tinyint
		dest.SetInt(int64(sv))
		return nil
	default:
		return convertError()
	}
}

func convertToUint64(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case nil:
		dest.SetUint(0)
		return nil
	case uint64:
		dest.SetUint(src)
		return nil
	case int64:
		dest.SetUint(uint64(src))
		return nil
	case []byte:
		n, err := strconv.ParseUint(unsafeBytesToString(src), 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
		return nil
	case string:
		n, err := strconv.ParseUint(src, 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
		return nil
	default:
		return convertError()
	}
}

func convertToFloat64(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case nil:
		dest.SetFloat(0)
		return nil
	case float64:
		dest.SetFloat(src)
		return nil
	case []byte:
		f, err := strconv.ParseFloat(unsafeBytesToString(src), 64)
		if err != nil {
			return err
		}
		dest.SetFloat(f)
		return nil
	case string:
		f, err := strconv.ParseFloat(src, 64)
		if err != nil {
			return err
		}
		dest.SetFloat(f)
		return nil
	default:
		return convertError()
	}
}

func convertToString(dest reflect.Value, src interface{}) error {
	if src == nil {
		return nil
	}
	switch src := src.(type) {
	case nil:
		dest.SetString("")
		return nil
	case string:
		dest.SetString(src)
		return nil
	case []byte:
		dest.SetString(string(src))
		return nil
	case time.Time:
		dest.SetString(src.Format(time.RFC3339Nano))
		return nil
	case int64:
		dest.SetString(strconv.FormatInt(src, 10))
		return nil
	case uint64:
		dest.SetString(strconv.FormatUint(src, 10))
		return nil
	case float64:
		dest.SetString(strconv.FormatFloat(src, 'G', -1, 64))
		return nil
	default:
		return convertError()
	}
}

func convertToBytes(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case nil:
		dest.SetBytes(nil)
		return nil
	case string:
		dest.SetBytes([]byte(src))
		return nil
	case []byte:
		clone := make([]byte, len(src))
		// 必须要调用copy函数，不然的话，保存的引用
		// 等下一次查询的时候，值会被覆盖
		copy(clone, src)
		dest.SetBytes(clone)
		return nil
	default:
		return convertError()
	}
}

func convertToTime(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case nil:
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = time.Time{}
		return nil
	case time.Time:
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = src
		return nil
	case string:
		srcTime, err := parseTime(src)
		if err != nil {
			return err
		}
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = srcTime
		return nil
	case []byte:
		srcTime, err := parseTime(unsafeBytesToString(src))
		if err != nil {
			return err
		}
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = srcTime
		return nil
	default:
		return convertError()
	}
}

func convertToGTime(dest reflect.Value, src interface{}) error {
	switch src := src.(type) {
	case time.Time:
		destTime := dest.Addr().Interface().(*gtime.Time)
		destTime.Time = src
		return nil
	case string:
		srcTime, err := parseTime(src)
		if err != nil {
			return err
		}
		destTime := dest.Addr().Interface().(*gtime.Time)
		destTime.Time = srcTime
		return nil
	case []byte:
		srcTime, err := parseTime(unsafeBytesToString(src))
		if err != nil {
			return err
		}
		destTime := dest.Addr().Interface().(*gtime.Time)
		destTime.Time = srcTime
		return nil
	default:
		return convertError()
	}
}

func convertToScanner(dest reflect.Value, src interface{}) error {
	return dest.Interface().(sql.Scanner).Scan(src)
}

func convertToJsonUnmarshal(dest reflect.Value, src interface{}) error {
	if src == nil {
		return nil
	}
	b, err := toBytes(src)
	if err != nil {
		return err
	}
	return dest.Interface().(json.Unmarshaler).UnmarshalJSON(b)
}

func convertToJSON(dest reflect.Value, src interface{}) error {
	if src == nil {
		return nil
	}
	b, err := toBytes(src)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, dest.Addr().Interface())
}

// 代表不支持数据表字段的类型转换到go类型
// 具体的错误信息在上层的Scan那里抛出
func convertError() error {
	return fmt.Errorf("unsupported types")
}
