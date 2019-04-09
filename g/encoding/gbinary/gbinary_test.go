// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbinary_test

import (
	"github.com/gogf/gf/g/encoding/gbinary"
	"github.com/gogf/gf/g/test/gtest"
	"math"
	"testing"
)

var testData = map[string]interface{}{
	//"nil":         nil,
	"int":         int(123),
	"int8":        int8(-99),
	"int8.max":    math.MaxInt8,
	"int16":       int16(123),
	"int16.max":   math.MaxInt16,
	"int32":       int32(-199),
	"int32.max":   math.MaxInt32,
	"int64":       int64(123),
	"uint":        uint(123),
	"uint8":       uint8(123),
	"uint8.max":   math.MaxUint8,
	"uint16":      uint16(9999),
	"uint16.max":  math.MaxUint16,
	"uint32":      uint32(123),
	"uint64":      uint64(123),
	"bool.true":   true,
	"bool.false":  false,
	"string":      "hehe haha",
	"byte":        []byte("hehe haha"),
	"float32":     float32(123.456),
	"float32.max": math.MaxFloat32,
	"float64":     float64(123.456),
}

func TestEncodeAndDecode(t *testing.T) {
	for k, v := range testData {
		ve := gbinary.Encode(v)
		ve1 := gbinary.EncodeByLength(len(ve), v)

		//t.Logf("%s:%v, encoded:%v\n", k, v, ve)
		switch v.(type) {
		case int:
			gtest.Assert(gbinary.DecodeToInt(ve), v)
			gtest.Assert(gbinary.DecodeToInt(ve1), v)
		case int8:
			gtest.Assert(gbinary.DecodeToInt8(ve), v)
			gtest.Assert(gbinary.DecodeToInt8(ve1), v)
		case int16:
			gtest.Assert(gbinary.DecodeToInt16(ve), v)
			gtest.Assert(gbinary.DecodeToInt16(ve1), v)
		case int32:
			gtest.Assert(gbinary.DecodeToInt32(ve), v)
			gtest.Assert(gbinary.DecodeToInt32(ve1), v)
		case int64:
			gtest.Assert(gbinary.DecodeToInt64(ve), v)
			gtest.Assert(gbinary.DecodeToInt64(ve1), v)
		case uint:
			gtest.Assert(gbinary.DecodeToUint(ve), v)
			gtest.Assert(gbinary.DecodeToUint(ve1), v)
		case uint8:
			gtest.Assert(gbinary.DecodeToUint8(ve), v)
			gtest.Assert(gbinary.DecodeToUint8(ve1), v)
		case uint16:
			gtest.Assert(gbinary.DecodeToUint16(ve1), v)
			gtest.Assert(gbinary.DecodeToUint16(ve), v)
		case uint32:
			gtest.Assert(gbinary.DecodeToUint32(ve1), v)
			gtest.Assert(gbinary.DecodeToUint32(ve), v)
		case uint64:
			gtest.Assert(gbinary.DecodeToUint64(ve), v)
			gtest.Assert(gbinary.DecodeToUint64(ve1), v)
		case bool:
			gtest.Assert(gbinary.DecodeToBool(ve), v)
			gtest.Assert(gbinary.DecodeToBool(ve1), v)
		case string:
			gtest.Assert(gbinary.DecodeToString(ve), v)
			gtest.Assert(gbinary.DecodeToString(ve1), v)
		case float32:
			gtest.Assert(gbinary.DecodeToFloat32(ve), v)
			gtest.Assert(gbinary.DecodeToFloat32(ve1), v)
		case float64:
			gtest.Assert(gbinary.DecodeToFloat64(ve), v)
			gtest.Assert(gbinary.DecodeToFloat64(ve1), v)
		default:
			if v == nil {
				continue
			}
			res := make([]byte, len(ve))
			err := gbinary.Decode(ve, res)
			if err != nil {
				t.Errorf("test data: %s, %v, error:%v", k, v, err)
			}
			gtest.Assert(res, v)
		}
	}
}

type User struct {
	Name string
	Age  int
	Url  string
}

func TestEncodeStruct(t *testing.T) {
	user := User{"wenzi1", 999, "www.baidu.com"}
	ve := gbinary.Encode(user)
	s := gbinary.DecodeToString(ve)
	gtest.Assert(string(s), s)
}

var testBitData = []int{0, 99, 122, 129, 222, 999, 22322}

func TestBits(t *testing.T) {
	for i := range testBitData {
		bits := make([]gbinary.Bit, 0)
		res := gbinary.EncodeBits(bits, testBitData[i], 64)

		gtest.Assert(gbinary.DecodeBits(res), testBitData[i])
		gtest.Assert(gbinary.DecodeBitsToUint(res), uint(testBitData[i]))

		gtest.Assert(gbinary.DecodeBytesToBits(gbinary.EncodeBitsToBytes(res)), res)
	}

}
