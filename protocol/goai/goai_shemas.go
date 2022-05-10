// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/container/gmap"
)

type Schemas struct {
	refs *gmap.ListMap // map[string]SchemaRef
}

func createSchemas() Schemas {
	return Schemas{
		refs: gmap.NewListMap(),
	}
}

func (s *Schemas) init() {
	if s.refs == nil {
		s.refs = gmap.NewListMap()
	}
}

func (s *Schemas) Get(name string) *SchemaRef {
	s.init()
	value := s.refs.Get(name)
	if value != nil {
		ref := value.(SchemaRef)
		return &ref
	}
	return nil
}

func (s *Schemas) Set(name string, ref SchemaRef) {
	s.init()
	s.refs.Set(name, ref)
}

func (s *Schemas) Removes(names []interface{}) {
	s.init()
	s.refs.Removes(names)
}

func (s *Schemas) Map() map[string]SchemaRef {
	s.init()
	m := make(map[string]SchemaRef)
	s.refs.Iterator(func(key, value interface{}) bool {
		m[key.(string)] = value.(SchemaRef)
		return true
	})
	return m
}

func (s *Schemas) Iterator(f func(key string, ref SchemaRef) bool) {
	s.init()
	s.refs.Iterator(func(key, value interface{}) bool {
		return f(key.(string), value.(SchemaRef))
	})
}

func (s Schemas) MarshalJSON() ([]byte, error) {
	s.init()
	return s.refs.MarshalJSON()
}
