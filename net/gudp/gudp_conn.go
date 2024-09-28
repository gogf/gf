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

// localConn provides common operations for udp connection.
type localConn struct {
	*net.UDPConn           // Underlying UDP connection.
	deadlineRecv time.Time // Timeout point for reading data.
	deadlineSend time.Time // Timeout point for writing data.
}

const (
	defaultRetryInterval  = 100 * time.Millisecond // Retry interval.
	defaultReadBufferSize = 1024                   // (Byte)Buffer size.
)

// Retry holds the retry options.
// TODO replace with standalone retry package.
type Retry struct {
	Count    int           // Max retry count.
	Interval time.Duration // Retry interval.
}

// Recv receives and returns data from remote address.
// The parameter `buffer` is used for customizing the receiving buffer size.
// If `buffer` <= 0, it uses the default buffer size, which is 1024 byte.
//
// There's package border in UDP protocol, we can receive a complete package if specified
// buffer size is big enough. VERY NOTE that we should receive the complete package in once
// or else the leftover package data would be dropped.
func (c *localConn) Recv(buffer int, retry ...Retry) ([]byte, *net.UDPAddr, error) {
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
	return data[:size], remoteAddr, err
}

// SetDeadline sets the read and write deadlines associated with the connection.
func (c *localConn) SetDeadline(t time.Time) (err error) {
	if err = c.UDPConn.SetDeadline(t); err == nil {
		c.deadlineRecv = t
		c.deadlineSend = t
	} else {
		err = gerror.Wrapf(err, `SetDeadline for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineRecv sets the read deadline associated with the connection.
func (c *localConn) SetDeadlineRecv(t time.Time) (err error) {
	if err = c.SetReadDeadline(t); err == nil {
		c.deadlineRecv = t
	} else {
		err = gerror.Wrapf(err, `SetDeadlineRecv for connection failed with "%s"`, t)
	}
	return err
}

// SetDeadlineSend sets the deadline of sending for current connection.
func (c *localConn) SetDeadlineSend(t time.Time) (err error) {
	if err = c.SetWriteDeadline(t); err == nil {
		c.deadlineSend = t
	} else {
		err = gerror.Wrapf(err, `SetDeadlineSend for connection failed with "%s"`, t)
	}
	return err
}
