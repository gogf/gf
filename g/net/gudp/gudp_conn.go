// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gudp

import (
    "net"
    "time"
    "io"
)

// 封装的链接对象
type Conn struct {
    *net.UDPConn
    raddr *net.UDPAddr
}

const (
    gDEFAULT_RETRY_INTERVAL   = 100   // (毫秒)默认重试时间间隔
    gDEFAULT_READ_BUFFER_SIZE = 1024  // 默认数据读取缓冲区大小
)

type Retry struct {
    Count    int  // 重试次数
    Interval int  // 重试间隔(毫秒)
}

// 创建TCP链接
func NewConn(raddr string, laddr...string) (*Conn, error) {
    if conn, err := NewNetConn(raddr, laddr...); err == nil {
        return NewConnByNetConn(conn), nil
    } else {
        return nil, err
    }
}

// 将*net.UDPConn对象转换为*Conn对象
func NewConnByNetConn(udp *net.UDPConn) *Conn {
    return &Conn {
        UDPConn : udp,
    }
}

// 发送数据
func (c *Conn) Send(data []byte, retry...Retry) error {
    var err     error
    var size    int
    var length  int
    for {
        if c.raddr != nil {
            size, err = c.WriteToUDP(data, c.raddr)
        } else {
            size, err = c.Write(data)
        }
        if err != nil {
            // 链接已关闭
            if err == io.EOF {
                return err
            }
            // 其他错误，重试之后仍不能成功
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
            length += size
            if length == len(data) {
                return nil
            }
        }
    }
    return nil
}

// 接收数据
func (c *Conn) Receive(length int, retry...Retry) ([]byte, error) {
    var err     error         // 读取错误
    var size    int           // 读取长度
    var index   int           // 已读取长度
    var buffer  []byte        // 读取缓冲区

    if length > 0 {
        buffer = make([]byte, length)
    } else {
        buffer = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
    }
    for {
        size, c.raddr, err = c.ReadFromUDP(buffer[index:])
        if size > 0 {
            index += size
            if length > 0 {
                // 如果指定了读取大小，那么必须读取到指定长度才返回
                if index == length {
                    break
                }
            } else {
                // 如果长度超过了自定义的读取缓冲区，那么自动增长
                if index >= gDEFAULT_READ_BUFFER_SIZE {
                    buffer = append(buffer, make([]byte, gDEFAULT_READ_BUFFER_SIZE)...)
                }
                if size < gDEFAULT_READ_BUFFER_SIZE {
                    break
                }
            }
        }
        if err != nil {
            // 链接已关闭
            if err == io.EOF {
                break
            }
            if len(retry) > 0 {
                // 其他错误，重试之后仍不能成功
                if retry[0].Count == 0 {
                    break
                }
                retry[0].Count--
                if retry[0].Interval == 0 {
                    retry[0].Interval = gDEFAULT_RETRY_INTERVAL
                }
                time.Sleep(time.Duration(retry[0].Interval) * time.Millisecond)
                continue
            }
            break
        }
    }
    return buffer[:index], err
}

// 发送数据并等待接收返回数据
func (c *Conn) SendReceive(data []byte, receive int, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.Receive(receive, retry...)
    } else {
        return nil, err
    }
}