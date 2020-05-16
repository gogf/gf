// Copyright 2018 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

// checkLuHn checks <value> with LUHN algorithm.
// It's usually used for bank card number validation.
func checkLuHn(value string) bool {
	var (
		sum     = 0
		nDigits = len(value)
		parity  = nDigits % 2
	)
	for i := 0; i < nDigits; i++ {
		var digit = int(value[i] - 48)
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
