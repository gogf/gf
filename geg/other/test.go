package main

import (
    "fmt"
    "gitee.com/johng/gf/g/os/gfpool"
    "os"
    "time"
)

func main() {
    for {
        f, err := gfpool.Open("/home/john/temp/log", os.O_RDWR, 0666)
        fmt.Println(err)
        _, err = f.WriteString("123")
        fmt.Println(err)
        //f.Close()
        time.Sleep(time.Second)
    }
}