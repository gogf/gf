package main

import (
	"github.com/gogf/gf/g"
)

func main() {
	db := g.DB()
	conditions := g.Map{
		"nickname like ?":    "%T%",
		"id between ? and ?": g.Slice{1, 3},
		"id >= ?":            1,
		"create_time > ?":    0,
		"id in(?)":           g.Slice{1, 2, 3},
	}
	db.Table("user").Where(conditions).OrderBy("id asc").All()
}
