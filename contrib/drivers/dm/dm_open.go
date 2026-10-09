// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// dm_open.go builds the DM data source name and opens the underlying SQL database.

package dm

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
)

// dmTimezoneOffsetLimitMinutes is the maximum absolute offset accepted by DM.
const dmTimezoneOffsetLimitMinutes = 12 * 60

// Open creates and returns an underlying sql.DB object for DM.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	var (
		underlyingDriverName = "dm"
	)
	if config.Name == "" {
		return nil, fmt.Errorf(
			`dm.Open failed for driver "%s" without DB Name`, underlyingDriverName,
		)
	}
	source, err := buildDataSourceName(config, time.Now())
	if err != nil {
		return nil, err
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

// buildDataSourceName creates a DM DSN and converts the configured IANA timezone
// to the fixed offset in minutes accepted by the DM driver. Daylight-saving zones
// use the offset active when the DSN is built.
func buildDataSourceName(config *gdb.ConfigNode, now time.Time) (string, error) {
	var domain string
	if config.Port != "" {
		domain = fmt.Sprintf("%s:%s", config.Host, config.Port)
	} else {
		domain = config.Host
	}
	source := fmt.Sprintf(
		"dm://%s:%s@%s/%s?charset=%s&schema=%s",
		config.User, config.Pass, domain, config.Name, config.Charset, config.Name,
	)
	if config.Timezone != "" {
		location, err := time.LoadLocation(config.Timezone)
		if err != nil {
			return "", gerror.WrapCodef(
				gcode.CodeInvalidParameter, err,
				`invalid timezone "%s"`, config.Timezone,
			)
		}
		_, offset := now.In(location).Zone()
		offsetMinutes := offset / 60
		if offsetMinutes < -dmTimezoneOffsetLimitMinutes || offsetMinutes > dmTimezoneOffsetLimitMinutes {
			return "", gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`timezone "%s" has an offset of %d minutes; DM supports offsets from -%d to %d minutes`,
				config.Timezone, offsetMinutes, dmTimezoneOffsetLimitMinutes, dmTimezoneOffsetLimitMinutes,
			)
		}
		source = fmt.Sprintf("%s&timeZone=%d", source, offsetMinutes)
	}
	if config.Extra != "" {
		source = fmt.Sprintf("%s&%s", source, config.Extra)
	}
	return source, nil
}
