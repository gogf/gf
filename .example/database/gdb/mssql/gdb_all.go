package main

import (
	"fmt"
	"github.com/gogf/gf/v2/os/gtime"

	//_ "github.com/denisenkom/go-mssqldb"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	type Table2 struct {
		Id         string     `orm:"id;pr" json:"id"`              //ID
		CreateTime gtime.Time `orm:"createtime" json:"createtime"` //创建时间
		UpdateTime gtime.Time `orm:"updatetime" json:"updatetime"` //更新时间
	}
	var table2 Table2
	err := g.DB().Model("table2").Where("id", 1).Scan(&table2)
	if err != nil {
		panic(err)
	}
	fmt.Println(table2.CreateTime)
}
