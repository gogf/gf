// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

// Run calls custom function that bound to this command.
func (c *Command) Run(ctx context.Context) error {
	// Parse command arguments and options using default algorithm.
	parser, err := Parse(nil)
	if err != nil {
		return err
	}
	args := parser.GetArgAll()
	if len(args) == 1 {
		if c.HelpFunc != nil {
			return c.HelpFunc(ctx, parser)
		}
		return c.defaultHelpFunc(ctx, parser)
	}

	// Exclude the root binary name.
	args = args[1:]

	// Find the matched command and run it.
	if subCommand := c.searchCommand(args); subCommand != nil {
		return subCommand.doRun(ctx, parser)
	}

	// Print error and help command if no command found.
	fmt.Printf("Error: command not found for \"%s\"\n\n", gstr.Join(args, " "))
	c.Print()

	return nil
}

func (c *Command) doRun(ctx context.Context, parser *Parser) (err error) {
	// Add built-in help option, just for info only.
	c.Options = append(c.Options, defaultHelpOption)
	// Check built-in help command.
	if parser.ContainsOpt(helpOptionName) || parser.ContainsOpt(helpOptionNameShort) {
		if c.HelpFunc != nil {
			return c.HelpFunc(ctx, parser)
		}
		return c.defaultHelpFunc(ctx, parser)
	}
	// Reparse the arguments for current command configuration.
	parser, err = c.reParse(ctx, parser)
	if err != nil {
		return err
	}
	// Registered command function calling.
	return c.Func(ctx, parser)
}

func (c *Command) reParse(ctx context.Context, parser *Parser) (*Parser, error) {
	// It seems just has built-in help option, it so does nothing.
	if len(c.Options) == 1 {
		return parser, nil
	}
	var (
		optionKey        string
		supportedOptions = make(map[string]bool)
	)
	for _, option := range c.Options {
		if option.Short != "" {
			optionKey = fmt.Sprintf(`%s.%s`, option.Name, option.Short)
		} else {
			optionKey = option.Name
		}
		supportedOptions[optionKey] = option.NeedValue
	}
	return Parse(supportedOptions)
}

func (c *Command) searchCommand(args []string) *Command {
	if len(args) == 0 {
		return nil
	}
	for _, cmd := range c.commands {
		if cmd.Name == args[0] {
			leftArgs := args[1:]
			if len(leftArgs) == 0 {
				return &cmd
			}
			return cmd.searchCommand(leftArgs)
		}
	}
	return nil
}
