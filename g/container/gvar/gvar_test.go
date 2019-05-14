// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"bytes"
	"encoding/binary"
	"testing"

	"github.com/gogf/gf/g/container/gvar"
	"github.com/gogf/gf/g/test/gtest"
)

func TestSet(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New("old", true)
		objOneOld, _ := objOne.Set("new").(string)
		gtest.Assert(objOneOld, "old")

		objTwo := gvar.New("old", false)
		objTwoOld, _ := objTwo.Set("new").(string)
		gtest.Assert(objTwoOld, "old")
	})
}

func TestVal(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Val().(int)
		gtest.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Val().(int)
		gtest.Assert(objTwoOld, 1)
	})
}
func TestInterface(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Interface().(int)
		gtest.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Interface().(int)
		gtest.Assert(objTwoOld, 1)
	})
}
func TestIsNil(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(nil, true)
		gtest.Assert(objOne.IsNil(), true)

		objTwo := gvar.New("noNil", false)
		gtest.Assert(objTwo.IsNil(), false)

	})
}

func TestBytes(t *testing.T) {
	gtest.Case(t, func() {
		x := int32(1)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, x)

		objOne := gvar.New(bytesBuffer.Bytes(), true)

		bBuf := bytes.NewBuffer(objOne.Bytes())
		var y int32
		binary.Read(bBuf, binary.BigEndian, &y)

		gtest.Assert(x, y)

	})
}

func TestString(t *testing.T) {
	gtest.Case(t, func() {
		var str string = "hello"
		objOne := gvar.New(str, true)
		gtest.Assert(objOne.String(), str)

	})
}
func TestBool(t *testing.T) {
	gtest.Case(t, func() {
		var ok bool = true
		objOne := gvar.New(ok, true)
		gtest.Assert(objOne.Bool(), ok)

		ok = false
		objTwo := gvar.New(ok, true)
		gtest.Assert(objTwo.Bool(), ok)

	})
}

func TestInt(t *testing.T) {
	gtest.Case(t, func() {
		var num int = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int(), num)

	})
}

func TestInt8(t *testing.T) {
	gtest.Case(t, func() {
		var num int8 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int8(), num)

	})
}

func TestInt16(t *testing.T) {
	gtest.Case(t, func() {
		var num int16 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int16(), num)

	})
}

func TestInt32(t *testing.T) {
	gtest.Case(t, func() {
		var num int32 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int32(), num)

	})
}

func TestInt64(t *testing.T) {
	gtest.Case(t, func() {
		var num int64 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int64(), num)

	})
}

func TestUint(t *testing.T) {
	gtest.Case(t, func() {
		var num uint = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint(), num)

	})
}

func TestUint8(t *testing.T) {
	gtest.Case(t, func() {
		var num uint8 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint8(), num)

	})
}

func TestUint16(t *testing.T) {
	gtest.Case(t, func() {
		var num uint16 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint16(), num)

	})
}

func TestUint32(t *testing.T) {
	gtest.Case(t, func() {
		var num uint32 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint32(), num)

	})
}

func TestUint64(t *testing.T) {
	gtest.Case(t, func() {
		var num uint64 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint64(), num)

	})
}
func TestFloat32(t *testing.T) {
	gtest.Case(t, func() {
		var num float32 = 1.1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Float32(), num)

	})
}

func TestFloat64(t *testing.T) {
	gtest.Case(t, func() {
		var num float64 = 1.1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Float64(), num)

	})
}

func TestInts(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Ints()[0], arr[0])
	})
}
func TestFloats(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Floats()[0], arr[0])
	})
}
func TestStrings(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []string{"hello", "world"}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Strings()[0], arr[0])
	})
}

func TestTime(t *testing.T) {
	gtest.Case(t, func() {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		gtest.Assert(objOne.Time().Unix(), timeUnix)
	})
}

type StTest struct {
	Test int
}

func TestStruct(t *testing.T) {
	gtest.Case(t, func() {
		Kv := make(map[string]int, 1)
		Kv["Test"] = 100

		testObj := &StTest{}

		objOne := gvar.New(Kv, true)

		objOne.Struct(testObj)

		gtest.Assert(testObj.Test, Kv["Test"])
	})
}
