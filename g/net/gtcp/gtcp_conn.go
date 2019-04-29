// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
    "net"
    "time"
    "io"
    "bufio"
    "bytes"
)

// 封装的链接对象
type Conn struct {
    conn           net.Conn      // 底层tcp对象
    reader         *bufio.Reader // 当前链接的缓冲读取对象
    buffer         []byte        // 读取缓冲区(用于数据读取时的缓冲区处理)
    recvDeadline   time.Time     // 读取超时时间
    sendDeadline   time.Time     // 写入超时时间
    recvBufferWait time.Duration // 读取全部缓冲区数据时，读取完毕后的写入等待间隔
}

const (
	// 读取全部缓冲数据时，没有缓冲数据时的等待间隔
    gRECV_ALL_WAIT_TIMEOUT = time.Millisecond
)

// 创建TCP链接
func NewConn(addr string, timeout...int) (*Conn, error) {
    if conn, err := NewNetConn(addr, timeout...); err == nil {
        return NewConnByNetConn(conn), nil
    } else {
        return nil, err
    }
}

// 将net.Conn接口对象转换为*gtcp.Conn对象
func NewConnByNetConn(conn net.Conn) *Conn {
    return &Conn {
        conn           : conn,
        reader         : bufio.NewReader(conn),
        recvDeadline   : time.Time{},
        sendDeadline   : time.Time{},
        recvBufferWait : gRECV_ALL_WAIT_TIMEOUT,
    }
}

// 关闭连接
func (c *Conn) Close() {
    c.conn.Close()
}

// 发送数据
func (c *Conn) Send(data []byte, retry...Retry) error {
    length := 0
    for {
        n, err := c.conn.Write(data)
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
func (c *Conn) Recv(length int, retry...Retry) ([]byte, error) {
    var err        error  // 读取错误
    var size       int    // 读取长度
    var index      int    // 已读取长度
    var buffer     []byte // 读取缓冲区
    var bufferWait bool   // 是否设置读取的超时时间

    if length > 0 {
        buffer = make([]byte, length)
    } else {
        buffer = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
    }

    for {
        // 缓冲区数据写入等待处理。
        // 如果已经读取到数据(这点很关键，表明缓冲区已经有数据，剩下的操作就是将所有数据读取完毕)，
        // 那么可以设置读取全部缓冲数据的超时时间；如果没有接收到任何数据，那么将会进入读取阻塞(或者自定义的超时阻塞);
        // 仅对读取全部缓冲区数据操作有效
        if length <= 0 && index > 0 {
            bufferWait = true
            c.conn.SetReadDeadline(time.Now().Add(c.recvBufferWait))
        }
        size, err = c.reader.Read(buffer[index:])
        if size > 0 {
            index += size
            if length > 0 {
                // 如果指定了读取大小，那么必须读取到指定长度才返回
                if index == length {
                    break
                }
            } else {
                if index >= gDEFAULT_READ_BUFFER_SIZE {
	                // 如果长度超过了自定义的读取缓冲区，那么自动增长
                    buffer = append(buffer, make([]byte, gDEFAULT_READ_BUFFER_SIZE)...)
                } else {
                	// 如果第一次读取的数据并未达到缓冲变量长度，那么直接返回
                	if !bufferWait {
						break
	                }
                }
            }
        }
        if err != nil {
            // 链接已关闭
            if err == io.EOF {
                break
            }
            // 判断数据是否全部读取完毕(由于超时机制的存在，获取的数据完整性不可靠)
            if bufferWait && isTimeout(err) {
                c.conn.SetReadDeadline(c.recvDeadline)
                err = nil
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

// 按行读取数据，阻塞读取，直到完成一行读取位置(末尾以'\n'结尾，返回数据不包含换行符)
func (c *Conn) RecvLine(retry...Retry) ([]byte, error) {
    var err    error
    var buffer []byte
    data := make([]byte, 0)
    for {
        buffer, err = c.Recv(1, retry...)
        if len(buffer) > 0 {
            data = append(data, buffer...)
            if buffer[0] == '\n' {
                break
            }
        }
        if err != nil {
            break
        }
    }
    if len(data) > 0 {
        data = bytes.TrimRight(data, "\n\r")
    }
    return data, err
}

// 带超时时间的数据获取
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry...Retry) ([]byte, error) {
    c.SetRecvDeadline(time.Now().Add(timeout))
    defer c.SetRecvDeadline(time.Time{})
    return c.Recv(length, retry...)
}

// 带超时时间的数据发送
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry...Retry) error {
    c.SetSendDeadline(time.Now().Add(timeout))
    defer c.SetSendDeadline(time.Time{})
    return c.Send(data, retry...)
}

// 发送数据并等待接收返回数据
func (c *Conn) SendRecv(data []byte, receive int, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.Recv(receive, retry...)
    } else {
        return nil, err
    }
}

// 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.RecvWithTimeout(receive, timeout, retry...)
    } else {
        return nil, err
    }
}

func (c *Conn) SetDeadline(t time.Time) error {
    err := c.conn.SetDeadline(t)
    if err == nil {
        c.recvDeadline = t
        c.sendDeadline = t
    }
    return err
}

func (c *Conn) SetRecvDeadline(t time.Time) error {
    err := c.conn.SetReadDeadline(t)
    if err == nil {
        c.recvDeadline = t
    }
    return err
}

func (c *Conn) SetSendDeadline(t time.Time) error {
    err := c.conn.SetWriteDeadline(t)
    if err == nil {
        c.sendDeadline = t
    }
    return err
}

// 读取全部缓冲区数据时，读取完毕后的写入等待间隔，如果超过该等待时间后仍无可读数据，那么读取操作返回。
// 该时间间隔不能设置得太大，会影响Recv读取时长(默认为1毫秒)。
func (c *Conn) SetRecvBufferWait(bufferWaitDuration time.Duration) {
    c.recvBufferWait = bufferWaitDuration
}

func (c *Conn) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}