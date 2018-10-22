package main

import (
    "fmt"
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gcron"
    "time"
)

func main() {
    gcron.Add("0 30 * * * *", func() { fmt.Println("Every hour on the half hour") })
    gcron.Add("* * * * * *",  func() { fmt.Println("Every second") })
    gcron.Add("@hourly",      func() { fmt.Println("Every hour") })
    gcron.Add("@every 1h30m", func() { fmt.Println("Every hour thirty") })
    g.Dump(gcron.Entries())
    time.Sleep(3*time.Second)
}