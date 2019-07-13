// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"time"

	"github.com/gogf/gf/g/errors/gerror"
)

// 简单协议: (方法覆盖)发送数据
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

// 简单协议: (方法覆盖)接收数据
func (c *PoolConn) RecvPkg(option ...PkgOption) ([]byte, error) {
	data, err := c.Conn.RecvPkg(option...)
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return data, err
}

// 简单协议: (方法覆盖)带超时时间的数据获取
func (c *PoolConn) RecvPkgWithTimeout(timeout time.Duration, option ...PkgOption) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer func() {
		err = gerror.Wrap(c.SetRecvDeadline(time.Time{}), "SetRecvDeadline error")
	}()
	data, err = c.RecvPkg(option...)
	return
}

// 简单协议: (方法覆盖)带超时时间的数据发送
func (c *PoolConn) SendPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer func() {
		err = gerror.Wrap(c.SetSendDeadline(time.Time{}), "SetSendDeadline error")
	}()
	err = c.SendPkg(data, option...)
	return
}

// 简单协议: (方法覆盖)发送数据并等待接收返回数据
func (c *PoolConn) SendRecvPkg(data []byte, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkg(option...)
	} else {
		return nil, err
	}
}

// 简单协议: (方法覆盖)发送数据并等待接收返回数据(带返回超时等待时间)
func (c *PoolConn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, option ...PkgOption) ([]byte, error) {
	if err := c.SendPkg(data, option...); err == nil {
		return c.RecvPkgWithTimeout(timeout, option...)
	} else {
		return nil, err
	}
}
