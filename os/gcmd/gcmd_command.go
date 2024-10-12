// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gcmd

import (
	"context"

	"github.com/gogf/gf/v2/container/gset"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/text/gstr"
)

// Command holds the info about an argument that can handle custom logic.
type Command struct {
	Name          string        // Command name(case-sensitive).
	Usage         string        // A brief line description about its usage, eg: gf build main.go [OPTION]
	Brief         string        // A brief info that describes what this command will do.
	Description   string        // A detailed description.
	Arguments     []Argument    // Argument array, configuring how this command act.
	Func          Function      // Custom function.
	FuncWithValue FuncWithValue // Custom function with output parameters that can interact with command caller.
	HelpFunc      Function      // Custom help function.
	Examples      string        // Usage examples.
	Additional    string        // Additional info about this command, which will be appended to the end of help info.
	Strict        bool          // Strict parsing options, which means it returns error if invalid option given.
	CaseSensitive bool          // CaseSensitive parsing options, which means it parses input options in case-sensitive way.
	Config        string        // Config node name, which also retrieves the values from config component along with command line.
	internalCommandAttributes
}

type internalCommandAttributes struct {
	parent   *Command   // Parent command for internal usage.
	commands []*Command // Sub commands of this command.
}

// Function is a custom command callback function that is bound to a certain argument.
type Function func(ctx context.Context, parser *Parser) (err error)

// FuncWithValue is similar like Func but with output parameters that can interact with command caller.
type FuncWithValue func(ctx context.Context, parser *Parser) (out interface{}, err error)

// Argument is the command value that are used by certain command.
type Argument struct {
	Name   string // Option name.
	Short  string // Option short.
	Brief  string // Brief info about this Option, which is used in help info.
	IsArg  bool   // IsArg marks this argument taking value from command line argument instead of option.
	Orphan bool   // Whether this Option having or having no value bound to it.
}

var (
	// defaultHelpOption is the default help option that will be automatically added to each command.
	defaultHelpOption = Argument{
		Name:   `help`,
		Short:  `h`,
		Brief:  `more information about this command`,
		Orphan: true,
	}
)

// CommandFromCtx retrieves and returns Command from context.
func CommandFromCtx(ctx context.Context) *Command {
	if v := ctx.Value(CtxKeyCommand); v != nil {
		if p, ok := v.(*Command); ok {
			return p
		}
	}
	return nil
}

// AddCommand adds one or more sub-commands to current command.
func (c *Command) AddCommand(commands ...*Command) error {
	for _, cmd := range commands {
		if err := c.doAddCommand(cmd); err != nil {
			return err
		}
	}
	return nil
}

// doAddCommand adds one sub-command to current command.
func (c *Command) doAddCommand(command *Command) error {
	command.Name = gstr.Trim(command.Name)
	if command.Name == "" {
		return gerror.New("command name should not be empty")
	}
	// Repeated check.
	var (
		commandNameSet = gset.NewStrSet()
	)
	for _, cmd := range c.commands {
		commandNameSet.Add(cmd.Name)
	}
	if commandNameSet.Contains(command.Name) {
		return gerror.Newf(`command "%s" is already added to command "%s"`, command.Name, c.Name)
	}
	// Add the given command to its sub-commands array.
	command.parent = c
	c.commands = append(c.commands, command)
	return nil
}

// AddObject adds one or more sub-commands to current command using struct object.
func (c *Command) AddObject(objects ...interface{}) error {
	var commands []*Command
	for _, object := range objects {
		rootCommand, err := NewFromObject(object)
		if err != nil {
			return err
		}
		commands = append(commands, rootCommand)
	}
	return c.AddCommand(commands...)
}
