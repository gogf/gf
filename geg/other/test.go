package main

import (
    "gitee.com/johng/gf/g/container/garray"
    "fmt"
)

type S struct {

}



func main() {
    var source = []string{"59705a2c1fd50736a4c768a1", "597a95ff1fd5073e48bb2272", "597a960f1fd5073e48bb2274"}
    var CacheChannelKeys = garray.NewSortedStringArray(3)
    CacheChannelKeys.Add(source...)
    fmt.Println(CacheChannelKeys.Len())
}

