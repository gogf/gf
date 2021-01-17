// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"bytes"
	"encoding/binary"
	"github.com/gogf/gf/util/gconv"
	"testing"
	"time"

	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/test/gtest"
)

func Test_Set(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var v gvar.Var
		v.Set(123.456)
		t.Assert(v.Val(), 123.456)
	})
	gtest.C(t, func(t *gtest.T) {
		var v gvar.Var
		v.Set(123.456)
		t.Assert(v.Val(), 123.456)
	})
	gtest.C(t, func(t *gtest.T) {
		v := gvar.Create(123.456)
		t.Assert(v.Val(), 123.456)
	})
	gtest.C(t, func(t *gtest.T) {
		objOne := gvar.New("old", true)
		objOneOld, _ := objOne.Set("new").(string)
		t.Assert(objOneOld, "old")

		objTwo := gvar.New("old", false)
		objTwoOld, _ := objTwo.Set("new").(string)
		t.Assert(objTwoOld, "old")
	})
}

func Test_Val(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Val().(int)
		t.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Val().(int)
		t.Assert(objTwoOld, 1)
	})
}
func Test_Interface(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		objOne := gvar.New(1, true)
		objOneOld, _ := objOne.Interface().(int)
		t.Assert(objOneOld, 1)

		objTwo := gvar.New(1, false)
		objTwoOld, _ := objTwo.Interface().(int)
		t.Assert(objTwoOld, 1)
	})
}
func Test_IsNil(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		objOne := gvar.New(nil, true)
		t.Assert(objOne.IsNil(), true)

		objTwo := gvar.New("noNil", false)
		t.Assert(objTwo.IsNil(), false)

	})
}

func Test_Bytes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		x := int32(1)
		bytesBuffer := bytes.NewBuffer([]byte{})
		binary.Write(bytesBuffer, binary.BigEndian, x)

		objOne := gvar.New(bytesBuffer.Bytes(), true)

		bBuf := bytes.NewBuffer(objOne.Bytes())
		var y int32
		binary.Read(bBuf, binary.BigEndian, &y)

		t.Assert(x, y)

	})
}

func Test_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var str string = "hello"
		objOne := gvar.New(str, true)
		t.Assert(objOne.String(), str)

	})
}
func Test_Bool(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var ok bool = true
		objOne := gvar.New(ok, true)
		t.Assert(objOne.Bool(), ok)

		ok = false
		objTwo := gvar.New(ok, true)
		t.Assert(objTwo.Bool(), ok)

	})
}

func Test_Int(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num int = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Int(), num)

	})
}

func Test_Int8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num int8 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Int8(), num)

	})
}

func Test_Int16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num int16 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Int16(), num)

	})
}

func Test_Int32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num int32 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Int32(), num)

	})
}

func Test_Int64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num int64 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Int64(), num)

	})
}

func Test_Uint(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num uint = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Uint(), num)

	})
}

func Test_Uint8(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num uint8 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Uint8(), num)

	})
}

func Test_Uint16(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num uint16 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Uint16(), num)

	})
}

func Test_Uint32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num uint32 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Uint32(), num)

	})
}

func Test_Uint64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num uint64 = 1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Uint64(), num)

	})
}
func Test_Float32(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num float32 = 1.1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Float32(), num)

	})
}

func Test_Float64(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var num float64 = 1.1
		objOne := gvar.New(num, true)
		t.Assert(objOne.Float64(), num)

	})
}

func Test_Ints(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Ints()[0], arr[0])
	})
}
func Test_Floats(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []float64{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Floats()[0], arr[0])
	})
}
func Test_Strings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []string{"hello", "world"}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Strings()[0], arr[0])
	})
}

func Test_Interfaces(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Interfaces(), arr)
	})
}

func Test_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, true)
		t.Assert(objOne.Slice(), arr)
	})
}

func Test_Array(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		t.Assert(objOne.Array(), arr)
	})
}

func Test_Vars(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var arr = []int{1, 2, 3, 4, 5}
		objOne := gvar.New(arr, false)
		t.Assert(len(objOne.Vars()), 5)
		t.Assert(objOne.Vars()[0].Int(), 1)
		t.Assert(objOne.Vars()[4].Int(), 5)
	})
}

func Test_Time(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		t.Assert(objOne.Time().Unix(), timeUnix)
	})
}

func Test_GTime(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		t.Assert(objOne.GTime().Unix(), timeUnix)
	})
}

func Test_Duration(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var timeUnix int64 = 1556242660
		objOne := gvar.New(timeUnix, true)
		t.Assert(objOne.Duration(), time.Duration(timeUnix))
	})
}

func Test_UnmarshalValue(t *testing.T) {
	type V struct {
		Name string
		Var  *gvar.Var
	}
	gtest.C(t, func(t *gtest.T) {
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.Assert(err, nil)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
}
