// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
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
func (m *Model) All(where ...any) (Result, error) {
	var ctx = m.GetCtx()
	return m.doGetAll(ctx, SelectTypeDefault, false, where...)
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
func (m *Model) AllAndCount(useFieldForCount bool) (result Result, totalCount int, err error) {
	// Clone the model for counting
	countModel := m.Clone()

	// Decide how to build the COUNT() expression:
	// - If caller explicitly wants to use the single field expression for counting,
	//   honor it (e.g. Fields("DISTINCT col") with useFieldForCount = true).
	// - Otherwise, clear fields to let Count() use its default COUNT(1),
	//   avoiding invalid COUNT(field1, field2, ...) with multiple fields,
	//   or incorrect COUNT(DISTINCT 1) when Distinct() is set.
	if useFieldForCount && len(m.fields) == 1 {
		countModel.fields = m.fields
	} else {
		countModel.fields = nil
	}
	if len(m.pageCacheOption) > 0 {
		countModel = countModel.Cache(m.pageCacheOption[0])
	}

	// Get the total count of records
	totalCount, err = countModel.Count()
	if err != nil {
		return
	}

	// If the total count is 0, there are no records to retrieve, so return early
	if totalCount == 0 {
		return
	}

	resultModel := m.Clone()
	if len(m.pageCacheOption) > 1 {
		resultModel = resultModel.Cache(m.pageCacheOption[1])
	}

	// Retrieve all records
	result, err = resultModel.doGetAll(m.GetCtx(), SelectTypeDefault, false)
	return
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
func (m *Model) One(where ...any) (Record, error) {
	var ctx = m.GetCtx()
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).One()
	}
	all, err := m.doGetAll(ctx, SelectTypeDefault, true)
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
func (m *Model) Array(fieldsAndWhere ...any) (Array, error) {
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Array()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Array()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Array()
		}
	}

	var (
		field string
		core  = m.db.GetCore()
		ctx   = core.injectInternalColumn(m.GetCtx())
	)
	all, err := m.doGetAll(ctx, SelectTypeArray, false)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		internalData := core.getInternalColumnFromCtx(ctx)
		if internalData == nil {
			return nil, gerror.NewCode(
				gcode.CodeInternalError,
				`query count error: the internal context data is missing. there's internal issue should be fixed`,
			)
		}
		// If FirstResultColumn present, it returns the value of the first record of the first field.
		// It means it use no cache mechanism, while cache mechanism makes `internalData` missing.
		field = internalData.FirstResultColumn
		if field == "" {
			// Fields number check.
			var recordFields = m.getRecordFields(all[0])
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
func (m *Model) doStruct(pointer any, where ...any) error {
	model := m
	// Auto selecting fields by struct attributes.
	if len(model.fieldsEx) == 0 && len(model.fields) == 0 {
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
func (m *Model) doStructs(pointer any, where ...any) error {
	model := m
	// Auto selecting fields by struct attributes.
	if len(model.fieldsEx) == 0 && len(model.fields) == 0 {
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

// isExpandingField reports whether the given field entry can expand to multiple
// result columns at the SQL level. This includes:
//   - gdb.Raw / *gdb.Raw expressions (e.g. Raw("name,age")),
//   - the wildcard "*",
//   - table-qualified wildcards like "a.*" or "table_name.*".
//
// SQL aggregate functions such as COUNT(*) are NOT considered expanding because
// they always produce exactly one column.
func isExpandingField(field any) bool {
	var s string
	switch v := field.(type) {
	case Raw:
		s = string(v)
	case *Raw:
		if v == nil {
			return false
		}
		s = string(*v)
	case string:
		s = v
	default:
		return false
	}
	return isExpandingFieldStr(s)
}

// isExpandingFieldStr checks whether a field string can expand to multiple
// result columns. It detects wildcards ("*", "a.*", also with surrounding
// spaces) and top-level commas ("name,age") while ignoring commas nested inside
// parentheses (e.g. function argument lists like "COUNT(DISTINCT `id`,`username`)")
// and commas inside quoted identifiers or string literals.
func isExpandingFieldStr(s string) bool {
	s = gstr.Trim(s)
	if s == "*" || gstr.HasSuffix(s, ".*") {
		return true
	}
	var (
		depth      int
		quote      rune
		quoteClose rune
	)
	for _, ch := range s {
		// Inside a quoted identifier or string literal, only the closing quote matters.
		if quote != 0 {
			if ch == quoteClose {
				quote = 0
			}
			continue
		}
		switch ch {
		case '`', '"', '\'':
			quote = ch
			quoteClose = ch
		case '[':
			quote = ch
			quoteClose = ']'
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				return true
			}
		}
	}
	return false
}

// isSingleFieldSpecified reports whether the model has exactly one explicit
// field that resolves to a single column. FieldsEx is not treated as a
// single-field specification.
//
// It rejects:
//   - gdb.Raw entries that may expand to multiple columns,
//   - wildcard "*" and table-qualified wildcards like "a.*",
//   - a field that FieldsEx filters out, which would leave an empty projection.
func (m *Model) isSingleFieldSpecified() bool {
	if len(m.fields) != 1 {
		return false
	}
	if isExpandingField(m.fields[0]) {
		return false
	}
	return !m.isFieldExcludedByFieldsEx(m.fields[0])
}

// isFieldExcludedByFieldsEx reports whether `field` is filtered out by FieldsEx.
// It matches the same way getFieldsFiltered does: the field string is compared
// as written (after stripping the "-" prefix used by Fields) against the
// FieldsEx entries.
func (m *Model) isFieldExcludedByFieldsEx(field any) bool {
	if len(m.fieldsEx) == 0 {
		return false
	}
	fieldStr := gconv.String(field)
	fieldStr = fieldStr[gstr.PosR(fieldStr, "-")+1:]
	return gset.NewStrSetFrom(gconv.Strings(m.fieldsEx)).Contains(fieldStr)
}

// validateFieldsExSingleColumn checks whether FieldsEx leaves exactly one table
// column, which is required for scanning into a basic type or a slice of basic
// types. Value()/Array() cannot be used as a post-query column-count check
// here, because they return FirstResultColumn before that validation runs.
func (m *Model) validateFieldsExSingleColumn() error {
	if m.db == nil {
		return nil
	}
	// In join statements, FieldsEx filters the fields of all joined tables, so the
	// remaining column count must be computed the same way getFieldsFiltered does,
	// instead of only counting the fields of the main table. Sub-query references
	// contribute one "prefix.*" entry here; such expanding columns are still caught
	// after the query by the column-count check in scanBasicType.
	if len(m.tableRefs) > 0 {
		fieldsArray := m.getJoinedTableFields(gset.NewStrSetFrom(gconv.Strings(m.fieldsEx)))
		if len(fieldsArray) != 1 {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`Scan into basic type requires exactly 1 field specified via Fields(), but FieldsEx leaves %d columns after filtering`,
				len(fieldsArray),
			)
		}
		return nil
	}
	tableFields, err := m.TableFields(m.tablesInit)
	if err != nil {
		return err
	}
	if len(tableFields) == 0 {
		return gerror.Newf(`empty table fields for table "%s"`, m.tablesInit)
	}
	var (
		remaining     = len(tableFields)
		excludedNames = make(map[string]bool, len(tableFields))
	)
	for _, excluded := range m.fieldsEx {
		for _, name := range gstr.SplitAndTrim(gconv.String(excluded), ",") {
			if name == "" {
				continue
			}
			if _, ok := tableFields[name]; ok && !excludedNames[name] {
				excludedNames[name] = true
				remaining--
			}
		}
	}
	if remaining != 1 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`Scan into basic type requires exactly 1 field specified via Fields(), but FieldsEx leaves %d columns after filtering`,
			remaining,
		)
	}
	return nil
}

// validateSingleFieldSpecified checks that the model has a field specification
// suitable for a single-column scan into a basic type or a slice of basic types.
// It accepts both an explicit single Fields entry and FieldsEx, because FieldsEx
// can narrow the result to a single column even though m.fields is empty.
func (m *Model) validateSingleFieldSpecified() error {
	if m.isSingleFieldSpecified() {
		return nil
	}
	if len(m.fields) == 1 {
		// The single field is filtered out by FieldsEx, which leaves an empty
		// projection and produces invalid SQL like "SELECT  FROM ...".
		if m.isFieldExcludedByFieldsEx(m.fields[0]) {
			return gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`Scan into basic type requires exactly 1 field, but field %v is excluded by FieldsEx`,
				m.fields[0],
			)
		}
		// Reaching here with exactly one field means it is an expanding field
		// (gdb.Raw or wildcard like "*" / "a.*").
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`Scan into basic type does not support expanding field %v (gdb.Raw or wildcard); use an explicit field name instead`,
			m.fields[0],
		)
	}
	if len(m.fields) == 0 && len(m.fieldsEx) > 0 {
		return m.validateFieldsExSingleColumn()
	}
	return gerror.NewCodef(
		gcode.CodeInvalidParameter,
		`Scan into basic type requires exactly 1 field specified via Fields(), but got %d`,
		len(m.fields),
	)
}

// Scan automatically calls Struct or Structs function according to the type of parameter `pointer`.
// It calls function doStruct if `pointer` is type of *struct/**struct.
// It calls function doStructs if `pointer` is type of *[]struct/*[]*struct.
//
// It also supports scanning a single field result into a basic type (int, float, string, bool, etc.),
// a pointer to a basic type, a slice of basic types, or a slice of pointers to basic types. In these
// cases exactly one field must be specified via Fields (or a field-producing helper such as FieldSum),
// otherwise an error is returned. FieldsEx is accepted only when it leaves exactly one column.
// `[]byte` and `[]uint8` targets are treated as a single column value that keeps its raw bytes,
// not as a slice of scalar values.
//
// The optional parameter `where` is the same as the parameter of Model.Where function,  see Model.Where.
//
// Note that it returns sql.ErrNoRows if the given parameter `pointer` pointed to a variable that has
// default value and there's no record retrieved with the given conditions from table.
//
// # No row versus NULL column value
//
// Scanning into a variable of basic type distinguishes the two cases:
//   - no record matches: sql.ErrNoRows is returned when the target holds a default value,
//     consistent with Record.Struct and Result.Structs;
//   - a record matches but the column value is NULL: the target is left untouched and
//     no error is returned, which follows gconv.Scan semantics for a nil source value.
//     Note that a pointer target ends up as a nil pointer in this case, and a slice of
//     pointers keeps a nil element for that row.
//
// Scanning into a slice of basic types also follows Result.Structs: no record leaves the
// target untouched and returns sql.ErrNoRows only when the target already holds values.
// The result must still have exactly one column, which is checked after the query as
// well, so a projection that the pre-query validation cannot recognize as multi-column
// is rejected instead of silently keeping only the first column.
//
// Example:
//
//	user := new(User)
//	err  := db.Model("user").Where("id", 1).Scan(user)
//
//	user := (*User)(nil)
//	err  := db.Model("user").Where("id", 1).Scan(&user)
//
//	users := ([]User)(nil)
//	err   := db.Model("user").Scan(&users)
//
//	users := ([]*User)(nil)
//	err   := db.Model("user").Scan(&users)
//
//	var name string
//	err   := db.Model("user").Fields("name").Where("id", 1).Scan(&name)
//
//	var names []string
//	err   := db.Model("user").Fields("name").Scan(&names)
//
//	var names2 []*string
//	err   := db.Model("user").Fields("name").Scan(&names2)
//
//	var avatar []byte
//	err   := db.Model("user").Fields("avatar").Where("id", 1).Scan(&avatar)
func (m *Model) Scan(pointer any, where ...any) error {
	reflectInfo := reflection.OriginTypeAndKind(pointer)
	if reflectInfo.InputKind != reflect.Pointer {
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`the parameter "pointer" for function Scan should type of pointer`,
		)
	}
	switch reflectInfo.OriginKind {
	case reflect.Slice, reflect.Array:
		// OriginType has all outer pointers stripped, so Elem() gives the slice element type.
		elemType := reflectInfo.OriginType.Elem()
		for elemType.Kind() == reflect.Pointer {
			elemType = elemType.Elem()
		}

		// `[]byte` and `[]uint8` hold one column raw value instead of a slice of scalar
		// values. Their element kind is `uint8`, so the multi-row scalar branch would
		// decode each row into one byte and silently truncate every value.
		if reflectInfo.OriginKind == reflect.Slice && reflectInfo.OriginType.Elem().Kind() == reflect.Uint8 {
			return m.scanBasicType(pointer, where...)
		}

		// Only route to the basic-type slice branch when the element is a true scalar type
		// (int/uint/float/bool/string). Map/record slices and other non-struct kinds must
		// fall through to doStructs; otherwise []Record or similar results would be
		// misclassified as a slice of basic values.
		if elemType != nil && reflection.IsBasicKind(elemType.Kind()) {
			return m.scanSliceOfBasicType(pointer, where...)
		}
		return m.doStructs(pointer, where...)

	case reflect.Struct, reflect.Invalid:
		return m.doStruct(pointer, where...)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Bool, reflect.String:
		return m.scanBasicType(pointer, where...)
	default:
		return gerror.NewCode(
			gcode.CodeInvalidParameter,
			`element of parameter "pointer" for function Scan should type of struct/*struct/[]struct/[]*struct, or basic type like int/string/float64/bool`,
		)
	}
}

// scanBasicType scans a single column value into a variable of basic type, or into a pointer
// to a basic type. It keeps the no-row and NULL value semantics aligned with Record.Struct:
// no row returns sql.ErrNoRows for a non-nil target, while a NULL column value leaves the
// target untouched without error.
func (m *Model) scanBasicType(pointer any, where ...any) error {
	if err := m.validateSingleFieldSpecified(); err != nil {
		return err
	}
	// Use a clone to avoid mutating the original model.
	model := m.Clone()
	if len(where) > 0 {
		model = model.Where(where[0], where[1:]...)
	}
	// One() is used instead of Value() because Value() returns nil both for no matching
	// record and for a NULL column value, which makes the two cases indistinguishable.
	record, err := model.One()
	if err != nil {
		return err
	}
	if record == nil {
		// No record matches.
		if !empty.IsNil(pointer, true) {
			return sql.ErrNoRows
		}
		return nil
	}
	if len(record) != 1 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid fields for "Scan" operation, result fields number "%d"%s, but expect one`,
			len(record),
			gjson.MustEncodeString(m.getRecordFields(record)),
		)
	}
	for _, value := range record {
		// A record matches but the column value is NULL, which leaves the target untouched.
		if value == nil {
			return nil
		}
		return value.Scan(pointer)
	}
	return nil
}

// scanSliceOfBasicType scans a single column result into a slice of basic types, or a slice
// of pointers to basic types. The no-row behavior is aligned with Result.Structs: it returns
// sql.ErrNoRows and keeps the target untouched when the target already holds values, and
// keeps a nil slice as nil otherwise.
func (m *Model) scanSliceOfBasicType(pointer any, where ...any) error {
	if err := m.validateSingleFieldSpecified(); err != nil {
		return err
	}
	// Use a clone to avoid mutating the original model. Do not pass the field
	// as an argument to Array() because it already lives in m.fields and
	// Array() would append it again (Fields is append-only).
	model := m.Clone()
	if len(where) > 0 {
		model = model.Where(where[0], where[1:]...)
	}
	// doGetAll is used instead of Array() because Array() returns the first result
	// column without checking the column count, which would silently drop the extra
	// columns of a projection that the pre-query validation failed to recognize.
	all, err := model.doGetAll(model.GetCtx(), SelectTypeArray, false)
	if err != nil {
		return err
	}
	if len(all) == 0 {
		// Keep the no-row semantics identical to Result.Structs.
		if !empty.IsEmpty(pointer, true) {
			return sql.ErrNoRows
		}
		return nil
	}
	recordFields := m.getRecordFields(all[0])
	if len(recordFields) != 1 {
		return gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid fields for "Scan" operation, result fields number "%d"%s, but expect one`,
			len(recordFields),
			gjson.MustEncodeString(recordFields),
		)
	}
	return all.Array(recordFields[0]).Scan(pointer)
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
func (m *Model) ScanAndCount(pointer any, totalCount *int, useFieldForCount bool) (err error) {
	// support Fields with *, example: .Fields("a.*, b.name"). Count sql is select count(1) from xxx
	countModel := m.Clone()
	// Decide how to build the COUNT() expression:
	// - If caller explicitly wants to use the single field expression for counting,
	//   honor it (e.g. Fields("DISTINCT col") with useFieldForCount = true).
	// - Otherwise, clear fields to let Count() use its default COUNT(1),
	//   avoiding invalid COUNT(field1, field2, ...) with multiple fields,
	//   or incorrect COUNT(DISTINCT 1) when Distinct() is set.
	if useFieldForCount && len(m.fields) == 1 {
		countModel.fields = m.fields
	} else {
		countModel.fields = nil
	}
	if len(m.pageCacheOption) > 0 {
		countModel = countModel.Cache(m.pageCacheOption[0])
	}
	// Get the total count of records
	*totalCount, err = countModel.Count()
	if err != nil {
		return err
	}

	// If the total count is 0, there are no records to retrieve, so return early
	if *totalCount == 0 {
		return
	}
	scanModel := m.Clone()
	if len(m.pageCacheOption) > 1 {
		scanModel = scanModel.Cache(m.pageCacheOption[1])
	}
	err = scanModel.Scan(pointer)
	return
}

// ScanList converts `r` to struct slice which contains other complex struct attributes.
// Note that the parameter `listPointer` should be type of *[]struct/*[]*struct.
//
// See Result.ScanList.
func (m *Model) ScanList(structSlicePointer any, bindToAttrName string, relationAttrNameAndFields ...string) (err error) {
	var result Result
	out, err := checkGetSliceElementInfoForScanList(structSlicePointer, bindToAttrName)
	if err != nil {
		return err
	}
	if len(m.fields) > 0 || len(m.fieldsEx) != 0 {
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
	return m.doScanListWithResult(result, structSlicePointer, bindToAttrName, relationAttrName, relationFields)
}

func (m *Model) doScanListWithResult(
	result Result, structSlicePointer any, bindToAttrName, relationAttrName, relationFields string,
) error {
	out, err := checkGetSliceElementInfoForScanList(structSlicePointer, bindToAttrName)
	if err != nil {
		return err
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

// Value retrieves a specified record value from table and returns the result as interface type.
// It returns nil if there's no record found with the given conditions from table.
//
// If the optional parameter `fieldsAndWhere` is given, the fieldsAndWhere[0] is the selected fields
// and fieldsAndWhere[1:] is treated as where condition fields.
// Also see Model.Fields and Model.Where functions.
func (m *Model) Value(fieldsAndWhere ...any) (Value, error) {
	var (
		core = m.db.GetCore()
		ctx  = core.injectInternalColumn(m.GetCtx())
	)
	if len(fieldsAndWhere) > 0 {
		if len(fieldsAndWhere) > 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1], fieldsAndWhere[2:]...).Value()
		} else if len(fieldsAndWhere) == 2 {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Where(fieldsAndWhere[1]).Value()
		} else {
			return m.Fields(gconv.String(fieldsAndWhere[0])).Value()
		}
	}
	var (
		sqlWithHolder, holderArgs = m.getFormattedSqlAndArgs(ctx, SelectTypeValue, true)
		all, err                  = m.doGetAllBySql(ctx, SelectTypeValue, sqlWithHolder, holderArgs...)
	)
	if err != nil {
		return nil, err
	}
	if len(all) > 0 {
		internalData := core.getInternalColumnFromCtx(ctx)
		if internalData == nil {
			return nil, gerror.NewCode(
				gcode.CodeInternalError,
				`query count error: the internal context data is missing. there's internal issue should be fixed`,
			)
		}
		// If FirstResultColumn present, it returns the value of the first record of the first field.
		// It means it use no cache mechanism, while cache mechanism makes `internalData` missing.
		if v, ok := all[0][internalData.FirstResultColumn]; ok {
			return v, nil
		}
		// Fields number check.
		var recordFields = m.getRecordFields(all[0])
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
func (m *Model) Count(where ...any) (int, error) {
	var (
		core = m.db.GetCore()
		ctx  = core.injectInternalColumn(m.GetCtx())
	)
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Count()
	}
	var (
		sqlWithHolder, holderArgs = m.getFormattedSqlAndArgs(ctx, SelectTypeCount, false)
		all, err                  = m.doGetAllBySql(ctx, SelectTypeCount, sqlWithHolder, holderArgs...)
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
		var recordFields = m.getRecordFields(all[0])
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
func (m *Model) Exist(where ...any) (bool, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).Exist()
	}
	one, err := m.Fields(Raw("1")).One()
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
// Note: Negative values are treated as zero.
func (m *Model) Limit(limit ...int) *Model {
	model := m.getModel()
	switch len(limit) {
	case 1:
		if limit[0] < 0 {
			limit[0] = 0
		}
		model.limit = limit[0]
	case 2:
		if limit[0] < 0 {
			limit[0] = 0
		}
		if limit[1] < 0 {
			limit[1] = 0
		}
		model.start = limit[0]
		model.limit = limit[1]
	}
	return model
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
// Note: Negative values are treated as zero.
func (m *Model) Offset(offset int) *Model {
	model := m.getModel()
	if offset < 0 {
		offset = 0
	}
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
// Note: Negative limit values are treated as zero.
func (m *Model) Page(page, limit int) *Model {
	model := m.getModel()
	if page <= 0 {
		page = 1
	}
	if limit < 0 {
		limit = 0
	}
	model.start = (page - 1) * limit
	model.limit = limit
	return model
}

// Having sets the having statement for the model.
// The parameters of this function usage are as the same as function Where.
// See Where.
func (m *Model) Having(having any, args ...any) *Model {
	model := m.getModel()
	model.having = []any{
		having, args,
	}
	return model
}

// doGetAll does "SELECT FROM ..." statement for the model.
// It retrieves the records from table and returns the result as slice type.
// It returns nil if there's no record retrieved with the given conditions from table.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
// The optional parameter `where` is the same as the parameter of Model.Where function,
// see Model.Where.
func (m *Model) doGetAll(ctx context.Context, selectType SelectType, limit1 bool, where ...any) (Result, error) {
	if len(where) > 0 {
		return m.Where(where[0], where[1:]...).All()
	}
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
				link: m.getLink(false),
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
			conditionWhere, conditionExtra, conditionArgs := m.formatCondition(ctx, false, true)
			sqlWithHolder = fmt.Sprintf(
				"SELECT %s FROM (%s%s) AS T",
				queryFields, m.rawSql, conditionWhere+conditionExtra,
			)
			return sqlWithHolder, conditionArgs
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
			m.distinct, m.getFieldsFiltered(), m.tables, conditionWhere+conditionExtra,
		)
		return sqlWithHolder, conditionArgs
	}
}

func (m *Model) getHolderAndArgsAsSubModel(ctx context.Context) (holder string, args []any) {
	holder, args = m.getFormattedSqlAndArgs(
		ctx, SelectTypeDefault, false,
	)
	args = m.mergeArguments(args)
	return
}

func (m *Model) getAutoPrefix() string {
	autoPrefix := ""
	if gstr.Contains(m.tables, " JOIN ") {
		autoPrefix = m.QuoteWord(
			m.db.GetCore().guessPrimaryTableName(m.tablesInit),
		)
	}
	return autoPrefix
}

func (m *Model) getFieldsAsStr() string {
	var (
		fieldsStr string
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
				field = m.QuoteWord(field)
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
func (m *Model) getFieldsFiltered() string {
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
	} else if len(m.tableRefs) > 0 {
		// Filter fields of all joined tables with fieldEx.
		fieldsArray = m.getJoinedTableFields(fieldsExSet)
	} else {
		// A single table with an alias, like "user u" or "user AS u", is still a single
		// table operation: its fields are looked up by the table name without the alias.
		// Only a table expression that is not a single identifier after removing the
		// alias, like "user,user_detail" or "user JOIN user_detail", is rejected.
		if tableName := m.db.GetCore().tableNameForFields(m.tables); !isSqlIdentifier(tableName, true) {
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
		newFields += m.QuoteWord(k)
	}
	return newFields
}

// getJoinedTableFields returns the prefixed fields of all tables in join statement in sequence.
// Fields without prefix in `fieldsExSet` are excluded from the main table, just like FieldsEx.
// It uses "prefix.*" for sub-query whose fields cannot be retrieved.
func (m *Model) getJoinedTableFields(fieldsExSet *gset.StrSet) []string {
	fieldsArray := make([]string, 0, 8)
	for i, ref := range m.tableRefs {
		if ref.prefix == "" {
			panic("function FieldsEx supports only single table operations")
		}
		if ref.table == "" {
			fieldsArray = append(fieldsArray, ref.prefix+".*")
			continue
		}
		tableFields, err := m.TableFields(ref.table)
		if err != nil {
			panic(err)
		}
		if len(tableFields) == 0 {
			panic(fmt.Sprintf(`empty table fields for table %q`, ref.table))
		}
		names := make([]string, len(tableFields))
		for k, v := range tableFields {
			names[v.Index] = k
		}
		excluded := m.getJoinedTableFieldsEx(ref, i == 0, fieldsExSet)
		for _, name := range names {
			if excluded.Contains(name) {
				continue
			}
			fieldsArray = append(fieldsArray, ref.prefix+"."+m.QuoteWord(name))
		}
	}
	return fieldsArray
}

// getJoinedTableFieldsEx returns the excluded field names of table `ref` from `fieldsExSet`,
// folded like the table fields. The prefix and the field of each excluded item are matched
// separately, so that "user.id", "`gf_user`.`id`" and "db.items.secret" are all recognized.
// Fields without prefix are excluded only if `isMain` is true.
func (m *Model) getJoinedTableFieldsEx(ref tableRef, isMain bool, fieldsExSet *gset.StrSet) *gset.StrSet {
	var (
		excluded  = gset.NewStrSet()
		refPrefix = m.unquoteIdentifier(ref.prefix)
	)
	fieldsExSet.Iterator(func(v string) bool {
		var (
			pos   = gstr.PosR(v, ".")
			field = m.unquoteIdentifier(v[pos+1:])
		)
		if pos < 0 {
			if !isMain {
				return true
			}
		} else {
			prefix := v[:pos]
			// Table name without alias may be given without table prefix.
			if m.unquoteIdentifier(prefix) != refPrefix &&
				(ref.alias || m.unquoteIdentifier(m.db.GetCore().QuotePrefixTableName(prefix)) != refPrefix) {
				return true
			}
		}
		excluded.Add(field, m.db.FoldIdentifier(field))
		return true
	})
	return excluded
}

// unquoteIdentifier removes the security chars from each part of the dot-separated identifier `s`.
func (m *Model) unquoteIdentifier(s string) string {
	charL, charR := m.db.GetChars()
	array := gstr.SplitAndTrim(s, ".")
	for i, v := range array {
		array[i] = gstr.Trim(v, charL+charR)
	}
	return gstr.Join(array, ".")
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
	conditionWhere, conditionArgs = m.whereBuilder.Build()
	softDeletingCondition := m.softTimeMaintainer().GetDeleteCondition(ctx)
	if m.rawSql != "" && conditionWhere != "" {
		if gstr.ContainsI(m.rawSql, sqlWhere) {
			conditionWhere = sqlAnd + conditionWhere
		} else {
			conditionWhere = sqlWhere + conditionWhere
		}
	} else if !m.unscoped && softDeletingCondition != "" {
		if conditionWhere == "" {
			conditionWhere = fmt.Sprintf(` WHERE %s`, softDeletingCondition)
		} else {
			conditionWhere = fmt.Sprintf(` WHERE (%s) AND %s`, conditionWhere, softDeletingCondition)
		}
	} else {
		if conditionWhere != "" {
			conditionWhere = sqlWhere + conditionWhere
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
