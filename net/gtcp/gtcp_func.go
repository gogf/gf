// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"crypto/rand"
	"crypto/tls"
	"net"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	defaultConnTimeout    = 30 * time.Second       // Default connection timeout.
	defaultRetryInternal  = 100 * time.Millisecond // Default retry interval.
	defaultReadBufferSize = 128                    // (Byte) Buffer size for reading.
)

type Retry struct {
	Count    int           // Retry count.
	Interval time.Duration // Retry interval.
}

// NewNetConn creates and returns a net.Conn with given address like "127.0.0.1:80".
// The optional parameter `timeout` specifies the timeout for dialing connection.
func NewNetConn(address string, timeout ...time.Duration) (net.Conn, error) {
	var (
		network  = `tcp`
		duration = defaultConnTimeout
	)
	if len(timeout) > 0 {
		duration = timeout[0]
	}
	conn, err := net.DialTimeout(network, address, duration)
	if err != nil {
		err = gerror.Wrapf(
			err,
			`net.DialTimeout failed with network "%s", address "%s", timeout "%s"`,
			network, address, duration,
		)
	}
	return conn, err
}

// NewNetConnTLS creates and returns a TLS net.Conn with given address like "127.0.0.1:80".
// The optional parameter `timeout` specifies the timeout for dialing connection.
func NewNetConnTLS(address string, tlsConfig *tls.Config, timeout ...time.Duration) (net.Conn, error) {
	var (
		network = `tcp`
		dialer  = &net.Dialer{
			Timeout: defaultConnTimeout,
		}
	)
	if len(timeout) > 0 {
		dialer.Timeout = timeout[0]
	}
	conn, err := tls.DialWithDialer(dialer, network, address, tlsConfig)
	if err != nil {
		err = gerror.Wrapf(
			err,
			`tls.DialWithDialer failed with network "%s", address "%s", timeout "%s", tlsConfig "%v"`,
			network, address, dialer.Timeout, tlsConfig,
		)
	}
	return conn, err
}

// NewNetConnKeyCrt creates and returns a TLS net.Conn with given TLS certificate and key files
// and address like "127.0.0.1:80". The optional parameter `timeout` specifies the timeout for
// dialing connection.
func NewNetConnKeyCrt(addr, crtFile, keyFile string, timeout ...time.Duration) (net.Conn, error) {
	tlsConfig, err := LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewNetConnTLS(addr, tlsConfig, timeout...)
}

// Send creates connection to `address`, writes `data` to the connection and then closes the connection.
// The optional parameter `retry` specifies the retry policy when fails in writing data.
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendRecv creates connection to `address`, writes `data` to the connection, receives response
// and then closes the connection.
//
// The parameter `length` specifies the bytes count waiting to receive. It receives all buffer content
// and returns if `length` is -1.
//
// The optional parameter `retry` specifies the retry policy when fails in writing data.
func SendRecv(address string, data []byte, length int, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecv(data, length, retry...)
}

// SendWithTimeout does Send logic with writing timeout limitation.
func SendWithTimeout(address string, data []byte, timeout time.Duration, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.SendWithTimeout(data, timeout, retry...)
}

// SendRecvWithTimeout does SendRecv logic with reading timeout limitation.
func SendRecvWithTimeout(address string, data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecvWithTimeout(data, receive, timeout, retry...)
}

// isTimeout checks whether given `err` is a timeout error.
func isTimeout(err error) bool {
	if err == nil {
		return false
	}
	if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
		return true
	}
	return false
}

// LoadKeyCrt creates and returns a TLS configuration object with given certificate and key files.
func LoadKeyCrt(crtFile, keyFile string) (*tls.Config, error) {
	crtPath, err := gfile.Search(crtFile)
	if err != nil {
		return nil, err
	}
	keyPath, err := gfile.Search(keyFile)
	if err != nil {
		return nil, err
	}
	crt, err := tls.LoadX509KeyPair(crtPath, keyPath)
	if err != nil {
		return nil, gerror.Wrapf(err,
			`tls.LoadX509KeyPair failed for certFile "%s" and keyFile "%s"`,
			crtPath, keyPath,
		)
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}

// MustGetFreePort performs as GetFreePort, but it panics is any error occurs.
func MustGetFreePort() int {
	port, err := GetFreePort()
	if err != nil {
		panic(err)
	}
	return port
}

// GetFreePort retrieves and returns a port that is free.
func GetFreePort() (port int, err error) {
	var (
		network = `tcp`
		address = `:0`
	)
	resolvedAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ResolveTCPAddr failed for network "%s", address "%s"`,
			network, address,
		)
	}
	l, err := net.ListenTCP(network, resolvedAddr)
	if err != nil {
		return 0, gerror.Wrapf(
			err,
			`net.ListenTCP failed for network "%s", address "%s"`,
			network, resolvedAddr.String(),
		)
	}
	port = l.Addr().(*net.TCPAddr).Port
	_ = l.Close()
	return
}

// GetFreePorts retrieves and returns specified number of ports that are free.
func GetFreePorts(count int) (ports []int, err error) {
	var (
		network = `tcp`
		address = `:0`
	)
	for i := 0; i < count; i++ {
		resolvedAddr, err := net.ResolveTCPAddr(network, address)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ResolveTCPAddr failed for network "%s", address "%s"`,
				network, address,
			)
		}
		l, err := net.ListenTCP(network, resolvedAddr)
		if err != nil {
			return nil, gerror.Wrapf(
				err,
				`net.ListenTCP failed for network "%s", address "%s"`,
				network, resolvedAddr.String(),
			)
		}
		ports = append(ports, l.Addr().(*net.TCPAddr).Port)
		_ = l.Close()
	}
	return ports, nil
}
