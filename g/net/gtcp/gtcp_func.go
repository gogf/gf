// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
    "net"
    "time"
)

const (
    gDEFAULT_RETRY_INTERVAL   = 100   // (毫秒)默认重试时间间隔
    gDEFAULT_READ_BUFFER_SIZE = 1024  // 默认数据读取缓冲区大小
)

type Retry struct {
    Count    int  // 重试次数
    Interval int  // 重试间隔(毫秒)
}

// 常见的二进制数据校验方式，生成校验结果
func Checksum(buffer []byte) uint32 {
    var checksum uint32
    for _, b := range buffer {
        checksum += uint32(b)
    }
    return checksum
}

// 创建原生TCP链接, addr地址格式形如：127.0.0.1:80
func NewNetConn(addr string, timeout...int) (net.Conn, error) {
    if len(timeout) > 0 {
        return net.DialTimeout("tcp", addr, time.Duration(timeout[0]) * time.Millisecond)
    } else {
        return net.Dial("tcp", addr)
    }
}

// (面向短链接)发送数据
func Send(addr string, data []byte, retry...Retry) error {
    conn, err := NewConn(addr)
    if err != nil {
        return err
    }
    defer conn.Close()
    return conn.Send(data, retry...)
}

// (面向短链接)发送数据并等待接收返回数据
func SendRecv(addr string, data []byte, receive int, retry...Retry) ([]byte, error) {
	conn, err := NewConn(addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecv(data, receive, retry...)
}

// (面向短链接)带超时时间的数据发送
func SendWithTimeout(addr string, data []byte, timeout time.Duration, retry...Retry) error {
    conn, err := NewConn(addr)
    if err != nil {
        return err
    }
    defer conn.Close()
    return conn.SendWithTimeout(data, timeout, retry...)
}

// (面向短链接)发送数据并等待接收返回数据(带返回超时等待时间)
func SendRecvWithTimeout(addr string, data []byte, receive int, timeout time.Duration, retry...Retry) ([]byte, error) {
    conn, err := NewConn(addr)
    if err != nil {
        return nil, err
    }
    defer conn.Close()
    return conn.SendRecvWithTimeout(data, receive, timeout, retry...)
}

// 判断是否是超时错误
func isTimeout(err error) bool {
    if err == nil {
        return false
    }
    if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
        return true
    }
    return false
}