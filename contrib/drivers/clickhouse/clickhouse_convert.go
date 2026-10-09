// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql/driver"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gogf/gf/v2/os/gtime"
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
		// otherwise the interface data is committed
		valuer, ok := itemValue.(driver.Valuer)
		if !ok {
			return itemValue, nil
		}
		convertedValue, err := valuer.Value()
		if err != nil {
			return nil, err
		}
		return convertedValue, nil
	}
}

// ConvertValueForLocal converts value to local Golang type of value according field type name from database.
// The underlying driver returns the values of a DateTime column declared without a time zone in the
// time zone of the server; they are read in the local time zone instead, as MySQL does with loc=Local.
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
	return d.Core.ConvertValueForLocal(ctx, fieldType, fieldValue)
}

// isZonelessDateTime reports whether `fieldType` is a DateTime or DateTime64 type declared without
// a time zone, like `DateTime`, `DateTime64(3)` or `Nullable(DateTime)`.
func isZonelessDateTime(fieldType string) bool {
	for _, wrapper := range []string{"Nullable(", "LowCardinality("} {
		if strings.HasPrefix(fieldType, wrapper) && strings.HasSuffix(fieldType, ")") {
			fieldType = fieldType[len(wrapper) : len(fieldType)-1]
		}
	}
	switch {
	case fieldType == "DateTime":
		return true
	case strings.HasPrefix(fieldType, "DateTime64(") && strings.HasSuffix(fieldType, ")"):
		return !strings.Contains(fieldType, ",")
	default:
		return false
	}
}
