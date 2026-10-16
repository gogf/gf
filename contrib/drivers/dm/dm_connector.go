// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql/driver"
)

// loginSemaphore lets one connection log in at a time. The underlying driver keeps the
// Diffie-Hellman parameters of the login in package variables that every login writes, so that
// connections logging in at the same time race.
var loginSemaphore = make(chan struct{}, 1)

// connector is the driver.Connector of a database pool, which opens each connection with a
// connector of the underlying driver of its own, and one at a time, see loginSemaphore. The
// underlying driver writes the settings that the server returns at login, like whether it
// compresses messages, to the connector that opens the connection, and every connection of that
// connector reads them, so that connections sharing a connector race.
type connector struct {
	dsn    string
	driver driver.DriverContext
}

// Connect opens a connection with a connector of its own, once no other connection logs in.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	select {
	case loginSemaphore <- struct{}{}:
	case <-ctx.Done():
		return nil, ctx.Err()
	}
	defer func() { <-loginSemaphore }()
	underlyingConnector, err := c.driver.OpenConnector(c.dsn)
	if err != nil {
		return nil, err
	}
	return underlyingConnector.Connect(ctx)
}

// Driver returns the underlying driver.
func (c *connector) Driver() driver.Driver {
	return c.driver.(driver.Driver)
}
