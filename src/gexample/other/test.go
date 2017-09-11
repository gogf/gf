package main

import (
    "fmt"
    "g/net/ghttp"
)

type ST struct {
    I int64
}


func main() {
    for i := 0; i< 1000; i++ {
        go func() {
            for i := 0; i< 100; i ++ {
                r := ghttp.Post("http://127.0.0.1:4168/kv", fmt.Sprintf("{\"key_%d\":\"value_%d\"}", i, i))
                r.Close()
            }
        }()
    }



    select {

    }

}