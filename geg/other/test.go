package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
    "gitee.com/johng/gf/g/os/gfile"
)

func main() {

    gtime.SetInterval(10*time.Millisecond, func() bool {
        path := "./temp.txt"
        gfile.PutBinContentsAppend(path, []byte("1"))
        fmt.Println(gfile.MTimeMillisecond(path))
        return true
    })

    time.Sleep(time.Hour)

}
