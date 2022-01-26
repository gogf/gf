// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"io"
	"net"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// Conn handles the UDP connection.
type Conn struct {
	*net.UDPConn                    // Underlying UDP connection.
	remoteAddr        *net.UDPAddr  // Remote address.
	receiveDeadline   time.Time     // Timeout point for reading data.
	sendDeadline      time.Time     // Timeout point for writing data.
	receiveBufferWait time.Duration // Interval duration for reading buffer.
}

const (
	defaultRetryInterval  = 100 * time.Millisecond // Retry interval.
	defaultReadBufferSize = 1024                   // (Byte)Buffer size.
	receiveAllWaitTimeout = time.Millisecond       // Default interval for reading buffer.
)

type Retry struct {
	Count    int           // Max retry count.
	Interval time.Duration // Retry interval.
}

// NewConn creates UDP connection to `remoteAddress`.
// The optional parameter `localAddress` specifies the local address for connection.
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
		UDPConn:           udp,
		receiveDeadline:   time.Time{},
		sendDeadline:      time.Time{},
		receiveBufferWait: receiveAllWaitTimeout,
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
				err = gerror.Wrap(err, `Write data failed`)
				return err
			}
			if len(retry) > 0 {
				retry[0].Count--
				if retry[0].Interval == 0 {
					retry[0].Interval = defaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
			}
		} else {
			return nil
		}
	}
}

// Recv receives and returns data from remote address.
// The parameter `buffer` is used for customizing the receiving buffer size. If `buffer` <= 0,
// it uses the default buffer size, which is 1024 byte.
//
// There's package border in UDP protocol, we can receive a complete package if specified
// buffer size is big enough. VERY NOTE that we should receive the complete package in once
// or else the leftover package data would be dropped.
func (c *Conn) Recv(buffer int, retry ...Retry) ([]byte, error) {
	var (
		err        error        // Reading error
		size       int          // Reading size
		data       []byte       // Buffer object
		remoteAddr *net.UDPAddr // Current remote address for reading
	)
	if buffer > 0 {
		data = make([]byte, buffer)
	} else {
		data = make([]byte, defaultReadBufferSize)
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
					retry[0].Interval = defaultRetryInterval
				}
				time.Sleep(retry[0].Interval)
				continue
			}
			err = gerror.Wrap(err, `ReadFromUDP failed`)
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
	if err = c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to connection with timeout.
func (c *Conn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err = c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
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

func (c *Conn) SetDeadline(t time.Time) (err error) {
	if err = c.UDPConn.SetDeadline(t); err == nil {
		c.receiveDeadline = t
		c.sendDeadline = t
	} else {
		err = gerror.Wrapf(err, `SetDeadline for connection failed with "%s"`, t)
	}
	return err
}

func (c *Conn) SetRecvDeadline(t time.Time) (err error) {
	if err = c.SetReadDeadline(t); err == nil {
		c.receiveDeadline = t
	} else {
		err = gerror.Wrapf(err, `SetReadDeadline for connection failed with "%s"`, t)
	}
	return err
}

func (c *Conn) SetSendDeadline(t time.Time) (err error) {
	if err = c.SetWriteDeadline(t); err == nil {
		c.sendDeadline = t
	} else {
		err = gerror.Wrapf(err, `SetWriteDeadline for connection failed with "%s"`, t)
	}
	return err
}

// SetRecvBufferWait sets the buffer waiting timeout when reading all data from connection.
// The waiting duration cannot be too long which might delay receiving data from remote address.
func (c *Conn) SetRecvBufferWait(d time.Duration) {
	c.receiveBufferWait = d
}

// RemoteAddr returns the remote address of current UDP connection.
// Note that it cannot use c.conn.RemoteAddr() as it is nil.
func (c *Conn) RemoteAddr() net.Addr {
	// return c.conn.RemoteAddr()
	return c.remoteAddr
}
