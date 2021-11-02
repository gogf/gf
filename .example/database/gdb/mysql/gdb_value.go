package main

import (
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
)

func main() {
	one, err := g.Model("carlist c").
		LeftJoin("cardetail d", "c.postid=d.carid").
		Where("c.postid", "142039140032006").
		Fields("c.*,d.*").One()
	fmt.Println(err)
	g.Dump(one)
}
