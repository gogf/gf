// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql/driver"
	"reflect"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gogf/gf/v2/database/gdb"
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
		return itemValue, nil

	case uuid.UUID:
		return itemValue, nil

	case *time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		return itemValue, nil

	case gtime.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}
		// for gtime type, needs to get time.Time
		return itemValue.Time, nil

	case *gtime.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		// for gtime type, needs to get time.Time
		return itemValue.Time, nil

	case decimal.Decimal:
		return itemValue, nil

	case *decimal.Decimal:
		if itemValue != nil {
			return *itemValue, nil
		}
		return nil, nil

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
// String column as the core converts it. Any other value, like a slice for an Array column, is
// returned as it is.
func (d *Driver) convertValueByFieldType(ctx context.Context, fieldType string, value any) (any, error) {
	switch value.(type) {
	case nil, []byte, gdb.Raw, *gdb.Raw, gdb.Counter, *gdb.Counter:
		return value, nil
	}
	var (
		family = fieldTypeFamily(fieldType)
		kind   = reflect.Indirect(reflect.ValueOf(value)).Kind()
		scalar = kind == reflect.String || kind == reflect.Bool ||
			(kind >= reflect.Int && kind <= reflect.Float64)
	)
	switch family {
	case "string", "fixedstring":
		switch kind {
		case reflect.Map, reflect.Slice, reflect.Array, reflect.Struct, reflect.Invalid:
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
	}
	if !scalar {
		return value, nil
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
		return gconv.Bool(value), nil
	case "decimal", "decimal32", "decimal64", "decimal128", "decimal256":
		return decimal.NewFromString(gconv.String(value))
	case "date", "date32", "datetime", "datetime64":
		if s, ok := value.(string); ok {
			t, err := gtime.StrToTime(s)
			if err != nil {
				return nil, err
			}
			return t.Time, nil
		}
	}
	return value, nil
}

// fieldTypeFamily returns the type family of `fieldType` in lower case, without its
// arguments and its Nullable or LowCardinality wrapper, like `decimal` for `Nullable(Decimal(5, 2))`.
func fieldTypeFamily(fieldType string) string {
	family, _, _ := strings.Cut(unwrapNullable(fieldType), "(")
	return strings.ToLower(strings.TrimSpace(family))
}

// decimalScale returns the scale of the Decimal type `fieldType`, like 2 for `Decimal(5, 2)` or
// 4 for `Decimal64(4)`, and false if `fieldType` is not a Decimal type.
func decimalScale(fieldType string) (int32, bool) {
	var (
		family     = fieldTypeFamily(fieldType)
		_, args, _ = strings.Cut(unwrapNullable(fieldType), "(")
		params     = strings.Split(strings.TrimSuffix(args, ")"), ",")
	)
	switch family {
	case "decimal":
		if len(params) == 2 {
			return gconv.Int32(strings.TrimSpace(params[1])), true
		}
		return 0, true
	case "decimal32", "decimal64", "decimal128", "decimal256":
		return gconv.Int32(strings.TrimSpace(params[0])), true
	default:
		return 0, false
	}
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
	if scale, ok := decimalScale(fieldType); ok {
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
