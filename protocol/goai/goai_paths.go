// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package goai

import (
	"github.com/gogf/gf/v2/container/gmap"
)

// Paths are specified by OpenAPI/Swagger standard version 3.0.
type Paths struct {
	paths *gmap.ListMap // map[string]Path
}

func (p *Paths) init() {
	if p.paths == nil {
		p.paths = gmap.NewListMap()
	}
}

func (p *Paths) Get(name string) *Path {
	p.init()
	value := p.paths.Get(name)
	if value != nil {
		path := value.(Path)
		return &path
	}
	return nil
}

func (p *Paths) Set(name string, path Path) {
	p.init()
	p.paths.Set(name, path)
}

func (p *Paths) Map() map[string]Path {
	p.init()
	m := make(map[string]Path)
	p.paths.Iterator(func(key, value interface{}) bool {
		m[key.(string)] = value.(Path)
		return true
	})
	return m
}

func (p Paths) MarshalJSON() ([]byte, error) {
	p.init()
	return p.paths.MarshalJSON()
}
