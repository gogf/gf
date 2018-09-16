package main

import (
    "gitee.com/johng/gf/g/util/gutil"
)

func main() {
    gutil.Dump(map[interface{}]interface{} {
        1 : "john",
    })
}
