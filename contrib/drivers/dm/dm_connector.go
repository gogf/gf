// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql/driver"
	"net"
	"net/url"
	"strconv"
	"strings"
	"time"

	"gitee.com/chunanyong/dm"
)

const (
	// dialName is the name of dialContext among the dial functions of the underlying driver.
	dialName = "gogf_gf_dm"
	// defaultConnectTimeout is the timeout of a dial without the connectTimeout option, which is
	// the default of the option in the underlying driver.
	defaultConnectTimeout = 5 * time.Second
)

// loginSemaphore lets one connection log in at a time. The underlying driver keeps the
// Diffie-Hellman parameters of the login in package variables that every login writes, so that
// connections logging in at the same time race.
var loginSemaphore = make(chan struct{}, 1)

// netDialContext dials the server for dialContext, with the TCP keep-alive period of the
// underlying driver.
var netDialContext = (&net.Dialer{KeepAlive: 2 * time.Second}).DialContext

func init() {
	dm.RegisterDialContext(dialName, dialContext)
}

// loginKey is the context key of the login of the connection that the context opens.
type loginKey struct{}

// login is the login of a connection: the timeout of its dial, and whether it holds the turn of
// loginSemaphore.
type login struct {
	dialTimeout time.Duration
	holding     bool
}

// dialContext dials `addr` for the underlying driver within the dial timeout of the login of
// `ctx`, and then takes the turn of the login, so that a server that cannot be reached holds up
// its own login only.
func dialContext(ctx context.Context, addr string) (net.Conn, error) {
	dialCtx := ctx
	if l, ok := ctx.Value(loginKey{}).(*login); ok && l.dialTimeout > 0 {
		var cancel context.CancelFunc
		dialCtx, cancel = context.WithTimeout(ctx, l.dialTimeout)
		defer cancel()
	}
	conn, err := netDialContext(dialCtx, "tcp", addr)
	if err != nil {
		return nil, err
	}
	if err = takeLoginTurn(ctx); err != nil {
		_ = conn.Close()
		return nil, err
	}
	return conn, nil
}

// takeLoginTurn takes the turn of loginSemaphore for the login of `ctx`, unless it holds it
// already, waiting until no other connection logs in or `ctx` is done.
func takeLoginTurn(ctx context.Context) error {
	l, ok := ctx.Value(loginKey{}).(*login)
	if !ok || l.holding {
		return nil
	}
	select {
	case loginSemaphore <- struct{}{}:
		l.holding = true
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// connector is the driver.Connector of a database pool, which opens each connection with a
// connector of the underlying driver of its own. The underlying driver writes the settings that
// the server returns at login, like whether it compresses messages, to the connector that opens
// the connection, and every connection of that connector reads them, so that connections sharing
// a connector race.
//
// A connection dials the server through dialContext, which lets it log in once no other
// connection does, see loginSemaphore. A data source name naming a dial function of its own with
// the dialName option keeps it, and its connections take the turn before they dial.
type connector struct {
	dsn         string
	driver      driver.DriverContext
	dialTimeout time.Duration
	ownDial     bool
}

// newConnector returns the connector of the data source name `dsn` for the underlying driver
// `d`. Its connections dial within the timeout of the connectTimeout option, in milliseconds.
func newConnector(dsn string, d driver.DriverContext) *connector {
	c := &connector{dsn: dsn, driver: d, dialTimeout: defaultConnectTimeout}
	_, query, _ := strings.Cut(dsn, "?")
	options, _ := url.ParseQuery(query)
	for key, values := range options {
		switch {
		case strings.EqualFold(key, "dialName"):
			c.ownDial = true
		case strings.EqualFold(key, "connectTimeout") && len(values) > 0:
			if milliseconds, err := strconv.Atoi(values[0]); err == nil && milliseconds >= 0 {
				c.dialTimeout = time.Duration(milliseconds) * time.Millisecond
			}
		}
	}
	if !c.ownDial {
		separator := "?"
		if strings.Contains(dsn, "?") {
			separator = "&"
		}
		c.dsn += separator + "dialName=" + dialName
	}
	return c
}

// Connect opens a connection with a connector of its own, see connector.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	l := &login{dialTimeout: c.dialTimeout}
	ctx = context.WithValue(ctx, loginKey{}, l)
	defer func() {
		if l.holding {
			<-loginSemaphore
		}
	}()
	if c.ownDial {
		if err := takeLoginTurn(ctx); err != nil {
			return nil, err
		}
	}
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
