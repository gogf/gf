// Copyright 2018 gf Author(https://gitee.com/johng/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gf.

package gtcp

import (
    "net"
    "time"
    "io"
    "bufio"
    "fmt"
)

// 封装的链接对象
type Conn struct {
    net.Conn
    reader *bufio.Reader // 当前链接的缓冲读取对象
}

// 创建TCP链接
func NewConn(addr string, timeout...int) (*Conn, error) {
    if conn, err := NewNetConn(addr, timeout...); err == nil {
        return &Conn {
            Conn : conn,
        }, nil
    } else {
        return nil, err
    }
}

// 将net.Conn接口对象转换为*gtcp.Conn对象(注意递归影响，因为*gtcp.Conn本身也实现了net.Conn接口)
func NewConnByNetConn(conn net.Conn) *Conn {
    return &Conn { Conn: conn }
}

// 发送数据
func (c *Conn) Send(data []byte, retry...Retry) error {
    length := 0
    for {
        n, err := c.Write(data)
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
            length += n
            if length == len(data) {
                return nil
            }
        }
    }
}

// 获取数据，指定读取的数据长度(length < 1表示获取所有可读数据)，以及重试策略(retry)
// 需要注意：
// 1、往往在socket通信中需要指定固定的数据结构，并在设定对应的长度字段，并在读取数据时便于区分包大小；
// 2、当length < 1时表示获取缓冲区所有的数据，但是可能会引起包解析问题(可能出现非完整的包情况)，因此需要解析端注意解析策略；
func (c *Conn) Receive(length int, retry...Retry) ([]byte, error) {
    var err    error  // 读取错误
    var size   int    // 读取长度
    var index  int    // 已读取长度
    var buffer []byte // 读取缓冲区

    if c.reader == nil {
        c.reader = bufio.NewReader(c)
    }
    if length > 0 {
        buffer = make([]byte, length)
    } else {
        buffer = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
    }

    for {
        time.Sleep(time.Second)
        fmt.Println(c.reader.Buffered())
        size, err = c.reader.Read(buffer[index:])
        if size > 0 {
            index += size
            if length > 0 {
                // 如果指定了读取大小，那么必须读取到指定长度才返回
                if index == length {
                    break
                }
            } else {
                // 否则读取所有缓冲区数据，直到没有可读数据为止
                //if c.reader.Buffered() < 1 {
                //    break
                //}
                // 如果长度超过了自定义的读取缓冲区，那么自动增长
                if index >= gDEFAULT_READ_BUFFER_SIZE {
                    buffer = append(buffer, make([]byte, gDEFAULT_READ_BUFFER_SIZE)...)
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

// 带超时时间的数据获取
func (c *Conn) ReceiveWithTimeout(length int, timeout time.Duration, retry...Retry) ([]byte, error) {
    c.SetReadDeadline(time.Now().Add(timeout))
    defer c.SetReadDeadline(time.Time{})
    return c.Receive(length, retry...)
}

// 带超时时间的数据发送
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry...Retry) error {
    c.SetWriteDeadline(time.Now().Add(timeout))
    defer c.SetWriteDeadline(time.Time{})
    return c.Send(data, retry...)
}

// 发送数据并等待接收返回数据
func (c *Conn) SendReceive(data []byte, receive int, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.Receive(receive, retry...)
    } else {
        return nil, err
    }
}

// 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendReceiveWithTimeout(data []byte, receive int, timeout time.Duration, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.ReceiveWithTimeout(receive, timeout, retry...)
    } else {
        return nil, err
    }
}