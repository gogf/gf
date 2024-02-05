// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs_test

import (
	"testing"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gstructs"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user User
		m, _ := gstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = gstructs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithBase struct {
			Id   int
			Name string
			Base `params:"base"`
		}
		user := new(UserWithBase)
		m, _ := gstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{
			"base":      "Base",
			"password1": "Pass1",
			"password2": "Pass2",
		})
	})

	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Pass1 string `params:"password1"`
			Pass2 string `params:"password2"`
		}
		type UserWithEmbeddedAttribute struct {
			Id   int
			Name string
			Base
		}
		type UserWithoutEmbeddedAttribute struct {
			Id   int
			Name string
			Pass Base
		}
		user1 := new(UserWithEmbeddedAttribute)
		user2 := new(UserWithoutEmbeddedAttribute)
		m, _ := gstructs.TagMapName(user1, []string{"params"})
		t.Assert(m, g.Map{"password1": "Pass1", "password2": "Pass2"})
		m, _ = gstructs.TagMapName(user2, []string{"params"})
		t.Assert(m, g.Map{})
	})
}

func Test_StructOfNilPointer(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gstructs.TagMapName(user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"params"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})

		m, _ = gstructs.TagMapName(&user, []string{"params", "my-tag1"})
		t.Assert(m, g.Map{"name": "Name", "pass": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"my-tag1", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass1": "Pass"})
		m, _ = gstructs.TagMapName(&user, []string{"my-tag2", "params"})
		t.Assert(m, g.Map{"name": "Name", "pass2": "Pass"})
	})
}

func Test_Fields(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		fields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         user,
			RecursiveOption: 0,
		})
		t.Assert(len(fields), 3)
		t.Assert(fields[0].Name(), "Id")
		t.Assert(fields[1].Name(), "Name")
		t.Assert(fields[1].Tag("params"), "name")
		t.Assert(fields[2].Name(), "Pass")
		t.Assert(fields[2].Tag("my-tag1"), "pass1")
		t.Assert(fields[2].Tag("my-tag2"), "pass2")
		t.Assert(fields[2].Tag("params"), "pass")
	})
}

func Test_Fields_WithEmbedded1(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
			Age  int
		}
		type A struct {
			Site  string
			B     // Should be put here to validate its index.
			Score int64
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Site`)
		t.Assert(r[1].Name(), `Name`)
		t.Assert(r[2].Name(), `Age`)
		t.Assert(r[3].Name(), `Score`)
	})
}

func Test_Fields_WithEmbedded2(t *testing.T) {
	type MetaNode struct {
		Id          uint   `orm:"id,primary"  description:""`
		Capacity    string `orm:"capacity"    description:"Capacity string"`
		Allocatable string `orm:"allocatable" description:"Allocatable string"`
		Status      string `orm:"status"      description:"Status string"`
	}
	type MetaNodeZone struct {
		Nodes    uint
		Clusters uint
		Disk     uint
		Cpu      uint
		Memory   uint
		Zone     string
	}

	type MetaNodeItem struct {
		MetaNode
		Capacity    []MetaNodeZone `dc:"Capacity []MetaNodeZone"`
		Allocatable []MetaNodeZone `dc:"Allocatable []MetaNodeZone"`
	}

	gtest.C(t, func(t *gtest.T) {
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(MetaNodeItem),
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Id`)
		t.Assert(r[1].Name(), `Capacity`)
		t.Assert(r[1].TagStr(), `dc:"Capacity []MetaNodeZone"`)
		t.Assert(r[2].Name(), `Allocatable`)
		t.Assert(r[2].TagStr(), `dc:"Allocatable []MetaNodeZone"`)
		t.Assert(r[3].Name(), `Status`)
	})
}

// Filter repeated fields when there is embedded struct.
func Test_Fields_WithEmbedded_Filter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
			Age  int
		}
		type A struct {
			Name  string
			Site  string
			Age   string
			B     // Should be put here to validate its index.
			Score int64
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: gstructs.RecursiveOptionEmbeddedNoTag,
		})
		t.AssertNil(err)
		t.Assert(len(r), 4)
		t.Assert(r[0].Name(), `Name`)
		t.Assert(r[1].Name(), `Site`)
		t.Assert(r[2].Name(), `Age`)
		t.Assert(r[3].Name(), `Score`)
	})
}

func Test_FieldMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gstructs.FieldMap(gstructs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: []string{"params"},
			RecursiveOption:  gstructs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, false)
		_, ok = m["name"]
		t.Assert(ok, true)
		_, ok = m["Pass"]
		t.Assert(ok, false)
		_, ok = m["pass"]
		t.Assert(ok, true)
	})
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id   int
			Name string `params:"name"`
			Pass string `my-tag1:"pass1" my-tag2:"pass2" params:"pass"`
		}
		var user *User
		m, _ := gstructs.FieldMap(gstructs.FieldMapInput{
			Pointer:          user,
			PriorityTagArray: nil,
			RecursiveOption:  gstructs.RecursiveOptionEmbedded,
		})
		t.Assert(len(m), 3)
		_, ok := m["Id"]
		t.Assert(ok, true)
		_, ok = m["Name"]
		t.Assert(ok, true)
		_, ok = m["name"]
		t.Assert(ok, false)
		_, ok = m["Pass"]
		t.Assert(ok, true)
		_, ok = m["pass"]
		t.Assert(ok, false)
	})
}

func Test_StructType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := gstructs.StructType(new(A))
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/v2/os/gstructs_test/gstructs_test.A`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			B
		}
		r, err := gstructs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/v2/os/gstructs_test/gstructs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
		}
		r, err := gstructs.StructType(new(A).B)
		t.AssertNil(err)
		t.Assert(r.String(), `gstructs_test.B`)
	})
	// Error.
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			*B
			Id int
		}
		_, err := gstructs.StructType(new(A).Id)
		t.AssertNE(err, nil)
	})
}

func Test_StructTypeBySlice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := gstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/v2/os/gstructs_test/gstructs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array []B
		}
		r, err := gstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/v2/os/gstructs_test/gstructs_test.B`)
	})
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Name string
		}
		type A struct {
			Array *[]B
		}
		r, err := gstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.Signature(), `github.com/gogf/gf/v2/os/gstructs_test/gstructs_test.B`)
	})
}

func TestType_FieldKeys(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type B struct {
			Id   int
			Name string
		}
		type A struct {
			Array []*B
		}
		r, err := gstructs.StructType(new(A).Array)
		t.AssertNil(err)
		t.Assert(r.FieldKeys(), g.Slice{"Id", "Name"})
	})
}

func TestType_TagMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Id   int    `d:"123" description:"I love gf"`
			Name string `v:"required" description:"应用Id"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagMap()["d"], `123`)
		t.Assert(r[0].TagMap()["description"], `I love gf`)
		t.Assert(r[1].TagMap()["v"], `required`)
		t.Assert(r[1].TagMap()["description"], `应用Id`)
	})
}

func TestType_TagJsonName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name string `json:"name,omitempty"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 1)
		t.Assert(r[0].TagJsonName(), `name`)
	})
}

func TestType_TagDefault(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `default:"john"`
			Name2 string `d:"john"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagDefault(), `john`)
		t.Assert(r[1].TagDefault(), `john`)
	})
}

func TestType_TagParam(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `param:"name"`
			Name2 string `p:"name"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagParam(), `name`)
		t.Assert(r[1].TagParam(), `name`)
	})
}

func TestType_TagValid(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `valid:"required"`
			Name2 string `v:"required"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagValid(), `required`)
		t.Assert(r[1].TagValid(), `required`)
	})
}

func TestType_TagDescription(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `description:"my name"`
			Name2 string `des:"my name"`
			Name3 string `dc:"my name"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0].TagDescription(), `my name`)
		t.Assert(r[1].TagDescription(), `my name`)
		t.Assert(r[2].TagDescription(), `my name`)
	})
}

func TestType_TagSummary(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `summary:"my name"`
			Name2 string `sum:"my name"`
			Name3 string `sm:"my name"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 3)
		t.Assert(r[0].TagSummary(), `my name`)
		t.Assert(r[1].TagSummary(), `my name`)
		t.Assert(r[2].TagSummary(), `my name`)
	})
}

func TestType_TagAdditional(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `additional:"my name"`
			Name2 string `ad:"my name"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagAdditional(), `my name`)
		t.Assert(r[1].TagAdditional(), `my name`)
	})
}

func TestType_TagExample(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type A struct {
			Name  string `example:"john"`
			Name2 string `eg:"john"`
		}
		r, err := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         new(A),
			RecursiveOption: 0,
		})
		t.AssertNil(err)

		t.Assert(len(r), 2)
		t.Assert(r[0].TagExample(), `john`)
		t.Assert(r[1].TagExample(), `john`)
	})
}

func Test_Fields_TagPriorityName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Name  string `gconv:"name_gconv" c:"name_c"`
			Age   uint   `p:"name_p" param:"age_param"`
			Pass  string `json:"pass_json"`
			IsMen bool
		}
		var user *User
		fields, _ := gstructs.Fields(gstructs.FieldsInput{
			Pointer:         user,
			RecursiveOption: 0,
		})
		t.Assert(fields[0].TagPriorityName(), "name_gconv")
		t.Assert(fields[1].TagPriorityName(), "age_param")
		t.Assert(fields[2].TagPriorityName(), "pass_json")
		t.Assert(fields[3].TagPriorityName(), "IsMen")
	})
}
