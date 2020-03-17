// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"github.com/gogf/gf/util/gconv"
	"math"
	"testing"
	"time"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/test/gtest"
)

func Test_Set(t *testing.T) {
	gtest.Case(t, func() {
		var v gvar.Var
		v.Set(123.456)
		gtest.Assert(v.Val(), 123.456)
	})
	gtest.Case(t, func() {
		var v gvar.Var
		v.Set(123.456)
		gtest.Assert(v.Val(), 123.456)
	})
	gtest.Case(t, func() {
		v := gvar.Create(123.456)
		gtest.Assert(v.Val(), 123.456)
	})
	gtest.Case(t, func() {
		objOne := gvar.New("old", true)
		objOneOld, _ := objOne.Set("new").(string)
		gtest.Assert(objOneOld, "old")

		objTwo := gvar.New("old", false)
		objTwoOld, _ := objTwo.Set("new").(string)
		gtest.Assert(objTwoOld, "old")
	})
}

func Test_Val(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Val().(int)
		gtest.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Val().(int)
		gtest.Assert(objTwoOld, 1)
	})
}
func Test_Interface(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Interface().(int)
		gtest.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Interface().(int)
		gtest.Assert(objTwoOld, 1)
	})
}
func Test_IsNil(t *testing.T) {
	gtest.Case(t, func() {
		objOne := gvar.New(nil, true)
		gtest.Assert(objOne.IsNil(), true)

		objTwo := gvar.New("noNil", false)
		gtest.Assert(objTwo.IsNil(), false)

	})
}

func Test_Bytes(t *testing.T) {
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

func Test_String(t *testing.T) {
	gtest.Case(t, func() {
		var str string = "hello"
		objOne := gvar.New(str, true)
		gtest.Assert(objOne.String(), str)

	})
}
func Test_Bool(t *testing.T) {
	gtest.Case(t, func() {
		var ok bool = true
		objOne := gvar.New(ok, true)
		gtest.Assert(objOne.Bool(), ok)

		ok = false
		objTwo := gvar.New(ok, true)
		gtest.Assert(objTwo.Bool(), ok)

	})
}

func Test_Int(t *testing.T) {
	gtest.Case(t, func() {
		var num int = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int(), num)

	})
}

func Test_Int8(t *testing.T) {
	gtest.Case(t, func() {
		var num int8 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int8(), num)

	})
}

func Test_Int16(t *testing.T) {
	gtest.Case(t, func() {
		var num int16 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int16(), num)

	})
}

func Test_Int32(t *testing.T) {
	gtest.Case(t, func() {
		var num int32 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int32(), num)

	})
}

func Test_Int64(t *testing.T) {
	gtest.Case(t, func() {
		var num int64 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Int64(), num)

	})
}

func Test_Uint(t *testing.T) {
	gtest.Case(t, func() {
		var num uint = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint(), num)

	})
}

func Test_Uint8(t *testing.T) {
	gtest.Case(t, func() {
		var num uint8 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint8(), num)

	})
}

func Test_Uint16(t *testing.T) {
	gtest.Case(t, func() {
		var num uint16 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint16(), num)

	})
}

func Test_Uint32(t *testing.T) {
	gtest.Case(t, func() {
		var num uint32 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint32(), num)

	})
}

func Test_Uint64(t *testing.T) {
	gtest.Case(t, func() {
		var num uint64 = 1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Uint64(), num)

	})
}
func Test_Float32(t *testing.T) {
	gtest.Case(t, func() {
		var num float32 = 1.1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Float32(), num)

	})
}

func Test_Float64(t *testing.T) {
	gtest.Case(t, func() {
		var num float64 = 1.1
		objOne := gvar.New(num, true)
		gtest.Assert(objOne.Float64(), num)

	})
}

func Test_Ints(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Ints()[0], arr[0])
	})
}
func Test_Floats(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Floats()[0], arr[0])
	})
}
func Test_Strings(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []string{"hello", "world"}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Strings()[0], arr[0])
	})
}

func Test_Interfaces(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Interfaces(), arr)
	})
}

func Test_Slice(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		gtest.Assert(objOne.Slice(), arr)
	})
}

func Test_Array(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		gtest.Assert(objOne.Array(), arr)
	})
}

func Test_Vars(t *testing.T) {
	gtest.Case(t, func() {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		gtest.Assert(len(objOne.Vars()), 5)
		gtest.Assert(objOne.Vars()[0].Int(), 1)
		gtest.Assert(objOne.Vars()[4].Int(), 5)
	})
}

func Test_Time(t *testing.T) {
	gtest.Case(t, func() {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		gtest.Assert(objOne.Time().Unix(), timeUnix)
	})
}

func Test_GTime(t *testing.T) {
	gtest.Case(t, func() {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		gtest.Assert(objOne.GTime().Unix(), timeUnix)
	})
}

func Test_Duration(t *testing.T) {
	gtest.Case(t, func() {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		gtest.Assert(objOne.Duration(), time.Duration(timeUnix))
	})
}

func Test_Map(t *testing.T) {
	gtest.Case(t, func() {
		m := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		objOne := gvar.New(m, true)
		gtest.Assert(objOne.Map()["k1"], m["k1"])
		gtest.Assert(objOne.Map()["k2"], m["k2"])
	})
}

func Test_Struct(t *testing.T) {
	gtest.Case(t, func() {
		type StTest struct {
			Test int
		}

		Kv := make(map[string]int, 1)
		Kv["Test"] = 100

		testObj := &StTest{}

		objOne := gvar.New(Kv, true)

		objOne.Struct(testObj)

		gtest.Assert(testObj.Test, Kv["Test"])
	})
	gtest.Case(t, func() {
		type StTest struct {
			Test int8
		}
		o := &StTest{}
		v := gvar.New(g.Slice{"Test", "-25"})
		v.Struct(o)
		gtest.Assert(o.Test, -25)
	})
}

func Test_Json(t *testing.T) {
	// Marshal
	gtest.Case(t, func() {
		s := "i love gf"
		v := gvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})

	gtest.Case(t, func() {
		s := int64(math.MaxInt64)
		v := gvar.New(s)
		b1, err1 := json.Marshal(v)
		b2, err2 := json.Marshal(s)
		gtest.Assert(err1, err2)
		gtest.Assert(b1, b2)
	})

	// Unmarshal
	gtest.Case(t, func() {
		s := "i love gf"
		v := gvar.New(nil)
		b, err := json.Marshal(s)
		gtest.Assert(err, nil)

		err = json.Unmarshal(b, v)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), s)
	})

	gtest.Case(t, func() {
		var v gvar.Var
		s := "i love gf"
		b, err := json.Marshal(s)
		gtest.Assert(err, nil)

		err = json.Unmarshal(b, &v)
		gtest.Assert(err, nil)
		gtest.Assert(v.String(), s)
	})
}

func Test_UnmarshalValue(t *testing.T) {
	type T struct {
		Name string
		Var  *gvar.Var
	}
	gtest.Case(t, func() {
		var t *T
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &t)
		gtest.Assert(err, nil)
		gtest.Assert(t.Name, "john")
		gtest.Assert(t.Var.String(), "v")
	})
}
