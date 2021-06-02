package main

import (
	"context"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/i18n/gi18n"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	type User struct {
		Name    string `v:"required#ReuiredUserName"`
		Type    int    `v:"required#ReuiredUserType"`
		Project string `v:"size:10#MustSize"`
	}
	var (
		data = g.Map{
			"name":    "john",
			"project": "gf",
		}
		user  = User{}
		ctxEn = gi18n.WithLanguage(context.TODO(), "en")
		ctxCh = gi18n.WithLanguage(context.TODO(), "zh-CN")
	)

	if err := gconv.Scan(data, &user); err != nil {
		panic(err)
	}
	// 英文
	if err := g.Validator().Ctx(ctxEn).Data(data).CheckStruct(user); err != nil {
		g.Dump(err.String())
	}
	// 中文
	if err := g.Validator().Ctx(ctxCh).Data(data).CheckStruct(user); err != nil {
		g.Dump(err.String())
	}
}
