// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcfg

import (
	"os"

	"github.com/gogf/gf/container/garray"

	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// ExpandValueEnv returns value of convert with environment variable.
//
// Return environment variable if value start with "${" and end with "}".
// Return default value if environment variable is empty or not exist.
//
// It accept value formats "${env}" , "${env||}}" , "${env||defaultValue}" , "defaultvalue".
// Examples:
//	v1 := gcfg.ExpandValueEnv("${GOROOT}")			    // return the GOROOT environment variable.
//	v2 := gcfg.ExpandValueEnv("${GOGF||/usr/local/go}")	// return the default value "/usr/local/go/".
//	v3 := gcfg.ExpandValueEnv("gogf")				    // return the value "gogf".
func ExpandValueEnv(value string) (realValue string) {
	realValue = value

	vLen := len(value)
	// Need"${}" string length gt 3
	if vLen < 3 {
		return
	}
	// Need start with "${" and end with "}", then return.
	if value[0] != '$' || value[1] != '{' || value[vLen-1] != '}' {
		return
	}

	key := ""
	defaultV := ""
	// value start with "${"
	for i := 2; i < vLen; i++ {
		if value[i] == '|' && (i+1 < vLen && value[i+1] == '|') {
			key = value[2:i]
			defaultV = value[i+2 : vLen-1] // other string is default value.
			break
		} else if value[i] == '}' {
			key = value[2:i]
			break
		}
	}

	realValue = os.Getenv(key)
	if realValue == "" {
		realValue = defaultV
	}

	return
}

// ExpandValueEnvForStr convert all string value with environment variable.
func ExpandValueEnvForStr(c string) string {
	// match example "${GOGF||/usr/local/go}"
	patten := `\"\$\{.*?\}\"`
	envs, err := gregex.MatchAll(patten, []byte(c))
	if err != nil {
		return c
	}
	if len(envs) > 0 {
		envMap := gmap.New(true)
		for _, env := range envs {
			if len(env) > 0 {
				envStr := string(env[0])
				realValue := ExpandValueEnv(gstr.Trim(envStr, "\""))
				realValue = EnvStrParse(realValue)
				envMap.Set(envStr, realValue)
			}
		}
		if envMap.Size() > 0 {
			envMap.Iterator(func(k, v interface{}) bool {
				c = gstr.Replace(c, k.(string), v.(string), -1)
				return true
			})
		}
	}
	return c
}

// EnvStrParse returns the toml value represented by the string.
//
// Bool
// true, TRUE, True return string true; false,FALSE, False return string false.
//
// Int、float、array
// return source string.
//
// Any other value return double quoted string.
func EnvStrParse(value string) string {
	if len(value) > 0 {
		// Double quotes are considered to be strings
		patterDQ := `^\".*?\"$`
		if gregex.IsMatch(patterDQ, []byte(value)) {
			return value
		}

		// parse bool
		trueArr := garray.NewStrArrayFrom([]string{"true", "TRUE", "True"})
		if trueArr.Contains(value) {
			return "true"
		}
		falseArr := garray.NewStrArrayFrom([]string{"false", "FALSE", "False"})
		if falseArr.Contains(value) {
			return "false"
		}

		// parse int or float
		patterIF := `^(\d+)$`
		intFloatStr := gstr.Replace(gstr.TrimLeft(value, "-"), ".", "", 1)
		if gregex.IsMatch(patterIF, []byte(intFloatStr)) {
			return value
		}

		// parse array
		// []
		pattenArr := `^\[.*?\]$`
		if gregex.IsMatch(pattenArr, []byte(value)) {
			return value
		}

		return "\"" + gconv.String(value) + "\""
	}
	return value
}
