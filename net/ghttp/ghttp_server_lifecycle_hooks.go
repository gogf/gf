// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp

import "github.com/gogf/gf/v2/os/gctx"

// SetBeforeStart sets the hook function that is called before the server starts.
// Multiple hooks can be registered and they will be executed in the order they were registered.
// If any hook returns an error, the server will not start and the error will be returned.
//
// Example:
//
//	s := ghttp.GetServer()
//	s.SetBeforeStart(func(s *ghttp.Server) error {
//	    fmt.Println("Server is about to start")
//	    return nil
//	})
func (s *Server) SetBeforeStart(hook ServerHookFunc) {
	s.beforeStartHooks = append(s.beforeStartHooks, hook)
}

// SetAfterStart sets the hook function that is called after the server starts successfully.
// Multiple hooks can be registered and they will be executed in the order they were registered.
// The hooks are executed after all server listeners have been created and started.
// If any hook returns an error, the error will be logged but the server will continue running.
//
// Example:
//
//	s := ghttp.GetServer()
//	s.SetAfterStart(func(s *ghttp.Server) error {
//	    fmt.Println("Server started successfully")
//	    return nil
//	})
func (s *Server) SetAfterStart(hook ServerHookFunc) {
	s.afterStartHooks = append(s.afterStartHooks, hook)
}

// executeBeforeStartHooks executes all registered before-start hooks.
// It returns an error if any hook fails, which will prevent the server from starting.
func (s *Server) executeBeforeStartHooks() error {
	for _, hook := range s.beforeStartHooks {
		if err := hook(s); err != nil {
			return err
		}
	}
	return nil
}

// executeAfterStartHooks executes all registered after-start hooks.
// Errors from hooks are logged but do not stop the server.
func (s *Server) executeAfterStartHooks() {
	var ctx = gctx.GetInitCtx()
	for _, hook := range s.afterStartHooks {
		if err := hook(s); err != nil {
			s.Logger().Errorf(ctx, `after-start hook error: %+v`, err)
		}
	}
}
