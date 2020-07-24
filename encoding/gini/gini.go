// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gini provides accessing and converting for INI content.
package gini

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"github.com/gogf/gf/internal/json"
	"io"
	"strings"
)

// Decode converts INI format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	fieldMap := make(map[string]interface{})

	a := bytes.NewReader(data)
	r := bufio.NewReader(a)
	var section string
	var lastSection string
	var haveSection bool
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		lineStr := strings.TrimSpace(string(line))
		if len(lineStr) == 0 {
			continue
		}

		if lineStr[0] == ';' || lineStr[0] == '#' {
			continue
		}

		sectionBeginPos := strings.Index(lineStr, "[")
		sectionEndPos := strings.Index(lineStr, "]")

		if sectionBeginPos >= 0 && sectionEndPos >= 2 {
			section = lineStr[sectionBeginPos+1 : sectionEndPos]

			if lastSection == "" {
				lastSection = section
			} else if lastSection != section {
				lastSection = section
				fieldMap = make(map[string]interface{})
			}
			haveSection = true
		} else if haveSection == false {
			continue
		}

		if strings.Contains(lineStr, "=") && haveSection {
			values := strings.Split(lineStr, "=")
			fieldMap[strings.TrimSpace(values[0])] = strings.TrimSpace(strings.Join(values[1:], ""))
			res[section] = fieldMap
		}
	}

	if haveSection == false {
		return nil, errors.New("failed to parse INI file, section not found")
	}
	return res, nil
}

// Encode converts map to INI format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	w := new(bytes.Buffer)

	w.WriteString("; this ini file is produced by package gini\n")
	for k, v := range data {
		n, err := w.WriteString(fmt.Sprintf("[%s]\n", k))
		if err != nil || n == 0 {
			return nil, fmt.Errorf("write data failed. %v", err)
		}
		for kk, vv := range v.(map[string]interface{}) {
			n, err := w.WriteString(fmt.Sprintf("%s=%s\n", kk, vv.(string)))
			if err != nil || n == 0 {
				return nil, fmt.Errorf("write data failed. %v", err)
			}
		}
	}
	res = make([]byte, w.Len())
	n, err := w.Read(res)
	if err != nil || n == 0 {
		return nil, fmt.Errorf("write data failed. %v", err)
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
