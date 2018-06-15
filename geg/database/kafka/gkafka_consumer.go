package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
)

func main () {
    client := gkafka.New(gkafka.Config{
        GroupId : "group_1",
        Servers : "localhost:9092",
        Topics  : "abc",
    })
    defer client.Close()
    for {
       msg, err := client.Receive()
       fmt.Println(err)
       fmt.Println(string(msg.Value))
    }
}
