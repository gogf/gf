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
	"strings"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"
)

// recordingDriver is a driver.DriverContext recording the connectors it opens.
type recordingDriver struct {
	connectors []*recordingConnector
}

func (d *recordingDriver) Open(string) (driver.Conn, error) {
	return nil, driver.ErrSkip
}

func (d *recordingDriver) OpenConnector(dsn string) (driver.Connector, error) {
	c := &recordingConnector{dsn: dsn, driver: d}
	d.connectors = append(d.connectors, c)
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
		c := &connector{dsn: "dm://u:p@h:5236", driver: underlying}
		for i := 0; i < 3; i++ {
			_, err := c.Connect(context.Background())
			t.AssertNil(err)
		}
		t.Assert(len(underlying.connectors), 3)
		for _, opened := range underlying.connectors {
			t.Assert(opened.dsn, "dm://u:p@h:5236")
			t.Assert(opened.connections, 1)
		}
		t.Assert(c.Driver(), underlying)
	})
}

// blockingDriver is a driver.DriverContext and the driver.Connector it opens, whose Connect
// reports that it is entered and then blocks until release is closed. OpenConnector fails with
// err if it is not nil.
type blockingDriver struct {
	entered chan struct{}
	release chan struct{}
	err     error
}

func (d *blockingDriver) Open(string) (driver.Conn, error) {
	return nil, driver.ErrSkip
}

func (d *blockingDriver) OpenConnector(string) (driver.Connector, error) {
	if d.err != nil {
		return nil, d.err
	}
	return d, nil
}

func (d *blockingDriver) Connect(ctx context.Context) (driver.Conn, error) {
	// The underlying driver logs in once dialContext has dialed and taken the turn.
	if err := takeLoginTurn(ctx); err != nil {
		return nil, err
	}
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
	gtest.C(t, func(t *gtest.T) {
		var (
			underlying = &blockingDriver{entered: make(chan struct{}), release: make(chan struct{})}
			c          = &connector{dsn: "dm://u:p@h:5236", driver: underlying}
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
	})
}

// Test_connector_OpenConnectorError tests that a connection failing to open its connector
// returns the error and gives the turn to the next connection.
func Test_connector_OpenConnectorError(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			openErr    = errors.New("open connector")
			underlying = &blockingDriver{entered: make(chan struct{}, 1), release: make(chan struct{}), err: openErr}
			c          = &connector{dsn: "dm://u:p@h:5236", driver: underlying}
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

// dialingDriver is a driver.DriverContext and the driver.Connector it opens, whose Connect dials
// through dialContext the address of the data source name, as the underlying driver does, and
// reports the connection that dialed.
type dialingDriver struct {
	dsn       string
	connected chan string
}

func (d *dialingDriver) Open(string) (driver.Conn, error) {
	return nil, driver.ErrSkip
}

func (d *dialingDriver) OpenConnector(dsn string) (driver.Connector, error) {
	return &dialingDriver{dsn: dsn, connected: d.connected}, nil
}

func (d *dialingDriver) Connect(ctx context.Context) (driver.Conn, error) {
	addr, _, _ := strings.Cut(strings.TrimPrefix(d.dsn, "dm://"), "?")
	conn, err := dialContext(ctx, addr)
	if err != nil {
		return nil, err
	}
	_ = conn.Close()
	d.connected <- addr
	return nil, nil
}

func (d *dialingDriver) Driver() driver.Driver {
	return d
}

// Test_connector_DialOutsideTurn tests that a connection dialing a server that does not answer
// does not hold up the login of another connection, and that its dial ends with its timeout.
func Test_connector_DialOutsideTurn(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer func(dial func(ctx context.Context, network, addr string) (net.Conn, error)) {
			netDialContext = dial
		}(netDialContext)
		netDialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == "unreachable:5236" {
				<-ctx.Done()
				return nil, ctx.Err()
			}
			client, server := net.Pipe()
			_ = server.Close()
			return client, nil
		}
		var (
			underlying  = &dialingDriver{connected: make(chan string, 2)}
			unreachable = newConnector("dm://unreachable:5236?connectTimeout=300", underlying)
			reachable   = newConnector("dm://reachable:5236", underlying)
			failed      = make(chan error, 1)
			start       = time.Now()
		)
		go func() {
			_, err := unreachable.Connect(context.Background())
			failed <- err
		}()
		time.Sleep(50 * time.Millisecond)
		_, err := reachable.Connect(context.Background())
		t.AssertNil(err)
		t.Assert(<-underlying.connected, "reachable:5236")
		t.AssertLT(time.Since(start), 300*time.Millisecond)

		t.Assert(<-failed, context.DeadlineExceeded)
		t.AssertGE(time.Since(start), 300*time.Millisecond)
		t.Assert(len(loginSemaphore), 0)
	})
}

// Test_newConnector tests the dial function and the dial timeout of the connector of a data
// source name.
func Test_newConnector(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var cases = []struct {
			dsn         string
			expectDsn   string
			dialTimeout time.Duration
			ownDial     bool
		}{
			{"dm://u:p@h:5236/S?schema=S", "dm://u:p@h:5236/S?schema=S&dialName=" + dialName, defaultConnectTimeout, false},
			{"dm://u:p@h:5236", "dm://u:p@h:5236?dialName=" + dialName, defaultConnectTimeout, false},
			{"dm://u:p@h:5236/S?connectTimeout=1500", "dm://u:p@h:5236/S?connectTimeout=1500&dialName=" + dialName, 1500 * time.Millisecond, false},
			{"dm://u:p@h:5236/S?CONNECTTIMEOUT=0", "dm://u:p@h:5236/S?CONNECTTIMEOUT=0&dialName=" + dialName, 0, false},
			{"dm://u:p@h:5236/S?connectTimeout=x", "dm://u:p@h:5236/S?connectTimeout=x&dialName=" + dialName, defaultConnectTimeout, false},
			{"dm://u:p@h:5236/S?dialName=mine", "dm://u:p@h:5236/S?dialName=mine", defaultConnectTimeout, true},
		}
		for _, c := range cases {
			opened := newConnector(c.dsn, &recordingDriver{})
			t.Assert(opened.dsn, c.expectDsn)
			t.Assert(opened.dialTimeout, c.dialTimeout)
			t.Assert(opened.ownDial, c.ownDial)
		}
	})
}

// Test_connector_OwnDialTakesTurnFirst tests that a connection dialing through a dial function of
// its own logs in once no other connection does, as it does not dial through dialContext.
func Test_connector_OwnDialTakesTurnFirst(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			underlying = &recordingDriver{}
			c          = newConnector("dm://u:p@h:5236?dialName=mine", underlying)
		)
		loginSemaphore <- struct{}{}
		waitCtx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
		defer cancel()
		_, err := c.Connect(waitCtx)
		t.Assert(err, context.DeadlineExceeded)
		t.Assert(len(underlying.connectors), 0)

		<-loginSemaphore
		_, err = c.Connect(context.Background())
		t.AssertNil(err)
		t.Assert(len(underlying.connectors), 1)
		t.Assert(len(loginSemaphore), 0)
	})
}
