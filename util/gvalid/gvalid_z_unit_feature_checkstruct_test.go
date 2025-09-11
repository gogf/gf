// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gconv"
)

func Test_CheckStruct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"@required|length:6,16",
			"@between:18,30",
		}
		msgs := map[string]any{
			"Name": map[string]string{
				"required": "名称不能为空",
				"length":   "名称长度为{min}到{max}个字符",
			},
			"Age": "年龄为18到30周岁",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"Name@required|length:6,16#名称不能为空",
			"Age@between:18,30",
		}
		msgs := map[string]any{
			"Name": map[string]string{
				"required": "名称不能为空",
				"length":   "名称长度为{min}到{max}个字符",
			},
			"Age": "年龄为18到30周岁",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名称长度为6到16个字符")
		t.Assert(err.Maps()["Age"]["between"], "年龄为18到30周岁")
	})

	gtest.C(t, func(t *gtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := []string{
			"Name@required|length:6,16#名称不能为空|",
			"Age@between:18,30",
		}
		msgs := map[string]any{
			"Name": map[string]string{
				"required": "名称不能为空",
				"length":   "名称长度为{min}到{max}个字符",
			},
			"Age": "年龄为18到30周岁",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名称长度为6到16个字符")
		t.Assert(err.Maps()["Age"]["between"], "年龄为18到30周岁")
	})

	gtest.C(t, func(t *gtest.T) {
		type Object struct {
			Name string
			Age  int
		}
		rules := map[string]string{
			"Name": "required|length:6,16",
			"Age":  "between:18,30",
		}
		msgs := map[string]any{
			"Name": map[string]string{
				"required": "名称不能为空",
				"length":   "名称长度为{min}到{max}个字符",
			},
			"Age": "年龄为18到30周岁",
		}
		obj := &Object{"john", 16}
		err := g.Validator().Data(obj).Rules(rules).Messages(msgs).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["Name"]["required"], "")
		t.Assert(err.Maps()["Name"]["length"], "名称长度为6到16个字符")
		t.Assert(err.Maps()["Age"]["between"], "年龄为18到30周岁")
	})

	gtest.C(t, func(t *gtest.T) {
		type LoginRequest struct {
			Username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 2)
		t.Assert(err.Maps()["username"]["required"], "用户名不能为空")
		t.Assert(err.Maps()["password"]["required"], "登录密码不能为空")
	})

	gtest.C(t, func(t *gtest.T) {
		type LoginRequest struct {
			Username string `json:"username" valid:"@required#用户名不能为空"`
			Password string `json:"password" valid:"@required#登录密码不能为空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNil(err)
	})

	gtest.C(t, func(t *gtest.T) {
		type LoginRequest struct {
			username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		var login LoginRequest
		err := g.Validator().Data(login).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["password"]["required"], "登录密码不能为空")
	})

	// gvalid tag
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能为空"`
			Age      int    `valid:"age@required#年龄不能为空"`
			Username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能为空")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能为空"`
			Age      int    `valid:"age@required#年龄不能为空"`
			Username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}

		rules := []string{
			"username@required#用户名不能为空",
		}

		err := g.Validator().Data(user).Rules(rules).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能为空")
	})

	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#ID不能为空"`
			Age      int    `valid:"age@required#年龄不能为空"`
			Username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
	})

	// valid tag
	gtest.C(t, func(t *gtest.T) {
		type User struct {
			Id       int    `valid:"uid@required|min:10#|ID不能为空"`
			Age      int    `valid:"age@required#年龄不能为空"`
			Username string `json:"username" valid:"username@required#用户名不能为空"`
			Password string `json:"password" valid:"password@required#登录密码不能为空"`
		}
		user := &User{
			Id:       1,
			Username: "john",
			Password: "123456",
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(len(err.Maps()), 1)
		t.Assert(err.Maps()["uid"]["min"], "ID不能为空")
	})
}

func Test_CheckStruct_EmbeddedObject_Attribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Time *gtime.Time
		}
		type Object struct {
			Base
			Name string
			Type int
		}
		rules := map[string]string{
			"Name": "required",
			"Type": "required",
		}
		ruleMsg := map[string]any{
			"Name": "名称必填",
			"Type": "类型必填",
		}
		obj := &Object{}
		obj.Type = 1
		obj.Name = "john"
		obj.Time = gtime.Now()
		err := g.Validator().Data(obj).Rules(rules).Messages(ruleMsg).Run(context.TODO())
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Base struct {
			Name string
			Type int
		}
		type Object struct {
			Base Base
			Name string
			Type int
		}
		rules := map[string]string{
			"Name": "required",
			"Type": "required",
		}
		ruleMsg := map[string]any{
			"Name": "名称必填",
			"Type": "类型必填",
		}
		obj := &Object{}
		obj.Type = 1
		obj.Name = "john"
		err := g.Validator().Data(obj).Rules(rules).Messages(ruleMsg).Run(context.TODO())
		t.AssertNil(err)
	})
}

func Test_CheckStruct_With_EmbeddedObject(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `valid:"password1@required|same:password2#请输入您的密码|您两次输入的密码不一致"`
			Pass2 string `valid:"password2@required|same:password1#请再次输入您的密码|您两次输入的密码不一致"`
		}
		type User struct {
			Id   int
			Name string `valid:"name@required#请输入您的姓名"`
			Pass
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["name"], g.Map{"required": "请输入您的姓名"})
		t.Assert(err.Maps()["password1"], g.Map{"same": "您两次输入的密码不一致"})
		t.Assert(err.Maps()["password2"], g.Map{"same": "您两次输入的密码不一致"})
	})
}

func Test_CheckStruct_With_StructAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Pass struct {
			Pass1 string `valid:"password1@required|same:password2#请输入您的密码|您两次输入的密码不一致"`
			Pass2 string `valid:"password2@required|same:password1#请再次输入您的密码|您两次输入的密码不一致"`
		}
		type User struct {
			Pass
			Id   int
			Name string `valid:"name@required#请输入您的姓名"`
		}
		user := &User{
			Name: "",
			Pass: Pass{
				Pass1: "1",
				Pass2: "2",
			},
		}
		err := g.Validator().Data(user).Run(context.TODO())
		t.AssertNE(err, nil)
		t.Assert(err.Maps()["name"], g.Map{"required": "请输入您的姓名"})
		t.Assert(err.Maps()["password1"], g.Map{"same": "您两次输入的密码不一致"})
		t.Assert(err.Maps()["password2"], g.Map{"same": "您两次输入的密码不一致"})
	})
}

func Test_CheckStruct_Optional(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page      int    `v:"required|min:1         # page is required"`
			Size      int    `v:"required|between:1,100 # size is required"`
			ProjectId string `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page      int       `v:"required|min:1         # page is required"`
			Size      int       `v:"required|between:1,100 # size is required"`
			ProjectId *gvar.Var `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page      int `v:"required|min:1         # page is required"`
			Size      int `v:"required|between:1,100 # size is required"`
			ProjectId int `v:"between:1,10000        # project id must between {min}, {max}"`
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.Assert(err.String(), "project id must between 1, 10000")
	})
}

func Test_CheckStruct_NoTag(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Page      int
			Size      int
			ProjectId string
		}
		obj := &Params{
			Page: 1,
			Size: 10,
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNil(err)
	})
}

func Test_CheckStruct_InvalidRule(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Name  string
			Age   uint
			Phone string `v:"mobile"`
		}
		obj := &Params{
			Name:  "john",
			Age:   18,
			Phone: "123",
		}
		err := g.Validator().Data(obj).Run(context.TODO())
		t.AssertNE(err, nil)
	})
}

func TestValidator_CheckStructWithData(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid      int64  `v:"required"`
			Nickname string `v:"required-with:uid"`
		}
		data := UserApiSearch{
			Uid:      1,
			Nickname: "john",
		}
		t.Assert(
			g.Validator().Data(data).Assoc(
				g.Map{"uid": 1, "nickname": "john"},
			).Run(context.TODO()),
			nil,
		)
	})
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid      int64  `v:"required"`
			Nickname string `v:"required-with:uid"`
		}
		data := UserApiSearch{}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid      int64  `json:"uid" v:"required"`
			Nickname string `json:"nickname" v:"required-with:Uid"`
		}
		data := UserApiSearch{
			Uid: 1,
		}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})

	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid       int64       `json:"uid"`
			Nickname  string      `json:"nickname" v:"required-with:Uid"`
			StartTime *gtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: nil,
			EndTime:   nil,
		}
		t.Assert(g.Validator().Data(data).Assoc(g.Map{}).Run(context.TODO()), nil)
	})
	gtest.C(t, func(t *gtest.T) {
		type UserApiSearch struct {
			Uid       int64       `json:"uid"`
			Nickname  string      `json:"nickname" v:"required-with:Uid"`
			StartTime *gtime.Time `json:"start_time" v:"required-with:EndTime"`
			EndTime   *gtime.Time `json:"end_time" v:"required-with:StartTime"`
		}
		data := UserApiSearch{
			StartTime: gtime.Now(),
			EndTime:   nil,
		}
		t.AssertNE(g.Validator().Data(data).Assoc(g.Map{"start_time": gtime.Now()}).Run(context.TODO()), nil)
	})
}

func Test_CheckStruct_PointerAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			Name string
			Age  *uint `v:"min:18"`
		}
		req := &Req{
			Name: "john",
			Age:  gconv.PtrUint(0),
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Age value `0` must be equal or greater than 18")
	})
	gtest.C(t, func(t *gtest.T) {
		type Req struct {
			Name string `v:"min-length:3"`
			Age  *uint  `v:"min:18"`
		}
		req := &Req{
			Name: "j",
			Age:  gconv.PtrUint(19),
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Name value `j` length must be equal or greater than 3")
	})
	gtest.C(t, func(t *gtest.T) {
		type Params struct {
			Age *uint `v:"min:18"`
		}
		type Req struct {
			Name   string
			Params *Params
		}
		req := &Req{
			Name: "john",
			Params: &Params{
				Age: gconv.PtrUint(0),
			},
		}
		err := g.Validator().Data(req).Run(context.TODO())
		t.Assert(err.String(), "The Age value `0` must be equal or greater than 18")
	})
}
