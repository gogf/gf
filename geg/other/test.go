package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregex"
)

func main() {
    fmt.Println(gregex.IsMatchString("g/os/glog/glog.+$", "g/os/glog/glog_logger.go"))
}

