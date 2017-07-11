package main

import (
    "fmt"
    "g/encoding/gjson"
)


func main() {
    //json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"items":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
    json := `[{"CityId":18,"CityName":"西安","ProvinceId":27,"CityOrder":1},{"CityId":53,"CityName":"广州","ProvinceId":27,"CityOrder":1}]`
    //json := `{"name"  :  "中国",  "age" : 31, "items":[1,2,3]}`
    //json := `[["a","b","c"],["d","e","f"]]`
    //json := `["a","b","c"]`
    //json   := `
    //[1,{"a":2},
    //{"a":{}},
    //{"a":[]},
    //{"a":[{}]},
    //{"{[a" : "\"2,:3," a ":33}]"}]` // 错误的json
    //json := `["a","b","c"`        // 错误的json
    //json := `,{ "name"  :  "中国",  "age" : 31, "items":[1,2]:}` //错误的json
    p, err := gjson.Decode(&json)
    if err == nil {
        //p.Print()
        //fmt.Println(p.Get("0"))
        fmt.Println(p.GetMap("0"))
    } else {
        fmt.Println(err)
    }
    //fmt.Println()
    //fmt.Println()
    ////v := make(map[string]interface{})
    ////i := 31
    ////j := "john"
    ////v["age"]  = i
    ////v["name"] = make(map[string]interface{})
    ////t := v["name"]
    ////t.(map[string]interface{})["n"] = j
    ////
    ////fmt.Println(v)
    //var s struct{
    //    v interface{}
    //    p interface{}
    //}
    //v  := make(map[string]interface{})
    //s.v = v
    //s.p = &v
    //c  := (*s.p.(*map[string]interface{}))
    //c["name1"] = "john1"
    //
    //t          := make(map[string]interface{})
    //c["/"]      = t
    //s.p         = &t
    //t["name2"]  = "john2"
    //
    //c2         := (*s.p.(*map[string]interface{}))
    //c2["name3"] = "john3"
    //
    ////t2[2] = 100
    //fmt.Println(s.v)


    //a := map[string]interface{} {
    //    "name" : "john",
    //    "list" : []interface{}{
    //        1,2,3, "fuck",
    //    },
    //    "item" : map[string]string {
    //        "n1" : "v1",
    //        "n2" : "v2",
    //        "n3" : "v3",
    //    },
    //}
    //fmt.Println(json.M)

    //
    //var a = []int{1,2,3}
    //var b = []int{4,5,6, 7,8}
    //cc := make([]int, len(a) + 12)
    //a = cc
    //copy(a, b)
    //fmt.Println(a)
    //fmt.Println(b)
}