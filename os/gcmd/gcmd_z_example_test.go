// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcmd_test

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/genv"
)

func ExampleInit() {
	gcmd.Init("gf", "build", "main.go", "-o=gf.exe", "-y")
	fmt.Printf(`%#v`, gcmd.GetArgAll())

	// Output:
	// []string{"gf", "build", "main.go"}
}

func ExampleGetArg() {
	gcmd.Init("gf", "build", "main.go", "-o=gf.exe", "-y")
	fmt.Printf(
		`Arg[0]: "%v", Arg[1]: "%v", Arg[2]: "%v", Arg[3]: "%v"`,
		gcmd.GetArg(0), gcmd.GetArg(1), gcmd.GetArg(2), gcmd.GetArg(3),
	)

	// Output:
	// Arg[0]: "gf", Arg[1]: "build", Arg[2]: "main.go", Arg[3]: ""
}

func ExampleGetArgAll() {
	gcmd.Init("gf", "build", "main.go", "-o=gf.exe", "-y")
	fmt.Printf(`%#v`, gcmd.GetArgAll())

	// Output:
	// []string{"gf", "build", "main.go"}
}

func ExampleGetOpt() {
	gcmd.Init("gf", "build", "main.go", "-o=gf.exe", "-y")
	fmt.Printf(
		`Opt["o"]: "%v", Opt["y"]: "%v", Opt["d"]: "%v"`,
		gcmd.GetOpt("o"), gcmd.GetOpt("y"), gcmd.GetOpt("d", "default value"),
	)

	// Output:
	// Opt["o"]: "gf.exe", Opt["y"]: "", Opt["d"]: "default value"
}

func ExampleGetOptAll() {
	gcmd.Init("gf", "build", "main.go", "-o=gf.exe", "-y")
	fmt.Printf(`%#v`, gcmd.GetOptAll())

	// May Output:
	// map[string]string{"o":"gf.exe", "y":""}
}

func ExampleGetOptWithEnv() {
	fmt.Printf("Opt[gf.test]:%s\n", gcmd.GetOptWithEnv("gf.test"))
	_ = genv.Set("GF_TEST", "YES")
	fmt.Printf("Opt[gf.test]:%s\n", gcmd.GetOptWithEnv("gf.test"))

	// Output:
	// Opt[gf.test]:
	// Opt[gf.test]:YES
}

func ExampleParse() {
	os.Args = []string{"gf", "build", "main.go", "-o=gf.exe", "-y"}
	p, err := gcmd.Parse(g.MapStrBool{
		"o,output": true,
		"y,yes":    false,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(p.GetOpt("o"))
	fmt.Println(p.GetOpt("output"))
	fmt.Println(p.GetOpt("y") != nil)
	fmt.Println(p.GetOpt("yes") != nil)
	fmt.Println(p.GetOpt("none") != nil)

	// Output:
	// gf.exe
	// gf.exe
	// true
	// true
	// false
}
