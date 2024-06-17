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
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/gogf/gf/v2"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/genv"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// Run calls custom function in os.Args that bound to this command.
// It exits this process with exit code 1 if any error occurs.
func (c *Command) Run(ctx context.Context) {
	_ = c.RunWithValue(ctx)
}

// RunWithValue calls custom function in os.Args that bound to this command with value output.
// It exits this process with exit code 1 if any error occurs.
func (c *Command) RunWithValue(ctx context.Context) (value interface{}) {
	value, err := c.RunWithValueError(ctx)
	if err != nil {
		var (
			code   = gerror.Code(err)
			detail = code.Detail()
			buffer = bytes.NewBuffer(nil)
		)
		if code.Code() == gcode.CodeNotFound.Code() {
			buffer.WriteString(fmt.Sprintf("ERROR: %s\n", gstr.Trim(err.Error())))
			if lastCmd, ok := detail.(*Command); ok {
				lastCmd.PrintTo(buffer)
			} else {
				c.PrintTo(buffer)
			}
		} else {
			buffer.WriteString(fmt.Sprintf("%+v\n", err))
		}
		if gtrace.GetTraceID(ctx) == "" {
			fmt.Println(buffer.String())
			os.Exit(1)
		}
		glog.Stack(false).Fatal(ctx, buffer.String())
	}
	return value
}

// RunWithError calls custom function in os.Args that bound to this command with error output.
func (c *Command) RunWithError(ctx context.Context) (err error) {
	_, err = c.RunWithValueError(ctx)
	return
}

// RunWithValueError calls custom function in os.Args that bound to this command with value and error output.
func (c *Command) RunWithValueError(ctx context.Context) (value interface{}, err error) {
	return c.RunWithSpecificArgs(ctx, os.Args)
}

// RunWithSpecificArgs calls custom function in specific args that bound to this command with value and error output.
func (c *Command) RunWithSpecificArgs(ctx context.Context, args []string) (value interface{}, err error) {
	if len(args) == 0 {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "args can not be empty!")
	}
	parser, err := ParseArgs(args, nil)
	if err != nil {
		return nil, err
	}
	parsedArgs := parser.GetArgAll()
	if len(parsedArgs) == 1 {
		return c.doRun(ctx, args, parser)
	}

	// Exclude the root binary name.
	parsedArgs = parsedArgs[1:]

	// Find the matched command and run it.
	lastCmd, foundCmd, newCtx := c.searchCommand(ctx, parsedArgs, 0)
	if foundCmd != nil {
		return foundCmd.doRun(newCtx, args, parser)
	}

	// Print error and help command if no command found.
	err = gerror.NewCodef(
		gcode.WithCode(gcode.CodeNotFound, lastCmd),
		`command "%s" not found for command "%s", command line: %s`,
		gstr.Join(parsedArgs, " "),
		c.Name,
		gstr.Join(args, " "),
	)
	return
}

func (c *Command) doRun(ctx context.Context, args []string, parser *Parser) (value interface{}, err error) {
	defer func() {
		if exception := recover(); exception != nil {
			if v, ok := exception.(error); ok && gerror.HasStack(v) {
				err = v
			} else {
				err = gerror.NewCodef(gcode.CodeInternalPanic, "exception recovered: %+v", exception)
			}
		}
	}()

	ctx = context.WithValue(ctx, CtxKeyCommand, c)
	// Check built-in help command.
	if parser.GetOpt(helpOptionName) != nil || parser.GetOpt(helpOptionNameShort) != nil {
		if c.HelpFunc != nil {
			return nil, c.HelpFunc(ctx, parser)
		}
		return nil, c.defaultHelpFunc(ctx, parser)
	}
	// OpenTelemetry for command.
	var (
		span trace.Span
		tr   = otel.GetTracerProvider().Tracer(
			tracingInstrumentName,
			trace.WithInstrumentationVersion(gf.VERSION),
		)
	)
	ctx, span = tr.Start(
		otel.GetTextMapPropagator().Extract(
			ctx,
			propagation.MapCarrier(genv.Map()),
		),
		gstr.Join(os.Args, " "),
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()
	span.SetAttributes(gtrace.CommonLabels()...)
	// Reparse the original arguments for current command configuration.
	parser, err = c.reParse(ctx, args, parser)
	if err != nil {
		return nil, err
	}
	// Registered command function calling.
	if c.Func != nil {
		return nil, c.Func(ctx, parser)
	}
	if c.FuncWithValue != nil {
		return c.FuncWithValue(ctx, parser)
	}
	// If no function defined in current command, it then prints help.
	if c.HelpFunc != nil {
		return nil, c.HelpFunc(ctx, parser)
	}
	return nil, c.defaultHelpFunc(ctx, parser)
}

// reParse parses the original arguments using option configuration of current command.
func (c *Command) reParse(ctx context.Context, args []string, parser *Parser) (*Parser, error) {
	if len(c.Arguments) == 0 {
		return parser, nil
	}
	var (
		optionKey        string
		supportedOptions = make(map[string]bool)
	)
	for _, arg := range c.Arguments {
		if arg.IsArg {
			continue
		}
		if arg.Short != "" {
			optionKey = fmt.Sprintf(`%s,%s`, arg.Name, arg.Short)
		} else {
			optionKey = arg.Name
		}
		supportedOptions[optionKey] = !arg.Orphan
	}
	parser, err := ParseArgs(args, supportedOptions, ParserOption{
		CaseSensitive: c.CaseSensitive,
		Strict:        c.Strict,
	})
	if err != nil {
		return nil, err
	}
	// Retrieve option values from config component if it has "config" tag.
	if c.Config != "" && gcfg.Instance().Available(ctx) {
		value, err := gcfg.Instance().Get(ctx, c.Config)
		if err != nil {
			return nil, err
		}
		configMap := value.Map()
		for optionName := range parser.supportedOptions {
			// The command line has the high priority.
			if parser.GetOpt(optionName) != nil {
				continue
			}
			// Merge the config value into parser.
			foundKey, foundValue := gutil.MapPossibleItemByKey(configMap, optionName)
			if foundKey != "" {
				parser.parsedOptions[optionName] = gconv.String(foundValue)
			}
		}
	}
	return parser, nil
}

// searchCommand recursively searches the command according given arguments.
func (c *Command) searchCommand(
	ctx context.Context, args []string, fromArgIndex int,
) (lastCmd, foundCmd *Command, newCtx context.Context) {
	if len(args) == 0 {
		return c, nil, ctx
	}
	for _, cmd := range c.commands {
		// Recursively searching the command.
		// String comparison case-sensitive.
		if cmd.Name == args[0] {
			leftArgs := args[1:]
			// If this command needs argument,
			// it then gives all its left arguments to it using arg index marks.
			//
			// Note that the args here (using default args parsing) could be different with the args
			// that are parsed in command.
			if cmd.hasArgumentFromIndex() || len(leftArgs) == 0 {
				ctx = context.WithValue(ctx, CtxKeyArgumentsIndex, fromArgIndex+1)
				return c, cmd, ctx
			}
			return cmd.searchCommand(ctx, leftArgs, fromArgIndex+1)
		}
	}
	return c, nil, ctx
}

func (c *Command) hasArgumentFromIndex() bool {
	for _, arg := range c.Arguments {
		if arg.IsArg {
			return true
		}
	}
	return false
}

func (c *Command) hasArgumentFromOption() bool {
	for _, arg := range c.Arguments {
		if !arg.IsArg {
			return true
		}
	}
	return false
}
