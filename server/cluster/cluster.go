package cluster

import (
	"context"

	"google.golang.org/grpc"
	"mlib.com/mrun"
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
