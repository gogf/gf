// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import "time"

// SendPkg sends a package containing <data> to <address> and closes the connection.
// The optional parameter <option> specifies the package options for sending.
func SendPkg(address string, data []byte, option ...PkgOption) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.SendPkg(data, option...)
}

// SendRecvPkg sends a package containing <data> to <address>, receives the response
// and closes the connection. The optional parameter <option> specifies the package options for sending.
func SendRecvPkg(address string, data []byte, option ...PkgOption) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecvPkg(data, option...)
}

// SendPkgWithTimeout sends a package containing <data> to <address> with timeout limitation
// and closes the connection. The optional parameter <option> specifies the package options for sending.
func SendPkgWithTimeout(address string, data []byte, timeout time.Duration, option ...PkgOption) error {
	conn, err := NewConn(address)
	if err != nil {
		return err
	}
	defer conn.Close()
	return conn.SendPkgWithTimeout(data, timeout, option...)
}

// SendRecvPkgWithTimeout sends a package containing <data> to <address>, receives the response with timeout limitation
// and closes the connection. The optional parameter <option> specifies the package options for sending.
func SendRecvPkgWithTimeout(address string, data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error) {
	conn, err := NewConn(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	return conn.SendRecvPkgWithTimeout(data, timeout, option...)
}
