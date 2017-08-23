package main

import (
    "fmt"
    "g/encoding/gjson"
    "gapp/gluster/gluster"
)



func main() {
    var st gluster.ServiceStruct
    err := gjson.DecodeTo("{\"names\":{\"1\":1}}", &st)
    fmt.Println(err)
    fmt.Println(st)
}