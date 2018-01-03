// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.
//

package ghttp

import (
    "net/http"
    "gitee.com/johng/gf/g/encoding/gjson"
    "sync"
)

// 服务端请求返回对象
type Response struct {
    http.ResponseWriter
    bufmu  sync.RWMutex // 缓冲区互斥锁
    buffer []byte       // 每个请求的返回数据缓冲区
}

// 返回的固定JSON数据结构
type ResponseJson struct {
    Result  int    `json:"result"`  // 标识消息状态
    Message string `json:"message"` // 消息使用string存储
    Data    []byte `json:"data"`    // 二进制数据(不管什么数据结构)
}

// 返回信息(byte)
func (r *Response) Write(content []byte) {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    r.buffer = append(r.buffer, content...)
}

// 返回信息(string)
func (r *Response) WriteString(content string) {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    r.buffer = append(r.buffer, content...)
}

// 返回固定格式的json
func (r *Response) WriteJson(result int, message string, data []byte) error {
    r.Header().Set("Content-Type", "application/json")
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    if jsonstr, err := gjson.Encode(ResponseJson{ result, message, data }); err != nil {
        return err
    } else {
        r.buffer = append(r.buffer, jsonstr...)
    }
    return nil
}

// 获取当前缓冲区中的数据
func (r *Response) Buffer() []byte {
    r.bufmu.RLock()
    defer r.bufmu.RUnlock()
    return r.buffer
}

// 清空缓冲区内容
func (r *Response) ClearBuffer() {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    r.buffer = make([]byte, 0)
}

// 输出缓冲区数据到客户端
func (r *Response) OutputBuffer() {
    r.bufmu.Lock()
    defer r.bufmu.Unlock()
    if len(r.buffer) > 0 {
        r.ResponseWriter.Write(r.buffer)
        r.buffer = make([]byte, 0)
    }
}
