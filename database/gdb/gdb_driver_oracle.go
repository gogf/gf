// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
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
	"database/sql"
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
	"reflect"
	"strconv"
	"strings"

	"github.com/gogf/gf/text/gregex"
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
		source = fmt.Sprintf("%s/%s@%s", config.User, config.Pass, config.Name)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("oci8", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// GetChars returns the security char for this type of database.
func (d *DriverOracle) GetChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// HandleSqlBeforeExec deals with the sql string before commits it to underlying sql driver.
func (d *DriverOracle) HandleSqlBeforeExec(query string) string {
	var index int
	// Convert place holder char '?' to string ":x".
	str, _ := gregex.ReplaceStringFunc("\\?", query, func(s string) string {
		index++
		return fmt.Sprintf(":%d", index)
	})
	str, _ = gregex.ReplaceString("\"", "", str)
	return d.parseSql(str)
}

func (d *DriverOracle) parseSql(sql string) string {
	patten := `^\s*(?i)(SELECT)|(LIMIT\s*(\d+)\s*,\s*(\d+))`
	if gregex.IsMatchString(patten, sql) == false {
		return sql
	}

	res, err := gregex.MatchAllString(patten, sql)
	if err != nil {
		return ""
	}

	index := 0
	keyword := strings.TrimSpace(res[index][0])
	keyword = strings.ToUpper(keyword)

	index++
	switch keyword {
	case "SELECT":
		// 不含LIMIT关键字则不处理
		if len(res) < 2 || (strings.HasPrefix(res[index][0], "LIMIT") == false && strings.HasPrefix(res[index][0], "limit") == false) {
			break
		}

		// 取limit前面的字符串
		if gregex.IsMatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql) == false {
			break
		}

		queryExpr, _ := gregex.MatchString("((?i)SELECT)(.+)((?i)LIMIT)", sql)
		if len(queryExpr) != 4 || strings.EqualFold(queryExpr[1], "SELECT") == false || strings.EqualFold(queryExpr[3], "LIMIT") == false {
			break
		}

		// 取limit后面的取值范围
		first, limit := 0, 0
		for i := 1; i < len(res[index]); i++ {
			if len(strings.TrimSpace(res[index][i])) == 0 {
				continue
			}

			if strings.HasPrefix(res[index][i], "LIMIT") || strings.HasPrefix(res[index][i], "limit") {
				first, _ = strconv.Atoi(res[index][i+1])
				limit, _ = strconv.Atoi(res[index][i+2])
				break
			}
		}

		// 也可以使用between,据说这种写法的性能会比between好点,里层SQL中的ROWNUM_ >= limit可以缩小查询后的数据集规模
		sql = fmt.Sprintf(
			"SELECT * FROM (SELECT GFORM.*, ROWNUM ROWNUM_ FROM (%s %s) GFORM WHERE ROWNUM <= %d) WHERE ROWNUM_ >= %d",
			queryExpr[1], queryExpr[2], limit, first,
		)
	}
	return sql
}

// Tables retrieves and returns the tables of current schema.
// Note that it ignores the parameter <schema> in oracle database, as it is not necessary.
func (d *DriverOracle) Tables(schema ...string) (tables []string, err error) {
	var result Result
	result, err = d.DB.DoGetAll(nil, "SELECT TABLE_NAME FROM USER_TABLES ORDER BY TABLE_NAME")
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
func (d *DriverOracle) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function TableFields supports only single table operations")
	}
	checkSchema := d.DB.GetSchema()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v := d.DB.GetCache().GetOrSetFunc(
		fmt.Sprintf(`oracle_table_fields_%s_%s`, table, checkSchema),
		func() interface{} {
			result := (Result)(nil)
			result, err = d.DB.GetAll(fmt.Sprintf(`
			SELECT COLUMN_NAME AS FIELD, CASE DATA_TYPE 
			    WHEN 'NUMBER' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
				WHEN 'FLOAT' THEN DATA_TYPE||'('||DATA_PRECISION||','||DATA_SCALE||')' 
				ELSE DATA_TYPE||'('||DATA_LENGTH||')' END AS TYPE  
			FROM USER_TAB_COLUMNS WHERE TABLE_NAME = '%s' ORDER BY COLUMN_ID`, strings.ToUpper(table)))
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
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}

func (d *DriverOracle) getTableUniqueIndex(table string) (fields map[string]map[string]string, err error) {
	table = strings.ToUpper(table)
	v := d.DB.GetCache().GetOrSetFunc("table_unique_index_"+table, func() interface{} {
		res := (Result)(nil)
		res, err = d.DB.GetAll(fmt.Sprintf(`
		SELECT INDEX_NAME,COLUMN_NAME,CHAR_LENGTH FROM USER_IND_COLUMNS 
		WHERE TABLE_NAME = '%s' 
		AND INDEX_NAME IN(SELECT INDEX_NAME FROM USER_INDEXES WHERE TABLE_NAME='%s' AND UNIQUENESS='UNIQUE') 
		ORDER BY INDEX_NAME,COLUMN_POSITION`, table, table))
		if err != nil {
			return nil
		}
		fields := make(map[string]map[string]string)
		for _, v := range res {
			mm := make(map[string]string)
			mm[v["COLUMN_NAME"].String()] = v["CHAR_LENGTH"].String()
			fields[v["INDEX_NAME"].String()] = mm
		}
		return fields
	}, 0)
	if err == nil {
		fields = v.(map[string]map[string]string)
	}
	return
}

func (d *DriverOracle) DoInsert(link dbLink, table string, data interface{}, option int, batch ...int) (result sql.Result, err error) {
	var fields []string
	var values []string
	var params []interface{}
	var dataMap Map
	rv := reflect.ValueOf(data)
	kind := rv.Kind()
	if kind == reflect.Ptr {
		rv = rv.Elem()
		kind = rv.Kind()
	}
	switch kind {
	case reflect.Slice:
		fallthrough
	case reflect.Array:
		return d.DB.DoBatchInsert(link, table, data, option, batch...)
	case reflect.Map:
		fallthrough
	case reflect.Struct:
		dataMap = varToMapDeep(data)
	default:
		return result, errors.New(fmt.Sprint("unsupported data type:", kind))
	}

	indexs := make([]string, 0)
	indexMap := make(map[string]string)
	indexExists := false
	if option != gINSERT_OPTION_DEFAULT {
		index, err := d.getTableUniqueIndex(table)
		if err != nil {
			return nil, err
		}

		if len(index) > 0 {
			for _, v := range index {
				for k, _ := range v {
					indexs = append(indexs, k)
				}
				indexMap = v
				indexExists = true
				break
			}
		}

	}

	subSqlStr := make([]string, 0)
	onStr := make([]string, 0)
	updateStr := make([]string, 0)

	charL, charR := d.DB.GetChars()
	for k, v := range dataMap {
		k = strings.ToUpper(k)

		// 操作类型为REPLACE/SAVE时且存在唯一索引才使用merge，否则使用insert
		if (option == gINSERT_OPTION_REPLACE || option == gINSERT_OPTION_SAVE) && indexExists {
			fields = append(fields, tableAlias1+"."+charL+k+charR)
			values = append(values, tableAlias2+"."+charL+k+charR)
			params = append(params, v)

			subSqlStr = append(subSqlStr, fmt.Sprintf("%s?%s %s", charL, charR, k))

			//merge中的on子句中由唯一索引组成,update子句中不含唯一索引
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
		if link, err = d.DB.Master(); err != nil {
			return nil, err
		}
	}

	if indexExists && option != gINSERT_OPTION_DEFAULT {
		switch option {
		case gINSERT_OPTION_REPLACE:
			fallthrough
		case gINSERT_OPTION_SAVE:
			tmp := fmt.Sprintf(
				"MERGE INTO %s %s USING(SELECT %s FROM DUAL) %s ON(%s) WHEN MATCHED THEN UPDATE SET %s WHEN NOT MATCHED THEN INSERT (%s) VALUES(%s)",
				table, tableAlias1, strings.Join(subSqlStr, ","), tableAlias2,
				strings.Join(onStr, "AND"), strings.Join(updateStr, ","), strings.Join(fields, ","), strings.Join(values, ","),
			)
			return d.DB.DoExec(link, tmp, params...)
		case gINSERT_OPTION_IGNORE:
			return d.DB.DoExec(link,
				fmt.Sprintf(
					"INSERT /*+ IGNORE_ROW_ON_DUPKEY_INDEX(%s(%s)) */ INTO %s(%s) VALUES(%s)",
					table, strings.Join(indexs, ","), table, strings.Join(fields, ","), strings.Join(values, ","),
				),
				params...)
		}
	}

	return d.DB.DoExec(
		link,
		fmt.Sprintf(
			"INSERT INTO %s(%s) VALUES(%s)",
			table, strings.Join(fields, ","), strings.Join(values, ","),
		),
		params...)
}

func (d *DriverOracle) DoBatchInsert(link dbLink, table string, list interface{}, option int, batch ...int) (result sql.Result, err error) {
	var keys []string
	var values []string
	var params []interface{}
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
		rv := reflect.ValueOf(list)
		kind := rv.Kind()
		if kind == reflect.Ptr {
			rv = rv.Elem()
			kind = rv.Kind()
		}
		switch kind {
		// 如果是slice，那么转换为List类型
		case reflect.Slice:
			fallthrough
		case reflect.Array:
			listMap = make(List, rv.Len())
			for i := 0; i < rv.Len(); i++ {
				listMap[i] = varToMapDeep(rv.Index(i).Interface())
			}
		case reflect.Map:
			fallthrough
		case reflect.Struct:
			listMap = List{Map(varToMapDeep(list))}
		default:
			return result, errors.New(fmt.Sprint("unsupported list type:", kind))
		}
	}
	// 判断长度
	if len(listMap) < 1 {
		return result, errors.New("empty data list")
	}
	if link == nil {
		if link, err = d.DB.Master(); err != nil {
			return
		}
	}
	// 首先获取字段名称及记录长度
	holders := []string(nil)
	for k, _ := range listMap[0] {
		keys = append(keys, k)
		holders = append(holders, "?")
	}
	batchResult := new(batchSqlResult)
	charL, charR := d.DB.GetChars()
	keyStr := charL + strings.Join(keys, charL+","+charR) + charR
	valueHolderStr := strings.Join(holders, ",")

	// 当操作类型非insert时调用单笔的insert功能
	if option != gINSERT_OPTION_DEFAULT {
		for _, v := range listMap {
			r, err := d.DB.DoInsert(link, table, v, option, 1)
			if err != nil {
				return r, err
			}

			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.lastResult = r
				batchResult.rowsAffected += n
			}
		}
		return batchResult, nil
	}

	// 构造批量写入数据格式(注意map的遍历是无序的)
	batchNum := gDEFAULT_BATCH_NUM
	if len(batch) > 0 {
		batchNum = batch[0]
	}

	intoStr := make([]string, 0) //组装into语句
	for i := 0; i < len(listMap); i++ {
		for _, k := range keys {
			params = append(params, listMap[i][k])
		}
		values = append(values, valueHolderStr)
		intoStr = append(intoStr, fmt.Sprintf(" INTO %s(%s) VALUES(%s) ", table, keyStr, valueHolderStr))
		if len(intoStr) == batchNum {
			r, err := d.DB.DoExec(link, fmt.Sprintf("INSERT ALL %s SELECT * FROM DUAL", strings.Join(intoStr, " ")), params...)
			if err != nil {
				return r, err
			}
			if n, err := r.RowsAffected(); err != nil {
				return r, err
			} else {
				batchResult.lastResult = r
				batchResult.rowsAffected += n
			}
			params = params[:0]
			intoStr = intoStr[:0]
		}
	}
	// 处理最后不构成指定批量的数据
	if len(intoStr) > 0 {
		r, err := d.DB.DoExec(link, fmt.Sprintf("INSERT ALL %s SELECT * FROM DUAL", strings.Join(intoStr, " ")), params...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.lastResult = r
			batchResult.rowsAffected += n
		}
	}
	return batchResult, nil
}
