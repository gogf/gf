// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gmetric

import (
	"sync"

	"github.com/gogf/gf/v2/text/gregex"
)

// SetGlobalAttributesOption binds the global attributes to certain instrument.
type SetGlobalAttributesOption struct {
	Instrument        string // Instrument specifies the instrument name.
	InstrumentVersion string // Instrument specifies the instrument version.
	InstrumentPattern string // InstrumentPattern specifies instrument by regular expression on Instrument name.
}

// GetGlobalAttributesOption binds the global attributes to certain instrument.
type GetGlobalAttributesOption struct {
	Instrument        string // Instrument specifies the instrument name.
	InstrumentVersion string // Instrument specifies the instrument version.
}

type globalAttributeItem struct {
	Attributes
	SetGlobalAttributesOption
}

var (
	globalAttributesMu sync.Mutex
	// globalAttributes stores the global attributes to a map.
	globalAttributes = make([]globalAttributeItem, 0)
)

// SetGlobalAttributes appends global attributes according `SetGlobalAttributesOption`.
// It appends global attributes to all metrics if given `SetGlobalAttributesOption` is empty.
// It appends global attributes to certain instrument by given `SetGlobalAttributesOption`.
func SetGlobalAttributes(attrs Attributes, option ...SetGlobalAttributesOption) {
	globalAttributesMu.Lock()
	defer globalAttributesMu.Unlock()
	var usedOption SetGlobalAttributesOption
	if len(option) > 0 {
		usedOption = option[0]
	}
	globalAttributes = append(
		globalAttributes, globalAttributeItem{
			Attributes:                attrs,
			SetGlobalAttributesOption: usedOption,
		},
	)
}

// GetGlobalAttributes retrieves and returns the global attributes by `GetGlobalAttributesOption`.
// It returns the global attributes if given `GetGlobalAttributesOption` is empty.
// It returns global attributes of certain instrument if `GetGlobalAttributesOption` is not empty.
func GetGlobalAttributes(option GetGlobalAttributesOption) Attributes {
	globalAttributesMu.Lock()
	defer globalAttributesMu.Unlock()
	var attributes = make(Attributes, 0)
	for _, attrItem := range globalAttributes {
		if option.InstrumentVersion != "" && attrItem.InstrumentVersion != option.InstrumentVersion {
			continue
		}
		if attrItem.InstrumentPattern == "" {
			if attrItem.Instrument != option.Instrument {
				continue
			}
		} else {
			if !gregex.IsMatchString(attrItem.InstrumentPattern, option.Instrument) {
				continue
			}
		}
		attributes = append(attributes, attrItem.Attributes...)
	}
	return attributes
}
