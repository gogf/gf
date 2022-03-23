// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/crypto/gmd5"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// All does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) All(where ...interface{}) (Result, error) {
	return m.doGetAll(false, where...)
}

// doGetAll does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) doGetAll(limit1 bool, where ...interface{}) (Result, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).All()
	}
	sqlWithHolder, holderArgs := m.getFormattedSqlAndArgs(queryTypeNormal, limit1)
	return m.doGetAllBySql(queryTypeNormal, sqlWithHolder, holderArgs...)
}

// getFieldsFiltered checks the fields and fieldsEx attributes, filters and returns the fields that will
// really be committed to underlying database driver.
func (m *Model) getFieldsFiltered() string {
	if m.fieldsEx == "" {
		// No filtering.
		if !gstr.Contains(m.fields, ".") && !gstr.Contains(m.fields, " ") {
			return m.db.GetCore().QuoteString(m.fields)
		}
		return m.fields
	}
	var (
		fieldsArray []string
		fieldsExSet = gset.NewStrSetFrom(gstr.SplitAndTrim(m.fieldsEx, ","))
	)
	if m.fields != "*" {
		// Filter custom fields with fieldEx.
		fieldsArray = make([]string, 0, 8)
		for _, v := range gstr.SplitAndTrim(m.fields, ",") {
			fieldsArray = append(fieldsArray, v[gstr.PosR(v, "-")+1:])
		}
	} else {
		if gstr.Contains(m.tables, " ") {
			panic("function FieldsEx supports only single table operations")
		}
		// Filter table fields with fieldEx.
		tableFields, err := m.TableFields(m.tablesInit)
		if err != nil {
			panic(err)
		}
		if len(tableFields) == 0 {
			panic(fmt.Sprintf(`empty table fields for table "%s"`, m.tables))
		}
		fieldsArray = make([]string, len(tableFields))
		for k, v := range tableFields {
			fieldsArray[v.Index] = k
		}
	}
	newFields := ""
	for _, k := range fieldsArray {
		if fieldsExSet.Contains(k) {
			continue
		}
		if len(newFields) > 0 {
			newFields += ","
		}
		newFields += m.db.GetCore().QuoteWord(k)
	}
	return newFields
}

// Chunk iterates the query result with given `size` and `handler` function.
func (m *Model) Chunk(size int, handler ChunkHandler) {
	page := m.start
	if page <= 0 {
		page = 1
	}
	model := m
	for {
		model = model.Page(page, size)
		data, err := model.All()
		if err != nil {
			handler(nil, err)
			break
		}
		if len(data) == 0 {
			break
		}
		if handler(data, err) == false {
			break
		}
		if len(data) < size {
			break
		}
		page++
	}
}

// One retrieves one record from table and returns the result as map type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) One(where ...interface{}) (Record, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).One()
	}
	all, err := m.doGetAll(true)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		return all[0], nil
	}
	return nil, nil
}

// Value retrieves a specified record value from table and returns the result as interface type.
// It returns nil if there's no record found with the given conditions from table.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Value(fieldsAndWhere ...interface{}) (Value, error) {
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Value()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Value()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Value()
		}
	}
	one, err := m.One()
	if err != nil {
		return gvar.New(nil), err
	}
	for _, v := range one {
		return v, nil
	}
	return gvar.New(nil), nil
}

// Array queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Array(fieldsAndWhere ...interface{}) ([]Value, error) {
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Array()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Array()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Array()
		}
	}
	all, err := m.All()
	if err != nil {
		return nil, err
	}
	return all.Array(), nil
}

// Struct retrieves one record from table and converts it into given struct.
// The parameter `pointer` should be type of *struct/**struct. If type **struct is given,
// it can create the struct internally during converting.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// user := new(User)
// err  := db.Model("user").Where("id", 1).Scan(user)
//
// user := (*User)(nil)
// err  := db.Model("user").Where("id", 1).Scan(&user).
func (m *Model) doStruct(pointer interface{}, where ...interface{}) error {
	model := m
	// Auto selecting fields by struct attributes.
	if model.fieldsEx == "" && (model.fields == "" || model.fields == "*") {
		if v, ok := pointer.(reflect.Value); ok {
			model = m.Fields(v.Interface())
		} else {
			model = m.Fields(pointer)
		}
	}
	one, err := model.One(where...)
	if err != nil {
		return err
	}
	if err = one.Struct(pointer); err != nil {
		return err
	}
	return model.doWithScanStruct(pointer)
}

// Structs retrieves records from table and converts them into given struct slice.
// The parameter `pointer` should be type of *[]struct/*[]*struct. It can create and fill the struct
// slice internally during converting.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// users := ([]User)(nil)
// err   := db.Model("user").Scan(&users)
//
// users := ([]*User)(nil)
// err   := db.Model("user").Scan(&users).
func (m *Model) doStructs(pointer interface{}, where ...interface{}) error {
	model := m
	// Auto selecting fields by struct attributes.
	if model.fieldsEx == "" && (model.fields == "" || model.fields == "*") {
		if v, ok := pointer.(reflect.Value); ok {
			model = m.Fields(
				reflect.New(
					v.Type().Elem(),
				).Interface(),
			)
		} else {
			model = m.Fields(
				reflect.New(
					reflect.ValueOf(pointer).Elem().Type().Elem(),
				).Interface(),
			)
		}
	}
	all, err := model.All(where...)
	if err != nil {
		return err
	}
	if err = all.Structs(pointer); err != nil {
		return err
	}
	return model.doWithScanStructs(pointer)
}

// Scan automatically calls Struct or Structs function according to the type of parameter `pointer`.
// It calls function doStruct if `pointer` is type of *struct/**struct.
// It calls function doStructs if `pointer` is type of *[]struct/*[]*struct.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,  see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// Example:
// user := new(User)
// err  := db.Model("user").Where("id", 1).Scan(user)
//
// user := (*User)(nil)
// err  := db.Model("user").Where("id", 1).Scan(&user)
//
// users := ([]User)(nil)
// err   := db.Model("user").Scan(&users)
//
// users := ([]*User)(nil)
// err   := db.Model("user").Scan(&users).
func (m *Model) Scan(pointer interface{}, where ...interface{}) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for function Scan should type of pointer`,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Slice, reflect.Array:
		return m.doStructs(pointer, where...)

	case reflect.Struct, reflect.Invalid:
		return m.doStruct(pointer, where...)

	default:
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`element of parameter "pointer" for function Scan should type of struct/*struct/[]struct/[]*struct`,
		)
	}
}

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
//
// See Result.ScanList.
func (m *Model) ScanList(structSlicePointer interface{}, bindToAttrName string, relationAttrNameAndFields ...string) (err error) {
	var result Result
	out, err := checkGetSliceElementInfoForScanList(structSlicePointer, bindToAttrName)
	if err != nil {
		return err
	}
	if m.fields != defaultFields || m.fieldsEx != "" {
		// There are custom fields.
		result, err = m.All()
	} else {
		// Filter fields using temporary created struct using reflect.New.
		result, err = m.Fields(reflect.New(out.BindToAttrType).Interface()).All()
	}
	if err != nil {
		return err
	}
	var (
		relationAttrName string
		relationFields   string
	)
	switch len(relationAttrNameAndFields) {
	case 2:
		relationAttrName = relationAttrNameAndFields[0]
		relationFields = relationAttrNameAndFields[1]
	case 1:
		relationFields = relationAttrNameAndFields[0]
	}
	return doScanList(doScanListInput{
		Model:              m,
		Result:             result,
		StructSlicePointer: structSlicePointer,
		StructSliceValue:   out.SliceReflectValue,
		BindToAttrName:     bindToAttrName,
		RelationAttrName:   relationAttrName,
		RelationFields:     relationFields,
	})
}

// Count does "SELECT COUNT(x) FROM ..." statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Count(where ...interface{}) (int, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Count()
	}
	var (
		sqlWithHolder, holderArgs = m.getFormattedSqlAndArgs(queryTypeCount, false)
		list, err                 = m.doGetAllBySql(queryTypeCount, sqlWithHolder, holderArgs...)
	)
	if err != nil {
		return 0, err
	}
	if len(list) > 0 {
		for _, v := range list[0] {
			return v.Int(), nil
		}
	}
	return 0, nil
}

// CountColumn does "SELECT COUNT(x) FROM ..." statement for the model.
func (m *Model) CountColumn(column string) (int, error) {
	if len(column) == 0 {
		return 0, nil
	}
	return m.Fields(column).Count()
}

// Min does "SELECT MIN(x) FROM ..." statement for the model.
func (m *Model) Min(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`MIN(%s)`, m.QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Max does "SELECT MAX(x) FROM ..." statement for the model.
func (m *Model) Max(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`MAX(%s)`, m.QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Avg does "SELECT AVG(x) FROM ..." statement for the model.
func (m *Model) Avg(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`AVG(%s)`, m.QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Sum does "SELECT SUM(x) FROM ..." statement for the model.
func (m *Model) Sum(column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	value, err := m.Fields(fmt.Sprintf(`SUM(%s)`, m.QuoteWord(column))).Value()
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Union does "(SELECT xxx FROM xxx) UNION (SELECT xxx FROM xxx) ..." statement for the model.
func (m *Model) Union(unions ...*Model) *Model {
	return m.db.Union(unions...)
}

// UnionAll does "(SELECT xxx FROM xxx) UNION ALL (SELECT xxx FROM xxx) ..." statement for the model.
func (m *Model) UnionAll(unions ...*Model) *Model {
	return m.db.UnionAll(unions...)
}

// Limit sets the "LIMIT" statement for the model.
// The parameter `limit` can be either one or two number, if passed two number is passed,
// it then sets "LIMIT limit[0],limit[1]" statement for the model, or else it sets "LIMIT limit[0]"
// statement.
func (m *Model) Limit(limit ...int) *Model {
	model := m.getModel()
	switch len(limit) {
	case 1:
		model.limit = limit[0]
	case 2:
		model.start = limit[0]
		model.limit = limit[1]
	}
	return model
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
func (m *Model) Offset(offset int) *Model {
	model := m.getModel()
	model.offset = offset
	return model
}

// Distinct forces the query to only return distinct results.
func (m *Model) Distinct() *Model {
	model := m.getModel()
	model.distinct = "DISTINCT "
	return model
}

// Page sets the paging number for the model.
// The parameter `page` is started from 1 for paging.
// Note that, it differs that the Limit function starts from 0 for "LIMIT" statement.
func (m *Model) Page(page, limit int) *Model {
	model := m.getModel()
	if page <= 0 {
		page = 1
	}
	model.start = (page - 1) * limit
	model.limit = limit
	return model
}

// Having sets the having statement for the model.
// The parameters of this function usage are as the same as function Where.
// See Where.
func (m *Model) Having(having interface{}, args ...interface{}) *Model {
	model := m.getModel()
	model.having = []interface{}{
		having, args,
	}
	return model
}

// doGetAllBySql does the select statement on the database.
func (m *Model) doGetAllBySql(queryType int, sql string, args ...interface{}) (result Result, err error) {
	var (
		ok       bool
		ctx      = m.GetCtx()
		cacheKey = ""
		cacheObj = m.db.GetCache()
	)
	// Retrieve from cache.
	if m.cacheEnabled && m.tx == nil {
		cacheKey = m.cacheOption.Name
		if len(cacheKey) == 0 {
			cacheKey = fmt.Sprintf(
				`GCache@Schema(%s):%s`,
				m.db.GetSchema(),
				gmd5.MustEncryptString(sql+", @PARAMS:"+gconv.String(args)),
			)
		}
		if v, _ := cacheObj.Get(ctx, cacheKey); !v.IsNil() {
			if result, ok = v.Val().(Result); ok {
				// In-memory cache.
				return result, nil
			}
			// Other cache, it needs conversion.
			if err = json.UnmarshalUseNumber(v.Bytes(), &result); err != nil {
				return nil, err
			} else {
				return result, nil
			}
		}
	}

	in := &HookSelectInput{
		internalParamHookSelect: internalParamHookSelect{
			internalParamHook: internalParamHook{
				link:  m.getLink(false),
				model: m,
			},
			handler: m.hookHandler.Select,
		},
		Table:            m.tables,
		Sql:              sql,
		Args:             m.mergeArguments(args),
		IsCountStatement: queryType == queryTypeCount,
	}
	result, err = in.Next(m.GetCtx())

	// Cache the result.
	if cacheKey != "" && err == nil {
		if m.cacheOption.Duration < 0 {
			if _, errCache := cacheObj.Remove(ctx, cacheKey); errCache != nil {
				intlog.Errorf(m.GetCtx(), `%+v`, errCache)
			}
		} else {
			// In case of Cache Penetration.
			if result.IsEmpty() && m.cacheOption.Force {
				result = Result{}
			}
			if errCache := cacheObj.Set(ctx, cacheKey, result, m.cacheOption.Duration); errCache != nil {
				intlog.Errorf(m.GetCtx(), `%+v`, errCache)
			}
		}
	}
	return result, err
}

func (m *Model) getFormattedSqlAndArgs(queryType int, limit1 bool) (sqlWithHolder string, holderArgs []interface{}) {
	switch queryType {
	case queryTypeCount:
		countFields := "COUNT(1)"
		if m.fields != "" && m.fields != "*" {
			// DO NOT quote the m.fields here, in case of fields like:
			// DISTINCT t.user_id uid
			countFields = fmt.Sprintf(`COUNT(%s%s)`, m.distinct, m.fields)
		}
		// Raw SQL Model.
		if m.rawSql != "" {
			sqlWithHolder = fmt.Sprintf("SELECT %s FROM (%s) AS T", countFields, m.rawSql)
			return sqlWithHolder, nil
		}
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(false, true)
		sqlWithHolder = fmt.Sprintf("SELECT %s FROM %s%s", countFields, m.tables, conditionWhere+conditionExtra)
		if len(m.groupBy) > 0 {
			sqlWithHolder = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", sqlWithHolder)
		}
		return sqlWithHolder, conditionArgs

	default:
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(limit1, false)
		// Raw SQL Model, especially for UNION/UNION ALL featured SQL.
		if m.rawSql != "" {
			sqlWithHolder = fmt.Sprintf(
				"%s%s",
				m.rawSql,
				conditionWhere+conditionExtra,
			)
			return sqlWithHolder, conditionArgs
		}
		// DO NOT quote the m.fields where, in case of fields like:
		// DISTINCT t.user_id uid
		sqlWithHolder = fmt.Sprintf(
			"SELECT %s%s FROM %s%s",
			m.distinct,
			m.getFieldsFiltered(),
			m.tables,
			conditionWhere+conditionExtra,
		)
		return sqlWithHolder, conditionArgs
	}
}

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the `m`.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(limit1 bool, isCountStatement bool) (conditionWhere string, conditionExtra string, conditionArgs []interface{}) {
	autoPrefix := ""
	if gstr.Contains(m.tables, " JOIN ") {
		autoPrefix = m.db.GetCore().QuoteWord(
			m.db.GetCore().guessPrimaryTableName(m.tablesInit),
		)
	}
	var (
		tableForMappingAndFiltering = m.tables
	)
	if len(m.whereHolder) > 0 {
		for _, holder := range m.whereHolder {
			tableForMappingAndFiltering = m.tables
			if holder.Prefix == "" {
				holder.Prefix = autoPrefix
			}

			switch holder.Operator {
			case whereHolderOperatorWhere:
				if conditionWhere == "" {
					newWhere, newArgs := formatWhereHolder(m.db, formatWhereHolderInput{
						ModelWhereHolder: holder,
						OmitNil:          m.option&optionOmitNilWhere > 0,
						OmitEmpty:        m.option&optionOmitEmptyWhere > 0,
						Schema:           m.schema,
						Table:            tableForMappingAndFiltering,
					})
					if len(newWhere) > 0 {
						conditionWhere = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case whereHolderOperatorAnd:
				newWhere, newArgs := formatWhereHolder(m.db, formatWhereHolderInput{
					ModelWhereHolder: holder,
					OmitNil:          m.option&optionOmitNilWhere > 0,
					OmitEmpty:        m.option&optionOmitEmptyWhere > 0,
					Schema:           m.schema,
					Table:            tableForMappingAndFiltering,
				})
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s AND (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) AND (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}

			case whereHolderOperatorOr:
				newWhere, newArgs := formatWhereHolder(m.db, formatWhereHolderInput{
					ModelWhereHolder: holder,
					OmitNil:          m.option&optionOmitNilWhere > 0,
					OmitEmpty:        m.option&optionOmitEmptyWhere > 0,
					Schema:           m.schema,
					Table:            tableForMappingAndFiltering,
				})
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s OR (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) OR (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}
			}
		}
	}
	// Soft deletion.
	softDeletingCondition := m.getConditionForSoftDeleting()
	if m.rawSql != "" && conditionWhere != "" {
		if gstr.ContainsI(m.rawSql, " WHERE ") {
			conditionWhere = " AND " + conditionWhere
		} else {
			conditionWhere = " WHERE " + conditionWhere
		}
	} else if !m.unscoped && softDeletingCondition != "" {
		if conditionWhere == "" {
			conditionWhere = fmt.Sprintf(` WHERE %s`, softDeletingCondition)
		} else {
			conditionWhere = fmt.Sprintf(` WHERE (%s) AND %s`, conditionWhere, softDeletingCondition)
		}
	} else {
		if conditionWhere != "" {
			conditionWhere = " WHERE " + conditionWhere
		}
	}

	// GROUP BY.
	if m.groupBy != "" {
		conditionExtra += " GROUP BY " + m.groupBy
	}
	// HAVING.
	if len(m.having) > 0 {
		havingHolder := ModelWhereHolder{
			Where:  m.having[0],
			Args:   gconv.Interfaces(m.having[1]),
			Prefix: autoPrefix,
		}
		havingStr, havingArgs := formatWhereHolder(m.db, formatWhereHolderInput{
			ModelWhereHolder: havingHolder,
			OmitNil:          m.option&optionOmitNilWhere > 0,
			OmitEmpty:        m.option&optionOmitEmptyWhere > 0,
			Schema:           m.schema,
			Table:            m.tables,
		})
		if len(havingStr) > 0 {
			conditionExtra += " HAVING " + havingStr
			conditionArgs = append(conditionArgs, havingArgs...)
		}
	}
	// ORDER BY.
	if m.orderBy != "" {
		conditionExtra += " ORDER BY " + m.orderBy
	}
	// LIMIT.
	if !isCountStatement {
		if m.limit != 0 {
			if m.start >= 0 {
				conditionExtra += fmt.Sprintf(" LIMIT %d,%d", m.start, m.limit)
			} else {
				conditionExtra += fmt.Sprintf(" LIMIT %d", m.limit)
			}
		} else if limit1 {
			conditionExtra += " LIMIT 1"
		}

		if m.offset >= 0 {
			conditionExtra += fmt.Sprintf(" OFFSET %d", m.offset)
		}
	}

	if m.lockInfo != "" {
		conditionExtra += " " + m.lockInfo
	}
	return
}
