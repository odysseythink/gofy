package cluster

import (
	"fmt"
	"reflect"
	"sync"

	"google.golang.org/grpc"
)

var (
	grpcDescMu sync.RWMutex
	grpcDescs  = map[string]*grpc.ServiceDesc{}
)

// RegisterGrpcServiceDesc publishes a *grpc.ServiceDesc so it can later be
// looked up by its fully-qualified service name (e.g. "pbapi.admin").
// Call this from the pb package's init() for each generated ServiceDesc.
func RegisterGrpcServiceDesc(sd *grpc.ServiceDesc) {
	grpcDescMu.Lock()
	defer grpcDescMu.Unlock()
	grpcDescs[sd.ServiceName] = sd
}

// RegisterServiceByName looks up the ServiceDesc previously published via
// RegisterGrpcServiceDesc and registers impl on the given grpc.Server.
// It replaces the former forked *grpc.Server.RegisterServiceWithoutDesc.
func RegisterServiceByName(s *grpc.Server, impl any, serviceName string) error {
	grpcDescMu.RLock()
	sd, ok := grpcDescs[serviceName]
	grpcDescMu.RUnlock()
	if !ok {
		return fmt.Errorf("cluster: grpc ServiceDesc %q not registered", serviceName)
	}
	if impl == nil {
		return fmt.Errorf("cluster: nil impl for service %q", serviceName)
	}
	ht := reflect.TypeOf(sd.HandlerType).Elem()
	if !reflect.TypeOf(impl).Implements(ht) {
		return fmt.Errorf("cluster: %T does not implement %v (service %q)", impl, ht, serviceName)
	}
	s.RegisterService(sd, impl)
	return nil
}
