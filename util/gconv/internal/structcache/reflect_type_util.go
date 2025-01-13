// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structcache

import (
	"reflect"
	"time"

	"github.com/gogf/gf/v2/os/gtime"
)

var (
	intType   = reflect.TypeOf(0)
	int8Type  = reflect.TypeOf(int8(0))
	int16Type = reflect.TypeOf(int16(0))
	int32Type = reflect.TypeOf(int32(0))
	int64Type = reflect.TypeOf(int64(0))

	uintType   = reflect.TypeOf(uint(0))
	uint8Type  = reflect.TypeOf(uint8(0))
	uint16Type = reflect.TypeOf(uint16(0))
	uint32Type = reflect.TypeOf(uint32(0))
	uint64Type = reflect.TypeOf(uint64(0))

	float32Type = reflect.TypeOf(float32(0))
	float64Type = reflect.TypeOf(float64(0))

	stringType = reflect.TypeOf("")
	bytesType  = reflect.TypeOf([]byte{})

	boolType = reflect.TypeOf(false)

	timeType  = reflect.TypeOf((*time.Time)(nil)).Elem()
	gtimeType = reflect.TypeOf((*gtime.Time)(nil)).Elem()
)
