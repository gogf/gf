// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package structs

import "reflect"

// Tag returns the value associated with key in the tag string. If there is no
// such key in the tag, Tag returns the empty string.
func (f *Field) Tag(key string) string {
	return f.Field.Tag.Get(key)
}

// TagLookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise, the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string. If the tag does not have the conventional format,
// the value returned by Lookup is unspecified.
func (f *Field) TagLookup(key string) (value string, ok bool) {
	return f.Field.Tag.Lookup(key)
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f *Field) IsEmbedded() bool {
	return f.Field.Anonymous
}

// TagStr returns the tag string of the field.
func (f *Field) TagStr() string {
	return string(f.Field.Tag)
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.Field.PkgPath == ""
}

// Name returns the name of the given field
func (f *Field) Name() string {
	return f.Field.Name
}

// Type returns the type of the given field
func (f *Field) Type() Type {
	return Type{
		Type: f.Field.Type,
	}
}

// Kind returns the reflect.Kind for Value of Field `f`.
func (f *Field) Kind() reflect.Kind {
	return f.Value.Kind()
}

// OriginalKind retrieves and returns the original reflect.Kind for Value of Field `f`.
func (f *Field) OriginalKind() reflect.Kind {
	var (
		kind  = f.Value.Kind()
		value = f.Value
	)
	for kind == reflect.Ptr {
		value = value.Elem()
		kind = value.Kind()
	}
	return kind
}

const (
	RecursiveOptionNone          = 0 // No recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbedded      = 1 // Recursively retrieving fields as map if the field is an embedded struct.
	RecursiveOptionEmbeddedNoTag = 2 // Recursively retrieving fields as map if the field is an embedded struct and the field has no tag.
)

type FieldsInput struct {
	// Pointer should be type of struct/*struct.
	Pointer interface{}

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption int
}

type FieldMapInput struct {
	// Pointer should be type of struct/*struct.
	Pointer interface{}

	// PriorityTagArray specifies the priority tag array for retrieving from high to low.
	// If it's given `nil`, it returns map[name]*Field, of which the `name` is attribute name.
	PriorityTagArray []string

	// RecursiveOption specifies the way retrieving the fields recursively if the attribute
	// is an embedded struct. It is RecursiveOptionNone in default.
	RecursiveOption int
}

// Fields retrieves and returns the fields of `pointer` as slice.
func Fields(in FieldsInput) ([]*Field, error) {
	var (
		fieldFilterMap  = make(map[string]struct{})
		retrievedFields = make([]*Field, 0)
	)
	rangeFields, err := getFieldValues(in.Pointer)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(rangeFields); i++ {
		field := rangeFields[i]
		if _, ok := fieldFilterMap[field.Name()]; ok {
			continue
		}
		// Only retrieve exported attributes.
		if !field.IsExported() {
			continue
		}
		if field.IsEmbedded() {
			if in.RecursiveOption != RecursiveOptionNone {
				switch in.RecursiveOption {
				case RecursiveOptionEmbeddedNoTag:
					if field.TagStr() != "" {
						break
					}
					fallthrough
				case RecursiveOptionEmbedded:
					structFields, err := Fields(FieldsInput{
						Pointer:         field.Value,
						RecursiveOption: in.RecursiveOption,
					})
					if err != nil {
						return nil, err
					}
					rangeFields = append(rangeFields, structFields...)
				}
			}
			continue
		}
		fieldFilterMap[field.Name()] = struct{}{}
		retrievedFields = append(retrievedFields, field)
	}
	return retrievedFields, nil
}

// FieldMap retrieves and returns struct field as map[name/tag]*Field from `pointer`.
//
// The parameter `pointer` should be type of struct/*struct.
//
// The parameter `priority` specifies the priority tag array for retrieving from high to low.
// If it's given `nil`, it returns map[name]*Field, of which the `name` is attribute name.
//
// The parameter `recursive` specifies the whether retrieving the fields recursively if the attribute
// is an embedded struct.
//
// Note that it only retrieves the exported attributes with first letter up-case from struct.
func FieldMap(in FieldMapInput) (map[string]*Field, error) {
	fields, err := getFieldValues(in.Pointer)
	if err != nil {
		return nil, err
	}
	var (
		tagValue = ""
		mapField = make(map[string]*Field)
	)
	for _, field := range fields {
		// Only retrieve exported attributes.
		if !field.IsExported() {
			continue
		}
		tagValue = ""
		for _, p := range in.PriorityTagArray {
			tagValue = field.Tag(p)
			if tagValue != "" && tagValue != "-" {
				break
			}
		}
		tempField := field
		tempField.TagValue = tagValue
		if tagValue != "" {
			mapField[tagValue] = tempField
		} else {
			if in.RecursiveOption != RecursiveOptionNone && field.IsEmbedded() {
				switch in.RecursiveOption {
				case RecursiveOptionEmbeddedNoTag:
					if field.TagStr() != "" {
						mapField[field.Name()] = tempField
						break
					}
					fallthrough
				case RecursiveOptionEmbedded:
					m, err := FieldMap(FieldMapInput{
						Pointer:          field.Value,
						PriorityTagArray: in.PriorityTagArray,
						RecursiveOption:  in.RecursiveOption,
					})
					if err != nil {
						return nil, err
					}
					for k, v := range m {
						if _, ok := mapField[k]; !ok {
							tempV := v
							mapField[k] = tempV
						}
					}
				}
			} else {
				mapField[field.Name()] = tempField
			}
		}
	}
	return mapField, nil
}
