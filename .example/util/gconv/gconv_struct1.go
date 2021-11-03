package main

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type User struct {
	Uid      int
	Name     string
	Site_Url string
	NickName string
	Pass1    string `gconv:"password1"`
	Pass2    string `gconv:"password2"`
}

func main() {
	user := (*User)(nil)

	// 使用默认映射规则绑定属性值到对象
	user = new(User)
	params1 := g.Map{
		"uid":       1,
		"Name":      "john",
		"siteurl":   "https://goframe.org",
		"nick_name": "johng",
		"PASS1":     "123",
		"PASS2":     "456",
	}
	if err := gconv.Struct(params1, user); err == nil {
		g.Dump(user)
	}

	// 使用struct tag映射绑定属性值到对象
	user = new(User)
	params2 := g.Map{
		"uid":       2,
		"name":      "smith",
		"site-url":  "https://goframe.org",
		"nick name": "johng",
		"password1": "111",
		"password2": "222",
	}
	if err := gconv.Struct(params2, user); err == nil {
		g.Dump(user)
	}
}
