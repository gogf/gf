// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gudp

import (
	"net"

	"github.com/gogf/gf/v2/errors/gerror"
)

// NewNetConn creates and returns a *net.UDPConn with given addresses.
func NewNetConn(remoteAddress string, localAddress ...string) (*net.UDPConn, error) {
	var (
		err        error
		remoteAddr *net.UDPAddr
		localAddr  *net.UDPAddr
	)
	remoteAddr, err = net.ResolveUDPAddr("udp", remoteAddress)
	if err != nil {
		err = gerror.Wrapf(err, `net.ResolveUDPAddr failed for address "%s"`, remoteAddress)
		return nil, err
	}
	if len(localAddress) > 0 {
		localAddr, err = net.ResolveUDPAddr("udp", localAddress[0])
		if err != nil {
			err = gerror.Wrapf(err, `net.ResolveUDPAddr failed for address "%s"`, localAddress[0])
			return nil, err
		}
	}
	conn, err := net.DialUDP("udp", localAddr, remoteAddr)
	if err != nil {
		err = gerror.Wrapf(err, `net.DialUDP failed for local "%s", remote "%s"`, localAddr.String(), remoteAddr.String())
		return nil, err
	}
	return conn, nil
}

// Send writes data to `address` using UDP connection and then closes the connection.
// Note that it is used for short connection usage.
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendRecv writes data to `address` using UDP connection, reads response and then closes the connection.
// Note that it is used for short connection usage.
func SendRecv(address string, data []byte, receive int, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecv(data, receive, retry...)
}
