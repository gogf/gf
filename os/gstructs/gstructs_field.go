// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gstructs

import (
	"reflect"

	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/internal/utils"
	"github.com/gogf/gf/v2/util/gtag"
)

// Tag returns the value associated with key in the tag string. If there is no
// such key in the tag, Tag returns the empty string.
func (f *Field) Tag(key string) string {
	s := f.Field.Tag.Get(key)
	if s != "" {
		s = gtag.Parse(s)
	}
	return s
}

// TagLookup returns the value associated with key in the tag string.
// If the key is present in the tag the value (which may be empty)
// is returned. Otherwise, the returned value will be the empty string.
// The ok return value reports whether the value was explicitly set in
// the tag string. If the tag does not have the conventional format,
// the value returned by Lookup is unspecified.
func (f *Field) TagLookup(key string) (value string, ok bool) {
	value, ok = f.Field.Tag.Lookup(key)
	if ok && value != "" {
		value = gtag.Parse(value)
	}
	return
}

// IsEmbedded returns true if the given field is an anonymous field (embedded)
func (f *Field) IsEmbedded() bool {
	return f.Field.Anonymous
}

// TagStr returns the tag string of the field.
func (f *Field) TagStr() string {
	return string(f.Field.Tag)
}

// TagMap returns all the tag of the field along with its value string as map.
func (f *Field) TagMap() map[string]string {
	var (
		data = ParseTag(f.TagStr())
	)
	for k, v := range data {
		data[k] = utils.StripSlashes(gtag.Parse(v))
	}
	return data
}

// IsExported returns true if the given field is exported.
func (f *Field) IsExported() bool {
	return f.Field.PkgPath == ""
}

// Name returns the name of the given field.
func (f *Field) Name() string {
	return f.Field.Name
}

// Type returns the type of the given field.
// Note that this Type is not reflect.Type. If you need reflect.Type, please use Field.Type().Type.
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
		reflectType = f.Value.Type()
		reflectKind = reflectType.Kind()
	)
	for reflectKind == reflect.Ptr {
		reflectType = reflectType.Elem()
		reflectKind = reflectType.Kind()
	}

	return reflectKind
}

// OriginalValue retrieves and returns the original reflect.Value of Field `f`.
func (f *Field) OriginalValue() reflect.Value {
	var (
		reflectValue = f.Value
		reflectType  = reflectValue.Type()
		reflectKind  = reflectType.Kind()
	)

	for reflectKind == reflect.Ptr && !f.IsNil() {
		reflectValue = reflectValue.Elem()
		reflectKind = reflectValue.Type().Kind()
	}

	return reflectValue
}

// IsEmpty checks and returns whether the value of this Field is empty.
func (f *Field) IsEmpty() bool {
	return empty.IsEmpty(f.Value)
}

// IsNil checks and returns whether the value of this Field is nil.
func (f *Field) IsNil(traceSource ...bool) bool {
	return empty.IsNil(f.Value, traceSource...)
}
