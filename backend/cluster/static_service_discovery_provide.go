package cluster

import (
	"context"
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
