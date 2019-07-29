package main

import (
	"fmt"

	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/encoding/gjson"
	"github.com/gogf/gf/g/util/gconv"
)

type Item struct {
	GroupId    int
	Interval   string
	MetricName string
	Url        string
}

func main() {
<<<<<<< HEAD
	latestVersion := g.NewVar(nil, true)
	fmt.Println(latestVersion.IsNil())
=======
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
>>>>>>> master
}
