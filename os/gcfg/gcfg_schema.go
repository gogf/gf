// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gcfg provides configuration management functionality for GoFrame.
// This file implements configuration schema registry for visual editing support.
package gcfg

import (
	"reflect"
	"strings"
	"sync"
)

// FieldSchema describes metadata for a single configuration field,
// extracted from struct tags (json, d, v, dc) via reflection.
type FieldSchema struct {
	Name        string   `json:"name"`              // Go struct field name
	JsonKey     string   `json:"jsonKey"`            // JSON/YAML key from json tag
	Type        string   `json:"type"`               // Field type: string, int, bool, duration, etc.
	Default     string   `json:"default"`            // Default value from `d` tag
	Rule        string   `json:"rule"`               // Validation rule from `v` tag
	Description string   `json:"description"`        // English description from `dc` tag
	I18nKey     string   `json:"i18nKey"`            // I18n key extracted from `dc` tag (i18n:xxx)
	Group       string   `json:"group"`              // Logical group (Basic, Logging, Cookie, etc.)
	Options     []string `json:"options,omitempty"`   // Enum options if applicable
}

// ModuleSchema describes the configuration schema for one module.
type ModuleSchema struct {
	Name       string         `json:"name"`       // Module name: server, database, redis, logger, viewer
	ConfigNode string         `json:"configNode"` // Config file node name
	Fields     []*FieldSchema `json:"fields"`     // All field schemas
	Groups     []string       `json:"groups"`     // Ordered unique group names
}

// SchemaRegistry is the global registry for all module configuration schemas.
// It is thread-safe and supports concurrent registration and retrieval.
type SchemaRegistry struct {
	mu      sync.RWMutex
	schemas map[string]*ModuleSchema
	order   []string // maintains registration order
}

// globalSchemaRegistry is the package-level global schema registry instance.
var globalSchemaRegistry = NewSchemaRegistry()

// NewSchemaRegistry creates and returns a new SchemaRegistry instance.
func NewSchemaRegistry() *SchemaRegistry {
	return &SchemaRegistry{
		schemas: make(map[string]*ModuleSchema),
	}
}

// RegisterSchema registers a module's configuration struct type to the global registry.
func RegisterSchema(name, configNode string, configStruct any, groupMap map[string]string) {
	globalSchemaRegistry.Register(name, configNode, configStruct, groupMap)
}

// GetSchema returns the ModuleSchema for a given module name from the global registry.
func GetSchema(name string) (*ModuleSchema, bool) {
	return globalSchemaRegistry.Get(name)
}

// GetAllSchemas returns all registered module schemas from the global registry.
func GetAllSchemas() []*ModuleSchema {
	return globalSchemaRegistry.GetAll()
}

// GetGlobalRegistry returns the package-level global schema registry.
func GetGlobalRegistry() *SchemaRegistry {
	return globalSchemaRegistry
}

// Register registers a module's configuration schema to this registry.
func (r *SchemaRegistry) Register(name, configNode string, configStruct any, groupMap map[string]string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	fields := scanStructTags(configStruct, groupMap)
	groups := extractGroups(fields)

	schema := &ModuleSchema{
		Name:       name,
		ConfigNode: configNode,
		Fields:     fields,
		Groups:     groups,
	}

	if _, exists := r.schemas[name]; !exists {
		r.order = append(r.order, name)
	}
	r.schemas[name] = schema
}

// Get returns the ModuleSchema for a given module name.
func (r *SchemaRegistry) Get(name string) (*ModuleSchema, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	schema, ok := r.schemas[name]
	return schema, ok
}

// GetAll returns all registered module schemas in registration order.
func (r *SchemaRegistry) GetAll() []*ModuleSchema {
	r.mu.RLock()
	defer r.mu.RUnlock()
	result := make([]*ModuleSchema, 0, len(r.order))
	for _, name := range r.order {
		if schema, ok := r.schemas[name]; ok {
			result = append(result, schema)
		}
	}
	return result
}

// scanStructTags scans a struct type via reflection and extracts FieldSchema
// from struct tags (json, d, v, dc).
func scanStructTags(configType any, groupMap map[string]string) []*FieldSchema {
	t := reflect.TypeOf(configType)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return nil
	}
	return scanStructFields(t, groupMap, "")
}

// scanStructFields recursively scans struct fields and returns FieldSchema list.
func scanStructFields(t reflect.Type, groupMap map[string]string, prefix string) []*FieldSchema {
	var fields []*FieldSchema
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Skip unexported fields.
		if !field.IsExported() {
			continue
		}

		// Handle embedded structs: recurse into them.
		if field.Anonymous {
			ft := field.Type
			if ft.Kind() == reflect.Ptr {
				ft = ft.Elem()
			}
			if ft.Kind() == reflect.Struct {
				fields = append(fields, scanStructFields(ft, groupMap, prefix)...)
			}
			continue
		}

		// Skip fields whose type is interface, func, chan, or complex struct.
		ft := field.Type
		if ft.Kind() == reflect.Ptr {
			ft = ft.Elem()
		}
		switch ft.Kind() {
		case reflect.Interface, reflect.Func, reflect.Chan:
			continue
		case reflect.Struct:
			// Allow structs from the "time" package (e.g. time.Time);
			// note that time.Duration is int64 and is handled in the Int64 case above.
			if ft.PkgPath() != "" && ft.PkgPath() != "time" {
				continue
			}
		}

		fs := parseFieldSchema(field, groupMap, prefix)
		if fs != nil {
			fields = append(fields, fs)
		}
	}
	return fields
}

// parseFieldSchema parses a single struct field into a FieldSchema.
func parseFieldSchema(field reflect.StructField, groupMap map[string]string, prefix string) *FieldSchema {
	// Get json key.
	jsonKey := ""
	if jsonTag := field.Tag.Get("json"); jsonTag != "" {
		parts := strings.Split(jsonTag, ",")
		if parts[0] != "-" {
			jsonKey = parts[0]
		} else {
			// Skip fields with json:"-"
			return nil
		}
	}
	if jsonKey == "" {
		jsonKey = lowerFirst(field.Name)
	}

	if prefix != "" {
		jsonKey = prefix + "." + jsonKey
	}

	typeName := fieldTypeName(field.Type)
	defaultVal := field.Tag.Get("d")
	rule := field.Tag.Get("v")
	description, i18nKey := parseDcTag(field.Tag.Get("dc"))

	group := "Other"
	if groupMap != nil {
		if g, ok := groupMap[field.Name]; ok {
			group = g
		}
	}

	return &FieldSchema{
		Name:        field.Name,
		JsonKey:     jsonKey,
		Type:        typeName,
		Default:     defaultVal,
		Rule:        rule,
		Description: description,
		I18nKey:     i18nKey,
		Group:       group,
	}
}

// parseDcTag parses the `dc` tag value into description and i18n key.
func parseDcTag(dc string) (description, i18nKey string) {
	if dc == "" {
		return "", ""
	}
	parts := strings.SplitN(dc, "|", 2)
	description = strings.TrimSpace(parts[0])
	if len(parts) > 1 {
		suffix := strings.TrimSpace(parts[1])
		if strings.HasPrefix(suffix, "i18n:") {
			i18nKey = strings.TrimPrefix(suffix, "i18n:")
		}
	}
	return
}

// fieldTypeName returns a human-readable type name for a reflect.Type.
func fieldTypeName(t reflect.Type) string {
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	switch t.Kind() {
	case reflect.String:
		return "string"
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if t.PkgPath() == "time" && t.Name() == "Duration" {
			return "duration"
		}
		return "int"
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return "int"
	case reflect.Float32, reflect.Float64:
		return "float"
	case reflect.Slice:
		return "[]" + fieldTypeName(t.Elem())
	case reflect.Map:
		return "map"
	default:
		return t.String()
	}
}

// lowerFirst returns the string with first character lowered.
func lowerFirst(s string) string {
	if s == "" {
		return s
	}
	return strings.ToLower(s[:1]) + s[1:]
}

// extractGroups returns ordered unique group names from field schemas.
func extractGroups(fields []*FieldSchema) []string {
	seen := make(map[string]bool)
	var groups []string
	for _, f := range fields {
		if !seen[f.Group] {
			seen[f.Group] = true
			groups = append(groups, f.Group)
		}
	}
	return groups
}
