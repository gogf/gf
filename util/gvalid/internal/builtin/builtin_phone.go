// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"

	"github.com/gogf/gf/v2/text/gregex"
)

// 1. China Mobile:
//    134, 135, 136, 137, 138, 139, 150, 151, 152, 157, 158, 159, 182, 183, 184, 187, 188,
//    178(4G), 147(Net)；
//    172
//
// 2. China Unicom:
//    130, 131, 132, 155, 156, 185, 186 ,176(4G), 145(Net), 175
//
// 3. China Telecom:
//    133, 153, 180, 181, 189, 177(4G)
//
// 4. Satelite:
//    1349
//
// 5. Virtual:
//    170, 173
//
// 6. 2018:
//    16x, 19x

type RulePhone struct{}

func init() {
	Register(&RulePhone{})
}

func (r *RulePhone) Name() string {
	return "phone"
}

func (r *RulePhone) Message() string {
	return "The {attribute} value `{value}` is not a valid phone number"
}

func (r *RulePhone) Run(in RunInput) error {
	ok := gregex.IsMatchString(
		`^13[\d]{9}$|^14[5,7]{1}\d{8}$|^15[^4]{1}\d{8}$|^16[\d]{9}$|^17[0,2,3,5,6,7,8]{1}\d{8}$|^18[\d]{9}$|^19[\d]{9}$`,
		in.Value.String(),
	)
	if ok {
		return nil
	}
	return errors.New(in.Message)
}
