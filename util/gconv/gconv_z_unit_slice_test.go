// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value := 123.456
		t.AssertEQ(gconv.Bytes("123"), []byte("123"))
		t.AssertEQ(gconv.Bytes([]interface{}{1}), []byte{1})
		t.AssertEQ(gconv.Bytes([]interface{}{300}), []byte("[300]"))
		t.AssertEQ(gconv.Strings(value), []string{"123.456"})
		t.AssertEQ(gconv.SliceStr(value), []string{"123.456"})
		t.AssertEQ(gconv.SliceInt(value), []int{123})
		t.AssertEQ(gconv.SliceUint(value), []uint{123})
		t.AssertEQ(gconv.SliceUint32(value), []uint32{123})
		t.AssertEQ(gconv.SliceUint64(value), []uint64{123})
		t.AssertEQ(gconv.SliceInt32(value), []int32{123})
		t.AssertEQ(gconv.SliceInt64(value), []int64{123})
		t.AssertEQ(gconv.Ints(value), []int{123})
		t.AssertEQ(gconv.SliceFloat(value), []float64{123.456})
		t.AssertEQ(gconv.Floats(value), []float64{123.456})
		t.AssertEQ(gconv.SliceFloat32(value), []float32{123.456})
		t.AssertEQ(gconv.SliceFloat64(value), []float64{123.456})
		t.AssertEQ(gconv.Interfaces(value), []interface{}{123.456})
		t.AssertEQ(gconv.SliceAny(" [26, 27] "), []interface{}{26, 27})
	})
	gtest.C(t, func(t *gtest.T) {
		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.SliceInt64(s), []int64{1, 2})
	})
}

func Test_Slice_Ints(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Ints(nil), nil)
		t.AssertEQ(gconv.Ints("[26, 27]"), []int{26, 27})
		t.AssertEQ(gconv.Ints(" [26, 27] "), []int{26, 27})
		t.AssertEQ(gconv.Ints([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int{0, 0})
		t.AssertEQ(gconv.Ints([]bool{true, false}), []int{1, 0})
		t.AssertEQ(gconv.Ints([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []int{1, 2})
	})
}

func Test_Slice_Int32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Int32s(nil), nil)
		t.AssertEQ(gconv.Int32s(" [26, 27] "), []int32{26, 27})
		t.AssertEQ(gconv.Int32s([]string{"1", "2"}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]int{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]int8{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]int16{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]int32{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]int64{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]uint{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]uint8{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int32{0, 0})
		t.AssertEQ(gconv.Int32s([]uint16{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]uint32{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]uint64{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]bool{true, false}), []int32{1, 0})
		t.AssertEQ(gconv.Int32s([]float32{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([]float64{1, 2}), []int32{1, 2})
		t.AssertEQ(gconv.Int32s([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []int32{1, 2})

		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.SliceInt32(s), []int32{1, 2})
	})
}

func Test_Slice_Int64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Int64s(nil), nil)
		t.AssertEQ(gconv.Int64s(" [26, 27] "), []int64{26, 27})
		t.AssertEQ(gconv.Int64s([]string{"1", "2"}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]int{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]int8{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]int16{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]int32{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]int64{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]uint{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]uint8{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []int64{0, 0})
		t.AssertEQ(gconv.Int64s([]uint16{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]uint32{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]uint64{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]bool{true, false}), []int64{1, 0})
		t.AssertEQ(gconv.Int64s([]float32{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([]float64{1, 2}), []int64{1, 2})
		t.AssertEQ(gconv.Int64s([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []int64{1, 2})

		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.Int64s(s), []int64{1, 2})
	})
}

func Test_Slice_Uints(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uints(nil), nil)
		t.AssertEQ(gconv.Uints("1"), []uint{1})
		t.AssertEQ(gconv.Uints(" [26, 27] "), []uint{26, 27})
		t.AssertEQ(gconv.Uints([]string{"1", "2"}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]int{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]int8{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]int16{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]int32{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]int64{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]uint{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]uint8{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint{0, 0})
		t.AssertEQ(gconv.Uints([]uint16{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]uint32{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]uint64{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]bool{true, false}), []uint{1, 0})
		t.AssertEQ(gconv.Uints([]float32{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([]float64{1, 2}), []uint{1, 2})
		t.AssertEQ(gconv.Uints([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []uint{1, 2})

		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.Uints(s), []uint{1, 2})
	})
}

func Test_Slice_Uint32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uint32s(nil), nil)
		t.AssertEQ(gconv.Uint32s("1"), []uint32{1})
		t.AssertEQ(gconv.Uint32s(" [26, 27] "), []uint32{26, 27})
		t.AssertEQ(gconv.Uint32s([]string{"1", "2"}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]int{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]int8{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]int16{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]int32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]int64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]uint{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]uint8{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint32{0, 0})
		t.AssertEQ(gconv.Uint32s([]uint16{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]uint32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]uint64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]bool{true, false}), []uint32{1, 0})
		t.AssertEQ(gconv.Uint32s([]float32{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([]float64{1, 2}), []uint32{1, 2})
		t.AssertEQ(gconv.Uint32s([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []uint32{1, 2})

		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.Uint32s(s), []uint32{1, 2})
	})
}

func Test_Slice_Uint64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uint64s(nil), nil)
		t.AssertEQ(gconv.Uint64s("1"), []uint64{1})
		t.AssertEQ(gconv.Uint64s(" [26, 27] "), []uint64{26, 27})
		t.AssertEQ(gconv.Uint64s([]string{"1", "2"}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]int{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]int8{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]int16{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]int32{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]int64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]uint{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]uint8{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []uint64{0, 0})
		t.AssertEQ(gconv.Uint64s([]uint16{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]uint64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]uint64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]bool{true, false}), []uint64{1, 0})
		t.AssertEQ(gconv.Uint64s([]float32{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([]float64{1, 2}), []uint64{1, 2})
		t.AssertEQ(gconv.Uint64s([][]byte{[]byte{byte(1)}, []byte{byte(2)}}), []uint64{1, 2})

		s := []*gvar.Var{
			gvar.New(1),
			gvar.New(2),
		}
		t.AssertEQ(gconv.Uint64s(s), []uint64{1, 2})
	})
}

func Test_Slice_Float32s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Float32s("123.4"), []float32{123.4})
		t.AssertEQ(gconv.Float32s([]string{"123.4", "123.5"}), []float32{123.4, 123.5})
		t.AssertEQ(gconv.Float32s([]int{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]int8{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]int16{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]int32{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]int64{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]uint{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]uint8{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []float32{0, 0})
		t.AssertEQ(gconv.Float32s([]uint16{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]uint32{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]uint64{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]bool{true, false}), []float32{0, 0})
		t.AssertEQ(gconv.Float32s([]float32{123}), []float32{123})
		t.AssertEQ(gconv.Float32s([]float64{123}), []float32{123})

		s := []*gvar.Var{
			gvar.New(1.1),
			gvar.New(2.1),
		}
		t.AssertEQ(gconv.SliceFloat32(s), []float32{1.1, 2.1})
	})
}

func Test_Slice_Float64s(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Float64s("123.4"), []float64{123.4})
		t.AssertEQ(gconv.Float64s([]string{"123.4", "123.5"}), []float64{123.4, 123.5})
		t.AssertEQ(gconv.Float64s([]int{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]int8{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]int16{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]int32{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]int64{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]uint{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]uint8{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`)), []float64{0, 0})
		t.AssertEQ(gconv.Float64s([]uint16{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]uint32{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]uint64{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]bool{true, false}), []float64{0, 0})
		t.AssertEQ(gconv.Float64s([]float32{123}), []float64{123})
		t.AssertEQ(gconv.Float64s([]float64{123}), []float64{123})
	})
}

func Test_Slice_Empty(t *testing.T) {
	// Int.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Ints(""), []int{})
		t.Assert(gconv.Ints(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Int32s(""), []int32{})
		t.Assert(gconv.Int32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Int64s(""), []int64{})
		t.Assert(gconv.Int64s(nil), nil)
	})
	// Uint.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uints(""), []uint{})
		t.Assert(gconv.Uints(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uint32s(""), []uint32{})
		t.Assert(gconv.Uint32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Uint64s(""), []uint64{})
		t.Assert(gconv.Uint64s(nil), nil)
	})
	// Float.
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Floats(""), []float64{})
		t.Assert(gconv.Floats(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Float32s(""), []float32{})
		t.Assert(gconv.Float32s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Float64s(""), []float64{})
		t.Assert(gconv.Float64s(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Strings(""), []string{})
		t.Assert(gconv.Strings(nil), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.SliceAny(""), []interface{}{})
		t.Assert(gconv.SliceAny(nil), nil)
	})
}

func Test_Strings(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := []*g.Var{
			g.NewVar(1),
			g.NewVar(2),
			g.NewVar(3),
		}
		t.AssertEQ(gconv.Strings(array), []string{"1", "2", "3"})

		t.AssertEQ(gconv.Strings([]uint8(`["1","2"]`)), []string{"1", "2"})
		t.AssertEQ(gconv.Strings([][]byte{{byte(0)}, {byte(1)}}), []string{"\u0000", "\u0001"})
	})
	// https://github.com/gogf/gf/issues/1750
	gtest.C(t, func(t *gtest.T) {
		t.AssertEQ(gconv.Strings("123"), []string{"123"})
	})
}

func Test_Slice_Interfaces(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := gconv.Interfaces([]uint8(`[{"id": 1, "name":"john"},{"id": 2, "name":"huang"}]`))
		t.Assert(len(array), 2)
		t.Assert(array[0].(g.Map)["id"], 1)
		t.Assert(array[0].(g.Map)["name"], "john")
	})
	// map
	gtest.C(t, func(t *gtest.T) {
		array := gconv.Interfaces(g.Map{
			"id":   1,
			"name": "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(g.Map)["id"], 1)
		t.Assert(array[0].(g.Map)["name"], "john")
	})
	// struct
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Id   int `json:"id"`
			Name string
		}
		array := gconv.Interfaces(&A{
			Id:   1,
			Name: "john",
		})
		t.Assert(len(array), 1)
		t.Assert(array[0].(*A).Id, 1)
		t.Assert(array[0].(*A).Name, "john")
	})
}

func Test_Slice_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int    `json:"id"`
		name string `json:"name"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john"}
		array := gconv.Interfaces(user)
		t.Assert(len(array), 1)
		t.Assert(array[0].(*User).Id, 1)
		t.Assert(array[0].(*User).name, "john")
	})
}

func Test_Slice_Structs(t *testing.T) {
	type Base struct {
		Age int
	}
	type User struct {
		Id   int
		Name string
		Base
	}

	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.Structs(params, &users)
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})

	gtest.C(t, func(t *gtest.T) {
		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.StructsTag(params, &users, "")
		t.AssertNil(err)
		t.Assert(len(users), 2)
		t.Assert(users[0].Id, params[0]["id"])
		t.Assert(users[0].Name, params[0]["name"])
		t.Assert(users[0].Age, 18)

		t.Assert(users[1].Id, params[1]["id"])
		t.Assert(users[1].Name, params[1]["name"])
		t.Assert(users[1].Age, 20)
	})
}
