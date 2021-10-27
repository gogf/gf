package main

import (
	"fmt"
	"regexp"

	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	timeRegex, err := regexp.Compile(gtime.TIME_REAGEX_PATTERN1)
	if err != nil {
		panic(err)
	}
	array := []string{
		"2017-12-14 04:51:34 +0805 LMT",
		"2006-01-02T15:04:05Z07:00",
		"2014-01-17T01:19:15+08:00",
		"2018-02-09T20:46:17.897Z",
		"2018-02-09 20:46:17.897",
		"2018-02-09T20:46:17Z",
		"2018-02-09 20:46:17",
		"2018/10/31 - 16:38:46",
		"2018-02-09",
		"2017/12/14 04:51:34 +0805 LMT",
		"2018/02/09 12:00:15",
		"18/02/09 12:16",
		"18/02/09 12",
		"18/02/09 +0805 LMT",
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
