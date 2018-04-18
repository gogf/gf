// Copyright 2012-2018 Charles Banning. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file
// io.Reader --> map[string]interface{} or JSON string
// nothing magic - just implements generic Go case

package x2j

import (
	"encoding/json"
	"io"

	"github.com/clbanning/mxj"
)

// ToMap() - parse a XML io.Reader to a map[string]interface{}
func ToMap(rdr io.Reader, recast ...bool) (map[string]interface{}, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	return mxj.NewMapXmlReader(rdr, r)
}

// ToJson() - parse a XML io.Reader to a JSON string
func ToJson(rdr io.Reader, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := mxj.NewMapXmlReader(rdr, r)
	if m == nil || merr != nil {
		return "", merr
	}

	b, berr := json.Marshal(m)
	if berr != nil {
		return "", berr
	}

	return string(b), nil
}

// ToJsonIndent - the pretty form of ReaderToJson
func ToJsonIndent(rdr io.Reader, recast ...bool) (string, error) {
	var r bool
	if len(recast) == 1 {
		r = recast[0]
	}
	m, merr := mxj.NewMapXmlReader(rdr, r)
	if m == nil || merr != nil {
		return "", merr
	}

	b, berr := json.MarshalIndent(m, "", "  ")
	if berr != nil {
		return "", berr
	}

	// NOTE: don't have to worry about safe JSON marshaling with json.Marshal, since '<' and '>" are reservedin XML.
	return string(b), nil
}

// ReaderValuesFromTagPath - io.Reader version of ValuesFromTagPath()
func ReaderValuesFromTagPath(rdr io.Reader, path string, getAttrs ...bool) ([]interface{}, error) {
	var a bool
	if len(getAttrs) == 1 {
		a = getAttrs[0]
	}
	m, err := mxj.NewMapXmlReader(rdr)
	if err != nil {
		return nil, err
	}

	return ValuesFromKeyPath(m, path, a), nil
}

// ReaderValuesForTag - io.Reader version of ValuesForTag()
func ReaderValuesForTag(rdr io.Reader, tag string) ([]interface{}, error) {
	m, err := mxj.NewMapXmlReader(rdr)
	if err != nil {
		return nil, err
	}

	return ValuesForKey(m, tag), nil
}
