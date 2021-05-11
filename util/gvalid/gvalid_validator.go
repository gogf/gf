// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import "context"

// Validator is the validation manager.
type Validator struct {
	i18nLang string          // I18n language.
	ctx      context.Context // Context containing custom context variables.
}

// New creates and returns a new Validator.
func New() *Validator {
	return &Validator{}
}

// Clone creates and returns a new Validator which is a shallow copy of current one.
func (v *Validator) Clone() *Validator {
	newValidator := New()
	*newValidator = *v
	return newValidator
}

// I18n is a chaining operation function which sets the I18n language for next validation.
func (v *Validator) I18n(language string) *Validator {
	newValidator := v.Clone()
	newValidator.i18nLang = language
	return newValidator
}

// Ctx is a chaining operation function which sets the context for next validation.
func (v *Validator) Ctx(ctx context.Context) *Validator {
	newValidator := v.Clone()
	newValidator.ctx = ctx
	return newValidator
}
