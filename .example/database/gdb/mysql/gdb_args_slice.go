package main

import (
	"github.com/jin502437344/gf/frame/g"
)

func main() {
	db := g.DB()

	db.Table("user").Where("nickname like ? and passport like ?", g.Slice{"T3", "t3"}).OrderBy("id asc").All()

	conditions := g.Map{
		"nickname like ?":    "%T%",
		"id between ? and ?": g.Slice{1, 3},
		"id >= ?":            1,
		"create_time > ?":    0,
		"id in(?)":           g.Slice{1, 2, 3},
	}
	db.Table("user").Where(conditions).OrderBy("id asc").All()

	var params []interface{}
	db.Table("user").Where("1=1", params).OrderBy("id asc").All()
}
