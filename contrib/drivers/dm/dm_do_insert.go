// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const (
	// ctxKeyReturningColumn marks an insert whose identity column DoCommit returns.
	ctxKeyReturningColumn gctx.StrKey = "dm_returning_column"

	// fieldExtraIdentity is the TableField.Extra of an identity column, as MySQL reports it for
	// an AUTO_INCREMENT column.
	fieldExtraIdentity = "auto_increment"
)

// DoInsert inserts or updates data for given table.
// The list parameter must contain at least one record, which was previously validated.
func (d *Driver) DoInsert(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	switch option.InsertOption {
	case gdb.InsertOptionSave:
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionReplace:
		// dm does not support REPLACE INTO syntax, use SAVE instead.
		return d.doSave(ctx, link, table, list, option)

	case gdb.InsertOptionIgnore:
		// dm does not support INSERT IGNORE syntax, use MERGE instead.
		return d.doInsertIgnore(ctx, link, table, list, option)

	default:
		return d.doInsertDefault(ctx, link, table, list, option)
	}
}

// doInsertDefault inserts `list` with the default insert option, so that the last insert id is
// the value of the identity column of the last record, as MySQL reports the value of its
// AUTO_INCREMENT column: the value the record gives, or else the value the server generates,
// which RETURNING gets for a record giving a raw expression or 0, as the underlying driver
// reports the ROWID of such a record. It is 0 for a table without identity column.
func (d *Driver) doInsertDefault(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (sql.Result, error) {
	identity, err := d.identityColumn(ctx, table)
	if err != nil {
		return nil, err
	}
	if identity == "" {
		result, err := d.Core.DoInsert(ctx, link, table, list, option)
		if err != nil {
			return result, err
		}
		return &Result{Result: result}, nil
	}
	value, ok := recordValueOfField(list[len(list)-1], identity)
	if !ok || value == nil {
		return d.Core.DoInsert(ctx, link, table, list, option)
	}
	lastInsertId, generated, err := identityValue(value, identity)
	if err != nil {
		return nil, err
	}
	if generated {
		return d.doInsertReturning(ctx, link, table, list, option, identity)
	}
	result, err := d.Core.DoInsert(ctx, link, table, list, option)
	if err != nil {
		return result, err
	}
	return &Result{Result: result, lastInsertId: lastInsertId}, nil
}

// identityValue returns the value `value` that a record gives the identity column `identity`, or
// whether the server generates the value instead, as it does for a raw expression and for 0.
func identityValue(value any, identity string) (id int64, generated bool, err error) {
	if _, ok := value.(gdb.Raw); ok {
		return 0, true, nil
	}
	if id, err = strconv.ParseInt(gconv.String(value), 10, 64); err != nil {
		return 0, false, gerror.WrapCodef(
			gcode.CodeInvalidParameter, err, `invalid value "%v" of identity column "%s"`, value, identity,
		)
	}
	return id, id == 0, nil
}

// doInsertReturning inserts `list` with the last record in a statement of its own, which returns
// the value of the identity column `identity` as the last insert id.
func (d *Driver) doInsertReturning(
	ctx context.Context, link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, identity string,
) (sql.Result, error) {
	var affected int64
	option.InsertOption = gdb.InsertOptionDefault
	if len(list) > 1 {
		result, err := d.Core.DoInsert(ctx, link, table, list[:len(list)-1], option)
		if err != nil {
			return result, err
		}
		if affected, err = result.RowsAffected(); err != nil {
			return result, err
		}
	}
	ctx = context.WithValue(ctx, ctxKeyReturningColumn, identity)
	result, err := d.Core.DoInsert(ctx, link, table, list[len(list)-1:], option)
	if err != nil {
		return result, err
	}
	n, err := result.RowsAffected()
	if err != nil {
		return result, err
	}
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		return result, err
	}
	return &Result{Result: &gdb.SqlResult{Result: result, Affected: affected + n}, lastInsertId: lastInsertId}, nil
}

// identityColumn returns the name of the identity column of `table`, or empty if it has none.
func (d *Driver) identityColumn(ctx context.Context, table string) (string, error) {
	fields, err := d.GetCore().GetDB().TableFields(ctx, d.GetTableNameForFields(table))
	if err != nil {
		return "", err
	}
	for _, field := range fields {
		if field.Extra == fieldExtraIdentity {
			return field.Name, nil
		}
	}
	return "", nil
}

// recordValueOfField returns the value of the field `name` in `record`, whose keys may differ
// from the field name in letter case.
func recordValueOfField(record gdb.Map, name string) (any, bool) {
	if value, ok := record[name]; ok {
		return value, true
	}
	for key, value := range record {
		if strings.EqualFold(key, name) {
			return value, true
		}
	}
	return nil, false
}

// tableWithoutAlias returns `table` without the alias that it may give the table, like
// "user u" or "user AS u", which a MERGE statement does not accept.
func tableWithoutAlias(table string) string {
	words := strings.Fields(table)
	if len(words) == 2 || len(words) == 3 && strings.EqualFold(words[1], "AS") {
		return words[0]
	}
	return table
}

// doSave support upsert for dm
func (d *Driver) doSave(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, true)
}

// doInsertIgnore implements INSERT IGNORE operation using MERGE statement for DM database.
// It only inserts records when there's no conflict on primary/unique keys.
func (d *Driver) doInsertIgnore(ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption,
) (result sql.Result, err error) {
	return d.doMergeInsert(ctx, link, table, list, option, false)
}

// doMergeInsert implements MERGE-based insert operations for DM database, with a MERGE statement
// for each record, as MySQL inserts or updates the records one after another. The last insert id
// is the value of the identity column of the last record, as for doInsertDefault: a record that
// leaves the identity to the server, being the conflict key, matches no record and is inserted
// with RETURNING.
// When withUpdate is true, it performs upsert (insert or update).
// When withUpdate is false, it performs insert ignore (insert only when no conflict).
func (d *Driver) doMergeInsert(
	ctx context.Context,
	link gdb.Link, table string, list gdb.List, option gdb.DoInsertOption, withUpdate bool,
) (result sql.Result, err error) {
	table = tableWithoutAlias(table)
	// If OnConflict is not specified, automatically get the primary key of the table
	conflictKeys := option.OnConflict
	if len(conflictKeys) == 0 {
		primaryKeys, err := d.Core.GetPrimaryKeys(ctx, table)
		if err != nil {
			return nil, gerror.WrapCode(
				gcode.CodeInternalError,
				err,
				`failed to get primary keys for table`,
			)
		}
		foundPrimaryKey := false
		for _, primaryKey := range primaryKeys {
			for dataKey := range list[0] {
				if strings.EqualFold(dataKey, primaryKey) {
					foundPrimaryKey = true
					break
				}
			}
			if foundPrimaryKey {
				break
			}
		}
		if !foundPrimaryKey {
			return nil, gerror.NewCodef(
				gcode.CodeMissingParameter,
				`Replace/Save/InsertIgnore operation requires conflict detection: `+
					`either specify OnConflict() columns or ensure table '%s' has a primary key in the data`,
				table,
			)
		}
		sort.Strings(primaryKeys)
		conflictKeys = primaryKeys
	}

	var (
		charL, charR   = d.GetChars()
		conflictKeySet = gset.NewStrSet(false)
		batchResult    = new(gdb.SqlResult)
	)

	// conflictKeys slice type conv to set type
	for _, conflictKey := range conflictKeys {
		conflictKeySet.Add(gstr.ToUpper(conflictKey))
	}

	identity, err := d.identityColumn(ctx, table)
	if err != nil {
		return nil, err
	}
	var (
		lastInsertId   int64
		mergesIdentity = identity != "" && len(conflictKeys) == 1 && strings.EqualFold(conflictKeys[0], identity)
	)
	for _, one := range list {
		if identity != "" {
			value, _ := recordValueOfField(one, identity)
			if _, isRaw := value.(gdb.Raw); mergesIdentity && !isRaw && (value == nil || gconv.Int64(value) == 0) {
				// The record matches no record by an identity that it leaves to the server to generate.
				r, err := d.doInsertReturning(ctx, link, table, gdb.List{one}, option, identity)
				if err != nil {
					return r, err
				}
				if lastInsertId, err = r.LastInsertId(); err != nil {
					return r, err
				}
				n, err := r.RowsAffected()
				if err != nil {
					return r, err
				}
				batchResult.Result = r
				batchResult.Affected += n
				continue
			}
			lastInsertId = gconv.Int64(value)
		}
		var (
			oneLen = len(one)
			keys   = make([]string, 0, oneLen)

			// queryHolders:	Handle data with Holder that need to be merged
			// queryValues:		Handle data that need to be merged
			// insertKeys:		Handle valid keys that need to be inserted
			// insertValues:	Handle values that need to be inserted
			// updateValues:	Handle values that need to be updated (only when withUpdate=true)
			queryHolders = make([]string, oneLen)
			queryValues  = make([]any, 0, oneLen)
			insertKeys   = make([]string, oneLen)
			insertValues = make([]string, oneLen)
			updateValues []string
		)
		for key := range one {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for index, key := range keys {
			keyWithChar := charL + key + charR
			if s, ok := one[key].(gdb.Raw); ok {
				queryHolders[index] = fmt.Sprintf("%s AS %s", gconv.String(s), keyWithChar)
			} else {
				queryHolders[index] = fmt.Sprintf("? AS %s", keyWithChar)
				queryValues = append(queryValues, one[key])
			}
			insertKeys[index] = keyWithChar
			insertValues[index] = fmt.Sprintf("T2.%s", keyWithChar)
		}
		// Build updateValues only when withUpdate is true
		if withUpdate {
			updateValues = d.formatMergeUpdateValues(keys, conflictKeySet, option)
		}

		sqlStr := parseSqlForMerge(
			d.QuotePrefixTableName(table), queryHolders, insertKeys, insertValues, updateValues, conflictKeys,
		)
		r, err := d.DoExec(ctx, link, sqlStr, queryValues...)
		if err != nil {
			return r, err
		}
		if n, err := r.RowsAffected(); err != nil {
			return r, err
		} else {
			batchResult.Result = r
			batchResult.Affected += n
		}
	}
	return &Result{Result: batchResult, lastInsertId: lastInsertId}, nil
}

// formatMergeUpdateValues returns the assignments of the MERGE UPDATE SET clause for a record
// with the given `keys`. It follows OnDuplicate/OnDuplicateEx of `option` if specified, or else
// updates every key except conflict keys and, unless the whole record is replaced, soft created
// fields. A conflict key assigned from its own source column is left out, as it cannot change
// the matched row.
func (d *Driver) formatMergeUpdateValues(
	keys []string, conflictKeySet *gset.StrSet, option gdb.DoInsertOption,
) (updateValues []string) {
	charL, charR := d.GetChars()
	if option.OnDuplicateStr != "" {
		return []string{option.OnDuplicateStr}
	}
	if len(option.OnDuplicateMap) > 0 {
		updateKeys := make([]string, 0, len(option.OnDuplicateMap))
		for key := range option.OnDuplicateMap {
			updateKeys = append(updateKeys, key)
		}
		sort.Strings(updateKeys)
		for _, key := range updateKeys {
			keyWithChar := charL + key + charR
			switch value := option.OnDuplicateMap[key].(type) {
			case gdb.Raw, *gdb.Raw:
				updateValues = append(updateValues, fmt.Sprintf(`T1.%s = %s`, keyWithChar, gconv.String(value)))

			case gdb.Counter, *gdb.Counter:
				var counter gdb.Counter
				switch v := value.(type) {
				case gdb.Counter:
					counter = v
				case *gdb.Counter:
					counter = *v
				}
				operator, columnVal := "+", counter.Value
				if columnVal < 0 {
					operator, columnVal = "-", -columnVal
				}
				updateValues = append(updateValues, fmt.Sprintf(
					`T1.%s = T1.%s%s%s`,
					keyWithChar, charL+counter.Field+charR, operator, gconv.String(columnVal),
				))

			default:
				column := gconv.String(value)
				if conflictKeySet.Contains(gstr.ToUpper(key)) && strings.EqualFold(key, column) {
					continue
				}
				updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, charL+column+charR))
			}
		}
		return updateValues
	}
	// Filter conflict keys, and soft created fields unless the whole record is replaced.
	for _, key := range keys {
		if conflictKeySet.Contains(gstr.ToUpper(key)) {
			continue
		}
		if option.InsertOption != gdb.InsertOptionReplace && d.IsSoftCreatedFieldName(key) {
			continue
		}
		keyWithChar := charL + key + charR
		updateValues = append(updateValues, fmt.Sprintf(`T1.%s = T2.%s`, keyWithChar, keyWithChar))
	}
	return updateValues
}

// parseSqlForMerge generates MERGE statement for DM database.
// When updateValues is empty, it only inserts (INSERT IGNORE behavior).
// When updateValues is provided, it performs upsert (INSERT or UPDATE).
// Examples:
// - INSERT IGNORE: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...)
// - UPSERT: MERGE INTO table T1 USING (...) T2 ON (...) WHEN NOT MATCHED THEN INSERT(...) VALUES (...) WHEN MATCHED THEN UPDATE SET ...
func parseSqlForMerge(table string,
	queryHolders, insertKeys, insertValues, updateValues, duplicateKey []string,
) (sqlStr string) {
	var (
		queryHolderStr  = strings.Join(queryHolders, ",")
		insertKeyStr    = strings.Join(insertKeys, ",")
		insertValueStr  = strings.Join(insertValues, ",")
		duplicateKeyStr string
	)

	// Build ON condition
	for index, keys := range duplicateKey {
		if index != 0 {
			duplicateKeyStr += " AND "
		}
		duplicateKeyStr += fmt.Sprintf("T1.%s = T2.%s", keys, keys)
	}

	// Build SQL based on whether UPDATE is needed
	pattern := gstr.Trim(`MERGE INTO %s T1 USING (SELECT %s FROM DUAL) T2 ON (%s) WHEN NOT MATCHED THEN INSERT(%s) VALUES (%s)`)
	if len(updateValues) > 0 {
		// Upsert: INSERT or UPDATE
		pattern += ` WHEN MATCHED THEN UPDATE SET %s`
		return fmt.Sprintf(
			pattern, table, queryHolderStr, duplicateKeyStr, insertKeyStr, insertValueStr,
			strings.Join(updateValues, ","),
		)
	}
	// Insert Ignore: INSERT only
	return fmt.Sprintf(pattern, table, queryHolderStr, duplicateKeyStr, insertKeyStr, insertValueStr)
}
