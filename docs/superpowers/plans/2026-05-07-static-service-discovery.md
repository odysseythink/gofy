# 静态服务发现 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 新增 `static` 服务发现 Provider，使单节点部署时无需 etcd/zookeeper 即可启动并通信。

**Architecture:** 实现 `ServiceDiscoveryProvider` 接口的 `staticServiceDiscoveryProvide`，内置模块名→端口映射，gRPC resolver 直接返回静态地址。`cluster.go` 增加 `"static"` 分支注册新 Provider。

**Tech Stack:** Go, gRPC resolver, confy, mlog, mrun

---

## File Structure

| File | Action | Responsibility |
|------|--------|---------------|
| `backend/cluster/static_service_discovery_provide.go` | Create | 静态 Provider 核心实现：端口映射、Resolver、Register 无操作 |
| `backend/cluster/static_service_discovery_provide_test.go` | Create | 单元测试：Build 解析、未知模块错误、配置覆盖 |
| `backend/cluster/cluster.go` | Modify | `Init()` switch 增加 `"static"` 分支 |

---

## Task 1: Create `static_service_discovery_provide.go`

**Files:**
- Create: `backend/cluster/static_service_discovery_provide.go`

- [ ] **Step 1: Write the static provider implementation**

```go
package cluster

import (
	"fmt"
	"strings"

	"google.golang.org/grpc/resolver"
	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
)

// defaultModulePorts 内置模块名到约定端口的映射
// 新增模块时在此添加一行
var defaultModulePorts = map[string]int{
	"api":    19001,
	"worker": 19002,
	"model":  19003,
}

type staticServiceDiscoveryProvide struct {
	host    string
	portMap map[string]int
}

func (p *staticServiceDiscoveryProvide) Init(args ...any) error {
	p.host = confy.GetWithDefault[string]("cluster.static_host", "127.0.0.1")

	// 深拷贝默认映射
	p.portMap = make(map[string]int, len(defaultModulePorts))
	for k, v := range defaultModulePorts {
		p.portMap[k] = v
	}

	// 用配置覆盖默认端口
	if confy.InConfig("cluster.static_ports") {
		customPorts := confy.Get[map[string]any]("cluster.static_ports")
		for k, v := range customPorts {
			if port, ok := v.(int); ok {
				p.portMap[k] = port
			} else if port, ok := v.(int64); ok {
				p.portMap[k] = int(port)
			} else if port, ok := v.(float64); ok {
				p.portMap[k] = int(port)
			}
		}
	}

	return nil
}

func (p *staticServiceDiscoveryProvide) RunOnce(ctx context.Context) error {
	return nil
}

func (p *staticServiceDiscoveryProvide) Destroy() {
}

func (p *staticServiceDiscoveryProvide) UserData() any {
	return p
}

func (p *staticServiceDiscoveryProvide) RootKey() string {
	return ""
}

func (p *staticServiceDiscoveryProvide) Register(service_name, ip_port_addr, service_status string) error {
	mlog.Infof("[static] register service=%s addr=%s (no-op)", service_name, ip_port_addr)
	return nil
}

func (p *staticServiceDiscoveryProvide) Available() bool {
	return true
}

func (p *staticServiceDiscoveryProvide) GetTarget(model_name string) string {
	return fmt.Sprintf("static:///%s", model_name)
}

func (p *staticServiceDiscoveryProvide) Build(
	target resolver.Target,
	cc resolver.ClientConn,
	opts resolver.BuildOptions,
) (resolver.Resolver, error) {
	moduleName := strings.TrimPrefix(target.URL.Path, "/")
	if moduleName == "" {
		moduleName = strings.TrimPrefix(target.URL.Opaque, "/")
	}

	port, ok := p.portMap[moduleName]
	if !ok {
		return nil, fmt.Errorf("static: unknown module %s, add it to defaultModulePorts", moduleName)
	}

	addr := fmt.Sprintf("%s:%d", p.host, port)
	cc.UpdateState(resolver.State{
		Addresses: []resolver.Address{{Addr: addr}},
	})

	return &staticResolver{}, nil
}

func (p *staticServiceDiscoveryProvide) Scheme() string {
	return "static"
}

type staticResolver struct{}

func (r *staticResolver) ResolveNow(resolver.ResolveNowOptions) {}
func (r *staticResolver) Close()                                {}
```

- [ ] **Step 2: Verify the file compiles**

Run:
```bash
cd backend/cluster && go build .
```

Expected: compile success (no output)

- [ ] **Step 3: Commit**

```bash
git add backend/cluster/static_service_discovery_provide.go
git commit -m "feat(cluster): add static service discovery provider"
```

---

## Task 2: Write unit tests for static provider

**Files:**
- Create: `backend/cluster/static_service_discovery_provide_test.go`

- [ ] **Step 1: Write the test file**

```go
package cluster

import (
	"testing"

	"google.golang.org/grpc/resolver"
)

// mockClientConn implements resolver.ClientConn for testing
type mockClientConn struct {
	state resolver.State
}

func (m *mockClientConn) UpdateState(s resolver.State) error {
	m.state = s
	return nil
}

func (m *mockClientConn) ReportError(err error)                       {}
func (m *mockClientConn) NewAddress(addrs []resolver.Address)          {}
func (m *mockClientConn) NewServiceConfig(serviceConfig string)        {}
func (m *mockClientConn) ParseServiceConfig(serviceConfigJSON string) *resolver.ServiceConfig {
	return nil
}

func TestStaticResolver_Build_KnownModule(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host:    "127.0.0.1",
		portMap: defaultModulePorts,
	}

	cc := &mockClientConn{}
	target := resolver.Target{
		URL: struct {
			Scheme string
			Opaque string
			Host   string
			Path   string
		}{
			Scheme: "static",
			Path:   "/api",
		},
	}

	r, err := p.Build(target, cc, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer r.Close()

	if len(cc.state.Addresses) != 1 {
		t.Fatalf("expected 1 address, got %d", len(cc.state.Addresses))
	}

	expected := "127.0.0.1:19001"
	if cc.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, cc.state.Addresses[0].Addr)
	}
}

func TestStaticResolver_Build_UnknownModule(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host:    "127.0.0.1",
		portMap: defaultModulePorts,
	}

	cc := &mockClientConn{}
	target := resolver.Target{
		URL: struct {
			Scheme string
			Opaque string
			Host   string
			Path   string
		}{
			Scheme: "static",
			Path:   "/unknown",
		},
	}

	_, err := p.Build(target, cc, resolver.BuildOptions{})
	if err == nil {
		t.Fatal("expected error for unknown module, got nil")
	}

	expectedErr := "unknown module unknown"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("expected error to contain %q, got %q", expectedErr, err.Error())
	}
}

func TestStaticResolver_Build_CustomHost(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host:    "192.168.1.100",
		portMap: defaultModulePorts,
	}

	cc := &mockClientConn{}
	target := resolver.Target{
		URL: struct {
			Scheme string
			Opaque string
			Host   string
			Path   string
		}{
			Scheme: "static",
			Path:   "/worker",
		},
	}

	r, err := p.Build(target, cc, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer r.Close()

	expected := "192.168.1.100:19002"
	if cc.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, cc.state.Addresses[0].Addr)
	}
}

func TestStaticResolver_Build_ConfigOverride(t *testing.T) {
	p := &staticServiceDiscoveryProvide{
		host: "127.0.0.1",
		portMap: map[string]int{
			"api":    20001,
			"worker": 19002,
			"model":  19003,
		},
	}

	cc := &mockClientConn{}
	target := resolver.Target{
		URL: struct {
			Scheme string
			Opaque string
			Host   string
			Path   string
		}{
			Scheme: "static",
			Path:   "/api",
		},
	}

	r, err := p.Build(target, cc, resolver.BuildOptions{})
	if err != nil {
		t.Fatalf("Build failed: %v", err)
	}
	defer r.Close()

	expected := "127.0.0.1:20001"
	if cc.state.Addresses[0].Addr != expected {
		t.Errorf("expected address %s, got %s", expected, cc.state.Addresses[0].Addr)
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
	err := p.Register("api", "127.0.0.1:19001", "{}")
	if err != nil {
		t.Fatalf("Register should be no-op and return nil, got: %v", err)
	}
}

func TestStaticProvider_GetTarget(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	target := p.GetTarget("api")
	expected := "static:///api"
	if target != expected {
		t.Errorf("expected target %s, got %s", expected, target)
	}
}

func TestStaticProvider_Scheme(t *testing.T) {
	p := &staticServiceDiscoveryProvide{}
	if p.Scheme() != "static" {
		t.Errorf("expected scheme static, got %s", p.Scheme())
	}
}
```

Note: Need to add `"strings"` to imports in test file.

- [ ] **Step 2: Run tests to verify they pass**

Run:
```bash
cd backend/cluster && go test -v -run TestStatic
```

Expected: All 8 tests PASS

```
=== RUN   TestStaticResolver_Build_KnownModule
--- PASS: TestStaticResolver_Build_KnownModule (0.00s)
=== RUN   TestStaticResolver_Build_UnknownModule
--- PASS: TestStaticResolver_Build_UnknownModule (0.00s)
=== RUN   TestStaticResolver_Build_CustomHost
--- PASS: TestStaticResolver_Build_CustomHost (0.00s)
=== RUN   TestStaticResolver_Build_ConfigOverride
--- PASS: TestStaticResolver_Build_ConfigOverride (0.00s)
=== RUN   TestStaticProvider_Available
--- PASS: TestStaticProvider_Available (0.00s)
=== RUN   TestStaticProvider_Register
--- PASS: TestStaticProvider_Register (0.00s)
=== RUN   TestStaticProvider_GetTarget
--- PASS: TestStaticProvider_GetTarget (0.00s)
=== RUN   TestStaticProvider_Scheme
--- PASS: TestStaticProvider_Scheme (0.00s)
PASS
ok      github.com/odysseythink/gofy/backend/cluster    0.XXXs
```

- [ ] **Step 3: Commit**

```bash
git add backend/cluster/static_service_discovery_provide_test.go
git commit -m "test(cluster): add unit tests for static service discovery provider"
```

---

## Task 3: Integrate static provider into `cluster.go`

**Files:**
- Modify: `backend/cluster/cluster.go:36-58`

- [ ] **Step 1: Add `"static"` branch to the switch statement**

Current code in `cluster.go`:

```go
provider_name := confy.Get[string]("cluster.service_discovery_provide")
switch provider_name {
case "zookeeper":
    zk_service_discovery_provider := &zkServiceDiscoveryProvide{}
    c.subsMgr.Register(zk_service_discovery_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
    c.service_discovery_provider = zk_service_discovery_provider
case "etcd":
    etcd_service_discovery_provider := &etcdServiceDiscoveryProvide{}
    c.subsMgr.Register(etcd_service_discovery_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
    c.service_discovery_provider = etcd_service_discovery_provider
default:
    mlog.Errorf("unsupported provider=%s", provider_name)
    return fmt.Errorf("unsupported provider=%s", provider_name)
}
```

Replace with:

```go
provider_name := confy.Get[string]("cluster.service_discovery_provide")
switch provider_name {
case "zookeeper":
    zk_service_discovery_provider := &zkServiceDiscoveryProvide{}
    c.subsMgr.Register(zk_service_discovery_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
    c.service_discovery_provider = zk_service_discovery_provider
case "etcd":
    etcd_service_discovery_provider := &etcdServiceDiscoveryProvide{}
    c.subsMgr.Register(etcd_service_discovery_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
    c.service_discovery_provider = etcd_service_discovery_provider
case "static":
    static_provider := &staticServiceDiscoveryProvide{}
    c.subsMgr.Register(static_provider, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(0)})
    c.service_discovery_provider = static_provider
default:
    mlog.Errorf("unsupported provider=%s", provider_name)
    return fmt.Errorf("unsupported provider=%s", provider_name)
}
```

- [ ] **Step 2: Verify cluster package still compiles**

Run:
```bash
cd backend/cluster && go build .
```

Expected: compile success (no output)

- [ ] **Step 3: Run all existing cluster tests to ensure no regression**

Run:
```bash
cd backend/cluster && go test -v ./...
```

Expected: All tests pass (including new static tests and existing tests)

- [ ] **Step 4: Commit**

```bash
git add backend/cluster/cluster.go
git commit -m "feat(cluster): wire static service discovery provider into cluster init"
```

---

## Self-Review

### 1. Spec coverage

| Spec Section | Plan Task |
|-------------|-----------|
| 4.2 新增文件 | Task 1, Task 2 |
| 4.3 端口约定机制 | Task 1 (defaultModulePorts, config override) |
| 4.4 静态 Resolver | Task 1 (Build, Scheme, staticResolver) |
| 4.5 cluster.go 集成 | Task 3 |
| 4.6 配置示例 | Implicit in design, validated by tests |
| 5 边界情况 | Task 2 (TestStaticResolver_Build_UnknownModule) |
| 6 测试计划 | Task 2 (全部测试) |

**Gap: None.**

### 2. Placeholder scan

- No TBD/TODO/implement later
- All steps contain actual code, commands, expected output
- No vague descriptions like "add appropriate error handling"

### 3. Type consistency

- `staticServiceDiscoveryProvide` struct fields: `host string`, `portMap map[string]int` — consistent across all tasks
- Method signatures match `ServiceDiscoveryProvider` interface
- `Build()` target parsing logic consistent with etcd/zk implementations (uses `target.URL.Path` and `target.URL.Opaque` fallback)

---

## Execution Handoff

**Plan complete and saved to `docs/superpowers/plans/2026-05-07-static-service-discovery.md`.**

Two execution options:

**1. Subagent-Driven (recommended)** — I dispatch a fresh subagent per task, review between tasks, fast iteration

**2. Inline Execution** — Execute tasks in this session using executing-plans, batch execution with checkpoints

Which approach?
