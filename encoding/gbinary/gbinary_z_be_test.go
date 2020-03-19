// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbinary_test

import (
	"testing"

	"github.com/gogf/gf/encoding/gbinary"
	"github.com/gogf/gf/test/gtest"
)

func Test_BeEncodeAndBeDecode(t *testing.T) {
	for k, v := range testData {
		ve := gbinary.BeEncode(v)
		ve1 := gbinary.BeEncodeByLength(len(ve), v)

		//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
		switch v.(type) {
		case int:
			t.Assert(gbinary.BeDecodeToInt(ve), v)
			t.Assert(gbinary.BeDecodeToInt(ve1), v)
		case int8:
			t.Assert(gbinary.BeDecodeToInt8(ve), v)
			t.Assert(gbinary.BeDecodeToInt8(ve1), v)
		case int16:
			t.Assert(gbinary.BeDecodeToInt16(ve), v)
			t.Assert(gbinary.BeDecodeToInt16(ve1), v)
		case int32:
			t.Assert(gbinary.BeDecodeToInt32(ve), v)
			t.Assert(gbinary.BeDecodeToInt32(ve1), v)
		case int64:
			t.Assert(gbinary.BeDecodeToInt64(ve), v)
			t.Assert(gbinary.BeDecodeToInt64(ve1), v)
		case uint:
			t.Assert(gbinary.BeDecodeToUint(ve), v)
			t.Assert(gbinary.BeDecodeToUint(ve1), v)
		case uint8:
			t.Assert(gbinary.BeDecodeToUint8(ve), v)
			t.Assert(gbinary.BeDecodeToUint8(ve1), v)
		case uint16:
			t.Assert(gbinary.BeDecodeToUint16(ve1), v)
			t.Assert(gbinary.BeDecodeToUint16(ve), v)
		case uint32:
			t.Assert(gbinary.BeDecodeToUint32(ve1), v)
			t.Assert(gbinary.BeDecodeToUint32(ve), v)
		case uint64:
			t.Assert(gbinary.BeDecodeToUint64(ve), v)
			t.Assert(gbinary.BeDecodeToUint64(ve1), v)
		case bool:
			t.Assert(gbinary.BeDecodeToBool(ve), v)
			t.Assert(gbinary.BeDecodeToBool(ve1), v)
		case string:
			t.Assert(gbinary.BeDecodeToString(ve), v)
			t.Assert(gbinary.BeDecodeToString(ve1), v)
		case float32:
			t.Assert(gbinary.BeDecodeToFloat32(ve), v)
			t.Assert(gbinary.BeDecodeToFloat32(ve1), v)
		case float64:
			t.Assert(gbinary.BeDecodeToFloat64(ve), v)
			t.Assert(gbinary.BeDecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gbinary.BeDecode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			t.Assert(res, v)
		}
	}
}

func Test_BeEncodeStruct(t *testing.T) {
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gbinary.BeEncode(user)
	s := gbinary.BeDecodeToString(ve)
	t.Assert(string(s), s)
}
