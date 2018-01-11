package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/encoding/gjson"
)

func main() {
    data := `{
        "users" : {
                "count" : 100,
                "list"  : [
                    {"name" : "小明",  "score" : 60},
                    {"name" : "John", "score" : 99.5}
                ]
            }
    }`
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        fmt.Println("John Score:", j.GetFloat32("users.list.1.score"))
    }
}