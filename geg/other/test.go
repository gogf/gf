package main

import (
    "github.com/Shopify/sarama"
    "log"
    "sync"
)

func main() {
    group := "test-group-100"
    topic := "test"
    config := sarama.NewConfig()
    config.Version = sarama.V0_10_2_0
    client, err := sarama.NewClient([]string{"localhost:9092"}, config)
    if err != nil {
        log.Fatalln(err)
    }

    offsetManager, err := sarama.NewOffsetManagerFromClient(group, client)
    if err != nil {
        log.Fatalln(err)
    }

    pids, err := client.Partitions(topic)
    if err != nil {
        log.Fatalln(err)
    }

    consumer, err := sarama.NewConsumerFromClient(client)
    if err != nil {
        log.Fatalln(err)
    }

    defer consumer.Close()

    wg := &sync.WaitGroup{}

    for _, v := range pids {
        wg.Add(1)
        go consume(wg, consumer, offsetManager, v)
    }

    wg.Wait()
}

func consume(wg *sync.WaitGroup, c sarama.Consumer, om sarama.OffsetManager, p int32) {
    defer wg.Done()

    pom, err := om.ManagePartition("test", p)
    if err != nil {
        log.Fatalln(err)
    }
    defer pom.Close()

    offset, _ := pom.NextOffset()
    if offset == -1 {
        offset = sarama.OffsetOldest
    }

    pc, err := c.ConsumePartition("test", p, 6)
    if err != nil {
        log.Fatalln(err)
    }
    defer pc.Close()

    for msg := range pc.Messages() {
        log.Printf("[%v] Consumed message offset %v\n", p, msg.Offset)
        pom.MarkOffset(msg.Offset + 1, "")
    }
}