// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

import (
	"math"
	"reflect"

	"github.com/gogf/gf/v2/encoding/gbinary"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/util/gconv/internal/localinterface"
)

// Bytes converts `any` to []byte.
func (c *Converter) Bytes(any any) ([]byte, error) {
	if empty.IsNil(any) {
		return nil, nil
	}
	switch value := any.(type) {
	case string:
		return []byte(value), nil

	case []byte:
		return value, nil

	default:
		if f, ok := value.(localinterface.IBytes); ok {
			return f.Bytes(), nil
		}
		originValueAndKind := reflection.OriginValueAndKind(any)
		switch originValueAndKind.OriginKind {
		case reflect.Map:
			bytes, err := json.Marshal(any)
			if err != nil {
				return nil, err
			}
			return bytes, nil

		case reflect.Array, reflect.Slice:
			var (
				ok    = true
				bytes = make([]byte, originValueAndKind.OriginValue.Len())
			)
			for i := range bytes {
				int32Value, err := c.Int32(originValueAndKind.OriginValue.Index(i).Interface())
				if err != nil {
					return nil, err
				}
				if int32Value < 0 || int32Value > math.MaxUint8 {
					ok = false
					break
				}
				bytes[i] = byte(int32Value)
			}
			if ok {
				return bytes, nil
			}
		default:
		}
		return gbinary.Encode(any), nil
	}
}
