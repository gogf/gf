package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
)

// 使用默认映射规则绑定属性值到对象
func main() {
	type User struct {
		Uid     int
		Name    string
		SiteUrl string
		Pass1   string
		Pass2   string
	}
	user := new(User)
	params := g.Map{
		"uid":      1,
		"Name":     "john",
		"site_url": "https://goframe.org",
		"PASS1":    "123",
		"PASS2":    "456",
	}
	if err := gconv.Struct(params, user); err == nil {
		fmt.Println(user)
	}
}
