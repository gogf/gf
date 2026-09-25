// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Option definitions for app initialization.

package gapp

import "context"

// Option is the interface for app initialization options.
// It is applied during App.Boot() to configure application-level
// concerns before servers start.
//
// Implementations can return an optional cleanup function from Apply,
// which will be called during App.Stop(ctx, graceful) in reverse registration order.
type Option interface {
	// Apply initializes resources for the application.
	// It returns an optional cleanup function that will be called
	// during App.Stop(ctx, graceful), or nil if no cleanup is needed.
	Apply(ctx context.Context, app *App) (func(ctx context.Context), error)
}

// optionFunc is an adapter to allow the use of ordinary functions as Options.
type optionFunc func(ctx context.Context, app *App) (func(ctx context.Context), error)

// Apply implements the Option interface.
func (f optionFunc) Apply(ctx context.Context, app *App) (func(ctx context.Context), error) {
	return f(ctx, app)
}

// NewOption creates an Option from a simple initialization function.
// The function runs once during Boot(). No cleanup is registered.
func NewOption(f func(ctx context.Context, app *App)) Option {
	return optionFunc(func(ctx context.Context, app *App) (func(ctx context.Context), error) {
		f(ctx, app)
		return nil, nil
	})
}

// NewOptionWithHook creates an Option from an initialization function
// that also returns an optional cleanup function.
// The init function runs during Boot(), and the returned cleanup
// function (if non-nil) runs during App.Stop(ctx, graceful) in reverse registration order.
func NewOptionWithHook(f func(ctx context.Context, app *App) (func(ctx context.Context), error)) Option {
	return optionFunc(f)
}
