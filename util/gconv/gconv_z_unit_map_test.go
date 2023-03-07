// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gconv_test

import (
	"encoding/json"
	"testing"

	"gopkg.in/yaml.v3"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

func Test_Map_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		m1 := map[string]string{
			"k": "v",
		}
		m2 := map[int]string{
			3: "v",
		}
		m3 := map[float64]float32{
			1.22: 3.1,
		}
		t.Assert(gconv.Map(m1), g.Map{
			"k": "v",
		})
		t.Assert(gconv.Map(m2), g.Map{
			"3": "v",
		})
		t.Assert(gconv.Map(m3), g.Map{
			"1.22": "3.1",
		})
		t.Assert(gconv.Map(`{"name":"goframe"}`), g.Map{
			"name": "goframe",
		})
		t.Assert(gconv.Map(`{"name":"goframe"`), nil)
		t.Assert(gconv.Map(`{goframe}`), nil)
		t.Assert(gconv.Map([]byte(`{"name":"goframe"}`)), g.Map{
			"name": "goframe",
		})
		t.Assert(gconv.Map([]byte(`{"name":"goframe"`)), nil)
		t.Assert(gconv.Map([]byte(`{goframe}`)), nil)
	})
}

func Test_Map_Slice(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		slice1 := g.Slice{"1", "2", "3", "4"}
		slice2 := g.Slice{"1", "2", "3"}
		slice3 := g.Slice{}
		t.Assert(gconv.Map(slice1), g.Map{
			"1": "2",
			"3": "4",
		})
		t.Assert(gconv.Map(slice2), g.Map{
			"1": "2",
			"3": nil,
		})
		t.Assert(gconv.Map(slice3), g.Map{})
	})
}

func Test_Maps_Basic(t *testing.T) {
	params := g.Slice{
		g.Map{"id": 100, "name": "john"},
		g.Map{"id": 200, "name": "smith"},
	}
	gtest.C(t, func(t *gtest.T) {
		list := gconv.Maps(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gtest.C(t, func(t *gtest.T) {
		list := gconv.SliceMap(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gtest.C(t, func(t *gtest.T) {
		list := gconv.SliceMapDeep(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Age int
		}
		type User struct {
			Id   int
			Name string
			Base
		}

		users := make([]User, 0)
		params := []g.Map{
			{"id": 1, "name": "john", "age": 18},
			{"id": 2, "name": "smith", "age": 20},
		}
		err := gconv.SliceStruct(params, &users)
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

func Test_Maps_JsonStr(t *testing.T) {
	jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`
	gtest.C(t, func(t *gtest.T) {
		list := gconv.Maps(jsonStr)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gconv.Maps([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.Maps(`[id]`), nil)
		t.Assert(gconv.Maps(`test`), nil)
		t.Assert(gconv.Maps([]byte(`[id]`)), nil)
		t.Assert(gconv.Maps([]byte(`test`)), nil)
	})
}

func Test_Map_StructWithGConvTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `gconv:"-"`
			NickName string `gconv:"nickname, omitempty"`
			Pass1    string `gconv:"password1"`
			Pass2    string `gconv:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gconv.Map(user1)
		map2 := gconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithJsonTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `json:"-"`
			NickName string `json:"nickname, omitempty"`
			Pass1    string `json:"password1"`
			Pass2    string `json:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gconv.Map(user1)
		map2 := gconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_StructWithCTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Uid      int
			Name     string
			SiteUrl  string `c:"-"`
			NickName string `c:"nickname, omitempty"`
			Pass1    string `c:"password1"`
			Pass2    string `c:"password2"`
		}
		user1 := User{
			Uid:     100,
			Name:    "john",
			SiteUrl: "https://goframe.org",
			Pass1:   "123",
			Pass2:   "456",
		}
		user2 := &user1
		map1 := gconv.Map(user1)
		map2 := gconv.Map(user2)
		t.Assert(map1["Uid"], 100)
		t.Assert(map1["Name"], "john")
		t.Assert(map1["SiteUrl"], nil)
		t.Assert(map1["NickName"], nil)
		t.Assert(map1["nickname"], nil)
		t.Assert(map1["password1"], "123")
		t.Assert(map1["password2"], "456")

		t.Assert(map2["Uid"], 100)
		t.Assert(map2["Name"], "john")
		t.Assert(map2["SiteUrl"], nil)
		t.Assert(map2["NickName"], nil)
		t.Assert(map2["nickname"], nil)
		t.Assert(map2["password1"], "123")
		t.Assert(map2["password2"], "456")
	})
}

func Test_Map_PrivateAttribute(t *testing.T) {
	type User struct {
		Id   int
		name string
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{1, "john"}
		t.Assert(gconv.Map(user), g.Map{"Id": 1})
	})
}

func Test_Map_Embedded(t *testing.T) {
	type Base struct {
		Id int
	}
	type User struct {
		Base
		Name string
	}
	type UserDetail struct {
		User
		Brief string
	}
	gtest.C(t, func(t *gtest.T) {
		user := &User{}
		user.Id = 1
		user.Name = "john"

		m := gconv.Map(user)
		t.Assert(len(m), 2)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
	})
	gtest.C(t, func(t *gtest.T) {
		user := &UserDetail{}
		user.Id = 1
		user.Name = "john"
		user.Brief = "john guo"

		m := gconv.Map(user)
		t.Assert(len(m), 3)
		t.Assert(m["Id"], user.Id)
		t.Assert(m["Name"], user.Name)
		t.Assert(m["Brief"], user.Brief)
	})
}

func Test_Map_Embedded2(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gconv.Map(user)
		t.Assert(m["id"], "100")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "2019")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gconv.MapDeep(user)
		t.Assert(m["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], user.CreateTime)
	})
}

func Test_MapDeep2(t *testing.T) {
	type A struct {
		F string
		G string
	}

	type B struct {
		A
		H string
	}

	type C struct {
		A A
		F string
	}

	type D struct {
		I A
		F string
	}

	gtest.C(t, func(t *gtest.T) {
		b := new(B)
		c := new(C)
		d := new(D)
		mb := gconv.MapDeep(b)
		mc := gconv.MapDeep(c)
		md := gconv.MapDeep(d)
		t.Assert(gutil.MapContains(mb, "F"), true)
		t.Assert(gutil.MapContains(mb, "G"), true)
		t.Assert(gutil.MapContains(mb, "H"), true)
		t.Assert(gutil.MapContains(mc, "A"), true)
		t.Assert(gutil.MapContains(mc, "F"), true)
		t.Assert(gutil.MapContains(mc, "G"), false)
		t.Assert(gutil.MapContains(md, "F"), true)
		t.Assert(gutil.MapContains(md, "I"), true)
		t.Assert(gutil.MapContains(md, "H"), false)
		t.Assert(gutil.MapContains(md, "G"), false)
	})
}

func Test_MapDeep3(t *testing.T) {
	type Base struct {
		Id   int    `c:"id"`
		Date string `c:"date"`
	}
	type User struct {
		UserBase Base   `c:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}

	gtest.C(t, func(t *gtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := gconv.MapDeep(user)
		t.Assert(m, g.Map{
			"base": g.Map{
				"id":   user.UserBase.Id,
				"date": user.UserBase.Date,
			},
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})

	gtest.C(t, func(t *gtest.T) {
		user := &User{
			UserBase: Base{
				Id:   1,
				Date: "2019-10-01",
			},
			Passport: "john",
			Password: "123456",
			Nickname: "JohnGuo",
		}
		m := gconv.Map(user)
		t.Assert(m, g.Map{
			"base":     user.UserBase,
			"passport": user.Passport,
			"password": user.Password,
			"nickname": user.Nickname,
		})
	})
}

func Test_MapDeepWithAttributeTag(t *testing.T) {
	type Ids struct {
		Id  int `c:"id"`
		Uid int `c:"uid"`
	}
	type Base struct {
		Ids        `json:"ids"`
		CreateTime string `c:"create_time"`
	}
	type User struct {
		Base     `json:"base"`
		Passport string `c:"passport"`
		Password string `c:"password"`
		Nickname string `c:"nickname"`
	}
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gconv.Map(user)
		t.Assert(m["id"], "")
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["create_time"], "")
	})
	gtest.C(t, func(t *gtest.T) {
		user := new(User)
		user.Id = 100
		user.Nickname = "john"
		user.CreateTime = "2019"
		m := gconv.MapDeep(user)
		t.Assert(m["base"].(map[string]interface{})["ids"].(map[string]interface{})["id"], user.Id)
		t.Assert(m["nickname"], user.Nickname)
		t.Assert(m["base"].(map[string]interface{})["create_time"], user.CreateTime)
	})
}

func Test_MapDeepWithNestedMapAnyAny(t *testing.T) {
	type User struct {
		ExtraAttributes g.Map `c:"extra_attributes"`
	}

	gtest.C(t, func(t *gtest.T) {
		user := &User{
			ExtraAttributes: g.Map{
				"simple_attribute": 123,
				"map_string_attribute": g.Map{
					"inner_value": 456,
				},
				"map_interface_attribute": g.MapAnyAny{
					"inner_value": 456,
					123:           "integer_key_should_be_converted_to_string",
				},
			},
		}
		m := gconv.MapDeep(user)
		t.Assert(m, g.Map{
			"extra_attributes": g.Map{
				"simple_attribute": 123,
				"map_string_attribute": g.Map{
					"inner_value": user.ExtraAttributes["map_string_attribute"].(g.Map)["inner_value"],
				},
				"map_interface_attribute": g.Map{
					"inner_value": user.ExtraAttributes["map_interface_attribute"].(g.MapAnyAny)["inner_value"],
					"123":         "integer_key_should_be_converted_to_string",
				},
			},
		})
	})

	type Outer struct {
		OuterStruct map[string]interface{} `c:"outer_struct" yaml:"outer_struct"`
		Field3      map[string]interface{} `c:"field3" yaml:"field3"`
	}

	gtest.C(t, func(t *gtest.T) {
		problemYaml := []byte(`
outer_struct:
  field1: &anchor1
    inner1: 123
    inner2: 345
  field2: 
    inner3: 456
    inner4: 789
    <<: *anchor1
field3:
  123: integer_key
`)
		parsed := &Outer{}

		err := yaml.Unmarshal(problemYaml, parsed)
		t.AssertNil(err)

		_, err = json.Marshal(parsed)
		t.Assert(err.Error(), "json: unsupported type: map[interface {}]interface {}")

		converted := gconv.MapDeep(parsed)
		jsonData, err := json.Marshal(converted)
		t.AssertNil(err)

		t.Assert(string(jsonData), `{"field3":{"123":"integer_key"},"outer_struct":{"field1":{"inner1":123,"inner2":345},"field2":{"inner1":123,"inner2":345,"inner3":456,"inner4":789}}}`)
	})
}

func TestMapStrStr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapStrStr(map[string]string{"k": "v"}), map[string]string{"k": "v"})
		t.Assert(gconv.MapStrStr(`{}`), nil)
	})
}

func TestMapStrStrDeep(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapStrStrDeep(map[string]string{"k": "v"}), map[string]string{"k": "v"})
		t.Assert(gconv.MapStrStrDeep(`{"k":"v"}`), map[string]string{"k": "v"})
		t.Assert(gconv.MapStrStrDeep(`{}`), nil)
	})
}

func TestMapsDeep(t *testing.T) {
	jsonStr := `[{"id":100, "name":"john"},{"id":200, "name":"smith"}]`
	params := g.Slice{
		g.Map{"id": 100, "name": "john"},
		g.Map{"id": 200, "name": "smith"},
	}

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapsDeep(nil), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		list := gconv.MapsDeep(params)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gtest.C(t, func(t *gtest.T) {
		list := gconv.MapsDeep(jsonStr)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gconv.MapsDeep([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})

	gtest.C(t, func(t *gtest.T) {
		t.Assert(gconv.MapsDeep(`[id]`), nil)
		t.Assert(gconv.MapsDeep(`test`), nil)
		t.Assert(gconv.MapsDeep([]byte(`[id]`)), nil)
		t.Assert(gconv.MapsDeep([]byte(`test`)), nil)
		t.Assert(gconv.MapsDeep([]string{}), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		stringInterfaceMapList := make([]map[string]interface{}, 0)
		stringInterfaceMapList = append(stringInterfaceMapList, map[string]interface{}{"id": 100})
		stringInterfaceMapList = append(stringInterfaceMapList, map[string]interface{}{"id": 200})
		list := gconv.MapsDeep(stringInterfaceMapList)
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)

		list = gconv.MapsDeep([]byte(jsonStr))
		t.Assert(len(list), 2)
		t.Assert(list[0]["id"], 100)
		t.Assert(list[1]["id"], 200)
	})
}

func TestMapWithJsonOmitEmpty(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S struct {
			Key   string      `json:",omitempty"`
			Value interface{} `json:",omitempty"`
		}
		s := S{
			Key:   "",
			Value: 1,
		}
		t.Assert(gconv.Map(s), g.Map{"Value": 1})
	})
}
