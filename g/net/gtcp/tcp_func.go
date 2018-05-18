// Copyright 2017 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtcp

import (
    "io"
    "net"
    "time"
)

const (
    gDEFAULT_RETRY_INTERVAL = 100 // 默认重试时间间隔
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

// 获取数据
func Receive(conn net.Conn, retry...Retry) []byte {
    size   := 1024
    data   := make([]byte, 0)
    for {
        buffer      := make([]byte, size)
        length, err := conn.Read(buffer)
        if length < 1 && err != nil {
            if err == io.EOF || len(retry) == 0 || retry[0].Count == 0 {
                break
            }
            if len(retry) > 0 {
                retry[0].Count--
                if retry[0].Interval == 0 {
                    retry[0].Interval = gDEFAULT_RETRY_INTERVAL
                }
                time.Sleep(time.Duration(retry[0].Interval) * time.Millisecond)
            }
        } else {
            data = append(data, buffer[0:length]...)
            if err == io.EOF {
                break
            }
        }
    }
    return data
}

// 带超时时间的数据获取
func ReceiveWithTimeout(conn net.Conn, timeout time.Duration, retry...Retry) []byte {
    conn.SetReadDeadline(time.Now().Add(timeout))
    return Receive(conn, retry...)
}

// 发送数据
func Send(conn net.Conn, data []byte, retry...Retry) error {
    for {
        _, err := conn.Write(data)
        if err != nil {
            if len(retry) == 0 || retry[0].Count == 0 {
                return err
            }
            if len(retry) > 0 {
                retry[0].Count--
                if retry[0].Interval == 0 {
                    retry[0].Interval = gDEFAULT_RETRY_INTERVAL
                }
                time.Sleep(time.Duration(retry[0].Interval) * time.Millisecond)
            }
        } else {
            return nil
        }
    }
}

// 带超时时间的数据发送
func SendWithTimeout(conn net.Conn, data []byte, timeout time.Duration, retry...Retry) error {
    conn.SetWriteDeadline(time.Now().Add(timeout))
    return Send(conn, data, retry...)
}