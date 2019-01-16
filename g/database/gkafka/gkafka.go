// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

// Package gkafka provides producer and consumer client for kafka server.
//
// Kafka客户端.
package gkafka

import (
    "gitee.com/johng/gf/g/os/glog"
    "gitee.com/johng/gf/third/github.com/Shopify/sarama"
    "gitee.com/johng/gf/third/github.com/johng-cn/sarama-cluster"
    "strings"
    "time"
)

var (
    // 当使用Topics方法获取所有topic后，进行过滤忽略的topic，多个以','号分隔
    ignoreTopics = map[string]bool {
        "__consumer_offsets" : true,
    }
)

// kafka Client based on sarama.Config
type Config struct {
    GroupId        string // group id for consumer.
    Servers        string // server list, multiple servers joined by ','.
    Topics         string // topic list, multiple topics joined by ','.
    AutoMarkOffset bool   // auto mark message read after consumer message from server
    sarama.Config
}

// Kafka Client(Consumer/SyncProducer/AsyncProducer)
type Client struct {
    Config        *Config
    consumer      *cluster.Consumer
    rawConsumer   sarama.Consumer
    syncProducer  sarama.SyncProducer
    asyncProducer sarama.AsyncProducer
}

// Kafka Message.
type Message struct {
    Value          []byte
    Key            []byte
    Topic          string
    Partition      int
    Offset         int
    client         *Client
    consumerMsg    *sarama.ConsumerMessage
}


// New a kafka client.
func NewClient(config *Config) *Client {
    return &Client {
        Config : config,
    }
}

// New a default configuration object.
func NewConfig() *Config {
    config       := &Config{}
    config.Config = *sarama.NewConfig()

    // default config for consumer
    config.Consumer.Return.Errors          = true
    config.Consumer.Offsets.Initial        = sarama.OffsetOldest
    config.Consumer.Offsets.CommitInterval = 1 * time.Second

    // default config for producer
    config.Producer.Return.Errors          = true
    config.Producer.Return.Successes       = true
    config.Producer.Timeout                = 5 * time.Second

    config.AutoMarkOffset                  = true
    return config
}

// Close client.
func (client *Client) Close() {
    if client.rawConsumer != nil {
        client.rawConsumer.Close()
    }
    if client.consumer != nil {
        client.consumer.Close()
    }
    if client.syncProducer != nil {
        client.syncProducer.Close()
    }
    if client.asyncProducer != nil {
        client.asyncProducer.Close()
    }
}

// Get all topics from kafka server.
// 这里创建独立的消费客户端获取topics，获取完之后销毁该客户端对象。
func (client *Client) Topics() ([]string, error) {
    if c, err := sarama.NewConsumer(strings.Split(client.Config.Servers, ","), &client.Config.Config); err != nil {
        return nil, err
    } else {
        if topics, err := c.Topics(); err == nil {
            for k, v := range topics {
                if _, ok := ignoreTopics[v]; ok {
                    topics = append(topics[ : k], topics[k + 1 : ]...)
                }
            }
            c.Close()
            return topics, nil
        } else {
            return nil, err
        }
    }
}

// 初始化内部消费客户端
func (client *Client) initConsumer() error {
    if client.consumer == nil {
        config       := cluster.NewConfig()
        config.Config = client.Config.Config
        config.Group.Return.Notifications = false
        c, err := cluster.NewConsumer(strings.Split(client.Config.Servers, ","), client.Config.GroupId, strings.Split(client.Config.Topics, ","), config)
        if err != nil {
            return err
        } else {
            client.consumer = c
        }
    }
    return nil
}

// 标记指定topic 分区开始读取位置
func (client *Client) MarkOffset(topic string, partition int, offset int, metadata...string) error {
    if err := client.initConsumer(); err != nil {
        return err
    }
    meta := ""
    if len(metadata) > 0 {
        meta = metadata[0]
    }
    client.consumer.MarkPartitionOffset(topic, int32(partition), int64(offset), meta)
    return nil
}


// Receive message from kafka from specified topics in config, in BLOCKING way, gkafka will handle offset tracking automatically.
func (client *Client) Receive() (*Message, error) {
    if err := client.initConsumer(); err != nil {
        return nil, err
    }
    errorsChan  := client.consumer.Errors()
    notifyChan  := client.consumer.Notifications()
    messageChan := client.consumer.Messages()
    for {
        select {
            case msg := <- messageChan:
                if client.Config.AutoMarkOffset {
                    client.consumer.MarkOffset(msg, "")
                }
                return &Message {
                    Value       : msg.Value,
                    Key         : msg.Key,
                    Topic       : msg.Topic,
                    Partition   : int(msg.Partition),
                    Offset      : int(msg.Offset),
                    client      : client,
                    consumerMsg : msg,
                }, nil

            case err := <-errorsChan:
                if err != nil {
                    return nil, err
                }

            case <-notifyChan:
        }
    }
}

// Send data to kafka in synchronized way.
func (client *Client) SyncSend(message *Message) error {
    if client.syncProducer == nil {
        if p, err := sarama.NewSyncProducer(strings.Split(client.Config.Servers, ","), &client.Config.Config); err != nil {
            return err
        } else {
            client.syncProducer = p
        }
    }
    for _, topic := range strings.Split(client.Config.Topics, ",") {
        msg := messageToProducerMessage(message)
        msg.Topic = topic
        if _, _, err := client.syncProducer.SendMessage(msg); err != nil {
            return err
        }
    }
    return nil
}

// Send data to kafka in asynchronized way(concurrent safe).
func (client *Client) AsyncSend(message *Message) error {
    if client.asyncProducer == nil {
        if p, err := sarama.NewAsyncProducer(strings.Split(client.Config.Servers, ","), &client.Config.Config); err != nil {
            return err
        } else {
            client.asyncProducer = p
            go func(p sarama.AsyncProducer) {
               errors  := p.Errors()
               success := p.Successes()
               for {
                   select {
                       case err := <-errors:
                           if err != nil {
                               glog.Error(err)
                           }
                       case <-success:
                   }
               }
            }(client.asyncProducer)
        }
    }

    for _, topic := range strings.Split(client.Config.Topics, ",") {
        msg      := messageToProducerMessage(message)
        msg.Topic = topic
        client.asyncProducer.Input() <- msg
    }
    return nil
}

// Convert *gkafka.Message to *sarama.ProducerMessage
func messageToProducerMessage(message *Message) *sarama.ProducerMessage {
    return &sarama.ProducerMessage {
        Topic     : message.Topic,
        Key       : sarama.ByteEncoder(message.Key),
        Value     : sarama.ByteEncoder(message.Value),
        Partition : int32(message.Partition),
        Offset    : int64(message.Offset),
    }
}
