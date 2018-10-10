package main

import (
    "gitee.com/johng/gf/g/database/gkafka"
    "fmt"
    "gitee.com/johng/gf/g/os/gtime"
    "time"
)

// 创建kafka生产客户端
func newKafkaClientProducer(topic string) *gkafka.Client {
    kafkaConfig               := gkafka.NewConfig()
    kafkaConfig.Servers        = "localhost:9092"
    kafkaConfig.AutoMarkOffset = false
    kafkaConfig.Topics         = topic
    return gkafka.NewClient(kafkaConfig)
}

func main () {
    client := newKafkaClientProducer("test")
    defer client.Close()
    for {
        if err := client.SyncSend(&gkafka.Message{Value: []byte(gtime.Now().String())}); err != nil {
            fmt.Println(err)
        }
        time.Sleep(time.Second)
    }
}
