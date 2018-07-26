package main

import (
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {
    gfile.PutContentsAppend("/tmp/test", "1")
    gfile.PutContentsAppend("/tmp/test", "2")
    gfile.PutContentsAppend("/tmp/test", "3")
}