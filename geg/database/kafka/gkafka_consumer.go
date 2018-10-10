package main

import (
    "fmt"
    "gitee.com/johng/gf/g/database/gkafka"
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
    group  := "test-group"
    topic  := "test"
    client := newKafkaClientConsumer(topic, group)
    defer client.Close()

    // 标记开始读取的offset位置
    client.MarkOffset(topic, 0, 6)
    for {
        if msg, err := client.Receive(); err != nil {
            fmt.Println(err)
            break
        } else {
            fmt.Println(msg.Partition, msg.Offset, string(msg.Value))
            msg.MarkOffset()
        }
    }
}
