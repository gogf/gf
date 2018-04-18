package main

import (
    "fmt"
    "github.com/clbanning/mxj"
)

func main() {
    m := make(map[string]interface{})
    m["m"] = map[string]string {
        "k" : "v",
    }
    b, _ := mxj.Map(m).Xml()
    fmt.Println(string(b))

    // expect {"m":{"k":"v"}} , but I got >UNKNOWN/>
}

