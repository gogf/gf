// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"bytes"
	"context"
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

// Print prints help info to stdout for current command.
func (c *Command) Print() {
	var (
		prefix  = gstr.Repeat(" ", 4)
		buffer  = bytes.NewBuffer(nil)
		options = make([]Option, len(c.Options))
	)
	// Copy options for printing.
	copy(options, c.Options)
	// Add built-in help option, just for info only.
	options = append(options, defaultHelpOption)

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
		var (
			maxSpaceLength = 0
		)
		for _, cmd := range c.commands {
			if len(cmd.Name) > maxSpaceLength {
				maxSpaceLength = len(cmd.Name)
			}
		}
		for _, cmd := range c.commands {
			// Add "..." to brief for those commands that also have sub-commands.
			if len(cmd.commands) > 0 {
				cmd.Brief = gstr.TrimRight(cmd.Brief, ".") + "..."
			}
			var (
				spaceLength    = maxSpaceLength - len(cmd.Name)
				lineStr        = fmt.Sprintf("%s%s%s%s\n", prefix, cmd.Name, gstr.Repeat(" ", spaceLength+4), cmd.Brief)
				wordwrapPrefix = gstr.Repeat(" ", len(prefix+cmd.Name)+spaceLength+4)
			)
			lineStr = gstr.WordWrap(lineStr, maxLineChars, "\n"+wordwrapPrefix)
			buffer.WriteString(lineStr)
		}
		buffer.WriteString("\n")
	}

	// Option.
	if len(options) > 0 {
		buffer.WriteString("OPTION\n")
		var (
			nameStr        string
			maxSpaceLength = 0
		)
		for _, option := range options {
			if option.Short != "" {
				nameStr = fmt.Sprintf("-%s,\t--%s", option.Short, option.Name)
			} else {
				nameStr = fmt.Sprintf("-/--%s", option.Name)
			}
			if len(nameStr) > maxSpaceLength {
				maxSpaceLength = len(nameStr)
			}
		}
		for _, option := range options {
			if option.Short != "" {
				nameStr = fmt.Sprintf("-%s,\t--%s", option.Short, option.Name)
			} else {
				nameStr = fmt.Sprintf("-/--%s", option.Name)
			}
			var (
				spaceLength    = maxSpaceLength - len(nameStr)
				lineStr        = fmt.Sprintf("%s%s%s%s\n", prefix, nameStr, gstr.Repeat(" ", spaceLength+4), option.Brief)
				wordwrapPrefix = gstr.Repeat(" ", len(prefix+nameStr)+spaceLength+4)
			)
			lineStr = gstr.WordWrap(lineStr, maxLineChars, "\n"+wordwrapPrefix)
			buffer.WriteString(lineStr)
		}
		buffer.WriteString("\n")
	}

	// Example.
	if c.Examples != "" {
		buffer.WriteString("EXAMPLE\n")
		for _, line := range gstr.SplitAndTrim(c.Examples, "\n") {
			buffer.WriteString(prefix)
			buffer.WriteString(gstr.WordWrap(gstr.Trim(line), maxLineChars, "\n"+prefix))
			buffer.WriteString("\n")
		}
		buffer.WriteString("\n")
	}

	// Description.
	if c.Description != "" {
		buffer.WriteString("DESCRIPTION\n")
		for _, line := range gstr.SplitAndTrim(c.Description, "\n") {
			buffer.WriteString(prefix)
			buffer.WriteString(gstr.WordWrap(gstr.Trim(line), maxLineChars, "\n"+prefix))
			buffer.WriteString("\n")
		}
		buffer.WriteString("\n")
	}

	// Additional.
	if c.Additional != "" {
		lineStr := gstr.WordWrap(gstr.Trim(c.Additional), maxLineChars, "\n")
		buffer.WriteString(lineStr)
		buffer.WriteString("\n")
	}

	fmt.Println(buffer.String())
}

func (c *Command) defaultHelpFunc(ctx context.Context, parser *Parser) error {
	// Print command help info to stdout.
	c.Print()
	return nil
}
