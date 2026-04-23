// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package builtin

import (
	"errors"
	"strings"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/internal/empty"
	"github.com/gogf/gf/v2/util/gconv"
	"github.com/gogf/gf/v2/util/gutil"
)

// RuleExcludedIf implements `excluded-if` rule:
// Current field must be empty if all given field/value pairs are equal.
//
// Format: excluded-if:field,value,...
type RuleExcludedIf struct{}

func init() {
	Register(RuleExcludedIf{})
	Register(RuleExcludedUnless{})
	Register(RuleExcludedWith{})
	Register(RuleExcludedWithAll{})
	Register(RuleExcludedWithout{})
	Register(RuleExcludedWithoutAll{})
}

func (r RuleExcludedIf) Name() string {
	return "excluded-if"
}

func (r RuleExcludedIf) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedIf) Run(in RunInput) error {
	matched, err := matchFieldValuePairs(in, r.Name())
	if err != nil {
		return err
	}
	if matched && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

// RuleExcludedUnless implements `excluded-unless` rule:
// Current field must be empty unless all given field/value pairs are equal.
//
// Format: excluded-unless:field,value,...
type RuleExcludedUnless struct{}

func (r RuleExcludedUnless) Name() string {
	return "excluded-unless"
}

func (r RuleExcludedUnless) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedUnless) Run(in RunInput) error {
	matched, err := matchFieldValuePairs(in, r.Name())
	if err != nil {
		return err
	}
	if !matched && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

// RuleExcludedWith implements `excluded-with` rule:
// Current field must be empty if any given fields are not empty.
//
// Format: excluded-with:field1,field2,...
type RuleExcludedWith struct{}

func (r RuleExcludedWith) Name() string {
	return "excluded-with"
}

func (r RuleExcludedWith) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedWith) Run(in RunInput) error {
	if anyFieldHasValue(in) && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

// RuleExcludedWithAll implements `excluded-with-all` rule:
// Current field must be empty if all given fields are not empty.
//
// Format: excluded-with-all:field1,field2,...
type RuleExcludedWithAll struct{}

func (r RuleExcludedWithAll) Name() string {
	return "excluded-with-all"
}

func (r RuleExcludedWithAll) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedWithAll) Run(in RunInput) error {
	if allFieldsHaveValue(in) && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

// RuleExcludedWithout implements `excluded-without` rule:
// Current field must be empty if any given fields are empty.
//
// Format: excluded-without:field1,field2,...
type RuleExcludedWithout struct{}

func (r RuleExcludedWithout) Name() string {
	return "excluded-without"
}

func (r RuleExcludedWithout) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedWithout) Run(in RunInput) error {
	if anyFieldMissingValue(in) && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

// RuleExcludedWithoutAll implements `excluded-without-all` rule:
// Current field must be empty if all given fields are empty.
//
// Format: excluded-without-all:field1,field2,...
type RuleExcludedWithoutAll struct{}

func (r RuleExcludedWithoutAll) Name() string {
	return "excluded-without-all"
}

func (r RuleExcludedWithoutAll) Message() string {
	return "The {field} field is excluded"
}

func (r RuleExcludedWithoutAll) Run(in RunInput) error {
	if allFieldsMissingValue(in) && !isRequiredEmpty(in.Value.Val()) {
		return errors.New(in.Message)
	}
	return nil
}

func matchFieldValuePairs(in RunInput, ruleName string) (bool, error) {
	var (
		array   = strings.Split(in.RulePattern, ",")
		dataMap = in.Data.Map()
	)
	if len(array)%2 != 0 {
		return false, gerror.NewCodef(
			gcode.CodeInvalidParameter,
			`invalid "%s" rule pattern: %s`,
			ruleName,
			in.RulePattern,
		)
	}
	if len(array) == 0 {
		return false, nil
	}
	for i := 0; i < len(array); i += 2 {
		var (
			fieldName = array[i]
			expect    = array[i+1]
			_, found  = gutil.MapPossibleItemByKey(dataMap, fieldName)
		)
		if in.Option.CaseInsensitive {
			if !strings.EqualFold(expect, gconv.String(found)) {
				return false, nil
			}
		} else if strings.Compare(expect, gconv.String(found)) != 0 {
			return false, nil
		}
	}
	return true, nil
}

func anyFieldHasValue(in RunInput) bool {
	for _, fieldName := range strings.Split(in.RulePattern, ",") {
		_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), fieldName)
		if !empty.IsEmpty(foundValue) {
			return true
		}
	}
	return false
}

func allFieldsHaveValue(in RunInput) bool {
	array := strings.Split(in.RulePattern, ",")
	if len(array) == 0 {
		return false
	}
	for _, fieldName := range array {
		_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), fieldName)
		if empty.IsEmpty(foundValue) {
			return false
		}
	}
	return true
}

func anyFieldMissingValue(in RunInput) bool {
	for _, fieldName := range strings.Split(in.RulePattern, ",") {
		_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), fieldName)
		if empty.IsEmpty(foundValue) {
			return true
		}
	}
	return false
}

func allFieldsMissingValue(in RunInput) bool {
	array := strings.Split(in.RulePattern, ",")
	if len(array) == 0 {
		return false
	}
	for _, fieldName := range array {
		_, foundValue := gutil.MapPossibleItemByKey(in.Data.Map(), fieldName)
		if !empty.IsEmpty(foundValue) {
			return false
		}
	}
	return true
}
