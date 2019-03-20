package main

import (
    "github.com/gogf/gf/g"
    "github.com/gogf/gf/g/text/gregex"
)

func main() {
    s := "127.0.0.1:6379,1,nhytaf176tg?maxIdle=1&maxActive=0&idleTimeout=60&maxConnLifetime=60"
    array, err := gregex.MatchString(`(.+):(\d+),{0,1}(\d*),{0,1}(.*)\?(.+)`, s)
    g.Dump(err)
    g.Dump(array)
}