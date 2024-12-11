// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mssql

import (
	"database/sql"
	"fmt"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
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
			source += fmt.Sprintf(`;%s=%s`, k, v)
		}
	}
	return source, nil
}
