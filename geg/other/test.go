package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    data :=
        `{
            "users" : {
                "count" : 100
            }
        }`
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        j.Set("users.count",  1)
        j.Set("users.list",  []string{"John", "小明"})
        j.Set("users.0",     "a")
        fmt.Println(j.Get("users.count"))
        fmt.Println(j.Get("users.count"))
        //fmt.Println(j.Get("users.list"))
        c, _ := j.ToJson()
        fmt.Println(string(c))
    }
}