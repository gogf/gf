// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvar_test

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
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

		objOne = nil
		t.Assert(objOne.Val(), nil)
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

func Test_UnmarshalJson(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type V struct {
			Name string
			Var  *gvar.Var
		}
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
	gtest.C(t, func(t *gtest.T) {
		type V struct {
			Name string
			Var  gvar.Var
		}
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
}

func Test_UnmarshalValue(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type V struct {
			Name string
			Var  *gvar.Var
		}
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
	gtest.C(t, func(t *gtest.T) {
		type V struct {
			Name string
			Var  gvar.Var
		}
		var v *V
		err := gconv.Struct(map[string]interface{}{
			"name": "john",
			"var":  "v",
		}, &v)
		t.AssertNil(err)
		t.Assert(v.Name, "john")
		t.Assert(v.Var.String(), "v")
	})
}

func Test_Copy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		srcVar := gvar.New(src)
		dstVar := srcVar.Copy()
		t.Assert(srcVar.Map(), src)
		t.Assert(dstVar.Map(), src)

		dstVar.Map()["k3"] = "v3"
		t.Assert(srcVar.Map(), g.Map{
			"k1": "v1",
			"k2": "v2",
		})
		t.Assert(dstVar.Map(), g.Map{
			"k1": "v1",
			"k2": "v2",
			"k3": "v3",
		})
	})
}

func Test_DeepCopy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		src := g.Map{
			"k1": "v1",
			"k2": "v2",
		}
		srcVar := gvar.New(src)
		copyVar := srcVar.DeepCopy().(*gvar.Var)
		copyVar.Set(g.Map{
			"k3": "v3",
			"k4": "v4",
		})
		t.AssertNE(srcVar, copyVar)

		srcVar = nil
		t.AssertNil(srcVar.DeepCopy())
	})
}
