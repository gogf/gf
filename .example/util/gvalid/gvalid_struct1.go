package main

import (
	"context"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gvalid"
)

type User struct {
	Uid   int    `gvalid:"uid      @integer|min:1#用户UID不能为空"`
	Name  string `gvalid:"name     @required|length:6,30#请输入用户名称|用户名称长度非法"`
	Pass1 string `gvalid:"password1@required|password3"`
	Pass2 string `gvalid:"password2@required|password3|same:password1#||两次密码不一致，请重新输入"`
}

func main() {
	user := &User{
		Name:  "john",
		Pass1: "Abc123!@#",
		Pass2: "123",
	}

	// 使用结构体定义的校验规则和错误提示进行校验
	g.Dump(gvalid.CheckStruct(context.TODO(), user, nil).Map())

	// 自定义校验规则和错误提示，对定义的特定校验规则和错误提示进行覆盖
	rules := map[string]string{
		"Uid": "required",
	}
	msgs := map[string]interface{}{
		"Pass2": map[string]string{
			"password3": "名称不能为空",
		},
	}
	g.Dump(gvalid.CheckStruct(context.TODO(), user, rules, msgs).Map())
}
