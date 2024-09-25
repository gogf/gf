package gdb

import (
	"context"
	"database/sql"
)

var (
	// Compile detection implementation
	_ ModelInterface = (*Model)(nil)
)

type ModelInterface interface {
	ModelSelectInterface
	ModelExecInterface
	ModelOmitInterface
	ModelUtilInterface
	ModelWhereInterface
	// 主要用于Model.Clone
	// 当Safe为true时，每次链式调用需要深拷贝Model
	// 同时也需要更新DefaultHookModelInterfaceImpl里面的Model
	setModel(model *Model)
}

type ModelOmitInterface interface {
	// gdb_model_option.go
	OmitEmpty() *Model
	OmitEmptyWhere() *Model
	OmitEmptyData() *Model
	OmitNil() *Model
	OmitNilWhere() *Model
	OmitNilData() *Model
}

type ModelExecInterface interface {
	// gdb_model_insert.go
	Batch(batch int) *Model
	Data(data ...interface{}) *Model
	OnConflict(onConflict ...interface{}) *Model
	OnDuplicate(onDuplicate ...interface{}) *Model
	OnDuplicateEx(onDuplicateEx ...interface{}) *Model
	Insert(data ...interface{}) (result sql.Result, err error)
	InsertAndGetId(data ...interface{}) (lastInsertId int64, err error)
	InsertIgnore(data ...interface{}) (result sql.Result, err error)
	Replace(data ...interface{}) (result sql.Result, err error)
	Save(data ...interface{}) (result sql.Result, err error)

	// gdb_model_delete.go
	Delete(where ...interface{}) (result sql.Result, err error)

	// gdb_model_update.go
	Update(dataAndWhere ...interface{}) (result sql.Result, err error)
	UpdateAndGetAffected(dataAndWhere ...interface{}) (affected int64, err error)
	Increment(column string, amount interface{}) (sql.Result, error)
	Decrement(column string, amount interface{}) (sql.Result, error)
}

type ModelSelectResult interface {
	All(where ...interface{}) (Result, error)
	AllAndCount(useFieldForCount bool) (result Result, totalCount int, err error)
	Chunk(size int, handler ChunkHandler)
	One(where ...interface{}) (Record, error)
	Array(fieldsAndWhere ...interface{}) ([]Value, error)
	Scan(pointer interface{}, where ...interface{}) error
	ScanAndCount(pointer interface{}, totalCount *int, useFieldForCount bool) (err error)
	ScanList(structSlicePointer interface{}, bindToAttrName string, relationAttrNameAndFields ...string) (err error)
	Value(fieldsAndWhere ...interface{}) (Value, error)

	Count(where ...interface{}) (int, error)
	CountColumn(column string) (int, error)
	Min(column string) (float64, error)
	Max(column string) (float64, error)
	Avg(column string) (float64, error)
	Sum(column string) (float64, error)
}

type ModelSelectJoin interface {
	LeftJoin(tableOrSubQueryAndJoinConditions ...string) *Model
	RightJoin(tableOrSubQueryAndJoinConditions ...string) *Model
	InnerJoin(tableOrSubQueryAndJoinConditions ...string) *Model
	LeftJoinOnField(table, field string) *Model
	RightJoinOnField(table, field string) *Model
	InnerJoinOnField(table, field string) *Model
	LeftJoinOnFields(table, firstField, operator, secondField string) *Model
	RightJoinOnFields(table, firstField, operator, secondField string) *Model
	InnerJoinOnFields(table, firstField, operator, secondField string) *Model
}

type ModelSelectField interface {
	Fields(fieldNamesOrMapStruct ...interface{}) *Model
	FieldsPrefix(prefixOrAlias string, fieldNamesOrMapStruct ...interface{}) *Model
	FieldsEx(fieldNamesOrMapStruct ...interface{}) *Model
	FieldsExPrefix(prefixOrAlias string, fieldNamesOrMapStruct ...interface{}) *Model
	FieldCount(column string, as ...string) *Model
	FieldSum(column string, as ...string) *Model
	FieldMin(column string, as ...string) *Model
	FieldMax(column string, as ...string) *Model
	FieldAvg(column string, as ...string) *Model
	GetFieldsStr(prefix ...string) string
	GetFieldsExStr(fields string, prefix ...string) (string, error)
	HasField(field string) (bool, error)
}

type ModelSelectOrderGroup interface {
	Order(orderBy ...interface{}) *Model
	OrderAsc(column string) *Model
	OrderDesc(column string) *Model
	OrderRandom() *Model
	Group(groupBy ...string) *Model
}

type ModelSelectPage interface {
	Limit(limit ...int) *Model
	Offset(offset int) *Model
	Page(page, limit int) *Model
}

type ModelSelectInterface interface {
	ModelSelectResult
	ModelSelectJoin
	ModelSelectField
	ModelSelectOrderGroup
	ModelSelectPage

	Union(unions ...*Model) *Model
	UnionAll(unions ...*Model) *Model

	With(objects ...interface{}) *Model
	WithAll() *Model
	Having(having interface{}, args ...interface{}) *Model
	Distinct() *Model
}

type ModelUtilInterface interface {
	// gdb_model.go
	Raw(rawSql string, args ...interface{}) *Model
	Partition(partitions ...string) *Model
	Ctx(ctx context.Context) *Model
	GetCtx() context.Context
	As(as string) *Model
	DB(db DB) *Model
	TX(tx TX) *Model
	Schema(schema string) *Model
	Clone() *Model
	Master() *Model
	Slave() *Model
	Safe(safe ...bool) *Model
	Args(args ...interface{}) *Model
	Handler(handlers ...ModelHandler) *Model
	// gdb_model_hook.go
	Hook(hook HookHandler) *Model

	// gdb_model_soft_time.go
	SoftTime(option SoftTimeOption) *Model
	Unscoped() *Model

	// gdb_model_lock.go
	LockUpdate() *Model
	LockShared() *Model
	// gdb_model_transaction.go
	Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error)
	// gdb_model_cache.go
	Cache(option CacheOption) *Model
	// gdb_model_utility.go
	QuoteWord(s string) string
	TableFields(tableStr string, schema ...string) (fields map[string]*TableField, err error)
}

type ModelWhereInterface interface {
	// gdb_model_where.go
	Where(where interface{}, args ...interface{}) *Model
	Wheref(format string, args ...interface{}) *Model
	WherePri(where interface{}, args ...interface{}) *Model
	WhereLT(column string, value interface{}) *Model
	WhereLTE(column string, value interface{}) *Model
	WhereGT(column string, value interface{}) *Model
	WhereGTE(column string, value interface{}) *Model
	WhereBetween(column string, min, max interface{}) *Model
	WhereLike(column string, like string) *Model
	WhereIn(column string, in interface{}) *Model
	WhereNull(columns ...string) *Model
	WhereNotBetween(column string, min, max interface{}) *Model
	WhereNotLike(column string, like interface{}) *Model
	WhereNot(column string, value interface{}) *Model
	WhereNotIn(column string, in interface{}) *Model
	WhereNotNull(columns ...string) *Model

	// gdb_model_whereor.go
	WhereOr(where interface{}, args ...interface{}) *Model
	WhereOrf(format string, args ...interface{}) *Model
	WhereOrLT(column string, value interface{}) *Model
	WhereOrLTE(column string, value interface{}) *Model
	WhereOrGT(column string, value interface{}) *Model
	WhereOrGTE(column string, value interface{}) *Model
	WhereOrBetween(column string, min, max interface{}) *Model
	WhereOrLike(column string, like interface{}) *Model
	WhereOrIn(column string, in interface{}) *Model
	WhereOrNull(columns ...string) *Model
	WhereOrNotBetween(column string, min, max interface{}) *Model
	WhereOrNotLike(column string, like interface{}) *Model
	WhereOrNot(column string, value interface{}) *Model
	WhereOrNotIn(column string, in interface{}) *Model
	WhereOrNotNull(columns ...string) *Model

	// gdb_model_where_prefix.go
	WherePrefix(prefix string, where interface{}, args ...interface{}) *Model
	WherePrefixLT(prefix string, column string, value interface{}) *Model
	WherePrefixLTE(prefix string, column string, value interface{}) *Model
	WherePrefixGT(prefix string, column string, value interface{}) *Model
	WherePrefixGTE(prefix string, column string, value interface{}) *Model
	WherePrefixBetween(prefix string, column string, min, max interface{}) *Model
	WherePrefixLike(prefix string, column string, like interface{}) *Model
	WherePrefixIn(prefix string, column string, in interface{}) *Model
	WherePrefixNull(prefix string, columns ...string) *Model
	WherePrefixNotBetween(prefix string, column string, min, max interface{}) *Model
	WherePrefixNotLike(prefix string, column string, like interface{}) *Model
	WherePrefixNot(prefix string, column string, value interface{}) *Model
	WherePrefixNotIn(prefix string, column string, in interface{}) *Model
	WherePrefixNotNull(prefix string, columns ...string) *Model

	// gdb_model_whereor_prefix.go
	WhereOrPrefix(prefix string, where interface{}, args ...interface{}) *Model
	WhereOrPrefixLT(prefix string, column string, value interface{}) *Model
	WhereOrPrefixLTE(prefix string, column string, value interface{}) *Model
	WhereOrPrefixGT(prefix string, column string, value interface{}) *Model
	WhereOrPrefixGTE(prefix string, column string, value interface{}) *Model
	WhereOrPrefixBetween(prefix string, column string, min, max interface{}) *Model
	WhereOrPrefixLike(prefix string, column string, like interface{}) *Model
	WhereOrPrefixIn(prefix string, column string, in interface{}) *Model
	WhereOrPrefixNull(prefix string, columns ...string) *Model
	WhereOrPrefixNotBetween(prefix string, column string, min, max interface{}) *Model
	WhereOrPrefixNotLike(prefix string, column string, like interface{}) *Model
	WhereOrPrefixNotIn(prefix string, column string, in interface{}) *Model
	WhereOrPrefixNotNull(prefix string, columns ...string) *Model
	WhereOrPrefixNot(prefix string, column string, value interface{}) *Model
}
