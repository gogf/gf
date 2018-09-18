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
    group  := "test-group-11"
    client := newKafkaClientConsumer("test", group)
    defer client.Close()

    for {
        fmt.Println(group + " reading...")
        for {
            if msg, err := client.Receive(); err != nil {
                fmt.Println(err)
            } else {
                fmt.Println(msg.Partition, msg.Offset, string(msg.Value))
                msg.MarkOffset()
            }
        }
        time.Sleep(3*time.Second)
    }

}
