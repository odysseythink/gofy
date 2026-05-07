package cluster

import (
	"context"
	"errors"
	"fmt"

	"github.com/odysseythink/confy"
	"github.com/odysseythink/mlog"
	"github.com/odysseythink/mrun"

	"google.golang.org/grpc"
)

type Cluster struct {
	service_discovery_provider ServiceDiscoveryProvider
	server
	subsMgr mrun.ModuleMgr
}

func (c *Cluster) RunOnce(ctx context.Context) error {

	return nil
}

func (c *Cluster) Destroy() {
	c.subsMgr.Destroy()
}

func (c *Cluster) UserData() any {
	return c
}

func (c *Cluster) GetRpcClientByModule(modulename string) *grpc.ClientConn {
	return c.server.GetRpcClientByModule(modulename)
}
func (c *Cluster) Init(args ...any) error {
	if len(args) < 1 {
		mlog.Errorf("in windows platform, first arg must be rpc service implement instance")
		return errors.New("in windows platform, first arg must be rpc service implement instance")
	}
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
	new_args := append([]any{c.service_discovery_provider}, args...)
	c.subsMgr.Register(&c.server, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(1)}, new_args...)
	return c.subsMgr.Init()
}
