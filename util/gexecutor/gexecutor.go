// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package gexecutor provides a generic executor that allows executing functions
// with before and after hooks in a chainable way.
package gexecutor

import (
	"context"
	"errors"
)

// ErrMainFuncNotSet is returned when the main function is not set before calling Do method.
var ErrMainFuncNotSet = errors.New("main function is not set")

// Executor represents a generic executor that can execute a main function with optional
// before and after hooks. It supports generic input and output types.
// T is the input type, R is the output type.
type Executor[T any, R any] struct {
	// input holds the input data for the execution
	input T
	// mainFunc is the main function to be executed
	mainFunc func(ctx context.Context, input T) (R, error)
	// beforeFunc is an optional function that runs before the main function
	beforeFunc func(ctx context.Context, input T)
	// afterFunc is an optional function that runs after the main function
	afterFunc func(ctx context.Context, result R)
}

// New creates and returns a new Executor instance with the given input.
// T is the input type, R is the output type.
func New[T any, R any](input T) *Executor[T, R] {
	return &Executor[T, R]{input: input}
}

// WithMain sets the main function for the executor.
// The main function takes a context and input of type T, and returns a result of type R and an error.
// Returns a new Executor instance with the main function set.
func (e *Executor[T, R]) WithMain(f func(context.Context, T) (R, error)) *Executor[T, R] {
	return &Executor[T, R]{
		input:      e.input,
		mainFunc:   f,
		beforeFunc: e.beforeFunc,
		afterFunc:  e.afterFunc,
	}
}

// WithBefore sets the before hook function for the executor.
// The before function runs before the main function and receives the context and input.
// Returns a new Executor instance with the before function set.
func (e *Executor[T, R]) WithBefore(f func(context.Context, T)) *Executor[T, R] {
	return &Executor[T, R]{
		input:      e.input,
		mainFunc:   e.mainFunc,
		beforeFunc: f,
		afterFunc:  e.afterFunc,
	}
}

// WithAfter sets the after hook function for the executor.
// The after function runs after the main function and receives the context and result.
// Returns a new Executor instance with the after function set.
func (e *Executor[T, R]) WithAfter(f func(context.Context, R)) *Executor[T, R] {
	return &Executor[T, R]{
		input:      e.input,
		mainFunc:   e.mainFunc,
		beforeFunc: e.beforeFunc,
		afterFunc:  f,
	}
}

// Do executes the executor with the given context.
// It first runs the before function (if set), then the main function (required),
// and finally the after function (if set).
// Returns the result of the main function and any error that occurred.
func (e *Executor[T, R]) Do(ctx context.Context) (R, error) {
	var zero R
	if ctx == nil {
		ctx = context.Background()
	}

	if e.beforeFunc != nil {
		e.beforeFunc(ctx, e.input)
	}

	if e.mainFunc == nil {
		return zero, ErrMainFuncNotSet
	}
	result, err := e.mainFunc(ctx, e.input)
	if err != nil {
		return zero, err
	}

	if e.afterFunc != nil {
		e.afterFunc(ctx, result)
	}

	return result, nil
}
