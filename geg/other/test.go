package main

import (
    "bytes"
    "encoding/json"
    "fmt"
)

func main() {
    value := interface{}(nil)
    data  := []byte(`{"n": 123456789}`)
    decoder := json.NewDecoder(bytes.NewReader(data))
    decoder.UseNumber()
    err := decoder.Decode(&value)
    //err   := json.Unmarshal(data, &value)
    fmt.Println(err)
    fmt.Println(value)
}