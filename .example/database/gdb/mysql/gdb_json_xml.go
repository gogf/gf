package main

import (
	"fmt"
	"github.com/gogf/gf/os/gctx"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/encoding/gparser"
	"github.com/gogf/gf/frame/g"
)

func main() {
	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "3306",
		User:    "root",
		Pass:    "12345678",
		Name:    "test",
		Type:    "mysql",
		Role:    "master",
		Charset: "utf8",
	})
	var (
		db  = g.DB()
		ctx = gctx.New()
	)
	one, err := db.Ctx(ctx).Model("user").Where("id=?", 1).One()
	if err != nil {
		panic(err)
	}

	// 使用内置方法转换为json/xml
	fmt.Println(one.Json())
	fmt.Println(one.Xml())

	// 自定义方法方法转换为json/xml
	jsonContent, _ := gparser.VarToJson(one.Map())
	fmt.Println(string(jsonContent))
	xmlContent, _ := gparser.VarToXml(one.Map())
	fmt.Println(string(xmlContent))
}
