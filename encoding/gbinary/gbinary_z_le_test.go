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

func Test_LeEncodeAndLeDecode(t *testing.T) {
	for k, v := range testData {
		ve := gbinary.LeEncode(v)
		ve1 := gbinary.LeEncodeByLength(len(ve), v)

		//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
		switch v.(type) {
		case int:
			t.Assert(gbinary.LeDecodeToInt(ve), v)
			t.Assert(gbinary.LeDecodeToInt(ve1), v)
		case int8:
			t.Assert(gbinary.LeDecodeToInt8(ve), v)
			t.Assert(gbinary.LeDecodeToInt8(ve1), v)
		case int16:
			t.Assert(gbinary.LeDecodeToInt16(ve), v)
			t.Assert(gbinary.LeDecodeToInt16(ve1), v)
		case int32:
			t.Assert(gbinary.LeDecodeToInt32(ve), v)
			t.Assert(gbinary.LeDecodeToInt32(ve1), v)
		case int64:
			t.Assert(gbinary.LeDecodeToInt64(ve), v)
			t.Assert(gbinary.LeDecodeToInt64(ve1), v)
		case uint:
			t.Assert(gbinary.LeDecodeToUint(ve), v)
			t.Assert(gbinary.LeDecodeToUint(ve1), v)
		case uint8:
			t.Assert(gbinary.LeDecodeToUint8(ve), v)
			t.Assert(gbinary.LeDecodeToUint8(ve1), v)
		case uint16:
			t.Assert(gbinary.LeDecodeToUint16(ve1), v)
			t.Assert(gbinary.LeDecodeToUint16(ve), v)
		case uint32:
			t.Assert(gbinary.LeDecodeToUint32(ve1), v)
			t.Assert(gbinary.LeDecodeToUint32(ve), v)
		case uint64:
			t.Assert(gbinary.LeDecodeToUint64(ve), v)
			t.Assert(gbinary.LeDecodeToUint64(ve1), v)
		case bool:
			t.Assert(gbinary.LeDecodeToBool(ve), v)
			t.Assert(gbinary.LeDecodeToBool(ve1), v)
		case string:
			t.Assert(gbinary.LeDecodeToString(ve), v)
			t.Assert(gbinary.LeDecodeToString(ve1), v)
		case float32:
			t.Assert(gbinary.LeDecodeToFloat32(ve), v)
			t.Assert(gbinary.LeDecodeToFloat32(ve1), v)
		case float64:
			t.Assert(gbinary.LeDecodeToFloat64(ve), v)
			t.Assert(gbinary.LeDecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gbinary.LeDecode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			t.Assert(res, v)
		}
	}
}

func Test_LeEncodeStruct(t *testing.T) {
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gbinary.LeEncode(user)
	s := gbinary.LeDecodeToString(ve)
	t.Assert(string(s), s)
}
