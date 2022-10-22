// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gini provides accessing and converting for INI content.
package gini

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
)

// Decode converts INI format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	var (
		fieldMap    = make(map[string]interface{})
		bytesReader = bytes.NewReader(data)
		bufioReader = bufio.NewReader(bytesReader)
		section     string
		lastSection string
		haveSection bool
		line        string
	)

	for {
		line, err = bufioReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			err = gerror.Wrapf(err, `bufioReader.ReadString failed`)
			return nil, err
		}
		if line = strings.TrimSpace(line); len(line) == 0 {
			continue
		}

		if line[0] == ';' || line[0] == '#' {
			continue
		}
		var (
			sectionBeginPos = strings.Index(line, "[")
			sectionEndPos   = strings.Index(line, "]")
		)
		if sectionBeginPos >= 0 && sectionEndPos >= 2 {
			section = line[sectionBeginPos+1 : sectionEndPos]
			if lastSection == "" {
				lastSection = section
			} else if lastSection != section {
				lastSection = section
				fieldMap = make(map[string]interface{})
			}
			haveSection = true
		} else if !haveSection {
			continue
		}

		if strings.Contains(line, "=") && haveSection {
			values := strings.Split(line, "=")
			fieldMap[strings.TrimSpace(values[0])] = strings.TrimSpace(strings.Join(values[1:], "="))
			res[section] = fieldMap
		}
	}

	if !haveSection {
		return nil, gerror.NewCode(gcode.CodeInvalidParameter, "failed to parse INI file, section not found")
	}
	return res, nil
}

// Encode converts map to INI format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	var (
		n  int
		w  = new(bytes.Buffer)
		m  map[string]interface{}
		ok bool
	)
	for section, item := range data {
		// Section key-value pairs.
		if m, ok = item.(map[string]interface{}); ok {
			n, err = w.WriteString(fmt.Sprintf("[%s]\n", section))
			if err != nil || n == 0 {
				return nil, gerror.Wrapf(err, "w.WriteString failed")
			}
			for k, v := range m {
				if n, err = w.WriteString(fmt.Sprintf("%s=%v\n", k, v)); err != nil || n == 0 {
					return nil, gerror.Wrapf(err, "w.WriteString failed")
				}
			}
			continue
		}
		// Simple key-value pairs.
		for k, v := range data {
			if n, err = w.WriteString(fmt.Sprintf("%s=%v\n", k, v)); err != nil || n == 0 {
				return nil, gerror.Wrapf(err, "w.WriteString failed")
			}
		}
		break
	}
	res = make([]byte, w.Len())
	if n, err = w.Read(res); err != nil || n == 0 {
		return nil, gerror.Wrapf(err, "w.Read failed")
	}
	return res, nil
}

// ToJson convert INI format to JSON.
func ToJson(data []byte) (res []byte, err error) {
	iniMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(iniMap)
}
