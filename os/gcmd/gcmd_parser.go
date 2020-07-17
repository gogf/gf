// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"fmt"
	"github.com/gogf/gf/internal/json"
	"os"
	"strings"

	"github.com/gogf/gf/text/gstr"

	"errors"

	"github.com/gogf/gf/container/gvar"

	"github.com/gogf/gf/text/gregex"
)

// Parser for arguments.
type Parser struct {
	strict           bool              // Whether stops parsing and returns error if invalid option passed.
	parsedArgs       []string          // As name described.
	parsedOptions    map[string]string // As name described.
	passedOptions    map[string]bool   // User passed supported options.
	supportedOptions map[string]bool   // Option [option name : need argument].
	commandFuncMap   map[string]func() // Command function map for function handler.
}

// Parse creates and returns a new Parser with os.Args and supported options.
//
// Note that the parameter <supportedOptions> is as [option name: need argument], which means
// the value item of <supportedOptions> indicates whether corresponding option name needs argument or not.
//
// The optional parameter <strict> specifies whether stops parsing and returns error if invalid option passed.
func Parse(supportedOptions map[string]bool, strict ...bool) (*Parser, error) {
	return ParseWithArgs(os.Args, supportedOptions, strict...)
}

// ParseWithArgs creates and returns a new Parser with given arguments and supported options.
//
// Note that the parameter <supportedOptions> is as [option name: need argument], which means
// the value item of <supportedOptions> indicates whether corresponding option name needs argument or not.
//
// The optional parameter <strict> specifies whether stops parsing and returns error if invalid option passed.
func ParseWithArgs(args []string, supportedOptions map[string]bool, strict ...bool) (*Parser, error) {
	strictParsing := false
	if len(strict) > 0 {
		strictParsing = strict[0]
	}
	parser := &Parser{
		strict:           strictParsing,
		parsedArgs:       make([]string, 0),
		parsedOptions:    make(map[string]string),
		passedOptions:    supportedOptions,
		supportedOptions: make(map[string]bool),
		commandFuncMap:   make(map[string]func()),
	}
	for name, needArgument := range supportedOptions {
		for _, v := range strings.Split(name, ",") {
			parser.supportedOptions[strings.TrimSpace(v)] = needArgument
		}
	}

	for i := 0; i < len(args); {
		if option := parser.parseOption(args[i]); option != "" {
			array, _ := gregex.MatchString(`^(.+?)=(.+)$`, option)
			if len(array) == 3 {
				if parser.isOptionValid(array[1]) {
					parser.setOptionValue(array[1], array[2])
				}
			} else {
				if parser.isOptionValid(option) {
					if parser.isOptionNeedArgument(option) {
						if i < len(args)-1 {
							parser.setOptionValue(option, args[i+1])
							i += 2
							continue
						}
					} else {
						parser.setOptionValue(option, "")
						i++
						continue
					}
				} else {
					// Multiple options?
					if array := parser.parseMultiOption(option); len(array) > 0 {
						for _, v := range array {
							parser.setOptionValue(v, "")
						}
						i++
						continue
					} else if parser.strict {
						return nil, errors.New(fmt.Sprintf(`invalid option '%s'`, args[i]))
					}
				}
			}
		} else {
			parser.parsedArgs = append(parser.parsedArgs, args[i])
		}
		i++
	}
	return parser, nil
}

// parseMultiOption parses option to multiple valid options like: --dav.
// It returns nil if given option is not multi-option.
func (p *Parser) parseMultiOption(option string) []string {
	for i := 1; i <= len(option); i++ {
		s := option[:i]
		if p.isOptionValid(s) && !p.isOptionNeedArgument(s) {
			if i == len(option) {
				return []string{s}
			}
			array := p.parseMultiOption(option[i:])
			if len(array) == 0 {
				return nil
			}
			return append(array, s)
		}
	}
	return nil
}

func (p *Parser) parseOption(argument string) string {
	array, _ := gregex.MatchString(`^\-{1,2}(.+)$`, argument)
	if len(array) == 2 {
		return array[1]
	}
	return ""
}

func (p *Parser) isOptionValid(name string) bool {
	_, ok := p.supportedOptions[name]
	return ok
}

func (p *Parser) isOptionNeedArgument(name string) bool {
	return p.supportedOptions[name]
}

// setOptionValue sets the option value for name and according alias.
func (p *Parser) setOptionValue(name, value string) {
	for optionName, _ := range p.passedOptions {
		array := gstr.SplitAndTrim(optionName, ",")
		for _, v := range array {
			if strings.EqualFold(v, name) {
				for _, v := range array {
					p.parsedOptions[v] = value
				}
				return
			}
		}
	}
}

// GetOpt returns the option value named <name>.
func (p *Parser) GetOpt(name string, def ...string) string {
	if v, ok := p.parsedOptions[name]; ok {
		return v
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetOptVar returns the option value named <name> as gvar.Var.
func (p *Parser) GetOptVar(name string, def ...interface{}) *gvar.Var {
	if p.ContainsOpt(name) {
		return gvar.New(p.GetOpt(name))
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return gvar.New(nil)
}

// GetOptAll returns all parsed options.
func (p *Parser) GetOptAll() map[string]string {
	return p.parsedOptions
}

// ContainsOpt checks whether option named <name> exist in the arguments.
func (p *Parser) ContainsOpt(name string) bool {
	_, ok := p.parsedOptions[name]
	return ok
}

// GetArg returns the argument at <index>.
func (p *Parser) GetArg(index int, def ...string) string {
	if index < len(p.parsedArgs) {
		return p.parsedArgs[index]
	}
	if len(def) > 0 {
		return def[0]
	}
	return ""
}

// GetArgVar returns the argument at <index> as gvar.Var.
func (p *Parser) GetArgVar(index int, def ...string) *gvar.Var {
	return gvar.New(p.GetArg(index, def...))
}

// GetArgAll returns all parsed arguments.
func (p *Parser) GetArgAll() []string {
	return p.parsedArgs
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (p *Parser) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"parsedArgs":       p.parsedArgs,
		"parsedOptions":    p.parsedOptions,
		"passedOptions":    p.passedOptions,
		"supportedOptions": p.supportedOptions,
	})
}
