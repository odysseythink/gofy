package cluster

import "google.golang.org/grpc/resolver"

type ServiceDiscoveryProvider interface {
	resolver.Builder
	RootKey() string
	Register(service_name, ip_port_addr, service_status string) error
	Available() bool
	GetTarget(model_name string) string
}
