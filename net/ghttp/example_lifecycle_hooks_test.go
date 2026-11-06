// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
)

// ExampleServer_SetBeforeStart demonstrates how to use the SetBeforeStart hook.
// This hook is executed before the server starts listening.
func ExampleServer_SetBeforeStart() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World")
	})

	// Set a hook that executes before the server starts
	s.SetBeforeStart(func(s *ghttp.Server) error {
		fmt.Println("Server is preparing to start...")
		// Perform initialization tasks here, such as:
		// - Loading configuration
		// - Connecting to databases
		// - Initializing caches
		// - Validating required resources
		return nil
	})

	// If any before-start hook returns an error, the server will not start
	s.SetBeforeStart(func(s *ghttp.Server) error {
		// Return error to prevent server from starting
		// return fmt.Errorf("initialization failed")
		return nil
	})

	// Start the server
	// s.Run()
}

// ExampleServer_SetAfterStart demonstrates how to use the SetAfterStart hook.
// This hook is executed after the server has successfully started listening.
func ExampleServer_SetAfterStart() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("Hello World")
	})

	// Set a hook that executes after the server starts
	s.SetAfterStart(func(s *ghttp.Server) error {
		port := s.GetListenedPort()
		fmt.Printf("Server started successfully on port %d\n", port)
		// Perform post-startup tasks here, such as:
		// - Registering with service discovery
		// - Starting background jobs
		// - Sending startup notifications
		// - Logging server information
		return nil
	})

	// Multiple after-start hooks can be registered
	s.SetAfterStart(func(s *ghttp.Server) error {
		fmt.Println("Performing additional startup tasks...")
		return nil
	})

	// If an after-start hook returns an error, it will be logged
	// but the server will continue running
	s.SetAfterStart(func(s *ghttp.Server) error {
		// This error will be logged but won't stop the server
		// return fmt.Errorf("non-critical error")
		return nil
	})

	// Start the server
	// s.Run()
}

// Example_lifecycleHooks demonstrates a complete example
// using both before-start and after-start hooks together.
func Example_lifecycleHooks() {
	s := g.Server()
	s.SetAddr(":8080")

	// Register routes
	s.BindHandler("/health", func(r *ghttp.Request) {
		r.Response.WriteJson(g.Map{
			"status": "healthy",
		})
	})

	// Before-start hook: validate configuration
	s.SetBeforeStart(func(s *ghttp.Server) error {
		fmt.Println("1. Validating server configuration...")
		// Validate required configuration here
		return nil
	})

	// Before-start hook: initialize resources
	s.SetBeforeStart(func(s *ghttp.Server) error {
		fmt.Println("2. Initializing database connections...")
		// Initialize database connections here
		return nil
	})

	// After-start hook: register with service discovery
	s.SetAfterStart(func(s *ghttp.Server) error {
		port := s.GetListenedPort()
		fmt.Printf("3. Registering service on port %d with discovery...\n", port)
		// Register with service discovery here
		return nil
	})

	// After-start hook: log startup information
	s.SetAfterStart(func(s *ghttp.Server) error {
		fmt.Println("4. Server is ready to accept requests")
		return nil
	})

	// Start the server
	// s.Run()
}
