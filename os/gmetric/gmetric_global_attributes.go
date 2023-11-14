// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import "fmt"

type GlobalAttributesOption struct {
	Instrument        string
	InstrumentVersion string
}

var (
	globalAttributeMap = make(map[string]Attributes)
)

func SetGlobalAttributes(attrs Attributes, option ...GlobalAttributesOption) {
	var usedOption GlobalAttributesOption
	if len(option) > 0 {
		usedOption = option[0]
	}
	var mapKey = usedOption.String()
	if _, ok := globalAttributeMap[mapKey]; !ok {
		globalAttributeMap[mapKey] = make(Attributes, 0)
	}
	globalAttributeMap[mapKey] = append(globalAttributeMap[mapKey], attrs...)
}

func GetGlobalAttributes(option GlobalAttributesOption) Attributes {
	return globalAttributeMap[option.String()]
}

// String converts and returns GlobalAttributesOption as string.
func (o GlobalAttributesOption) String() string {
	if o.Instrument != "" {
		return fmt.Sprintf(`%s@%s`, o.Instrument, o.InstrumentVersion)
	}
	return ""
}
