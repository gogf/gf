package main

import (
    "gitee.com/johng/gf/g/database/gmq"
    "fmt"
)

func main() {
    mq, err := gmq.New("/tmp/gmq")
    if err != nil {
        fmt.Println(err)
    }

    //t1 := gtime.Microsecond()
    //for i := 0; i < 10; i++ {
    //    mq.Group("test").Push([]byte("gmq_message_" + strconv.Itoa(i)))
    //}
    //fmt.Println("push cost:", gtime.Microsecond() - t1)
    fmt.Println(string(mq.Group("test").Pop()))
    fmt.Println("length", mq.Group("test").Length())
    //fmt.Println(mq.Group("test").Add([]byte("gmq_message")))
    //fmt.Println(mq.Group("test").Remove(1000))
}