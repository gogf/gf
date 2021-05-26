package main

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcmd"
)

func main() {
	var err error
	c := &gcmd.Command{
		Name:        "gf",
		Description: `GoFrame Command Line Interface, which is your helpmate for building GoFrame application with convenience.`,
		Additional: `
Use 'gf help COMMAND' or 'gf COMMAND -h' for detail about a command, which has '...' in the tail of their comments.`,
	}
	// env
	commandEnv := gcmd.Command{
		Name:        "env",
		Brief:       "show current Golang environment variables",
		Description: "show current Golang environment variables",
		Func: func(parser *gcmd.Parser) {

		},
	}
	if err = c.AddCommand(commandEnv); err != nil {
		g.Log().Fatal(err)
	}
	// get
	commandGet := gcmd.Command{
		Name:        "get",
		Brief:       "install or update GF to system in default...",
		Description: "show current Golang environment variables",

		Examples: `
gf get github.com/gogf/gf
gf get github.com/gogf/gf@latest
gf get github.com/gogf/gf@master
gf get golang.org/x/sys
`,
		Func: func(parser *gcmd.Parser) {

		},
	}
	if err = c.AddCommand(commandGet); err != nil {
		g.Log().Fatal(err)
	}
	// build
	//-n, --name       output binary name
	//-v, --version    output binary version
	//-a, --arch       output binary architecture, multiple arch separated with ','
	//-s, --system     output binary system, multiple os separated with ','
	//-o, --output     output binary path, used when building single binary file
	//-p, --path       output binary directory path, default is './bin'
	//-e, --extra      extra custom "go build" options
	//-m, --mod        like "-mod" option of "go build", use "-m none" to disable go module
	//-c, --cgo        enable or disable cgo feature, it's disabled in default

	commandBuild := gcmd.Command{
		Name:  "build",
		Usage: "gf build FILE [OPTION]",
		Brief: "cross-building go project for lots of platforms...",
		Description: `
The "build" command is most commonly used command, which is designed as a powerful wrapper for
"go build" command for convenience cross-compiling usage.
It provides much more features for building binary:
1. Cross-Compiling for many platforms and architectures.
2. Configuration file support for compiling.
3. Build-In Variables.
`,
		Examples: `
gf build main.go
gf build main.go --swagger
gf build main.go --pack public,template
gf build main.go --cgo
gf build main.go -m none 
gf build main.go -n my-app -a all -s all
gf build main.go -n my-app -a amd64,386 -s linux -p .
gf build main.go -n my-app -v 1.0 -a amd64,386 -s linux,windows,darwin -p ./docker/bin
`,
		Func: func(parser *gcmd.Parser) {

		},
	}
	if err = c.AddCommand(commandBuild); err != nil {
		g.Log().Fatal(err)
	}
	c.Run()
}
