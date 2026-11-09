// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql/driver"
	"maps"
	"math"
	"math/big"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/util/gconv"
)

// ConvertValueForField converts value to the type of the record field.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	switch itemValue := fieldValue.(type) {
	case time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}
		return itemValue, checkTimeRange(fieldType, itemValue)

	case uuid.UUID:
		return itemValue, nil

	case *time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		return itemValue, checkTimeRange(fieldType, *itemValue)

	case gtime.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}
		// for gtime type, needs to get time.Time
		return itemValue.Time, checkTimeRange(fieldType, itemValue.Time)

	case *gtime.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		// for gtime type, needs to get time.Time
		return itemValue.Time, checkTimeRange(fieldType, itemValue.Time)

	case decimal.Decimal:
		if _, _, ok := decimalPrecisionScale(fieldType); ok {
			return convertDecimal(fieldType, itemValue)
		}
		return itemValue, nil

	case *decimal.Decimal:
		if itemValue == nil {
			return nil, nil
		}
		if _, _, ok := decimalPrecisionScale(fieldType); ok {
			return convertDecimal(fieldType, *itemValue)
		}
		return *itemValue, nil

	default:
		// if the other type implements valuer for the driver package
		// the converted result is used
		// otherwise the value is converted to the type of the field
		valuer, ok := itemValue.(driver.Valuer)
		if !ok {
			return d.convertValueByFieldType(ctx, fieldType, itemValue)
		}
		convertedValue, err := valuer.Value()
		if err != nil {
			return nil, err
		}
		return convertedValue, nil
	}
}

// convertValueByFieldType converts `value` to the Go type that the underlying driver binds to a
// column of type `fieldType`, as the underlying driver refuses to convert between types itself:
// a string or a number to the number, decimal, boolean or time of the column, a number or a
// boolean to the string of a String column, and a map, a slice or a struct to its JSON in a
// String column as the core converts it. A nil pointer becomes nil, and a value that is not one
// of its column, like a number out of the range of the column or a float with a fraction for an
// integer column, is an error of code gcode.CodeInvalidParameter. Any other value, like a slice
// for an Array column, is returned as it is.
func (d *Driver) convertValueByFieldType(ctx context.Context, fieldType string, value any) (any, error) {
	switch value.(type) {
	case nil, []byte, gdb.Raw, *gdb.Raw, gdb.Counter, *gdb.Counter:
		return value, nil
	}
	reflectValue := reflect.ValueOf(value)
	if reflectValue.Kind() == reflect.Pointer && reflectValue.IsNil() {
		return nil, nil
	}
	reflectValue = reflect.Indirect(reflectValue)
	var (
		family = fieldTypeFamily(fieldType)
		kind   = reflectValue.Kind()
		scalar = kind == reflect.String || kind == reflect.Bool ||
			(kind >= reflect.Int && kind <= reflect.Float64)
	)
	switch family {
	case "string", "fixedstring":
		switch kind {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct:
			converted, err := d.Core.ConvertValueForField(ctx, fieldType, value)
			if b, ok := converted.([]byte); ok {
				return string(b), err
			}
			return converted, err
		}
		if scalar {
			return gconv.String(value), nil
		}
		return value, nil
	case "int128", "uint128", "int256", "uint256":
		return convertBigInteger(family, fieldType, value, reflectValue)
	case "decimal", "decimal32", "decimal64", "decimal128", "decimal256":
		if scalar {
			return convertDecimal(fieldType, value)
		}
	}
	if !scalar {
		return value, nil
	}
	if kind == reflect.String {
		s := strings.TrimSpace(reflectValue.String())
		isNumber, err := checkNumberString(family, s)
		if err != nil {
			return nil, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `invalid value for column type "%s"`, fieldType)
		}
		if isNumber {
			value = s
		}
	} else if err := checkNumberRange(family, fieldType, reflectValue); err != nil {
		return nil, err
	}
	switch family {
	case "uint8":
		return gconv.Uint8(value), nil
	case "uint16":
		return gconv.Uint16(value), nil
	case "uint32":
		return gconv.Uint32(value), nil
	case "uint64":
		return gconv.Uint64(value), nil
	case "int8":
		return gconv.Int8(value), nil
	case "int16":
		return gconv.Int16(value), nil
	case "int32":
		return gconv.Int32(value), nil
	case "int64":
		return gconv.Int64(value), nil
	case "float32":
		return gconv.Float32(value), nil
	case "float64":
		return gconv.Float64(value), nil
	case "bool":
		if s, ok := value.(string); ok {
			return strconv.ParseBool(s)
		}
		return gconv.Bool(value), nil
	case "date", "date32", "datetime", "datetime64":
		if s, ok := value.(string); ok {
			t, err := gtime.StrToTime(s)
			if err != nil {
				return nil, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `invalid value for column type "%s"`, fieldType)
			}
			return t.Time, checkTimeRange(fieldType, t.Time)
		}
	}
	return value, nil
}

// checkNumberString reports whether `family` is an integer, float or boolean column type family,
// and returns an error if `s` is not a value of the family or is out of its range.
func checkNumberString(family, s string) (isNumber bool, err error) {
	isNumber = true
	switch family {
	case "uint8", "uint16", "uint32", "uint64":
		bitSize, _ := strconv.Atoi(strings.TrimPrefix(family, "uint"))
		_, err = strconv.ParseUint(s, 10, bitSize)
	case "int8", "int16", "int32", "int64":
		bitSize, _ := strconv.Atoi(strings.TrimPrefix(family, "int"))
		_, err = strconv.ParseInt(s, 10, bitSize)
	case "float32", "float64":
		bitSize, _ := strconv.Atoi(strings.TrimPrefix(family, "float"))
		_, err = strconv.ParseFloat(s, bitSize)
	case "bool":
		_, err = strconv.ParseBool(s)
	default:
		isNumber = false
	}
	return isNumber, err
}

// checkNumberRange returns an error of code gcode.CodeInvalidParameter if the integer, unsigned
// integer or float `v` is out of the range of the integer or Float32 column type family `family`,
// or has a fraction for an integer one, of the column type `fieldType`.
func checkNumberRange(family, fieldType string, v reflect.Value) error {
	if family == "float32" {
		if v.CanFloat() && !math.IsInf(v.Float(), 0) && math.Abs(v.Float()) > math.MaxFloat32 {
			return errOutOfRange(v.Interface(), fieldType)
		}
		return nil
	}
	bits, signed, ok := integerFamily(family)
	if !ok {
		return nil
	}
	var (
		maxValue = uint64(math.MaxUint64) >> (64 - bits)
		minValue int64
		inRange  bool
	)
	if signed {
		maxValue >>= 1
		minValue = -int64(maxValue) - 1
	}
	switch {
	case v.CanInt():
		inRange = v.Int() >= minValue && (v.Int() < 0 || uint64(v.Int()) <= maxValue)
	case v.CanUint():
		inRange = v.Uint() <= maxValue
	case v.CanFloat():
		if f := v.Float(); f != math.Trunc(f) {
			return errNotInteger(v.Interface(), fieldType)
		}
		inRange = v.Float() >= float64(minValue) && v.Float() < float64(maxValue)+1
	default:
		return nil
	}
	if !inRange {
		return errOutOfRange(v.Interface(), fieldType)
	}
	return nil
}

// convertBigInteger converts `value`, whose reflection value is `v`, to the *big.Int that the
// underlying driver binds to a column of the type `fieldType` of the Int128, UInt128, Int256 or
// UInt256 type family `family`, if it is an integer, an integral float, a decimal string or a
// big.Int. It returns an error of code gcode.CodeInvalidParameter if it is not an integer or is
// out of the range of the column, and any other value as it is.
func convertBigInteger(family, fieldType string, value any, v reflect.Value) (any, error) {
	var n *big.Int
	switch number := value.(type) {
	case *big.Int:
		n = number
	case big.Int:
		n = &number
	default:
		switch {
		case v.CanInt():
			n = big.NewInt(v.Int())
		case v.CanUint():
			n = new(big.Int).SetUint64(v.Uint())
		case v.CanFloat():
			if f := v.Float(); f != math.Trunc(f) || math.IsInf(f, 0) {
				return nil, errNotInteger(value, fieldType)
			}
			n, _ = big.NewFloat(v.Float()).Int(nil)
		case v.Kind() == reflect.String:
			var ok bool
			if n, ok = new(big.Int).SetString(strings.TrimSpace(v.String()), 10); !ok {
				return nil, errNotInteger(value, fieldType)
			}
		default:
			return value, nil
		}
	}
	bits, signed, _ := integerFamily(family)
	if !fitsInteger(n, bits, signed) {
		return nil, errOutOfRange(value, fieldType)
	}
	return n, nil
}

// fitsInteger reports whether `n` is in the range of a signed or unsigned integer of `bits` bits.
// A negative `n` fits a signed integer if its one's complement -n-1 fits the non-negative half.
func fitsInteger(n *big.Int, bits int, signed bool) bool {
	switch {
	case !signed:
		return n.Sign() >= 0 && n.BitLen() <= bits
	case n.Sign() >= 0:
		return n.BitLen() < bits
	default:
		return new(big.Int).Not(n).BitLen() < bits
	}
}

// integerFamily returns the size in bits and the signedness of the integer column type family
// `family`, and false if it is not an integer one.
func integerFamily(family string) (bits int, signed bool, ok bool) {
	switch family {
	case "int8", "uint8":
		bits = 8
	case "int16", "uint16":
		bits = 16
	case "int32", "uint32":
		bits = 32
	case "int64", "uint64":
		bits = 64
	case "int128", "uint128":
		bits = 128
	case "int256", "uint256":
		bits = 256
	default:
		return 0, false, false
	}
	return bits, !strings.HasPrefix(family, "u"), true
}

// convertDecimal converts `value` to a decimal.Decimal for a column of the Decimal type `fieldType`.
// If `fieldType` has its precision and scale, the decimal is rounded half away from zero to the
// scale, as MySQL and PostgreSQL do where ClickHouse truncates, and an error of code
// gcode.CodeInvalidParameter is returned if its rounded integer part exceeds the precision, as it
// is if `value` is not a number. The core passes ConvertValueForField the type of a column without
// its arguments, so inserts and updates check their values in DoInsert and DoUpdate, through
// convertDecimalValues.
func convertDecimal(fieldType string, value any) (any, error) {
	number, ok := value.(decimal.Decimal)
	if !ok {
		var err error
		if number, err = decimal.NewFromString(strings.TrimSpace(gconv.String(value))); err != nil {
			return nil, gerror.WrapCodef(gcode.CodeInvalidParameter, err, `invalid value for column type "%s"`, fieldType)
		}
	}
	precision, scale, ok := decimalPrecisionScale(fieldType)
	if !ok {
		return number, nil
	}
	number = number.Round(scale)
	if number.Abs().Cmp(decimal.New(1, precision-scale)) >= 0 {
		return nil, errOutOfRange(value, fieldType)
	}
	return number, nil
}

// convertDecimalValues returns `records` with the values of the Decimal columns of `fields`
// converted by ConvertValueForField with the types of the columns, which round them to the scale
// of their column and check them against its precision, as the core passes
// ConvertValueForField the type of a column without its precision and scale. The records are
// copied if `fields` has a Decimal column, and returned as they are otherwise.
func (d *Driver) convertDecimalValues(
	ctx context.Context, fields map[string]*gdb.TableField, records gdb.List,
) (gdb.List, error) {
	decimalTypes := make(map[string]string)
	for name, field := range fields {
		if _, _, ok := decimalPrecisionScale(field.Type); ok {
			decimalTypes[name] = field.Type
		}
	}
	if len(decimalTypes) == 0 {
		return records, nil
	}
	var (
		converted = make(gdb.List, len(records))
		err       error
	)
	for i, record := range records {
		converted[i] = maps.Clone(record)
		for name, fieldType := range decimalTypes {
			value, ok := record[name]
			if !ok {
				continue
			}
			if converted[i][name], err = d.ConvertValueForField(ctx, fieldType, value); err != nil {
				return nil, err
			}
		}
	}
	return converted, nil
}

// tableFields returns the fields of `table`, which DoInsert and DoUpdate write, from the cached
// TableFields of the database that `table` is prefixed with, or else of the configured database,
// as the core reads them to convert the values of every driver.
func (d *Driver) tableFields(ctx context.Context, table string) (map[string]*gdb.TableField, error) {
	database, name := splitTableName(d.QuotePrefixTableName(table))
	if database == "" {
		database = d.GetSchema()
	}
	return d.GetDB().TableFields(ctx, name, database)
}

var (
	minDateTime   = time.Unix(0, 0)
	maxDateTime   = time.Unix(math.MaxUint32, 0)
	minDateTime64 = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDateTime64 = time.Unix(0, math.MaxInt64)
	minDate       = time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate       = time.Date(2149, 6, 6, 0, 0, 0, 0, time.UTC)
	minDate32     = time.Date(1900, 1, 1, 0, 0, 0, 0, time.UTC)
	maxDate32     = time.Date(2299, 12, 31, 0, 0, 0, 0, time.UTC)
)

// checkTimeRange returns an error of code gcode.CodeInvalidParameter if the time `t` is out of
// the range of a column of the DateTime, DateTime64, Date or Date32 type `fieldType`, where
// ClickHouse or the underlying driver would wrap it around or clamp it. A Date or Date32 is
// compared by the date of `t` in its own time zone, as the underlying driver writes it. A
// DateTime64 begins on 1900-01-01 in the time zone of `t`, as ClickHouse counts it in the time
// zone of the column, and ends at the last time the underlying driver binds, before its
// documented end in 2299.
func checkTimeRange(fieldType string, t time.Time) error {
	var (
		first, last time.Time
		value       = t
	)
	switch fieldTypeFamily(fieldType) {
	case "datetime":
		first, last = minDateTime, maxDateTime
	case "datetime64":
		if dateOf(t).Before(minDateTime64) || t.After(maxDateTime64) {
			return errOutOfRange(t, fieldType)
		}
		return nil
	case "date":
		first, last, value = minDate, maxDate, dateOf(t)
	case "date32":
		first, last, value = minDate32, maxDate32, dateOf(t)
	default:
		return nil
	}
	if value.Before(first) || value.After(last) {
		return errOutOfRange(t, fieldType)
	}
	return nil
}

// dateOf returns the date of `t` in its own time zone as the midnight of that date in UTC.
func dateOf(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, time.UTC)
}

func errOutOfRange(value any, fieldType string) error {
	return gerror.NewCodef(gcode.CodeInvalidParameter, `value %v is out of the range of column type "%s"`, value, fieldType)
}

func errNotInteger(value any, fieldType string) error {
	return gerror.NewCodef(gcode.CodeInvalidParameter, `value %v is not an integer for column type "%s"`, value, fieldType)
}

// fieldTypeFamilies are the column type families, in lower case, whose values
// convertValueByFieldType converts.
var fieldTypeFamilies = [...]string{
	"string", "fixedstring", "bool",
	"uint8", "uint16", "uint32", "uint64", "uint128", "uint256",
	"int8", "int16", "int32", "int64", "int128", "int256",
	"float32", "float64",
	"decimal", "decimal32", "decimal64", "decimal128", "decimal256",
	"date", "date32", "datetime", "datetime64",
}

// fieldTypeFamily returns the type family of `fieldType` in lower case, without its arguments and
// its Nullable or LowCardinality wrapper, like `decimal` for `Nullable(Decimal(5, 2))`, or an empty
// string if it is not one of fieldTypeFamilies. It does not allocate, as it runs for every value
// written.
func fieldTypeFamily(fieldType string) string {
	family, _, _ := strings.Cut(unwrapNullable(fieldType), "(")
	family = strings.TrimSpace(family)
	for _, name := range fieldTypeFamilies {
		if strings.EqualFold(family, name) {
			return name
		}
	}
	return ""
}

// decimalPrecisionScale returns the precision and the scale of the Decimal type `fieldType`, like
// 5 and 2 for `Decimal(5, 2)` or 18 and 4 for `Decimal64(4)`, and false if `fieldType` is not a
// Decimal type with its arguments, like `Decimal` that the core passes to ConvertValueForField.
func decimalPrecisionScale(fieldType string) (precision, scale int32, ok bool) {
	name, args, hasArgs := strings.Cut(unwrapNullable(fieldType), "(")
	name, args = strings.TrimSpace(name), strings.TrimSpace(strings.TrimSuffix(args, ")"))
	switch {
	case !hasArgs:
		return 0, 0, false
	case strings.EqualFold(name, "decimal"):
		precisionArg, scaleArg, _ := strings.Cut(args, ",")
		return gconv.Int32(strings.TrimSpace(precisionArg)), gconv.Int32(strings.TrimSpace(scaleArg)), true
	case strings.EqualFold(name, "decimal32"):
		precision = 9
	case strings.EqualFold(name, "decimal64"):
		precision = 18
	case strings.EqualFold(name, "decimal128"):
		precision = 38
	case strings.EqualFold(name, "decimal256"):
		precision = 76
	default:
		return 0, 0, false
	}
	return precision, gconv.Int32(args), true
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The underlying driver returns the values of a DateTime column declared without a time zone in the
// time zone of the server; they are read in the local time zone instead, as MySQL does with loc=Local.
// A decimal keeps the scale of its column, like "0.00" for Decimal(5, 2).
func (d *Driver) ConvertValueForLocal(ctx context.Context, fieldType string, fieldValue any) (any, error) {
	if isZonelessDateTime(fieldType) {
		switch t := fieldValue.(type) {
		case time.Time:
			fieldValue = t.In(time.Local)
		case *time.Time:
			if t != nil {
				fieldValue = t.In(time.Local)
			}
		}
	}
	if _, scale, ok := decimalPrecisionScale(fieldType); ok {
		switch v := fieldValue.(type) {
		case string:
			if value, err := decimal.NewFromString(v); err == nil {
				fieldValue = value.StringFixed(scale)
			}
		case decimal.Decimal:
			fieldValue = v.StringFixed(scale)
		case *decimal.Decimal:
			if v != nil {
				fieldValue = v.StringFixed(scale)
			}
		}
	}
	return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
}

// isZonelessDateTime reports whether `fieldType` is a DateTime or DateTime64 type declared without
// a time zone, like `DateTime`, `DateTime64(3)` or `Nullable(DateTime)`.
func isZonelessDateTime(fieldType string) bool {
	fieldType = unwrapNullable(fieldType)
	switch {
	case fieldType == "DateTime":
		return true
	case strings.HasPrefix(fieldType, "DateTime64(") && strings.HasSuffix(fieldType, ")"):
		return !strings.Contains(fieldType, ",")
	default:
		return false
	}
}
