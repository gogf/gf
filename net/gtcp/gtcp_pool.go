// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtcp

import (
	"time"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gpool"
)

// PoolConn is a connection with pool feature for TCP.
// Note that it is NOT a pool or connection manager,
// it is just a TCP connection object.
type PoolConn struct {
	*Conn              // Underlying connection object.
	pool   *gpool.Pool // Connection pool, which is not a really connection pool, but a connection reusable pool.
	status int         // Status of current connection, which is used to mark this connection usable or not.
}

const (
	gDEFAULT_POOL_EXPIRE = 10 * time.Second // Default TTL for connection in the pool.
	gCONN_STATUS_UNKNOWN = 0                // Means it is unknown it's connective or not.
	gCONN_STATUS_ACTIVE  = 1                // Means it is now connective.
	gCONN_STATUS_ERROR   = 2                // Means it should be closed and removed from pool.
)

var (
	// addressPoolMap is a mapping for address to its pool object.
	addressPoolMap = gmap.NewStrAnyMap(true)
)

// NewPoolConn creates and returns a connection with pool feature.
func NewPoolConn(addr string, timeout ...time.Duration) (*PoolConn, error) {
	v := addressPoolMap.GetOrSetFuncLock(addr, func() interface{} {
		var pool *gpool.Pool
		pool = gpool.New(gDEFAULT_POOL_EXPIRE, func() (interface{}, error) {
			if conn, err := NewConn(addr, timeout...); err == nil {
				return &PoolConn{conn, pool, gCONN_STATUS_ACTIVE}, nil
			} else {
				return nil, err
			}
		})
		return pool
	})
	if v, err := v.(*gpool.Pool).Get(); err == nil {
		return v.(*PoolConn), nil
	} else {
		return nil, err
	}
}

// Close puts back the connection to the pool if it's active,
// or closes the connection if it's not active.
//
// Note that, if <c> calls Close function closing itself, <c> can not
// be used again.
func (c *PoolConn) Close() error {
	if c.pool != nil && c.status == gCONN_STATUS_ACTIVE {
		c.status = gCONN_STATUS_UNKNOWN
		c.pool.Put(c)
	} else {
		return c.Conn.Close()
	}
	return nil
}

// Send writes data to the connection. It retrieves a new connection from its pool if it fails
// writing data.
func (c *PoolConn) Send(data []byte, retry ...Retry) error {
	err := c.Conn.Send(data, retry...)
	if err != nil && c.status == gCONN_STATUS_UNKNOWN {
		if v, e := c.pool.Get(); e == nil {
			c.Conn = v.(*PoolConn).Conn
			err = c.Send(data, retry...)
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

// Recv receives data from the connection.
func (c *PoolConn) Recv(length int, retry ...Retry) ([]byte, error) {
	data, err := c.Conn.Recv(length, retry...)
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return data, err
}

// RecvLine reads data from the connection until reads char '\n'.
// Note that the returned result does not contain the last char '\n'.
func (c *PoolConn) RecvLine(retry ...Retry) ([]byte, error) {
	data, err := c.Conn.RecvLine(retry...)
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return data, err
}

// RecvTil reads data from the connection until reads bytes <til>.
// Note that the returned result contains the last bytes <til>.
func (c *PoolConn) RecvTil(til []byte, retry ...Retry) ([]byte, error) {
	data, err := c.Conn.RecvTil(til, retry...)
	if err != nil {
		c.status = gCONN_STATUS_ERROR
	} else {
		c.status = gCONN_STATUS_ACTIVE
	}
	return data, err
}

// RecvWithTimeout reads data from the connection with timeout.
func (c *PoolConn) RecvWithTimeout(length int, timeout time.Duration, retry ...Retry) (data []byte, err error) {
	if err := c.SetRecvDeadline(time.Now().Add(timeout)); err != nil {
		return nil, err
	}
	defer c.SetRecvDeadline(time.Time{})
	data, err = c.Recv(length, retry...)
	return
}

// SendWithTimeout writes data to the connection with timeout.
func (c *PoolConn) SendWithTimeout(data []byte, timeout time.Duration, retry ...Retry) (err error) {
	if err := c.SetSendDeadline(time.Now().Add(timeout)); err != nil {
		return err
	}
	defer c.SetSendDeadline(time.Time{})
	err = c.Send(data, retry...)
	return
}

// SendRecv writes data to the connection and blocks reading response.
func (c *PoolConn) SendRecv(data []byte, receive int, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.Recv(receive, retry...)
	} else {
		return nil, err
	}
}

// SendRecvWithTimeout writes data to the connection and reads response with timeout.
func (c *PoolConn) SendRecvWithTimeout(data []byte, receive int, timeout time.Duration, retry ...Retry) ([]byte, error) {
	if err := c.Send(data, retry...); err == nil {
		return c.RecvWithTimeout(receive, timeout, retry...)
	} else {
		return nil, err
	}
}
