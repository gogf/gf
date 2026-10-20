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
	"sync"
	"time"

	"gitee.com/chunanyong/dm"
)

const (
	// dialNamePrefix is the prefix of the names of the dial functions of the login slots, see
	// loginSlots.
	dialNamePrefix = "gogf_gf_dm_login_"
	// defaultConnectTimeout is the timeout of a dial without the connectTimeout option, which is
	// the default of the option in the underlying driver.
	defaultConnectTimeout = 5 * time.Second
)

// loginSemaphore lets one connection log in at a time. The underlying driver keeps the
// Diffie-Hellman parameters of the login in package variables that every login writes, so that
// connections logging in at the same time race.
var loginSemaphore = make(chan struct{}, 1)

// netDialContext dials the server for a login, with the TCP keep-alive period of the underlying
// driver.
var netDialContext = (&net.Dialer{KeepAlive: 2 * time.Second}).DialContext

// loginSlots holds the logins in progress, each in a slot whose dial function, registered with the
// underlying driver under the name of the slot, dials for that login, as the underlying driver
// passes a dial function the address only, with a context of its own. A slot that a login leaves
// is taken by the next one, so that there are as many slots as logins running at the same time.
var loginSlots struct {
	sync.Mutex
	logins []*login
	free   []int
}

// login is the login of a connection, which holds the turn of loginSemaphore but while it dials:
// the context and the timeout of its dial.
type login struct {
	ctx         context.Context
	dialTimeout time.Duration
}

// dialName returns the name of the dial function of the login slot `slot`.
func dialName(slot int) string {
	return dialNamePrefix + strconv.Itoa(slot)
}

// beginLogin puts the login `l` in a free login slot, registering the dial function of a new slot
// if none is free, and returns the slot.
func beginLogin(l *login) int {
	loginSlots.Lock()
	defer loginSlots.Unlock()
	var slot int
	if n := len(loginSlots.free); n > 0 {
		slot = loginSlots.free[n-1]
		loginSlots.free = loginSlots.free[:n-1]
	} else {
		slot = len(loginSlots.logins)
		loginSlots.logins = append(loginSlots.logins, nil)
		dm.RegisterDialContext(dialName(slot), func(_ context.Context, addr string) (net.Conn, error) {
			return dialSlot(slot, addr)
		})
	}
	loginSlots.logins[slot] = l
	return slot
}

// endLogin frees the login slot `slot`.
func endLogin(slot int) {
	loginSlots.Lock()
	defer loginSlots.Unlock()
	loginSlots.logins[slot] = nil
	loginSlots.free = append(loginSlots.free, slot)
}

// dialSlot dials `addr` for the login of the login slot `slot` within its dial timeout, leaving the
// turn of the login to the other connections while it dials, so that a server that cannot be
// reached holds up its own login only. It takes the turn back once the dial ends, before the
// underlying driver goes on. Without a login in the slot, it dials within the default timeout.
func dialSlot(slot int, addr string) (net.Conn, error) {
	loginSlots.Lock()
	l := loginSlots.logins[slot]
	loginSlots.Unlock()
	ctx, timeout := context.Background(), defaultConnectTimeout
	if l != nil {
		ctx, timeout = l.ctx, l.dialTimeout
	}
	if timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	if l != nil {
		<-loginSemaphore
		defer func() {
			loginSemaphore <- struct{}{}
		}()
	}
	return netDialContext(ctx, "tcp", addr)
}

// takeLoginTurn takes the turn of loginSemaphore, waiting until no other connection logs in or
// `ctx` is done.
func takeLoginTurn(ctx context.Context) error {
	select {
	case loginSemaphore <- struct{}{}:
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
// A connection opens its connector and logs in once no other connection does, see loginSemaphore,
// as the underlying driver also writes package variables when it opens a connector, which its
// logins read. It dials the server through the dial function of a login slot, which leaves the
// turn to the other connections while it dials. A data source name naming a dial function of its
// own with the dialName option keeps it, and with the driverReconnect or rwSeparate option, under
// which the underlying driver dials outside the login of a connection, the underlying driver dials
// as it does; their connections keep the turn while they dial.
type connector struct {
	dsn         string
	driver      driver.DriverContext
	dialTimeout time.Duration
	keepTurn    bool
	// wallClock is the configured Timezone, or else the local time zone, in which a connection's
	// session opens with the offset at the moment, unless the timeZone option of the data source
	// name overrides it; see Connect and sessionTimes.
	wallClock *time.Location
}

// newConnector returns the connector of the data source name `dsn` for the underlying driver
// `d`. Its connections dial within the timeout of the connectTimeout option, in milliseconds.
func newConnector(dsn string, d driver.DriverContext) *connector {
	c := &connector{dsn: dsn, driver: d, dialTimeout: defaultConnectTimeout, wallClock: time.Local}
	_, query, _ := strings.Cut(dsn, "?")
	options, _ := url.ParseQuery(query)
	for key, values := range options {
		var value string
		if len(values) > 0 {
			value = values[0]
		}
		switch strings.ToLower(key) {
		case "dialname":
			c.keepTurn = true
		case "driverreconnect":
			if reconnect, err := strconv.ParseBool(value); err == nil && reconnect {
				c.keepTurn = true
			}
		case "rwseparate":
			if mode, err := strconv.Atoi(value); err == nil && mode > 0 {
				c.keepTurn = true
			}
		case "connecttimeout":
			if milliseconds, err := strconv.Atoi(value); err == nil && milliseconds >= 0 {
				c.dialTimeout = time.Duration(milliseconds) * time.Millisecond
			}
		}
	}
	return c
}

// Connect opens a connection with a connector of its own, see connector, and wraps it so that it
// binds a time argument by the time zone of its own session, see sessionTimes.
func (c *connector) Connect(ctx context.Context) (driver.Conn, error) {
	if err := takeLoginTurn(ctx); err != nil {
		return nil, err
	}
	defer func() {
		<-loginSemaphore
	}()
	dsn := c.dsn
	minutes, timeZoneGiven := timeZoneOption(c.dsn)
	if !c.keepTurn {
		slot := beginLogin(&login{ctx: ctx, dialTimeout: c.dialTimeout})
		defer endLogin(slot)
		dsn = appendOption(dsn, "dialName="+dialName(slot))
	}
	if !timeZoneGiven {
		minutes = offsetMinutes(c.wallClock, time.Now())
		dsn = appendOption(dsn, "timeZone="+strconv.Itoa(minutes))
	}
	underlyingConnector, err := c.driver.OpenConnector(dsn)
	if err != nil {
		return nil, err
	}
	conn, err := underlyingConnector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	session := time.FixedZone("", minutes*60)
	return newTimedConn(conn, sessionTimes{wallClock: c.wallClock, session: session}), nil
}

// appendOption returns the data source name `dsn` with the option `option`.
func appendOption(dsn, option string) string {
	if strings.Contains(dsn, "?") {
		return dsn + "&" + option
	}
	return dsn + "?" + option
}

// Driver returns the underlying driver.
func (c *connector) Driver() driver.Driver {
	return c.driver.(driver.Driver)
}
