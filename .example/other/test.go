package main

import (
	"fmt"

	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
)

type Item struct {
	GroupId    int
	Interval   string
	MetricName string
	Url        string
}

func main() {
	j, err := gjson.Load("config.toml")
	if err != nil {
		panic(err)
	}
	m := j.GetMap("active-pulling")
	//g.Dump(m)

	newm := make(map[string][]Item)
	err = gconv.MapStructs(m, &newm)
	fmt.Println(err)
	g.Dump(newm)

}
