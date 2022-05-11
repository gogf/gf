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
	"io/ioutil"
	"os"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/spf13/viper"
)

// Decode converts properties format to map.
func Decode(data []byte) (res map[string]interface{}, err error) {
	res = make(map[string]interface{})
	var (
		bytesReader = bytes.NewReader(data)
		bufioReader = bufio.NewReader(bytesReader)
	)
	vp := viper.New()
	vp.SetConfigType("properties")
	if err = vp.ReadConfig(bufioReader); err != nil {
		err = gerror.Wrapf(err, `viper ReadConfog failed`)
		return nil, err
	}
	res = vp.AllSettings()
	return res, nil
}

// Encode converts map to properties format.
func Encode(data map[string]interface{}) (res []byte, err error) {
	var (
		//w  = new(bytes.Buffer)
		vp          = viper.New()
		tmpFileName = fmt.Sprintf("vp_tmp_config%s", gtime.Now().Format("YmdHis"))
	)
	vp.SetConfigName(tmpFileName)
	vp.SetConfigType("properties")
	vp.AddConfigPath(".")
	vp.MergeConfigMap(data)
	if err = vp.SafeWriteConfig(); err != nil {
		err = gerror.Wrapf(err, `viper WriteConfog failed`)
		return nil, err
	}
	tmpFileNameP := tmpFileName + ".properties"
	res, err = ioutil.ReadFile(tmpFileNameP)
	defer os.Remove(tmpFileNameP)

	if err != nil {
		err = gerror.Wrapf(err, `Read viper tmp file failed`)
		return nil, err
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
