// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"

	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(
	ctx context.Context, link gdb.Link, sql string, args []interface{},
) (newSql string, newArgs []interface{}, err error) {
	// There should be no need to capitalize, because it has been done from field processing before
	newSql, _ = gregex.ReplaceString(`["\n\t]`, "", sql)
	newSql = gstr.ReplaceI(gstr.ReplaceI(newSql, "GROUP_CONCAT", "LISTAGG"), "SEPARATOR", ",")

	// TODO The current approach is too rough. We should deal with the GROUP_CONCAT function and the
	// parsing of the index field from within the select from match.
	// （GROUP_CONCAT DM  does not approve; index cannot be used as a query column name, and security characters need to be added, such as "index"）
	l, r := d.GetChars()
	if strings.Contains(newSql, "INDEX") || strings.Contains(newSql, "index") {
		if !(strings.Contains(newSql, "_INDEX") || strings.Contains(newSql, "_index")) {
			newSql = gstr.ReplaceI(newSql, "INDEX", l+"INDEX"+r)
		}
	}

	// TODO i tried to do but it never work：
	// array, err := gregex.MatchAllString(`SELECT (.*INDEX.*) FROM .*`, newSql)
	// g.Dump("err:", err)
	// g.Dump("array:", array)
	// g.Dump("array:", array[0][1])

	// newSql, err = gregex.ReplaceString(`SELECT (.*INDEX.*) FROM .*`, l+"INDEX"+r, newSql)
	// g.Dump("err:", err)
	// g.Dump("newSql:", newSql)

	// re, err := regexp.Compile(`.*SELECT (.*INDEX.*) FROM .*`)
	// newSql = re.ReplaceAllStringFunc(newSql, func(data string) string {
	// 	fmt.Println("data:", data)
	// 	return data
	// })

	return d.Core.DoFilter(
		ctx,
		link,
		newSql,
		args,
	)
}
