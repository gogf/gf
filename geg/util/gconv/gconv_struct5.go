package main

import (
	"fmt"
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/util/gconv"
)

func main() {
	type Score struct {
		Name   string
		Result int
	}
	type User struct {
		Scores []Score
	}

	user := new(User)
	scores := map[string]interface{}{
		"Scores": map[string]interface{}{
			"Name":   "john",
			"Result": 100,
		},
	}

	// 嵌套struct转换，属性为slice类型，数值为map类型
	if err := gconv.Struct(scores, user); err != nil {
		fmt.Println(err)
	} else {
		g.Dump(user)
	}
}
