// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://gitee.com/johng/gp.

package gparser

// SetSplitChar sets the separator char for hierarchical data access.
func (p *Parser) SetSplitChar(char byte) {
	p.json.SetSplitChar(char)
}

// SetViolenceCheck enables/disables violence check for hierarchical data access.
func (p *Parser) SetViolenceCheck(check bool) {
	p.json.SetViolenceCheck(check)
}