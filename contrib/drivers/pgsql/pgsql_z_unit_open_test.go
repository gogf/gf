// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"

	"github.com/gogf/gf/contrib/drivers/pgsql/v2"
)

// Test_Open tests the Open method with various configurations
func Test_Open_WithNamespace(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User:      "postgres",
			Pass:      "12345678",
			Host:      "127.0.0.1",
			Port:      "5432",
			Name:      "test",
			Namespace: "public",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_WithTimezone tests Open with timezone configuration
func Test_Open_WithTimezone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User:     "postgres",
			Pass:     "12345678",
			Host:     "127.0.0.1",
			Port:     "5432",
			Name:     "test",
			Timezone: "Asia/Shanghai",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_WithExtra tests Open with extra configuration
func Test_Open_WithExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User:  "postgres",
			Pass:  "12345678",
			Host:  "127.0.0.1",
			Port:  "5432",
			Name:  "test",
			Extra: "connect_timeout=10",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_WithInvalidExtra tests Open with invalid extra configuration
func Test_Open_WithInvalidExtra(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User: "postgres",
			Pass: "12345678",
			Host: "127.0.0.1",
			Port: "5432",
			Name: "test",
			// Invalid extra format with invalid URL encoding that will cause parse error
			Extra: "%Q=%Q&b",
		}
		_, err := driver.Open(config)
		t.AssertNE(err, nil)
	})
}

// Test_Open_WithFullConfig tests Open with all configuration options
func Test_Open_WithFullConfig(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User:      "postgres",
			Pass:      "12345678",
			Host:      "127.0.0.1",
			Port:      "5432",
			Name:      "test",
			Namespace: "public",
			Timezone:  "UTC",
			Extra:     "connect_timeout=10",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_WithoutPort tests Open without port
func Test_Open_WithoutPort(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User: "postgres",
			Pass: "12345678",
			Host: "127.0.0.1",
			Name: "test",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_WithoutName tests Open without database name
func Test_Open_WithoutName(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User: "postgres",
			Pass: "12345678",
			Host: "127.0.0.1",
			Port: "5432",
		}
		db, err := driver.Open(config)
		t.AssertNil(err)
		t.AssertNE(db, nil)
		if db != nil {
			db.Close()
		}
	})
}

// Test_Open_InvalidHost tests Open with invalid host
func Test_Open_InvalidHost(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		driver := pgsql.Driver{}
		config := &gdb.ConfigNode{
			User: "postgres",
			Pass: "12345678",
			Host: "invalid_host_that_does_not_exist",
			Port: "5432",
			Name: "test",
		}
		// Note: sql.Open doesn't actually connect, so no error here
		// The error would occur when actually using the connection
		db, err := driver.Open(config)
		t.AssertNil(err)
		if db != nil {
			db.Close()
		}
	})
}
