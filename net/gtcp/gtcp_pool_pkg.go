// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"time"
)

// SendPkg sends a package containing <data> to the connection.
// The optional parameter <option> specifies the package options for sending.
func (c *PoolConn) SendPkg(data []byte, option ...PkgOption) (err error) {
	if err = c.Conn.SendPkg(data, option...); err != nil && c.status == gCONN_STATUS_UNKNOWN {
		if v, e := c.pool.NewFunc(); e == nil {
			c.Conn = v.(*PoolConn).Conn
			err = c.Conn.SendPkg(data, option...)
		} else {
			err = e
		}
	}
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return err
}

// RecvPkg receives package from connection using simple package protocol.
// The optional parameter <option> specifies the package options for receiving.
func (c *PoolConn) RecvPkg(option ...PkgOption) ([]byte, error) {
	data, err := c.Conn.RecvPkg(option...)
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return data, err
}

// RecvPkgWithTimeout reads data from connection with timeout using simple package protocol.
func (c *PoolConn) RecvPkgWithTimeout(timeout time.Duration, option ...PkgOption) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	data, err = c.RecvPkg(option...)
	return
}

// SendPkgWithTimeout writes data to connection with timeout using simple package protocol.
func (c *PoolConn) SendPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	err = c.SendPkg(data, option...)
	return
}

// SendRecvPkg writes data to connection and blocks reading response using simple package protocol.
func (c *PoolConn) SendRecvPkg(data []byte, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkg(option...)
	} else {
		return nil, err
	}
}

// RecvPkgWithTimeout reads data from connection with timeout using simple package protocol.
func (c *PoolConn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkgWithTimeout(timeout, option...)
	} else {
		return nil, err
	}
}
