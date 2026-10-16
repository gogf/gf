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

func (d *blockingDriver) Connect(context.Context) (driver.Conn, error) {
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
