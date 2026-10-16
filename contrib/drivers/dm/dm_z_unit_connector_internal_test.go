// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql/driver"
	"testing"

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
