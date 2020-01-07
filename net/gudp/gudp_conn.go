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
)

// Conn handles the UDP connection.
type Conn struct {
	*net.UDPConn                 // Underlying UDP connection.
	remoteAddr     *net.UDPAddr  // Remote address.
	recvDeadline   time.Time     // Timeout point for reading data.
	sendDeadline   time.Time     // Timeout point for writing data.
	recvBufferWait time.Duration // Interval duration for reading buffer.
}

const (
	gDEFAULT_RETRY_INTERVAL   = 100 * time.Millisecond // Retry interval.
	gDEFAULT_READ_BUFFER_SIZE = 64                     // (KB)Buffer size.
	gRECV_ALL_WAIT_TIMEOUT    = time.Millisecond       // Default interval for reading buffer.
)

type Retry struct {
	Count    int           // Max retry count.
	Interval time.Duration // Retry interval.
}

// NewConn creates UDP connection to <remoteAddress>.
// The optional parameter <localAddress> specifies the local address for connection.
func NewConn(remoteAddress string, localAddress ...string) (*Conn, error) {
	if conn, err := NewNetConn(remoteAddress, localAddress...); err == nil {
		return NewConnByNetConn(conn), nil
	} else {
		return nil, err
	}
}

// NewConnByNetConn creates a UDP connection object with given *net.UDPConn object.
func NewConnByNetConn(udp *net.UDPConn) *Conn {
	return &Conn{
		UDPConn:        udp,
		recvDeadline:   time.Time{},
		sendDeadline:   time.Time{},
		recvBufferWait: gRECV_ALL_WAIT_TIMEOUT,
	}
}

// Send writes data to remote address.
func (c *Conn) Send(data []byte, retry ...Retry) (err error) {
	for {
		if c.remoteAddr != nil {
			_, err = c.WriteToUDP(data, c.remoteAddr)
		} else {
			_, err = c.Write(data)
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				return err
			}
			// Still failed even after retrying.
			if len(retry) == 0 || retry[0].Count == 0 {
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = gDEFAULT_RETRY_INTERVAL
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv receives data from remote address.
//
// Note that,
// 1. There's package border in UDP protocol, so we can receive a complete package if it specifies length < 0.
// 2. If length = 0, it means it receives the data from current buffer and returns immediately.
// 3. If length > 0, it means it blocks reading data from connection until length size was received.
func (c *Conn) Recv(length int, retry ...Retry) ([]byte, error) {
	var err error               // Reading error.
	var size int                // Reading size.
	var index int               // Received size.
	var buffer []byte           // Buffer object.
	var bufferWait bool         // Whether buffer reading timeout set.
	var remoteAddr *net.UDPAddr // Current remote address for reading.

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
		size, remoteAddr, err = c.ReadFromUDP(buffer[index:])
		if err == nil {
			c.remoteAddr = remoteAddr
		}
		if size > 0 {
			index += size
			if length > 0 {
				// It reads til <length> size if <length> is specified.
				if index == length {
					break
				}
			} else {
				if index >= gDEFAULT_READ_BUFFER_SIZE {
					// If it exceeds the buffer size, it then automatically increases its buffer size.
					buffer = append(buffer, make([]byte, gDEFAULT_READ_BUFFER_SIZE)...)
				} else {
					// It returns immediately if received size is lesser than buffer size.
					if !bufferWait {
						break
					}
				}
			}
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
				break
			}
			// Re-set the timeout when reading data.
			if bufferWait && isTimeout(err) {
				if err = c.SetReadDeadline(c.recvDeadline); err != nil {
					return nil, err
				}
				err = nil
				break
			}
			if len(retry) > 0 {
				// It fails even it retried.
				if retry[0].Count == 0 {
					break
				}
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = gDEFAULT_RETRY_INTERVAL
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			break
		}
		// Just read once from buffer.
		if length == 0 {
			break
		}
	}
	return buffer[:index], err
}

// SendRecv writes data to connection and blocks reading response.
func (c *Conn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.Recv(receive, retry...)
	} else {
		return nil, err
	}
}

// RecvWithTimeout reads data from remote address with timeout.
func (c *Conn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to connection with timeout.
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	err = c.Send(data, retry...)
	return
}

// SendRecvWithTimeout writes data to connection and reads response with timeout.
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

// SetRecvBufferWait sets the buffer waiting timeout when reading all data from connection.
// The waiting duration cannot be too long which might delay receiving data from remote address.
func (c *Conn) SetRecvBufferWait(d time.Duration) {
	c.recvBufferWait = d
}

// RemoteAddr returns the remote address of current UDP connection.
// Note that it cannot use c.conn.RemoteAddr() as it's nil.
func (c *Conn) RemoteAddr() net.Addr {
	//return c.conn.RemoteAddr()
	return c.remoteAddr
}
