package main

import (
	"github.com/gogf/gf/frame/g"
)

func main() {
	// error!
	r, err := g.DB().Model("user").Where(g.Map{
		"or": g.Map{
			"nickname":       "jim",
			"create_time > ": "2019-10-01",
		},
		"and": g.Map{
			"nickname":       "tom",
			"create_time > ": "2019-10-01",
		},
	}).All()
	if err != nil {
		panic(err)
	}
	g.Dump(r)

}
