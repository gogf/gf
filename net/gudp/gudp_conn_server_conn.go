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

// ServerConn holds the server side connection.
type ServerConn struct {
	*localConn
}

// NewServerConn creates an udp connection that listens to `localAddress`.
func NewServerConn(listenedConn *net.UDPConn) *ServerConn {
	return &ServerConn{
		localConn: &localConn{
			UDPConn: listenedConn,
		},
	}
}

// Send writes data to remote address.
func (c *ServerConn) Send(data []byte, remoteAddr *net.UDPAddr, retry ...Retry) (err error) {
	for {
		_, err = c.WriteToUDP(data, remoteAddr)
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
