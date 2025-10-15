//go:build !windows
// +build !windows

package cluster

import (
	"fmt"

	"mlib.com/confy"
	"mlib.com/mlog"
	"mlib.com/mrun"
)

func (c *Cluster) Init(args ...any) error {
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
	new_args := append([]any{c.service_discovery_provider}, args...)
	c.subsMgr.Register(&c.server, []mrun.ModuleMgrOption{mrun.NewPriorityModuleMgrOption(1)}, new_args...)
	return c.subsMgr.Init()
}
