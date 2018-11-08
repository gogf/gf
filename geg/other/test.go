package main

import (
    "fmt"
    "os"
    "time"
)

func main() {
    for {
        stat, err := os.Stat("/home/john/temp/log")
        fmt.Println(err)
        fmt.Println(stat.Size())
        time.Sleep(time.Second)
    }
}