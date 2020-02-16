package main

import (
	"fmt"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	// select * from table1 where field1='1111'
	// and (
	//   (field2='2' and field3='3')
	//   or
	//   (field2='21' and field3='31')
	//   or
	//   (field2='22' and field3='32')
	// )
	//g.DB().Table("table1").
	//	Where("field1", "1111").
	//	And(g.Map{"field2": 2, "field3": 3})
	//fmt.Println(gconv.GTime("2020-01-01 12:01:00").String())
	//t := gconv.Convert("2020-01-01 12:01:00", "gtime.Time").(gtime.Time)
	t := gconv.Convert(1989, "gtime.Time").(gtime.Time)
	fmt.Println(t.String())
	fmt.Println(gconv.String(t))
}
