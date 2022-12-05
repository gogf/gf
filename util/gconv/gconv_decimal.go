/**
// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
*/

package gconv

import (
	"github.com/shopspring/decimal"
)

// Decimal converts <any> to decimal.
func Decimal(any interface{}) decimal.Decimal {
	if any == nil {
		return decimal.Zero
	}
	switch value := any.(type) {
	case decimal.Decimal:
		return value
	case int32:
		return decimal.NewFromInt32(value)
	case int64:
		return decimal.NewFromInt(value)
	case float32:
		return decimal.NewFromFloat32(value)
	case float64:
		return decimal.NewFromFloat(value)
	case []byte:
		v, _ := decimal.NewFromString(String(value))
		return v
	default:
		v, _ := decimal.NewFromString(String(value))
		return v
	}
}
