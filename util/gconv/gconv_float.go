// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv

import (
	"strconv"

	"github.com/gogf/gf/v2/encoding/gbinary"
)

// Float32 converts `any` to float32.
func Float32(any interface{}) float32 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return value
	case float64:
		return float32(value)
	case []byte:
		return gbinary.DecodeToFloat32(value)
	default:
		if f, ok := value.(iFloat32); ok {
			return f.Float32()
		}
		v, _ := strconv.ParseFloat(String(any), 64)
		return float32(v)
	}
}

// Float64 converts `any` to float64.
func Float64(any interface{}) float64 {
	if any == nil {
		return 0
	}
	switch value := any.(type) {
	case float32:
		return float64(value)
	case float64:
		return value
	case []byte:
		return gbinary.DecodeToFloat64(value)
	default:
		if f, ok := value.(iFloat64); ok {
			return f.Float64()
		}
		v, _ := strconv.ParseFloat(String(any), 64)
		return v
	}
}
