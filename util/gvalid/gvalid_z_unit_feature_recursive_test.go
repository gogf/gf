// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_CheckStruct_Recursive_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id   int
			Name string `v:"required"`
			Pass Pass
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := g.Validator().Data(user).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `1` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `2` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursive_Struct_WithData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id   int
			Name string `v:"required"`
			Pass Pass
		}
		user := &User{}
		data := g.Map{
			"Name": "john",
			"Pass": g.Map{
				"Pass1": 100,
				"Pass2": 200,
			},
		}
		err := g.Validator().Data(user).Assoc(data).Run(ctx)
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], nil)
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `100` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `200` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursive_SliceStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id     int
			Name   string `v:"required"`
			Passes []Pass
		}
		user := &User{
			Name: "",
			Passes: []Pass{
				{
					Pass1: "1",
					Pass2: "2",
				},
				{
					Pass1: "3",
					Pass2: "4",
				},
			},
		}
		err := g.Validator().Data(user).Run(ctx)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `3` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `4` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursive_SliceStruct_Bail(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id     int
			Name   string `v:"required"`
			Passes []Pass
		}
		user := &User{
			Name: "",
			Passes: []Pass{
				{
					Pass1: "1",
					Pass2: "2",
				},
				{
					Pass1: "3",
					Pass2: "4",
				},
			},
		}
		err := g.Validator().Bail().Data(user).Run(ctx)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], nil)
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `1` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], nil)
	})
}

func Test_CheckStruct_Recursive_SliceStruct_Required(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id     int
			Name   string `v:"required"`
			Passes []Pass
		}
		user := &User{}
		err := g.Validator().Data(user).Run(ctx)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.Assert(err.Maps()["Pass1"], nil)
		t.Assert(err.Maps()["Pass2"], nil)
	})
}

func Test_CheckStruct_Recursive_MapStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		type User struct {
			Id     int
			Name   string `v:"required"`
			Passes map[string]Pass
		}
		user := &User{
			Name: "",
			Passes: map[string]Pass{
				"test1": {
					Pass1: "1",
					Pass2: "2",
				},
				"test2": {
					Pass1: "3",
					Pass2: "4",
				},
			},
		}
		err := g.Validator().Data(user).Run(ctx)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], g.Map{"required": "The Name field is required"})
		t.AssertNE(err.Maps()["Pass1"], nil)
		t.AssertNE(err.Maps()["Pass2"], nil)
	})
}

func Test_CheckMap_Recursive_SliceStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `v:"required|same:Pass2"`
			Pass2 string `v:"required|same:Pass1"`
		}
		user := g.Map{
			"Name": "",
			"Pass": []Pass{
				{
					Pass1: "1",
					Pass2: "2",
				},
				{
					Pass1: "3",
					Pass2: "4",
				},
			},
		}
		err := g.Validator().Data(user).Run(ctx)
		g.Dump(err.Items())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["Name"], nil)
		t.Assert(err.Maps()["Pass1"], g.Map{"same": "The Pass1 value `3` must be the same as field Pass2"})
		t.Assert(err.Maps()["Pass2"], g.Map{"same": "The Pass2 value `4` must be the same as field Pass1"})
	})
}

func Test_CheckStruct_Recursively_SliceAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required#Student Name is required"`
			Age  int    `v:"required"`
		}
		type Teacher struct {
			Name     string    `v:"required#Teacher Name is required"`
			Students []Student `v:"required"`
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name":     "john",
				"students": `[]`,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `The Students field is required`)
	})

	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required#Student Name is required"`
			Age  int    `v:"required"`
		}
		type Teacher struct {
			Name     string `v:"required#Teacher Name is required"`
			Students []Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name": "john",
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, ``)
	})

	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required#Student Name is required"`
			Age  int    `v:"required"`
		}
		type Teacher struct {
			Name     string    `v:"required#Teacher Name is required"`
			Students []Student `v:"required"`
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name":     "john",
				"students": `[{"age":2}, {"name":"jack", "age":4}]`,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `Student Name is required`)
	})

	// https://github.com/gogf/gf/issues/1864
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required"`
			Age  int
		}
		type Teacher struct {
			Name     string
			Students []*Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name":     "john",
				"students": `[{"age":2},{"name":"jack", "age":4}]`,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `The Name field is required`)
	})
}

func Test_CheckStruct_Recursively_SliceAttribute_WithTypeAlias(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type ParamsItemBase struct {
			Component string `v:"required" dc:"组件名称"`
			Params    string `v:"required" dc:"配置参数(一般是JSON)"`
			Version   uint64 `v:"required" dc:"参数版本"`
		}
		type ParamsItem = ParamsItemBase
		type ParamsModifyReq struct {
			Revision  uint64       `v:"required"`
			BizParams []ParamsItem `v:"required"`
		}
		var (
			req  = ParamsModifyReq{}
			data = g.Map{
				"Revision":  "1",
				"BizParams": `[{}]`,
			}
		)
		err := g.Validator().Assoc(data).Data(req).Run(ctx)
		t.Assert(err, `The Component field is required; The Params field is required; The Version field is required`)
	})
}

func Test_CheckStruct_Recursively_MapAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required#Student Name is required"`
			Age  int    `v:"required"`
		}
		type Teacher struct {
			Name     string             `v:"required#Teacher Name is required"`
			Students map[string]Student `v:"required"`
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name":     "john",
				"students": `{"john":{"age":18}}`,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `Student Name is required`)
	})
}

// https://github.com/gogf/gf/issues/1983
func Test_Issue1983(t *testing.T) {
	// Error as the attribute Student in Teacher is an initialized struct, which has default value.
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required"`
			Age  int
		}
		type Teacher struct {
			Students Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"students": nil,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `The Name field is required`)
	})
	// The same as upper, it is not affected by association values.
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required"`
			Age  int
		}
		type Teacher struct {
			Students Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"name":     "john",
				"students": nil,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, `The Name field is required`)
	})
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required"`
			Age  int
		}
		type Teacher struct {
			Students *Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"students": nil,
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.AssertNil(err)
	})
}

// https://github.com/gogf/gf/issues/1921
func Test_Issue1921(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type SearchOption struct {
			Size int `v:"max:100"`
		}
		type SearchReq struct {
			Option *SearchOption `json:"option,omitempty"`
		}

		var (
			req = SearchReq{
				Option: &SearchOption{
					Size: 10000,
				},
			}
		)
		err := g.Validator().Data(req).Run(ctx)
		t.Assert(err, "The Size value `10000` must be equal or lesser than 100")
	})
}

// https://github.com/gogf/gf/issues/2011
func Test_Issue2011(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Student struct {
			Name string `v:"required|min-length:6"`
			Age  int
		}
		type Teacher struct {
			Student *Student
		}
		var (
			teacher = Teacher{}
			data    = g.Map{
				"student": g.Map{
					"name": "john",
				},
			}
		)
		err := g.Validator().Assoc(data).Data(teacher).Run(ctx)
		t.Assert(err, "The Name value `john` length must be equal or greater than 6")
	})
}
