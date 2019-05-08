// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
    "time"
    "github.com/gogf/gf/g/container/gmap"
    "github.com/gogf/gf/g/container/gpool"
)

// 链接池链接对象
type PoolConn struct {
    *Conn              // 继承底层链接接口对象
    pool   *gpool.Pool // 对应的链接池对象
    status int         // 当前对象的状态，主要用于失败重连判断
}

const (
    gDEFAULT_POOL_EXPIRE = 60000 // (毫秒)默认链接对象过期时间
    gCONN_STATUS_UNKNOWN = 0     // 未知，表示未经过连通性操作;
    gCONN_STATUS_ACTIVE  = 1     // 正常，表示已经经过连通性操作
    gCONN_STATUS_ERROR   = 2     // 错误，表示该接口操作产生了错误，不应当被循环使用了

)

var (
    // 连接池对象map，键名为地址端口，键值为对应的连接池对象
    pools = gmap.NewStrAnyMap()
)

// 创建TCP链接池对象
func NewPoolConn(addr string, timeout...int) (*PoolConn, error) {
    var pool *gpool.Pool
    if v := pools.Get(addr); v == nil {
        pools.LockFunc(func(m map[string]interface{}) {
            if v, ok := m[addr]; ok {
                pool = v.(*gpool.Pool)
            } else {
                pool = gpool.New(gDEFAULT_POOL_EXPIRE, func() (interface{}, error) {
                    if conn, err := NewConn(addr, timeout...); err == nil {
                        return &PoolConn { conn, pool, gCONN_STATUS_ACTIVE }, nil
                    } else {
                        return nil, err
                    }
                })
                m[addr] = pool
            }
        })
    } else {
        pool = v.(*gpool.Pool)
    }

    if v, err := pool.Get(); err == nil {
        return v.(*PoolConn), nil
    } else {
        return nil, err
    }
}

// (方法覆盖)覆盖底层接口对象的Close方法
func (c *PoolConn) Close() error {
    if c.pool != nil && c.status == gCONN_STATUS_ACTIVE {
        c.status = gCONN_STATUS_UNKNOWN
        c.pool.Put(c)
    } else {
        c.Conn.Close()
    }
    return nil
}

// (方法覆盖)发送数据
func (c *PoolConn) Send(data []byte, retry...Retry) error {
    var err error
    if err = c.Conn.Send(data, retry...); err != nil && c.status == gCONN_STATUS_UNKNOWN {
        if v, e := c.pool.NewFunc(); e == nil {
            c.Conn = v.(*PoolConn).Conn
            err    = c.Conn.Send(data, retry...)
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

// (方法覆盖)接收数据
func (c *PoolConn) Recv(length int, retry...Retry) ([]byte, error) {
    data, err := c.Conn.Recv(length, retry...)
    if err != nil {
        c.status = gCONN_STATUS_ERROR
    } else {
        c.status = gCONN_STATUS_ACTIVE
    }
    return data, err
}

// (方法覆盖)按行读取数据，阻塞读取，直到完成一行读取位置(末尾以'\n'结尾，返回数据不包含换行符)
func (c *PoolConn) RecvLine(retry...Retry) ([]byte, error) {
    data, err := c.Conn.RecvLine(retry...)
    if err != nil {
        c.status = gCONN_STATUS_ERROR
    } else {
        c.status = gCONN_STATUS_ACTIVE
    }
    return data, err
}

// (方法覆盖)带超时时间的数据获取
func (c *PoolConn) RecvWithTimeout(length int, timeout time.Duration, retry...Retry) ([]byte, error) {
    c.SetRecvDeadline(time.Now().Add(timeout))
    defer c.SetRecvDeadline(time.Time{})
    return c.Recv(length, retry...)
}

// (方法覆盖)带超时时间的数据发送
func (c *PoolConn) SendWithTimeout(data []byte, timeout time.Duration, retry...Retry) error {
    c.SetSendDeadline(time.Now().Add(timeout))
    defer c.SetSendDeadline(time.Time{})
    return c.Send(data, retry...)
}

// (方法覆盖)发送数据并等待接收返回数据
func (c *PoolConn) SendRecv(data []byte, receive int, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.Recv(receive, retry...)
    } else {
        return nil, err
    }
}

// (方法覆盖)发送数据并等待接收返回数据(带返回超时等待时间)
func (c *PoolConn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry...Retry) ([]byte, error) {
    if err := c.Send(data, retry...); err == nil {
        return c.RecvWithTimeout(receive, timeout, retry...)
    } else {
        return nil, err
    }
}