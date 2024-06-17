// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle

import (
	"database/sql"
	"strings"

	gora "github.com/sijms/go-ora/v2"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/util/gconv"
)

// Open creates and returns an underlying sql.DB object for oracle.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "oracle"
	)

	options := map[string]string{
		"CONNECTION TIMEOUT": "60",
		"PREFETCH_ROWS":      "25",
	}

	if config.Debug {
		options["TRACE FILE"] = "oracle_trace.log"
	}
	// [username:[password]@]host[:port][/service_name][?param1=value1&...&paramN=valueN]
	if config.Link != "" {
		// ============================================================================
		// Deprecated from v2.2.0.
		// ============================================================================
		source = config.Link
		// Custom changing the schema in runtime.
		if config.Name != "" {
			source, _ = gregex.ReplaceString(`@(.+?)/([\w\.\-]+)+`, "@$1/"+config.Name, source)
		}
	} else {
		if config.Extra != "" {
			// fix #3226
			list := strings.Split(config.Extra, "&")
			for _, v := range list {
				kv := strings.Split(v, "=")
				if len(kv) == 2 {
					options[kv[0]] = kv[1]
				}
			}
		}
		source = gora.BuildUrl(
			config.Host, gconv.Int(config.Port), config.Name, config.User, config.Pass, options,
		)
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
