// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package sqlitecgo implements gdb.Driver, which supports operations for database SQLite.
//
// Note:
//  1. Using sqlitecgo is for building a 32-bit Windows operating system
//  2. You need to set the environment variable CGO_ENABLED=1 and make sure that GCC is installed
//     on your path. windows gcc: https://jmeubank.github.io/tdm-gcc/
package sqlitecgo

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/encoding/gurl"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// Driver is the driver for sqlite database.
type Driver struct {
	gdb.DB
}

var (
	sqliteDriver = sqlite.New()
)

func init() {
	if err := gdb.Register(`sqlite`, New()); err != nil {
		panic(err)
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for SQLite.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for sqlite.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	db, err := sqliteDriver.New(core, node)
	if err != nil {
		return nil, err
	}
	return &Driver{
		DB: db,
	}, nil
}

// Open creates and returns an underlying sql.DB object for sqlite.
// https://github.com/mattn/go-sglite3
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "sqlite3"
	)
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
	} else {
		source = config.Name
	}
	// It searches the source file to locate its absolute path..
	if absolutePath, _ := gfile.Search(source); absolutePath != "" {
		source = absolutePath
	}

	// Multiple PRAGMAs can be specified, e.g.:
	// path/to/some.db?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)
	if config.Extra != "" {
		var (
			options  string
			extraMap map[string]interface{}
		)
		if extraMap, err = gstr.Parse(config.Extra); err != nil {
			return nil, err
		}
		for k, v := range extraMap {
			if options != "" {
				options += "&"
			}
			options += fmt.Sprintf(`_pragma=%s(%s)`, k, gurl.Encode(gconv.String(v)))
		}
		if len(options) > 1 {
			source += "?" + options
		}
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}
