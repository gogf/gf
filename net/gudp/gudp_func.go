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
		network    = `udp`
	)
	remoteAddr, err = net.ResolveUDPAddr(network, remoteAddress)
	if err != nil {
		return nil, gerror.Wrapf(
			err,
			`net.ResolveUDPAddr failed for network "%s", address "%s"`,
			network, remoteAddress,
		)
	}
	if len(localAddress) > 0 {
		localAddr, err = net.ResolveUDPAddr(network, localAddress[0])
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ResolveUDPAddr failed for network "%s", address "%s"`,
				network, localAddress[0],
			)
		}
	}
	conn, err := net.DialUDP(network, localAddr, remoteAddr)
	if err != nil {
		return nil, gerror.Wrapf(
			err,
			`net.DialUDP failed for network "%s", local "%s", remote "%s"`,
			network, localAddr.String(), remoteAddr.String(),
		)
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

// GetFreePort retrieves and returns a port that is free.
func GetFreePort() (port int, err error) {
	var (
		network = `udp`
		address = `:0`
	)
	resolvedAddr, err := net.ResolveUDPAddr(network, address)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ResolveUDPAddr failed for network "%s", address "%s"`,
			network, address,
		)
	}
	l, err := net.ListenUDP(network, resolvedAddr)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ListenUDP failed for network "%s", address "%s"`,
			network, resolvedAddr.String(),
		)
	}
	port = l.LocalAddr().(*net.UDPAddr).Port
	_ = l.Close()
	return
}

// GetFreePorts retrieves and returns specified number of ports that are free.
func GetFreePorts(count int) (ports []int, err error) {
	var (
		network = `udp`
		address = `:0`
	)
	for i := 0; i < count; i++ {
		resolvedAddr, err := net.ResolveUDPAddr(network, address)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ResolveUDPAddr failed for network "%s", address "%s"`,
				network, address,
			)
		}
		l, err := net.ListenUDP(network, resolvedAddr)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ListenUDP failed for network "%s", address "%s"`,
				network, resolvedAddr.String(),
			)
		}
		ports = append(ports, l.LocalAddr().(*net.UDPAddr).Port)
		_ = l.Close()
	}
	return ports, nil
}
