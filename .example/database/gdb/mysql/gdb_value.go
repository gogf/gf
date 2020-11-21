package main

import (
	"github.com/gogf/gf/frame/g"
	"time"
)

func main() {
	db := g.DB()
	db.SetDebug(true)

	t1, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:32")
	t2, _ := time.Parse("2006-01-02 15:04:05", "2020-10-27 19:03:34")
	u, err := g.DB().Table("orders").Where("updated_at>? and updated_at<?", t1, t2).One()
	//u, err := g.DB().Table("orders").Where("updated_at>? and updated_at<?", gtime.New("2020-10-27 19:03:32"), gtime.New("2020-10-27 19:03:34")).One()
	//u, err := g.DB().Table("orders").Fields("id").Where("updated_at>'2020-10-27 19:03:32' and updated_at<'2020-10-27 19:03:34'").Value()
	g.Dump(u, err)
}
