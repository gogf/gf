package gdb

import (
	"context"
	"fmt"
	"reflect"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/internal/reflection"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gutil"
)

// TenantValueType defines the type of tenant value
type TenantValueType string

// Constants for different tenant value types
const (
	ArrayOrSliceType TenantValueType = "ArrayOrSliceType"
	BaseType         TenantValueType = "BaseType"
	NullType         TenantValueType = "NullType"
)

// CtxKeyForTenant defines the context key for tenant ID field and value
type CtxKeyForTenant string

// Context keys for tenant ID field and value
const (
	CtxKeyForTenantIdField CtxKeyForTenant = "CtxKeyForTenantIdField"
	CtxKeyForTenantIdValue CtxKeyForTenant = "CtxKeyForTenantIdValue"
)

// WithTenantIdField sets the tenant ID field name into context
func WithTenantIdField(ctx context.Context, field string) context.Context {
	return context.WithValue(ctx, CtxKeyForTenantIdField, field)
}

// WithTenantIdValue sets the tenant ID value into context
func WithTenantIdValue(ctx context.Context, value any) context.Context {
	return context.WithValue(ctx, CtxKeyForTenantIdValue, value)
}

// WithTenantIdFieldAndValue sets the tenant ID field and value into context
func WithTenantIdFieldAndValue(ctx context.Context, field string, value any) context.Context {
	return WithTenantIdField(WithTenantIdValue(ctx, value), field)
}

// DefaultGetTenantIdFieldValue retrieves tenant ID field and value from context
func DefaultGetTenantIdFieldValue(ctx context.Context) (field string, value any) {
	value = ctx.Value(CtxKeyForTenantIdValue)
	if f := ctx.Value(CtxKeyForTenantIdField); f != nil {
		if a, ok := f.(string); ok {
			return a, value
		}
	}
	return
}

// TenantOption provides configuration options for multi-tenancy support.
type TenantOption struct {
	Enable                    bool                                                // Enable controls whether the tenant feature is enabled.
	PropagateToJoins          bool                                                // PropagateToJoins determines whether tenant conditions should be applied to joined tables.
	InsertNullValue           bool                                                // InsertNullValue  refers to whether to insert a null value for the tenant ID field if the value is nil at insert
	GetTenantIdFieldValueFunc func(ctx context.Context) (field string, value any) // GetTenantIdFieldValueFunc is a function that retrieves the tenant ID field and value for a given context.
}

// Tenant enables multi-tenancy support with the given options
func (m *Model) Tenant(options ...TenantOption) *Model {
	model := m.getModel()
	if len(options) > 0 {
		model.tenantOption = options[0]
		return model
	}
	model.tenantOption.Enable = true
	return model
}

// UnTenant disables multi-tenancy support
func (m *Model) UnTenant() *Model {
	model := m.getModel()
	model.tenantOption = TenantOption{}
	return model
}

// tenantMaintainer creates and returns a TenantMaintainer instance
func (m *Model) tenantMaintainer() *TenantMaintainer {
	return &TenantMaintainer{
		Model: m,
	}
}

// TenantMaintainer handles tenant-related operations
type TenantMaintainer struct {
	*Model
}

// AppendTenantCondition appends tenant condition to the model based on tenant settings
func (tm *TenantMaintainer) AppendTenantCondition(ctx context.Context) {
	if !tm.tenantOption.Enable {
		return
	}
	tenantCondition, tenantConditionArgs, tenantValueType := tm.tenantMaintainer().getWhereConditionForTenant(ctx)
	// Apply different conditions based on tenant value type
	switch tenantValueType {
	case ArrayOrSliceType:
		// Handle array or slice type values with IN condition
		tenantCondition.Iterator(func(k int, v string) bool {
			if value, found := tenantConditionArgs.Get(k); found {
				tm.WhereIn(v, value)
			}
			return true
		})
	case NullType:
		// Handle null values with IS NULL condition
		tenantCondition.Iterator(func(k int, v string) bool {
			tm.WhereNull(v)
			return true
		})
	case BaseType:
		// Handle basic type values with equal condition
		tenantCondition.Iterator(func(k int, v string) bool {
			if value, found := tenantConditionArgs.Get(k); found {
				tm.Wheref(v, value)
			}
			return true
		})
	}
}

// getWhereConditionForTenant generates WHERE conditions for tenant filtering
func (tm *TenantMaintainer) getWhereConditionForTenant(ctx context.Context) (*garray.StrArray, *garray.Array, TenantValueType) {
	var (
		tenantIdField string
		tenantIdValue any
	)
	// Get tenant ID field and value using custom function or default function
	if tm.tenantOption.GetTenantIdFieldValueFunc == nil {
		tenantIdField, tenantIdValue = DefaultGetTenantIdFieldValue(ctx)
	} else {
		tenantIdField, tenantIdValue = tm.tenantOption.GetTenantIdFieldValueFunc(ctx)
	}
	if tenantIdField == "" {
		return nil, nil, ""
	}
	conditionArray := garray.NewStrArray()
	argArray := garray.NewArray()
	tenantValueType := tm.getTenantValueType(tenantIdValue)
	// Handle JOIN queries
	if gstr.Contains(tm.tables, " JOIN ") {
		// Extract main table from JOIN query
		tableMatch, _ := gregex.MatchString(`(.+?) [A-Z]+ JOIN`, tm.tables)
		if c := tm.getCondition(ctx, tableMatch[1], tenantIdField, tenantValueType); c != "" {
			conditionArray.Append(c)
			if tenantValueType != NullType {
				argArray.Append(tenantIdValue)
			}
		}
		// Apply tenant condition to joined tables if PropagateToJoins is enabled
		if tm.tenantOption.PropagateToJoins {
			tableMatches, _ := gregex.MatchAllString(`JOIN ([^()]+?) ON`, tm.tables)
			for _, match := range tableMatches {
				if c := tm.getCondition(ctx, match[1], tenantIdField, tenantValueType); c != "" {
					conditionArray.Append(c)
					if tenantValueType != NullType {
						argArray.Append(tenantIdValue)
					}
				}
			}
		}
	}
	// Handle comma-separated multiple tables
	if conditionArray.Len() == 0 && gstr.Contains(tm.tables, ",") {
		for _, s := range gstr.SplitAndTrim(tm.tables, ",") {
			if c := tm.getCondition(ctx, s, tenantIdField, tenantValueType); c != "" {
				conditionArray.Append(c)
				if tenantValueType != NullType {
					argArray.Append(tenantIdValue)
				}
			}
		}
	}
	if conditionArray.Len() > 0 {
		return conditionArray, argArray, tenantValueType
	}
	// Only one table
	if c := tm.getCondition(ctx, tm.tablesInit, tenantIdField, tenantValueType); c != "" {
		conditionArray.Append(c)
		if tenantValueType != NullType {
			argArray.Append(tenantIdValue)
		}
	}
	return conditionArray, argArray, tenantValueType
}

// getTenantValueType determines the type of tenant value
func (tm *TenantMaintainer) getTenantValueType(value any) TenantValueType {
	if value == nil {
		return NullType
	}
	reflectInfo := reflection.OriginValueAndKind(value)
	switch reflectInfo.OriginKind {
	case reflect.Array, reflect.Slice:
		return ArrayOrSliceType
	default:
		return BaseType
	}
}

// getCondition generates tenant condition field name for a specific table string
func (tm *TenantMaintainer) getCondition(ctx context.Context, tableName string, tenantIdField string, tenantValueType TenantValueType) string {
	var (
		table  string
		schema string
		array1 = gstr.SplitAndTrim(tableName, " ")
		array2 = gstr.SplitAndTrim(array1[0], ".")
	)
	// Parse schema and table name
	if len(array2) >= 2 {
		table = array2[1]
		schema = array2[0]
	} else {
		table = array2[0]
	}
	// Check if tenant field exists in the table
	if !tm.existFieldName(ctx, schema, table, tenantIdField) {
		return ""
	}
	// Generate condition with appropriate field prefix
	if len(array1) >= 3 {
		return tm.getFieldNameCondition(array1[2], tenantIdField, tenantValueType)
	}
	if len(array1) >= 2 {
		return tm.getFieldNameCondition(array1[1], tenantIdField, tenantValueType)
	}
	return tm.getFieldNameCondition(table, tenantIdField, tenantValueType)
}

// getFieldNameCondition generates condition field string based on field prefix and field name
func (tm *TenantMaintainer) getFieldNameCondition(fieldPrefix, fieldName string, tenantValueType TenantValueType) string {
	var (
		quotedFieldPrefix = tm.db.GetCore().QuoteWord(fieldPrefix)
		quotedFieldName   = tm.db.GetCore().QuoteWord(fieldName)
	)
	// Construct full field name with prefix if available
	if quotedFieldPrefix != "" {
		quotedFieldName = fmt.Sprintf(`%s.%s`, quotedFieldPrefix, quotedFieldName)
	}
	// Generate condition based on value type
	switch tenantValueType {
	case BaseType:
		return fmt.Sprintf(`%s = ?`, quotedFieldName)
	default:
		return quotedFieldName
	}
}

// existFieldName checks if a field exists in the specified table
func (tm *TenantMaintainer) existFieldName(ctx context.Context, schema string, table string, tenantIdField string) bool {
	group := tm.db.GetGroup()
	charL, charR := tm.db.GetChars()
	key := genTableFieldsCacheKey(group, gutil.GetOrDefaultStr(tm.db.GetSchema(), schema), gstr.Trim(table, charL+charR))
	v, err := tm.db.GetCore().GetInnerMemCache().Get(ctx, key)
	if err != nil {
		return false
	}
	if !v.IsNil() {
		if fields, ok := v.Val().(map[string]*TableField); ok {
			if _, ok := fields[tenantIdField]; ok {
				return true
			}
		}
	}
	return false
}

// GetFieldNameAndValueForCreate returns the field name and value for creating a record
func (tm *TenantMaintainer) GetFieldNameAndValueForCreate(ctx context.Context, schema string, table string) (string, any) {
	if !tm.tenantOption.Enable {
		return "", nil
	}
	var (
		tenantIdField string
		tenantIdValue any
	)
	// Get tenant ID field and value using custom function or default function
	if tm.tenantOption.GetTenantIdFieldValueFunc == nil {
		tenantIdField, tenantIdValue = DefaultGetTenantIdFieldValue(ctx)
	} else {
		tenantIdField, tenantIdValue = tm.tenantOption.GetTenantIdFieldValueFunc(ctx)
	}
	if tenantIdField == "" {
		return "", nil
	}
	if tm.existFieldName(ctx, schema, table, tenantIdField) {
		charL, charR := tm.db.GetChars()
		return gstr.Trim(tenantIdField, charL+charR), tenantIdValue
	}
	return "", nil
}
