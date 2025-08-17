// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"context"
	"database/sql/driver"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"

	"github.com/gogf/gf/v2/os/gtime"
)

// ConvertValueForField converts value to the type of the record field.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	switch itemValue := fieldValue.(type) {
	case time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}

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
		// for gtime type, needs to get time.Time
		if itemValue != nil {
			return itemValue.Time, nil
		}
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}

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
	return fieldValue, nil
}
