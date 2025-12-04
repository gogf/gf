// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package converter

// Rune converts `any` to rune.
func (c *Converter) Rune(anyInput any) (rune, error) {
	if v, ok := anyInput.(rune); ok {
		return v, nil
	}
	v, err := c.Int32(anyInput)
	if err != nil {
		return 0, err
	}
	return v, nil
}

// Runes converts `any` to []rune.
func (c *Converter) Runes(anyInput any) ([]rune, error) {
	if v, ok := anyInput.([]rune); ok {
		return v, nil
	}
	s, err := c.String(anyInput)
	if err != nil {
		return nil, err
	}
	return []rune(s), nil
}
