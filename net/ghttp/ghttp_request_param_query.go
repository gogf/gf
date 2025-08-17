// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import (
	"net/http"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/util/gconv"
)

// SetQuery sets custom query value with key-value pairs.
func (r *Request) SetQuery(key string, value interface{}) {
	r.parseQuery()
	if r.queryMap == nil {
		r.queryMap = make(map[string]interface{})
	}
	r.queryMap[key] = value
}

// GetQuery retrieves and return parameter with the given name `key` from query string
// and request body. It returns `def` if `key` does not exist in the query and `def` is given,
// or else it returns nil.
//
// Note that if there are multiple parameters with the same name, the parameters are retrieved
// and overwrote in order of priority: query > body.
func (r *Request) GetQuery(key string, def ...interface{}) *gvar.Var {
	r.parseQuery()
	if len(r.queryMap) > 0 {
		if value, ok := r.queryMap[key]; ok {
			return gvar.New(value)
		}
	}
	if r.Method == http.MethodGet {
		r.parseBody()
	}
	if len(r.bodyMap) > 0 {
		if v, ok := r.bodyMap[key]; ok {
			return gvar.New(v)
		}
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetQueryMap retrieves and returns all parameters passed from the client using HTTP GET method
// as the map. The parameter `kvMap` specifies the keys retrieving from client parameters,
// the associated values are the default values if the client does not pass.
//
// Note that if there are multiple parameters with the same name, the parameters are retrieved and overwrote
// in order of priority: query > body.
func (r *Request) GetQueryMap(kvMap ...map[string]interface{}) map[string]interface{} {
	r.parseQuery()
	if r.Method == http.MethodGet {
		r.parseBody()
	}
	var m map[string]interface{}
	if len(kvMap) > 0 && kvMap[0] != nil {
		if len(r.queryMap) == 0 && len(r.bodyMap) == 0 {
			return kvMap[0]
		}
		m = make(map[string]interface{}, len(kvMap[0]))
		if len(r.bodyMap) > 0 {
			for k, v := range kvMap[0] {
				if postValue, ok := r.bodyMap[k]; ok {
					m[k] = postValue
				} else {
					m[k] = v
				}
			}
		}
		if len(r.queryMap) > 0 {
			for k, v := range kvMap[0] {
				if postValue, ok := r.queryMap[k]; ok {
					m[k] = postValue
				} else {
					m[k] = v
				}
			}
		}
	} else {
		m = make(map[string]interface{}, len(r.queryMap)+len(r.bodyMap))
		for k, v := range r.bodyMap {
			m[k] = v
		}
		for k, v := range r.queryMap {
			m[k] = v
		}
	}
	return m
}

// GetQueryMapStrStr retrieves and returns all parameters passed from the client using the HTTP GET method as a
//
//	map[string]string. The parameter `kvMap` specifies the keys
//
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrStr(kvMap ...map[string]interface{}) map[string]string {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]string, len(queryMap))
		for k, v := range queryMap {
			m[k] = gconv.String(v)
		}
		return m
	}
	return nil
}

// GetQueryMapStrVar retrieves and returns all parameters passed from the client using the HTTP GET method
// as map[string]*gvar.Var. The parameter `kvMap` specifies the keys
// retrieving from client parameters, the associated values are the default values if the client
// does not pass.
func (r *Request) GetQueryMapStrVar(kvMap ...map[string]interface{}) map[string]*gvar.Var {
	queryMap := r.GetQueryMap(kvMap...)
	if len(queryMap) > 0 {
		m := make(map[string]*gvar.Var, len(queryMap))
		for k, v := range queryMap {
			m[k] = gvar.New(v)
		}
		return m
	}
	return nil
}

// GetQueryStruct retrieves all parameters passed from the client using the HTTP GET method
// and converts them to a given struct object. Note that the parameter `pointer` is a pointer
// to the struct object. The optional parameter `mapping` is used to specify the key to
// attribute mapping.
func (r *Request) GetQueryStruct(pointer interface{}, mapping ...map[string]string) error {
	_, err := r.doGetQueryStruct(pointer, mapping...)
	return err
}

func (r *Request) doGetQueryStruct(pointer interface{}, mapping ...map[string]string) (data map[string]interface{}, err error) {
	r.parseQuery()
	data = r.GetQueryMap()
	if data == nil {
		data = map[string]interface{}{}
	}
	if err = r.mergeDefaultStructValue(data, pointer); err != nil {
		return data, nil
	}
	return data, gconv.Struct(data, pointer, mapping...)
}
