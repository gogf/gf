// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// Command holds the info about an argument that can handle custom logic.
type Command struct {
	Name        string    // Command name(case-sensitive).
	Usage       string    // A brief line description about its usage, eg: gf build main.go [OPTION]
	Brief       string    // A brief info that describes what this command will do.
	Description string    // A detailed description.
	Options     []Option  // Option array, configuring how this command act.
	Func        Function  // Custom function.
	HelpFunc    Function  // Custom help function
	Examples    string    // Usage examples.
	Additional  string    // Additional custom info about this command.
	parent      *Command  // Parent command for internal usage.
	commands    []Command // Sub commands of this command.
}

// Function is a custom command callback function that is bound to a certain argument.
type Function func(ctx context.Context, parser *Parser) (err error)

// Option is the command value that is specified by a name or shor name.
// An Option can have or have no value bound to it.
type Option struct {
	Name      string // Option name.
	Short     string // Option short.
	Brief     string // Brief info about this Option, which is used in help info.
	NeedValue bool   // Whether this Option having or having no value bound to it.
}

var (
	// defaultHelpOption is the default help option that will be automatically added to each command.
	defaultHelpOption = Option{
		Name:      `help`,
		Short:     `h`,
		Brief:     `more information about this command`,
		NeedValue: false,
	}
)

// Add adds one or more sub-commands to current command.
func (c *Command) Add(commands ...Command) error {
	for _, cmd := range commands {
		cmd.Name = gstr.Trim(cmd.Name)
		if cmd.Name == "" {
			return gerror.New("command name should not be empty")
		}
		if cmd.Func == nil {
			return gerror.New("command function should not be empty")
		}
		cmd.parent = c
		c.commands = append(c.commands, cmd)
	}
	return nil
}
