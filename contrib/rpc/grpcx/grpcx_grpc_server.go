// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Implements gRPC serving and service registration lifecycle management.

package grpcx

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"sort"
	"sync"
	"time"

	"google.golang.org/grpc"

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gproc"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
)

// GrpcServer is the server for GRPC protocol.
type GrpcServer struct {
	Server             *grpc.Server
	config             *GrpcServerConfig
	listenerMu         sync.RWMutex          // Guards listener, Server replacement, serverStopped, and registrations.
	listener           net.Listener          // Non-nil only while the current Server is serving.
	serverStopped      bool                  // Server was stopped after serving and must be rebuilt before the next start.
	registrations      []serviceRegistration // RegisterService calls replayed onto a rebuilt Server.
	services           []gsvc.Service        // Service definitions retained across registration attempts.
	registeredServices []gsvc.Service        // Registry handles used exclusively for deregistration.
	waitGroup          sync.WaitGroup
	registrar          gsvc.Registrar
	serviceMu          sync.Mutex
}

// Service implements gsvc.Service interface.
type Service struct {
	gsvc.Service
	Endpoints gsvc.Endpoints
}

// serviceRegistration records one RegisterService call for replay onto a rebuilt grpc.Server.
type serviceRegistration struct {
	desc *grpc.ServiceDesc
	impl any
}

// Compile-time check that GrpcServer can be passed to generated Register*Server functions.
var _ grpc.ServiceRegistrar = (*GrpcServer)(nil)

// New creates and returns a grpc server.
func (s modServer) New(conf ...*GrpcServerConfig) *GrpcServer {
	autoLoadAndRegisterFileRegistry()

	var (
		ctx    = gctx.GetInitCtx()
		config *GrpcServerConfig
	)
	if len(conf) > 0 {
		config = conf[0]
	} else {
		config = s.NewConfig()
	}
	if config.Address == "" {
		config.Address = defaultListenAddress
	}
	if !gstr.Contains(config.Address, ":") {
		g.Log().Fatal(ctx, "invalid service address, should contain listening port")
	}
	if config.Logger == nil {
		config.Logger = glog.New()
	}
	grpcServer := &GrpcServer{
		config:    config,
		registrar: gsvc.GetRegistry(),
	}
	grpcServer.config.Options = append([]grpc.ServerOption{
		s.ChainUnary(
			s.UnaryTracing,
			grpcServer.UnaryLogger,
			s.UnaryRecover,
			s.UnaryAllowNilRes,
			s.UnaryError,
		),
		s.ChainStream(
			s.StreamTracing,
		),
	}, grpcServer.config.Options...)
	grpcServer.Server = grpc.NewServer(grpcServer.config.Options...)
	return grpcServer
}

// Service binds service list to current server.
// Server will automatically register the service list after it starts.
func (s *GrpcServer) Service(services ...gsvc.Service) {
	s.serviceMu.Lock()
	defer s.serviceMu.Unlock()
	s.services = append(s.services, services...)
}

// RegisterService registers a service implementation on the underlying grpc.Server
// and records it, so a start after Stop can replay it onto a fresh grpc.Server.
// It implements grpc.ServiceRegistrar, so generated code can register through
// GrpcServer, for example pb.RegisterGreeterServer(s, &Controller{}).
// Services registered on the Server field directly cannot be replayed, and a
// restart after a stop that served them returns an error.
func (s *GrpcServer) RegisterService(desc *grpc.ServiceDesc, impl any) {
	s.listenerMu.Lock()
	defer s.listenerMu.Unlock()
	s.Server.RegisterService(desc, impl)
	s.registrations = append(s.registrations, serviceRegistration{desc: desc, impl: impl})
}

// rebuildServerLocked replaces a stopped grpc.Server, which grpc-go cannot serve again,
// with a fresh one carrying the recorded registrations; the caller holds listenerMu.
func (s *GrpcServer) rebuildServerLocked() error {
	recorded := make(map[string]struct{}, len(s.registrations))
	for _, registration := range s.registrations {
		recorded[registration.desc.ServiceName] = struct{}{}
	}
	var unrecorded []string
	for name := range s.Server.GetServiceInfo() {
		if _, ok := recorded[name]; !ok {
			unrecorded = append(unrecorded, name)
		}
	}
	if len(unrecorded) > 0 {
		// Serving without these services would report a healthy start that answers Unimplemented.
		sort.Strings(unrecorded)
		return gerror.NewCodef(
			gcode.CodeInvalidOperation,
			`grpc server cannot restart after stop: services %v were registered on GrpcServer.Server directly; register them through GrpcServer.RegisterService, for example pb.RegisterXxxServer(s, impl)`,
			unrecorded,
		)
	}
	server := grpc.NewServer(s.config.Options...)
	for _, registration := range s.registrations {
		server.RegisterService(registration.desc, registration.impl)
	}
	s.Server = server
	s.serverStopped = false
	return nil
}

// detachServer clears the listener as a stop begins, so the next start can bind again,
// and reports whether the returned Server was serving and therefore needs grpc-go's stop.
// A Server that never served keeps its registrations and stays usable.
func (s *GrpcServer) detachServer() (*grpc.Server, bool) {
	s.listenerMu.Lock()
	defer s.listenerMu.Unlock()
	if s.listener == nil {
		return s.Server, false
	}
	s.listener = nil
	s.serverStopped = true
	return s.Server, true
}

// serveFailureAction selects how an unexpected Server.Serve error is reported.
type serveFailureAction int

const (
	// serveFailureFatal terminates the process. Used by blocking Run().
	serveFailureFatal serveFailureAction = iota
	// serveFailureLog logs the error and returns. Used by StartManaged().
	serveFailureLog
)

// serve binds the listener and initializes asynchronous serving and service registration.
// It does not block on OS signal handling and is intended for external lifecycle managers.
// failureAction controls whether an unexpected Server.Serve error terminates the process
// (blocking Run) or is only logged (StartManaged).
func (s *GrpcServer) serve(failureAction serveFailureAction) error {
	ctx := gctx.GetInitCtx()

	s.listenerMu.Lock()
	if s.listener != nil {
		s.listenerMu.Unlock()
		return gerror.NewCode(gcode.CodeInvalidOperation, "grpc server already started")
	}
	if s.serverStopped {
		if err := s.rebuildServerLocked(); err != nil {
			s.listenerMu.Unlock()
			return err
		}
	}
	listener, err := net.Listen("tcp", s.config.Address)
	if err != nil {
		s.listenerMu.Unlock()
		return err
	}
	s.listener = listener
	// Serve the instance bound to this listener even if a later restart replaces s.Server.
	server := s.Server
	s.listenerMu.Unlock()

	if failureAction == serveFailureFatal {
		// Legacy registrars may synchronously verify an RPC before accepting registration.
		go s.doServeAsynchronously(ctx, server, listener, failureAction)
	}
	if err := s.doServiceRegister(failureAction); err != nil {
		// Managed starts register before Serve because a stopped grpc.Server cannot be restarted.
		if closeErr := listener.Close(); closeErr != nil {
			s.Logger().Errorf(ctx, `close listener after registration failure: %+v`, closeErr)
		}
		s.listenerMu.Lock()
		s.listener = nil
		s.listenerMu.Unlock()
		return err
	}
	if failureAction == serveFailureLog {
		go s.doServeAsynchronously(ctx, server, listener, failureAction)
	}
	s.Logger().Infof(
		ctx,
		"pid[%d]: grpc server started listening on [%s]",
		gproc.Pid(), s.GetListenedAddress(),
	)
	return nil
}

// Run starts the server in blocking way.
func (s *GrpcServer) Run() {
	ctx := gctx.GetInitCtx()
	if err := s.serve(serveFailureFatal); err != nil {
		s.Logger().Fatalf(ctx, `%+v`, err)
	}
	s.doSignalListen()
}

// StartManaged starts serving under external lifecycle management without signal handling.
func (s *GrpcServer) StartManaged() error {
	return s.serve(serveFailureLog)
}

// doServeAsynchronously serves gRPC connections of server on the given listener.
// Unexpected Serve errors terminate the process when failureAction is
// serveFailureFatal so blocking Run() cannot sit in signal wait with a dead
// listener. StartManaged uses serveFailureLog so the caller retains lifecycle control.
func (s *GrpcServer) doServeAsynchronously(ctx context.Context, server *grpc.Server, listener net.Listener, failureAction serveFailureAction) {
	if err := server.Serve(listener); err != nil && err != grpc.ErrServerStopped {
		if failureAction == serveFailureFatal {
			s.Logger().Fatalf(ctx, `%+v`, err)
		} else {
			s.Logger().Errorf(ctx, `grpc server serve error: %+v`, err)
		}
	}
}

// doSignalListen does signal listening and handling for gracefully shutdown.
func (s *GrpcServer) doSignalListen() {
	var ctx = context.Background()
	gproc.AddSigHandlerShutdown(func(sig os.Signal) {
		s.Logger().Infof(ctx, "signal received: %s, gracefully shutting down", sig.String())
		// Deregister services when shutdown signal triggers.
		s.doServiceDeregister()
		time.Sleep(time.Second)
		s.Stop()
	})
	gproc.Listen()
	// Deregister services when process ends.
	s.doServiceDeregister()
}

// Logger is alias of GetLogger.
func (s *GrpcServer) Logger() *glog.Logger {
	return s.config.Logger
}

// doServiceRegister registers current service to Registry.
// Register failures terminate the process on the blocking Run path.
// StartManaged returns the error so the caller can stop other servers.
func (s *GrpcServer) doServiceRegister(failureAction serveFailureAction) error {
	if s.registrar == nil {
		return nil
	}
	s.serviceMu.Lock()
	defer s.serviceMu.Unlock()
	ctx := gctx.GetInitCtx()
	if err := s.doServiceDeregisterLocked(ctx); err != nil {
		return gerror.Wrap(err, "cleanup previous service registration failed")
	}
	if len(s.services) == 0 {
		s.services = []gsvc.Service{&gsvc.LocalService{
			Name:     s.config.Name,
			Metadata: gsvc.Metadata{},
		}}
	}
	var (
		protocol           = `grpc`
		registeredServices = make([]gsvc.Service, 0, len(s.services))
	)
	// Register services after the listener is bound.
	for _, service := range s.services {
		service = &gsvc.LocalService{
			Name:      service.GetName(),
			Endpoints: s.calculateListenedEndpoints(ctx),
			Metadata:  service.GetMetadata(),
		}
		service.GetMetadata().Sets(gsvc.Metadata{
			gsvc.MDProtocol: protocol,
		})
		s.Logger().Debugf(ctx, `service register: %+v`, service)
		if len(service.GetEndpoints()) == 0 {
			s.Logger().Warningf(ctx, `no endpoints found to register service, abort service registering`)
			return nil
		}
		registered, err := s.registrar.Register(ctx, service)
		if err != nil {
			// Retain handles whose rollback fails so shutdown or the next start can retry.
			s.registeredServices = registeredServices
			err = errors.Join(err, s.doServiceDeregisterLocked(ctx))
			if failureAction == serveFailureFatal {
				s.Logger().Fatalf(ctx, `%+v`, err)
			}
			return err
		}
		registeredServices = append(registeredServices, registered)
	}
	s.registeredServices = registeredServices
	return nil
}

// doServiceDeregister de-registers current service from Registry.
func (s *GrpcServer) doServiceDeregister() {
	if s.registrar == nil {
		return
	}
	s.serviceMu.Lock()
	defer s.serviceMu.Unlock()
	ctx := gctx.GetInitCtx()
	if err := s.doServiceDeregisterLocked(ctx); err != nil {
		s.Logger().Errorf(ctx, `%+v`, err)
	}
}

// doServiceDeregisterLocked removes registrations while retaining failed handles; the caller holds serviceMu.
func (s *GrpcServer) doServiceDeregisterLocked(ctx context.Context) error {
	var (
		failedServices []gsvc.Service
		deregisterErr  error
	)
	for _, service := range s.registeredServices {
		s.Logger().Debugf(ctx, `service deregister: %+v`, service)
		if err := s.registrar.Deregister(ctx, service); err != nil {
			failedServices = append(failedServices, service)
			deregisterErr = errors.Join(deregisterErr, err)
		}
	}
	s.registeredServices = failedServices
	return deregisterErr
}

// Start starts the server in no-blocking way.
func (s *GrpcServer) Start() {
	s.waitGroup.Add(1)
	go s.doStartAsynchronously()
}

func (s *GrpcServer) doStartAsynchronously() {
	defer s.waitGroup.Done()
	s.Run()
}

// Wait works with Start, which blocks current goroutine until the server stops.
func (s *GrpcServer) Wait() {
	s.waitGroup.Wait()
}

// Stop gracefully stops the server.
// The server can be started again afterwards; see RegisterService.
func (s *GrpcServer) Stop() {
	server, serving := s.detachServer()
	s.doServiceDeregister()
	if serving {
		server.GracefulStop()
	}
}

// StopForceful forcibly stops the server and de-registers services from the registry.
// The server can be started again afterwards; see RegisterService.
func (s *GrpcServer) StopForceful() {
	server, serving := s.detachServer()
	s.doServiceDeregister()
	if serving {
		server.Stop()
	}
}

// GetConfig returns the configuration of current Server.
func (s *GrpcServer) GetConfig() *GrpcServerConfig {
	return s.config
}

// GetListenedAddress retrieves and returns the address string which are listened by current server.
func (s *GrpcServer) GetListenedAddress() string {
	if !gstr.Contains(s.config.Address, FreePortAddress) {
		return s.config.Address
	}
	var (
		address      = s.config.Address
		listenedPort = s.GetListenedPort()
	)
	address = gstr.Replace(address, FreePortAddress, fmt.Sprintf(`:%d`, listenedPort))
	return address
}

// GetListenedPort retrieves and returns one port which is listened to by current server.
func (s *GrpcServer) GetListenedPort() int {
	s.listenerMu.RLock()
	defer s.listenerMu.RUnlock()
	if ln := s.listener; ln != nil {
		return ln.Addr().(*net.TCPAddr).Port
	}
	return -1
}

func (s *GrpcServer) calculateListenedEndpoints(ctx context.Context) gsvc.Endpoints {
	var (
		configAddr = s.config.Address
		endpoints  = make(gsvc.Endpoints, 0)
		addresses  = s.config.Endpoints
	)
	if len(addresses) == 0 {
		addresses = gstr.SplitAndTrim(configAddr, ",")
	}
	for _, address := range addresses {
		var (
			addrArray     = gstr.Split(address, ":")
			listenedIps   []string
			listenedPorts []int
		)
		if len(addrArray) == 1 {
			configItemName := "address"
			if len(s.config.Endpoints) != 0 {
				configItemName = "endpoint"
			}
			panic(gerror.NewCodef(
				gcode.CodeInvalidConfiguration,
				`invalid "%s" configuration "%s", missing port`,
				configItemName, address,
			))
		}
		// IPs.
		switch addrArray[0] {
		case "0.0.0.0", "":
			intranetIps, err := gipv4.GetIntranetIpArray()
			if err != nil {
				s.Logger().Errorf(ctx, `error retrieving intranet ip: %+v`, err)
				return nil
			}
			if len(intranetIps) != 0 {
				listenedIps = intranetIps
				break
			}
			// If no intranet ips found, it uses all ips that can be retrieved,
			// it may include internet ip.
			allIps, err := gipv4.GetIpArray()
			if err != nil {
				s.Logger().Errorf(ctx, `error retrieving ip from current node: %+v`, err)
				return nil
			}
			s.Logger().Noticef(
				ctx,
				`no intranet ip found, using internet ip to register service: %v`,
				allIps,
			)
			listenedIps = allIps
		default:
			listenedIps = []string{addrArray[0]}
		}
		// Ports.
		switch addrArray[1] {
		case "0":
			listenedPorts = []int{s.GetListenedPort()}
		default:
			listenedPorts = []int{gconv.Int(addrArray[1])}
		}
		for _, ip := range listenedIps {
			for _, port := range listenedPorts {
				endpoints = append(
					endpoints,
					gsvc.NewEndpoint(fmt.Sprintf(`%s:%d`, ip, port)),
				)
			}
		}
	}
	return endpoints
}
