package main

import (
    "fmt"
    "sync"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/net/ghttp"
    "gitee.com/johng/gf/g/container/gtype"
)

func main() {
    clientMax  := 10
    requestMax := 1000
    failureNum := gtype.NewInt64()
    successNum := gtype.NewInt64()
    startTime  := gtime.Millisecond()

    wg := sync.WaitGroup{}
    for i := 0; i < clientMax; i++ {
        wg.Add(1)
        go func(clientId int) {
            for j := 0; j < requestMax; j++ {
                if c, e := ghttp.Get("http://127.0.0.1/"); e == nil {
                    c.Close()
                    successNum.Add(1)
                } else {
                    failureNum.Add(1)
                }
            }
            wg.Done()
        }(i)
    }
    wg.Wait()

    fmt.Printf("time spent: %d ms, success:%d, failure: %d\n",
        gtime.Millisecond() - startTime, successNum.Val(), failureNum.Val())
}