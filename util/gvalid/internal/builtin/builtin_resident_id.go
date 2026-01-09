// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gogf/gf/v2/text/gregex"
)

// RuleResidentId implements `resident-id` rule:
// Resident id number.
//
// Format: resident-id
type RuleResidentId struct{}

func init() {
	Register(RuleResidentId{})
}

func (r RuleResidentId) Name() string {
	return "resident-id"
}

func (r RuleResidentId) Message() string {
	return "The {field} value `{value}` is not a valid resident id number"
}

func (r RuleResidentId) Run(in RunInput) error {
	if r.checkResidentId(in.Value.String()) {
		return nil
	}
	return errors.New(in.Message)
}

// checkResidentId checks whether given id a china resident id number.
//
// xxxxxx yyyy MM dd 375 0  18 digits
// xxxxxx   yy MM dd  75 0  15 digits
//
// Region:     [1-9]\d{5}
// First two digits of year: (18|19|([23]\d))  1800-2399
// Last two digits of year: \d{2}
// Month:     ((0[1-9])|(10|11|12))
// Day:       (([0-2][1-9])|10|20|30|31) Leap year cannot prohibit 29+
//
// Three sequential digits: \d{3}
// Two sequential digits: \d{2}
// Check code:   [0-9Xx]
//
// 18 digits: ^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$
// 15 digits: ^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$
//
// Total:
// (^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)
func (r RuleResidentId) checkResidentId(id string) bool {
	id = strings.ToUpper(strings.TrimSpace(id))
	if len(id) != 18 {
		return false
	}
	var (
		weightFactor = []int{7, 9, 10, 5, 8, 4, 2, 1, 6, 3, 7, 9, 10, 5, 8, 4, 2}
		checkCode    = []byte{'1', '0', 'X', '9', '8', '7', '6', '5', '4', '3', '2'}
		last         = id[17]
		num          = 0
	)
	for i := 0; i < 17; i++ {
		tmp, err := strconv.Atoi(string(id[i]))
		if err != nil {
			return false
		}
		num = num + tmp*weightFactor[i]
	}
	if checkCode[num%11] != last {
		return false
	}

	return gregex.IsMatchString(
		`(^[1-9]\d{5}(18|19|([23]\d))\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}[0-9Xx]$)|(^[1-9]\d{5}\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\d{3}$)`,
		id,
	)
}
