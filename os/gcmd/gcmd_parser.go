// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"
	"os"
	"strings"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/text/gregex"
	"github.com/gogf/gf/v2/text/gstr"
)

// ParserOption manages the parsing options.
type ParserOption struct {
	CaseSensitive bool // Marks options parsing in case-sensitive way.
	Strict        bool // Whether stops parsing and returns error if invalid option passed.
}

// Parser for arguments.
type Parser struct {
	option           ParserOption      // Parse option.
	parsedArgs       []string          // As name described.
	parsedOptions    map[string]string // As name described.
	passedOptions    map[string]bool   // User passed supported options, like: map[string]bool{"name,n":true}
	supportedOptions map[string]bool   // Option [OptionName:WhetherNeedArgument], like: map[string]bool{"name":true, "n":true}
	commandFuncMap   map[string]func() // Command function map for function handler.
}

// ParserFromCtx retrieves and returns Parser from context.
func ParserFromCtx(ctx context.Context) *Parser {
	if v := ctx.Value(CtxKeyParser); v != nil {
		if p, ok := v.(*Parser); ok {
			return p
		}
	}
	return nil
}

// Parse creates and returns a new Parser with os.Args and supported options.
//
// Note that the parameter `supportedOptions` is as [option name: need argument], which means
// the value item of `supportedOptions` indicates whether corresponding option name needs argument or not.
//
// The optional parameter `strict` specifies whether stops parsing and returns error if invalid option passed.
func Parse(supportedOptions map[string]bool, option ...ParserOption) (*Parser, error) {
	if supportedOptions == nil {
		command.Init(os.Args...)
		return &Parser{
			parsedArgs:    GetArgAll(),
			parsedOptions: GetOptAll(),
		}, nil
	}
	return ParseArgs(os.Args, supportedOptions, option...)
}

// ParseArgs creates and returns a new Parser with given arguments and supported options.
//
// Note that the parameter `supportedOptions` is as [option name: need argument], which means
// the value item of `supportedOptions` indicates whether corresponding option name needs argument or not.
//
// The optional parameter `strict` specifies whether stops parsing and returns error if invalid option passed.
func ParseArgs(args []string, supportedOptions map[string]bool, option ...ParserOption) (*Parser, error) {
	if supportedOptions == nil {
		command.Init(args...)
		return &Parser{
			parsedArgs:    GetArgAll(),
			parsedOptions: GetOptAll(),
		}, nil
	}
	var parserOption ParserOption
	if len(option) > 0 {
		parserOption = option[0]
	}
	parser := &Parser{
		option:           parserOption,
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
					if array = parser.parseMultiOption(option); len(array) > 0 {
						for _, v := range array {
							parser.setOptionValue(v, "")
						}
						i++
						continue
					} else if parser.option.Strict {
						return nil, gerror.NewCodef(gcode.CodeInvalidParameter, `invalid option '%s'`, args[i])
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
	// Case-Sensitive.
	if p.option.CaseSensitive {
		_, ok := p.supportedOptions[name]
		return ok
	}
	// Case-InSensitive.
	for optionName := range p.supportedOptions {
		if gstr.Equal(optionName, name) {
			return true
		}
	}
	return false
}

func (p *Parser) isOptionNeedArgument(name string) bool {
	return p.supportedOptions[name]
}

// setOptionValue sets the option value for name and according alias.
func (p *Parser) setOptionValue(name, value string) {
	// Accurate option name match.
	for optionName := range p.passedOptions {
		optionNameAndShort := gstr.SplitAndTrim(optionName, ",")
		for _, optionNameItem := range optionNameAndShort {
			if optionNameItem == name {
				for _, v := range optionNameAndShort {
					p.parsedOptions[v] = value
				}
				return
			}
		}
	}
	// Fuzzy option name match.
	for optionName := range p.passedOptions {
		optionNameAndShort := gstr.SplitAndTrim(optionName, ",")
		for _, optionNameItem := range optionNameAndShort {
			if strings.EqualFold(optionNameItem, name) {
				for _, v := range optionNameAndShort {
					p.parsedOptions[v] = value
				}
				return
			}
		}
	}
}

// GetOpt returns the option value named `name` as gvar.Var.
func (p *Parser) GetOpt(name string, def ...interface{}) *gvar.Var {
	if p == nil {
		return nil
	}
	if v, ok := p.parsedOptions[name]; ok {
		return gvar.New(v)
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetOptAll returns all parsed options.
func (p *Parser) GetOptAll() map[string]string {
	if p == nil {
		return nil
	}
	return p.parsedOptions
}

// GetArg returns the argument at `index` as gvar.Var.
func (p *Parser) GetArg(index int, def ...string) *gvar.Var {
	if p == nil {
		return nil
	}
	if index >= 0 && index < len(p.parsedArgs) {
		return gvar.New(p.parsedArgs[index])
	}
	if len(def) > 0 {
		return gvar.New(def[0])
	}
	return nil
}

// GetArgAll returns all parsed arguments.
func (p *Parser) GetArgAll() []string {
	if p == nil {
		return nil
	}
	return p.parsedArgs
}

// MarshalJSON implements the interface MarshalJSON for json.Marshal.
func (p Parser) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"parsedArgs":       p.parsedArgs,
		"parsedOptions":    p.parsedOptions,
		"passedOptions":    p.passedOptions,
		"supportedOptions": p.supportedOptions,
	})
}
