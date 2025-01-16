// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"database/sql"
	"encoding/json"
	"errors"
	"reflect"
	"strconv"
	"sync"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

type iUnmarshalValue interface {
	UnmarshalValue(any) error
}

var (
	convertFuncMap map[reflect.Kind]fieldConvertFunc
	converterMap   sync.Map
)

var (
	bytesType           = reflect.TypeOf((*[]byte)(nil)).Elem()
	timePtrType         = reflect.TypeOf((*time.Time)(nil))
	timeType            = timePtrType.Elem()
	jsonRawMessageType  = reflect.TypeOf((*json.RawMessage)(nil)).Elem()
	jsonUnmarshalerType = reflect.TypeOf((*json.Unmarshaler)(nil)).Elem()
	unmarshalValueType  = reflect.TypeOf((*iUnmarshalValue)(nil)).Elem()
	scannerType         = reflect.TypeOf((*sql.Scanner)(nil)).Elem()
	gtimePtrType        = reflect.TypeOf((*gtime.Time)(nil))
	gtimeType           = gtimePtrType.Elem()
)

func init() {
	convertFuncMap = map[reflect.Kind]fieldConvertFunc{
		reflect.Bool:          convertToBool,
		reflect.Int:           convertToInt64,
		reflect.Int8:          convertToInt64,
		reflect.Int16:         convertToInt64,
		reflect.Int32:         convertToInt64,
		reflect.Int64:         convertToInt64,
		reflect.Uint:          convertToUint64,
		reflect.Uint8:         convertToUint64,
		reflect.Uint16:        convertToUint64,
		reflect.Uint32:        convertToUint64,
		reflect.Uint64:        convertToUint64,
		reflect.Uintptr:       convertToUint64,
		reflect.Float32:       convertToFloat64,
		reflect.Float64:       convertToFloat64,
		reflect.Complex64:     nil,
		reflect.Complex128:    nil,
		reflect.Array:         nil,
		reflect.String:        convertToString,
		reflect.Ptr:           nil,
		reflect.UnsafePointer: nil,
	}
}

// Mainly the bit type of MySQL
func getBitConvertFunc(typ reflect.Type, deref int) fieldConvertFunc {
	// Supports up to one level of pointers
	if deref > 1 {
		return nil
	}
	// When entering this method for the first time
	if deref == 0 {
		fn := checkTypeImplSqlScanner(typ)
		if fn != nil {
			if typ.Kind() == reflect.Ptr {
				return ptrSqlScannerConvert
			}
			return fn
		}
	}

	if typ.Kind() == reflect.Ptr {
		fn := getBitConvertFunc(typ.Elem(), deref+1)
		if fn != nil {
			return ptrConverter(fn)
		}
		return nil
	}

	switch typ.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return func(dst reflect.Value, src any) error {
			return bitArrayConvertToNumber[int64](dst.SetInt, src)
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return func(dst reflect.Value, src any) error {
			return bitArrayConvertToNumber[uint64](dst.SetUint, src)
		}
	case reflect.Float32, reflect.Float64:
		return func(dst reflect.Value, src any) error {
			return bitArrayConvertToNumber[float64](dst.SetFloat, src)
		}
	case reflect.String:
		return convertToString
	default:
		if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8 {
			// []byte
			return convertToBytes
		}
		return nil
	}
}

// []byte{0xff,0xff} => 0xffff
func bitArrayConvertToNumber[T int64 | uint64 | float64](setFn func(value T), src any) error {
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

func getDecimalConvertFunc(typ reflect.Type, deref int) fieldConvertFunc {
	// Supports up to one level of pointers
	if deref > 1 {
		return nil
	}
	// When entering this method for the first time
	if deref == 0 {
		fn := checkTypeImplSqlScanner(typ)
		if fn != nil {
			if typ.Kind() == reflect.Ptr {
				return ptrSqlScannerConvert
			}
			return fn
		}
	}

	if typ.Kind() == reflect.Ptr {
		fn := getDecimalConvertFunc(typ.Elem(), deref+1)
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
	case reflect.String:
		return convertToString
	default:
		if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8 {
			// []byte
			return convertToBytes
		}
		return nil
	}
}

func decimalConvertFunc[T int64 | uint64](set func(value T), src any) error {
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
	case float32:
		set(T(sv))
	case float64:
		set(T(sv))
	default:
		return convertError()
	}
	return nil
}

func ptrConverter(fn fieldConvertFunc) fieldConvertFunc {
	return func(dest reflect.Value, src interface{}) error {
		if src == nil {
			if dest.IsNil() == false {
				// If the struct field is not nil, reinitialize
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

// deref represents the number of times the pointer has been dereferenced.
// If it is greater than 1, it returns nil directly
// impl Indicates whether the type implements the following two interfaces
// 1. sql.Scanner
// 2. json.Unmarshaler
func getConverter(typ reflect.Type, deref int) (fn fieldConvertFunc, impl bool) {
	// Secondary pointer assignment is not supported
	if deref > 1 {
		return nil, false
	}
	if v, ok := converterMap.Load(typ); ok {
		return v.(fieldConvertFunc), false
	}
	fn, impl = getConvertFunc(typ, deref)
	if v, ok := converterMap.LoadOrStore(typ, fn); ok {
		return v.(fieldConvertFunc), impl
	}
	return fn, impl
}

// deref represents the number of times the pointer has been dereferenced.
// If it is greater than 1, it returns nil directly
// implInterface Indicates whether the type implements the following two interfaces, and if so,
// it will not be wrapped in any other way, and will be returned directly
// 1. sql.Scanner
// 2. json.Unmarshaler
func getConvertFunc(typ reflect.Type, deref int) (fn fieldConvertFunc, implInterface bool) {
	kind := typ.Kind()
	if deref == 0 {
		fn = checkTypeImplSqlScanner(typ)
		if fn != nil {
			if kind == reflect.Ptr {
				return ptrSqlScannerConvert, true
			}
			return fn, true
		}
	}

	if kind == reflect.Ptr {
		fn, implInterface = getConverter(typ.Elem(), deref+1)
		if fn != nil {
			// If the type implements the [sql.Scanner] or [json.Unmarshaler] interface,
			// there is no need for [ptrConverter] to wrap it once
			if implInterface {
				return fn, true
			}
			// If those two interfaces are not implemented, [ptrConverter] needs to be wrapped and returned once,
			// For example, if the database driver returns an `int`, but the structure field type `*int`,
			// it will return `ptrConverter(convertToInt64)`
			// After the initialization of the pointer is completed in ptrConverter,
			// the callback convertToInt64 completes the final assignment
			return ptrConverter(fn), false
		}
	}

	switch typ {
	case bytesType:
		return convertToBytes, false
	case timeType:
		return convertToTime, false
	case jsonRawMessageType:
		return convertToBytes, false
	}

	if typ.Kind() == reflect.Slice && typ.Elem().Kind() == reflect.Uint8 {
		return convertToBytes, false
	}
	// Covers, Prop, Sku, regardless of whether these three types are struct, slice, map, or other custom types,
	// If the [json.Unmarshaler] interface is implemented, the user-defined UnmarshalJSON method will be called
	// tip:  But there is an exception to custom types, which is the custom []byte type.
	//       Even if the [json.Unmarshaler] interface is implemented, it will not be called
	// If not implemented, it is necessary to ensure that the database driver returns either a `[]byte` or a `string` type
	// If it is not a `[]byte` or a `string` type, it will return an error
	//
	//  type MyBytes []byte // Assuming it implements [json.Unmarshaler], it will not call
	//
	// 	type GiftEntity struct {
	//		Covers    Covers
	//		Props     []Prop
	//		Skus      map[int]Sku
	//	}
	switch typ.Kind() {
	case reflect.Slice, reflect.Map, reflect.Struct:
		fn := checkTypeImplJsonUnmarshaler(typ)
		// If the type implements json.Unmarshaler interface
		if fn != nil {
			if deref > 0 {
				return ptrUnmarshalJsonConvert, true
			}
			return unmarshalJsonConvert, true
		}
		return convertToJSON, false
	}

	return convertFuncMap[kind], false
}

func checkTypeImplSqlScanner(typ reflect.Type) fieldConvertFunc {
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)
	}
	if typ.Implements(scannerType) {
		return sqlScannerConvert
	}
	return nil
}

func ptrSqlScannerConvert(dest reflect.Value, src interface{}) error {
	typ := dest.Type()
	if src == nil {
		if dest.IsNil() == false {
			// If the struct field is not nil, reinitialize
			dest.Set(reflect.New(typ.Elem()))
		}
		return nil
	}
	if dest.IsNil() {
		dest.Set(reflect.New(typ.Elem()))
	}
	return dest.Interface().(sql.Scanner).Scan(src)
}

func sqlScannerConvert(dest reflect.Value, src interface{}) error {
	return dest.Addr().Interface().(sql.Scanner).Scan(src)
}

func convertToBool(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		dest.SetBool(false)
		return nil
	case bool:
		dest.SetBool(sv)
		return nil
	case int64:
		dest.SetBool(sv != 0)
		return nil
	case []byte:
		f, err := strconv.ParseBool(unsafeBytesToString(sv))
		if err != nil {
			return err
		}
		dest.SetBool(f)
		return nil
	case string:
		f, err := strconv.ParseBool(sv)
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
	default:
		// Most drivers will go to the top branch,
		// and only a few drivers will go to the bottom branch
		switch sv := src.(type) {
		case int8: // dm tinyint
			dest.SetInt(int64(sv))
			return nil
		case int16: // dm smallint
			dest.SetInt(int64(sv))
			return nil
		case int32: // dm int
			dest.SetInt(int64(sv))
			return nil
		}
		return convertError()
	}
}

func convertToUint64(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		dest.SetUint(0)
		return nil
	case uint64:
		dest.SetUint(sv)
		return nil
	case int64:
		dest.SetUint(uint64(sv))
		return nil
	case []byte:
		n, err := strconv.ParseUint(unsafeBytesToString(sv), 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
		return nil
	case string:
		n, err := strconv.ParseUint(sv, 10, 64)
		if err != nil {
			return err
		}
		dest.SetUint(n)
		return nil
	default:
		// Most drivers will go to the top branch,
		// and only a few drivers will go to the bottom branch
		// clickhouse
		switch sv := src.(type) {
		case uint8:
			dest.SetUint(uint64(sv))
			return nil
		case uint16:
			dest.SetUint(uint64(sv))
			return nil
		case uint32:
			dest.SetUint(uint64(sv))
			return nil
		}
		return convertError()
	}
}

func convertToFloat64(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		dest.SetFloat(0)
		return nil
	case float64:
		dest.SetFloat(sv)
		return nil
	case float32:
		dest.SetFloat(float64(sv))
		return nil
	case []byte:
		f, err := strconv.ParseFloat(unsafeBytesToString(sv), 64)
		if err != nil {
			return err
		}
		dest.SetFloat(f)
		return nil
	case string:
		f, err := strconv.ParseFloat(sv, 64)
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
	switch sv := src.(type) {
	case nil:
		dest.SetString("")
		return nil
	case string:
		dest.SetString(sv)
		return nil
	case []byte:
		dest.SetString(string(sv))
		return nil
	case time.Time:
		dest.SetString(sv.Format(time.RFC3339Nano))
		return nil
	case int64:
		dest.SetString(strconv.FormatInt(sv, 10))
		return nil
	case uint64:
		dest.SetString(strconv.FormatUint(sv, 10))
		return nil
	case float64:
		dest.SetString(strconv.FormatFloat(sv, 'G', -1, 64))
		return nil
	default:
		// Most drivers will go to the top branch,
		// and only a few drivers will go to the bottom branch
		switch sv := src.(type) {
		case int8:
			dest.SetString(strconv.FormatInt(int64(sv), 10))
			return nil
		case int16:
			dest.SetString(strconv.FormatInt(int64(sv), 10))
			return nil
		case int32:
			dest.SetString(strconv.FormatInt(int64(sv), 10))
			return nil
		case uint8:
			dest.SetString(strconv.FormatUint(uint64(sv), 10))
			return nil
		case uint16:
			dest.SetString(strconv.FormatUint(uint64(sv), 10))
			return nil
		case uint32:
			dest.SetString(strconv.FormatUint(uint64(sv), 10))
			return nil
		case float32:
			dest.SetString(strconv.FormatFloat(float64(sv), 'G', -1, 64))
			return nil
		}
		return convertError()
	}
}

func convertToBytes(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		dest.SetBytes(nil)
		return nil
	case string:
		dest.SetBytes([]byte(sv))
		return nil
	case []byte:
		clone := make([]byte, len(sv))
		// The copy function must be called, otherwise the reference is saved
		// The next time you query, the value will be overwritten
		copy(clone, sv)
		dest.SetBytes(clone)
		return nil
	default:
		return convertError()
	}
}

func convertToTime(dest reflect.Value, src interface{}) error {
	switch sv := src.(type) {
	case nil:
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = time.Time{}
		return nil
	case time.Time:
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = sv
		return nil
	case string:
		srcTime, err := parseTime(sv)
		if err != nil {
			return err
		}
		destTime := dest.Addr().Interface().(*time.Time)
		*destTime = srcTime
		return nil
	case []byte:
		srcTime, err := parseTime(unsafeBytesToString(sv))
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

func checkTypeImplJsonUnmarshaler(typ reflect.Type) fieldConvertFunc {
	if typ.Kind() != reflect.Ptr {
		typ = reflect.PointerTo(typ)
	}
	if typ.Implements(jsonUnmarshalerType) {
		return unmarshalJsonConvert
	}
	return nil
}

func ptrUnmarshalJsonConvert(dest reflect.Value, src any) error {
	if src == nil {
		if dest.IsNil() == false {
			// If the struct field is not nil, reinitialize
			dest.Set(reflect.New(dest.Type().Elem()))
		}
		return nil
	}
	b, err := toBytes(src)
	if err != nil {
		return err
	}
	if dest.IsNil() {
		dest.Set(reflect.New(dest.Type().Elem()))
	}
	return dest.Interface().(json.Unmarshaler).UnmarshalJSON(b)
}

func unmarshalJsonConvert(dest reflect.Value, src interface{}) error {
	if src == nil {
		return nil
	}
	b, err := toBytes(src)
	if err != nil {
		return err
	}
	return dest.Addr().Interface().(json.Unmarshaler).UnmarshalJSON(b)
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

// Indicates that the type of a data table field is not supported and the type is converted to the go type
// The specific error message is thrown at the upper Scan
func convertError() error {
	return errors.New("unsupported types")
}
