// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse

import (
	"database/sql"
	"fmt"
	"net/url"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Open creates and returns an underlying sql.DB object for clickhouse.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var source string
	// clickhouse://username:password@host1:9000,host2:9000/database?dial_timeout=200ms&max_execution_time=60
	if config.Pass != "" {
		source = fmt.Sprintf(
			"clickhouse://%s:%s@%s:%s/%s?debug=%t",
			config.User, url.PathEscape(config.Pass),
			config.Host, config.Port, config.Name, config.Debug,
		)
	} else {
		source = fmt.Sprintf(
			"clickhouse://%s@%s:%s/%s?debug=%t",
			config.User, config.Host, config.Port, config.Name, config.Debug,
		)
	}
	if config.Extra != "" {
		source = fmt.Sprintf("%s&%s", source, config.Extra)
	}
	if db, err = sql.Open(driverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, driverName, source,
		)
		return nil, err
	}
	return
}
