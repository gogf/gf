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
	gDEFAULT_READ_BUFFER_SIZE = 1024                   // (Byte)Buffer size.
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

// Recv receives and returns data from remote address.
// The parameter <buffer> is used for customizing the receiving buffer size. If <buffer> <= 0,
// it uses the default buffer size, which is 1024 byte.
//
// There's package border in UDP protocol, we can receive a complete package if specified
// buffer size is big enough. VERY NOTE that we should receive the complete package in once
// or else the leftover package data would be dropped.
func (c *Conn) Recv(buffer int, retry ...Retry) ([]byte, error) {
	var err error               // Reading error.
	var size int                // Reading size.
	var data []byte             // Buffer object.
	var remoteAddr *net.UDPAddr // Current remote address for reading.
	if buffer > 0 {
		data = make([]byte, buffer)
	} else {
		data = make([]byte, gDEFAULT_READ_BUFFER_SIZE)
	}
	for {
		size, remoteAddr, err = c.ReadFromUDP(data)
		if err == nil {
			c.remoteAddr = remoteAddr
		}
		if err != nil {
			// Connection closed.
			if err == io.EOF {
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
		break
	}
	return data[:size], err
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
