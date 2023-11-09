// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package dm implements gdb.Driver, which supports operations for database DM.
package dm

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"
	"strings"
	"time"

	_ "gitee.com/chunanyong/dm"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

type Driver struct {
	*gdb.Core
}

const (
	quoteChar = `"`
)

func init() {
	var (
		err         error
		driverObj   = New()
		driverNames = g.SliceStr{"dm"}
	)
	for _, driverName := range driverNames {
		if err = gdb.Register(driverName, driverObj); err != nil {
			panic(err)
		}
	}
}

// New create and returns a driver that implements gdb.Driver, which supports operations for dm.
func New() gdb.Driver {
	return &Driver{}
}

// New creates and returns a database object for dm.
func (d *Driver) New(core *gdb.Core, node *gdb.ConfigNode) (gdb.DB, error) {
	return &Driver{
		Core: core,
	}, nil
}

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
	source = fmt.Sprintf(
		"dm://%s:%s@%s:%s/%s?charset=%s&schema=%s",
		config.User, config.Pass, config.Host, config.Port, config.Name, config.Charset, config.Name,
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

// GetChars returns the security char for this type of database.
func (d *Driver) GetChars() (charLeft string, charRight string) {
	return quoteChar, quoteChar
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
func (d *Driver) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result gdb.Result
	// When schema is empty, return the default link
	link, err := d.SlaveLink(schema...)
	if err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	result, err = d.DoSelect(
		ctx, link, `SELECT * FROM ALL_TABLES`,
	)
	if err != nil {
		return
	}
	for _, m := range result {
		if v, ok := m["IOT_NAME"]; ok {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*gdb.TableField, err error) {
	var (
		result gdb.Result
		link   gdb.Link
		// When no schema is specified, the configuration item is returned by default
		usedSchema = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
	)
	// When usedSchema is empty, return the default link
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	result, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(
			`SELECT * FROM ALL_TAB_COLUMNS WHERE Table_Name= '%s' AND OWNER = '%s'`,
			strings.ToUpper(table),
			strings.ToUpper(d.GetSchema()),
		),
	)
	if err != nil {
		return nil, err
	}
	fields = make(map[string]*gdb.TableField)
	for i, m := range result {
		// m[NULLABLE] returns "N" "Y"
		// "N" means not null
		// "Y" means could be null
		var nullable bool
		if m["NULLABLE"].String() != "N" {
			nullable = true
		}
		fields[m["COLUMN_NAME"].String()] = &gdb.TableField{
			Index:   i,
			Name:    m["COLUMN_NAME"].String(),
			Type:    m["DATA_TYPE"].String(),
			Null:    nullable,
			Default: m["DATA_DEFAULT"].Val(),
			// Key:     m["Key"].String(),
			// Extra:   m["Extra"].String(),
			// Comment: m["Comment"].String(),
		}
	}
	return fields, nil
}

// ConvertValueForField converts value to the type of the record field.
func (d *Driver) ConvertValueForField(ctx context.Context, fieldType string, fieldValue interface{}) (interface{}, error) {
	switch itemValue := fieldValue.(type) {
	// dm does not support time.Time, it so here converts it to time string that it supports.
	case time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue.IsZero() {
			return nil, nil
		}
		return gtime.New(itemValue).String(), nil

	// dm does not support time.Time, it so here converts it to time string that it supports.
	case *time.Time:
		// If the time is zero, it then updates it to nil,
		// which will insert/update the value to database as "null".
		if itemValue == nil || itemValue.IsZero() {
			return nil, nil
		}
		return gtime.New(itemValue).String(), nil
	}

	return fieldValue, nil
}

// DoFilter deals with the sql string before commits it to underlying sql driver.
func (d *Driver) DoFilter(ctx context.Context, link gdb.Link, sql string, args []interface{}) (newSql string, newArgs []interface{}, err error) {
	// There should be no need to capitalize, because it has been done from field processing before
	newSql, _ = gregex.ReplaceString(`["\n\t]`, "", sql)
	newSql = gstr.ReplaceI(gstr.ReplaceI(newSql, "GROUP_CONCAT", "LISTAGG"), "SEPARATOR", ",")

	// TODO The current approach is too rough. We should deal with the GROUP_CONCAT function and the parsing of the index field from within the select from match.
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

// TODO I originally wanted to only convert keywords in select
// 但是我发现 DoQuery 中会对 sql 会对 " " 达梦的安全字符 进行 / 转义，最后还是导致达梦无法正常解析
// However, I found that DoQuery() will perform / escape on sql with " " Dameng's safe characters, which ultimately caused Dameng to be unable to parse normally.
// But processing in DoFilter() is OK
// func (d *Driver) DoQuery(ctx context.Context, link gdb.Link, sql string, args ...interface{}) (gdb.Result, error) {
// 	l, r := d.GetChars()
// 	new := gstr.ReplaceI(sql, "INDEX", l+"INDEX"+r)
// 	g.Dump("new:", new)
// 	return d.Core.DoQuery(
// 		ctx,
// 		link,
// 		new,
// 		args,
// 	)
// }

// DoInsert inserts or updates data forF given table.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionReplace:
		// TODO:: Should be Supported
		return nil, gerror.NewCode(
			gcode.CodeNotSupported, `Replace operation is not supported by dm driver`,
		)

	case gdb.InsertOptionSave:
		// This syntax currently only supports design tables whose primary key is ID.
		listLength := len(list)
		if listLength == 0 {
			return nil, gerror.NewCode(
				gcode.CodeInvalidRequest, `Save operation list is empty by dm driver`,
			)
		}
		var (
			keysSort     []string
			charL, charR = d.GetChars()
		)
		// Column names need to be aligned in the syntax
		for k := range list[0] {
			keysSort = append(keysSort, k)
		}
		var char = struct {
			charL        string
			charR        string
			valueCharL   string
			valueCharR   string
			duplicateKey string
			keys         []string
		}{
			charL:      charL,
			charR:      charR,
			valueCharL: "'",
			valueCharR: "'",
			// TODO:: Need to dynamically set the primary key of the table
			duplicateKey: "ID",
			keys:         keysSort,
		}

		// insertKeys:   Handle valid keys that need to be inserted and updated
		// insertValues: Handle values that need to be inserted
		// updateValues: Handle values that need to be updated
		// queryValues:  Handle only one insert with column name
		insertKeys, insertValues, updateValues, queryValues := parseValue(list[0], char)
		// unionValues: Handling values that need to be inserted and updated
		unionValues := parseUnion(list[1:], char)

		batchResult := new(gdb.SqlResult)
		// parseSql():
		// MERGE INTO {{table}} T1
		// USING ( SELECT {{queryValues}} FROM DUAL
		// {{unionValues}} ) T2
		// ON (T1.{{duplicateKey}} = T2.{{duplicateKey}})
		// WHEN NOT MATCHED THEN
		// INSERT {{insertKeys}} VALUES {{insertValues}}
		// WHEN MATCHED THEN
		// UPDATE SET {{updateValues}}
		sqlStr := parseSql(
			insertKeys, insertValues, updateValues, queryValues, unionValues, table, char.duplicateKey,
		)
		r, err := d.DoExec(ctx, link, sqlStr)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
		return batchResult, nil
	}
	return d.Core.DoInsert(ctx, link, table, list, option)
}

func parseValue(listOne gdb.Map, char struct {
	charL        string
	charR        string
	valueCharL   string
	valueCharR   string
	duplicateKey string
	keys         []string
}) (insertKeys []string, insertValues []string, updateValues []string, queryValues []string) {
	for _, column := range char.keys {
		if listOne[column] == nil {
			// remove unassigned struct object
			continue
		}
		insertKeys = append(insertKeys, char.charL+column+char.charR)
		insertValues = append(insertValues, "T2."+char.charL+column+char.charR)
		if column != char.duplicateKey {
			updateValues = append(
				updateValues,
				fmt.Sprintf(`T1.%s = T2.%s`, char.charL+column+char.charR, char.charL+column+char.charR),
			)
		}

		saveValue := gconv.String(listOne[column])
		queryValues = append(
			queryValues,
			fmt.Sprintf(
				char.valueCharL+"%s"+char.valueCharR+" AS "+char.charL+"%s"+char.charR,
				saveValue, column,
			),
		)
	}
	return
}

func parseUnion(list gdb.List, char struct {
	charL        string
	charR        string
	valueCharL   string
	valueCharR   string
	duplicateKey string
	keys         []string
}) (unionValues []string) {
	for _, mapper := range list {
		var saveValue []string
		for _, column := range char.keys {
			if mapper[column] == nil {
				continue
			}
			// va := reflect.ValueOf(mapper[column])
			// ty := reflect.TypeOf(mapper[column])
			// switch ty.Kind() {
			// case reflect.String:
			// 	saveValue = append(saveValue, char.valueCharL+va.String()+char.valueCharR)

			// case reflect.Int:
			// 	saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))

			// case reflect.Int64:
			// 	saveValue = append(saveValue, strconv.FormatInt(va.Int(), 10))

			// default:
			// 	// The fish has no chance getting here.
			// 	// Nothing to do.
			// }
			saveValue = append(saveValue,
				fmt.Sprintf(
					char.valueCharL+"%s"+char.valueCharR,
					gconv.String(mapper[column]),
				))
		}
		unionValues = append(
			unionValues,
			fmt.Sprintf(`UNION ALL SELECT %s FROM DUAL`, strings.Join(saveValue, ",")),
		)
	}
	return
}

func parseSql(
	insertKeys, insertValues, updateValues, queryValues, unionValues []string, table, duplicateKey string,
) (sqlStr string) {
	var (
		queryValueStr  = strings.Join(queryValues, ",")
		unionValueStr  = strings.Join(unionValues, " ")
		insertKeyStr   = strings.Join(insertKeys, ",")
		insertValueStr = strings.Join(insertValues, ",")
		updateValueStr = strings.Join(updateValues, ",")
		pattern        = gstr.Trim(`
MERGE INTO %s T1 USING (SELECT %s FROM DUAL %s) T2 ON %s 
WHEN NOT MATCHED 
THEN 
INSERT(%s) VALUES (%s) 
WHEN MATCHED 
THEN 
UPDATE SET %s; 
COMMIT;
`)
	)
	return fmt.Sprintf(
		pattern,
		table, queryValueStr, unionValueStr,
		fmt.Sprintf("(T1.%s = T2.%s)", duplicateKey, duplicateKey),
		insertKeyStr, insertValueStr, updateValueStr,
	)
}
