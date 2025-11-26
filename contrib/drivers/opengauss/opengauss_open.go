// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package opengauss

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// Open creates and returns an underlying sql.DB object for openGauss.
// https://pkg.go.dev/gitee.com/opengauss/openGauss-connector-go-pq
func (d *Driver) Open(config *gdb.ConfigNode) (db *sql.DB, err error) {
	source, err := configNodeToSource(config)
	if err != nil {
		return nil, err
	}
	underlyingDriverName := "opengauss"
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
	var (
		parts    []string
		addParam = func(key string, value any) {
			parts = append(parts, fmt.Sprintf(`%s=%s`, key, quoteConnectionValue(fmt.Sprint(value))))
		}
		addOptionalParam = func(key string, value any) {
			valueStr := fmt.Sprint(value)
			if valueStr == "" {
				return
			}
			addParam(key, valueStr)
		}
	)

	addParam("user", config.User)
	addParam("password", config.Pass)
	addParam("host", config.Host)
	addParam("sslmode", "disable")
	addOptionalParam("port", config.Port)
	addOptionalParam("dbname", config.Name)
	addOptionalParam("search_path", config.Namespace)
	addOptionalParam("timezone", config.Timezone)
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
			addOptionalParam(k, v)
		}
	}
	return strings.Join(parts, " "), nil
}

func quoteConnectionValue(value string) string {
	value = strings.ReplaceAll(value, `\`, `\\`)
	value = strings.ReplaceAll(value, `'`, `\\'`)
	return fmt.Sprintf(`'%s'`, value)
}
