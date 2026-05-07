package cluster

import (
	"net/url"
	"testing"

	"github.com/odysseythink/confy"
	"google.golang.org/grpc/resolver"
	"google.golang.org/grpc/serviceconfig"
)

// mockClientConn implements resolver.ClientConn for testing.
type mockClientConn struct {
	state resolver.State
}

func (m *mockClientConn) UpdateState(s resolver.State) error {
	m.state = s
	return nil
}

func (m *mockClientConn) ReportError(err error) {}

func (m *mockClientConn) NewAddress(addresses []resolver.Address) {}

func (m *mockClientConn) ParseServiceConfig(serviceConfigJSON string) *serviceconfig.ParseResult {
	return nil
}

func makeTarget(moduleName string) resolver.Target {
	return resolver.Target{
		URL: url.URL{
			Scheme: "static",
			Path:   "/" + moduleName,
		},
	}
}

func TestStaticResolver_Build_KnownModule(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host: "127.0.0.1",
		portMap: map[string]int{
			"api":    19001,
			"worker": 19002,
			"model":  19003,
		},
	}

	mockCC := &mockClientConn{}
	target := makeTarget("api")

	_, err := p.Build(target, mockCC, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	if len(mockCC.state.Addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(mockCC.state.Addresses))
	}

	expected := "127.0.0.1:19001"
	if mockCC.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, mockCC.state.Addresses[0].Addr)
	}
}

func TestStaticResolver_Build_UnknownModule(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host: "127.0.0.1",
		portMap: map[string]int{
			"api":    19001,
			"worker": 19002,
			"model":  19003,
		},
	}

	mockCC := &mockClientConn{}
	target := makeTarget("unknown")

	_, err := p.Build(target, mockCC, resolver.BuildOptions{})
	if err == nil {
		t.Fatal("expected error for unknown module, got nil")
	}

	expectedErr := "static: unknown module unknown, add it to defaultModulePorts"
	if err.Error() != expectedErr {
		t.Errorf("expected error %q, got %q", expectedErr, err.Error())
	}
}

func TestStaticResolver_Build_CustomHost(t *testing.T) {
	confy.Reset()
	defer confy.Reset()

	err := confy.MergeConfigMap(map[string]any{
		"cluster.static_host": "192.168.1.100",
	})
	if err != nil {
		t.Fatalf("failed to merge config: %v", err)
	}

	p := &staticServiceDiscoveryProvide{}
	if err := p.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mockCC := &mockClientConn{}
	target := makeTarget("worker")

	_, err = p.Build(target, mockCC, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expected := "192.168.1.100:19002"
	if mockCC.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, mockCC.state.Addresses[0].Addr)
	}
}

func TestStaticResolver_Build_ConfigOverride(t *testing.T) {
	confy.Reset()
	defer confy.Reset()

	err := confy.MergeConfigMap(map[string]any{
		"cluster.static_ports": map[string]any{
			"api": 29001,
		},
	})
	if err != nil {
		t.Fatalf("failed to merge config: %v", err)
	}

	p := &staticServiceDiscoveryProvide{}
	if err := p.Init(); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	mockCC := &mockClientConn{}
	target := makeTarget("api")

	_, err = p.Build(target, mockCC, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}

	expected := "127.0.0.1:29001"
	if mockCC.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, mockCC.state.Addresses[0].Addr)
	}
}

func TestStaticProvider_Available(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	if !p.Available() {
		t.Error("expected Available() to return true")
	}
}

func TestStaticProvider_Register(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	err := p.Register("test-service", "127.0.0.1:8080", "online")
	if err != nil {
		t.Errorf("expected Register to return nil, got %v", err)
	}
}

func TestStaticProvider_GetTarget(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	got := p.GetTarget("api")
	expected := "static:///api"
	if got != expected {
		t.Errorf("expected GetTarget to return %q, got %q", expected, got)
	}
}

func TestStaticProvider_Scheme(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	got := p.Scheme()
	expected := "static"
	if got != expected {
		t.Errorf("expected Scheme to return %q, got %q", expected, got)
	}
}
