// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"context"
	"fmt"
	"reflect"
	"runtime"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/intlog"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gtag"
	"github.com/gogf/gf/v2/util/gutil"
)

const parseRuleForeach = "foreach"

// ParseFunc is the custom function for request parameter pre-processing.
//
// The function is called before request parameter binding and validation.
// It can transform the input value and return the new value for later binding.
type ParseFunc func(ctx context.Context, in ParseFuncInput) (any, error)

// ParseFuncInput holds the input parameters passed to custom parse function ParseFunc.
type ParseFuncInput struct {
	// Rule specifies the parse rule string, like "trim-space", "trim-prefix:demo", etc.
	Rule string

	// Name specifies the rule name of Rule, like "trim-space" and "trim-prefix".
	Name string

	// Pattern specifies the rule parameter pattern of Rule, like "demo" from "trim-prefix:demo".
	Pattern string

	// Field specifies the field path for this rule to process.
	Field string

	// FieldType specifies the type of the field.
	FieldType reflect.Type

	// Value specifies the current field value for this rule to process.
	Value any

	// Data specifies the current request data map that contains this field value.
	Data map[string]any

	// Request specifies the current HTTP request.
	Request *Request
}

type parseRuleItem struct {
	Rule    string
	Name    string
	Pattern string
}

var (
	// customParseFuncMap stores the custom parse functions.
	// map[Rule]ParseFunc
	customParseFuncMap = make(map[string]ParseFunc)
)

func init() {
	RegisterParseRuleByMap(map[string]ParseFunc{
		"trim-space":   parseRuleTrimSpace,
		"trim-left":    parseRuleTrimLeft,
		"trim-right":   parseRuleTrimRight,
		"trim-prefix":  parseRuleTrimPrefix,
		"trim-suffix":  parseRuleTrimSuffix,
		"trim":         parseRuleTrim,
		"lower":        parseRuleLower,
		"upper":        parseRuleUpper,
		"title":        parseRuleTitle,
		"replace":      parseRuleReplace,
		"squash-space": parseRuleSquashSpace,
		"remove-space": parseRuleRemoveSpace,
		"empty-to-nil": parseRuleEmptyToNil,
	})
}

// RegisterParseRule registers custom parse rule and function for package.
//
// The custom parse rule can be used in struct tag `parse`, for example:
// `parse:"trim-space|your-rule:param"`.
func RegisterParseRule(rule string, f ParseFunc) {
	if customParseFuncMap[rule] != nil {
		intlog.PrintFunc(context.TODO(), func() string {
			return fmt.Sprintf(
				`parse rule "%s" is overwritten by function "%s"`,
				rule, runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name(),
			)
		})
	}
	customParseFuncMap[rule] = f
}

// RegisterParseRuleByMap registers custom parse rules using map for package.
func RegisterParseRuleByMap(m map[string]ParseFunc) {
	for k, v := range m {
		customParseFuncMap[k] = v
	}
}

// GetRegisteredParseRuleMap returns all the custom registered parse rules and associated functions.
func GetRegisteredParseRuleMap() map[string]ParseFunc {
	if len(customParseFuncMap) == 0 {
		return nil
	}
	ruleMap := make(map[string]ParseFunc)
	for k, v := range customParseFuncMap {
		ruleMap[k] = v
	}
	return ruleMap
}

// DeleteParseRule deletes custom defined parse one or more rules and associated functions from global package.
func DeleteParseRule(rules ...string) {
	for _, rule := range rules {
		delete(customParseFuncMap, rule)
	}
}

func (r *Request) doParseRequestData(data map[string]any, pointer any, mapping ...map[string]string) error {
	if len(data) == 0 || pointer == nil {
		return nil
	}
	structType, err := getParseStructType(pointer)
	if err != nil {
		return err
	}
	var customMapping map[string]string
	if len(mapping) > 0 {
		customMapping = mapping[0]
	}
	return r.doParseMapByStruct(data, structType, customMapping, "")
}

func (r *Request) doParseArrayData(pointer any) ([]map[string]any, error) {
	structType, err := getParseStructArrayItemType(pointer)
	if err != nil {
		return nil, err
	}
	var data []map[string]any
	if err = json.UnmarshalUseNumber(r.GetBody(), &data); err != nil {
		return nil, err
	}
	for i := range data {
		if data[i] == nil {
			data[i] = map[string]any{}
		}
		if err = r.doParseMapByStruct(data[i], structType, nil, fmt.Sprintf("[%d]", i)); err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (r *Request) doParseMapByStruct(
	data map[string]any, structType reflect.Type, mapping map[string]string, fieldPrefix string,
) error {
	structType = indirectToType(structType)
	if structType.Kind() != reflect.Struct || len(data) == 0 {
		return nil
	}
	for i := 0; i < structType.NumField(); i++ {
		structField := structType.Field(i)
		if !isExportedStructField(structField) {
			continue
		}
		fieldPath := joinParseFieldPath(fieldPrefix, structField.Name)
		if structField.Anonymous && structField.Tag == "" {
			if err := r.doParseMapByStruct(data, structField.Type, nil, fieldPrefix); err != nil {
				return err
			}
			continue
		}
		foundKey, foundValue, found := findParseValueFromMap(data, structField, mapping)
		if parseTag := strings.TrimSpace(structField.Tag.Get(gtag.ParseRule)); parseTag != "" {
			rules, err := parseRuleItems(parseTag)
			if err != nil {
				return err
			}
			if found {
				parsedValue, err := r.doParseRuleItems(foundValue, data, structField, fieldPath, rules)
				if err != nil {
					return err
				}
				data[foundKey] = parsedValue
				foundValue = parsedValue
			}
		}
		if found {
			parsedNestedValue, err := r.doParseNestedValue(foundValue, structField.Type, fieldPath)
			if err != nil {
				return err
			}
			data[foundKey] = parsedNestedValue
		}
	}
	return nil
}

func (r *Request) doParseNestedValue(value any, fieldType reflect.Type, fieldPath string) (any, error) {
	indirectType := indirectToType(fieldType)
	switch indirectType.Kind() {
	case reflect.Struct:
		nestedMap, ok := value.(map[string]any)
		if !ok {
			return value, nil
		}
		if err := r.doParseMapByStruct(nestedMap, indirectType, nil, fieldPath); err != nil {
			return value, err
		}
	case reflect.Slice, reflect.Array:
		elemType := indirectToType(indirectType.Elem())
		if elemType.Kind() != reflect.Struct {
			return value, nil
		}
		switch arrayValue := value.(type) {
		case []map[string]any:
			for i := range arrayValue {
				if err := r.doParseMapByStruct(
					arrayValue[i], elemType, nil, fmt.Sprintf("%s[%d]", fieldPath, i),
				); err != nil {
					return value, err
				}
			}
		case []any:
			for i, item := range arrayValue {
				itemMap, ok := item.(map[string]any)
				if !ok {
					continue
				}
				if err := r.doParseMapByStruct(
					itemMap, elemType, nil, fmt.Sprintf("%s[%d]", fieldPath, i),
				); err != nil {
					return value, err
				}
			}
		}
	}
	return value, nil
}

func (r *Request) doParseRuleItems(
	value any, data map[string]any, structField reflect.StructField, fieldPath string, rules []parseRuleItem,
) (any, error) {
	currentValue := value
	for i, rule := range rules {
		if rule.Name == parseRuleForeach {
			return r.doParseForeachRule(currentValue, data, structField, fieldPath, rules[i+1:])
		}
		parseFunc := customParseFuncMap[rule.Name]
		if parseFunc == nil {
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`parse rule "%s" for field "%s" is not registered`,
				rule.Name, fieldPath,
			)
		}
		var err error
		currentValue, err = parseFunc(r.Context(), ParseFuncInput{
			Rule:      rule.Rule,
			Name:      rule.Name,
			Pattern:   rule.Pattern,
			Field:     fieldPath,
			FieldType: structField.Type,
			Value:     currentValue,
			Data:      data,
			Request:   r,
		})
		if err != nil {
			return nil, err
		}
	}
	return currentValue, nil
}

func (r *Request) doParseForeachRule(
	value any, data map[string]any, structField reflect.StructField, fieldPath string, rules []parseRuleItem,
) (any, error) {
	if len(rules) == 0 || value == nil {
		return value, nil
	}
	reflectValue := reflect.ValueOf(value)
	reflectKind := reflectValue.Kind()
	if reflectKind != reflect.Slice && reflectKind != reflect.Array {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`parse rule "%s" for field "%s" requires slice/array value, but got "%T"`,
			parseRuleForeach, fieldPath, value,
		)
	}
	parsedValues := make([]any, reflectValue.Len())
	for i := 0; i < reflectValue.Len(); i++ {
		parsedItem, err := r.doParseRuleItems(
			reflectValue.Index(i).Interface(),
			data,
			structField,
			fmt.Sprintf("%s[%d]", fieldPath, i),
			rules,
		)
		if err != nil {
			return nil, err
		}
		parsedValues[i] = parsedItem
	}
	return rebuildParsedArrayValue(value, parsedValues)
}

func parseRuleItems(tagValue string) ([]parseRuleItem, error) {
	array := strings.Split(tagValue, "|")
	items := make([]parseRuleItem, 0, len(array))
	for _, item := range array {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		name, pattern, _ := strings.Cut(item, ":")
		name = strings.TrimSpace(strings.ToLower(name))
		if name == "" {
			return nil, gerror.NewCode(
				gcode.CodeInvalidParameter,
				"parse rule name cannot be empty",
			)
		}
		items = append(items, parseRuleItem{
			Rule:    item,
			Name:    name,
			Pattern: strings.TrimSpace(pattern),
		})
	}
	return items, nil
}

func findParseValueFromMap(
	data map[string]any, structField reflect.StructField, mapping map[string]string,
) (foundKey string, foundValue any, found bool) {
	for _, candidate := range buildParseLookupKeys(structField, mapping) {
		if candidate == "" {
			continue
		}
		if foundKey, foundValue = gutil.MapPossibleItemByKey(data, candidate); foundKey != "" {
			return foundKey, foundValue, true
		}
	}
	return "", nil, false
}

func buildParseLookupKeys(structField reflect.StructField, mapping map[string]string) []string {
	keys := make([]string, 0, len(gtag.StructTagPriority)+2)
	if len(mapping) > 0 {
		for paramKey, attrName := range mapping {
			if attrName == structField.Name {
				keys = append(keys, paramKey)
			}
		}
	}
	for _, tagName := range gtag.StructTagPriority {
		if tagValue := structField.Tag.Get(tagName); tagValue != "" {
			tagValue = strings.Split(tagValue, ",")[0]
			if tagValue != "" && tagValue != "-" {
				keys = append(keys, tagValue)
			}
		}
	}
	keys = append(keys, structField.Name)
	return uniqueParseLookupKeys(keys)
}

func joinParseFieldPath(prefix, name string) string {
	if prefix == "" {
		return name
	}
	if strings.HasPrefix(name, "[") {
		return prefix + name
	}
	return prefix + "." + name
}

func uniqueParseLookupKeys(keys []string) []string {
	if len(keys) == 0 {
		return nil
	}
	var (
		uniqueKeys = make([]string, 0, len(keys))
		existsMap  = make(map[string]struct{}, len(keys))
	)
	for _, key := range keys {
		if key == "" {
			continue
		}
		if _, ok := existsMap[key]; ok {
			continue
		}
		existsMap[key] = struct{}{}
		uniqueKeys = append(uniqueKeys, key)
	}
	return uniqueKeys
}

func getParseStructType(pointer any) (reflect.Type, error) {
	reflectType := reflect.TypeOf(pointer)
	if reflectType == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "destination pointer cannot be nil")
	}
	if reflectType.Kind() != reflect.Pointer {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`destination pointer should be type of "*struct", but got "%v"`,
			reflectType.Kind(),
		)
	}
	return indirectToType(reflectType.Elem()), nil
}

func getParseStructArrayItemType(pointer any) (reflect.Type, error) {
	reflectType := reflect.TypeOf(pointer)
	if reflectType == nil {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "destination pointer cannot be nil")
	}
	if reflectType.Kind() != reflect.Pointer {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`destination pointer should be type of "*[]struct" or "*[]*struct", but got "%v"`,
			reflectType.Kind(),
		)
	}
	reflectType = indirectToType(reflectType.Elem())
	if reflectType.Kind() != reflect.Slice && reflectType.Kind() != reflect.Array {
		return nil, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`destination pointer should be type of "*[]struct" or "*[]*struct", but got "%v"`,
			reflectType.Kind(),
		)
	}
	return indirectToType(reflectType.Elem()), nil
}

func indirectToType(reflectType reflect.Type) reflect.Type {
	for reflectType != nil && reflectType.Kind() == reflect.Pointer {
		reflectType = reflectType.Elem()
	}
	return reflectType
}

func isExportedStructField(structField reflect.StructField) bool {
	return structField.PkgPath == ""
}

func rebuildParsedArrayValue(value any, parsedValues []any) (any, error) {
	reflectValue := reflect.ValueOf(value)
	switch reflectValue.Kind() {
	case reflect.Slice:
		sliceValue := reflect.MakeSlice(reflectValue.Type(), len(parsedValues), len(parsedValues))
		for i, item := range parsedValues {
			if item == nil {
				sliceValue.Index(i).Set(reflect.Zero(reflectValue.Type().Elem()))
				continue
			}
			convertedValue, err := convertParsedArrayItemToType(item, reflectValue.Type().Elem())
			if err != nil {
				return nil, err
			}
			sliceValue.Index(i).Set(convertedValue)
		}
		return sliceValue.Interface(), nil
	case reflect.Array:
		arrayValue := reflect.New(reflectValue.Type()).Elem()
		for i := 0; i < len(parsedValues) && i < arrayValue.Len(); i++ {
			if parsedValues[i] == nil {
				arrayValue.Index(i).Set(reflect.Zero(reflectValue.Type().Elem()))
				continue
			}
			convertedValue, err := convertParsedArrayItemToType(parsedValues[i], reflectValue.Type().Elem())
			if err != nil {
				return nil, err
			}
			arrayValue.Index(i).Set(convertedValue)
		}
		return arrayValue.Interface(), nil
	}
	return value, nil
}

func convertParsedArrayItemToType(value any, targetType reflect.Type) (reflect.Value, error) {
	if value == nil {
		return reflect.Zero(targetType), nil
	}
	if reflectValue := reflect.ValueOf(value); reflectValue.IsValid() {
		if reflectValue.Type().AssignableTo(targetType) {
			return reflectValue, nil
		}
		if reflectValue.Type().ConvertibleTo(targetType) {
			return reflectValue.Convert(targetType), nil
		}
	}
	targetValue := reflect.New(targetType)
	if err := gconv.Scan(value, targetValue.Interface()); err != nil {
		return reflect.Value{}, err
	}
	return targetValue.Elem(), nil
}

func getParseStringValue(in ParseFuncInput) (string, bool, error) {
	if in.Value == nil {
		return "", false, nil
	}
	switch value := in.Value.(type) {
	case string:
		return value, true, nil
	default:
		return "", false, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`parse rule "%s" for field "%s" does not support value type "%T"`,
			in.Rule, in.Field, in.Value,
		)
	}
}

func parseRuleTrimSpace(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.TrimSpace(value), nil
}

func parseRuleTrimLeft(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.TrimLeft(value, in.Pattern), nil
}

func parseRuleTrimRight(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.TrimRight(value, in.Pattern), nil
}

func parseRuleTrimPrefix(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.TrimPrefix(value, in.Pattern), nil
}

func parseRuleTrimSuffix(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.TrimSuffix(value, in.Pattern), nil
}

func parseRuleTrim(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.Trim(value, in.Pattern), nil
}

func parseRuleLower(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.ToLower(value), nil
}

func parseRuleUpper(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.ToUpper(value), nil
}

func parseRuleTitle(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return gstr.UcWords(value), nil
}

func parseRuleReplace(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	oldValue, newValue, ok := strings.Cut(in.Pattern, ",")
	if !ok {
		oldValue = in.Pattern
	}
	return strings.ReplaceAll(value, oldValue, newValue), nil
}

func parseRuleSquashSpace(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.Join(strings.Fields(value), " "), nil
}

func parseRuleRemoveSpace(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	return strings.Join(strings.Fields(value), ""), nil
}

func parseRuleEmptyToNil(_ context.Context, in ParseFuncInput) (any, error) {
	value, ok, err := getParseStringValue(in)
	if !ok || err != nil {
		return in.Value, err
	}
	if value == "" {
		return nil, nil
	}
	return value, nil
}
