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

// 设置数据
func testSet() {
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
        c, _ := j.ToJson()
        fmt.Println(string(c))
    }
}

// 将Json数据转换为其他数据格式
func testConvert() {
    data :=
        `{
            "users" : {
                "count" : 100,
                "list"  : ["John", "小明"]
            }
        }`
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        c, _ := j.ToJson()
        fmt.Println("JSON:")
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("XML:")
        c, _ = j.ToXmlIndent()
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("YAML:")
        c, _ = j.ToYaml()
        fmt.Println(string(c))
        fmt.Println("======================")

        fmt.Println("TOML:")
        c, _ = j.ToToml()
        fmt.Println(string(c))
    }
}

func main() {
    testSet()
}