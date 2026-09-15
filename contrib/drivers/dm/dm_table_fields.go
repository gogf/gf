// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/util/gutil"
)

// escapeSingleQuote escapes single quotes in the string to prevent SQL injection.
// In SQL, single quotes are escaped by doubling them (two single quotes).
func escapeSingleQuote(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

// lookupRecordValue returns the value of the given column from the record, ignoring
// the case of the column name, together with whether the column is present.
//
// The underlying DM driver accepts the DSN option `columnNameCase=lower|upper`,
// which rewrites the case of every column name of a result set. As the metadata
// queries of this package are written with uppercase column names, a plain map
// lookup returns nothing once the option is set to `lower`. Looking the column up
// case-insensitively keeps the metadata readable for any value of that option.
func lookupRecordValue(record gdb.Record, column string) (gdb.Value, bool) {
	if value, ok := record[column]; ok {
		return value, true
	}
	for name, value := range record {
		if strings.EqualFold(name, column) {
			return value, true
		}
	}
	return nil, false
}

// recordValue returns the value of the given column from the record, ignoring the
// case of the column name. It returns a nil value when the column is absent, which
// is still safe to read from.
func recordValue(record gdb.Record, column string) gdb.Value {
	value, _ := lookupRecordValue(record, column)
	return value
}

const (
	tableFieldsSqlTmp         = `SELECT c.COLUMN_NAME, c.DATA_TYPE, c.DATA_LENGTH, c.DATA_DEFAULT, c.NULLABLE, cc.COMMENTS FROM ALL_TAB_COLUMNS c LEFT JOIN ALL_COL_COMMENTS cc ON c.COLUMN_NAME = cc.COLUMN_NAME AND c.TABLE_NAME = cc.TABLE_NAME AND c.OWNER = cc.OWNER WHERE c.TABLE_NAME = '%s' AND c.OWNER = '%s'`
	tableFieldsPkSqlSchemaTmp = `SELECT COLS.COLUMN_NAME AS PRIMARY_KEY_COLUMN FROM USER_CONSTRAINTS CONS JOIN USER_CONS_COLUMNS COLS ON CONS.CONSTRAINT_NAME = COLS.CONSTRAINT_NAME WHERE CONS.TABLE_NAME = '%s' AND CONS.CONSTRAINT_TYPE = 'P'`
	tableFieldsPkSqlDBATmp    = `SELECT COLS.COLUMN_NAME AS PRIMARY_KEY_COLUMN FROM DBA_CONSTRAINTS CONS JOIN DBA_CONS_COLUMNS COLS ON CONS.CONSTRAINT_NAME = COLS.CONSTRAINT_NAME WHERE CONS.TABLE_NAME = '%s' AND CONS.OWNER = '%s' AND CONS.CONSTRAINT_TYPE = 'P'`
)

// tableNameCandidatesForMetadata returns table name candidates for DM metadata queries.
func tableNameCandidatesForMetadata(table string) []string {
	table = strings.Trim(table, quoteChar)
	upperTable := strings.ToUpper(table)
	if upperTable == table {
		return []string{table}
	}
	return []string{upperTable, table}
}

// schemaNameForMetadata returns the schema owner name used in DM metadata queries.
func schemaNameForMetadata(schema string) string {
	return strings.ToUpper(strings.Trim(schema, quoteChar))
}

// TableFields retrieves and returns the fields' information of specified table of current schema.
func (d *Driver) TableFields(
	ctx context.Context, table string, schema ...string,
) (fields map[string]*gdb.TableField, err error) {
	var (
		result   gdb.Result
		pkResult gdb.Result
		link     gdb.Link
		// When no schema is specified, the configuration item is returned by default
		usedSchema        = gutil.GetOrDefaultStr(d.GetSchema(), schema...)
		usedMetadataTable string
		usedMetadataOwner = schemaNameForMetadata(usedSchema)
	)
	// When usedSchema is empty, return the default link
	if link, err = d.SlaveLink(usedSchema); err != nil {
		return nil, err
	}
	// The link has been distinguished and no longer needs to judge the owner
	for _, candidate := range tableNameCandidatesForMetadata(table) {
		usedMetadataTable = candidate
		result, err = d.DoSelect(
			ctx, link,
			fmt.Sprintf(
				tableFieldsSqlTmp,
				escapeSingleQuote(usedMetadataTable),
				escapeSingleQuote(usedMetadataOwner),
			),
		)
		if err != nil {
			return nil, err
		}
		if !result.IsEmpty() {
			break
		}
	}
	// Query the primary key field
	pkResult, err = d.DoSelect(
		ctx, link,
		fmt.Sprintf(tableFieldsPkSqlSchemaTmp, escapeSingleQuote(usedMetadataTable)),
	)
	if err != nil {
		return nil, err
	}
	if pkResult.IsEmpty() {
		pkResult, err = d.DoSelect(
			ctx, link,
			fmt.Sprintf(tableFieldsPkSqlDBATmp, escapeSingleQuote(usedMetadataTable), escapeSingleQuote(usedMetadataOwner)),
		)
		if err != nil {
			return nil, err
		}
	}
	fields = make(map[string]*gdb.TableField)
	pkFields := gmap.NewStrStrMap()
	for _, pk := range pkResult {
		pkFields.Set(recordValue(pk, "PRIMARY_KEY_COLUMN").String(), "PRI")
	}
	for i, m := range result {
		// NULLABLE returns "N" "Y"
		// "N" means not null
		// "Y" means could be null
		var nullable bool
		if recordValue(m, "NULLABLE").String() != "N" {
			nullable = true
		}

		// Build field type with length/precision
		// For NUMBER(p,s): use DATA_PRECISION and DATA_SCALE
		// For VARCHAR2/CHAR: use DATA_LENGTH
		var (
			fieldType  string
			columnName = recordValue(m, "COLUMN_NAME").String()
			dataType   = recordValue(m, "DATA_TYPE").String()
			dataLength = recordValue(m, "DATA_LENGTH").Int()
		)
		if dataLength > 0 {
			fieldType = fmt.Sprintf("%s(%d)", dataType, dataLength)
		} else {
			fieldType = dataType
		}
		fields[columnName] = &gdb.TableField{
			Index:   i,
			Name:    columnName,
			Type:    fieldType,
			Null:    nullable,
			Default: recordValue(m, "DATA_DEFAULT").Val(),
			Key:     pkFields.Get(columnName),
			// Extra:   recordValue(m, "Extra").String(),
			Comment: recordValue(m, "COMMENTS").String(),
		}
	}
	return fields, nil
}
