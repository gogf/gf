package main

import (
    "fmt"
    "gitee.com/johng/gf/g/database/gkafka"
    "time"
)

// 创建kafka消费客户端
func newKafkaClientConsumer(topic, group string) *gkafka.Client {
    kafkaConfig               := gkafka.NewConfig()
    kafkaConfig.Servers        = "localhost:9092"
    kafkaConfig.AutoMarkOffset = false
    kafkaConfig.Topics         = topic
    kafkaConfig.GroupId        = group
    return gkafka.NewClient(kafkaConfig)
}

func main () {
    client := newKafkaClientConsumer("test", "test-group-1")
    defer client.Close()

    for {
        fmt.Println("reading...")
        for i := 1; i < 10; i++ {
            if msg, err := client.Receive(); err != nil {
                fmt.Println(err)
            } else {
                fmt.Println(string(msg.Value))
            }
        }
        time.Sleep(3*time.Second)
    }

}
