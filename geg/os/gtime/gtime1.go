package main

import (
    "regexp"
    "fmt"
)

func main() {
    timeRegex, err := regexp.Compile(`(\d{4}-\d{2}-\d{2})[\sT]{0,1}(\d{2}:\d{2}:\d{2}){0,1}\.{0,1}(\d{0,9})([\sZ]{0,1})([\+-]{0,1})([:\d]*)`)
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
        "2018-02-09",
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
