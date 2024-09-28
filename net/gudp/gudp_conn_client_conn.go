// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"io"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
)

// ClientConn holds the client side connection.
type ClientConn struct {
	*localConn
}

// NewClientConn creates UDP connection to `remoteAddress`.
// The optional parameter `localAddress` specifies the local address for connection.
func NewClientConn(remoteAddress string, localAddress ...string) (*ClientConn, error) {
	udpConn, err := NewNetConn(remoteAddress, localAddress...)
	if err != nil {
		return nil, err
	}
	return &ClientConn{
		localConn: &localConn{
			UDPConn: udpConn,
		},
	}, nil
}

// Send writes data to remote address.
func (c *ClientConn) Send(data []byte, retry ...Retry) (err error) {
	for {
		_, err = c.Write(data)
		if err == nil {
			return nil
		}
		// Connection closed.
		if err == io.EOF {
			return err
		}
		// Still failed even after retrying.
		if len(retry) == 0 || retry[0].Count == 0 {
			return gerror.Wrap(err, `Write data failed`)
		}
		if len(retry) > 0 {
			retry[0].Count--
			if retry[0].Interval == 0 {
				retry[0].Interval = defaultRetryInterval
			}
			time.Sleep(retry[0].Interval)
			continue
		}
		return err
	}
}

// SendRecv writes data to connection and blocks reading response.
func (c *ClientConn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err != nil {
		return nil, err
	}
	result, _, err := c.Recv(receive, retry...)
	return result, err
}
