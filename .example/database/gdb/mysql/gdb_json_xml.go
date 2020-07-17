package main

import (
	"fmt"

	"github.com/jin502437344/gf/database/gdb"
	"github.com/jin502437344/gf/encoding/gparser"
	"github.com/jin502437344/gf/frame/g"
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
	db := g.DB()
	one, err := db.Table("user").Where("id=?", 1).One()
	if err != nil {
		panic(err)
	}

	// 使用内置方法转换为json/xml
	fmt.Println(one.ToJson())
	fmt.Println(one.ToXml())

	// 自定义方法方法转换为json/xml
	jsonContent, _ := gparser.VarToJson(one.ToMap())
	fmt.Println(string(jsonContent))
	xmlContent, _ := gparser.VarToXml(one.ToMap())
	fmt.Println(string(xmlContent))
}
