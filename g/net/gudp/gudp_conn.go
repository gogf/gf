// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"io"
	"net"
	"time"

	"github.com/gogf/gf/g/errors/gerror"
)

// 封装的UDP链接对象
type Conn struct {
	*net.UDPConn                 // 底层链接对象
	raddr          *net.UDPAddr  // 远程地址
	recvDeadline   time.Time     // 读取超时时间
	sendDeadline   time.Time     // 写入超时时间
	recvBufferWait time.Duration // 读取全部缓冲区数据时，读取完毕后的写入等待间隔
}

const (
	gDEFAULT_RETRY_INTERVAL   = 100              // (毫秒)默认重试时间间隔
	gDEFAULT_READ_BUFFER_SIZE = 64               // (KB)默认数据读取缓冲区大小
	gRECV_ALL_WAIT_TIMEOUT    = time.Millisecond // 读取全部缓冲数据时，没有缓冲数据时的等待间隔
)

type Retry struct {
	Count    int // 重试次数
	Interval int // 重试间隔(毫秒)
}

// 创建TCP链接
func NewConn(raddr string, laddr ...string) (*Conn, error) {
	if conn, err := NewNetConn(raddr, laddr...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// 将*net.UDPConn对象转换为*Conn对象
func NewConnByNetConn(udp *net.UDPConn) *Conn {
	return &Conn{
		UDPConn:        udp,
		recvDeadline:   time.Time{},
		sendDeadline:   time.Time{},
		recvBufferWait: gRECV_ALL_WAIT_TIMEOUT,
	}
}

// 发送数据
func (c *Conn) Send(data []byte, retry ...Retry) (err error) {
	for {
		if c.raddr != nil {
			_, err = c.WriteToUDP(data, c.raddr)
		} else {
			_, err = c.Write(data)
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
			return nil
		}
	}
}

// 接收UDP协议数据.
//
// 注意事项：
// 1、UDP协议存在消息边界，因此使用 length < 0 可以获取缓冲区所有消息包数据，即一个完整包；
// 2、当length = 0时，表示获取当前的缓冲区数据，获取一次后立即返回；
func (c *Conn) Recv(length int, retry ...Retry) ([]byte, error) {
	var err error          // 读取错误
	var size int           // 读取长度
	var index int          // 已读取长度
	var raddr *net.UDPAddr // 当前读取的远程地址
	var buffer []byte      // 读取缓冲区
	var bufferWait bool    // 是否设置读取的超时时间

	if length > 0 {
		buffer = make([]byte, length)
	} else {
		buffer = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
	}

	for {
		if length < 0 && index > 0 {
			bufferWait = true
			if err = c.SetReadDeadline(time.Now().Add(c.recvBufferWait)); err != nil {
				return nil, err
			}
		}
		size, raddr, err = c.ReadFromUDP(buffer[index:])
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
				if err = c.SetReadDeadline(c.recvDeadline); err != nil {
					return nil, err
				}
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
		// 只获取一次数据
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

// 发送数据并等待接收返回数据
func (c *Conn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.Recv(receive, retry...)
	} else {
		return nil, err
	}
}

// 带超时时间的数据获取
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer func() {
		err = gerror.Wrap(c.SetRecvDeadline(time.Time{}), "SetRecvDeadline error")
	}()
	data, err = c.Recv(length, retry...)
	return
}

// 带超时时间的数据发送
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer func() {
		err = gerror.Wrap(c.SetSendDeadline(time.Time{}), "SetSendDeadline error")
	}()
	err = c.Send(data, retry...)
	return
}

// 发送数据并等待接收返回数据(带返回超时等待时间)
func (c *Conn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.RecvWithTimeout(receive, timeout, retry...)
	} else {
		return nil, err
	}
}

func (c *Conn) SetDeadline(t time.Time) error {
	err := c.UDPConn.SetDeadline(t)
	if err == nil {
		c.recvDeadline = t
		c.sendDeadline = t
	}
	return err
}

func (c *Conn) SetRecvDeadline(t time.Time) error {
	err := c.SetReadDeadline(t)
	if err == nil {
		c.recvDeadline = t
	}
	return err
}

func (c *Conn) SetSendDeadline(t time.Time) error {
	err := c.SetWriteDeadline(t)
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

// 不能使用c.conn.RemoteAddr()，其返回为nil，
// 这里使用c.raddr获取远程连接地址。
func (c *Conn) RemoteAddr() net.Addr {
	//return c.conn.RemoteAddr()
	return c.raddr
}
