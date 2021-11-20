// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/command"
	"github.com/gogf/gf/v2/text/gstr"
)

type Command struct {
	parent      *Command
	commands    []Command
	options     []Option
	level       int
	Name        string
	Usage       string
	Short       string
	Brief       string
	Description string
	Func        func(parser *Parser)
	HelpFunc    func(parser *Parser)
	Examples    string
	Additional  string
}

type Option struct {
	Name        string
	Short       string
	Brief       string
	Description string
	NeedValue   bool
}

func (c *Command) Print() {
	prefix := gstr.Repeat(" ", 4)
	buffer := bytes.NewBuffer(nil)
	// Usage.
	if c.Usage != "" || c.Name != "" {
		buffer.WriteString("USAGE\n")
		buffer.WriteString(prefix)
		if c.Usage != "" {
			buffer.WriteString(c.Usage)
		} else {
			var (
				p    = c
				name = c.Name
			)
			for p.parent != nil {
				name = p.parent.Name + " " + name
				p = p.parent
			}
			buffer.WriteString(fmt.Sprintf(`%s ARGUMENT [OPTION]`, name))
		}
		buffer.WriteString("\n\n")
	}
	// Command.
	if len(c.commands) > 0 {
		buffer.WriteString("COMMAND\n")
		maxSpaceLength := 0
		for _, cmd := range c.commands {
			nameStr := cmd.Name + "/" + cmd.Short
			if len(nameStr) > maxSpaceLength {
				maxSpaceLength = len(nameStr)
			}
		}
		for _, cmd := range c.commands {
			nameStr := cmd.Name
			if cmd.Short != "" {
				nameStr += "/" + cmd.Short
			}
			var (
				spaceLength = maxSpaceLength - len(nameStr)
				lineStr     = fmt.Sprintf(
					"%s%s%s    %s\n",
					prefix, nameStr, gstr.Repeat(" ", spaceLength), cmd.Brief,
				)
			)
			lineStr = gstr.WordWrap(lineStr, maxLineChars, "\n")
			buffer.WriteString(lineStr)
		}
		buffer.WriteString("\n")
	}

	// Examples.
	if c.Examples != "" {
		buffer.WriteString("EXAMPLES\n")
		lineStr := gstr.WordWrap(gstr.Trim(c.Examples), maxLineChars, "\n")
		for _, line := range gstr.SplitAndTrim(lineStr, "\n") {
			buffer.WriteString(prefix)
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
		buffer.WriteString("\n")
	}
	// Description.
	if c.Description != "" {
		buffer.WriteString("DESCRIPTION\n")
		lineStr := gstr.WordWrap(gstr.Trim(c.Description), maxLineChars, "\n")
		for _, line := range gstr.SplitAndTrim(lineStr, "\n") {
			buffer.WriteString(prefix)
			buffer.WriteString(line)
			buffer.WriteString("\n")
		}
	}
	buffer.WriteString("\n")
	// Additional.
	if c.Additional != "" {
		lineStr := gstr.WordWrap(gstr.Trim(c.Additional), maxLineChars, "\n")
		buffer.WriteString(lineStr)
	}
	buffer.WriteString("\n")
	fmt.Println(buffer.String())
}

func (c *Command) AddCommand(command ...Command) error {
	for _, cmd := range command {
		cmd.Name = gstr.Trim(cmd.Name)
		if cmd.Name == "" {
			return gerror.New("command name should not be empty")
		}
		if cmd.Func == nil {
			return gerror.New("command function should not be empty")
		}
		cmd.parent = c
		cmd.level = c.level + 1
		c.commands = append(c.commands, cmd)
	}
	return nil
}

func (c *Command) AddOption(option ...Option) error {
	for _, opt := range option {
		opt.Name = gstr.Trim(opt.Name)
		if opt.Name == "" {
			return gerror.New("option name should not be empty")
		}
	}
	c.options = append(c.options, option...)
	return nil
}

func (c *Command) Run() {
	// Find the matched command and run it.
	argument := GetArg(c.level + 1)
	if !argument.IsEmpty() {
		if len(c.commands) > 0 {
			for _, cmd := range c.commands {
				if gstr.Equal(cmd.Name, argument.String()) {
					cmd.Run()
					return
				}
			}
		}
	}
	// Run current command function.
	var (
		err    error
		parser *Parser
	)
	if len(c.options) > 0 {
		optionParsingMap := make(map[string]bool, 0)
		// Add custom options to parser.
		for _, option := range c.options {
			optionParsingKey := option.Name
			if option.Short != "" {
				optionParsingKey += "," + option.Short
			}
			optionParsingMap[optionParsingKey] = option.NeedValue
		}
		// Add help option to parser.
		optionParsingMap[helpOptionName+","+helpOptionNameShort] = false
		parser, err = Parse(optionParsingMap)
	} else {
		parsedArgs, parsedOptions := command.ParseUsingDefaultAlgorithm(os.Args...)
		parser = &Parser{
			strict:        false,
			parsedArgs:    parsedArgs,
			parsedOptions: parsedOptions,
		}
	}
	if err != nil {
		fmt.Println("Error:", err)
	}
	if parser.ContainsOpt(helpOptionName) || parser.ContainsOpt(helpOptionNameShort) {
		if c.HelpFunc != nil {
			c.HelpFunc(parser)
		} else {
			c.Print()
		}
		return
	}
	c.Func(parser)
}
