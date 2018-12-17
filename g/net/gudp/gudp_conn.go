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
    conn          *net.UDPConn   // 底层链接对象
    raddr         *net.UDPAddr   // 远程地址
    recvDeadline   time.Time     // 读取超时时间
    sendDeadline   time.Time     // 写入超时时间
    recvBufferWait time.Duration // 读取全部缓冲区数据时，读取完毕后的写入等待间隔
}

const (
    gDEFAULT_RETRY_INTERVAL   = 100   // (毫秒)默认重试时间间隔
    gDEFAULT_READ_BUFFER_SIZE = 1024  // 默认数据读取缓冲区大小
    gRECV_ALL_WAIT_TIMEOUT    = time.Millisecond // 读取全部缓冲数据时，没有缓冲数据时的等待间隔
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
        conn           : udp,
        recvDeadline   : time.Time{},
        sendDeadline   : time.Time{},
        recvBufferWait : gRECV_ALL_WAIT_TIMEOUT,
    }
}

// 发送数据
func (c *Conn) Send(data []byte, retry...Retry) error {
    var err     error
    var size    int
    var length  int
    for {
        if c.raddr != nil {
            size, err = c.conn.WriteToUDP(data, c.raddr)
        } else {
            size, err = c.conn.Write(data)
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
}

// 接收数据
func (c *Conn) Recv(length int, retry...Retry) ([]byte, error) {
    var err        error        // 读取错误
    var size       int          // 读取长度
    var index      int          // 已读取长度
    var raddr      *net.UDPAddr // 当前读取的远程地址
    var buffer     []byte       // 读取缓冲区
    var bufferWait bool         // 是否设置读取的超时时间

    if length > 0 {
        buffer = make([]byte, length)
    } else {
        buffer = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
    }

    for {
        if length <= 0 && index > 0 {
            bufferWait = true
            c.conn.SetReadDeadline(time.Now().Add(c.recvBufferWait))
        }
        size, raddr, err = c.conn.ReadFromUDP(buffer[index:])
        if err == nil {
            c.raddr = raddr
        }
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

// 发送数据并等待接收返回数据
func (c *Conn) SendRecv(data []byte, receive int, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.Recv(receive, retry...)
    } else {
        return nil, err
    }
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
func (c *Conn) SetRecvBufferWait(d time.Duration) {
    c.recvBufferWait = d
}

func (c *Conn) LocalAddr() net.Addr {
    return c.conn.LocalAddr()
}

func (c *Conn) RemoteAddr() net.Addr {
    return c.conn.RemoteAddr()
}

func (c *Conn) Close() error {
    return c.conn.Close()
}