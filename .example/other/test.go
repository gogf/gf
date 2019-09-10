package main

import (
	"encoding/json"
	"strings"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/text/gregex"

	"github.com/gogf/gf/container/gset"

	"github.com/gogf/gf/os/gfile"
)

func main() {
	path1 := "/Users/john/Temp/downloaded_data_parsed.txt"
	path2 := "/Users/john/Temp/downloaded_data_parsed_mapping.txt"
	array := strings.Split(gfile.GetContents(path1), "\n")
	mapping := make(map[string]*gset.Set)
	for _, line := range array {
		array, _ := gregex.MatchString(`add group success \[\[(\d+),\[*(\d+)\]*`, line)
		if len(array) != 3 {
			g.Dump(line)
			g.Dump(array)
			continue
		}
		if _, ok := mapping[array[1]]; !ok {
			mapping[array[1]] = gset.New()
		}
		mapping[array[1]].Add(gconv.Interfaces(strings.Split(array[2], ","))...)
	}
	g.Dump(mapping)
	b, _ := json.Marshal(mapping)
	gfile.PutBytes(path2, b)
}
