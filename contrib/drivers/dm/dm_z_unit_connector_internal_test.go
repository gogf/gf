// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql/driver"
	"errors"
	"net"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// dialAsUnderlyingDriver dials as the underlying driver does for the data source name `dsn`:
// through the dial function of the login slot that its dialName option names.
func dialAsUnderlyingDriver(dsn string) (net.Conn, error) {
	addr, query, _ := strings.Cut(strings.TrimPrefix(dsn, "dm://"), "?")
	if i := strings.LastIndexByte(addr, '@'); i >= 0 {
		addr = addr[i+1:]
	}
	addr, _, _ = strings.Cut(addr, "/")
	options, _ := url.ParseQuery(query)
	slot, err := strconv.Atoi(strings.TrimPrefix(options.Get("dialName"), dialNamePrefix))
	if err != nil {
		return nil, err
	}
	return dialSlot(slot, addr)
}

// stubNetDial makes the logins dial a connection that the server closes at once, except for the
// address `unreachable`, which answers no dial, until the test ends.
func stubNetDial(t *testing.T, unreachable string) {
	dial := netDialContext
	t.Cleanup(func() {
		netDialContext = dial
	})
	netDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
		if addr == unreachable {
			<-ctx.Done()
			return nil, ctx.Err()
		}
		client, server := net.Pipe()
		_ = server.Close()
		return client, nil
	}
}

// recordingDriver is a driver.DriverContext recording the connectors it opens, and whether the
// turn of loginSemaphore is taken when it opens them.
type recordingDriver struct {
	connectors []*recordingConnector
	turnTaken  []bool
}

func (d *recordingDriver) Open(string) (driver.Conn, error) {
	return nil, driver.ErrSkip
}

func (d *recordingDriver) OpenConnector(dsn string) (driver.Connector, error) {
	c := &recordingConnector{dsn: dsn, driver: d}
	d.connectors = append(d.connectors, c)
	d.turnTaken = append(d.turnTaken, len(loginSemaphore) == 1)
	return c, nil
}

// recordingConnector is a driver.Connector counting the connections it opens.
type recordingConnector struct {
	dsn         string
	driver      *recordingDriver
	connections int
}

func (c *recordingConnector) Connect(context.Context) (driver.Conn, error) {
	c.connections++
	return nil, nil
}

func (c *recordingConnector) Driver() driver.Driver {
	return c.driver
}

func Test_connector_ConnectorPerConnection(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		underlying := &recordingDriver{}
		c := newConnector("dm://u:p@h:5236", underlying)
		for i := 0; i < 3; i++ {
			_, err := c.Connect(context.Background())
			t.AssertNil(err)
		}
		t.Assert(len(underlying.connectors), 3)
		for _, opened := range underlying.connectors {
			t.Assert(strings.HasPrefix(opened.dsn, "dm://u:p@h:5236?dialName="+dialNamePrefix), true)
			t.Assert(opened.connections, 1)
		}
		t.Assert(underlying.turnTaken, []bool{true, true, true})
		t.Assert(len(loginSemaphore), 0)
		t.Assert(c.Driver(), underlying)
	})
}

// blockingDriver is a driver.DriverContext and the driver.Connector it opens, whose Connect dials
// as the underlying driver does, reports that it is entered and then blocks until release is
// closed. OpenConnector fails with err if it is not nil.
type blockingDriver struct {
	dsn     string
	entered chan struct{}
	release chan struct{}
	err     error
}

func (d *blockingDriver) Open(string) (driver.Conn, error) {
	return nil, driver.ErrSkip
}

func (d *blockingDriver) OpenConnector(dsn string) (driver.Connector, error) {
	if d.err != nil {
		return nil, d.err
	}
	return &blockingDriver{dsn: dsn, entered: d.entered, release: d.release}, nil
}

func (d *blockingDriver) Connect(context.Context) (driver.Conn, error) {
	conn, err := dialAsUnderlyingDriver(d.dsn)
	if err != nil {
		return nil, err
	}
	_ = conn.Close()
	d.entered <- struct{}{}
	<-d.release
	return nil, nil
}

func (d *blockingDriver) Driver() driver.Driver {
	return d
}

// Test_connector_OneLoginAtATime tests that a connection logs in once the connection logging in
// has done so, and that a connection waiting for its turn returns once its context is done,
// without taking the turn.
func Test_connector_OneLoginAtATime(t *testing.T) {
	stubNetDial(t, "")
	gtest.C(t, func(t *gtest.T) {
		var (
			underlying = &blockingDriver{entered: make(chan struct{}), release: make(chan struct{})}
			c          = newConnector("dm://u:p@h:5236", underlying)
			first      = make(chan error, 1)
			next       = make(chan error, 1)
		)
		go func() {
			_, err := c.Connect(context.Background())
			first <- err
		}()
		<-underlying.entered

		waitCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := c.Connect(waitCtx)
		t.Assert(err, context.DeadlineExceeded)

		go func() {
			_, err := c.Connect(context.Background())
			next <- err
		}()
		select {
		case <-underlying.entered:
			t.Error("a connection logged in while another was logging in")
		case <-time.After(100 * time.Millisecond):
		}

		close(underlying.release)
		t.AssertNil(<-first)
		<-underlying.entered
		t.AssertNil(<-next)
		t.Assert(len(loginSemaphore), 0)
	})
}

// Test_connector_OpenConnectorError tests that a connection failing to open its connector
// returns the error and gives the turn to the next connection.
func Test_connector_OpenConnectorError(t *testing.T) {
	stubNetDial(t, "")
	gtest.C(t, func(t *gtest.T) {
		var (
			openErr    = errors.New("open connector")
			underlying = &blockingDriver{entered: make(chan struct{}, 1), release: make(chan struct{}), err: openErr}
			c          = newConnector("dm://u:p@h:5236", underlying)
		)
		_, err := c.Connect(context.Background())
		t.Assert(err, openErr)

		underlying.err = nil
		close(underlying.release)
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		defer cancel()
		_, err = c.Connect(ctx)
		t.AssertNil(err)
	})
}

// Test_connector_DialOutsideTurn tests that a connection dialing a server that does not answer
// does not hold up the login of another connection, and that its dial ends with its timeout.
func Test_connector_DialOutsideTurn(t *testing.T) {
	stubNetDial(t, "unreachable:5236")
	gtest.C(t, func(t *gtest.T) {
		var (
			underlying  = &blockingDriver{entered: make(chan struct{}, 2), release: make(chan struct{})}
			unreachable = newConnector("dm://u:p@unreachable:5236?connectTimeout=300", underlying)
			reachable   = newConnector("dm://u:p@reachable:5236", underlying)
			failed      = make(chan error, 1)
			start       = time.Now()
		)
		close(underlying.release)
		go func() {
			_, err := unreachable.Connect(context.Background())
			failed <- err
		}()
		time.Sleep(50 * time.Millisecond)
		_, err := reachable.Connect(context.Background())
		t.AssertNil(err)
		t.AssertLT(time.Since(start), 300*time.Millisecond)

		t.Assert(<-failed, context.DeadlineExceeded)
		t.AssertGE(time.Since(start), 300*time.Millisecond)
		t.Assert(len(underlying.entered), 1)
		t.Assert(len(loginSemaphore), 0)
	})
}

// Test_connector_KeepTurn tests that a connection keeping the turn while it dials, as it does with
// the dialName, driverReconnect or rwSeparate option, logs in once no other connection does, with
// the data source name as it is.
func Test_connector_KeepTurn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		for _, dsn := range []string{
			"dm://u:p@h:5236?dialName=mine",
			"dm://u:p@h:5236?driverReconnect=true",
			"dm://u:p@h:5236?rwSeparate=1",
		} {
			var (
				underlying = &recordingDriver{}
				c          = newConnector(dsn, underlying)
			)
			loginSemaphore <- struct{}{}
			waitCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
			_, err := c.Connect(waitCtx)
			cancel()
			t.Assert(err, context.DeadlineExceeded)
			t.Assert(len(underlying.connectors), 0)

			<-loginSemaphore
			_, err = c.Connect(context.Background())
			t.AssertNil(err)
			t.Assert(len(underlying.connectors), 1)
			t.Assert(underlying.connectors[0].dsn, dsn)
			t.Assert(len(loginSemaphore), 0)
		}
	})
}

// Test_newConnector tests the dial timeout of the connector of a data source name, and whether its
// connections keep the turn while they dial.
func Test_newConnector(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			dsn         string
			dialTimeout time.Duration
			keepTurn    bool
		}{
			{"dm://u:p@h:5236/S?schema=S", defaultConnectTimeout, false},
			{"dm://u:p@h:5236", defaultConnectTimeout, false},
			{"dm://u:p@h:5236/S?connectTimeout=1500", 1500 * time.Millisecond, false},
			{"dm://u:p@h:5236/S?CONNECTTIMEOUT=0", 0, false},
			{"dm://u:p@h:5236/S?connectTimeout=x", defaultConnectTimeout, false},
			{"dm://u:p@h:5236/S?dialName=mine", defaultConnectTimeout, true},
			{"dm://u:p@h:5236/S?driverReconnect=1", defaultConnectTimeout, true},
			{"dm://u:p@h:5236/S?driverReconnect=false", defaultConnectTimeout, false},
			{"dm://u:p@h:5236/S?rwSeparate=2", defaultConnectTimeout, true},
			{"dm://u:p@h:5236/S?rwSeparate=0", defaultConnectTimeout, false},
		}
		for _, c := range cases {
			opened := newConnector(c.dsn, &recordingDriver{})
			t.Assert(opened.dsn, c.dsn)
			t.Assert(opened.dialTimeout, c.dialTimeout)
			t.Assert(opened.keepTurn, c.keepTurn)
		}
	})
}

// Test_loginSlots tests that the logins running at the same time take different slots, and that a
// login takes a slot that an ended login has left rather than a new one.
func Test_loginSlots(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			first  = beginLogin(&login{ctx: context.Background()})
			second = beginLogin(&login{ctx: context.Background()})
		)
		t.AssertNE(first, second)
		endLogin(first)
		endLogin(second)

		loginSlots.Lock()
		slots := len(loginSlots.logins)
		loginSlots.Unlock()
		third, fourth := beginLogin(&login{ctx: context.Background()}), beginLogin(&login{ctx: context.Background()})
		t.Assert(map[int]bool{third: true, fourth: true}, map[int]bool{first: true, second: true})
		endLogin(third)
		endLogin(fourth)
		loginSlots.Lock()
		t.Assert(len(loginSlots.logins), slots)
		loginSlots.Unlock()
	})
}
