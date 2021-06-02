// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/mattn/go-oci8"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// DriverOracle is the driver for oracle database.
type DriverOracle struct {
	*Core
}

const (
	tableAlias1 = "GFORM1"
	tableAlias2 = "GFORM2"
)

// New creates and returns a database object for oracle.
// It implements the interface of gdb.Driver for extra database driver installation.
func (d *DriverOracle) New(core *Core, node *ConfigNode) (DB, error) {
	return &DriverOracle{
		Core: core,
	}, nil
}

// Open creates and returns a underlying sql.DB object for oracle.
func (d *DriverOracle) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"%s/%s@%s:%s/%s",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("oci8", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// FilteredLinkInfo retrieves and returns filtered `linkInfo` that can be using for
// logging or tracing purpose.
func (d *DriverOracle) FilteredLinkInfo() string {
	linkInfo := d.GetConfig().LinkInfo
	if linkInfo == "" {
		return ""
	}
	s, _ := gregex.ReplaceString(
		`(.+?)\s*/\s*(.+)\s*@\s*(.+)\s*:\s*(\d+)\s*/\s*(.+)`,
		`$1/xxx@$3:$4/$5`,
		linkInfo,
	)
	return s
}

// GetChars returns the security char for this type of database.
func (d *DriverOracle) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// HandleSqlBeforeCommit deals with the sql string before commits it to underlying sql driver.
func (d *DriverOracle) HandleSqlBeforeCommit(ctx context.Context, link Link, sql string, args []interface{}) (newSql string, newArgs []interface{}) {
	var index int
	// Convert place holder char '?' to string ":vx".
	newSql, _ = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf(":v%d", index)
	})
	newSql, _ = gregex.ReplaceString("\"", "", newSql)
	// Handle string datetime argument.
	for i, v := range args {
		if reflect.TypeOf(v).Kind() == reflect.String {
			valueStr := gconv.String(v)
			if gregex.IsMatchString(`^\d{4}-\d{2}-\d{2} \d{2}:\d{2}:\d{2}$`, valueStr) {
				//args[i] = fmt.Sprintf(`TO_DATE('%s','yyyy-MM-dd HH:MI:SS')`, valueStr)
				args[i], _ = time.ParseInLocation("2006-01-02 15:04:05", valueStr, time.Local)
			}
		}
	}
	newSql = d.parseSql(newSql)
	newArgs = args
	return
}

// parseSql does some replacement of the sql before commits it to underlying driver,
// for support of oracle server.
func (d *DriverOracle) parseSql(sql string) string {
	var (
		patten      = `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,{0,1}\s*(\d*))`
		allMatch, _ = gregex.MatchAllString(patten, sql)
	)
	if len(allMatch) == 0 {
		return sql
	}
	var (
		index   = 0
		keyword = strings.ToUpper(strings.TrimSpace(allMatch[index][0]))
	)
	index++
	switch keyword {
	case "SELECT":
		if len(allMatch) < 2 || strings.HasPrefix(allMatch[index][0], "LIMIT") == false {
			break
		}
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}
		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
		if len(queryExpr) != 4 ||
			strings.EqualFold(queryExpr[1], "SELECT") == false ||
			strings.EqualFold(queryExpr[3], "LIMIT") == false {
			break
		}
		first, limit := 0, 0
		for i := 1; i < len(allMatch[index]); i++ {
			if len(strings.TrimSpace(allMatch[index][i])) == 0 {
				continue
			}

			if strings.HasPrefix(allMatch[index][i], "LIMIT") {
				if allMatch[index][i+2] != "" {
					first, _ = strconv.Atoi(allMatch[index][i+1])
					limit, _ = strconv.Atoi(allMatch[index][i+2])
				} else {
					limit, _ = strconv.Atoi(allMatch[index][i+1])
				}
				break
			}
		}
		sql = fmt.Sprintf(
			"SELECT * FROM "+
				"(SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d)"+
				" WHERE ROWNUM_ >= %d",
			queryExpr[1], queryExpr[2], limit, first,
		)
	}
	return sql
}

// Tables retrieves and returns the tables of current schema.
// It's mainly used in cli tool chain for automatically generating the models.
// Note that it ignores the parameter `schema` in oracle database, as it is not necessary.
func (d *DriverOracle) Tables(ctx context.Context, schema ...string) (tables []string, err error) {
	var result Result
	result, err = d.DoGetAll(ctx, nil, "SELECT TABLE_NAME FROM USER_TABLES ORDER BY TABLE_NAME")
	if err != nil {
		return
	}
	for _, m := range result {
		for _, v := range m {
			tables = append(tables, v.String())
		}
	}
	return
}

// TableFields retrieves and returns the fields information of specified table of current schema.
//
// Also see DriverMysql.TableFields.
func (d *DriverOracle) TableFields(ctx context.Context, table string, schema ...string) (fields map[string]*TableField, err error) {
	charL, charR := d.GetChars()
	table = gstr.Trim(table, charL+charR)
	if gstr.Contains(table, " ") {
		return nil, gerror.New("function TableFields supports only single table operations")
	}
	useSchema := d.db.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		useSchema = schema[0]
	}
	tableFieldsCacheKey := fmt.Sprintf(
		`oracle_table_fields_%s_%s@group:%s`,
		table, useSchema, d.GetGroup(),
	)
	v := tableFieldsMap.GetOrSetFuncLock(tableFieldsCacheKey, func() interface{} {
		var (
			result       Result
			link, err    = d.SlaveLink(useSchema)
			structureSql = fmt.Sprintf(`
SELECT 
	COLUMN_NAME AS FIELD, 
	CASE DATA_TYPE  
	WHEN 'NUMBER' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
	WHEN 'FLOAT' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
	ELSE DATA_TYPE||'('||DATA_LENGTH||')' END AS TYPE  
FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s' ORDER BY COLUMN_ID`,
				strings.ToUpper(table),
			)
		)
		if err != nil {
			return nil
		}
		structureSql, _ = gregex.ReplaceString(`[\n\r\s]+`, " ", gstr.Trim(structureSql))
		result, err = d.DoGetAll(ctx, link, structureSql)
		if err != nil {
			return nil
		}
		fields = make(map[string]*TableField)
		for i, m := range result {
			fields[strings.ToLower(m["FIELD"].String())] = &TableField{
				Index: i,
				Name:  strings.ToLower(m["FIELD"].String()),
				Type:  strings.ToLower(m["TYPE"].String()),
			}
		}
		return fields
	})
	if v != nil {
		fields = v.(map[string]*TableField)
	}
	return
}

func (d *DriverOracle) getTableUniqueIndex(table string) (fields map[string]map[string]string, err error) {
	table = strings.ToUpper(table)
	v, _ := internalCache.GetOrSetFunc(
		"table_unique_index_"+table,
		func() (interface{}, error) {
			res := (Result)(nil)
			res, err = d.db.GetAll(fmt.Sprintf(`
		SELECT INDEX_NAME,COLUMN_NAME,CHAR_LENGTH FROM USER_IND_COLUMNS 
		WHERE TABLE_NAME = '%s' 
		AND INDEX_NAME IN(SELECT INDEX_NAME FROM USER_INDEXES WHERE TABLE_NAME='%s' AND UNIQUENESS='UNIQUE') 
		ORDER BY INDEX_NAME,COLUMN_POSITION`, table, table))
			if err != nil {
				return nil, err
			}
			fields := make(map[string]map[string]string)
			for _, v := range res {
				mm := make(map[string]string)
				mm[v["COLUMN_NAME"].String()] = v["CHAR_LENGTH"].String()
				fields[v["INDEX_NAME"].String()] = mm
			}
			return fields, nil
		}, 0)
	if err == nil {
		fields = v.(map[string]map[string]string)
	}
	return
}

func (d *DriverOracle) DoInsert(ctx context.Context, link Link, table string, data interface{}, option int, batch ...int) (result sql.Result, err error) {
	var (
		fields  []string
		values  []string
		params  []interface{}
		dataMap Map
		rv      = reflect.ValueOf(data)
		kind    = rv.Kind()
	)
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Slice, reflect.Array:
		return d.DoBatchInsert(ctx, link, table, data, option, batch...)
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		dataMap = ConvertDataForTableRecord(data)
	default:
		return result, gerror.New(fmt.Sprint("unsupported data type:", kind))
	}
	var (
		indexes     = make([]string, 0)
		indexMap    = make(map[string]string)
		indexExists = false
	)
	if option != insertOptionDefault {
		index, err := d.getTableUniqueIndex(table)
		if err != nil {
			return nil, err
		}

		if len(index) > 0 {
			for _, v := range index {
				for k, _ := range v {
					indexes = append(indexes, k)
				}
				indexMap = v
				indexExists = true
				break
			}
		}
	}
	var (
		subSqlStr = make([]string, 0)
		onStr     = make([]string, 0)
		updateStr = make([]string, 0)
	)
	charL, charR := d.db.GetChars()
	for k, v := range dataMap {
		k = strings.ToUpper(k)

		// 操作类型为REPLACE/SAVE时且存在唯一索引才使用merge，否则使用insert
		if (option == insertOptionReplace || option == insertOptionSave) && indexExists {
			fields = append(fields, tableAlias1+"."+charL+k+charR)
			values = append(values, tableAlias2+"."+charL+k+charR)
			params = append(params, v)
			subSqlStr = append(subSqlStr, fmt.Sprintf("%s?%s %s", charL, charR, k))
			//m erge中的on子句中由唯一索引组成, update子句中不含唯一索引
			if _, ok := indexMap[k]; ok {
				onStr = append(onStr, fmt.Sprintf("%s.%s = %s.%s ", tableAlias1, k, tableAlias2, k))
			} else {
				updateStr = append(updateStr, fmt.Sprintf("%s.%s = %s.%s ", tableAlias1, k, tableAlias2, k))
			}
		} else {
			fields = append(fields, charL+k+charR)
			values = append(values, "?")
			params = append(params, v)
		}
	}

	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return nil, err
		}
	}

	if indexExists && option != insertOptionDefault {
		switch option {
		case
			insertOptionReplace,
			insertOptionSave:
			tmp := fmt.Sprintf(
				"MERGE INTO %s %s USING(SELECT %s FROM DUAL) %s ON(%s) WHEN MATCHED THEN UPDATE SET %s WHEN NOT MATCHED THEN INSERT (%s) VALUES(%s)",
				table, tableAlias1, strings.Join(subSqlStr, ","), tableAlias2,
				strings.Join(onStr, "AND"), strings.Join(updateStr, ","), strings.Join(fields, ","), strings.Join(values, ","),
			)
			return d.DoExec(ctx, link, tmp, params...)

		case insertOptionIgnore:
			return d.DoExec(ctx, link, fmt.Sprintf(
				"INSERT /*+ IGNORE_ROW_ON_DUPKEY_INDEX(%s(%s)) */ INTO %s(%s) VALUES(%s)",
				table, strings.Join(indexes, ","), table, strings.Join(fields, ","), strings.Join(values, ","),
			), params...)
		}
	}

	return d.DoExec(ctx, link,
		fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES(%s)",
			table, strings.Join(fields, ","), strings.Join(values, ","),
		),
		params...)
}

func (d *DriverOracle) DoBatchInsert(ctx context.Context, link Link, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	var (
		keys   []string
		values []string
		params []interface{}
	)
	listMap := (List)(nil)
	switch v := list.(type) {
	case Result:
		listMap = v.List()
	case Record:
		listMap = List{v.Map()}
	case List:
		listMap = v
	case Map:
		listMap = List{v}
	default:
		var (
			rv   = reflect.ValueOf(list)
			kind = rv.Kind()
		)
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		case reflect.Slice, reflect.Array:
			listMap = make(List, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				listMap[i] = ConvertDataForTableRecord(rv.Index(i).Interface())
			}
		case reflect.Map:
			fallthrough
		case reflect.Struct:
			listMap = List{ConvertDataForTableRecord(list)}
		default:
			return result, gerror.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	if len(listMap) < 1 {
		return result, gerror.New("empty data list")
	}
	if link == nil {
		if link, err = d.MasterLink(); err != nil {
			return
		}
	}
	// Retrieve the table fields and length.
	holders := []string(nil)
	for k, _ := range listMap[0] {
		keys = append(keys, k)
		holders = append(holders, "?")
	}
	var (
		batchResult    = new(SqlResult)
		charL, charR   = d.db.GetChars()
		keyStr         = charL + strings.Join(keys, charL+","+charR) + charR
		valueHolderStr = strings.Join(holders, ",")
	)
	if option != insertOptionDefault {
		for _, v := range listMap {
			r, err := d.DoInsert(ctx, link, table, v, option, 1)
			if err != nil {
				return r, err
			}

			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.result = r
				batchResult.affected += n
			}
		}
		return batchResult, nil
	}

	batchNum := defaultBatchNumber
	if len(batch) > 0 {
		batchNum = batch[0]
	}
	// Format "INSERT...INTO..." statement.
	intoStr := make([]string, 0)
	for i := 0; i < len(listMap); i++ {
		for _, k := range keys {
			params = append(params, listMap[i][k])
		}
		values = append(values, valueHolderStr)
		intoStr = append(intoStr, fmt.Sprintf(" INTO %s(%s) VALUES(%s) ", table, keyStr, valueHolderStr))
		if len(intoStr) == batchNum {
			r, err := d.DoExec(ctx, link, fmt.Sprintf("INSERT ALL %s SELECT * FROM DUAL", strings.Join(intoStr, " ")), params...)
			if err != nil {
				return r, err
			}
			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.result = r
				batchResult.affected += n
			}
			params = params[:0]
			intoStr = intoStr[:0]
		}
	}
	// The leftover data.
	if len(intoStr) > 0 {
		r, err := d.DoExec(ctx, link, fmt.Sprintf("INSERT ALL %s SELECT * FROM DUAL", strings.Join(intoStr, " ")), params...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.result = r
			batchResult.affected += n
		}
	}
	return batchResult, nil
}
