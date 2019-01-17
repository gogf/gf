package main

import (
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g"
    "fmt"
)

func main() {
    type Score struct {
        Name   string
        Result int
    }
    type User1 struct {
        Scores Score
    }
    type User2 struct {
        Scores *Score
    }

    user1  := new(User1)
    user2  := new(User2)
    scores := map[string]interface{}{
        "Scores" : map[string]interface{}{
            "Name"   : "john",
            "Result" : 100,
        },
    }

    if err := gconv.Struct(scores, user1); err != nil {
        fmt.Println(err)
    } else {
        g.Dump(user1)
    }
    if err := gconv.Struct(scores, user2); err != nil {
        fmt.Println(err)
    } else {
        g.Dump(user2)
    }
}