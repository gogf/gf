// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gproperties provides accessing and converting for .properties content.

package gproperties

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
)

func genMap(m map[string]interface{}, keys []string, value interface{}) (err error) {
	if m == nil {
		m = make(map[string]interface{})
	}
	if len(keys) == 1 {
		m[keys[0]] = value
	} else {
		if m[keys[0]] == nil {
			m[keys[0]] = make(map[string]interface{})
		}
		if _, ok := m[keys[0]].(map[string]interface{}); ok == false {
			err = gerror.Wrapf(gerror.New(".properties data format error"), `.properties data format error`)
			return
		}
		err = genMap(m[keys[0]].(map[string]interface{}), keys[1:], value)
		if err != nil {
			return
		}
	}
	return
}

// Decode converts properties format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	var (
		m           = make(map[string]interface{})
		bytesReader = bytes.NewReader(data)
		bufioReader = bufio.NewReader(bytesReader)
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
		if strings.Contains(line, "=") {
			values := strings.Split(line, "=")
			m[strings.TrimSpace(values[0])] = strings.TrimSpace(strings.Join(values[1:], "="))
		}
	}

	//[0][name]  replace to .0.name

	for k, v := range m {
		keys := strings.Split(strings.ReplaceAll(strings.ReplaceAll(k, "[", "."), "]", ""), ".")
		genMap(res, keys, v)
	}
	//.properties is not support array officially
	return res, nil
}
func point(base string) string {
	if len(strings.TrimSpace(base)) > 0 {
		return "."
	}
	return ""
}
func arrayToProperties(data []interface{}, w *bytes.Buffer, base string) (err error) {
	for i, value := range data {
		switch t := value.(type) {
		//the basic types
		case bool, uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, string, nil:
			{
				n, err := w.WriteString(fmt.Sprintf("%s%s%d=%v\n", base, point(base), i, value))
				if err != nil || n == 0 {
					return gerror.Wrapf(err, "w.WriteString failed")
				}
			}
			break
			//array
		case []interface{}:
			if err = arrayToProperties(value.([]interface{}), w, fmt.Sprintf("%s%s%v", base, point(base), i)); err != nil {
				return
			}
			break
			//object
		case map[string]interface{}:
			if err = toProperties(value.(map[string]interface{}), w, fmt.Sprintf("%s%s%d", base, point(base), i)); err != nil {
				return
			}
			break
		default:
			return gerror.Wrapf(err, fmt.Sprintf("data type not support %v", t))
		}
	}
	return
}

func toProperties(data map[string]interface{}, w *bytes.Buffer, base string) (err error) {
	for key, value := range data {
		switch t := value.(type) {
		//base types
		case bool, uint, uint8, uint16, uint32, uint64, int, int8, int16, int32, int64, float32, float64, string, nil:
			{
				n, err := w.WriteString(fmt.Sprintf("%s%s%s=%v\n", base, point(base), key, value))
				if err != nil || n == 0 {
					return gerror.Wrapf(err, "w.WriteString failed")
				}
			}
			break
			//array
		case []interface{}:
			if err = arrayToProperties(value.([]interface{}), w, fmt.Sprintf("%s%s%v", base, point(base), key)); err != nil {
				return
			}
			break
			//object
		case map[string]interface{}:
			if err = toProperties(value.(map[string]interface{}), w, fmt.Sprintf("%s%s%s", base, point(base), key)); err != nil {
				return
			}
			break
		default:
			return gerror.Wrapf(err, fmt.Sprintf("data type not support %v", t))
		}
	}
	return
}

// Encode converts map to properties format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	var (
		w = new(bytes.Buffer)
	)
	err = toProperties(data, w, "")
	if err != nil {
		return nil, err
	}
	res = make([]byte, w.Len())
	if n, err := w.Read(res); err != nil || n == 0 {
		return nil, gerror.Wrapf(err, "w.Read failed")
	}
	return res, nil
}

// ToJson convert .properties format to JSON.
func ToJson(data []byte) (res []byte, err error) {
	iniMap, err := Decode(data)
	if err != nil {
		return nil, err
	}
	return json.Marshal(iniMap)
}
