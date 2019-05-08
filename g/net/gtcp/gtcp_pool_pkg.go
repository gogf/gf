// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"time"
)

// 简单协议: (方法覆盖)发送数据
func (c *PoolConn) SendPkg(data []byte, retry...Retry) (err error) {
    if err = c.Conn.SendPkg(data, retry...); err != nil && c.status == gCONN_STATUS_UNKNOWN {
        if v, e := c.pool.NewFunc(); e == nil {
            c.Conn = v.(*PoolConn).Conn
            err    = c.Conn.SendPkg(data, retry...)
        } else {
            err    = e
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
func (c *PoolConn) RecvPkg(retry...Retry) ([]byte, error) {
    data, err := c.Conn.RecvPkg(retry...)
    if err != nil {
        c.status = gCONN_STATUS_ERROR
    } else {
        c.status = gCONN_STATUS_ACTIVE
    }
    return data, err
}

// 简单协议: (方法覆盖)带超时时间的数据获取
func (c *PoolConn) RecvPkgWithTimeout(timeout time.Duration, retry...Retry) ([]byte, error) {
    c.SetRecvDeadline(time.Now().Add(timeout))
    defer c.SetRecvDeadline(time.Time{})
    return c.RecvPkg(retry...)
}

// 简单协议: (方法覆盖)带超时时间的数据发送
func (c *PoolConn) SendPkgWithTimeout(data []byte, timeout time.Duration, retry...Retry) error {
    c.SetSendDeadline(time.Now().Add(timeout))
    defer c.SetSendDeadline(time.Time{})
    return c.SendPkg(data, retry...)
}

// 简单协议: (方法覆盖)发送数据并等待接收返回数据
func (c *PoolConn) SendRecvPkg(data []byte, retry...Retry) ([]byte, error) {
    if err := c.SendPkg(data, retry...); err == nil {
        return c.RecvPkg(retry...)
    } else {
        return nil, err
    }
}

// 简单协议: (方法覆盖)发送数据并等待接收返回数据(带返回超时等待时间)
func (c *PoolConn) SendRecvPkgWithTimeout(data []byte, timeout time.Duration, retry...Retry) ([]byte, error) {
    if err := c.SendPkg(data, retry...); err == nil {
        return c.RecvPkgWithTimeout(timeout, retry...)
    } else {
        return nil, err
    }
}