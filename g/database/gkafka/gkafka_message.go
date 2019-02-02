// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gkafka

// 自动标记已读取
func (msg *Message) MarkOffset() {
    if msg.consumerMsg != nil && msg.client != nil && msg.client.consumer != nil {
        msg.client.consumer.MarkOffset(msg.consumerMsg, "")
    }
}