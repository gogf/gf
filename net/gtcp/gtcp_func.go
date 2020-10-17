// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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

	"github.com/gogf/gf/os/gfile"
)

const (
	gDEFAULT_CONN_TIMEOUT     = 30 * time.Second       // Default connection timeout.
	gDEFAULT_RETRY_INTERVAL   = 100 * time.Millisecond // Default retry interval.
	gDEFAULT_READ_BUFFER_SIZE = 128                    // (Byte) Buffer size for reading.
)

type Retry struct {
	Count    int           // Retry count.
	Interval time.Duration // Retry interval.
}

// NewNetConn creates and returns a net.Conn with given address like "127.0.0.1:80".
// The optional parameter <timeout> specifies the timeout for dialing connection.
func NewNetConn(addr string, timeout ...time.Duration) (net.Conn, error) {
	d := gDEFAULT_CONN_TIMEOUT
	if len(timeout) > 0 {
		d = timeout[0]
	}
	return net.DialTimeout("tcp", addr, d)
}

// NewNetConnTLS creates and returns a TLS net.Conn with given address like "127.0.0.1:80".
// The optional parameter <timeout> specifies the timeout for dialing connection.
func NewNetConnTLS(addr string, tlsConfig *tls.Config, timeout ...time.Duration) (net.Conn, error) {
	dialer := &net.Dialer{
		Timeout: gDEFAULT_CONN_TIMEOUT,
	}
	if len(timeout) > 0 {
		dialer.Timeout = timeout[0]
	}
	return tls.DialWithDialer(dialer, "tcp", addr, tlsConfig)
}

// NewNetConnKeyCrt creates and returns a TLS net.Conn with given TLS certificate and key files
// and address like "127.0.0.1:80". The optional parameter <timeout> specifies the timeout for
// dialing connection.
func NewNetConnKeyCrt(addr, crtFile, keyFile string, timeout ...time.Duration) (net.Conn, error) {
	tlsConfig, err := LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		return nil, err
	}
	return NewNetConnTLS(addr, tlsConfig, timeout...)
}

// Send creates connection to <address>, writes <data> to the connection and then closes the connection.
// The optional parameter <retry> specifies the retry policy when fails in writing data.
func Send(address string, data []byte, retry ...Retry) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.Send(data, retry...)
}

// SendRecv creates connection to <address>, writes <data> to the connection, receives response
// and then closes the connection.
//
// The parameter <length> specifies the bytes count waiting to receive. It receives all buffer content
// and returns if <length> is -1.
//
// The optional parameter <retry> specifies the retry policy when fails in writing data.
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

// isTimeout checks whether given <err> is a timeout error.
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
		return nil, err
	}
	tlsConfig := &tls.Config{}
	tlsConfig.Certificates = []tls.Certificate{crt}
	tlsConfig.Time = time.Now
	tlsConfig.Rand = rand.Reader
	return tlsConfig, nil
}
