// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// timezoneParamName is the connection parameter naming the time zone in which the underlying
	// driver reads values of the types without time zone, like datetime and datetime2.
	timezoneParamName = "timezone"
	// timezoneDefault is the local time zone, in which such values are written, used unless the
	// timezone configuration or the parameter sets another.
	timezoneDefault = "Local"
)

// Open creates and returns an underlying sql.DB object for mssql.
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	source, err := configNodeToSource(config)
	if err != nil {
		return nil, err
	}
	underlyingDriverName := "sqlserver"
	if db, err = sql.Open(underlyingDriverName, source); err != nil {
		err = gerror.WrapCodef(
			gcode.CodeDbOperationError, err,
			`sql.Open failed for driver "%s" by source "%s"`, underlyingDriverName, source,
		)
		return nil, err
	}
	return
}

// configNodeToLocation returns the time zone in which the underlying driver reads the values of the
// types without time zone: the one that the connection parameter timezoneParamName of `config`
// names, or else the one of its timezone configuration, or else the local time zone.
func configNodeToLocation(config *gdb.ConfigNode) (*time.Location, error) {
	timezone := configNodeToTimezone(config)
	if config.Extra != "" {
		extraMap, err := gstr.Parse(config.Extra)
		if err != nil {
			return nil, gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				err,
				`invalid extra configuration: %s`, config.Extra,
			)
		}
		for k, v := range extraMap {
			if strings.EqualFold(k, timezoneParamName) {
				timezone = gconv.String(v)
			}
		}
	}
	location, err := time.LoadLocation(timezone)
	if err != nil {
		return nil, gerror.WrapCodef(
			gcode.CodeInvalidParameter,
			err,
			`invalid time zone "%s"`, timezone,
		)
	}
	return location, nil
}

// configNodeToTimezone returns the time zone that the timezone configuration of `config` names,
// or timezoneDefault if it names none.
func configNodeToTimezone(config *gdb.ConfigNode) string {
	if config.Timezone != "" {
		return config.Timezone
	}
	return timezoneDefault
}

func configNodeToSource(config *gdb.ConfigNode) (string, error) {
	var source string
	source = fmt.Sprintf(
		"user id=%s;password=%s;server=%s;encrypt=disable",
		config.User, config.Pass, config.Host,
	)
	if config.Name != "" {
		source = fmt.Sprintf("%s;database=%s", source, config.Name)
	}
	if config.Port != "" {
		source = fmt.Sprintf("%s;port=%s", source, config.Port)
	}
	hasTimezone := false
	if config.Extra != "" {
		extraMap, err := gstr.Parse(config.Extra)
		if err != nil {
			return "", gerror.WrapCodef(
				gcode.CodeInvalidParameter,
				err,
				`invalid extra configuration: %s`, config.Extra,
			)
		}
		for k, v := range extraMap {
			if strings.EqualFold(k, timezoneParamName) {
				hasTimezone = true
			}
			source += fmt.Sprintf(`;%s=%s`, k, v)
		}
	}
	if !hasTimezone {
		source += fmt.Sprintf(`;%s=%s`, timezoneParamName, configNodeToTimezone(config))
	}
	return source, nil
}
