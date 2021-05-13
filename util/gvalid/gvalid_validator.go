// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid

import (
	"context"
	"github.com/gogf/gf/i18n/gi18n"
)

// Validator is the validation manager.
type Validator struct {
	ctx         context.Context // Context containing custom context variables.
	i18nManager *gi18n.Manager  // I18n manager for error message translation.

}

// New creates and returns a new Validator.
func New() *Validator {
	return &Validator{
		ctx:         context.TODO(),   // Initialize an empty context.
		i18nManager: gi18n.Instance(), // Use default i18n manager.
	}
}

// I18n sets the i18n manager for the validator.
func (v *Validator) I18n(i18nManager *gi18n.Manager) *Validator {
	v.i18nManager = i18nManager
	return v
}

// Ctx is a chaining operation function which sets the context for next validation.
func (v *Validator) Ctx(ctx context.Context) *Validator {
	v.ctx = ctx
	return v
}
