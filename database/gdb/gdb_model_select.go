// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v3/container/gset"
	"github.com/gogf/gf/v3/encoding/gjson"
	"github.com/gogf/gf/v3/errors/gcode"
	"github.com/gogf/gf/v3/errors/gerror"
	"github.com/gogf/gf/v3/internal/reflection"
	"github.com/gogf/gf/v3/text/gstr"
	"github.com/gogf/gf/v3/util/gconv"
)

// All does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) All(ctx context.Context) (Result, error) {
	model := m.callHandlers(ctx)
	return model.doGetAll(ctx, SelectTypeDefault, false)
}

// AllAndCount retrieves all records and the total count of records from the model.
// If useFieldForCount is true, it will use the fields specified in the model for counting;
// otherwise, it will use a constant value of 1 for counting.
// It returns the result as a slice of records, the total count of records, and an error if any.
// The where parameter is an optional list of conditions to use when retrieving records.
//
// Example:
//
//	var model Model
//	var result Result
//	var count int
//	where := []any{"name = ?", "John"}
//	result, count, err := model.AllAndCount(true)
//	if err != nil {
//	    // Handle error.
//	}
//	fmt.Println(result, count)
func (m *Model) AllAndCount(ctx context.Context, useFieldForCount bool) (result Result, totalCount int, err error) {
	var (
		allModel   = m.callHandlers(ctx)
		countModel = m.Clone().callHandlers(ctx)
	)

	// If useFieldForCount is false, set the fields to a constant value of 1 for counting
	if !useFieldForCount {
		countModel.fields = []any{Raw("1")}
	}

	// Get the total count of records
	totalCount, err = countModel.Count(ctx)
	if err != nil {
		return
	}

	// If the total count is 0, there are no records to retrieve, so return early
	if totalCount == 0 {
		return
	}

	// Retrieve all records
	result, err = allModel.doGetAll(ctx, SelectTypeDefault, false)
	return
}

// Chunk iterates the query result with given `size` and `handler` function.
func (m *Model) Chunk(ctx context.Context, size int, handler ChunkHandler) {
	var (
		model = m.callHandlers(ctx)
		page  = model.start
	)
	if page <= 0 {
		page = 1
	}
	for {
		model = model.Clone().Page(page, size)
		data, err := model.All(ctx)
		if err != nil {
			handler(nil, err)
			break
		}
		if len(data) == 0 {
			break
		}
		if !handler(data, err) {
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
func (m *Model) One(ctx context.Context) (Record, error) {
	model := m.callHandlers(ctx)
	all, err := model.doGetAll(ctx, SelectTypeDefault, true)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		return all[0], nil
	}
	return nil, nil
}

// Array queries and returns data values as slice from database.
// Note that if there are multiple columns in the result, it returns just one column values randomly.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Array(ctx context.Context) ([]Value, error) {
	var (
		field string
		model = m.callHandlers(ctx)
		core  = model.db.GetCore()
	)
	ctx = core.injectInternalColumnIntoCtx(ctx)

	all, err := model.doGetAll(ctx, SelectTypeArray, false)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		internalColumn := core.getInternalColumnFromCtx(ctx)
		if internalColumn == nil {
			return nil, gerror.NewCode(
				gcode.CodeInternalError,
				`query count error: the internal context data is missing. there's internal issue should be fixed`,
			)
		}
		// If FirstResultColumn present, it returns the value of the first record of the first field.
		// It means it use no cache mechanism, while cache mechanism makes `internalColumnData` missing.
		field = internalColumn.FirstResultColumn
		if field == "" {
			// Fields number check.
			var recordFields = model.getRecordFields(all[0])
			if len(recordFields) == 1 {
				field = recordFields[0]
			} else {
				// it returns error if there are multiple fields in the result record.
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`invalid fields for "Array" operation, result fields number "%d"%s, but expect one`,
					len(recordFields),
					gjson.MustEncodeString(recordFields),
				)
			}
		}
	}
	return all.Array(field), nil
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
func (m *Model) Scan(ctx context.Context, pointer any) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Ptr {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for function Scan should type of pointer`,
		)
	}

	switch reflectInfo.OriginKind {
	case reflect.Slice, reflect.Array:
		return m.callHandlers(ctx).doStructs(ctx, pointer)

	case reflect.Struct, reflect.Invalid:
		return m.callHandlers(ctx).doStruct(ctx, pointer)

	default:
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`element of parameter "pointer" for function Scan should type of struct/*struct/[]struct/[]*struct`,
		)
	}
}

// ScanAndCount scans a single record or record array that matches the given conditions and counts the total number
// of records that match those conditions.
//
// If `useFieldForCount` is true, it will use the fields specified in the model for counting;
// The `pointer` parameter is a pointer to a struct that the scanned data will be stored in.
// The `totalCount` parameter is a pointer to an integer that will be set to the total number of records that match the given conditions.
// The where parameter is an optional list of conditions to use when retrieving records.
//
// Example:
//
//	var count int
//	user := new(User)
//	err  := db.Model("user").Where("id", 1).ScanAndCount(user,&count,true)
//	fmt.Println(user, count)
//
// Example Join:
//
//	type User struct {
//		Id       int
//		Passport string
//		Name     string
//		Age      int
//	}
//	var users []User
//	var count int
//	db.Model(table).As("u1").
//		LeftJoin(tableName2, "u2", "u2.id=u1.id").
//		Fields("u1.passport,u1.id,u2.name,u2.age").
//		Where("u1.id<2").
//		ScanAndCount(&users, &count, false)
func (m *Model) ScanAndCount(ctx context.Context, pointer any, totalCount *int, useFieldForCount bool) (err error) {
	// support Fields with *, example: .Fields("a.*, b.name"). Count sql is select count(1) from xxx
	var (
		scanModel  = m.callHandlers(ctx)
		countModel = m.Clone().callHandlers(ctx)
	)
	// If useFieldForCount is false, set the fields to a constant value of 1 for counting
	if !useFieldForCount {
		countModel.fields = []any{Raw("1")}
	}

	// Get the total count of records
	*totalCount, err = countModel.Count(ctx)
	if err != nil {
		return err
	}

	// If the total count is 0, there are no records to retrieve, so return early
	if *totalCount == 0 {
		return
	}
	err = scanModel.Scan(ctx, pointer)
	return
}

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
//
// See Result.ScanList.
func (m *Model) ScanList(ctx context.Context, structSlicePointer any, bindToAttrName string, relationAttrNameAndFields ...string) (err error) {
	var (
		result Result
		model  = m.callHandlers(ctx)
	)
	out, err := checkGetSliceElementInfoForScanList(structSlicePointer, bindToAttrName)
	if err != nil {
		return err
	}
	if len(model.fields) > 0 || len(model.fieldsEx) != 0 {
		// There are custom fields.
		result, err = model.All(ctx)
	} else {
		// Filter fields using temporary created struct using reflect.New.
		result, err = model.Fields(reflect.New(out.BindToAttrType).Interface()).All(ctx)
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
	return doScanList(ctx, doScanListInput{
		Model:              model,
		Result:             result,
		StructSlicePointer: structSlicePointer,
		StructSliceValue:   out.SliceReflectValue,
		BindToAttrName:     bindToAttrName,
		RelationAttrName:   relationAttrName,
		RelationFields:     relationFields,
	})
}

// Value retrieves a specified record value from table and returns the result as interface type.
// It returns nil if there's no record found with the given conditions from table.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Value(ctx context.Context) (Value, error) {
	var (
		model = m.callHandlers(ctx)
		core  = model.db.GetCore()
	)
	ctx = core.injectInternalColumnIntoCtx(ctx)

	var (
		sqlWithHolder, holderArgs = model.getFormattedSqlAndArgs(ctx, SelectTypeValue, true)
		all, err                  = model.doGetAllBySql(ctx, SelectTypeValue, sqlWithHolder, holderArgs...)
	)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		internalColumn := core.getInternalColumnFromCtx(ctx)
		if internalColumn == nil {
			return nil, gerror.NewCode(
				gcode.CodeInternalError,
				`query count error: the internal context data is missing. there's internal issue should be fixed`,
			)
		}
		// If `FirstResultColumn` present, it returns the value of the first record of the first field.
		// It means it use no cache mechanism, while cache mechanism makes `internalColumnData` missing.
		if v, ok := all[0][internalColumn.FirstResultColumn]; ok {
			return v, nil
		}
		// Fields number check.
		var recordFields = model.getRecordFields(all[0])
		if len(recordFields) == 1 {
			for _, v := range all[0] {
				return v, nil
			}
		}
		// it returns error if there are multiple fields in the result record.
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid fields for "Value" operation, result fields number "%d"%s, but expect one`,
			len(recordFields),
			gjson.MustEncodeString(recordFields),
		)
	}
	return nil, nil
}

func (m *Model) getRecordFields(record Record) []string {
	if len(record) == 0 {
		return nil
	}
	var fields = make([]string, 0)
	for k := range record {
		fields = append(fields, k)
	}
	return fields
}

// Count does "SELECT COUNT(x) FROM ..." statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Count(ctx context.Context) (int, error) {
	model := m.callHandlers(ctx)

	var core = model.db.GetCore()
	ctx = core.injectInternalColumnIntoCtx(ctx)

	var (
		sqlWithHolder, holderArgs = model.getFormattedSqlAndArgs(ctx, SelectTypeCount, false)
		all, err                  = model.doGetAllBySql(ctx, SelectTypeCount, sqlWithHolder, holderArgs...)
	)
	if err != nil {
		return 0, err
	}
	if len(all) > 0 {
		internalData := core.getInternalColumnFromCtx(ctx)
		if internalData == nil {
			return 0, gerror.NewCode(
				gcode.CodeInternalError,
				`query count error: the internal context data is missing. there's internal issue should be fixed`,
			)
		}
		// If FirstResultColumn present, it returns the value of the first record of the first field.
		// It means it use no cache mechanism, while cache mechanism makes `internalData` missing.
		if v, ok := all[0][internalData.FirstResultColumn]; ok {
			return v.Int(), nil
		}
		// Fields number check.
		var recordFields = model.getRecordFields(all[0])
		if len(recordFields) == 1 {
			for _, v := range all[0] {
				return v.Int(), nil
			}
		}
		// it returns error if there are multiple fields in the result record.
		return 0, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid fields for "Count" operation, result fields number "%d"%s, but expect one`,
			len(recordFields),
			gjson.MustEncodeString(recordFields),
		)
	}
	return 0, nil
}

// Exist does "SELECT 1 FROM ... LIMIT 1" statement for the model.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) Exist(ctx context.Context) (bool, error) {
	model := m.callHandlers(ctx)
	one, err := model.Fields(Raw("1")).One(ctx)
	if err != nil {
		return false, err
	}
	for _, val := range one {
		if val.Bool() {
			return true, nil
		}
	}
	return false, nil
}

// CountColumn does "SELECT COUNT(x) FROM ..." statement for the model.
func (m *Model) CountColumn(ctx context.Context, column string) (int, error) {
	if len(column) == 0 {
		return 0, nil
	}
	model := m.callHandlers(ctx)
	return model.Fields(column).Count(ctx)
}

// Min does "SELECT MIN(x) FROM ..." statement for the model.
func (m *Model) Min(ctx context.Context, column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	model := m.callHandlers(ctx)
	value, err := model.
		Fields(fmt.Sprintf(`MIN(%s)`, model.QuoteWord(column))).
		Value(ctx)
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Max does "SELECT MAX(x) FROM ..." statement for the model.
func (m *Model) Max(ctx context.Context, column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	model := m.callHandlers(ctx)
	value, err := model.
		Fields(fmt.Sprintf(`MAX(%s)`, model.QuoteWord(column))).
		Value(ctx)
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Avg does "SELECT AVG(x) FROM ..." statement for the model.
func (m *Model) Avg(ctx context.Context, column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	model := m.callHandlers(ctx)
	value, err := model.
		Fields(fmt.Sprintf(`AVG(%s)`, model.QuoteWord(column))).
		Value(ctx)
	if err != nil {
		return 0, err
	}
	return value.Float64(), err
}

// Sum does "SELECT SUM(x) FROM ..." statement for the model.
func (m *Model) Sum(ctx context.Context, column string) (float64, error) {
	if len(column) == 0 {
		return 0, nil
	}
	model := m.callHandlers(ctx)
	value, err := model.
		Fields(fmt.Sprintf(`SUM(%s)`, model.QuoteWord(column))).
		Value(ctx)
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
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		switch len(limit) {
		case 1:
			model.limit = limit[0]
		case 2:
			model.start = limit[0]
			model.limit = limit[1]
		}
		return model
	})
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
func (m *Model) Offset(offset int) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.offset = offset
		return model
	})
}

// Distinct forces the query to only return distinct results.
func (m *Model) Distinct() *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.distinct = "DISTINCT "
		return model
	})
}

// Page sets the paging number for the model.
// The parameter `page` is started from 1 for paging.
// Note that, it differs that the Limit function starts from 0 for "LIMIT" statement.
func (m *Model) Page(page, limit int) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		if page <= 0 {
			page = 1
		}
		model.start = (page - 1) * limit
		model.limit = limit
		return model
	})
}

// Having sets the having statement for the model.
// The parameters of this function usage are as the same as function Where.
// See Where.
func (m *Model) Having(having any, args ...any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.having = []any{
			having, args,
		}
		return model
	})
}

// doGetAll does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) doGetAll(ctx context.Context, selectType SelectType, limit1 bool) (Result, error) {
	sqlWithHolder, holderArgs := m.getFormattedSqlAndArgs(ctx, selectType, limit1)
	return m.doGetAllBySql(ctx, selectType, sqlWithHolder, holderArgs...)
}

// doGetAllBySql does the select statement on the database.
func (m *Model) doGetAllBySql(
	ctx context.Context, selectType SelectType, sql string, args ...any,
) (result Result, err error) {
	if result, err = m.getSelectResultFromCache(ctx, sql, args...); err != nil || result != nil {
		return
	}
	in := &HookSelectInput{
		internalParamHookSelect: internalParamHookSelect{
			internalParamHook: internalParamHook{
				link: m.getLink(ctx, false),
			},
			handler: m.hookHandler.Select,
		},
		Model:      m,
		Table:      m.tables,
		Schema:     m.schema,
		Sql:        sql,
		Args:       m.mergeArguments(args),
		SelectType: selectType,
	}
	if result, err = in.Next(ctx); err != nil {
		return
	}

	err = m.saveSelectResultToCache(ctx, selectType, result, sql, args...)
	return
}

func (m *Model) getFormattedSqlAndArgs(
	ctx context.Context, selectType SelectType, limit1 bool,
) (sqlWithHolder string, holderArgs []any) {
	switch selectType {
	case SelectTypeCount:
		queryFields := "COUNT(1)"
		if len(m.fields) > 0 {
			// DO NOT quote the m.fields here, in case of fields like:
			// DISTINCT t.user_id uid
			queryFields = fmt.Sprintf(`COUNT(%s%s)`, m.distinct, m.getFieldsAsStr())
		}
		// Raw SQL Model.
		if m.rawSql != "" {
			sqlWithHolder = fmt.Sprintf("SELECT %s FROM (%s) AS T", queryFields, m.rawSql)
			return sqlWithHolder, nil
		}
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(ctx, false, true)
		sqlWithHolder = fmt.Sprintf("SELECT %s FROM %s%s", queryFields, m.tables, conditionWhere+conditionExtra)
		if len(m.groupBy) > 0 {
			sqlWithHolder = fmt.Sprintf("SELECT COUNT(1) FROM (%s) count_alias", sqlWithHolder)
		}
		return sqlWithHolder, conditionArgs

	default:
		conditionWhere, conditionExtra, conditionArgs := m.formatCondition(ctx, limit1, false)
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
			m.distinct, m.getFieldsFiltered(ctx), m.tables, conditionWhere+conditionExtra,
		)
		return sqlWithHolder, conditionArgs
	}
}

func (m *Model) getHolderAndArgsAsSubModel(ctx context.Context) (holder string, args []any) {
	model := m.callHandlers(ctx)

	holder, args = model.getFormattedSqlAndArgs(
		ctx, SelectTypeDefault, false,
	)
	args = model.mergeArguments(args)
	return
}

func (m *Model) getAutoPrefix() string {
	autoPrefix := ""
	if gstr.Contains(m.tables, " JOIN ") {
		autoPrefix = m.db.GetCore().QuoteWord(
			m.db.GetCore().guessPrimaryTableName(m.tablesInit),
		)
	}
	return autoPrefix
}

func (m *Model) getFieldsAsStr() string {
	var (
		fieldsStr string
		core      = m.db.GetCore()
	)
	for _, v := range m.fields {
		field := gconv.String(v)
		switch {
		case gstr.ContainsAny(field, "()"):
		case gstr.ContainsAny(field, ". "):
		default:
			switch v.(type) {
			case Raw, *Raw:
			default:
				field = core.QuoteString(field)
			}
		}
		if fieldsStr != "" {
			fieldsStr += ","
		}
		fieldsStr += field
	}
	return fieldsStr
}

// getFieldsFiltered checks the fields and fieldsEx attributes, filters and returns the fields that will
// really be committed to underlying database driver.
func (m *Model) getFieldsFiltered(ctx context.Context) string {
	if len(m.fieldsEx) == 0 && len(m.fields) == 0 {
		return defaultField
	}
	if len(m.fieldsEx) == 0 && len(m.fields) > 0 {
		return m.getFieldsAsStr()
	}
	var (
		fieldsArray []string
		fieldsExSet = gset.NewStrSetFrom(gconv.Strings(m.fieldsEx))
	)
	if len(m.fields) > 0 {
		// Filter custom fields with fieldEx.
		fieldsArray = make([]string, 0, 8)
		for _, v := range m.fields {
			field := gconv.String(v)
			fieldsArray = append(fieldsArray, field[gstr.PosR(field, "-")+1:])
		}
	} else {
		if gstr.Contains(m.tables, " ") {
			panic("function FieldsEx supports only single table operations")
		}
		// Filter table fields with fieldEx.
		tableFields, err := m.TableFields(ctx, m.tablesInit)
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

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the `m`.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(
	ctx context.Context, limit1 bool, isCountStatement bool,
) (conditionWhere string, conditionExtra string, conditionArgs []any) {
	var autoPrefix = m.getAutoPrefix()
	// GROUP BY.
	if m.groupBy != "" {
		conditionExtra += " GROUP BY " + m.groupBy
	}
	// WHERE
	conditionWhere, conditionArgs = m.whereBuilder.Build(ctx)
	softDeletingCondition := m.softTimeMaintainer().GetWhereConditionForDelete(ctx)
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
	// HAVING.
	if len(m.having) > 0 {
		havingHolder := WhereHolder{
			Where:  m.having[0],
			Args:   gconv.Interfaces(m.having[1]),
			Prefix: autoPrefix,
		}
		havingStr, havingArgs := formatWhereHolder(ctx, m.db, formatWhereHolderInput{
			WhereHolder: havingHolder,
			OmitNil:     m.option&optionOmitNilWhere > 0,
			OmitEmpty:   m.option&optionOmitEmptyWhere > 0,
			Schema:      m.schema,
			Table:       m.tables,
		})
		if len(havingStr) > 0 {
			conditionExtra += " HAVING " + havingStr
			conditionArgs = append(conditionArgs, havingArgs...)
		}
	}
	// ORDER BY.
	if !isCountStatement { // The count statement of sqlserver cannot contain the order by statement
		if m.orderBy != "" {
			conditionExtra += " ORDER BY " + m.orderBy
		}
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
func (m *Model) doStruct(ctx context.Context, pointer any) error {
	// Auto selecting fields by struct attributes.
	if len(m.fieldsEx) == 0 && len(m.fields) == 0 {
		if v, ok := pointer.(reflect.Value); ok {
			m.Fields(v.Interface())
		} else {
			m.Fields(pointer)
		}
	}
	one, err := m.One(ctx)
	if err != nil {
		return err
	}
	if err = one.Struct(pointer); err != nil {
		return err
	}
	return m.doWithScanStruct(ctx, pointer)
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
func (m *Model) doStructs(ctx context.Context, pointer any) error {
	// Auto selecting fields by struct attributes.
	if len(m.fieldsEx) == 0 && len(m.fields) == 0 {
		if v, ok := pointer.(reflect.Value); ok {
			m.Fields(
				reflect.New(
					v.Type().Elem(),
				).Interface(),
			)
		} else {
			m.Fields(
				reflect.New(
					reflect.ValueOf(pointer).Elem().Type().Elem(),
				).Interface(),
			)
		}
	}
	all, err := m.All(ctx)
	if err != nil {
		return err
	}
	if err = all.Structs(pointer); err != nil {
		return err
	}
	return m.doWithScanStructs(ctx, pointer)
}
