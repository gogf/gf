// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbinary_test

import (
	"testing"

	"github.com/gogf/gf/g/encoding/gbinary"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_BeEncodeAndBeDecode(t *testing.T) {
	for k, v := range testData {
		ve := gbinary.BeEncode(v)
		ve1 := gbinary.BeEncodeByLength(len(ve), v)

		//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
		switch v.(type) {
		case int:
			gtest.Assert(gbinary.BeDecodeToInt(ve), v)
			gtest.Assert(gbinary.BeDecodeToInt(ve1), v)
		case int8:
			gtest.Assert(gbinary.BeDecodeToInt8(ve), v)
			gtest.Assert(gbinary.BeDecodeToInt8(ve1), v)
		case int16:
			gtest.Assert(gbinary.BeDecodeToInt16(ve), v)
			gtest.Assert(gbinary.BeDecodeToInt16(ve1), v)
		case int32:
			gtest.Assert(gbinary.BeDecodeToInt32(ve), v)
			gtest.Assert(gbinary.BeDecodeToInt32(ve1), v)
		case int64:
			gtest.Assert(gbinary.BeDecodeToInt64(ve), v)
			gtest.Assert(gbinary.BeDecodeToInt64(ve1), v)
		case uint:
			gtest.Assert(gbinary.BeDecodeToUint(ve), v)
			gtest.Assert(gbinary.BeDecodeToUint(ve1), v)
		case uint8:
			gtest.Assert(gbinary.BeDecodeToUint8(ve), v)
			gtest.Assert(gbinary.BeDecodeToUint8(ve1), v)
		case uint16:
			gtest.Assert(gbinary.BeDecodeToUint16(ve1), v)
			gtest.Assert(gbinary.BeDecodeToUint16(ve), v)
		case uint32:
			gtest.Assert(gbinary.BeDecodeToUint32(ve1), v)
			gtest.Assert(gbinary.BeDecodeToUint32(ve), v)
		case uint64:
			gtest.Assert(gbinary.BeDecodeToUint64(ve), v)
			gtest.Assert(gbinary.BeDecodeToUint64(ve1), v)
		case bool:
			gtest.Assert(gbinary.BeDecodeToBool(ve), v)
			gtest.Assert(gbinary.BeDecodeToBool(ve1), v)
		case string:
			gtest.Assert(gbinary.BeDecodeToString(ve), v)
			gtest.Assert(gbinary.BeDecodeToString(ve1), v)
		case float32:
			gtest.Assert(gbinary.BeDecodeToFloat32(ve), v)
			gtest.Assert(gbinary.BeDecodeToFloat32(ve1), v)
		case float64:
			gtest.Assert(gbinary.BeDecodeToFloat64(ve), v)
			gtest.Assert(gbinary.BeDecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gbinary.BeDecode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			gtest.Assert(res, v)
		}
	}
}

func Test_BeEncodeStruct(t *testing.T) {
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gbinary.BeEncode(user)
	s := gbinary.BeDecodeToString(ve)
	gtest.Assert(string(s), s)
}
