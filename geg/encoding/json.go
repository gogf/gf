package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/g/encoding/gjson"
)

func getByPattern() {
    data :=
        `{
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

// 当键名存在"."号时，检索优先级：键名->层级，因此不会引起歧义
func testMultiDots() {
    data :=
        `{
            "users" : {
                "count" : 100
            },
            "users.count" : 101
        }`
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        fmt.Println("Users Count:", j.GetInt("users.count"))
    }
}

func main() {
    testMultiDots()
}