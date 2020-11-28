// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/util/gconv"
)

// SetForm sets custom form value with key-value pair.
func (r *Request) SetForm(key string, value interface{}) {
	r.parseForm()
	if r.formMap == nil {
		r.formMap = make(map[string]interface{})
	}
	r.formMap[key] = value
}

// GetForm retrieves and returns parameter <key> from form.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetForm(key string, def ...interface{}) interface{} {
	r.parseForm()
	if len(r.formMap) > 0 {
		if v, ok := r.formMap[key]; ok {
			return v
		}
	}
	if len(def) > 0 {
		return def[0]
	}
	return nil
}

// GetFormVar retrieves and returns parameter <key> from form as Var.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormVar(key string, def ...interface{}) *gvar.Var {
	return gvar.New(r.GetForm(key, def...))
}

// GetFormString retrieves and returns parameter <key> from form as string.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormString(key string, def ...interface{}) string {
	return r.GetFormVar(key, def...).String()
}

// GetFormBool retrieves and returns parameter <key> from form as bool.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormBool(key string, def ...interface{}) bool {
	return r.GetFormVar(key, def...).Bool()
}

// GetFormInt retrieves and returns parameter <key> from form as int.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormInt(key string, def ...interface{}) int {
	return r.GetFormVar(key, def...).Int()
}

// GetFormInt32 retrieves and returns parameter <key> from form as int32.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormInt32(key string, def ...interface{}) int32 {
	return r.GetFormVar(key, def...).Int32()
}

// GetFormInt64 retrieves and returns parameter <key> from form as int64.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormInt64(key string, def ...interface{}) int64 {
	return r.GetFormVar(key, def...).Int64()
}

// GetFormInts retrieves and returns parameter <key> from form as []int.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormInts(key string, def ...interface{}) []int {
	return r.GetFormVar(key, def...).Ints()
}

// GetFormUint retrieves and returns parameter <key> from form as uint.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormUint(key string, def ...interface{}) uint {
	return r.GetFormVar(key, def...).Uint()
}

// GetFormUint32 retrieves and returns parameter <key> from form as uint32.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormUint32(key string, def ...interface{}) uint32 {
	return r.GetFormVar(key, def...).Uint32()
}

// GetFormUint64 retrieves and returns parameter <key> from form as uint64.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormUint64(key string, def ...interface{}) uint64 {
	return r.GetFormVar(key, def...).Uint64()
}

// GetFormFloat32 retrieves and returns parameter <key> from form as float32.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormFloat32(key string, def ...interface{}) float32 {
	return r.GetFormVar(key, def...).Float32()
}

// GetFormFloat64 retrieves and returns parameter <key> from form as float64.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormFloat64(key string, def ...interface{}) float64 {
	return r.GetFormVar(key, def...).Float64()
}

// GetFormFloats retrieves and returns parameter <key> from form as []float64.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormFloats(key string, def ...interface{}) []float64 {
	return r.GetFormVar(key, def...).Floats()
}

// GetFormArray retrieves and returns parameter <key> from form as []string.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormArray(key string, def ...interface{}) []string {
	return r.GetFormVar(key, def...).Strings()
}

// GetFormStrings retrieves and returns parameter <key> from form as []string.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormStrings(key string, def ...interface{}) []string {
	return r.GetFormVar(key, def...).Strings()
}

// GetFormInterfaces retrieves and returns parameter <key> from form as []interface{}.
// It returns <def> if <key> does not exist in the form and <def> is given, or else it returns nil.
func (r *Request) GetFormInterfaces(key string, def ...interface{}) []interface{} {
	return r.GetFormVar(key, def...).Interfaces()
}

// GetFormMap retrieves and returns all form parameters passed from client as map.
// The parameter <kvMap> specifies the keys retrieving from client parameters,
// the associated values are the default values if the client does not pass.
func (r *Request) GetFormMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseForm()
	if len(kvMap) > 0 && kvMap[0] != nil {
		if len(r.formMap) == 0 {
			return kvMap[0]
		}
		m := make(map[string]interface{}, len(kvMap[0]))
		for k, defValue := range kvMap[0] {
			if postValue, ok := r.formMap[k]; ok {
				m[k] = postValue
			} else {
				m[k] = defValue
			}
		}
		return m
	} else {
		return r.formMap
	}
}

// GetFormMapStrStr retrieves and returns all form parameters passed from client as map[string]string.
// The parameter <kvMap> specifies the keys retrieving from client parameters, the associated values
// are the default values if the client does not pass.
func (r *Request) GetFormMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	postMap := r.GetFormMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]string, len(postMap))
		for k, v := range postMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetFormMapStrVar retrieves and returns all form parameters passed from client as map[string]*gvar.Var.
// The parameter <kvMap> specifies the keys retrieving from client parameters, the associated values
// are the default values if the client does not pass.
func (r *Request) GetFormMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	postMap := r.GetFormMap(kvMap...)
	if len(postMap) > 0 {
		m := make(map[string]*gvar.Var, len(postMap))
		for k, v := range postMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetFormStruct retrieves all form parameters passed from client and converts them to
// given struct object. Note that the parameter <pointer> is a pointer to the struct object.
// The optional parameter <mapping> is used to specify the key to attribute mapping.
func (r *Request) GetFormStruct(pointer interface{}, mapping ...map[string]string) error {
	r.parseForm()
	data := r.formMap
	if data == nil {
		data = map[string]interface{}{}
	}
	if err := r.mergeDefaultStructValue(data, pointer); err != nil {
		return nil
	}
	return gconv.Struct(data, pointer, mapping...)
}
