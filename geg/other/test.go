package main

import (
	"fmt"

	"github.com/gogf/gf/g"

	"github.com/gogf/gf/g/util/gconv"

	"github.com/gogf/gf/g/encoding/gparser"
)

func main() {
	type User struct {
		Uid      int
		Name     string
		SiteUrl  string `gconv:"-"`
		NickName string `gconv:"nickname, omitempty"`
		Pass1    string `gconv:"password1"`
		Pass2    string `gconv:"password2"`
	}

	g.Dump(gconv.Map(User{
		Uid:     100,
		Name:    "john",
		SiteUrl: "https://goframe.org",
		Pass1:   "123",
		Pass2:   "456",
	}))

	s, err := gparser.VarToJsonString(User{
		Uid:     100,
		Name:    "john",
		SiteUrl: "https://goframe.org",
		Pass1:   "123",
		Pass2:   "456",
	})
	fmt.Println(err)
	fmt.Println(s)
}
