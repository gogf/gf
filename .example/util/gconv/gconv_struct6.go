package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

func main() {
	type Score struct {
		Name   string
		Result int
	}
	type User struct {
		Scores []*Score
	}

	user := new(User)
	scores := g.Map{
		"Scores": g.Slice{
			g.Map{
				"Name":   "john",
				"Result": 100,
			},
			g.Map{
				"Name":   "smith",
				"Result": 60,
			},
		},
	}

	// 嵌套struct转换，属性为slice类型，数值为slice map类型
	if err := gconv.Struct(scores, user); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user)
	}
}
