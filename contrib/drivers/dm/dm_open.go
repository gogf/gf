// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"database/sql"
	"fmt"

	"net/url"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// Open creates and returns an underlying sql.DB object for pgsql.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		source               string
		underlyingDriverName = "dm"
	)
	if config.Name == "" {
		return nil, fmt.Errorf(
			`dm.Open failed for driver "%s" without DB Name`, underlyingDriverName,
		)
	}
	// Data Source Name of DM8:
	// dm://userName:password@ip:port/dbname
	// dm://userName:password@DW/dbname?DW=(192.168.1.1:5236,192.168.1.2:5236)
	var domain string
	if config.Port != "" {
		domain = fmt.Sprintf("%s:%s", config.Host, config.Port)
	} else {
		domain = config.Host
	}
	source = fmt.Sprintf(
		"dm://%s:%s@%s/%s?charset=%s&schema=%s",
		config.User, config.Pass, domain, config.Name, config.Charset, config.Name,
	)
	// Demo of timezone setting:
	// &loc=Asia/Shanghai
	if config.Timezone != "" {
		if strings.Contains(config.Timezone, "/") {
			config.Timezone = url.QueryEscape(config.Timezone)
		}
		source = fmt.Sprintf("%s&loc%s", source, config.Timezone)
	}
	if config.Extra != "" {
		source = fmt.Sprintf("%s&%s", source, config.Extra)
	}

	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`dm.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}
