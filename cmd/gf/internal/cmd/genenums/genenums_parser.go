// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package genenums

import (
	"go/constant"
	"go/types"
	"golang.org/x/tools/go/packages"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

const pkgLoadMode = 0xffffff

type EnumsParser struct {
	enums     []EnumItem
	parsedPkg map[string]struct{}
	prefix    string
}

type EnumItem struct {
	Name  string
	Value string
	Kind  constant.Kind // String/Int/Bool/Float/Complex/Unknown
	Type  string        // Pkg.ID + TypeName
}

var standardPackages = make(map[string]struct{})

func init() {
	stdPackages, err := packages.Load(nil, "std")
	if err != nil {
		panic(err)
	}
	for _, p := range stdPackages {
		standardPackages[p.ID] = struct{}{}
	}
}

func NewEnumsParser(prefix string) *EnumsParser {
	return &EnumsParser{
		enums:     make([]EnumItem, 0),
		parsedPkg: make(map[string]struct{}),
		prefix:    prefix,
	}
}

func (p *EnumsParser) ParsePackages(pkgs []*packages.Package) {
	for _, pkg := range pkgs {
		p.ParsePackage(pkg)
	}
}

func (p *EnumsParser) ParsePackage(pkg *packages.Package) {
	// Ignore std packages.
	if _, ok := standardPackages[pkg.ID]; ok {
		return
	}
	// Ignore pared packages.
	if _, ok := p.parsedPkg[pkg.ID]; ok {
		return
	}
	p.parsedPkg[pkg.ID] = struct{}{}
	// Only parse specified prefix.
	if p.prefix != "" {
		if !gstr.HasPrefix(pkg.ID, p.prefix) {
			return
		}
	}

	var (
		scope = pkg.Types.Scope()
		names = scope.Names()
	)
	for _, name := range names {
		con, ok := scope.Lookup(name).(*types.Const)
		if !ok {
			// Only constants can be enums.
			continue
		}
		if !con.Exported() {
			// Ignore unexported values.
			continue
		}

		var enumType = con.Type().String()
		if !gstr.Contains(enumType, "/") {
			// Ignore std types.
			continue
		}
		var (
			enumName  = con.Name()
			enumValue = con.Val().ExactString()
			enumKind  = con.Val().Kind()
		)
		if con.Val().Kind() == constant.String {
			enumValue = constant.StringVal(con.Val())
		}
		p.enums = append(p.enums, EnumItem{
			Name:  enumName,
			Value: enumValue,
			Type:  enumType,
			Kind:  enumKind,
		})

	}
	for _, im := range pkg.Imports {
		p.ParsePackage(im)
	}
}

func (p *EnumsParser) Export() string {
	var typeEnumMap = make(map[string][]interface{})
	for _, enum := range p.enums {
		if typeEnumMap[enum.Type] == nil {
			typeEnumMap[enum.Type] = make([]interface{}, 0)
		}
		var value interface{}
		switch enum.Kind {
		case constant.Int:
			value = gconv.Int64(enum.Value)
		case constant.String:
			value = enum.Value
		case constant.Float:
			value = gconv.Float64(enum.Value)
		case constant.Bool:
			value = gconv.Bool(enum.Value)
		default:
			value = enum.Value
		}
		typeEnumMap[enum.Type] = append(typeEnumMap[enum.Type], value)
	}
	return gjson.MustEncodeString(typeEnumMap)
}
