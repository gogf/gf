package main

import (
	"fmt"
	"github.com/jin502437344/gf/os/gtime"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	type Table2 struct {
		Id         string     `orm:"id;pr" json:"id"`              //ID
		Createtime gtime.Time `orm:"createtime" json:"createtime"` //创建时间
		Updatetime gtime.Time `orm:"updatetime" json:"updatetime"` //更新时间
	}
	var table2 Table2
	err := g.DB().Table("table2").Where("id=?", 1).Struct(&table2)
	if err != nil {
		panic(err)
	}
	fmt.Println(table2.Createtime)
}
