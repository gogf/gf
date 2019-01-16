package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/ghttp"
    "sync"
)

func main() {
    g := sync.WaitGroup{}
    c := ghttp.NewClient()
    for s1 := 97; s1 <= 122; s1++ {
        g.Add(1)
        go func(s1 int) {
            for s2 := 97; s2 <= 122; s2++ {
                for s3 := 97; s3 <= 122; s3++ {
                    for s4 := 97; s4 <= 122; s4++ {
                        url := "https://github.com/" + string(s1) + string(s2) + string(s3) + string(s4)
                        if r, _ := c.Get(url); r != nil {
                            if r.StatusCode != 200 {
                                fmt.Println(url, r.StatusCode)
                                r.Close()
                            }
                        }
                    }
                }
            }
        }(s1)
    }
    g.Wait()
}