package main

import (
    "github.com/gogf/gf/g/database/gkafka"
    "fmt"
)

func main () {
    config        := gkafka.NewConfig()
    config.Servers = "localhost:9092"

    client := gkafka.NewClient(config)
    defer client.Close()

    fmt.Println(client.Topics())
}
