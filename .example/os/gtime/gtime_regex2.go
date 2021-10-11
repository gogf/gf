package main

import (
	"fmt"
	"regexp"

	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	timeRegex, err := regexp.Compile(gtime.TIME_REAGEX_PATTERN2)
	if err != nil {
		panic(err)
	}
	array := []string{
		"01-Nov-2018 11:50:28 +0805 LMT",
		"01-Nov-2018T15:04:05Z07:00",
		"01-Nov-2018T01:19:15+08:00",
		"01-Nov-2018 11:50:28 +0805 LMT",
		"01/Nov/18 11:50:28",
		"01/Nov/2018 11:50:28",
		"01/Nov/2018:11:50:28",
		"01/Nov/2018",
	}
	for _, s := range array {
		fmt.Println(s)
		match := timeRegex.FindStringSubmatch(s)
		for k, v := range match {
			fmt.Println(k, v)
		}
		fmt.Println()
	}
}
