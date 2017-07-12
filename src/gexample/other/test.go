package main

import (
    "fmt"
    "g/core/gfunc"
    "g/database/gdb"
)


func main() {
    //s := "中国"
    //p := &s
    //m := map[string]string {
    //    "name" : string([]rune(*p)),
    //}
    //fmt.Println(rand.Int()%15)
    fmt.Println(gfunc.Rand(10, 100))
    gdb.Test()
    redis
    //fmt.Println(len(make(map[string]string)))
    //json := `{"name":"中国","age":31,"list":[["a","b","c"],["d","e","f"]],"item":{"title":"make\"he moon","name":"make'he moon","content":"'[}]{[}he moon"}}`
    //json := `{"name"  :  "中国",  "age" : 31, "items":[1,2,3]}`
    //json := `[["a","b","c"],["d","e","f"]]`
    //json := `["a","b","c"]`
    //jsonDecode(&json)
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


    //type JsonMap   map[string]interface{}
    //type JsonArray []interface{}
    //
    //a := JsonMap {
    //    "name" : "john",
    //    "list" : JsonArray {
    //        1,2,3, "fuck",
    //    },
    //    "items" : JsonMap {
    //        "n1" : "v1",
    //        "n2" : "v2",
    //        "n3" : "v3",
    //    },
    //}
    //fmt.Println(a["items"].(JsonMap)["n1"])

    //
    //var a = []int{1,2,3}
    //var b = []int{4,5,6, 7,8}
    //cc := make([]int, len(a) + 12)
    //a = cc
    //copy(a, b)
    //fmt.Println(a)
    //fmt.Println(b)
}